# 多阶段构建
# 第一阶段，391MB，编译前准备：go 和 git 工具、代码依赖库
FROM golang:1.12.4-alpine3.9  AS front
RUN set -xe && \
    apk add git && \
    go get -v github.com/go-redis/redis && \
    go get -v gopkg.in/yaml.v2

# 分成两次 RUN 目的是可复用上面的缓存，编译 go 代码
COPY . /go/src/github.com/wpxun/gorpcweb
RUN set -xe && \
    go install github.com/wpxun/gorpcweb


# 第二阶段，14.6MB；仅仅复制了可执行程序和程序的配置文件
FROM alpine:3.9
ENV GOM_VERSION   1904.1
COPY --from=front /go/bin /go/src/github.com/wpxun/gorpcweb/config.yaml /go/bin/
EXPOSE 80
WORKDIR /go/bin
CMD ["/go/bin/gorpcweb"]