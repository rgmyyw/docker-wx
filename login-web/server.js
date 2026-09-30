/**
 * docker-wx 登录台 v2 — 扫码登录网关(带事件日志 / 心跳保活 / 掉线推送)
 * 页面:  GET /             状态仪表盘 + 事件日志
 * 图片:  GET /qr           当前登录二维码(jpeg)
 * 数据:  GET /status       JSON 状态;GET /logs 最近事件
 * 操作:  POST /newqr /relogin /logout /test-notify
 * 掉线策略:登录后每 5 分钟心跳;连续 3 次失败 → 先自动尝试 62 二次登录,
 *          失败则推送钉钉+邮件并回到扫码等待;恢复后重置。
 */
const http = require('http');
const fs = require('fs');
const crypto = require('crypto');
const nodemailer = require('nodemailer');
const QRCode = require('qrcode');

const API = process.env.wxapi_url || 'http://wxapi:8057';
const PORT = process.env.port || 8058;
const INFO_PATH = process.env.info_path || '/data/login_info.json';
const CREDS_PATH = '/data/sms_creds.json';
const DEVICE_NAME = 'docker-wx-pad';
const PAGE_URL = process.env.page_url || `http://192.168.100.10:${PORT}/`;

const DING = { webhook: process.env.DINGTALK_WEBHOOK || '', secret: process.env.DINGTALK_SECRET || '' };
const SMTP = {
  host: process.env.SMTP_HOST || '', port: Number(process.env.SMTP_PORT || 465),
  user: process.env.SMTP_USER || '', pass: process.env.SMTP_PASS || '', to: process.env.SMTP_TO || process.env.SMTP_USER || '',
};

/* ---------------- 事件日志(内存环形 120 条 + stdout) ---------------- */
const LOGS = [];
const fmtLocal = (d) => `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${d.toLocaleTimeString('zh-CN', { hour12: false })}`;
function log(level, msg, detail) {
  const e = { ts: new Date().toISOString(), local: fmtLocal(new Date()), level, msg, detail };
  LOGS.push(e); if (LOGS.length > 120) LOGS.shift();
  const line = `[${e.local}][${level}] ${msg}${detail ? ` | ${detail}` : ''}`;
  if (level === 'error') console.error(line); else console.log(line);
}
const brief = (o) => { try { const s = JSON.stringify(o); return s && s.length > 500 ? s.slice(0, 500) + '…' : s; } catch { return String(o); } };

/* ---------------- 通知:钉钉(加签) + 邮件(nodemailer) ---------------- */
async function notifyDing(title, text) {
  if (!DING.webhook) return { ok: false, msg: '未配置 DINGTALK_WEBHOOK' };
  const ts = Date.now().toString();
  const sign = encodeURIComponent(crypto.createHmac('sha256', DING.secret).update(`${ts}\n${DING.secret}`).digest('base64'));
  const r = await fetch(`${DING.webhook}&timestamp=${ts}&sign=${sign}`, {
    method: 'POST', headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ msgtype: 'markdown', markdown: { title: title.slice(0, 30), text: `### ${title}\n\n${text}` } }),
  });
  const j = await r.json();
  return j.errcode === 0 ? { ok: true } : { ok: false, msg: `钉钉 errcode=${j.errcode} ${j.errmsg}` };
}
async function notifyMail(title, text) {
  if (!SMTP.host || !SMTP.user) return { ok: false, msg: '未配置 SMTP' };
  const t = nodemailer.createTransport({ host: SMTP.host, port: SMTP.port, secure: SMTP.port === 465, auth: { user: SMTP.user, pass: SMTP.pass } });
  await t.sendMail({ from: SMTP.user, to: SMTP.to, subject: title, text });
  return { ok: true };
}
async function notify(title, text) {
  const [d, m] = await Promise.allSettled([notifyDing(title, text), notifyMail(title, text)]);
  const rd = d.status === 'fulfilled' ? d.value : { ok: false, msg: d.reason?.message };
  const rm = m.status === 'fulfilled' ? m.value : { ok: false, msg: m.reason?.message };
  log(rd.ok ? 'info' : 'error', `钉钉推送${rd.ok ? '成功' : '失败:' + rd.msg}`);
  log(rm.ok ? 'info' : 'error', `邮件推送${rm.ok ? '成功' : '失败:' + rm.msg}`);
  return { dingtalk: rd, mail: rm };
}

/* ---------------- 状态与凭证 ---------------- */
const st = { phase: 'boot', msg: '初始化', uuid: '', qrB64: '', qrTs: 0, expireTs: 0, wxid: '', nick: '', headUrl: '', alias: '', uin: '', mobile: '', deviceId: '', loginTime: '', lastHbOk: '', hbFails: 0, offlineNotified: false, channel: 'Pad', sliderUrl: '' };
/* 取码通道:Pad=8.0.53 正式版;Padx=换版本号绕过 -106 验证;Pad1=云函数;Win/Mac=桌面端(风控策略不同) */
const CHANNELS = {
  Pad: '/api/Login/GetQRPad', Padx: '/api/Login/GetQRPadx', Pad1: '/api/Login/GetQRPad1',
  Win: '/api/Login/GetQRWin', Mac: '/api/Login/GetQRMac',
};
/* 短信验证登录(-106 解法):62dataSMSApply → 收码 → 62dataSMSVerify → 62data 完成 */
const sms = { phase: 'idle', msg: '', checkUrl: '', againUrl: '', cookie: '', data62: '', username: '', password: '', sliderUrl: '', qrPhase: '', qrUrl: '', qrCheck: '' };
/* 扫码验证设备(备用通道):62dataQRCodeApply → 微信扫码确认 → 62dataQRCodeVerify 轮询 */
let info = loadInfo();
function loadInfo() { try { return JSON.parse(fs.readFileSync(INFO_PATH, 'utf8')); } catch { return {}; } }
function saveInfo() { try { fs.writeFileSync(INFO_PATH, JSON.stringify(info, null, 2)); } catch {} }
function genDeviceId() { let s = ''; for (let i = 0; i < 15; i++) s += Math.floor(Math.random() * 10); return s; }
async function api(path, opt) {
  const ctl = new AbortController();
  const timer = setTimeout(() => ctl.abort(new Error('请求超时(60s)')), 60_000);
  try {
    const r = await fetch(`${API}${path}`, { ...opt, signal: ctl.signal });
    return await r.json();
  } finally { clearTimeout(timer); }
}

function deepFindWxid(obj, d = 0) {
  if (!obj || typeof obj !== 'object' || d > 5) return null;
  for (const v of Object.values(obj)) {
    if (typeof v === 'string' && /^wxid_[0-9a-zA-Z_-]{6,}$/.test(v)) return v;
    const r = deepFindWxid(v, d + 1); if (r) return r;
  }
  return null;
}

/* 拉取身份信息(昵称/头像/微信号/Uin/手机号) */
async function fetchProfile() {
  if (!st.wxid) return;
  try {
    const res = await api(`/api/Login/GetCacheInfo?wxid=${encodeURIComponent(st.wxid)}`, { method: 'POST' });
    if (res.Success && res.Data) {
      st.nick = res.Data.NickName || st.nick || '';
      st.headUrl = res.Data.HeadUrl || st.headUrl || '';
      st.alias = res.Data.Alais || st.alias || '';
      st.uin = String(res.Data.Uin || st.uin || '');
      st.mobile = res.Data.Mobile || st.mobile || '';
      info.nick = st.nick; saveInfo();
      log('info', `身份信息已刷新 昵称=${st.nick || '(空)'} 微信号=${st.alias || '(未设)'} Uin=${st.uin}`);
    }
  } catch (e) { log('warn', `身份信息拉取失败: ${e.message}`); }
}

/* ---------------- 取码 ---------------- */
async function newQR(force, reason) {
  if (st.phase === 'qr' && !force && Date.now() < st.expireTs) return;
  const deviceId = info.deviceId || genDeviceId();
  try {
    const qr = await api(CHANNELS[st.channel] || CHANNELS.Pad, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ DeviceID: deviceId, DeviceName: DEVICE_NAME }),
    });
    if (!qr.Success) {
      log('error', `取码失败: ${qr.Message}`, brief(qr));
      Object.assign(st, { phase: 'error', msg: `取码失败: ${qr.Message}` });
      return;
    }
    info.deviceId = deviceId; info.data62 = qr.Data62 || info.data62 || ''; saveInfo();
    Object.assign(st, {
      phase: 'qr', msg: '等待扫码', uuid: qr.Data.Uuid, qrB64: qr.Data.QrBase64,
      qrTs: Date.now(), expireTs: Date.now() + 4.5 * 60 * 1000, deviceId,
      wxid: info.wxid || '', loginTime: info.loginTime || '', hbFails: 0, offlineNotified: false, sliderUrl: '',
    });
    log('info', `取码成功 通道=${st.channel} uuid=${qr.Data.Uuid} 本地过期=${new Date(st.expireTs).toLocaleTimeString('zh-CN', { hour12: false })}${reason ? `(刷新原因:${reason})` : ''}`);
  } catch (e) {
    log('error', `取码异常: wxapi 不可达 ${e.message}`);
    Object.assign(st, { phase: 'error', msg: `wxapi 不可达: ${e.message}` });
  }
}

/* ---------------- 扫码轮询(全量日志,异常如实显示) ---------------- */
let lastPollLog = 0;
async function poll() {
  if (st.phase !== 'qr') return;
  if (Date.now() > st.expireTs) { await newQR(true, '二维码到期'); return; }
  let res;
  try {
    res = await api(`/api/Login/CheckQR?uuid=${encodeURIComponent(st.uuid)}`, { method: 'POST' });
  } catch (e) { log('error', `CheckQR 请求失败: ${e.message}`); return; }

  // 登录成功
  if (res.Message === '登陆成功' || (res.Success && res.Data && res.Data.AcctSectResp)) {
    const d = res.Data || {};
    const wxid = d.AcctSectResp?.UserName || deepFindWxid(d) || info.wxid || '';
    info.wxid = wxid; info.nick = d.AcctSectResp?.NickName || info.nick || '';
    info.loginTime = new Date().toISOString(); saveInfo();
      Object.assign(st, { phase: 'ok', msg: '登录成功', wxid, nick: info.nick, loginTime: info.loginTime, lastHbOk: info.loginTime, hbFails: 0, offlineNotified: false });
      log('info', `登录成功 wxid=${wxid} 昵称=${info.nick || '(未知)'}`);
      fetchProfile();
      return;
  }

  const s = res.Data && (res.Data.status ?? res.Data.Status);
  const now = Date.now();

  // uuid 会话在服务端消失 → 真过期,换码
  if (res.Code === -8 && typeof res.Message === 'string' && res.Message.includes('数据不存在')) {
    log('warn', `CheckQR: uuid 会话不存在(${res.Message}),自动换码`);
    await newQR(true, 'uuid 会话过期');
    return;
  }
  if (res.Code === -3) {
    st.msg = '触发验证码流程,请点「重新取码」';
    log('warn', 'CheckQR: 需要验证码(ticket)', brief(res.Data));
    return;
  }
  if (s === 1) {
    if (st.msg !== '已扫码,请在手机上点确认') log('info', 'CheckQR: 已扫码,等待手机确认');
    st.msg = '已扫码,请在手机上点确认';
    return;
  }
  if (res.Code === 0 && res.Success) { // status 0 = 未扫
    st.msg = '等待扫码';
    if (now - lastPollLog > 60_000) { log('debug', `CheckQR 轮询中(未扫码) uuid=${st.uuid}`); lastPollLog = now; }
    return;
  }
  // -106 分类处理:tcaptcha 滑块→展示给用户;版本过低→提示换通道;不再自动切 Padx(7.x 已被微信封)
  const ret106 = res.Data && res.Data.baseResponse && (res.Data.baseResponse.ret ?? res.Data.baseResponse.Ret);
  if (res.Message === '登陆异常' && ret106 === -106) {
    const em = (res.Data && res.Data.baseResponse && res.Data.baseResponse.errMsg && res.Data.baseResponse.errMsg.string) || '';
    const um = em.match(/<Url><!\[CDATA\[(.*?)\]\]><\/Url>/) || em.match(/<Url>(.*?)<\/Url>/);
    const cm = (em.match(/<Content><!\[CDATA\[(.*?)\]\]><\/Content>/) || em.match(/<Content>(.*?)<\/Content>/) || [])[1] || '';
    if (um && um[1] && /shminorshort|captcha/.test(um[1])) {
      st.msg = '需滑块验证:手机微信扫下方滑块码完成后,点「重新取码」再扫码登录';
      st.sliderUrl = um[1];
      log('warn', '扫码登录触发滑块验证,已展示验证码,完成后请重新取码再扫');
    } else if (cm.includes('版本过低') || cm.includes('升级')) {
      st.msg = `该通道版本已被微信封禁(${cm.slice(0, 20)}…),请用下拉切换通道`;
      log('warn', `扫码 -106 版本过低`, brief(res).slice(0, 200));
    } else {
      st.msg = `登录被拒(-106): ${cm || '环境验证'}`;
      log('warn', '扫码 -106', brief(res).slice(0, 200));
    }
    return;
  }
  // 其它一切异常:如实显示,绝不静默换码
  st.msg = `扫码异常(Code ${res.Code}): ${res.Message}`;
  if (now - lastPollLog > 10_000) { log('error', `CheckQR 异常响应`, brief(res)); lastPollLog = now; }
}

/* ---------------- 登录成功后:心跳保活 + 掉线检测 + 推送 ---------------- */
async function heartbeat() {
  if (st.phase !== 'ok' || !st.wxid) return;
  let res;
  try { res = await api(`/api/Login/HeartBeat?wxid=${encodeURIComponent(st.wxid)}`, { method: 'POST' }); }
  catch (e) { log('error', `心跳请求失败: ${e.message}`); st.hbFails++; return checkOffline(); }
  if (res.Success) {
    if (st.hbFails > 0) log('info', `心跳恢复正常`);
    st.hbFails = 0; st.lastHbOk = new Date().toISOString();
    return;
  }
  st.hbFails++;
  log('warn', `心跳失败(${st.hbFails}/3): ${res.Message}`, brief(res).slice(0, 200));
  checkOffline();
}
async function checkOffline() {
  if (st.hbFails < 3 || st.phase !== 'ok') return;
  log('error', `连续 ${st.hbFails} 次心跳失败,判定掉线 wxid=${st.wxid}`);
  // 先尝试 62 自动恢复
  const r = await relogin(true);
  if (r.ok) {
    await notify('docker-wx 登录已自动恢复', `62 二次登录成功\n\nwxid: ${st.wxid}\n时间: ${new Date().toLocaleString()}`);
    return;
  }
  st.phase = 'offline'; st.msg = '登录已失效,等待重新扫码';
  log('warn', `62 自动恢复失败: ${r.msg},回到扫码等待并推送通知`);
  await newQR(true, '掉线');
  st.phase = 'offline';
  if (!st.offlineNotified) {
    st.offlineNotified = true;
    await notify('docker-wx 微信登录已掉线',
      `wxid: ${st.wxid || info.wxid || '未知'}\n设备ID: ${st.deviceId || info.deviceId}\n最后心跳成功: ${st.lastHbOk || '无'}\n自动恢复: ${r.msg}\n\n请打开 ${PAGE_URL} 重新扫码登录`);
  }
}
setInterval(heartbeat, 5 * 60_000);

/* ---------------- 二次登录 / 退出 ---------------- */
async function relogin(auto) {
  if (!info.data62) return { ok: false, msg: '无 62 凭证(先扫码登录一次)' };
  try {
    const res = await api('/api/Login/62data', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Data62: info.data62, DeviceName: DEVICE_NAME }),
    });
    log(auto ? 'info' : 'info', `62data 响应`, brief(res).slice(0, 300));
    const wxid = res.Data?.AcctSectResp?.UserName || deepFindWxid(res.Data);
    if (res.Success && wxid) {
      info.wxid = wxid; info.loginTime = new Date().toISOString(); saveInfo();
      Object.assign(st, { phase: 'ok', msg: auto ? '已自动恢复(62)' : '二次登录成功(62)', wxid, nick: info.nick, loginTime: info.loginTime, lastHbOk: info.loginTime, hbFails: 0, offlineNotified: false });
      fetchProfile();
      return { ok: true, msg: `二次登录成功 ${wxid}` };
    }
    return { ok: false, msg: res.Message || '未知错误' };
  } catch (e) { return { ok: false, msg: `请求失败: ${e.message}` }; }
}
async function logout() {
  if (!info.wxid) return { ok: false, msg: '未登录' };
  log('info', `退出登录 ${info.wxid}`);
  try { await api(`/api/Login/LogOut?wxid=${encodeURIComponent(info.wxid)}`, { method: 'POST' }); } catch {}
  info.wxid = ''; info.loginTime = ''; saveInfo();
  Object.assign(st, { phase: 'qr', msg: '已退出,等待扫码', wxid: '', loginTime: '', hbFails: 0, offlineNotified: false });
  await newQR(true, '退出登录');
  return { ok: true, msg: '已退出' };
}

/* ---------------- 扫码验证设备(备用通道) ---------------- */
async function qrApply() {
  let c = {}; try { c = JSON.parse(fs.readFileSync(CREDS_PATH, 'utf8')); } catch {}
  const u = sms.username || c.username, p = sms.password || c.password;
  if (!u || !p) return { ok: false, msg: '请先输入账号密码(或已被记住)' };
  log('info', '扫码验证:申请验证二维码');
  try {
    const res = await api('/api/Login/62dataQRCodeApply', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ UserName: u, Password: p, Data62: info.data62 || '', DeviceName: DEVICE_NAME }),
    });
    log('info', '扫码验证:Apply 响应', brief(res).slice(0, 300));
    if (res.Success && res.Data && res.Data.QrUrl) {
      sms.qrUrl = res.Data.QrUrl; sms.qrCheck = res.Data.CheckUrl; sms.qrPhase = 'wait';
      info.data62 = res.Data62 || info.data62; saveInfo();
      sms.msg = '扫码验证:手机微信扫页面上的验证码并确认';
      return { ok: true, msg: '验证二维码已生成,手机微信扫码并确认' };
    }
    return { ok: false, msg: `生成失败: ${res.Message || '未知'}` };
  } catch (e) { return { ok: false, msg: `请求失败: ${e.message}` }; }
}
let qrBusy = false;
async function qrPoll() {
  if (sms.qrPhase !== 'wait' || !sms.qrCheck || qrBusy) return;
  qrBusy = true;
  try {
    const res = await api('/api/Login/62dataQRCodeVerify', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Url: sms.qrCheck }),
    });
    const raw = typeof res.Data === 'string' ? res.Data : JSON.stringify(res.Data);
    if (/window\.code=200|wx_code=/.test(raw)) {
      sms.qrPhase = 'ok'; sms.qrCheck = '';
      sms.msg = '✔ 设备验证通过!请点「申请验证码」收取短信';
      log('info', '扫码验证:设备验证通过');
      notify('docker-wx 设备扫码验证通过', '请回登录台点「申请验证码」收取短信完成登录');
    } else if (/window\.code=201/.test(raw)) {
      sms.msg = '扫码验证:已扫码,请在手机上点确认';
    }
  } catch {} finally { qrBusy = false; }
}
setInterval(qrPoll, 4000);

setInterval(poll, 3000);

/* ---------------- 短信验证登录(-106 解法) ---------------- */
async function smsApply(username, password) {
  if (!username || !password) return { ok: false, msg: '请输入微信账号与密码' };
  sms.username = username; sms.password = password; sms.phase = 'applying'; sms.msg = '正在申请短信验证…';
  try { fs.writeFileSync(CREDS_PATH, JSON.stringify({ username, password }, null, 2)); } catch {}
  log('info', `短信登录:申请验证(账号=${username.slice(0, 3)}***,密码不记录,已记住账号)`);
  try {
    const res = await api('/api/Login/62dataSMSApply', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ UserName: username, Password: password, DeviceName: DEVICE_NAME }),
    });
    log(res.Success ? 'info' : 'error', `短信登录:SMSApply 响应`, brief(res).slice(0, 300));
    if (res.Message === '已申请短信验证' && res.Data && res.Data.CheckUrl) {
      sms.checkUrl = res.Data.CheckUrl; sms.againUrl = res.Data.AgainUrl || ''; sms.cookie = res.Data.Cookie || '';
      sms.data62 = res.Data62 || '';
      sms.phase = 'applied'; sms.msg = '验证码已发送到该微信绑定的手机,请查收短信';
      sms.sliderUrl = '';
      log('info', '短信登录:验证码已申请,等待用户输入');
      return { ok: true, msg: '验证码已发送,请输入收到的短信验证码' };
    }
    // -106 环境检测:errMsg 带 Url = 滑块验证页,交给用户浏览器手动完成
    const errMsgXml = (res.Data && res.Data.baseResponse && res.Data.baseResponse.errMsg && res.Data.baseResponse.errMsg.string) || '';
    const retCode = res.Data && res.Data.baseResponse && res.Data.baseResponse.ret;
    const um = errMsgXml.match(/<Url><!\[CDATA\[(.*?)\]\]><\/Url>/) || errMsgXml.match(/<Url>(.*?)<\/Url>/);
    const contentM = errMsgXml.match(/<Content><!\[CDATA\[(.*?)\]\]><\/Content>/) || errMsgXml.match(/<Content>(.*?)<\/Content>/);
    const sliderU = um && um[1] && /shminorshort\.weixin\.qq\.com|wx\.qq\.com.*tcaptcha|captcha/.test(um[1]) ? um[1] : '';
    if (String(retCode) === '-106' && sliderU) {
      sms.phase = 'slider'; sms.msg = '需滑块安全验证:点下方链接在浏览器完成验证,然后回来重新点「申请验证码」';
      sms.sliderUrl = sliderU;
      log('warn', '短信登录:微信要求滑块验证(环境检测),验证链接已生成,待用户手动完成');
      return { ok: false, msg: '需滑块安全验证:请点页面下方链接完成滑块,再回来重新申请' };
    }
    sms.phase = 'idle'; sms.msg = ''; sms.sliderUrl = '';
    const cmsg = contentM && contentM[1] ? contentM[1] : (res.Message || '未知错误');
    return { ok: false, msg: `申请失败(ret=${retCode}): ${cmsg}` };
  } catch (e) {
    sms.phase = 'idle'; sms.msg = `申请异常: ${e.message}`;
    log('error', `短信登录:申请异常 ${e.message}`);
    return { ok: false, msg: `申请异常: ${e.message}(服务端可能崩溃,查看 docker logs wxapi)` };
  }
}
async function smsAgain() {
  if (!sms.againUrl) return { ok: false, msg: '尚未申请验证码' };
  try {
    const res = await api('/api/Login/62dataSMSAgain', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Url: sms.againUrl, Cookie: sms.cookie }),
    });
    log('info', `短信登录:重发验证码 ${res.Success ? '成功' : '失败 ' + res.Message}`);
    return res.Success ? { ok: true, msg: '已重发' } : { ok: false, msg: res.Message || '重发失败' };
  } catch (e) { return { ok: false, msg: `请求失败: ${e.message}` }; }
}
async function smsVerify(code) {
  if (sms.phase !== 'applied') return { ok: false, msg: '请先申请验证码' };
  if (!code) return { ok: false, msg: '请输入验证码' };
  sms.phase = 'verifying'; sms.msg = '正在提交验证码…';
  log('info', '短信登录:提交验证码');
  try {
    const v = await api('/api/Login/62dataSMSVerify', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Url: sms.checkUrl, Cookie: sms.cookie, Sms: code }),
    });
    log(v.Success ? 'info' : 'error', `短信登录:SMSVerify 响应`, brief(v).slice(0, 200));
    if (!v.Success) { sms.phase = 'applied'; return { ok: false, msg: `验证失败: ${v.Message}(可重试)` }; }
    // 验证通过 → 用同一 62 凭证重新登录完成会话
    log('info', '短信登录:验证通过,正在完成登录');
    const res = await api('/api/Login/62data', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ UserName: sms.username, Password: sms.password, Data62: sms.data62, DeviceName: DEVICE_NAME }),
    });
    log('info', '短信登录:62data 完成响应', brief(res).slice(0, 300));
    const wxid = res.Data?.AcctSectResp?.UserName || deepFindWxid(res.Data);
    if (res.Success && wxid) {
      info.deviceId = info.deviceId || '';
      info.wxid = wxid; info.data62 = sms.data62 || info.data62; info.loginTime = new Date().toISOString(); saveInfo();
      Object.assign(st, { phase: 'ok', msg: '登录成功(短信验证)', wxid, nick: res.Data?.AcctSectResp?.NickName || '', loginTime: info.loginTime, lastHbOk: info.loginTime, hbFails: 0, offlineNotified: false });
      sms.phase = 'idle'; sms.msg = ''; sms.password = '';
      log('info', `短信登录成功 wxid=${wxid}`);
      return { ok: true, msg: `登录成功 ${wxid}` };
    }
    sms.phase = 'applied';
    return { ok: false, msg: `最终登录未成功: ${res.Message}(若提示仍需验证,可重新申请)` };
  } catch (e) { sms.phase = 'applied'; return { ok: false, msg: `请求失败: ${e.message}` }; }
}
(async () => {
  log('info', `login-web 启动 api=${API} 钉钉=${DING.webhook ? '已配置' : '未配置'} SMTP=${SMTP.user ? '已配置' : '未配置'}`);
  if (info.wxid) {
    Object.assign(st, { phase: 'ok', msg: '已登录(历史会话,待心跳确认)', wxid: info.wxid, nick: info.nick || '', loginTime: info.loginTime || '' });
    fetchProfile();
    heartbeat();
  } else await newQR(true);
})();

/* ---------------- HTTP ---------------- */
function statusJson() { return JSON.stringify({ ...st, qrB64: undefined, serverTime: new Date().toISOString(), channels: Object.keys(CHANNELS), smsPhase: sms.phase, smsMsg: sms.msg, sliderUrl: (st.sliderUrl || sms.sliderUrl), qrPhase: sms.qrPhase, qrUrl: sms.qrUrl }); }

const PAGE = `<!doctype html><html lang="zh"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>微信登录台 · docker-wx</title>
<style>
:root{--bg:#0f1115;--card:#171a21;--line:#232833;--tx:#e6e9ef;--sub:#8b93a3;--ok:#3fb96f;--warn:#e0a23c;--err:#e05c5c;--acc:#4f8ef7}
*{box-sizing:border-box;margin:0;padding:0}
body{background:var(--bg);color:var(--tx);font:15px/1.6 -apple-system,"Segoe UI","PingFang SC","Microsoft YaHei",sans-serif;min-height:100vh;display:flex;align-items:center;justify-content:center;padding:24px}
.wrap{width:100%;max-width:980px}
h1{font-size:20px;font-weight:600;display:flex;align-items:center;gap:10px;margin-bottom:16px}
.dot{width:9px;height:9px;border-radius:50%;background:var(--sub)}
.dot.qr{background:var(--warn);box-shadow:0 0 8px var(--warn)}
.dot.ok{background:var(--ok);box-shadow:0 0 8px var(--ok)}
.dot.err,.dot.offline{background:var(--err);box-shadow:0 0 8px var(--err)}
.grid{display:grid;grid-template-columns:340px 1fr;gap:16px}
@media(max-width:760px){.grid{grid-template-columns:1fr}}
.card{background:var(--card);border:1px solid var(--line);border-radius:14px;padding:20px}
.card.wide{grid-column:1/-1;margin-top:16px}
.card h2{font-size:14px;color:var(--sub);font-weight:600;letter-spacing:.05em;margin-bottom:14px;text-transform:uppercase}
.qrbox{display:flex;flex-direction:column;align-items:center;gap:12px}
.qrbox img{width:264px;height:264px;border-radius:10px;background:#fff;padding:10px}
.tip{color:var(--sub);font-size:13px;text-align:center}
.tip b{color:var(--tx)}
.err-tip b{color:var(--err)}
.row{display:flex;justify-content:space-between;padding:10px 0;border-bottom:1px dashed var(--line);font-size:14px}
.row:last-child{border-bottom:0}
.row .k{color:var(--sub)}
.row .v{font-family:ui-monospace,Menlo,Consolas,monospace;word-break:break-all;text-align:right;max-width:60%}
.btns{display:flex;flex-wrap:wrap;gap:10px;margin-top:16px}
button{background:#1f2530;color:var(--tx);border:1px solid var(--line);border-radius:9px;padding:9px 16px;font-size:13px;cursor:pointer;transition:.15s}
button:hover{border-color:var(--acc);color:var(--acc)}
button.primary{background:var(--acc);border-color:var(--acc);color:#fff}
button.primary:hover{opacity:.88;color:#fff}
.foot{color:var(--sub);font-size:12px;text-align:center;margin-top:16px}
.ok-big{color:var(--ok);font-size:16px;font-weight:600;margin-bottom:10px}
.err-big{color:var(--err);font-size:16px;font-weight:600;margin-bottom:10px}
#logs{max-height:260px;overflow-y:auto;font:12px/1.9 ui-monospace,Menlo,Consolas,monospace;color:var(--sub)}
#logs .lv-info{color:var(--tx)}#logs .lv-warn{color:var(--warn)}#logs .lv-error{color:var(--err)}#logs .lv-debug{color:#5b6472}
#logs .t{color:#5b6472;margin-right:8px}
.tabs{display:flex;gap:8px;margin-bottom:14px}
.tab{flex:1;background:#1f2530;color:var(--sub);border:1px solid var(--line);border-radius:9px;padding:8px 10px;font-size:13px;cursor:pointer}
.tab.active{color:var(--tx);border-color:var(--acc)}
.smsform{display:flex;flex-direction:column;gap:10px;padding-top:4px}
.smsform input,.chsel select{background:#1f2530;color:var(--tx);border:1px solid var(--line);border-radius:9px;padding:10px 12px;font-size:14px;width:100%}
.chsel{display:flex;align-items:center;gap:8px;margin-top:12px;font-size:13px;color:var(--sub)}
.chsel select{flex:1;width:auto}
#smsMsg{color:var(--warn);font-size:13px;text-align:center;min-height:20px}
#toast{position:fixed;top:18px;left:50%;transform:translateX(-50%);background:#1f2530;border:1px solid var(--line);padding:10px 18px;border-radius:9px;font-size:13px;display:none;z-index:9}
</style></head><body><div class="wrap">
<h1><span class="dot" id="dot"></span>微信登录台 <span style="color:var(--sub);font-size:13px;font-weight:400">docker-wx · 安卓Pad 8.0.53</span></h1>
<div class="grid">
<div class="card"><h2>登录</h2>
<div class="tabs"><button class="tab active" id="tabBtn-qr" onclick="switchTab('qr')">扫码登录</button><button class="tab" id="tabBtn-sms" onclick="switchTab('sms')">短信登录(-106)</button></div>
<div id="tab-qr">
<div class="qrbox">
<div id="qrArea"><img id="qr" src="/qr"></div>
<div class="tip" id="qrTip">微信扫一扫,过期自动刷新</div>
</div>
<div class="chsel">取码通道 <select id="channel" onchange="chgChannel()"></select></div>
</div>
<div id="tab-sms" style="display:none">
<div class="smsform">
<input id="smsUser" placeholder="微信账号(手机号/QQ号/微信号)" autocomplete="off">
<input id="smsPass" type="password" placeholder="微信密码(仅用于本次登录,不存储)">
<button onclick="smsAct('apply')">申请验证码</button>
<input id="smsCode" placeholder="短信验证码" autocomplete="off">
<button class="primary" onclick="smsAct('verify')">验证并登录</button>
<button onclick="smsAct('again')">重发验证码</button>
<button onclick="smsAct('qrapply')">扫码验证设备</button>
<div id="smsMsg">扫码被 -106 拦截时用此方式:输入账号密码申请验证码,微信会发送短信到绑定手机</div>
<div class="tip" style="margin-top:6px"><span id="credTip">账号密码将记住在服务端(调试用,</span><a style="color:var(--acc);cursor:pointer" onclick="forgetCreds()">清除</a><span id="credTip2">)</span></div>
</div>
</div></div>
<div class="card"><h2>状态</h2><div id="stateBody"></div>
<div class="btns">
<button class="primary" onclick="act('newqr','确认重新取码?当前二维码作废')">重新取码</button>
<button onclick="act('relogin','用已保存的62凭证免扫码登录?')">二次登录</button>
<button onclick="act('logout','确认退出当前账号?')">退出登录</button>
<button onclick="act('test-notify','发送一条测试推送(钉钉+邮件)?')">测试推送</button>
</div></div>
<div class="card wide"><h2>事件日志</h2><div id="logs"></div></div>
</div>
<div class="foot">API: ${API} · 图片直链 <a style="color:var(--acc)" href="/qr" target="_blank">/qr</a> · 状态 <a style="color:var(--acc)" href="/status" target="_blank">/status</a> · 日志 <a style="color:var(--acc)" href="/logs" target="_blank">/logs</a></div>
</div><div id="toast"></div>
<script>
let lastUuid='';let lastLogTs='';let lastChannels='';let lastSmsHtml='';
function switchTab(k){document.getElementById('tab-qr').style.display=k==='qr'?'':'none';document.getElementById('tab-sms').style.display=k==='sms'?'':'none';document.getElementById('tabBtn-qr').className='tab'+(k==='qr'?' active':'');document.getElementById('tabBtn-sms').className='tab'+(k==='sms'?' active':'');}
async function fillCreds(){try{const c=await(await fetch('/sms/creds')).json();if(c.username&&c.password){document.getElementById('smsUser').value=c.username;document.getElementById('smsPass').value=c.password;}}catch(e){}}
fillCreds();
async function forgetCreds(){if(!confirm('清除记住的账号密码?'))return;try{await fetch('/sms/forget',{method:'POST'});document.getElementById('smsUser').value='';document.getElementById('smsPass').value='';toast('已清除');}catch(e){toast('失败');}}
async function chgChannel(){const ch=document.getElementById('channel').value;if(!ch)return;toast('切换通道 '+ch+' …');
 try{const r=await(await fetch('/channel',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({channel:ch})})).json();toast(r.msg);}catch(e){toast('请求失败');}}
async function smsAct(k){
 if(k==='apply'){const u=document.getElementById('smsUser').value.trim(),p=document.getElementById('smsPass').value;if(!u||!p){toast('请输入账号和密码');return;}
  document.getElementById('smsMsg').textContent='申请中…';
  try{const r=await(await fetch('/sms/apply',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({username:u,password:p})})).json();document.getElementById('smsMsg').textContent=r.msg||r.ok?'成功':'失败';toast(r.msg);}catch(e){document.getElementById('smsMsg').textContent='请求失败';}}
 if(k==='verify'){const c=document.getElementById('smsCode').value.trim();if(!c){toast('请输入验证码');return;}
  document.getElementById('smsMsg').textContent='验证中…';
  try{const r=await(await fetch('/sms/verify',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({code:c})})).json();document.getElementById('smsMsg').textContent=r.msg;toast(r.msg);}catch(e){document.getElementById('smsMsg').textContent='请求失败';}}
 if(k==='again'){try{const r=await(await fetch('/sms/again',{method:'POST'})).json();toast(r.msg);}catch(e){toast('请求失败');}}
 if(k==='qrapply'){toast('生成验证码…');
  try{const r=await(await fetch('/sms/qrapply',{method:'POST'})).json();toast(r.msg);}catch(e){toast('请求失败');}}
}
async function tick(){
 try{
  const [s,lg]=await Promise.all([(fetch('/status')).then(r=>r.json()),(fetch('/logs')).then(r=>r.json())]);
  const dot=document.getElementById('dot');
  dot.className='dot '+(s.phase==='ok'?'ok':(s.phase==='error'||s.phase==='offline')?'err':'qr');
  const body=document.getElementById('stateBody');
  const rows=(arr)=>arr.map(([k,v])=>'<div class="row"><span class="k">'+k+'</span><span class="v">'+(v||'—')+'</span></div>').join('');
  if(s.phase==='ok'){
    const avUrl=s.headUrl?(s.headUrl.indexOf('http://')===0?'https://'+s.headUrl.slice(7):s.headUrl):'';
    const av=avUrl?'<img src="'+avUrl.replace(/"/g,'&quot;')+'" style="width:54px;height:54px;border-radius:27px;background:#fff" onerror="this.style.display=&apos;none&apos;">':'';
    const maskMobile=s.mobile?s.mobile.slice(0,3)+'****'+s.mobile.slice(-4):'';
    body.innerHTML='<div style="display:flex;align-items:center;gap:12px;margin-bottom:12px">'+av+'<div><div style="font-size:17px;font-weight:600">'+(s.nick||'微信用户')+'</div><div style="color:var(--sub);font-size:13px">'+(s.alias?'微信号:'+s.alias+' · ':'')+'Uin:'+(s.uin||'—')+'</div></div></div><div class="ok-big">✔ '+s.msg+'</div>'+rows([['wxid',s.wxid],['手机号',maskMobile],['登录时间',s.loginTime?new Date(s.loginTime).toLocaleString():''],['最后心跳',s.lastHbOk?new Date(s.lastHbOk).toLocaleString():'待确认'],['设备ID',s.deviceId]]);
    document.getElementById('qrArea').innerHTML='<div style="width:264px;height:264px;display:flex;align-items:center;justify-content:center;border-radius:10px;background:rgba(63,185,111,.08);color:var(--ok);font-size:15px">已在线,无需扫码</div>';
  }else{
    const bad=(s.phase==='error'||s.phase==='offline');
    body.innerHTML=(bad?'<div class="err-big">✘ '+s.msg+'</div>':'')+rows([['状态',s.msg],['取码通道',s.channel],['wxid',s.wxid||'未登录'],['设备ID',s.deviceId],['二维码到期',s.expireTs?new Date(s.expireTs).toLocaleTimeString():'—'],['心跳失败次数',s.hbFails||0]]);
    document.getElementById('qrTip').innerHTML='<b>'+s.msg+'</b>'+(s.sliderUrl?'<div style="display:flex;flex-direction:column;align-items:center;gap:6px;margin-top:8px"><img src="/slider-qr" style="width:170px;height:170px;border-radius:10px;background:#fff;padding:6px"><span style="font-size:12px">滑块验证:手机微信扫此码并完成,然后点「重新取码」再扫登录码</span></div>':'');
    if(s.uuid&&s.uuid!==lastUuid){lastUuid=s.uuid;document.getElementById('qrArea').innerHTML='<img id="qr" src="/qr?ts='+Date.now()+'">';}
  }
  if(lg.ts!==lastLogTs){lastLogTs=lg.ts;
    document.getElementById('logs').innerHTML=lg.items.slice().reverse().map(e=>'<div class="lv-'+e.level+'"><span class="t">'+(e.local||e.ts.slice(11,19))+'</span>'+e.msg+(e.detail?'<div style="opacity:.55;word-break:break-all">'+e.detail+'</div>':'')+'</div>').join('');}
  const sel=document.getElementById('channel');
  if(s.channels&&s.channels.join()!==lastChannels){lastChannels=s.channels.join();sel.innerHTML=s.channels.map(c=>'<option'+(c===s.channel?' selected':'')+'>'+c+'</option>').join('');}
  else if(sel.value!==s.channel){sel.value=s.channel;}
  const sm=document.getElementById('smsMsg');
  const qrHtml=(s.qrPhase==='wait'&&s.qrUrl)?'<div style="display:flex;flex-direction:column;align-items:center;gap:6px;margin-top:8px"><img src="'+s.qrUrl+'" style="width:200px;height:200px;border-radius:10px;background:#fff;padding:8px"><span style="color:var(--sub);font-size:12px">手机微信扫此码并确认(设备验证)</span></div>':'';
  const sliderHtml=s.sliderUrl?'<a href="'+s.sliderUrl.replace(/"/g,'&quot;')+'" target="_blank" style="color:var(--acc);font-weight:600">👉 点此打开滑块验证页面(完成后回来重新申请)</a><br><div style="display:flex;flex-direction:column;align-items:center;gap:6px;margin-top:8px"><img src="/slider-qr" style="width:180px;height:180px;border-radius:10px;background:#fff;padding:8px"><span style="color:var(--sub);font-size:12px">手机微信扫此码打开验证页(需在微信内完成)</span></div>':'';
  const want=(s.smsMsg||'扫码被 -106 拦截时用此方式')+ (sliderHtml?'<br>'+sliderHtml:'')+(qrHtml?qrHtml:'');
  if(want!==lastSmsHtml){lastSmsHtml=want;sm.innerHTML=want;}
 }catch(e){}
 setTimeout(tick,3000);
}
tick();
function toast(m){const t=document.getElementById('toast');t.textContent=m;t.style.display='block';setTimeout(()=>t.style.display='none',3000);}
async function act(k,confirmMsg){if(confirmMsg&&!confirm(confirmMsg))return;toast('执行中...');
 try{const r=await(await fetch('/'+k,{method:'POST'})).json();toast(typeof r.msg==='string'?r.msg:(r.dingtalk||r.mail)?('钉钉:'+(r.dingtalk.ok?'成功':'失败')+' 邮件:'+(r.mail.ok?'成功':'失败')):'完成');}catch(e){toast('请求失败');}}
</script></body></html>`;

/* ---------------- smallcat 兼容层(QLScriptPublic wxapp 脚本取小程序 code) ----------------
   协议:POST /wx/code {appid, openid} → {code} 或 {status:false, message};POST /wx/refresh → no-op
   内部调 docker-wx /api/Wxapp/JSLogin(当前登录微信);openid 忽略(单账号),多账号时可扩展映射 */
async function wxCodeCompat(body) {
  const appid = body && body.appid;
  if (!appid) return { status: false, message: '缺少 appid' };
  const wxid = st.wxid || info.wxid;
  if (!wxid) return { status: false, message: '微信未登录,请先在登录台扫码' };
  try {
    const res = await api('/api/Wxapp/JSLogin', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Wxid: wxid, Appid: appid }),
    });
    const d = res.Data || {};
    const code = d.Code || d.code || (d.Data && (d.Data.Code || d.Data.code));
    if (res.Success && code) {
      log('info', `小程序取码成功 appid=${appid}`);
      return { code: String(code) };
    }
    log('warn', `小程序取码失败 appid=${appid}: ${res.Message}`);
    return { status: false, message: res.Message || 'JSLogin 失败' };
  } catch (e) {
    log('error', `小程序取码异常 appid=${appid}: ${e.message}`);
    return { status: false, message: e.message };
  }
}

http.createServer(async (req, res) => {
  const url = new URL(req.url, 'http://x');
  const json = (code, obj) => { res.writeHead(code, { 'Content-Type': 'application/json; charset=utf-8', 'Cache-Control': 'no-store' }); res.end(JSON.stringify(obj)); };
  if (req.method === 'GET' && url.pathname === '/') {
    res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' }); res.end(PAGE);
  } else if (req.method === 'GET' && url.pathname === '/qr') {
    if (st.qrB64) {
      res.writeHead(200, { 'Content-Type': 'image/jpeg', 'Cache-Control': 'no-store' });
      res.end(Buffer.from(st.qrB64.split(',')[1], 'base64'));
    } else {
      res.writeHead(503, { 'Content-Type': 'text/plain; charset=utf-8' });
      res.end(st.phase === 'ok' ? '已登录,无需二维码' : '二维码生成中,稍后刷新');
    }
  } else if (req.method === 'GET' && url.pathname === '/status') {
    json(200, JSON.parse(statusJson()));
  } else if (req.method === 'GET' && url.pathname === '/logs') {
    json(200, { ts: LOGS.length ? LOGS[LOGS.length - 1].ts : '', items: LOGS });
  } else if (req.method === 'GET' && url.pathname === '/slider-qr') {
    const su = st.sliderUrl || sms.sliderUrl;
    if (!su) { res.writeHead(404, { 'Content-Type': 'text/plain; charset=utf-8' }); res.end('当前无滑块验证链接'); return; }
    const buf = await QRCode.toBuffer(su, { width: 240, margin: 2 });
    res.writeHead(200, { 'Content-Type': 'image/png', 'Cache-Control': 'max-age=30' }); res.end(buf);
  } else if (req.method === 'POST' && url.pathname === '/newqr') {
    await newQR(true, '手动重新取码'); json(200, { ok: true, msg: '已重新取码' });
  } else if (req.method === 'POST' && url.pathname === '/relogin') {
    json(200, await relogin(false));
  } else if (req.method === 'POST' && url.pathname === '/logout') {
    json(200, await logout());
  } else if (req.method === 'POST' && url.pathname === '/test-notify') {
    log('info', '手动触发测试推送');
    json(200, await notify('docker-wx 登录台测试推送', `通道自检 ${new Date().toLocaleString('zh-CN', { hour12: false })}\n\n收到本条说明钉钉与邮件通道工作正常。`));
  } else if (req.method === 'POST' && url.pathname === '/channel') {
    let b = ''; for await (const c of req) b += c;
    const ch = JSON.parse(b || '{}').channel;
    if (!CHANNELS[ch]) { json(400, { ok: false, msg: `未知通道 ${ch}` }); return; }
    st.channel = ch;
    log('info', `手动切换取码通道 → ${ch}`);
    await newQR(true, `手动切换通道 ${ch}`);
    json(200, { ok: true, msg: `已切换到 ${ch} 通道并取新码` });
  } else if (req.method === 'GET' && url.pathname === '/sms/creds') {
    let c = {}; try { c = JSON.parse(fs.readFileSync(CREDS_PATH, 'utf8')); } catch {}
    json(200, { username: c.username || '', password: c.password || '' });
  } else if (req.method === 'POST' && url.pathname === '/sms/forget') {
    try { fs.unlinkSync(CREDS_PATH); } catch {}
    log('info', '已清除记住的账号密码');
    json(200, { ok: true, msg: '已清除' });
  } else if (req.method === 'POST' && url.pathname === '/wx/code') {
    let b = ''; for await (const c of req) b += c;
    let body = {}; try { body = JSON.parse(b || '{}'); } catch {}
    json(200, await wxCodeCompat(body));
  } else if (req.method === 'POST' && url.pathname === '/wx/refresh') {
    json(200, { status: true });
  } else if (req.method === 'POST' && url.pathname === '/sms/qrapply') {
    json(200, await qrApply());
  } else if (req.method === 'POST' && url.pathname === '/sms/apply') {
    let b = ''; for await (const c of req) b += c;
    const { username, password } = JSON.parse(b || '{}');
    json(200, await smsApply(username, password));
  } else if (req.method === 'POST' && url.pathname === '/sms/again') {
    json(200, await smsAgain());
  } else if (req.method === 'POST' && url.pathname === '/sms/verify') {
    let b = ''; for await (const c of req) b += c;
    const { code } = JSON.parse(b || '{}');
    json(200, await smsVerify(code));
  } else {
    res.writeHead(404, { 'Content-Type': 'text/plain; charset=utf-8' }); res.end('Not Found');
  }
}).listen(PORT, () => log('info', `login-web v2 on :${PORT}`));
