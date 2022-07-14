// Copyright 2020 Tetrate
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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liamawhite/licenser/pkg/license"
)

func Test_shouldSkip(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		// never ignore root dir
		{".", false},

		// ignore *ignore files themselves
		{".gitignore", true},
		{".gitattributes", true},
		{".licenserignore", true},

		// ignore according to .licenserignore files
		{"licenserignore/ignore.yaml", true},
		{"licenserignore/include.yaml", false},
		{"licenserignore/nested/ignore.yaml", true},
		{"licenserignore/nested/include.yaml", false},

		// ignore according to .gitignore files
		{"gitignore/ignore.yaml", true},
		{"gitignore/include.yaml", false},
		{"gitignore/nested/ignore.yaml", true},
		{"gitignore/nested/include.yaml", false},
	}
	processor := New("testdata", license.NewApache20(2020, "ASF"))
	for _, tt := range tests {
		tc := tt
		name := fmt.Sprintf("shouldSkip %s is %t", tc.path, tc.want)
		t.Run(name, func(t *testing.T) {
			path := filepath.Clean(filepath.Join("testdata", tc.path))
			assert.Equal(t, tc.want, processor.shouldSkip(path))
		})
	}
}
