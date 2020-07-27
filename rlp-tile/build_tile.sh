#!/usr/bin/env bash
set -e

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd ${dir}

pushd ../rlp
make boostrap

rm -f ../tile/resources/nozzle.zip
zip -r ../tile/resources/nozzle.zip *
popd

tile build
