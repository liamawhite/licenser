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

package file

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/liamawhite/licenser/pkg/license"
)

// New returns a new file Mutator
func New(license license.Handler) *Mutator {
	return &Mutator{license: license}
}

var _ Licenser = &Mutator{}

// Mutator mutates files
type Mutator struct {
	license license.Handler
}

// Apply the license to the path passed or print to stdout if dryRun
func (m *Mutator) Apply(path string, dryRun bool) bool {
	// If we can't detect language skip (return true)
	styled := m.styledLicense(path)
	if styled == nil {
		return true
	}
	contents := getFileContents(path)
	if contents == nil {
		return false
	}
	if !m.license.IsPresent(bytes.NewReader(contents)) {
		newContents := merge(styled, contents)
		if dryRun {
			fmt.Printf("%s\n", newContents)
		} else if err := os.WriteFile(path, newContents, 0644); err != nil { // nolint: gosec
			_, _ = fmt.Fprintf(os.Stderr, "error writing license to %v:%v", path, err)
		}
	}
	return true
}

// Verify returns true if the license is present in the file passed
func (m *Mutator) Verify(path string, _ bool) bool {
	contents := getFileContents(path)
	if contents == nil {
		return false
	}
	// If we can't detect language skip (return true)
	if style := identifyLanguageStyle(path); style == nil {
		return true
	}
	present := m.license.IsPresent(bytes.NewReader(contents))
	if !present {
		_, _ = fmt.Fprintf(os.Stderr, "license missing from %v\n", path)
	}
	return present
}

// this should probably be cached on a per language basis
func (m *Mutator) styledLicense(path string) []byte {
	style := identifyLanguageStyle(path)
	if style == nil {
		return nil
	}
	buf := bytes.NewBuffer([]byte{})

	// TODO: implement block styling
	if style.isBlock {

	} else {
		scanner := bufio.NewScanner(m.license.Reader())
		for scanner.Scan() {
			_, _ = buf.WriteString(style.comment)
			if len(scanner.Bytes()) != 0 {
				_, _ = buf.WriteString(" ")
			}
			_, _ = buf.Write(scanner.Bytes())
			_, _ = buf.WriteString("\n")
		}
	}
	return buf.Bytes()
}

// this is also pretty horrible but does the job
func merge(license, file []byte) []byte {
	result := bytes.NewBuffer([]byte{})
	fileScanner := bufio.NewScanner(bytes.NewReader(file))
	for fileScanner.Scan() {
		// If there's a #! preserve it
		if strings.Contains(fileScanner.Text(), "#!") {
			result.Write(fileScanner.Bytes())
			result.WriteString("\n\n")
			result.Write(license)
		} else {
			result.Write(license)
			result.WriteString("\n")
			result.Write(fileScanner.Bytes())
			result.WriteString("\n")
		}

		// Now that we've written the license just dump the rest
		for fileScanner.Scan() {
			result.Write(fileScanner.Bytes())
			result.WriteString("\n")
		}
	}
	return result.Bytes()
}

// This function has the potential to become an unwiedly mess, consider rethinking.
// TODO: Create a language interface that can be cycled through in order to identify the file as said language
// Interface should have a lightweight "looksLike" and then a more heavyweight "verify"
func identifyLanguageStyle(path string) *languageStyle {
	// This comparison is probably cheaper so do it first.
	if result := identifyFromExtension(filepath.Ext(path)); result != nil {
		return result
	}
	if match, _ := regexp.MatchString("\\..*rc", path); match {
		return commentStyles["shell"]
	}
	if match, _ := regexp.MatchString(".*Makefile$", path); match {
		return commentStyles["make"]
	}
	if match, _ := regexp.MatchString(".*Dockerfile(\\..*)?$", path); match {
		return commentStyles["docker"]
	}
	if match, _ := regexp.MatchString(".*BUILD(\\..*)?$|WORKSPACE(\\..*)?$", path); match {
		return commentStyles["bazel"]
	}
	_, _ = fmt.Fprintf(os.Stderr, "unable to identify language of %v\n", path)
	return nil
}

func identifyFromExtension(extension string) *languageStyle {
	style, ok := commonExtensions[extension]
	if !ok {
		return nil
	}
	return commentStyles[style]
}

func getFileContents(path string) []byte {
	// Reading an entire file into memory will be an issue with really large files...
	contents, err := os.ReadFile(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to read file %v\n", path)
	}
	return contents
}
