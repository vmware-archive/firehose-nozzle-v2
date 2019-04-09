package nozzle

import (
	"crypto/tls"
	"net/http"
	"net/url"
)

func GatewayMain(c *Config, uaaClient UAA, shipper LogShipper) error {
	token, err := uaaClient.GetAuthToken()
	if err != nil {
		panic(err)
	}

	gatewayURI := c.LogStreamUrl + "/v2/read?counter"
	transport := http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := http.Client{Transport: &transport}
	gatewayURL, err := url.Parse(gatewayURI)
	if err != nil {
		return err
	}

	response, err := client.Do(&http.Request{
		Header: map[string][]string{
			"Authorization": {token},
		},
		URL: gatewayURL,
	})

	if err != nil {
		return err
	}

	//todo: the following code is NOT a sane implementation, just committing WIP
	//todo: handle buffers in a sane way
	b := make([]byte, 1)
	payload := ""
	for {
		_, err := response.Body.Read(b)
		if err != nil {
			return err
		}
		if b[0] == []byte("\n")[0] {
			shipper.LogShip(payload)
			payload = ""
		} else {
			payload += string(b[0])
		}
	}

	return nil
}
