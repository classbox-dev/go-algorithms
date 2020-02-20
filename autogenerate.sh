#!/usr/bin/env bash
set -x

inotifywait -e close_write -m $1 |
while read -r directory events filename; do
  go generate ./...
done
