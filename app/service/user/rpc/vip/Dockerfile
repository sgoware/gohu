#获取golang
FROM golang:1.19.0-alpine3.16 as go

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct"

WORKDIR /server

# 将代码复制到容器中
#COPY etc /server/etc
COPY ./gohu-user-vip-rpc /server/gohu-user-vip-rpc

# 启动容器时运行的命令
CMD ["./gohu-user-vip-rpc"]
