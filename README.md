# firehose-nozzle-v2

## Connecting

There are two ways to connect get data from the v2 API.

* Directly connecting is doing via is a component called
[reverse log proxy](https://github.com/cloudfoundry/loggregator-release/tree/develop/jobs/reverse_log_proxy).
* The [reverse log proxy gateway](https://github.com/cloudfoundry/loggregator/blob/master/docs/rlp_gateway.md)
component provides an HTTP based way to connect to this API.

## Using RLP Gateway

The [RLP Gateway](https://github.com/cloudfoundry/loggregator/blob/master/docs/rlp_gateway.md)
makes v2 nozzle development much easier by:
* Adding back the ability to deploy applications as a nozzle
* Eliminating the need for mTLS
* Released [in PCF version 2.4](https://docs.pivotal.io/pivotalcf/2-4/pcf-release-notes/runtime-rn.html#-loggregator-v2-api-is-readable-through-rlp-gateway)
    
## Directly Connecting to RLP
### Authenticating

Authentication is done via mTLS. The certificate should be generated and signed by 
the [Ops Manager Root CA](https://docs.pivotal.io/pivotalcf/2-4/security/pcf-infrastructure/api-cert-rotation.html#-certificate-types).
Communication is done directly to the RLP over HTTP/2.
This endpoint is discoverable via
[its BOSH LINK](https://github.com/cloudfoundry/loggregator-release/blob/v105.1/jobs/reverse_log_proxy/spec#L21-L25)

### Connecting

The link is shared in a PCF deployment:
```yaml
...
provides: |
  reverse_log_proxy: {as: reverse_log_proxy, shared: true}
...
```

* In a full runtime deployment, this component is on the `loggregator_trafficcontroller` vm and listens on `:8082`
* In the small footprint PAS, it is deployed on the `control` vm and listens on `:8086`

Due to the protocol, this component is not exposed via go-router
(meaning that v2 nozzles need to use the gateway to be pushed as apps).

For development, the nozzle author can generate a certificate signed by OpsMan's root CA
[using its certificate API](https://docs.pivotal.io/pivotalcf/2-4/opsman-api/#certificates)
and the [om tool](https://github.com/pivotal-cf/om)
```
om -t https://pcf.example.com -k -u {user} -p ${pwd} \
    curl -x POST -p /api/v0/certificates/generate -d '{ "domains": ["*.example.com", "*.sub.example.com"] }'
```

The domain used in the certificate isn't important.

Downloading the CA cert
Admin -> Settings -> Advanced -> Download Root CA Cert

As mentioned above, the communication uses http/2.
In a PCF environment (rather than [cfdev](https://github.com/cloudfoundry-incubator/cfdev)
or some other tooling where the component would be directly accessible), 
on way to develop locally is to setup a ssh tunnel through OpsMan:
```
ssh -i [path to private key] -L 9000:[IP of loggregator_trafficcontroller or control VM]:[8082 or 8086] ubuntu@opsman.example.com
```

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

ssh -i [path to private key] \
  -L 9000:[IP of loggregator_trafficcontroller or control VM]:[8082 or 8086] \
  ubuntu@opsman.example.com
```

## With RLP Gateway

The [RLP Gateway](https://github.com/cloudfoundry/loggregator-release/tree/develop/jobs/reverse_log_proxy_gateway) makes
v2 nozzle development much easier by:
* Addding back the ability to deploy applications as a nozzle
* Eliminating the need for mTLS

Coming in PCF version **???**

**todo: detailed technical notes**
