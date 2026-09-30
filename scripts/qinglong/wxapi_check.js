/**
 * docker-wx (smallfawn) 对接自检脚本 — 青龙面板
 * 用法:手动运行或建定时任务;node 18+ 内置 fetch,无需额外依赖
 * 通过 = 青龙容器可正常访问微信协议 API
 */
const API = process.env.wxapi_url || 'http://wxapi:8057';

function randDeviceId(len = 15) {
  const chars = '0123456789';
  let s = '';
  for (let i = 0; i < len; i++) s += chars[Math.floor(Math.random() * chars.length)];
  return s;
}

(async () => {
  // 1. 服务可达性
  try {
    const r = await fetch(`${API}/`);
    console.log(`[1/2] 服务可达: ${r.status === 200 ? 'OK' : 'FAIL'} (HTTP ${r.status})`);
    if (r.status !== 200) return;
  } catch (e) {
    console.log(`[1/2] 服务可达: FAIL (${e.message})`);
    return;
  }

  // 2. 业务接口:获取登录二维码(安卓Pad 8.0.53 协议)
  try {
    const r = await fetch(`${API}/api/Login/GetQRPad`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ DeviceID: randDeviceId(), DeviceName: 'qinglong-bot' }),
    });
    const text = await r.text();
    let brief = text.slice(0, 300);
    try {
      const j = JSON.parse(text);
      brief = JSON.stringify(j).slice(0, 300);
    } catch {}
    console.log(`[2/2] GetQRPad 业务接口: HTTP ${r.status}`);
    console.log(`响应摘要: ${brief}`);
    console.log(r.ok ? '结论: 对接正常,可用扫码登录微信' : '结论: 接口异常,查看 wxapi 容器日志排查');
  } catch (e) {
    console.log(`[2/2] GetQRPad 业务接口: FAIL (${e.message})`);
  }
})();
