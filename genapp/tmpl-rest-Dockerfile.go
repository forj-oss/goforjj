package main

const template_rest_dockerfile = `FROM alpine:latest

WORKDIR /src

COPY ca_certificates/* /usr/local/share/ca-certificates/

RUN apk update && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates --fresh && \
    rm -f /var/cache/apk/*tar.gz && \
    adduser devops devops -D

COPY {{ .Yaml.Name }} /bin/{{ .Yaml.Name }}

USER devops

ENTRYPOINT ["/bin/{{.Yaml.Name}}"]

CMD ["--help"]`
