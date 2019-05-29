# 多阶段构建
# 第一阶段，360MB，编译前准备：go ，没有代码依赖库
FROM golang:1.12.4-alpine3.9  AS front

# 分成两次 RUN 目的是可复用上面的缓存，编译 go 代码
COPY . /go/src/github.com/wpxun/gorpcsvc
RUN set -xe && \
    go install github.com/wpxun/gorpcsvc


# 第二阶段，15.8MB；仅仅复制了可执行程序和程序的配置文件
FROM alpine:3.9
ENV GOM_VERSION   1904.1
COPY --from=front /go/bin /go/bin/
EXPOSE 1234
WORKDIR /go/bin
CMD ["/go/bin/gorpcsvc"]