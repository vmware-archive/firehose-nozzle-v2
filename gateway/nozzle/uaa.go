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

import "github.com/cloudfoundry-incubator/uaago"

//go:generate counterfeiter ./ UAA
type UAA interface {
	GetAuthToken() (string, error)
}

type uaa struct {
	uaaUser           string
	uaaPass           string
	skipSSLValidation bool
	uaaClient         *uaago.Client
}

func NewUAA(uaaURL string, uaaUser string, uaaPass string, skipSSLValidation bool) (UAA, error) {
	uaaClient, err := uaago.NewClient(uaaURL)
	if err != nil {
		return nil, err
	}

	return &uaa{
		uaaUser:           uaaUser,
		uaaPass:           uaaPass,
		skipSSLValidation: skipSSLValidation,
		uaaClient:         uaaClient,
	}, nil
}

func (uaa *uaa) GetAuthToken() (string, error) {
	token, _, err := uaa.uaaClient.GetAuthTokenWithExpiresIn(uaa.uaaUser, uaa.uaaPass, uaa.skipSSLValidation)
	return token, err
}
