// Copyright 2019 Liam White
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package license

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApache20_Reader(t *testing.T) {
	t.Run("Reader has correct bytes", func(t *testing.T) {
		a := NewApache20(2019, "Test")
		want, _ := ioutil.ReadFile("testdata/apache.golden")
		got, _ := ioutil.ReadAll(a.Reader())
		assert.Equal(t, want, got)
	})
}

func TestApache20_IsPresent(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		want      bool
	}{
		{
			name:      "License is present",
			inputFile: "testdata/apache.golden",
			want:      true,
		},
		{
			name:      "License is not present",
			inputFile: "testdata/nolicense.golden",
			want:      false,
		},
	}
	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			a := &Apache20{}
			inputReader, _ := os.Open(tc.inputFile)
			assert.Equal(t, tc.want, a.IsPresent(inputReader))
		})
	}
}
