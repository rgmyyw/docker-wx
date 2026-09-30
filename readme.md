# go docker服务

# 远程调试

1. 复制代码到服务器;
2. 服务器进入代码目录并执行dlv;
```bash
cd /root/go_project/wic-go/
dlv debug --headless --listen=:2346 --api-version=2
```
3.客户端运行go-remote远程调试

# http专用

cd /home/docker-go
docker build -f ./Dockerfile.http -t wic-go-http:1.0.1 .
docker service remove wic-go-http
docker service create --replicas 16 --network wic-business --name wic-go-http -p 8005:8005 wic-go-http:1.0.1

# tcp专用

cd /home/docker-go
docker build -f ./Dockerfile.tcp -t wic-go-tcp:1.0.1 .
docker service remove wic-go-tcp
docker service create --replicas 1 --network wic-business --name wic-go-tcp -p 8006:8006 wic-go-tcp:1.0.1


go env -w GOOS=linux GOARCH=amd64
go build -o wx main.go

# Docker Compose 一键部署(fork 增强)

相比上游原始 Dockerfile,本部署方式修复了若干问题并自带配套服务:

- 修复 `Dockerfile.http`:`CMD ["./main"]` 实际不存在,改为运行仓库自带的 `amd64` 预编译二进制(含执行位与上海时区);
- 自建 Redis:上游 `conf/app.conf` 默认指向作者外网 Redis(登录会话存到别人机器),改为 compose 内置 `wx-redis`,密码自定;
- `conf/app.conf` 以挂载方式覆盖进容器,改配置重启即生效,无需重建镜像;
- 内置扫码登录台 `wxlogin`(端口 **8058**):打开网页即见二维码,自动轮询、过期自动刷新、保存 62 凭证支持免扫码二次登录;
- `.dockerignore` 排除无关大文件(`.git` / `*.zip` / `*.exe` 等),镜像瘦身约 45MB。

## 部署步骤

```bash
cp .env.example .env                 # 编辑 REDIS_PASSWORD
sed -i "s/CHANGE_ME_REDIS_PASSWORD/<你的密码>/" conf/app.conf
docker compose up -d --build
```

服务一览:

| 服务 | 端口 | 说明 |
|---|---|---|
| wxapi | 8057 | 微信协议 API,接口带 `/api` 前缀,根路径为 Swagger UI |
| wxlogin | 8058 | 扫码登录台(`/` 页面、`/qr` 二维码图直链、`/status` JSON) |
| wx-redis | 仅内网 | 会话存储,appendonly 持久化 |

## 对接青龙面板

wxapi 与 qinglong 同机部署时,可将 wxapi 加入青龙网络(compose 中已预留 `qinglong_net` external 网络,按需修改),脚本内直接访问 `http://wxapi:8057`。示例脚本见 `scripts/qinglong/`(`wxapi_check.js` 链路自检、`login.js` 命令行扫码登录)。
