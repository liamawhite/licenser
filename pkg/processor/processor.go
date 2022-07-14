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

	"github.com/denormal/go-gitignore"

	"github.com/liamawhite/licenser/pkg/file"
	mutator "github.com/liamawhite/licenser/pkg/file"
	"github.com/liamawhite/licenser/pkg/license"
)

// licenserignoreFile is a name for *ignore files specific to licenser.
const licenserignoreFile = ".licenserignore"

// Processor finds all valid files and passes them to a file mutator to be handled
type Processor struct {
	startDirectory string

	mutator   file.Licenser
	visitFunc func(path string, dryRun bool) bool
	wg        sync.WaitGroup

	dryRun  bool
	success bool

	skipListGitIgnore      gitignore.GitIgnore
	skipListLicenserIgnore gitignore.GitIgnore
	skipListExtension      map[string]bool
}

// New creates a new file processor starting the the passed startDirectory
// and using the passed license to apply and verify files
func New(startDirectory string, license license.Handler) *Processor {
	return &Processor{
		startDirectory:         startDirectory,
		mutator:                mutator.New(license),
		skipListGitIgnore:      buildGitIgnoreSkip(startDirectory),
		skipListLicenserIgnore: buildLicenserIgnoreSkip(startDirectory),
		skipListExtension:      buildExtensionSkip(),
	}
}

// Apply tells the mutator to prepend the license to all walked files
func (p *Processor) Apply(recurse, dryRun bool) bool {
	p.dryRun = dryRun
	p.visitFunc = p.mutator.Apply
	return p.run(recurse)
}

// Verify tells the mutator to check that all walked files have a license
func (p *Processor) Verify(recurse bool) bool {
	p.visitFunc = p.mutator.Verify
	return p.run(recurse)
}

func (p *Processor) run(recurse bool) bool {
	p.success = true
	if recurse {
		if err := filepath.Walk(p.startDirectory, p.visit); err != nil {
			fmt.Fprintf(os.Stderr, "error walking filepath: %v", err)
		}
	} else {
		f, err := os.Open(p.startDirectory)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening info for directory %v:%v\n", p.startDirectory, err)
			return false
		}
		fileList, err := f.Readdir(0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading contents of directory %v:%v\n", p.startDirectory, err)
			return false
		}
		for _, file := range fileList {
			filePath := filepath.Join(p.startDirectory, file.Name())
			if err := p.visit(filePath, file, nil); err != nil {
				fmt.Fprintf(os.Stderr, "error visiting %q: %v", filePath, err)
			}
		}
	}
	p.wg.Wait()
	return p.success
}

func (p *Processor) visit(path string, f os.FileInfo, err error) error {
	if p.shouldSkip(path) {
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

func (p *Processor) shouldSkip(path string) bool {
	// skip predefined file types
	if _, ok := p.skipListExtension[filepath.Ext(path)]; ok {
		return true
	}
	// skip .git/**, .gitignore, .gitattributes, etc
	if strings.Contains(path, ".git") {
		return true
	}
	// don't skip start dir (it cannot be ignored by *ignore files anyway,
	// and `go-gitignore` library crashes on this use case)
	if path == p.startDirectory {
		return false
	}
	// skip according to .gitignore
	if match := p.skipListGitIgnore.Match(path); match != nil && match.Ignore() {
		return true
	}
	// skip .licenserignore
	if filepath.Base(path) == licenserignoreFile {
		return true
	}
	// skip according to .licenserignore
	if match := p.skipListLicenserIgnore.Match(path); match != nil && match.Ignore() {
		return true
	}
	return false
}

func buildGitIgnoreSkip(startDirectory string) gitignore.GitIgnore {
	gitignore, err := gitignore.NewRepository(startDirectory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading contents of .gitignore:%v\n", err)
	}
	return gitignore
}

func buildLicenserIgnoreSkip(startDirectory string) gitignore.GitIgnore {
	ignore, err := gitignore.NewRepositoryWithFile(startDirectory, licenserignoreFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading contents of %s:%v\n", licenserignoreFile, err)
	}
	return ignore
}

func buildExtensionSkip() map[string]bool {
	return map[string]bool{
		".md":     true,
		".golden": true,
	}
}
