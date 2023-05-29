#!/usr/bin/env bash

set -e

current_dir=$(pwd)

cleanup() {
  echo "Cleaning up..."
  # Kill the godoc process if it's still running
  if [[ -n $godoc_pid ]]; then
    kill -9 "$godoc_pid" || true
    wait "$godoc_pid" 2>/dev/null || true
  fi
  echo "Cleanup complete."
}

trap cleanup EXIT

cd stdlib
godoc -http=:6060 -links=false -templates=../godocs/static &
godoc_pid=$!

url="http://127.0.0.1:6060"
for i in {1..10}; do
  response=$(curl --write-out "%{http_code}" --silent --output /dev/null "$url") || true
  if [ "$response" -ne "302" ]; then
    echo "[$i] No response, retrying in 1 second..."
    sleep 1
  else
    echo "godoc is up!"
    break
  fi
done

#sleep 200

docs_dir=$(mktemp -d)
cd "$docs_dir"
wget -r -np -N -nH --cut-dirs=3 -E -p -k -e robots=off http://127.0.0.1:6060/pkg/hsecode.com/stdlib/v2 || true

cd "$current_dir"
rm -rf docs
mv "$docs_dir" docs
