// firehose-nozzle-v2
//
// Copyright (c) 2018-Present Pivotal Software, Inc. All Rights Reserved.
//
// This program and the accompanying materials are made available under the terms of the under the Apache License,
// Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may
// obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing permissions and
// limitations under the License.

package nozzle

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)


func Stream(ctx context.Context, c *Config, uaaClient UAA, shipper LogShipper) error {
	for ctx.Err() == nil {
		Receive(ctx, c, uaaClient, shipper)
	}
	return ctx.Err()
}

func Receive(ctx context.Context, c *Config, uaaClient UAA, shipper LogShipper) {
	token, err := uaaClient.GetAuthToken()
	if err != nil {
		fmt.Print(err)
		return
	}

	gatewayURI := c.LogStreamUrl + "/v2/read?" + strings.Join(c.Envelopes, "&")
	transport := http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: c.SkipSSLValidation}}
	client := http.Client{Transport: &transport}
	gatewayURL, err := url.Parse(gatewayURI)
	if err != nil {
		fmt.Print(err)
		return
	}

	response, err := client.Do(&http.Request{
		Header: map[string][]string{
			"Authorization": {token},
		},
		URL: gatewayURL,
	})
	if err != nil {
		fmt.Print(err)
		return
	}

	reader := bufio.NewReader(response.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print(err)
			return
		}

		line = strings.TrimSpace(line)
		select {
		case <-ctx.Done():
			return
		default:
			if len(line) > 0 {
				err = shipper.LogShip(line)
				if err != nil {
					fmt.Print(err)
					return
				}
			}
		}
	}
}
