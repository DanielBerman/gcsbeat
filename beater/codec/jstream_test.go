// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package codec

import (
	"strings"
	"testing"
)

func TestJsonStreamNextErr(t *testing.T) {
	cases := map[string]struct {
		Json      string
		Length    int
		ExpectErr bool
	}{
		"empty file": {
			Json:      "",
			Length:    0,
			ExpectErr: false,
		},
		"single element": {
			Json:      `{"foo":"bar"}`,
			Length:    1,
			ExpectErr: false,
		},
		"multi element": {
			Json:      `{"foo":"bar"}{"bar":"bazz"}`,
			Length:    2,
			ExpectErr: false,
		},
		"newline delimited": {
			Json: `{"foo":"bar"}
{"bar":"bazz"}`,
			Length:    2,
			ExpectErr: false,
		},
		"second corrupt": {
			Json:      `{"foo":"bar"}aaa`,
			Length:    1,
			ExpectErr: true,
		},
		"non object": {
			Json:      `"foo":`,
			Length:    0,
			ExpectErr: true,
		},
	}

	for tn, tc := range cases {
		reader := strings.NewReader(tc.Json)
		c := NewJsonStreamCodec("file/path", reader)

		if c.Err() != nil {
			t.Errorf("%q | Not expected to start with an error, got %q", tn, c.Err())
		}

		counter := 0
		for c.Next() {
			counter++
		}

		hasErr := c.Err() != nil
		if hasErr != tc.ExpectErr {
			t.Errorf("%q | Got error %d, expected? %v", tn, c.Err(), tc.ExpectErr)
		}

		if counter != tc.Length {
			t.Errorf("%q | Expected to decode %d objects, got %d", tn, tc.Length, counter)
		}

	}
}
