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
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"regexp"
	"strconv"
	"time"
)

const appID = "1706934d-a53b-4d60-bf19-16e69702f423"
const timeout = 30 * time.Second

type Batch struct {
	Batch []Entry `json:"batch"`
}

type Entry struct {
	Tags     map[string]interface{} `json:"tags"`
	EntryLog EntryLog               `json:"log"`
}

type EntryLog struct {
	Payload string `json:"payload"`
	Type    string `json:"type"`
}

type ParsingShipper struct {
	writer io.Writer
	t      *time.Timer
	buffer []int
	count  int
}

func NewParsingShipper(writer io.Writer) LogShipper {
	s := &ParsingShipper{
		writer: writer,
		t:      time.NewTimer(timeout),
	}
	go func() {
		<-s.t.C
		log.Fatal("No messages received before timeout")
	}()
	go func() {
		for range time.NewTicker(30 * time.Second).C {
			log.Printf("Checking in. Messages received: [%v]. Sequence %+v", s.count, s.buffer)
		}
	}()
	return s
}

func (s *ParsingShipper) LogShip(line string) error {
	jsonLog := line[5:] //string off "data:"

	var batch Batch
	err := json.Unmarshal([]byte(jsonLog), &batch)
	if err != nil {
		log.Printf("Error parsing: %s", jsonLog)
		log.Printf("%s", err.Error())
	}

	for _, entry := range batch.Batch {
		if entry.Tags["app_id"] == appID {
			err := s.handleEntry(entry)
			if err != nil {
				log.Printf("%s", err.Error())
			}
		}
	}

	return nil
}

func (s *ParsingShipper) handleEntry(line Entry) error {
	s.count++
	if line.EntryLog.Payload != "" {
		decoded, err := base64.StdEncoding.DecodeString(line.EntryLog.Payload)
		if err != nil {
			return err
		}

		r := regexp.MustCompile(`.*GENERATOR.*Identifier: .*\[(.*)\].*\[(.*)\]`)
		if r.Match(decoded) {

			capture := r.FindStringSubmatch(string(decoded))
			seq, err := strconv.Atoi(capture[2])
			if err != nil {
				return err
			}
			if seq == 0 {
				log.Printf("Sequence reset to 0")
				s.buffer = []int{seq}
			} else {
				s.buffer = append([]int{seq}, s.buffer[0:min(len(s.buffer), 5)]...)
				if len(s.buffer) > 1 {
					if s.buffer[0]-s.buffer[1] != 1 {
						log.Printf(">>>>> Missed item in sequence: %+v, %+v", seq, s.buffer)
					}
				}
			}

			s.t.Reset(timeout)
		}
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
