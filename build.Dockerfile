FROM golang:1.13.5-alpine3.11
LABEL maintainer="Max Kuznetsov <maks.kuznetsov@gmail.com>"
ARG USER_ID
ARG GROUP_ID
ENV \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE="on" \
    USER_NAME=sandbox
RUN \
    : "${USER_ID:?UID of the local unprivileged user must be passed via --build-arg}" && \
    : "${GROUP_ID:?GID of the local unprivileged group must be passed via --build-arg}" && \
    apk add --no-cache --update su-exec curl ca-certificates python3 wget git && \
    mkdir -p /home/$USER_NAME && \
    adduser -s /bin/sh -D -u ${USER_ID} -g ${GROUP_ID} $USER_NAME && \
    mkdir -p /in /out /opt/bin && \
    go get golang.org/x/tools/cmd/godoc@v0.0.0-20200110142700-428f1ab0ca03 && \
    rm -rf /var/cache/apk/*
WORKDIR /

COPY godocs /opt
COPY build /opt/bin
RUN \
    patch -N $(go env GOROOT)/src/fmt/print.go /opt/noprint.patch && \
    chmod -R +x /opt/bin && \
    ln -s /opt/bin/build.py /opt/bin/build
ENTRYPOINT ["/opt/bin/init.sh"]

COPY stdlib-linter /stdlib-linter
RUN cd /stdlib-linter && go install

COPY stdlib-tests /stdlib-tests
USER ${USER_ID}:${GROUP_ID}
