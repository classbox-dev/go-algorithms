#!/usr/bin/env ash
if [ "$(id -u)" = "0" ]; then
    echo "container cannot be run as root"
    exit 1
fi
export PATH=/opt/bin:/sandbox:$PATH
ulimit -v "$MEMORY_RUN"
timeout -s KILL "$TIMEOUT_RUN" "$@"
