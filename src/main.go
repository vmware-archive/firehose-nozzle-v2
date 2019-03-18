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

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/cloudfoundry-incubator/uaago"
)

var allSelectors = []*loggregator_v2.Selector{
	{
		Message: &loggregator_v2.Selector_Log{
			Log: &loggregator_v2.LogSelector{},
		},
	},
	{
		Message: &loggregator_v2.Selector_Counter{
			Counter: &loggregator_v2.CounterSelector{},
		},
	},
	{
		Message: &loggregator_v2.Selector_Gauge{
			Gauge: &loggregator_v2.GaugeSelector{},
		},
	},
	{
		Message: &loggregator_v2.Selector_Timer{
			Timer: &loggregator_v2.TimerSelector{},
		},
	},
	{
		Message: &loggregator_v2.Selector_Event{
			Event: &loggregator_v2.EventSelector{},
		},
	},
}

func newTLSConfig(caPath, certPath, keyPath, cn string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		ServerName:   cn,
		Certificates: []tls.Certificate{cert},
		//InsecureSkipVerify: true,
	}

	caCertBytes, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCertBytes); !ok {
		return nil, errors.New("cannot parse ca cert")
	}

	tlsConfig.RootCAs = caCertPool

	return tlsConfig, nil
}

func main() {
	if true {
		gatewayMain()
	} else {
		rlpMain()
	}
}

func gatewayMain() {
	uaaURI := "https://uaa.sys.cf.example.com"
	uaaClient, err := uaago.NewClient(uaaURI)
	if err != nil {
		panic(err)
	}
	token, expiresIn, err := uaaClient.GetAuthTokenWithExpiresIn("v2-nozzle-test", "captor1263_Winding", true)
	println(token)
	println(expiresIn)
	println(err)
	if err != nil {
		panic(err)
	}

	gatewayURI := "https://log-stream.sys.cf.example.com.com/v2/read?counter"
	transport := http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := http.Client{Transport: &transport}
	gatewayURL, err := url.Parse(gatewayURI)
	if err != nil {
		panic(err)
	}

	response, err := client.Do(&http.Request{
		Header: map[string][]string{
			"Authorization": {token},
		},
		URL: gatewayURL,
	})

	if err != nil {
		panic(err)
	}

	//todo: the following code is NOT a sane implementation, just committing WIP
	//todo: handle buffers in a sane way
	b := make([]byte, 1)
	payload := ""
	for {
		_, err := response.Body.Read(b)
		if err != nil {
			panic(err)
		}
		if b[0] == []byte("\n")[0] {
			println(payload)
			payload = ""
		} else {
			payload += string(b[0])
		}
	}
}

func rlpMain() {
	tlsConfig, err := newTLSConfig(
		os.Getenv("CA_CERT_PATH"), os.Getenv("CERT_PATH"),
		os.Getenv("KEY_PATH"), "reverselogproxy",
	)
	if err != nil {
		log.Fatal("Could not create TLS config", err)
	}

	loggr := log.New(os.Stderr, "[", log.LstdFlags)

	streamConnector := loggregator.NewEnvelopeStreamConnector(
		os.Getenv("LOGS_API_ADDR"),
		tlsConfig,
		loggregator.WithEnvelopeStreamLogger(loggr),
	)

	rx := streamConnector.Stream(context.Background(), &loggregator_v2.EgressBatchRequest{
		ShardId:   os.Getenv("SHARD_ID"),
		Selectors: allSelectors,
	})

	for {
		batch := rx()

		for _, e := range batch {
			log.Printf("%+v\n", e)
		}
	}
}
