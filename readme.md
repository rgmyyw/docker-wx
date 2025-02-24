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