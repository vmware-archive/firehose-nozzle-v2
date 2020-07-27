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
	"gateway/nozzle"
	"os"
)

func main() {
	c, err := nozzle.GetConfig()
	if err != nil {
		panic(err)
	}

	uaaClient, err := nozzle.NewUAA(c.UAAURL, c.UAAUser, c.UAAPass, true)
	if err != nil {
		panic(err)
	}

	//shipper := nozzle.NewSampleShipper(os.Stdout)
	shipper := nozzle.NewParsingShipper(os.Stdout)
	err = nozzle.Receive(c, uaaClient, shipper)
	if err != nil {
		panic(err)
	}
}
