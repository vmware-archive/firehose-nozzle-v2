#!/usr/bin/env bash
set -e

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd ${dir}

export local_port=9000
export remote_port=8086
export remote_ip=X.X.X.X
export bind_host_host_port=${local_port}:${remote_ip}:${remote_port}

export CA_CERT_PATH=./certs/root_ca_certificate
export CERT_PATH=./certs/my.crt
export KEY_PATH=./certs/my.key

export LOGS_API_ADDR=localhost:${local_port}

trap cleanup INT
function cleanup() {
    echo "Killing tunnel"
    kill $(pgrep -f "$bind_host_host_port")
    exit 0
}

echo "Starting tunnel"
echo ""

ssh -N -f \
    -i $HOME/pcf.example.opsman_rsa \
    -L ${bind_host_host_port} \
    ubuntu@pcf.example.com

echo "Starting nozzle"
echo ""

go run main.go

