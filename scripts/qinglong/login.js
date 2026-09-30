/**
 * docker-wx 微信扫码登录 — 青龙面板
 * 用法:node login.js(取二维码 → 手机扫码确认 → 自动完成登录)
 * 凭证保存 login_info.json(deviceId/data62 供二次登录,wxid 供业务接口)
 */
const API = process.env.wxapi_url || 'http://wxapi:8057';
const fs = require('fs');
const INFO = '/ql/data/scripts/wxapi/login_info.json';

function genDeviceId() {
  let s = '';
  for (let i = 0; i < 15; i++) s += Math.floor(Math.random() * 10);
  return s;
}

// 兜底:从响应里递归找 wxid_xxx / wxid 形态字符串
function deepFindWxid(obj, depth = 0) {
  if (!obj || typeof obj !== 'object' || depth > 5) return null;
  for (const v of Object.values(obj)) {
    if (typeof v === 'string' && /^wxid_[0-9a-zA-Z_-]{6,}$/.test(v)) return v;
    const r = deepFindWxid(v, depth + 1);
    if (r) return r;
  }
  return null;
}

(async () => {
  let info = {};
  try { info = JSON.parse(fs.readFileSync(INFO, 'utf8')); } catch {}
  const deviceId = info.deviceId || genDeviceId();

  const r = await fetch(`${API}/api/Login/GetQRPad`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ DeviceID: deviceId, DeviceName: 'qinglong-bot' }),
  });
  const qr = await r.json();
  if (!qr.Success) { console.log('取码失败:', qr.Message); process.exit(1); }

  const { Uuid, QrUrl, QrBase64, ExpiredTime } = qr.Data;
  fs.writeFileSync('/ql/data/scripts/wxapi/qrcode.b64', QrBase64);
  console.log('UUID:', Uuid);
  console.log('QRURL:', QrUrl);
  console.log('过期时间:', ExpiredTime);
  console.log('--- 等待扫码(有效期约 5 分钟,每 3 秒轮询)---');

  let seen = '';
  for (let i = 0; i < 100; i++) {
    await new Promise((r) => setTimeout(r, 3000));
    let res;
    try {
      const cr = await fetch(`${API}/api/Login/CheckQR?uuid=${encodeURIComponent(Uuid)}`, { method: 'POST' });
      res = await cr.json();
    } catch { continue; }

    if (res.Message === '登陆成功' || (res.Success && res.Data && res.Data.AcctSectResp)) {
      const wxid = res.Data?.AcctSectResp?.UserName || deepFindWxid(res.Data);
      fs.writeFileSync(INFO, JSON.stringify({
        deviceId, wxid, data62: qr.Data62 || info.data62 || '',
        loginTime: new Date().toISOString(), uuid: Uuid,
      }, null, 2));
      console.log('=== 登录成功 ===');
      console.log('WXID:', wxid);
      console.log('NICK:', res.Data?.AcctSectResp?.NickName || '(见响应)');
      console.log('凭证已保存:', INFO);
      process.exit(0);
    }

    const st = res.Data && res.Data.status;
    const tag = st === 1 ? '已扫码,请在手机上点击确认' : st === 0 || st === undefined ? '等待扫码' : `状态 ${st}`;
    if (tag !== seen) { console.log(`[${new Date().toLocaleTimeString()}] ${tag}`); seen = tag; }

    if (res.Code === -3) { console.log('触发验证码流程(ticket),本脚本未实现,退出'); process.exit(2); }
    if (res.Code === -8 && res.Message && res.Message.includes('过期')) { console.log('二维码已过期,请重新运行'); process.exit(1); }
  }
  console.log('超时未完成扫码,请重新运行取新码');
  process.exit(1);
})();
