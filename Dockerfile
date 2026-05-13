FROM golang:1.22 AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 复制项目中的源代码
COPY . .

# 将我们的代码编译成二进制可执行文件 blueblog_app
RUN go build -o blueblog .

###################
# 接下来创建一个小镜像
###################
FROM debian:latest


COPY ./scripts/wait.sh /
COPY ./web/templates /web/templates
COPY ./web/static /web/static
COPY ./conf /conf

# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=builder /build/blueblog /

RUN set -eux; \
    apt-get update; \
    apt-get install -y \
    --no-install-recommends \
    netcat; \
    chmod 755 wait.sh

# 声明服务端口
EXPOSE 8888

