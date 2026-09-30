/**
 * docker-wx 登录台 — 通用扫码登录网关
 * 页面:  GET /        状态仪表盘(自动刷新二维码/登录状态)
 * 图片:  GET /qr      当前登录二维码(jpeg,可直接在浏览器打开扫码)
 * 数据:  GET /status  JSON 状态
 * 操作:  POST /newqr  重新取码(切换账号)
 *        POST /relogin  62 凭证二次登录(免扫码)
 *        POST /logout   退出登录
 * 凭证: /data/login_info.json(deviceId/data62/wxid)
 */
const http = require('http');
const fs = require('fs');

const API = process.env.wxapi_url || 'http://wxapi:8057';
const PORT = process.env.port || 8058;
const INFO_PATH = process.env.info_path || '/data/login_info.json';
const DEVICE_NAME = 'docker-wx-pad';

const st = { phase: 'boot', msg: '初始化', uuid: '', qrB64: '', qrTs: 0, expireTs: 0, wxid: '', nick: '', headUrl: '', deviceId: '', loginTime: '', lastEvent: '' };
let info = loadInfo();

function loadInfo() { try { return JSON.parse(fs.readFileSync(INFO_PATH, 'utf8')); } catch { return {}; } }
function saveInfo() { try { fs.writeFileSync(INFO_PATH, JSON.stringify(info, null, 2)); } catch {} }
function genDeviceId() { let s = ''; for (let i = 0; i < 15; i++) s += Math.floor(Math.random() * 10); return s; }

async function api(path, opt) {
  const r = await fetch(`${API}${path}`, opt);
  return r.json();
}

async function newQR(force) {
  if (st.phase === 'qr' && !force && Date.now() < st.expireTs) return;
  const deviceId = info.deviceId || genDeviceId();
  try {
    const qr = await api('/api/Login/GetQRPad', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ DeviceID: deviceId, DeviceName: DEVICE_NAME }),
    });
    if (!qr.Success) { Object.assign(st, { phase: 'error', msg: `取码失败: ${qr.Message}` }); return; }
    info.deviceId = deviceId;
    info.data62 = qr.Data62 || info.data62 || '';
    saveInfo();
    Object.assign(st, {
      phase: 'qr', msg: '等待扫码', uuid: qr.Data.Uuid, qrB64: qr.Data.QrBase64,
      qrTs: Date.now(), expireTs: Date.now() + 4.5 * 60 * 1000, deviceId,
      wxid: info.wxid || '', loginTime: info.loginTime || '',
    });
  } catch (e) { Object.assign(st, { phase: 'error', msg: `wxapi 不可达: ${e.message}` }); }
}

function deepFindWxid(obj, d = 0) {
  if (!obj || typeof obj !== 'object' || d > 5) return null;
  for (const v of Object.values(obj)) {
    if (typeof v === 'string' && /^wxid_[0-9a-zA-Z_-]{6,}$/.test(v)) return v;
    const r = deepFindWxid(v, d + 1);
    if (r) return r;
  }
  return null;
}

async function poll() {
  if (st.phase !== 'qr') return;
  if (Date.now() > st.expireTs) { await newQR(true); return; }
  try {
    const res = await api(`/api/Login/CheckQR?uuid=${encodeURIComponent(st.uuid)}`, { method: 'POST' });
    if (res.Message === '登陆成功' || (res.Success && res.Data && res.Data.AcctSectResp)) {
      const d = res.Data || {};
      const wxid = d.AcctSectResp?.UserName || deepFindWxid(d) || info.wxid || '';
      info.wxid = wxid;
      info.nick = d.AcctSectResp?.NickName || info.nick || '';
      info.loginTime = new Date().toISOString();
      saveInfo();
      Object.assign(st, { phase: 'ok', msg: '登录成功', wxid, nick: info.nick, loginTime: info.loginTime, lastEvent: '' });
      console.log(`[login] ok wxid=${wxid}`);
      return;
    }
    const s = res.Data && res.Data.status;
    if (s === 1) st.msg = '已扫码,请在手机上点确认';
    else if (res.Code === -3) st.msg = '触发验证码流程,请点「重新取码」';
    else if (res.Code === -8) await newQR(true);
  } catch {}
}

async function relogin() {
  if (!info.data62) return { ok: false, msg: '无 62 凭证(先扫码登录一次)' };
  try {
    const res = await api('/api/Login/62data', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Data62: info.data62, DeviceName: DEVICE_NAME }),
    });
    const ok = res.Success && res.Data && (res.Data.AcctSectResp || deepFindWxid(res.Data));
    if (ok) {
      const wxid = res.Data.AcctSectResp?.UserName || deepFindWxid(res.Data);
      info.wxid = wxid; info.loginTime = new Date().toISOString(); saveInfo();
      Object.assign(st, { phase: 'ok', msg: '二次登录成功(62)', wxid, nick: info.nick, loginTime: info.loginTime });
      return { ok: true, msg: `二次登录成功 ${wxid}` };
    }
    return { ok: false, msg: `二次登录失败: ${res.Message || '未知'}(62 凭证可能已失效,请扫码)` };
  } catch (e) { return { ok: false, msg: `请求失败: ${e.message}` }; }
}

async function logout() {
  if (!info.wxid) return { ok: false, msg: '未登录' };
  try {
    await api(`/api/Login/LogOut?wxid=${encodeURIComponent(info.wxid)}`, { method: 'POST' });
  } catch {}
  const old = info.wxid;
  info.wxid = ''; info.loginTime = ''; saveInfo();
  Object.assign(st, { phase: 'qr', msg: '已退出,等待扫码', wxid: '', loginTime: '' });
  await newQR(true);
  return { ok: true, msg: `已退出 ${old}` };
}

setInterval(poll, 3000);
(async () => { if (info.wxid) Object.assign(st, { phase: 'ok', msg: '已登录(历史会话)', wxid: info.wxid, nick: info.nick || '', loginTime: info.loginTime || '' }); else await newQR(true); })();

function statusJson() {
  return JSON.stringify({ ...st, qrB64: undefined, serverTime: new Date().toISOString() });
}

const PAGE = `<!doctype html><html lang="zh"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>微信登录台 · docker-wx</title>
<style>
:root{--bg:#0f1115;--card:#171a21;--line:#232833;--tx:#e6e9ef;--sub:#8b93a3;--ok:#3fb96f;--warn:#e0a23c;--err:#e05c5c;--acc:#4f8ef7}
*{box-sizing:border-box;margin:0;padding:0}
body{background:var(--bg);color:var(--tx);font:15px/1.6 -apple-system,"Segoe UI","PingFang SC","Microsoft YaHei",sans-serif;min-height:100vh;display:flex;align-items:center;justify-content:center;padding:24px}
.wrap{width:100%;max-width:920px}
h1{font-size:20px;font-weight:600;display:flex;align-items:center;gap:10px;margin-bottom:16px}
.dot{width:9px;height:9px;border-radius:50%;background:var(--sub)}
.dot.qr{background:var(--warn);box-shadow:0 0 8px var(--warn)}
.dot.ok{background:var(--ok);box-shadow:0 0 8px var(--ok)}
.dot.err{background:var(--err)}
.grid{display:grid;grid-template-columns:340px 1fr;gap:16px}
@media(max-width:760px){.grid{grid-template-columns:1fr}}
.card{background:var(--card);border:1px solid var(--line);border-radius:14px;padding:20px}
.card h2{font-size:14px;color:var(--sub);font-weight:600;letter-spacing:.05em;margin-bottom:14px;text-transform:uppercase}
.qrbox{display:flex;flex-direction:column;align-items:center;gap:12px}
.qrbox img{width:264px;height:264px;border-radius:10px;background:#fff;padding:10px}
.qrbox img.err{display:flex;content:'';background:var(--card)}
.tip{color:var(--sub);font-size:13px;text-align:center}
.tip b{color:var(--tx)}
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
#toast{position:fixed;top:18px;left:50%;transform:translateX(-50%);background:#1f2530;border:1px solid var(--line);padding:10px 18px;border-radius:9px;font-size:13px;display:none;z-index:9}
</style></head><body><div class="wrap">
<h1><span class="dot" id="dot"></span>微信登录台 <span style="color:var(--sub);font-size:13px;font-weight:400">docker-wx · 安卓Pad 8.0.53</span></h1>
<div class="grid">
<div class="card"><h2>扫码登录</h2>
<div class="qrbox">
<div id="qrArea"><img id="qr" src="/qr"></div>
<div class="tip" id="qrTip">微信扫一扫,有效期约 5 分钟,过期自动刷新</div>
</div></div>
<div class="card"><h2>状态</h2><div id="stateBody"></div>
<div class="btns">
<button class="primary" onclick="act('newqr','确认重新取码?当前二维码作废')">重新取码</button>
<button onclick="act('relogin','用已保存的62凭证免扫码登录?')">二次登录</button>
<button onclick="act('logout','确认退出当前账号?')">退出登录</button>
</div></div></div>
<div class="foot">API: ${API} · 图片直链 <a style="color:var(--acc)" href="/qr" target="_blank">/qr</a> · 状态 <a style="color:var(--acc)" href="/status" target="_blank">/status</a></div>
</div><div id="toast"></div>
<script>
let lastUuid='';let lastPhase='';
async function tick(){
 try{
  const s=await (await fetch('/status')).json();
  const dot=document.getElementById('dot');
  dot.className='dot '+(s.phase==='ok'?'ok':s.phase==='error'?'err':'qr');
  const body=document.getElementById('stateBody');
  const rows=(arr)=>arr.map(([k,v])=>'<div class="row"><span class="k">'+k+'</span><span class="v">'+(v||'—')+'</span></div>').join('');
  if(s.phase==='ok'){
    body.innerHTML='<div class="ok-big">✔ '+s.msg+'</div>'+rows([['wxid',s.wxid],['昵称',s.nick],['登录时间',s.loginTime?new Date(s.loginTime).toLocaleString():''],['设备ID',s.deviceId]]);
    document.getElementById('qrArea').innerHTML='<div style="width:264px;height:264px;display:flex;align-items:center;justify-content:center;border-radius:10px;background:rgba(63,185,111,.08);color:var(--ok);font-size:15px">已在线,无需扫码</div>';
  }else{
    body.innerHTML=rows([['状态',s.msg],['wxid',s.wxid||'未登录'],['设备ID',s.deviceId],['二维码到期',s.expireTs?new Date(s.expireTs).toLocaleTimeString():'—']]);
    document.getElementById('qrTip').innerHTML='<b>'+ (s.phase==='error'?s.msg:(s.msg==='已扫码,请在手机上点确认'?s.msg:'微信扫一扫,过期自动刷新')) +'</b>';
    if(s.uuid&&s.uuid!==lastUuid){
      lastUuid=s.uuid;
      document.getElementById('qrArea').innerHTML='<img id="qr" src="/qr?ts='+Date.now()+'">';
    }
  }
  lastPhase=s.phase;
 }catch(e){}
 setTimeout(tick,3000);
}
tick();
function toast(m){const t=document.getElementById('toast');t.textContent=m;t.style.display='block';setTimeout(()=>t.style.display='none',2600);}
async function act(k,confirmMsg){if(confirmMsg&&!confirm(confirmMsg))return;toast('执行中...');
 try{const r=await(await fetch('/'+k,{method:'POST'})).json();toast(r.msg||r.ok?'完成':'失败');}catch(e){toast('请求失败');}}
</script></body></html>`;

http.createServer((req, res) => {
  const url = new URL(req.url, 'http://x');
  if (req.method === 'GET' && url.pathname === '/') {
    res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' }); res.end(PAGE);
  } else if (req.method === 'GET' && url.pathname === '/qr') {
    if (st.qrB64) {
      const b64 = st.qrB64.split(',')[1];
      const buf = Buffer.from(b64, 'base64');
      res.writeHead(200, { 'Content-Type': 'image/jpeg', 'Cache-Control': 'no-store' });
      res.end(buf);
    } else {
      res.writeHead(503, { 'Content-Type': 'text/plain; charset=utf-8' });
      res.end(st.phase === 'ok' ? '已登录,无需二维码;重新取码请访问页面操作' : '二维码生成中,稍后刷新');
    }
  } else if (req.method === 'GET' && url.pathname === '/status') {
    res.writeHead(200, { 'Content-Type': 'application/json; charset=utf-8', 'Cache-Control': 'no-store' });
    res.end(statusJson());
  } else if (req.method === 'POST' && url.pathname === '/newqr') {
    newQR(true).then(() => { res.writeHead(200, { 'Content-Type': 'application/json' }); res.end(JSON.stringify({ ok: true, msg: '已重新取码' })); });
  } else if (req.method === 'POST' && url.pathname === '/relogin') {
    relogin().then((r) => { res.writeHead(200, { 'Content-Type': 'application/json' }); res.end(JSON.stringify(r)); });
  } else if (req.method === 'POST' && url.pathname === '/logout') {
    logout().then((r) => { res.writeHead(200, { 'Content-Type': 'application/json' }); res.end(JSON.stringify(r)); });
  } else {
    res.writeHead(404, { 'Content-Type': 'text/plain; charset=utf-8' }); res.end('Not Found');
  }
}).listen(PORT, () => console.log(`login-web on :${PORT}, api=${API}`));
