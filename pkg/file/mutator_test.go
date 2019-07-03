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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_identifyLanguageStyle(t *testing.T) {
	noLanguage := "nil"
	tests := []struct {
		path string
		want string
	}{
		{"BUILD", "bazel"},
		{"BUILD.test", "bazel"},
		{"test/BUILD", "bazel"},
		{"BUILD/nope", noLanguage},
		{"test/BUILD/nope", noLanguage},
		{"WORKSPACE", "bazel"},
		{"WORKSPACE.test", "bazel"},
		{"test/WORKSPACE", "bazel"},
		{"WORKSPACE/nope", noLanguage},
		{"test/WORKSPACE/nope", noLanguage},

		{"test.c", "c"},
		{"test.cc", "c"},
		{"test.cpp", "c"},
		{"test.c++", "c"},
		{"test/test.c++", "c"},
		{"test.c++/test", noLanguage},

		{"Dockerfile", "docker"},
		{"Dockerfile.test", "docker"},
		{"test/Dockerfile", "docker"},
		{"Dockerfile/nope", noLanguage},
		{"test/Dockerfile/nope", noLanguage},

		{"test.go", "golang"},
		{"test/test.go", "golang"},
		{"test.go/test", noLanguage},

		{"Makefile", "make"},
		{"test/Makefile", "make"},
		{"Makefile/nope", noLanguage},
		{"test/Makefile/nope", noLanguage},

		{"test.proto", "protobuf"},
		{"test/test.proto", "protobuf"},
		{"test.proto/test", noLanguage},

		{"test.py", "python"},
		{"test/test.py", "python"},
		{"test.py/test", noLanguage},

		{"test.sh", "shell"},
		{"test/test.sh", "shell"},
		{"test.sh/test", noLanguage},
		{".bashrc", "shell"},
		{"test/.bashrc", "shell"},
		{"test/bashrc", noLanguage},

		{"test.yaml", "yaml"},
		{"test/test.yaml", "yaml"},
		{"test.yaml/test", noLanguage},
		{"test.yml", "yaml"},
		{"test/test.yml", "yaml"},
		{"test.yml/test", noLanguage},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%s is %s", tt.path, tt.want)
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, commentStyles[tt.want], identifyLanguageStyle(tt.path))
		})
	}
}
