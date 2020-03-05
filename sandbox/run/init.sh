#!/usr/bin/env ash
if [ "$(id -u)" = "0" ]; then
    echo "container cannot be run as root"
    exit 1
fi
export PATH=/opt/bin:/in:$PATH
ulimit -v 786432
timeout -s KILL $TIMEOUT "$@"
