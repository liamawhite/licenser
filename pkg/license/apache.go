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
	"bufio"
	"bytes"
	"io"
	"strings"
	"text/template"
)

var _ Handler = &Apache20{}

// NewApache20 creates a new Apache 2.0 license handler
func NewApache20(year int, owner string) *Apache20 {
	return &Apache20{Year: year, Owner: owner}
}

// Apache20 is an Apache 2.0 license handler
type Apache20 struct {
	Year  int
	Owner string

	licenseCache []byte
}

// Reader returns a reader populated with the Apache 2.0 license file prefix
func (a *Apache20) Reader() io.Reader {
	return bytes.NewReader(a.bytes())
}

// IsPresent verifies that an Apache 2.0 license is present in the reader passed.
func (a *Apache20) IsPresent(in io.Reader) bool {
	inScanner := bufio.NewScanner(in)
	// Check for presence of license in first 20 lines
	for i := 0; i < 20; i++ {
		if inScanner.Scan() {
			// We should definitely be more thorough here but this will do for now
			if strings.Contains(inScanner.Text(), "Licensed under the Apache License, Version 2.0") {
				return true
			}
		}
	}
	return false
}

func (a *Apache20) bytes() []byte {
	if a.licenseCache != nil {
		return copyBytes(a.licenseCache)
	}
	tmpl, _ := template.New("apache20").Parse(apache20Template)
	b := bytes.NewBuffer([]byte{})
	_ = tmpl.Execute(b, a)
	a.licenseCache = b.Bytes()
	return copyBytes(a.licenseCache)
}

// make copies so consumers of this interface can't mess with our cache
func copyBytes(in []byte) []byte {
	tmp := make([]byte, len(in))
	copy(tmp, in)
	return tmp
}

const apache20Template = `Copyright {{.Year}} {{.Owner}}

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
`
