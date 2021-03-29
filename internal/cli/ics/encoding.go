/*
Copyright 2021 The Knative Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ics

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Encode will encode a cloud event to ICS encoding form:
// Base64(zlib(minimal JSON)).
func Encode(ce cloudevents.Event) (string, error) {
	bb, err := json.Marshal(ce)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCouldntEncode, err)
	}
	var b bytes.Buffer
	encoder := base64.NewEncoder(base64.RawURLEncoding, &b)
	w := zlib.NewWriter(encoder)
	_, err = w.Write(bb)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCouldntEncode, err)
	}
	err = w.Close()
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCouldntEncode, err)
	}
	err = encoder.Close()
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCouldntEncode, err)
	}
	return b.String(), nil
}

// Decode will decode an event from ICS encoding.
func Decode(encoded string) (*cloudevents.Event, error) {
	r := strings.NewReader(encoded)
	decoder := base64.NewDecoder(base64.RawURLEncoding, r)
	reader, err := zlib.NewReader(decoder)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCouldntDecode, err)
	}
	bb, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCouldntDecode, err)
	}
	ce := &cloudevents.Event{}
	err = json.Unmarshal(bb, ce)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCouldntDecode, err)
	}
	return ce, nil
}
