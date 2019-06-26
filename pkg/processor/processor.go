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

package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/liamawhite/licenser/pkg/file"
	mutator "github.com/liamawhite/licenser/pkg/file"
	"github.com/liamawhite/licenser/pkg/license"
)

type processor struct {
	startDirectory string
	dryRun         bool

	wg      sync.WaitGroup
	mutator *file.Mutator
	success bool

	visitFunc func(path string, dryRun bool) bool
}

func New(startDirectory string, license license.Handler) *processor {
	return &processor{
		startDirectory: startDirectory,
		mutator:        mutator.New(license),
	}
}

func (p *processor) Apply(recurse, dryRun bool) bool {
	p.dryRun = dryRun
	p.visitFunc = p.mutator.AppendLicense
	return p.run(recurse)
}

func (p *processor) Verify(recurse bool) bool {
	p.visitFunc = p.mutator.VerifyLicense
	return p.run(recurse)
}

func (p *processor) run(recurse bool) bool {
	p.success = true
	if recurse {
		filepath.Walk(p.startDirectory, p.visit)
	} else {
		f, err := os.Open(p.startDirectory)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening info for directory %v:%v", p.startDirectory, err)
			return false
		}
		fileList, err := f.Readdir(0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading contents of directory %v:%v", p.startDirectory, err)
			return false
		}
		for _, file := range fileList {
			p.visit(filepath.Join(p.startDirectory, file.Name()), file, nil)
		}
	}
	p.wg.Wait()
	return p.success
}

func (p *processor) visit(path string, f os.FileInfo, err error) error {
	if shouldSkip(path) {
		return nil
	}
	if f.Mode().IsRegular() {
		p.wg.Add(1)
		go func(path string) {
			if !p.visitFunc(path, p.dryRun) {
				p.success = false
			}
			p.wg.Done()
		}(path)
	}
	return nil
}

func shouldSkip(path string) bool {
	if strings.Contains(path, ".git") {
		return true
	}
	if strings.Contains(path, ".md") {
		return true
	}
	return false
}
