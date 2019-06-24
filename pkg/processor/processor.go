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
	recurse        bool
	wg             sync.WaitGroup
	mutator        *file.Mutator
	dryRun         bool
}

func New(startDirectory string, recurse bool, license license.Handler) *processor {
	return &processor{
		startDirectory: startDirectory,
		recurse:        recurse,
		mutator:        mutator.New(license),
	}
}

func (p *processor) Run(dryRun bool) error {
	p.dryRun = dryRun
	if p.recurse {
		filepath.Walk(p.startDirectory, p.visit)
	} else {
		f, err := os.Open(p.startDirectory)
		if err != nil {
			return err
		}
		fileList, err := f.Readdir(0)
		if err != nil {
			return err
		}
		for _, file := range fileList {
			p.visit(filepath.Join(p.startDirectory, file.Name()), file, nil)
		}
	}
	p.wg.Wait()
	return nil
}

func (p *processor) visit(path string, f os.FileInfo, err error) error {
	if shouldSkip(path) {
		return nil
	}
	if !f.IsDir() {
		p.wg.Add(1)
		go func(path string) {
			p.mutator.AppendLicense(path, p.dryRun)
			p.wg.Done()
		}(path)
	}
	return nil
}

func shouldSkip(path string) bool {
	if strings.Contains(path, ".git") {
		return true
	}
	return false
}
