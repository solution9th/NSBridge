FROM alpine:3.9

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk update && apk add tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \ 
    && echo "Asia/Shanghai" > /etc/timezone

RUN apk add --update ca-certificates && rm -rf /var/cache/apk/*

COPY ./ns_bridge /main
COPY ./config.yml /config.yml

EXPOSE 8080

ENTRYPOINT ["/main"]