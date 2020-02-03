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
ENV LOGIN=sandbox UID=2000 TIMEOUT=60
RUN mkdir -p "/in" && \
    adduser -s /bin/sh -D -u ${UID} $LOGIN && \
    chown -R $LOGIN:$LOGIN "/in"
WORKDIR "/in"
COPY sandbox/run/init.sh /opt/bin/init.sh
RUN chmod +x /opt/bin/init.sh
ENTRYPOINT ["/opt/bin/init.sh"]
USER ${UID}
