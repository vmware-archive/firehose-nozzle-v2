#!/usr/bin/env bash
set -e

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd ${dir}

pushd ../gateway
zip -r nozzle.zip *
popd

mv ../gateway/nozzle.zip resources
tile build
