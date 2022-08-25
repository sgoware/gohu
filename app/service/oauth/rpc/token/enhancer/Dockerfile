#获取golang
FROM golang:1.19.0-alpine3.16 as go

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOPROXY="https://goproxy.cn,direct"
ENV APOLLO_APP_ID=$APOLLO_APP_ID
ENV	APOLLO_CLUSTER_NAME=$APOLLO_CLUSTER_NAME
ENV	APOLLO_IP=$APOLLO_IP
ENV	APOLLO_SECRET=$APOLLO_SECRET

WORKDIR /server

# 将代码复制到容器中
COPY ./gohu-oauth-token-enhancer-rpc /server/gohu-oauth-token-enhancer-rpc

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# 启动容器时运行的命令
CMD ["./gohu-oauth-token-enhancer-rpc"]
