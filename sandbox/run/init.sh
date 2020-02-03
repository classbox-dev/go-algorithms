#!/usr/bin/env ash
if [ "$(id -u)" = "0" ]; then
    echo "container cannot be run as root"
    exit 1
fi
export PATH=/opt/bin:/in:$PATH
ulimit -v 524288
timeout -s KILL $TIMEOUT "$@"
