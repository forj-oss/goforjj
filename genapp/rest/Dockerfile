FROM alpine:latest

WORKDIR /src

COPY ca_certificates/* /usr/local/share/ca-certificates/

RUN apk update && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates --fresh && \
    rm -f /var/cache/apk/*tar.gz && \
    adduser devops devops -D

COPY __MYPLUGINNAME__ /bin/__MYPLUGINNAME__

USER devops

ENTRYPOINT ["/bin/__MYPLUGINNAME__"]

CMD ["--help"]
