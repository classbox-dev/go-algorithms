#!/usr/bin/env ash
if [ "$(id -u)" = "0" ]; then
    echo "container cannot be run as root"
    exit 1
fi
export PATH=/opt/bin:$PATH
alias build="python3 /opt/bin/build"
exec "$@"
