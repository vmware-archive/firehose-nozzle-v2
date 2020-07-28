#!/usr/bin/env bash
set -e

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd ${dir}

pushd ../gateway
make boostrap

rm -f ../gateway-tile/resources/nozzle.zip
zip -r ../gateway-tile/resources/nozzle.zip *
popd

tile build
