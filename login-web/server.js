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
    acc.loginTime = new Date().toISOString(); acc.lastHbOk = acc.loginTime; acc.hbFails = 0; acc.offlineNotified = false; acc.offline = false;
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
  if (s === 1) {
    if (st.msg !== '已扫码,请在手机上点确认') { log('info', 'CheckQR: 已扫码,等待确认'); st.expireTs = Math.max(st.expireTs, Date.now() + 3 * 60_000); }
    st.msg = '已扫码,请在手机上点确认';
    return;
  }
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

/* ---------------- 小程序注册状态矩阵 ----------------
   青龙采集器每天上报(POST /registry/report),聚合存 /data/registry.json
   状态: ok=已注册/正常, unreg=未注册/需绑定手机号, fail=其他失败 */
const REG_PATH = '/data/registry.json';
function loadReg() { try { return JSON.parse(fs.readFileSync(REG_PATH, 'utf8')); } catch { return { updated: '', scripts: {} }; } }
function saveReg(r) { try { fs.writeFileSync(REG_PATH, JSON.stringify(r, null, 2)); } catch {} }

/* 一键检测:立即对所有账号心跳并返回结果 */
async function checkNow() {
  if (!accounts.length) return { ok: false, msg: '无账号', lines: [] };
  const startIdx = LOGS.length;
  log('info', `开始在线检测:共 ${accounts.length} 个账号`);
  for (const acc of accounts) {
    log('info', `[${acc.alias}] 检测 ${acc.nick || acc.wxid} …`);
    await heartbeat(acc);
  }
  const off = accounts.filter(a => a.offline).length;
  const r = `检测完成:${accounts.length} 个账号,${accounts.length - off} 在线${off ? ',' + off + ' 掉线' : ''}`;
  log(off ? 'error' : 'info', r);
  return { ok: !off, msg: r, lines: LOGS.slice(startIdx).map(e => ({ t: (e.local || '').slice(11), level: e.level, msg: e.msg })) };
}

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

/* ---------------- JSLogin 统一闸门 ----------------
   微信对 js-login 按【微信号】限频,额度跨小程序共享(jingjianx 同号 12s 连登 5 店后,
   该号换任何小程序取码全 -13000 的实证)。策略:同 wxid 串行 + 令牌桶(短期 3/20s、
   长期 8/10min)+ 超发排队 12s 后干净失败(不把额度打爆进冷却);-13000 = 整号 15 分钟
   快速失败(冷却期内换 appid 也拒发,不再白烧额度)。
   注:code 是一次性消费品(业务后端 jscode2session 即作废),不可跨登录复用;
   唯一安全复用 = 微信官方 getallphone 返回的配对 code(见 /wx/getphonenumber)。 */
const jsGate = {};         // wxid -> { tail, times[], coolUntil }
const JS_MIN_GAP = 1500;       // 同微信号两次 js-login 最小间隔
const JS_WIN1_MS = 20_000, JS_WIN1_MAX = 3;        // 短期突发上限
const JS_WIN2_MS = 600_000, JS_WIN2_MAX = 8;       // 长期总量上限
const JS_QUEUE_WAIT = 12_000;   // 令牌桶满时最长排队(脚本侧 HTTP 超时多为 15s+)

function sleep(ms) { return new Promise(r => setTimeout(r, ms)); }
function gateOf(wxid) {
  if (!jsGate[wxid]) jsGate[wxid] = { tail: Promise.resolve(), times: [], coolUntil: 0 };
  return jsGate[wxid];
}

/* 所有 JSLogin 必经:按 wxid 串行 + 节流;整号冷却中抛 {gate:true,message} */
function jsLoginViaGate(acc, appid) {
  const st = gateOf(acc.wxid);
  const run = (async () => {
    if (st.coolUntil > Date.now()) {
      const left = Math.ceil((st.coolUntil - Date.now()) / 60_000);
      throw { gate: true, message: `微信限频冷却中(约剩${left}分钟),该微信号全部小程序取码暂停` };
    }
    const last = st.times[st.times.length - 1];
    if (last) await sleep(Math.max(0, JS_MIN_GAP - (Date.now() - last)));
    const deadline = Date.now() + JS_QUEUE_WAIT;
    for (;;) {
      const now = Date.now();
      st.times = st.times.filter(t => now - t < JS_WIN2_MS);
      const short = st.times.filter(t => now - t < JS_WIN1_MS).length;
      if ((short < JS_WIN1_MAX && st.times.length < JS_WIN2_MAX) || now > deadline) break;
      await sleep(400);
    }
    st.times.push(Date.now());
    return await api('/api/Wxapp/JSLogin', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Wxid: acc.wxid, Appid: appid }),
    });
  })();
  st.tail = run.catch(() => {});
  return run;
}

/* 统一解析 JSLogin 响应:成功 {ok,code,sk,openid};-13000 置整号冷却 */
function parseJsLogin(acc, appid, res) {
  const d = res.Data || {};
  const code = d.code || d.Code || (d.Data && (d.Data.code || d.Data.Code));
  const sk = d.sessionKey || d.SessionKey || (d.Data && (d.Data.sessionKey || d.Data.SessionKey));
  const openid = d.openid || d.Openid;
  if (res.Success && code) return { ok: true, code: String(code), sk: sk ? String(sk) : '', openid: openid ? String(openid) : '' };
  const errCode = d.jsapiBaseresponse && d.jsapiBaseresponse.errcode;
  const err = (d.jsapiBaseresponse && (d.jsapiBaseresponse.errcode + ' ' + d.jsapiBaseresponse.errmsg)) || res.Message || '空code';
  if (String(errCode) === '-13000') {
    gateOf(acc.wxid).coolUntil = Date.now() + 15 * 60_000;
    log('warn', `[${acc.alias}] js-login 触发微信限频(-13000),微信号进入 15 分钟冷却(appid=${appid})`);
    return { ok: false, freq: true, message: '微信限频(-13000),该微信号已进入 15 分钟冷却(冷却期内所有小程序取码暂停)' };
  }
  return { ok: false, message: `${res.Message || ''} jsapi=${err}` };
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

/* /wx/getuserinfo:JSLogin 取 code + JSOperateWxData(Opt=2) 取 encryptedData/iv */
async function wxGetUserInfoCompat(body) {
  const appid = body && body.appid;
  if (!appid) return { status: false, message: '缺少 appid' };
  const key = body.openid ? String(body.openid).split('#')[0].trim() : '';
  let acc = null;
  if (key) acc = accounts.find(a => a.alias === key || a.wxid === key) || null;
  if (!acc) acc = accounts[0];
  if (!acc) return { status: false, message: '无已登录微信' };
  try {
    const login = await jsLoginViaGate(acc, appid);
    const p = parseJsLogin(acc, appid, login);
    if (!p.ok) return { status: false, message: `JSLogin 失败: ${p.message}` };
    const code = p.code;
    const oper = await api('/api/Wxapp/JSOperateWxData', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Wxid: acc.wxid, Appid: appid, Opt: 2, Data: '' }),
    });
    const b64 = oper.Data && oper.Data.data;
    if (!oper.Success || !b64) return { status: false, message: `getUserInfo 失败: ${(oper.Data && oper.Data.jsapiBaseresponse && oper.Data.jsapiBaseresponse.errmsg) || oper.Message}` };
    const inner = JSON.parse(Buffer.from(b64, 'base64').toString('utf-8'));
    if (!code || !inner.encryptedData || !inner.iv) return { status: false, message: '授权数据不完整' };
    let userInfo = {};
    try { userInfo = JSON.parse(inner.data); } catch {}
    log('info', `[${acc.alias}] getuserinfo 成功 appid=${appid}`);
    return { status: true, data: { code: String(code), encryptedData: inner.encryptedData, iv: inner.iv, signature: inner.signature, cloud_id: inner.cloud_id, userInfo } };
  } catch (e) { return { status: false, message: e && e.message }; }
}

/* ---------------- /wx/getphonenumber ----------------
   走 wxapi GetAllMobile(mmbiz customphone/getallphone):微信官方返回本账号绑定
   手机号的 encryptedData/iv/cloud_id(+配对 code),内容与在小程序里手点一次
   "允许"完全一致;不经 js-login,不占 -13000 限额。登录用的 wx.login code
   仍由脚本随后经 /wx/code 自取(标准小程序会话模型,session_key 同源可解)。 */
async function wxGetPhoneNumberCompat(body) {
  const appid = body && body.appid;
  if (!appid) return { status: false, message: '缺少 appid' };
  const key = body.openid ? String(body.openid).split('#')[0].trim() : '';
  let acc = null;
  if (key) acc = accounts.find(a => a.alias === key || a.wxid === key) || null;
  if (!acc) acc = accounts[0];
  if (!acc) return { status: false, message: '无已登录微信' };
  try {
    const res = await api('/api/Wxapp/GetAllMobile', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Wxid: acc.wxid, Appid: appid }),
    });
    const d = res.Data || {};
    const err = d.jsapiBaseresponse && d.jsapiBaseresponse.errcode;
    if (!res.Success || (err !== undefined && err !== null && String(err) !== '0')) {
      return { status: false, message: `获取手机号授权数据失败: ${(d.jsapiBaseresponse && d.jsapiBaseresponse.errmsg) || res.Message || ('errcode ' + err)}` };
    }
    // 首选内层 Data 字符串(wx_phone / custom_phone_list),兜底 ALLMobile 数组
    let items = [];
    try {
      const inner = typeof d.Data === 'string' ? JSON.parse(d.Data) : d.Data;
      if (inner && Array.isArray(inner.custom_phone_list)) items = inner.custom_phone_list;
      else if (inner && inner.wx_phone) items = [inner.wx_phone];
    } catch {}
    if (!items.length && Array.isArray(d.ALLMobile)) items = d.ALLMobile;
    const accMobile = String(acc.mobile || '').replace(/\D/g, '');
    const item = items.find(x => x && x.encryptedData && (!accMobile || String(x.mobile || '').replace(/\D/g, '') === accMobile))
      || items.find(x => x && x.encryptedData);
    if (!item) return { status: false, message: `微信未返回手机号授权数据(${items.length} 条记录均无 encryptedData,可能需先在小程序内授权一次)` };
    let phoneCode = '';
    try { const dd = typeof item.data === 'string' ? JSON.parse(item.data) : item.data; phoneCode = (dd && dd.code) || item.code || ''; } catch {}
    log('info', `[${acc.alias}] getphonenumber 成功(微信官方数据) appid=${appid} 手机号=${item.show_mobile || item.mobile}`);
    const raw = { mobile: item.mobile || '', show_mobile: item.show_mobile || '', encryptedData: item.encryptedData, iv: item.iv, cloud_id: item.cloud_id || '', code: phoneCode };
    return {
      status: true, code: phoneCode, phone: item.mobile || '', encryptedData: item.encryptedData, iv: item.iv, cloud_id: item.cloud_id || '',
      data: { code: phoneCode, phone: item.mobile || '', phoneNumber: item.mobile || '', encryptedData: item.encryptedData, iv: item.iv, cloud_id: item.cloud_id || '', raw },
    };
  } catch (e) { return { status: false, message: e && e.message }; }
}

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
  let lastMsg = '';
  for (let i = 0; i < 3; i++) {
    if (i) await sleep(6000);
    try {
      const res = await jsLoginViaGate(acc, appid);
      const p = parseJsLogin(acc, appid, res);
      if (p.ok) {
        if (i) log('info', `[${acc.alias}] 小程序取码成功(第${i + 1}次尝试) appid=${appid}`);
        else log('info', `[${acc.alias}] 小程序取码成功 appid=${appid}`);
        // 平铺 code + 嵌套 data.code 双路径:脚本解析两种风格并存(提现免费券.py 只读 data.data.code,camel.js 等读 code || data.code)
        return { status: true, code: p.code, data: { code: p.code } };
      }
      lastMsg = p.message;
      if (p.freq) return { status: false, message: p.message };
      log('warn', `[${acc.alias}] 小程序取码空code(尝试${i + 1}/3) appid=${appid}: ${lastMsg}`);
    } catch (e) {
      if (e && e.gate) return { status: false, message: e.message };
      lastMsg = e.message; log('error', `[${acc.alias}] 取码异常(尝试${i + 1}/3): ${e.message}`); }
  }
  return { status: false, message: `取码失败: ${lastMsg}` };
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
#modal,#checkModal{position:fixed;inset:0;background:rgba(0,0,0,.6);display:none;align-items:center;justify-content:center;z-index:50;padding:16px}
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
<button onclick="openCheck()">检测状态</button>
<button onclick="location.href=&apos;/registrypage&apos;">注册状态</button>
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
<div id="checkModal"><div class="mbox">
<div class="mhd"><b>在线状态检测</b><button class="small" onclick="closeCheck()">关闭</button></div>
<div class="success" id="checkBar"></div>
<div id="checkLogs" style="max-height:300px;overflow-y:auto;font:12px/1.9 ui-monospace,Menlo,Consolas,monospace;color:var(--sub);background:#12141a;border-radius:8px;padding:10px">准备检测…</div>
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
let checkOpen=false,checkStart='',checkResult='';
function openCheck(){checkOpen=true;checkResult='';document.getElementById('checkBar').style.display='none';document.getElementById('checkModal').style.display='flex';document.getElementById('checkLogs').innerHTML='正在逐账号检测…';
 fetch('/checknow',{method:'POST'}).then(r=>r.json()).then(j=>{
  const bar=document.getElementById('checkBar');bar.style.display='block';bar.style.background=j.ok?'rgba(63,185,111,.12)':'rgba(224,92,92,.12)';bar.style.borderColor=j.ok?'rgba(63,185,111,.4)':'rgba(224,92,92,.4)';bar.style.color=j.ok?'var(--ok)':'var(--err)';bar.textContent=(j.ok?'✅ ':'⚠️ ')+j.msg;
  const ch=(j.lines||[]).map(e=>'<div style="color:'+(e.level==='error'?'#e05c5c':e.level==='warn'?'#e0a23c':'#e6e9ef')+'"><span style="color:#5b6472;margin-right:6px">'+e.t+'</span>'+e.msg+'</div>').join('');
  const cb=document.getElementById('checkLogs');if(ch){cb.innerHTML=ch;cb.scrollTop=cb.scrollHeight;}else{cb.textContent=j.msg;}
 }).catch(e=>{document.getElementById('checkLogs').textContent='检测请求失败:'+e.message;});}
function closeCheck(){checkOpen=false;document.getElementById('checkModal').style.display='none';}
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

/* ---------------- 注册状态矩阵页 ---------------- */
const REG_PAGE = `<!doctype html><html lang="zh"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>小程序注册状态 · docker-wx</title>
<style>
:root{--bg:#0f1115;--card:#171a21;--line:#232833;--tx:#e6e9ef;--sub:#8b93a3;--ok:#3fb96f;--warn:#e0a23c;--err:#e05c5c;--acc:#4f8ef7;--dim:#5b6472}
*{box-sizing:border-box;margin:0;padding:0}
body{background:var(--bg);color:var(--tx);font:14px/1.6 -apple-system,"Segoe UI","PingFang SC",sans-serif;padding:22px}
.wrap{max-width:1100px;margin:0 auto}
.hd{display:flex;align-items:center;gap:8px;margin-bottom:8px;flex-wrap:wrap}
h1{font-size:19px;font-weight:600;flex:1}
.hd .meta{color:var(--sub);font-size:12px;width:100%}
button{background:#1f2530;color:var(--tx);border:1px solid var(--line);border-radius:8px;padding:6px 12px;font-size:12px;cursor:pointer}
button.on{border-color:var(--acc);color:var(--acc)}
.tabs{display:flex;gap:8px;margin:10px 0}
.stat{display:flex;gap:14px;font-size:12px;color:var(--sub);margin-bottom:10px;flex-wrap:wrap}
.stat b{color:var(--tx)}
#tableWrap{overflow:auto;-webkit-overflow-scrolling:touch;border:1px solid var(--line);border-radius:12px}
table{width:100%;border-collapse:collapse;background:var(--card)}
th{position:sticky;top:0;z-index:2}
td.l,th:first-child{position:sticky;left:0;z-index:1}
th:first-child{z-index:3}
th,td{padding:8px 10px;border-bottom:1px solid var(--line);text-align:center;font-size:13px}
th{color:var(--sub);font-weight:600;font-size:12px;background:#141822}
td.l{text-align:left;background:var(--card)}
tr:hover td{background:#1b2029}
tr:hover td.l{background:#1b2029}
@media(max-width:640px){
 body{padding:10px 10px 30px}
 h1{font-size:16px;margin-bottom:2px;flex-basis:100%;white-space:nowrap}
 .hd{gap:6px;margin-bottom:4px}
 .hd .meta{font-size:11px}
 button{padding:7px 9px;font-size:12px}
 .tabs{gap:6px;margin:6px 0}
 .stat{gap:8px;margin-bottom:4px;font-size:11px}
 .legend{font-size:11px;line-height:1.9;margin-bottom:6px}
 th,td{padding:5px 4px!important;font-size:12px}
 th{font-size:11px}
}
.ok{color:var(--ok)}.unreg{color:var(--warn)}.fail{color:var(--err)}.unknown{color:var(--dim)}
.bar{display:inline-block;width:8px;height:8px;border-radius:4px;margin-right:6px}
.empty{color:var(--sub);text-align:center;padding:50px 0}
.bdg{cursor:help;font-size:11px;white-space:nowrap}
.legend{font-size:12px;color:var(--sub);margin-bottom:10px}
.legend .bdg{font-size:13px}
.acc-card{background:var(--card);border:1px solid var(--line);border-radius:12px;padding:14px;margin-bottom:10px;display:flex;gap:16px;align-items:center;flex-wrap:wrap}
.acc-card .big{font-size:22px;font-weight:700}
</style></head><body><div class="wrap">
<div class="hd">
<h1>小程序注册状态矩阵</h1>
<button id="filterAll" class="on" onclick="setF('all')">全部</button>
<button id="filterUnreg" onclick="setF('unreg')">未注册</button>
<button id="filterSpec" onclick="setF('spec')">需处理</button>
<button onclick="location.href=&apos;/&apos;">返回</button>
<div class="meta" id="meta"></div>
</div>
<div class="tabs">
<button class="on" id="tab-m" onclick="setTab('m')">按小程序</button>
<button id="tab-a" onclick="setTab('a')">按账号</button>
</div>
<div class="stat" id="stat"></div>
<div class="legend" id="legend"></div>
<div id="body"></div>
</div>
<script>
const NEED_HAR={fuyouhui:'需抓包:填 fuyouhui_token',hisense_aijia:'需抓包:hisense_aijia_token/customerId/loginKey/refreshToken/sign_task_id(5个,可手机授权免抓但协议不支持)',hongsehuojian:'需抓包:填 hshj_ticket(App内已登录token,手机授权协议不支持)',longfor:'需抓包:填 longfor_dx_token',qqpcmgr:'需抓包:qqpcmgr_authCode/guid/lid/sdiaid/computer/ua/version 等8个',roki:'需抓包:填 roki_api_base',yichengtong:'需抓包:填 yichengtong_token'};
/* 手机号授权已由 wxlogin /wx/getphonenumber 走 GetAllMobile(微信官方 edata/iv/phoneCode)打通,
   原 9 个"需手机授权"脚本 2026-10-07 实测全通(haitian/jx/jdbclub/colorful/rytyn/wrn/yipiaoda/jyxe 签到成功,
   dfmfs 登录通但业务需扫产品红包码,移入 NEED_REG);NEED_PHONE 留空待未来新脚本 */
const NEED_PHONE={};
const NEED_REG={aiguo:'需注册:手机微信打开小程序完成注册/授权一次',aima:'需注册:爱玛会员俱乐部登录报110502用户不存在,需小程序内注册会员一次',dfmfs:'需注册:手机号授权登录已通,但业务要求扫产品顶部红包码激活/延长日常活动有效期后才能签到',ardywj:'需注册:爱康需完成手机号授权/注册(小程序内)',bydhy:'需注册:手机微信打开小程序完成注册/授权一次',dfrc:'需注册:手机微信打开小程序完成注册/授权一次',dsmmhy:'需注册:袋鼠妈妈有赞平台需注册会员',fmy:'需注册:手机微信打开小程序完成注册/授权一次',hougongfang:'需注册:手机微信打开小程序完成注册/授权一次',hxek:'需注册:手机微信打开小程序完成注册/授权一次',junpinhui:'需注册:手机微信打开小程序完成注册/授权一次',lmf:'需注册:手机微信打开小程序完成注册/授权一次',lthwy:'需注册:手机微信打开小程序完成注册/授权一次',mdhy:'需注册:美的M-VIP需注册会员',mobil:'需注册:美孚臻享俱乐部需注册会员',mtyl:'需注册:手机微信打开小程序完成注册/授权一次',nndj:'需注册:牛牛短剧需注册',nxdc:'需注册:奈雪需注册会员',olecs:'需注册:Ole需注册会员',parkson:'需注册:百盛呼啦圈需注册会员',qmsd:'需注册:全棉时代需注册会员',qqhyjlb:'需注册:洽洽会员俱乐部需注册会员',quanmianshidai:'需注册:手机微信打开小程序完成注册/授权一次',quncrm:'需注册:群脉平台需注册会员',qyqd:'需注册:手机微信打开小程序完成注册/授权一次',rio:'需注册:RIO微醺俱乐部需注册会员',rrk:'需注册:手机微信打开小程序完成注册/授权一次',sf:'需注册:顺丰需注册/绑定会员',shanyi:'需注册:手机微信打开小程序完成注册/授权一次',smgc:'需注册:SM广场需注册会员(多城市)',tjg:'需注册:手机微信打开小程序完成注册/授权一次',trsj:'需注册:甜润世界需注册会员',txq:'需注册:汤星球需注册会员',wanjiale:'需注册:万家乐需注册会员',wuyingyundiannao:'需注册:无影云电脑静默登录可能被安全验证拦,被拦时需手填 wuying_token',wx_xlxyh:'需注册:骁龙骁友会需注册',wzy:'需注册:喂自由需注册',xiaodangjia:'需注册:小铛家需注册',xinxianghui:'需注册:手机微信打开小程序完成注册/授权一次',xmsq:'需注册:小米社区需绑定账号',xzyy:'需注册:小紫有约需注册',yjlxh:'需注册:伊家乐享会(伊利)需注册会员',youzan:'需注册:有赞店铺(临水玉泉/TOI/七点五等)需注册会员',yuexihui:'需注册:中粮悦喜荟需注册会员',yzyj:'需注册:微盟onecrm需注册会员',jdcode:'需注册:非签到,采集京东JD_COOKIE;需在京东小程序内绑定京东账号',juziyingtao:'需注册:橘子樱桃需在小程序内完成手机号授权(业务层)',sinsin:'需注册:sinsin需在小程序内完成手机号授权(业务层)',maopu:'需注册:猫扑不代填资料/不代过手机号授权,需小程序内完成'};
function badge(k){const t=NEED_HAR[k]||NEED_PHONE[k]||NEED_REG[k];if(!t)return '';return ' <span class="bdg" title="'+t+'">'+(NEED_HAR[k]?'📡':NEED_PHONE[k]?'📲':'📱')+'</span>';}
let tab='m', filter='all', data=null;
function setF(f){filter=f;document.getElementById('filterAll').className=f==='all'?'on':'';document.getElementById('filterUnreg').className=f==='unreg'?'on':'';document.getElementById('filterSpec').className=f==='spec'?'on':'';render();}
function setTab(t){tab=t;document.getElementById('tab-m').className=t==='m'?'on':'';document.getElementById('tab-a').className=t==='a'?'on':'';render();}
const ICON={ok:'<span class="ok" title="已注册">✅</span>',unreg:'<span class="unreg" title="未注册(需手机授权)">❌</span>',fail:'<span class="fail" title="其他失败">⚠️</span>',unknown:'<span class="unknown" title="无数据">➖</span>'};
function fmtPer(p){return p?ICON[p]||ICON.unknown:ICON.unknown;}
function render(){
 if(!data)return;
 const scripts=Object.entries(data.scripts||{});
 const aliases=[...new Set(scripts.flatMap(([,v])=>Object.keys(v.per||{})))].sort();
 const disp=a=>{const m=(data.aliases||{})[a];if(!m)return a;return m.mobile?m.mobile.slice(0,3)+'****'+m.mobile.slice(-4):(m.wx||m.nick||a);};
 document.getElementById('meta').textContent='数据更新: '+(data.updated?new Date(data.updated).toLocaleString():'无')+' · 来自每日任务执行日志(20:05 采集)';
 let rows=scripts.filter(([k,v])=>{
   if(filter==='unreg')return Object.values(v.per||{}).some(x=>x==='unreg');
   if(filter==='spec')return !!(NEED_HAR[k]||NEED_PHONE[k]);
   return true;
 });
 rows.sort((a,b)=>{const ua=Object.values(a[1].per||{}).filter(x=>x==='unreg').length,ub=Object.values(b[1].per||{}).filter(x=>x==='unreg').length;return ub-ua||a[1].name.localeCompare(b[1].name);});
 let ok=0,un=0,fl=0;
 scripts.forEach(([k,v])=>Object.values(v.per||{}).forEach(x=>{if(x==='ok')ok++;else if(x==='unreg')un++;else if(x==='fail')fl++;}));
 document.getElementById('stat').innerHTML='<span>✅ 已注册: <b>'+ok+'</b></span><span>❌ 未注册: <b>'+un+'</b></span><span>⚠️ 其他失败: <b>'+fl+'</b></span><span>➖ 无数据: <b>'+(scripts.length?scripts.filter(([k,v])=>!Object.keys(v.per||{}).length).length:0)+'</b></span><span>脚本数: <b>'+scripts.length+'</b></span>';
 document.getElementById('legend').innerHTML='状态:✅已注册 ❌未注册 ⚠️失败 ➖无数据<br>类型: <span class="bdg" title="需手动抓包获取token填变量,悬停各行徽章看具体变量">📡 需抓包 '+Object.keys(NEED_HAR).length+'</span> · <span class="bdg" title="登录依赖手机号授权,协议层不支持,须手机微信内操作一次">📲 需手机授权 '+Object.keys(NEED_PHONE).length+'</span> · <span class="bdg" title="手机微信打开该小程序,完成注册/授权一次后脚本才有产出">📱 需注册 '+Object.keys(NEED_REG).length+'</span> · 无标记=打开即用';
 if(tab==='m'){
   const short = window.matchMedia('(max-width:640px)').matches;
   let h='<table><tr><th style="text-align:left">小程序</th>'+aliases.map(a=>'<th title="'+disp(a)+'">'+(short && (data.aliases||{})[a] && (data.aliases||{})[a].mobile ? (data.aliases||{})[a].mobile.slice(-4) : disp(a))+'</th>').join('')+'</tr>';
   for(const [k,v] of rows){
     h+='<tr><td class="l">'+(v.name||k)+badge(k)+'</td>'+aliases.map(a=>'<td>'+fmtPer((v.per||{})[a])+'</td>').join('')+'</tr>';
   }
   document.getElementById('body').innerHTML='<div id="tableWrap">'+h+'</table></div>';
 }else{
   let h='';
   for(const a of aliases){
     const mine=scripts.filter(([k,v])=>(v.per||{}).hasOwnProperty(a));
     const myOk=mine.filter(([k,v])=>v.per[a]==='ok').length;
     const myUn=mine.filter(([k,v])=>v.per[a]==='unreg').length;
     h+='<div class="acc-card"><div><div class="big">'+disp(a)+'</div><div style="color:var(--sub);font-size:12px">标识 '+a+' · 已注册 '+myOk+' · 未注册 '+myUn+' · 共 '+mine.length+'</div></div>';
     if(myUn){
       h+='<div style="flex:1;min-width:280px"><div style="font-size:12px;color:var(--sub)">未注册清单:</div><div style="font-size:12px;color:var(--warn)">'+mine.filter(([k,v])=>v.per[a]==='unreg').map(([k,v])=>(v.name||k)+badge(k)).slice(0,60).join('、')+'</div></div>';
     } else h+='<div style="color:var(--ok)">全部已注册 ✓</div>';
     h+='</div>';
   }
   document.getElementById('body').innerHTML=h||'<div class="empty">无数据</div>';
 }
}
async function tick(){
 try{data=await(await fetch('/registry')).json();render();}catch(e){}
 setTimeout(tick,30000);
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
  } else if (req.method === 'POST' && url.pathname === '/wx/getuserinfo') {
    json(200, await wxGetUserInfoCompat(await body()));
  } else if (req.method === 'POST' && url.pathname === '/wx/getphonenumber') {
    json(200, await wxGetPhoneNumberCompat(await body()));
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
  } else if (req.method === 'POST' && url.pathname === '/registry/report') {
    const b = await body();
    const reg = loadReg();
    reg.updated = new Date().toISOString();
    for (const item of (Array.isArray(b) ? b : [])) {
      if (!item || !item.script) continue;
      reg.scripts[item.script] = { name: item.name || item.script, per: item.per || {}, ts: new Date().toISOString() };
    }
    saveReg(reg);
    log('info', `注册状态上报:${(Array.isArray(b) ? b : []).length} 个脚本`);
    json(200, { ok: true });
  } else if (req.method === 'GET' && url.pathname === '/registry') {
    const reg = loadReg();
    const amap = {};
    for (const a of accounts) amap[a.alias] = { wx: a.aliasWx || '', nick: a.nick || '', mobile: a.mobile || '' };
    json(200, { ...reg, aliases: amap });
  } else if (req.method === 'GET' && url.pathname === '/registrypage') {
    res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' }); res.end(REG_PAGE);
  } else if (req.method === 'POST' && url.pathname === '/checknow') {
    json(200, await checkNow());
  } else if (req.method === 'POST' && url.pathname === '/test-notify') {
    log('info', '手动触发测试推送');
    json(200, await notify('docker-wx 登录台测试推送', `通道自检 ${fmtLocal(new Date())}`));
  } else {
    res.writeHead(404, { 'Content-Type': 'text/plain; charset=utf-8' }); res.end('Not Found');
  }
}).listen(PORT, () => log('info', `login-web v3 on :${PORT}`));
