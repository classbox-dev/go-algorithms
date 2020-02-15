FROM golang:1.13.5-alpine3.11
LABEL maintainer="Max Kuznetsov <maks.kuznetsov@gmail.com>"
ENV \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE="on" \
    LOGIN=sandbox \
    UID=2000
RUN \
    apk add --no-cache --update curl ca-certificates python3 py3-yaml py3-requests wget git && \
    adduser -s /bin/sh -D -u $UID $LOGIN && \
    mkdir -p /out /opt/bin && \
    chown -R $LOGIN:$LOGIN /out
RUN \
    go get golang.org/x/tools/cmd/godoc@v0.0.0-20200110142700-428f1ab0ca03 && \
    go get github.com/mkuznets/stdlib-linter@v0.3.1 && \
    go get github.com/cheekybits/genny@v1.0.0
WORKDIR /

COPY godocs /opt
COPY sandbox/build /opt/bin
RUN \
    patch -N $(go env GOROOT)/src/fmt/print.go /opt/noprint.patch && \
    chmod -R +x /opt/bin && \
    ln -s /opt/bin/build.py /opt/bin/build
ENTRYPOINT ["/opt/bin/init.sh"]

COPY stdlib-tests /stdlib-tests
COPY stdlib /stdlib
COPY linter_config.yaml /linter_config.yaml
USER ${UID}
