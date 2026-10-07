#!/usr/bin/env python3
# 订阅级任务均匀重排 + 标签维护:python3 respread-sub.py <订阅名>
# 1) 该订阅的"每日任务"(M H * * *)等间隔铺满 06:00-18:00;周/月任务不动
# 2) smallfawn 订阅顺带重打标签(wxapp/sf-other),视图按标签过滤免疫订阅重建
import sqlite3, json, sys
from collections import Counter

name = sys.argv[1] if len(sys.argv) > 1 else None
H0, H1 = 6, 18
if not name:
    print('[respread] 用法: respread-sub.py <订阅名>'); sys.exit(0)

db = sqlite3.connect('/ql/data/db/database.sqlite')
sub = db.execute('select id from Subscriptions where name=?', (name,)).fetchone()
if not sub:
    print(f'[respread:{name}] 订阅不存在'); sys.exit(0)
rows = db.execute('select id, name, schedule, command, labels from Crontabs where sub_id=? and isDisabled=0', (sub[0],)).fetchall()

# 标签维护(仅 smallfawn):wxapp / sf-other
tagged = 0
if name == 'smallfawn':
    for tid, tname, sch, cmd, labels in rows:
        try: arr = json.loads(labels or '[]')
        except Exception: arr = []
        tag = 'wxapp' if 'wxapp' in (cmd or '') else 'sf-other'
        if tag not in arr:
            arr.append(tag); tagged += 1
            db.execute('update Crontabs set labels=? where id=?', (json.dumps(arr), tid))
    if tagged:
        db.commit()
        print(f'[respread:{name}] 补标签 {tagged} 个')

def is_daily(s):
    p = (s or '').split()
    return len(p) == 5 and p[2] == '*' and p[3] == '*' and p[4] == '*'

daily = [(i, n) for i, n, s, c, l in rows if is_daily(s)]
if not daily:
    print(f'[respread:{name}] 无每日任务,跳过'); sys.exit(0)

daily.sort(key=lambda x: x[1].lower())
span = (H1 - H0) * 60
step = span / len(daily)
for i, (tid, _) in enumerate(daily):
    t = int(i * step)
    db.execute('update Crontabs set schedule=? where id=?', (f"{t % 60} {H0 + t // 60} * * *", tid))
db.commit()
slots = [f"{H0 + int(i*step)//60}:{int(i*step)%60:02d}" for i in range(len(daily))]
print(f'[respread:{name}] 重排 {len(daily)} 个每日任务: {slots[0]}~{slots[-1]} 同分钟最大 {max(Counter(slots).values())}')
