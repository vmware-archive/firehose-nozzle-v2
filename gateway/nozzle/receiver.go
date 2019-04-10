package nozzle

import (
	"bufio"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
)

func Receive(c *Config, uaaClient UAA, shipper LogShipper) error {
	token, err := uaaClient.GetAuthToken()
	if err != nil {
		return err
	}

	gatewayURI := c.LogStreamUrl + "/v2/read?counter"
	transport := http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: c.SkipSSLValidation}}
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

	reader := bufio.NewReader(response.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		line = strings.TrimSpace(line)
		if len(line) > 0 {
			shipper.LogShip(line)
		}
	}
}
