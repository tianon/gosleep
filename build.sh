#!/bin/bash
set -e

cd "$(dirname "$(readlink -f "$BASH_SOURCE")")"

set -x

docker build --pull -t gosleep:build -f Dockerfile.build .

rm -rf artifacts
docker run --rm gosleep:build tar -cC / artifacts | tar -xv
cd artifacts
sha256sum * | tee SHA256SUMS
file *
ls -lAFh

"./gosleep-$(dpkg --print-architecture | awk -F- '{ print $NF }')" --for 1s
