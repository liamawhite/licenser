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

var _ Handler = (*Generic)(nil)

// FromTemplateFile creates a new license handler that uses the given file
// as the template source.
func FromTemplateFile(templatePath string, markerText string, year int, owner string) *Generic {
	tmpl := template.Must(template.ParseFiles(templatePath))
	return &Generic{Template: tmpl, MarkerText: markerText, Year: year, Owner: owner}
}

// FromTemplateString creates a new license handler that uses the given template
// string to generate the license headers.
func FromTemplateString(templateStr string, markerText string, year int, owner string) *Generic {
	tmpl := template.Must(template.New("license").Parse(templateStr))
	return &Generic{Template: tmpl, MarkerText: markerText, Year: year, Owner: owner}
}

// Generic license handler that renders the license from the configured template
type Generic struct {
	Year       int
	Owner      string
	Template   *template.Template
	MarkerText string

	licenseCache []byte
}

// Reader returns a reader populated with the license file prefix
func (g *Generic) Reader() io.Reader {
	return bytes.NewReader(g.bytes())
}

// IsPresent verifies that the license is present in the reader passed.
func (g *Generic) IsPresent(in io.Reader) bool {
	inScanner := bufio.NewScanner(in)
	// Check for presence of license in first 20 lines
	for i := 0; i < 20; i++ {
		if inScanner.Scan() {
			// We should definitely be more thorough here but this will do for now
			if strings.Contains(inScanner.Text(), g.MarkerText) {
				return true
			}
		}
	}
	return false
}

func (g *Generic) bytes() []byte {
	if g.licenseCache != nil {
		return copyBytes(g.licenseCache)
	}
	b := bytes.NewBuffer([]byte{})
	_ = g.Template.Execute(b, g)
	g.licenseCache = b.Bytes()
	return copyBytes(g.licenseCache)
}

// copyBytes makes copies so consumers of this interface can't mess with our cache
func copyBytes(in []byte) []byte {
	tmp := make([]byte, len(in))
	copy(tmp, in)
	return tmp
}
