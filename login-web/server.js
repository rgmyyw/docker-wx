/**
 * docker-wx 登录台 v3 — 多账号扫码登录网关
 * 多账号模型:login_info.json v2 {version, accounts:[{alias,wxid,nick,...,deviceId,data62}]}
 * 每号独立 deviceId/62 凭证/心跳/掉线检测;扫码=追加新号;小程序取码按 alias 路由
 * 页面:  GET /  状态仪表盘+事件日志;GET /qr 二维码;GET /status /logs
 * 操作:  POST /newqr /logout {alias} /relogin {alias} /test-notify /channel {ch}
 * 短信:  POST /sms/apply|again|verify /sms/forget;GET /sms/creds
 * 设备验证: POST /sms/qrapply
 * 兼容层: POST /wx/code {appid, openid→alias} /wx/refresh(smallcat 协议)
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

/* ---------------- 事件日志 ---------------- */
const LOGS = [];
const fmtLocal = (d) => `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${d.toLocaleTimeString('zh-CN', { hour12: false })}`;
function log(level, msg, detail) {
  const e = { ts: new Date().toISOString(), local: fmtLocal(new Date()), level, msg, detail };
  LOGS.push(e); if (LOGS.length > 500) LOGS.shift();
  const line = `[${e.local}][${level}] ${msg}${detail ? ` | ${detail}` : ''}`;
  if (level === 'error') console.error(line); else console.log(line);
}
const brief = (o) => { try { const s = JSON.stringify(o); return s && s.length > 2000 ? s.slice(0, 2000) + '…' : s; } catch { return String(o); } };

/* ---------------- 通知 ---------------- */
async function notifyDing(title, text) {
  if (!DING.webhook) return { ok: false, msg: '未配置' };
  const ts = Date.now().toString();
  const sign = encodeURIComponent(crypto.createHmac('sha256', DING.secret).update(`${ts}\n${DING.secret}`).digest('base64'));
  const r = await fetch(`${DING.webhook}&timestamp=${ts}&sign=${sign}`, {
    method: 'POST', headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ msgtype: 'markdown', markdown: { title: title.slice(0, 30), text: `### ${title}\n\n${text}` } }),
  });
  const j = await r.json();
  return j.errcode === 0 ? { ok: true } : { ok: false, msg: `errcode=${j.errcode} ${j.errmsg}` };
}
async function notifyMail(title, text) {
  if (!SMTP.host || !SMTP.user) return { ok: false, msg: '未配置' };
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

/* ---------------- 存储(v2 多账号) ---------------- */
function rawInfo() { try { return JSON.parse(fs.readFileSync(INFO_PATH, 'utf8')); } catch { return null; } }
function writeInfo(o) { try { fs.writeFileSync(INFO_PATH, JSON.stringify(o, null, 2)); } catch (e) { log('error', `写 login_info 失败: ${e.message}`); } }
/* v1(单账号 {deviceId,data62,wxid,nick,loginTime}) → v2 迁移 */
function migrate(old) {
  if (old && Array.isArray(old.accounts)) return old;
  const accs = [];
  if (old && old.wxid) accs.push({
    alias: '1', wxid: old.wxid, nick: old.nick || '', headUrl: '', aliasWx: '', uin: '', mobile: '',
    deviceId: old.deviceId || '', data62: old.data62 || '', loginTime: old.loginTime || '',
  });
  return { version: 2, accounts: accs };
}
let info = migrate(rawInfo());
if (rawInfo() && !Array.isArray(rawInfo().accounts)) writeInfo(info);
function saveInfo() { writeInfo(info); }

function genDeviceId() { let s = ''; for (let i = 0; i < 15; i++) s += Math.floor(Math.random() * 10); return s; }
function nextAlias() { let n = 1; const used = new Set(info.accounts.map(a => a.alias)); while (used.has(String(n))) n++; return String(n); }

/* ---------------- 运行时账号 ---------------- */
const accounts = []; // {alias,wxid,nick,headUrl,aliasWx,uin,mobile,deviceId,data62,loginTime,lastHbOk,hbFails,offlineNotified}
function initAccounts() {
  accounts.length = 0;
  for (const a of info.accounts) accounts.push({ ...a, lastHbOk: '', hbFails: 0, offlineNotified: false });
}
function persistAccount(acc) {
  const i = info.accounts.findIndex(a => a.wxid === acc.wxid);
  const rec = {
    alias: acc.alias, wxid: acc.wxid, nick: acc.nick || '', headUrl: acc.headUrl || '',
    aliasWx: acc.aliasWx || '', uin: acc.uin || '', mobile: acc.mobile || '',
    deviceId: acc.deviceId || '', data62: acc.data62 || '', loginTime: acc.loginTime || '',
  };
  if (i >= 0) info.accounts[i] = rec; else info.accounts.push(rec);
  saveInfo();
}
function removeAccount(wxid) {
  const i = info.accounts.findIndex(a => a.wxid === wxid);
  if (i >= 0) { info.accounts.splice(i, 1); saveInfo(); }
  const j = accounts.findIndex(a => a.wxid === wxid);
  if (j >= 0) accounts.splice(j, 1);
}
function accLabel(a) { return `${a.nick || '微信用户'}(${a.alias})${a.mobile ? ' ' + a.mobile.slice(0, 3) + '****' + a.mobile.slice(-4) : ''}`; }

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

/* ---------------- 身份信息 ---------------- */
async function fetchProfile(acc) {
  if (!acc || !acc.wxid) return;
  try {
    const res = await api(`/api/Login/GetCacheInfo?wxid=${encodeURIComponent(acc.wxid)}`, { method: 'POST' });
    if (res.Success && res.Data) {
      acc.nick = res.Data.NickName || acc.nick || '';
      acc.headUrl = res.Data.HeadUrl || acc.headUrl || '';
      acc.aliasWx = res.Data.Alais || acc.aliasWx || '';
      acc.uin = String(res.Data.Uin || acc.uin || '');
      acc.mobile = res.Data.Mobile || acc.mobile || '';
      persistAccount(acc);
      log('info', `[${acc.alias}] 身份刷新 ${acc.nick || '(空)'} ${acc.aliasWx || ''}`);
    }
  } catch (e) { log('warn', `[${acc.alias}] 身份拉取失败: ${e.message}`); }
}

/* ---------------- 扫码会话(登录新号) ---------------- */
const CHANNELS = {
  Pad: '/api/Login/GetQRPad', Padx: '/api/Login/GetQRPadx', Pad1: '/api/Login/GetQRPad1',
  Win: '/api/Login/GetQRWin', Mac: '/api/Login/GetQRMac',
};
const st = { phase: 'boot', msg: '初始化', uuid: '', qrB64: '', qrTs: 0, expireTs: 0, deviceId: '', data62: '', channel: 'Pad', sliderUrl: '' };

async function newQR(force, reason) {
  if (st.phase === 'qr' && !force && Date.now() < st.expireTs) return;
  const deviceId = genDeviceId(); // 每次扫码会话独立设备,登录后归属该号
  try {
    const qr = await api(CHANNELS[st.channel] || CHANNELS.Pad, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ DeviceID: deviceId, DeviceName: DEVICE_NAME }),
    });
    if (!qr.Success) { log('error', `取码失败: ${qr.Message}`, brief(qr)); Object.assign(st, { phase: 'error', msg: `取码失败: ${qr.Message}` }); return; }
    Object.assign(st, {
      phase: 'qr', msg: '等待扫码(登录新号)', uuid: qr.Data.Uuid, qrB64: qr.Data.QrBase64, qrTs: Date.now(),
      expireTs: Date.now() + 4.5 * 60 * 1000, deviceId, data62: qr.Data62 || '', sliderUrl: '',
    });
    log('info', `取码成功 通道=${st.channel} uuid=${qr.Data.Uuid}${reason ? `(刷新:${reason})` : ''}`);
  } catch (e) { log('error', `取码异常: ${e.message}`); Object.assign(st, { phase: 'error', msg: `wxapi 不可达: ${e.message}` }); }
}

function addAccount(wxid, deviceId, data62) {
  let acc = accounts.find(a => a.wxid === wxid);
  if (acc) { // 同号重登:更新凭证
    acc.deviceId = deviceId || acc.deviceId; acc.data62 = data62 || acc.data62;
    acc.loginTime = new Date().toISOString(); acc.lastHbOk = acc.loginTime; acc.hbFails = 0; acc.offlineNotified = false;
    persistAccount(acc); fetchProfile(acc);
    log('info', `[${acc.alias}] ${accLabel(acc)} 重新登录`);
    return acc;
  }
  acc = { alias: nextAlias(), wxid, nick: '', headUrl: '', aliasWx: '', uin: '', mobile: '', deviceId: deviceId || '', data62: data62 || '', loginTime: new Date().toISOString(), lastHbOk: '', hbFails: 0, offlineNotified: false };
  accounts.push(acc); persistAccount(acc); fetchProfile(acc); heartbeat(acc);
  log('info', `[${acc.alias}] 新账号登录成功 wxid=${wxid},标识=${acc.alias}`);
  return acc;
}

let lastPollLog = 0;
async function poll() {
  if (st.phase !== 'qr') return;
  if (Date.now() > st.expireTs) { await newQR(true, '二维码到期'); return; }
  let res;
  try { res = await api(`/api/Login/CheckQR?uuid=${encodeURIComponent(st.uuid)}`, { method: 'POST' }); }
  catch (e) { log('error', `CheckQR 请求失败: ${e.message}`); return; }

  if (res.Message === '登陆成功' || (res.Success && res.Data && res.Data.AcctSectResp)) {
    const d = res.Data || {};
    const wxid = d.AcctSectResp?.UserName || deepFindWxid(d) || '';
    if (!wxid) { log('error', '登录成功但未解析出 wxid', brief(d).slice(0, 200)); return; }
    const acc = addAccount(wxid, st.deviceId, st.data62 || '');
    Object.assign(st, { phase: 'ok', msg: `登录成功:${accLabel(acc)}(可继续扫码添加下一号)`, wxidLast: wxid, sliderUrl: '' });
    log('info', `[${acc.alias}] 扫码登录成功 ${accLabel(acc)}`);
    setTimeout(() => { if (st.phase === 'ok') newQR(true, '继续添加账号'); }, 4000);
    return;
  }
  const s = res.Data && (res.Data.status ?? res.Data.Status);
  const now = Date.now();
  if (res.Code === -8 && typeof res.Message === 'string' && res.Message.includes('数据不存在')) {
    log('warn', `CheckQR: uuid 会话不存在,自动换码`); await newQR(true, 'uuid 过期'); return;
  }
  if (res.Code === -3) { st.msg = '触发验证码流程,请点「重新取码」'; log('warn', 'CheckQR: 需要验证码(ticket)'); return; }
  const ret106 = res.Data && res.Data.baseResponse && (res.Data.baseResponse.ret ?? res.Data.baseResponse.Ret);
  if (res.Message === '登陆异常' && ret106 === -106) {
    const em = (res.Data && res.Data.baseResponse && res.Data.baseResponse.errMsg && res.Data.baseResponse.errMsg.string) || '';
    const um = em.match(/<Url><!\[CDATA\[(.*?)\]\]><\/Url>/) || em.match(/<Url>(.*?)<\/Url>/);
    const cm = (em.match(/<Content><!\[CDATA\[(.*?)\]\]><\/Content>/) || em.match(/<Content>(.*?)<\/Content>/) || [])[1] || '';
    if (um && um[1] && /shminorshort|captcha/.test(um[1])) {
      st.msg = '需滑块验证:手机微信扫下方滑块码完成后,点「重新取码」再扫码登录';
      st.sliderUrl = um[1];
      log('warn', '扫码登录触发滑块验证,已展示验证码');
    } else if (cm.includes('版本过低') || cm.includes('升级')) {
      st.msg = `该通道版本已被微信封禁,请用下拉切换通道`; log('warn', '扫码 -106 版本过低');
    } else { st.msg = `登录被拒(-106): ${cm || '环境验证'}`; log('warn', '扫码 -106', brief(res).slice(0, 200)); }
    return;
  }
  if (s === 1) { if (st.msg !== '已扫码,请在手机上点确认') log('info', 'CheckQR: 已扫码,等待确认'); st.msg = '已扫码,请在手机上点确认'; return; }
  if (res.Code === 0 && res.Success) {
    st.msg = `等待扫码(登录新号${accounts.length ? `,已有 ${accounts.length} 个号` : ''})`;
    if (now - lastPollLog > 60_000) { log('debug', `CheckQR 轮询中 uuid=${st.uuid}`); lastPollLog = now; }
    return;
  }
  st.msg = `扫码异常(Code ${res.Code}): ${res.Message}`;
  if (now - lastPollLog > 10_000) { log('error', 'CheckQR 异常响应', brief(res)); lastPollLog = now; }
}

/* ---------------- 心跳/掉线(全部账号) ---------------- */
async function heartbeatAll() {
  for (const acc of accounts) { await heartbeat(acc); }
}
async function heartbeat(acc) {
  if (!acc.wxid) return;
  let res;
  try { res = await api(`/api/Login/HeartBeat?wxid=${encodeURIComponent(acc.wxid)}`, { method: 'POST' }); }
  catch (e) { log('error', `[${acc.alias}] 心跳请求失败: ${e.message}`); acc.hbFails++; await checkOffline(acc); return; }
  if (res.Success) { if (acc.hbFails > 0) log('info', `[${acc.alias}] 心跳恢复`); log('debug', `[${acc.alias}] 心跳OK`); acc.hbFails = 0; acc.lastHbOk = new Date().toISOString(); return; }
  acc.hbFails++;
  const definite = /退出|-13/.test(String(res.Message));
  log('warn', `[${acc.alias}] 心跳失败(${acc.hbFails}/2${definite ? ',确定性错误' : ''}): ${res.Message}`);
  if (definite) { acc.hbFails = 2; }
  await checkOffline(acc);
}
async function checkOffline(acc) {
  if (acc.hbFails < 2) return;
  log('error', `[${acc.alias}] 连续心跳失败,判定掉线 ${accLabel(acc)}`);
  const r = await reloginAcc(acc, true);
  if (r.ok) {
    await notify('docker-wx 账号已自动恢复', `${accLabel(acc)}\n62 二次登录成功\n时间: ${fmtLocal(new Date())}`);
    return;
  }
  acc.offline = true;
  if (!acc.offlineNotified) {
    acc.offlineNotified = true;
    await notify('docker-wx 微信账号已掉线', `${accLabel(acc)}\n自动恢复失败: ${r.msg}\n\n请打开 ${PAGE_URL} 重新扫码登录该号`);
  }
}
setInterval(heartbeatAll, 2 * 60_000);

/* ---------------- 二次登录/退出(按账号) ---------------- */
async function reloginAcc(acc, auto) {
  if (!acc.data62) return { ok: false, msg: '无 62 凭证(需扫码)' };
  try {
    const res = await api('/api/Login/62data', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Data62: acc.data62, DeviceName: DEVICE_NAME }),
    });
    const wxid = res.Data?.AcctSectResp?.UserName || deepFindWxid(res.Data);
    if (res.Success && wxid) {
      acc.wxid = wxid; acc.loginTime = new Date().toISOString(); acc.lastHbOk = acc.loginTime; acc.hbFails = 0; acc.offlineNotified = false; acc.offline = false;
      persistAccount(acc); fetchProfile(acc);
      log('info', `[${acc.alias}] ${auto ? '自动' : '手动'}62 登录成功`);
      return { ok: true, msg: `62 登录成功 ${accLabel(acc)}` };
    }
    return { ok: false, msg: res.Message || '未知错误' };
  } catch (e) { return { ok: false, msg: e.message }; }
}
function pickAcc(aliasOrWxid) {
  if (!aliasOrWxid && accounts.length) return accounts[0];
  return accounts.find(a => a.alias === String(aliasOrWxid) || a.wxid === aliasOrWxid) || null;
}

/* ---------------- 短信登录(独立,成功即 addAccount) ---------------- */
const sms = { phase: 'idle', msg: '', checkUrl: '', againUrl: '', cookie: '', data62: '', username: '', password: '', sliderUrl: '', qrPhase: '', qrUrl: '', qrCheck: '' };
async function smsApply(username, password) {
  if (!username || !password) return { ok: false, msg: '请输入微信账号与密码' };
  sms.username = username; sms.password = password; sms.phase = 'applying'; sms.msg = '正在申请短信验证…';
  try { fs.writeFileSync(CREDS_PATH, JSON.stringify({ username, password }, null, 2)); } catch {}
  log('info', `短信登录:申请验证(账号=${username.slice(0, 3)}***,已记住)`);
  try {
    const res = await api('/api/Login/62dataSMSApply', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ UserName: username, Password: password, DeviceName: DEVICE_NAME }),
    });
    log(res.Success ? 'info' : 'error', '短信登录:SMSApply 响应', brief(res).slice(0, 300));
    if (res.Message === '已申请短信验证' && res.Data && res.Data.CheckUrl) {
      Object.assign(sms, { checkUrl: res.Data.CheckUrl, againUrl: res.Data.AgainUrl || '', cookie: res.Data.Cookie || '', data62: res.Data62 || '', phase: 'applied', msg: '验证码已发送到绑定手机', sliderUrl: '' });
      return { ok: true, msg: '验证码已发送,请输入收到的短信验证码' };
    }
    const errMsgXml = (res.Data && res.Data.baseResponse && res.Data.baseResponse.errMsg && res.Data.baseResponse.errMsg.string) || '';
    const retCode = res.Data && res.Data.baseResponse && res.Data.baseResponse.ret;
    const um = errMsgXml.match(/<Url><!\[CDATA\[(.*?)\]\]><\/Url>/) || errMsgXml.match(/<Url>(.*?)<\/Url>/);
    const contentM = errMsgXml.match(/<Content><!\[CDATA\[(.*?)\]\]><\/Content>/) || errMsgXml.match(/<Content>(.*?)<\/Content>/);
    const sliderU = um && um[1] && /shminorshort|captcha/.test(um[1]) ? um[1] : '';
    if (String(retCode) === '-106' && sliderU) {
      sms.phase = 'slider'; sms.msg = '需滑块安全验证:点下方链接完成验证后重新申请'; sms.sliderUrl = sliderU;
      log('warn', '短信登录:微信要求滑块验证');
      return { ok: false, msg: '需滑块安全验证:请点页面下方链接完成滑块,再回来重新申请' };
    }
    sms.phase = 'idle'; sms.msg = ''; sms.sliderUrl = '';
    const cmsg = contentM && contentM[1] ? contentM[1] : (res.Message || '未知错误');
    return { ok: false, msg: `申请失败(ret=${retCode}): ${cmsg}` };
  } catch (e) {
    sms.phase = 'idle'; sms.msg = `申请异常: ${e.message}`;
    log('error', `短信登录:申请异常 ${e.message}`);
    return { ok: false, msg: `申请异常: ${e.message}` };
  }
}
async function smsAgain() {
  if (!sms.againUrl) return { ok: false, msg: '尚未申请验证码' };
  try {
    const res = await api('/api/Login/62dataSMSAgain', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Url: sms.againUrl, Cookie: sms.cookie }),
    });
    return res.Success ? { ok: true, msg: '已重发' } : { ok: false, msg: res.Message || '重发失败' };
  } catch (e) { return { ok: false, msg: e.message }; }
}
async function smsVerify(code) {
  if (sms.phase !== 'applied') return { ok: false, msg: '请先申请验证码' };
  if (!code) return { ok: false, msg: '请输入验证码' };
  sms.phase = 'verifying'; sms.msg = '正在提交验证码…';
  try {
    const v = await api('/api/Login/62dataSMSVerify', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Url: sms.checkUrl, Cookie: sms.cookie, Sms: code }),
    });
    if (!v.Success) { sms.phase = 'applied'; return { ok: false, msg: `验证失败: ${v.Message}` }; }
    const res = await api('/api/Login/62data', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ UserName: sms.username, Password: sms.password, Data62: sms.data62, DeviceName: DEVICE_NAME }),
    });
    const wxid = res.Data?.AcctSectResp?.UserName || deepFindWxid(res.Data);
    if (res.Success && wxid) {
      const acc = addAccount(wxid, '', sms.data62 || '');
      Object.assign(sms, { phase: 'idle', msg: '', password: '' });
      Object.assign(st, { phase: 'ok', msg: `短信登录成功:${accLabel(acc)}` });
      return { ok: true, msg: `登录成功 ${accLabel(acc)}` };
    }
    sms.phase = 'applied';
    return { ok: false, msg: `最终登录未成功: ${res.Message}` };
  } catch (e) { sms.phase = 'applied'; return { ok: false, msg: e.message }; }
}

/* ---------------- 扫码验证设备(备用) ---------------- */
async function qrApply() {
  let c = {}; try { c = JSON.parse(fs.readFileSync(CREDS_PATH, 'utf8')); } catch {}
  const u = sms.username || c.username, p = sms.password || c.password;
  if (!u || !p) return { ok: false, msg: '请先输入账号密码' };
  try {
    const res = await api('/api/Login/62dataQRCodeApply', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ UserName: u, Password: p, Data62: accounts[0]?.data62 || '', DeviceName: DEVICE_NAME }),
    });
    if (res.Success && res.Data && res.Data.QrUrl) {
      Object.assign(sms, { qrUrl: res.Data.QrUrl, qrCheck: res.Data.CheckUrl, qrPhase: 'wait', msg: '扫码验证:手机微信扫页面上的验证码并确认' });
      return { ok: true, msg: '验证二维码已生成' };
    }
    return { ok: false, msg: `生成失败: ${res.Message || '未知'}` };
  } catch (e) { return { ok: false, msg: e.message }; }
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
      sms.qrPhase = 'ok'; sms.qrCheck = ''; sms.msg = '✔ 设备验证通过!请点「申请验证码」';
      notify('docker-wx 设备扫码验证通过', '请回登录台点「申请验证码」完成登录');
    }
  } catch {} finally { qrBusy = false; }
}
setInterval(qrPoll, 4000);

/* ---------------- smallcat 兼容层(多账号路由) ----------------
   POST /wx/code {appid, openid} → openid 按 alias/wxid 匹配账号(空=第一个) */
async function wxCodeCompat(body) {
  const appid = body && body.appid;
  if (!appid) return { status: false, message: '缺少 appid' };
  const key = body.openid ? String(body.openid).split('#')[0].trim() : '';
  let acc;
  if (key) {
    acc = accounts.find(a => a.alias === key || a.wxid === key);
    if (!acc) return { status: false, message: `标识 ${key} 未登录(现有标识:${accounts.map(a => a.alias).join(',') || '无'})` };
  } else acc = accounts[0];
  if (!acc) return { status: false, message: '无已登录微信,请先在登录台扫码' };
  try {
    const res = await api('/api/Wxapp/JSLogin', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Wxid: acc.wxid, Appid: appid }),
    });
    const d = res.Data || {};
    const code = d.Code || d.code || (d.Data && (d.Data.Code || d.Data.code));
    if (res.Success && code) {
      log('info', `[${acc.alias}] 小程序取码成功 appid=${appid}`);
      return { code: String(code) };
    }
    log('warn', `[${acc.alias}] 小程序取码失败 appid=${appid}: ${res.Message}`);
    return { status: false, message: res.Message || 'JSLogin 失败' };
  } catch (e) { return { status: false, message: e.message }; }
}

/* ---------------- 启动 ---------------- */
setInterval(poll, 3000);
(async () => {
  log('info', `login-web v3(多账号) 启动 api=${API} 账号数=${info.accounts.length} 钉钉=${DING.webhook ? '有' : '无'} SMTP=${SMTP.user ? '有' : '无'}`);
  initAccounts();
  if (accounts.length) { for (const a of accounts) { fetchProfile(a); heartbeat(a); } }
  await newQR(true);
})();

/* ---------------- HTTP ---------------- */
function accBrief(a) {
  return { alias: a.alias, wxid: a.wxid, nick: a.nick, headUrl: (a.headUrl || '').replace(/^http:\/\//, 'https://'), aliasWx: a.aliasWx, uin: a.uin, mobile: a.mobile, loginTime: a.loginTime, lastHbOk: a.lastHbOk, hbFails: a.hbFails, offline: !!a.offline };
}
function statusJson() {
  return JSON.stringify({
    phase: st.phase, msg: st.msg, uuid: st.uuid, qrTs: st.qrTs, expireTs: st.expireTs, channel: st.channel,
    sliderUrl: (st.sliderUrl || sms.sliderUrl), serverTime: new Date().toISOString(),
    channels: Object.keys(CHANNELS), accounts: accounts.map(accBrief),
    smsPhase: sms.phase, smsMsg: sms.msg, qrPhase: sms.qrPhase, qrUrl: sms.qrUrl,
  });
}

const PAGE = `<!doctype html><html lang="zh"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>微信账号管理 · docker-wx</title>
<style>
:root{--bg:#0f1115;--card:#171a21;--line:#232833;--tx:#e6e9ef;--sub:#8b93a3;--ok:#3fb96f;--warn:#e0a23c;--err:#e05c5c;--acc:#4f8ef7}
*{box-sizing:border-box;margin:0;padding:0}
body{background:var(--bg);color:var(--tx);font:15px/1.6 -apple-system,"Segoe UI","PingFang SC","Microsoft YaHei",sans-serif;min-height:100vh;padding:28px 24px}
.wrap{max-width:1060px;margin:0 auto}
.hd{display:flex;align-items:center;gap:12px;margin-bottom:20px;flex-wrap:wrap}
h1{font-size:21px;font-weight:600;display:flex;align-items:center;gap:10px}
.dot{width:9px;height:9px;border-radius:50%;background:var(--sub)}
.dot.ok{background:var(--ok);box-shadow:0 0 8px var(--ok)}
.dot.err{background:var(--err)}
.hd .sp{flex:1}
.hd .summary{color:var(--sub);font-size:13px}
button{background:#1f2530;color:var(--tx);border:1px solid var(--line);border-radius:9px;padding:9px 16px;font-size:13px;cursor:pointer;transition:.15s}
button:hover{border-color:var(--acc);color:var(--acc)}
button.primary{background:var(--acc);border-color:var(--acc);color:#fff}
button.primary:hover{opacity:.88;color:#fff}
button.small{padding:5px 10px;font-size:12px}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(300px,1fr));gap:14px}
.acc{background:var(--card);border:1px solid var(--line);border-radius:14px;padding:16px;display:flex;gap:12px;align-items:center}
.acc img{width:46px;height:46px;border-radius:23px;background:#fff;flex:none}
.acc .info{flex:1;min-width:0}
.acc .nick{font-weight:600;font-size:15px;display:flex;align-items:center;gap:6px;flex-wrap:wrap}
.acc .meta{color:var(--sub);font-size:12px;word-break:break-all;margin-top:2px}
.acc .ops{display:flex;flex-direction:column;gap:6px}
.badge{display:inline-block;font-size:11px;padding:1px 8px;border-radius:8px}
.badge.on{background:rgba(63,185,111,.15);color:var(--ok)}
.badge.off{background:rgba(224,92,92,.15);color:var(--err)}
.alias{display:inline-block;font-size:11px;padding:1px 7px;border-radius:8px;background:rgba(79,142,247,.15);color:var(--acc)}
.empty{color:var(--sub);text-align:center;padding:60px 0;font-size:14px}
.foot{color:var(--sub);font-size:12px;text-align:center;margin-top:22px}
.foot a{color:var(--acc);text-decoration:none}
/* 弹窗 */
#modal{position:fixed;inset:0;background:rgba(0,0,0,.6);display:none;align-items:center;justify-content:center;z-index:50;padding:16px}
.mbox{background:var(--card);border:1px solid var(--line);border-radius:16px;width:100%;max-width:430px;max-height:92vh;overflow-y:auto;padding:20px}
.mhd{display:flex;align-items:center;margin-bottom:14px}
.mhd b{font-size:16px;flex:1}
.tabs{display:flex;gap:8px;margin-bottom:14px}
.tab{flex:1;background:#1f2530;color:var(--sub);border:1px solid var(--line);border-radius:9px;padding:8px;font-size:13px;cursor:pointer}
.tab.active{color:var(--tx);border-color:var(--acc)}
.qrbox{display:flex;flex-direction:column;align-items:center;gap:10px}
.qrbox img{width:238px;height:238px;border-radius:12px;background:#fff;padding:10px}
.tip{color:var(--sub);font-size:13px;text-align:center}
.tip b{color:var(--tx)}
.chsel{display:flex;align-items:center;gap:8px;margin-top:10px;font-size:13px;color:var(--sub);width:100%;justify-content:center}
.chsel select{background:#1f2530;color:var(--tx);border:1px solid var(--line);border-radius:8px;padding:6px 8px}
.smsform{display:flex;flex-direction:column;gap:10px}
.smsform input{background:#1f2530;color:var(--tx);border:1px solid var(--line);border-radius:9px;padding:10px 12px;font-size:14px;width:100%}
#smsMsg{color:var(--warn);font-size:13px;text-align:center;min-height:18px}
.sess{margin-top:14px;border-top:1px dashed var(--line);padding-top:10px}
.sess h4{font-size:12px;color:var(--sub);margin-bottom:6px;font-weight:600}
#sessLogs{max-height:150px;overflow-y:auto;font:11.5px/1.8 ui-monospace,Menlo,Consolas,monospace;color:var(--sub);background:#12141a;border-radius:8px;padding:8px}
#sessLogs .lv-info{color:var(--tx)}#sessLogs .lv-warn{color:var(--warn)}#sessLogs .lv-error{color:var(--err)}#sessLogs .lv-debug{color:#5b6472}
#sessLogs .t{color:#5b6472;margin-right:6px}
.success{background:rgba(63,185,111,.12);border:1px solid rgba(63,185,111,.4);color:var(--ok);border-radius:10px;padding:12px;text-align:center;font-weight:600;margin-bottom:10px;display:none}
#toast{position:fixed;top:18px;left:50%;transform:translateX(-50%);background:#1f2530;border:1px solid var(--line);padding:10px 18px;border-radius:9px;font-size:13px;display:none;z-index:99}
</style></head><body><div class="wrap">
<div class="hd">
<h1><span class="dot" id="dot"></span>微信账号管理 <span style="color:var(--sub);font-size:13px;font-weight:400">docker-wx</span></h1>
<span class="summary" id="summary"></span>
<span class="sp"></span>
<button class="primary" onclick="openModal()">＋ 添加账号</button>
<button onclick="location.href=&apos;/logspage&apos;">全量日志</button>
</div>
<div class="grid" id="accGrid"></div>
<div class="empty" id="empty" style="display:none">还没有账号,点右上「添加账号」扫码登录第一个</div>
<div class="foot">小程序脚本变量填账号标识(如 <b>1</b> 或 <b>1&amp;2</b>) · <a href="/status">status</a> · <a href="/logspage">日志</a></div>
</div>

<div id="modal"><div class="mbox">
<div class="mhd"><b>添加账号</b><button class="small" onclick="closeModal()">关闭</button></div>
<div class="success" id="okBar"></div>
<div class="tabs"><button class="tab active" id="tb-qr" onclick="mTab(&apos;qr&apos;)">扫码登录</button><button class="tab" id="tb-sms" onclick="mTab(&apos;sms&apos;)">短信登录</button></div>
<div id="m-qr">
<div class="qrbox">
<div id="qrArea"><img id="qr" src="/qr"></div>
<div class="tip" id="qrTip">微信扫一扫</div>
<div class="chsel">通道 <select id="channel" onchange="chgChannel()"></select></div>
</div>
</div>
<div id="m-sms" style="display:none">
<div class="smsform">
<input id="smsUser" placeholder="微信账号(手机号/QQ号/微信号)" autocomplete="off">
<input id="smsPass" type="password" placeholder="微信密码(仅登录用,记住于服务端)" autocomplete="off">
<button onclick="smsAct(&apos;apply&apos;)">申请验证码</button>
<input id="smsCode" placeholder="短信验证码" autocomplete="off">
<button class="primary" onclick="smsAct(&apos;verify&apos;)">验证并登录</button>
<button onclick="smsAct(&apos;again&apos;)">重发验证码</button>
<div id="smsMsg">被 -106 拦截时用此方式</div>
</div>
</div>
<div class="sess"><h4>本次会话日志</h4><div id="sessLogs"><div style="color:#5b6472">等待操作…</div></div></div>
</div></div>
<div id="toast"></div>

<script>
let lastUuid='',lastChannels='',lastSmsHtml='',lastAccHtml='',lastTipHtml='',sessStart='',modalOpen=false,doneShown=false;
function openModal(){modalOpen=true;doneShown=false;sessStart=new Date().toISOString();document.getElementById('modal').style.display='flex';}
function closeModal(){modalOpen=false;document.getElementById('modal').style.display='none';}
function mTab(k){document.getElementById('m-qr').style.display=k==='qr'?'':'none';document.getElementById('m-sms').style.display=k==='sms'?'':'none';document.getElementById('tb-qr').className='tab'+(k==='qr'?' active':'');document.getElementById('tb-sms').className='tab'+(k==='sms'?' active':'');}
async function fillCreds(){try{const c=await(await fetch('/sms/creds')).json();if(c.username&&c.password){document.getElementById('smsUser').value=c.username;document.getElementById('smsPass').value=c.password;}}catch(e){}}
async function chgChannel(){const ch=document.getElementById('channel').value;if(!ch)return;toast('切换 '+ch+' …');try{const r=await(await fetch('/channel',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({channel:ch})})).json();toast(r.msg);}catch(e){toast('失败');}}
async function smsAct(k){
 if(k==='apply'){const u=document.getElementById('smsUser').value.trim(),p=document.getElementById('smsPass').value;if(!u||!p){toast('请输入账号和密码');return;}
  document.getElementById('smsMsg').textContent='申请中…';
  try{const r=await(await fetch('/sms/apply',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({username:u,password:p})})).json();document.getElementById('smsMsg').textContent=r.msg;toast(r.msg);}catch(e){document.getElementById('smsMsg').textContent='请求失败';}}
 if(k==='verify'){const c=document.getElementById('smsCode').value.trim();if(!c){toast('请输入验证码');return;}
  document.getElementById('smsMsg').textContent='验证中…';
  try{const r=await(await fetch('/sms/verify',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({code:c})})).json();document.getElementById('smsMsg').textContent=r.msg;toast(r.msg);}catch(e){document.getElementById('smsMsg').textContent='请求失败';}}
 if(k==='again'){try{const r=await(await fetch('/sms/again',{method:'POST'})).json();toast(r.msg);}catch(e){toast('失败');}}
}
async function accAct(k,alias,confirmMsg){if(confirmMsg&&!confirm(confirmMsg))return;toast('执行中…');
 try{const r=await(await fetch('/'+k,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({alias})})).json();toast(r.msg);}catch(e){toast('请求失败');}}
function toast(m){const t=document.getElementById('toast');t.textContent=m;t.style.display='block';setTimeout(()=>{t.style.display='none'},2800);}
async function tick(){
 try{
  const [s,lg]=await Promise.all([(fetch('/status')).then(r=>r.json()),(fetch('/logs')).then(r=>r.json())]);
  document.getElementById('dot').className='dot '+(s.accounts.length?'ok':(s.phase==='error'?'err':'ok'));
  const on=s.accounts.filter(a=>!a.offline).length;
  document.getElementById('summary').textContent=s.accounts.length?('共 '+s.accounts.length+' 个账号,'+on+' 个在线'):'';
  document.getElementById('empty').style.display=s.accounts.length?'none':'';
  const accHtml=s.accounts.map(a=>{
   const av=a.headUrl?'<img src="'+a.headUrl.replace(/"/g,'')+'" onerror="this.style.visibility=&apos;hidden&apos;">':'<img style="visibility:hidden">';
   const badge=a.offline?'<span class="badge off">掉线</span>':'<span class="badge on">在线</span>';
   return '<div class="acc">'+av+'<div class="info"><div class="nick">'+(a.nick||'微信用户')+badge+'<span class="alias">标识 '+a.alias+'</span></div><div class="meta">'+a.wxid+(a.aliasWx?' · '+a.aliasWx:'')+(a.mobile?' · '+a.mobile.slice(0,3)+'****'+a.mobile.slice(-4):'')+'</div></div><div class="ops">'+(a.offline?'<button class="small" data-act="hint">已掉线·扫码恢复</button>':'')+'<button class="small" data-act="logout" data-alias="'+a.alias+'">退出</button></div></div>';
  }).join('');
  if(accHtml!==lastAccHtml){lastAccHtml=accHtml;document.getElementById('accGrid').innerHTML=accHtml;}
  if(modalOpen){
   if(s.uuid&&s.uuid!==lastUuid){lastUuid=s.uuid;document.getElementById('qrArea').innerHTML='<img id="qr" src="/qr?ts='+Date.now()+'">';}
   const sl=s.sliderUrl?'<br><img src="/slider-qr" style="width:150px;height:150px;border-radius:10px;background:#fff;padding:6px;margin-top:6px"><br><span style="font-size:12px">滑块:手机微信扫此码完成后点「重新取码」</span>':'';
   const wt='<b>'+s.msg+'</b>'+sl;
   if(wt!==lastTipHtml){lastTipHtml=wt;document.getElementById('qrTip').innerHTML=wt;}
   const sel=document.getElementById('channel');
   if(s.channels&&s.channels.join()!==lastChannels){lastChannels=s.channels.join();sel.innerHTML=s.channels.map(c=>'<option'+(c===s.channel?' selected':'')+'>'+c+'</option>').join('');}
   const wantSms=(s.smsMsg||' ')+(s.sliderUrl?'<br><a href=&quot;'+s.sliderUrl+'&quot; target="_blank" style="color:var(--acc)">打开滑块验证页</a>':'');
   if(wantSms!==lastSmsHtml){lastSmsHtml=wantSms;document.getElementById('smsMsg').innerHTML=wantSms;}
   const sess=(lg.items||[]).filter(e=>e.ts>=sessStart);
   const sh=sess.slice(-40).map(e=>'<div class="lv-'+e.level+'"><span class="t">'+(e.local||'').slice(11)+'</span>'+e.msg+'</div>').join('');
   const sl2=document.getElementById('sessLogs');
   if(sh&&sl2.innerHTML!==sh){sl2.innerHTML=sh;sl2.scrollTop=sl2.scrollHeight;}
   if(s.phase==='ok'&&!doneShown){doneShown=true;const bar=document.getElementById('okBar');bar.style.display='block';bar.textContent='✅ '+s.msg;setTimeout(()=>{bar.style.display='none';closeModal();},2500);}
  }
 }catch(e){}
 setTimeout(tick,3000);
}
document.getElementById('accGrid').addEventListener('click', function(e){
  var b = e.target.closest('button[data-act]');
  if (!b) return;
  var act = b.getAttribute('data-act'), alias = b.getAttribute('data-alias');
  if (act === 'hint') { toast('请扫码恢复该号'); openModal(); }
  if (act === 'relogin') { accAct('relogin', alias, ''); }
  if (act === 'logout') { if (confirm('退出标识 ' + alias + '?该号需重新扫码')) accAct('logout', alias, ''); }
});
fillCreds();tick();
</script></body></html>`;

/* ---------------- 全量日志独立页 ---------------- */
const LOG_PAGE = `<!doctype html><html lang="zh"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>全量日志 · docker-wx</title>
<style>
:root{--bg:#0f1115;--card:#171a21;--line:#232833;--tx:#e6e9ef;--sub:#8b93a3;--ok:#3fb96f;--warn:#e0a23c;--err:#e05c5c;--acc:#4f8ef7}
*{box-sizing:border-box;margin:0;padding:0}
body{background:var(--bg);color:var(--tx);font:14px/1.6 -apple-system,"Segoe UI","PingFang SC",sans-serif;padding:24px}
.wrap{max-width:1000px;margin:0 auto}
.hd{display:flex;align-items:center;gap:10px;margin-bottom:14px;flex-wrap:wrap}
h1{font-size:18px;font-weight:600;flex:1}
button{background:#1f2530;color:var(--tx);border:1px solid var(--line);border-radius:8px;padding:6px 12px;font-size:12px;cursor:pointer}
button.on{border-color:var(--acc);color:var(--acc)}
#logs{background:var(--card);border:1px solid var(--line);border-radius:12px;padding:14px;max-height:calc(100vh - 130px);overflow-y:auto;font:12px/1.9 ui-monospace,Menlo,Consolas,monospace;color:var(--sub)}
.lv-info{color:var(--tx)}.lv-warn{color:var(--warn)}.lv-error{color:var(--err)}.lv-debug{color:#5b6472}
.t{color:#5b6472;margin-right:8px}
.d{opacity:.55;word-break:break-all;margin-left:8px}
.meta{color:var(--sub);font-size:12px;margin-bottom:8px}
</style></head><body><div class="wrap">
<div class="hd">
<h1>事件日志(全部)</h1>
<button data-lv="" class="on" onclick="setLv(this)">全部</button>
<button data-lv="info" onclick="setLv(this)">info</button>
<button data-lv="warn" onclick="setLv(this)">warn</button>
<button data-lv="error" onclick="setLv(this)">error</button>
<button id="pauseBtn" onclick="togglePause()">暂停滚动</button>
<button onclick="location.href=&apos;/&apos;">返回账号管理</button>
</div>
<div class="meta" id="meta"></div>
<div id="logs"></div>
</div>
<script>
let lv='',paused=false,lastHtml='';
function setLv(b){lv=b.getAttribute('data-lv');document.querySelectorAll('.hd button[data-lv]').forEach(x=>x.className='');b.className='on';lastHtml='';}
function togglePause(){paused=!paused;document.getElementById('pauseBtn').className=paused?'on':'';document.getElementById('pauseBtn').textContent=paused?'恢复滚动':'暂停滚动';}
async function tick(){
 try{
  const lg=await(await fetch('/logs')).json();
  const items=(lg.items||[]).filter(e=>!lv||e.level===lv);
  document.getElementById('meta').textContent='最近 '+items.length+' 条(上限 500,内存滚动)';
  const h=items.map(e=>'<div class="lv-'+e.level+'"><span class="t">'+(e.local||e.ts)+'</span>'+e.msg+(e.detail?'<span class="d">'+e.detail+'</span>':'')+'</div>').join('');
  const box=document.getElementById('logs');
  if(h!==lastHtml){lastHtml=h;box.innerHTML=h;if(!paused)box.scrollTop=box.scrollHeight;}
 }catch(e){}
 setTimeout(tick,3000);
}
tick();
</script></body></html>`;

http.createServer(async (req, res) => {
  const _t0 = Date.now();
  res.on('finish', () => {
    const _p = req.url.split('?')[0];
    if (['/', '/status', '/logs', '/qr', '/slider-qr'].includes(_p)) return;
    log('info', `HTTP ${req.method} ${_p} -> ${res.statusCode} (${Date.now() - _t0}ms)`);
  });
  const url = new URL(req.url, 'http://x');
  const json = (code, obj) => { res.writeHead(code, { 'Content-Type': 'application/json; charset=utf-8', 'Cache-Control': 'no-store' }); res.end(JSON.stringify(obj)); };
  const body = async () => { let b = ''; for await (const c of req) b += c; try { return JSON.parse(b || '{}'); } catch { return {}; } };
  if (req.method === 'GET' && url.pathname === '/') {
    res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' }); res.end(PAGE);
  } else if (req.method === 'GET' && url.pathname === '/qr') {
    if (st.qrB64) { res.writeHead(200, { 'Content-Type': 'image/jpeg', 'Cache-Control': 'no-store' }); res.end(Buffer.from(st.qrB64.split(',')[1], 'base64')); }
    else { res.writeHead(503, { 'Content-Type': 'text/plain; charset=utf-8' }); res.end('二维码生成中'); }
  } else if (req.method === 'GET' && url.pathname === '/slider-qr') {
    const su = st.sliderUrl || sms.sliderUrl;
    if (!su) { res.writeHead(404, { 'Content-Type': 'text/plain; charset=utf-8' }); res.end('当前无滑块验证链接'); return; }
    const buf = await QRCode.toBuffer(su, { width: 240, margin: 2 });
    res.writeHead(200, { 'Content-Type': 'image/png', 'Cache-Control': 'max-age=30' }); res.end(buf);
  } else if (req.method === 'GET' && url.pathname === '/logspage') {
    res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' }); res.end(LOG_PAGE);
  } else if (req.method === 'GET' && url.pathname === '/status') {
    json(200, JSON.parse(statusJson()));
  } else if (req.method === 'GET' && url.pathname === '/logs') {
    json(200, { ts: LOGS.length ? LOGS[LOGS.length - 1].ts : '', items: LOGS });
  } else if (req.method === 'GET' && url.pathname === '/sms/creds') {
    let c = {}; try { c = JSON.parse(fs.readFileSync(CREDS_PATH, 'utf8')); } catch {}
    json(200, { username: c.username || '', password: c.password || '' });
  } else if (req.method === 'POST' && url.pathname === '/sms/forget') {
    try { fs.unlinkSync(CREDS_PATH); } catch {}
    json(200, { ok: true, msg: '已清除' });
  } else if (req.method === 'POST' && url.pathname === '/channel') {
    const { channel: ch } = await body();
    if (!CHANNELS[ch]) { json(400, { ok: false, msg: `未知通道 ${ch}` }); return; }
    st.channel = ch; await newQR(true, `手动切换 ${ch}`);
    json(200, { ok: true, msg: `已切换 ${ch} 并取新码` });
  } else if (req.method === 'POST' && url.pathname === '/wx/code') {
    json(200, await wxCodeCompat(await body()));
  } else if (req.method === 'POST' && url.pathname === '/wx/refresh') {
    json(200, { status: true });
  } else if (req.method === 'POST' && url.pathname === '/newqr') {
    await newQR(true, '手动重新取码'); json(200, { ok: true, msg: '已重新取码' });
  } else if (req.method === 'POST' && url.pathname === '/logout') {
    const { alias } = await body();
    const acc = pickAcc(alias);
    if (!acc) { json(200, { ok: false, msg: '账号不存在' }); return; }
    log('info', `[${acc.alias}] 退出登录`);
    try { await api(`/api/Login/LogOut?wxid=${encodeURIComponent(acc.wxid)}`, { method: 'POST' }); } catch {}
    removeAccount(acc.wxid);
    json(200, { ok: true, msg: `已退出 ${accLabel(acc)}` });
  } else if (req.method === 'POST' && url.pathname === '/relogin') {
    const { alias } = await body();
    const acc = pickAcc(alias);
    if (!acc) { json(200, { ok: false, msg: '账号不存在' }); return; }
    json(200, await reloginAcc(acc, false));
  } else if (req.method === 'POST' && url.pathname === '/sms/apply') {
    const { username, password } = await body();
    json(200, await smsApply(username, password));
  } else if (req.method === 'POST' && url.pathname === '/sms/again') {
    json(200, await smsAgain());
  } else if (req.method === 'POST' && url.pathname === '/sms/verify') {
    const { code } = await body();
    json(200, await smsVerify(code));
  } else if (req.method === 'POST' && url.pathname === '/sms/qrapply') {
    json(200, await qrApply());
  } else if (req.method === 'POST' && url.pathname === '/test-notify') {
    log('info', '手动触发测试推送');
    json(200, await notify('docker-wx 登录台测试推送', `通道自检 ${fmtLocal(new Date())}`));
  } else {
    res.writeHead(404, { 'Content-Type': 'text/plain; charset=utf-8' }); res.end('Not Found');
  }
}).listen(PORT, () => log('info', `login-web v3 on :${PORT}`));
