# firehose-nozzle-v2

## References

* v1 -> v2 mapping: https://github.com/cloudfoundry/loggregator-api/blob/master/README.md#v2---v1-mapping
* Envelope proto buff def https://github.com/cloudfoundry/loggregator-api/blob/master/v2/envelope.proto
* Example: https://github.com/cloudfoundry-incubator/refnozzle/blob/master/cmd/nozzle/main.go
* Example: https://github.com/cloudfoundry/go-loggregator/blob/master/examples/envelope_stream_connector/main.go
* https://github.com/cloudfoundry/log-cache-cli
* https://github.com/cloudfoundry/cf-drain-cli

## Reverse log proxy

The v2 egress API is a component called [reverse log proxy](https://github.com/cloudfoundry/loggregator-release/tree/develop/jobs/reverse_log_proxy).

Use the BOSH link here to get the IP/port, the link is shared:
```yaml
...
provides: |
  reverse_log_proxy: {as: reverse_log_proxy, shared: true}
...
```

* Full ERT: `loggregator_trafficcontroller` vm on :8082
* Small run time: `control` vm on :8086

## Without RLP Gateway
### Authenticating

via mTLS

Cert signed by OpsMan:
```
om -t https://pcf.example.com -k -u admin -p ${pwd} \
    curl -x POST -p /api/v0/certificates/generate -d '{ "domains": ["*.example.com", "*.sub.example.com"] }'
```

OpsMan CA cert:
Admin -> Settings -> Advanced -> Download Root CA Cert

### Local dev
The communication uses http/2, which Gorouter doesn't support. You have to be on-network;
is *not* accessible to applications running on the platform (requring packaging as a BOSH release)

For local dev, it's possible to setup a ssh tunnel through OpsMan:
```

ssh -i [path to private key] -L 9000:[IP of loggregator_trafficcontroller or control VM]:[8082 or 8086] ubuntu@opsman.example.com
```

## With RLP Gateway

The [RLP Gateway](https://github.com/cloudfoundry/loggregator-release/tree/develop/jobs/reverse_log_proxy_gateway) makes
v2 nozzle development much easier by:
* Addding back the ability to deploy applications as a nozzle
* Eliminating the need for mTLS

Coming in PCF version **???**

**todo: detailed technical notes**
