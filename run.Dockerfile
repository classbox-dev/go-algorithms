FROM ubuntu:18.04 as build
WORKDIR /
ADD sandbox/run/sources.txt /etc/apt/sources.list
RUN \
    apt-get update -qq && \
    apt-get install -y -q --no-upgrade --no-install-recommends \
        dpkg-dev \
        flex bison libbz2-dev libdw-dev libelf-dev systemtap-sdt-dev libaudit-dev \
        libssl-dev libslang2-dev libunwind-dev libiberty-dev binutils-dev \
        make build-essential \
        && \
    apt-get source linux-image-unsigned-4.15.0-50-generic && \
    apt-get clean
RUN \
    cd /linux-4.15.0/tools/perf && \
    make -j LDFLAGS=-static CFLAGS='-DNDEBUG -O3'

FROM alpine:3.11.2
COPY --from=build linux-4.15.0/tools/perf/perf /usr/local/bin/
ENV \
    USER_NAME=sandbox \
    USER_ID=2000 \
    SANDBOX_DIR="/sandbox" \
    TIMEOUT_RUN=10 \
    TIMEOUT_MEASURE=120 \
    MEMORY_RUN=524288 \
    MEMORY_MEASURE=524288
RUN mkdir -p /home/$USER_NAME ${SANDBOX_DIR} && \
    adduser -s /bin/sh -D -u ${USER_ID} $USER_NAME && \
    chown -R $USER_NAME:$USER_NAME ${SANDBOX_DIR}

WORKDIR ${SANDBOX_DIR}

COPY sandbox/run /opt/bin
RUN chmod +x /opt/bin/*
ENTRYPOINT ["/opt/bin/init.sh"]

USER ${USER_ID}
