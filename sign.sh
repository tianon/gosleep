#!/bin/bash
set -e

cd "$(dirname "$(readlink -f "$BASH_SOURCE")")/artifacts"

set -x
rm -f *.asc
for f in gosleep*; do
	gpg --output "$f.asc" --detach-sign "$f"
done
sha256sum gosleep* > SHA256SUMS
gpg --output SHA256SUMS.asc --detach-sign SHA256SUMS
ls -lAFh gosleep* SHA256SUMS*
