# firehose-nozzle-v2

## Connecting

There are two ways to connect get data from the v2 API.

* [Reverse Log Proxy Gateway](https://github.com/cloudfoundry/loggregator/blob/master/docs/rlp_gateway.md)
* [Reverse Log Proxy](https://github.com/cloudfoundry/loggregator-release/tree/develop/jobs/reverse_log_proxy)

Each is described below. Using the Gateway is the easier path to building a nozzle.

## Building a Nozzle Directly Connecting to RLP (Recommended Path)

Communication is done directly to the RLP over HTTP/2.
This endpoint is discoverable via
[its BOSH LINK](https://github.com/cloudfoundry/loggregator-release/blob/v105.1/jobs/reverse_log_proxy/spec#L21-L25)

The link is shared in a PCF deployment:
```yaml
...
provides: |
  reverse_log_proxy: {as: reverse_log_proxy, shared: true}
...
```

* In a full runtime deployment, this component is on the `loggregator_trafficcontroller` vm and listens on `:8082`
* In the small footprint PAS, it is deployed on the `control` vm and listens on `:8086`

Authentication is done via mTLS. The mTLS connection is authenticated by connecting with a certificate signed by 
the [Ops Manager Root CA](https://docs.pivotal.io/pivotalcf/2-4/security/pcf-infrastructure/api-cert-rotation.html#-certificate-types).

For development, the nozzle author can manually generate a certificate signed by OpsMan's root CA
[using its certificate API](https://docs.pivotal.io/pivotalcf/2-4/opsman-api/#certificates)
and the [om tool](https://github.com/pivotal-cf/om):

```
om -t https://pcf.example.com -k -u ${user} -p ${pwd} \
    curl -x POST -p /api/v0/certificates/generate -d '{ "domains": ["*.example.com", "*.sub.example.com"] }'
```

The domain used in the certificate does _not_ matter.

To get the
[root certificate from Ops Manager](https://docs.pivotal.io/pivotalcf/security/pcf-infrastructure/api-cert-rotation.html#-certificate-types),
download in advanced settings:
Admin -> Settings -> Advanced -> Download Root CA Cert

As mentioned above, the communication uses http/2.
In a PCF environment (rather than [cfdev](https://github.com/cloudfoundry-incubator/cfdev)
or some other tooling where the component would be directly accessible), 
one way to develop locally is to setup a ssh tunnel through OpsMan:

```
ssh -i [path to ssh private key] \
  -L 9000:[IP of loggregator_trafficcontroller or control VM]:[8082 or 8086] \
  ubuntu@opsman.example.com
```

`src/local_dev.template.sh` is a sample script that will run the nozzle, once 
the certificates are generated and put on disk.

## Building a Nozzle Using RLP Gateway

The RLP Gateway adds:
* Eliminating the need for mTLS
* Does not require any Loggregator specific libraries to process data

The gateway was released [in PCF version 2.4](https://docs.pivotal.io/pivotalcf/2-4/pcf-release-notes/runtime-rn.html#-loggregator-v2-api-is-readable-through-rlp-gateway)

The code in this repo is for bootstraping, but the client here:
https://github.com/cloudfoundry/go-loggregator/blob/master/rlp_gateway_client.go
should be the basis for a production nozzle.

### Authentication & Testing
To create a UAA user that can access the data, use
the [UAA CLI](https://docs.cloudfoundry.org/uaa/uaa-user-management.html).

Create the user:

```bash
uaac target https://uaa.sys.<pcf system domain> --skip-ssl-validation
uaac token client get admin -s <admin client secret>
uaac client add my-v2-nozzle \
  --name my-v2-nozzle \
  --secret <my-v2-nozzle client secret> \
  --authorized_grant_types client_credentials,refresh_token \
  --authorities logs.admin
```

To manually get a token:

```bash
uaac token client get v2-nozzle-test -s <my-v2-nozzle client secret>
uaac context
``` 

The RLP Gateway data can be tested with just `curl`. To view the 
data (newline delimited JSON payloads), copy the token and run:
```bash
export token=<my-v2-nozzle token from context>
curl -k -H "Authorization: $token" 'https://log-stream.sys.<pcf system domain>/v2/read?counter&gauge'
```

## Envelope Tags

Converting back to V1 envelopes may be desirable in in some cases. There are notes about it in the v2 -> v1 mapping doc below.
One important differences is the type of tags used in envelopes. In v1 tags have types while in v2 they are strings. If you want
to use v2 envelopes rather then converting back to v1 envelopes, you should add the flag `UsePreferredTags` when connecting to the
RLP(this is added automatically when connecting to the rlp gateway). In this case tags will be stored in tags as v2 style strings.
However, if you wish to convert the envelopes to v1 envelopes, you should not add this flag, 
tags will then be stored in `deprecated_tags` as `loggregator_v2.Value` items that will be converted when calling `toV1`.

## Tile
The `tile` directory packages the gateway version of the nozzle as an app deployed on the platform.

Build is done via [PCF Tile Generator](https://github.com/cf-platform-eng/tile-generator/)

## References

* v2 -> v1 mapping: https://github.com/cloudfoundry/loggregator-api/blob/master/README.md#v2---v1-mapping
* v2 reference example https://github.com/cloudfoundry-incubator/refnozzle
* Envelope proto buff def https://github.com/cloudfoundry/loggregator-api/blob/master/v2/envelope.proto
* Example: https://github.com/cloudfoundry/go-loggregator/blob/master/examples/envelope_stream_connector/main.go
* CLI plugin to stream v2 data https://github.com/cloudfoundry/log-cache-cli
* https://github.com/cloudfoundry/cf-drain-cli
