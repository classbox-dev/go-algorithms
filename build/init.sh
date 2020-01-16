#!/usr/bin/env sh
if [ "$(id -u)" = "0" ]; then
    echo "container cannot be run as root"
    exit 1
fi
export PATH=/opt/bin:$PATH
exec "$@"
