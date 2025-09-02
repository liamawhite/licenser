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

type languageStyle struct {

	// Will this language use block comments for the license?
	// If false, this will use single line comment style
	// WARNING: NOT YET IMPLEMENTED
	isBlock bool

	// The comment string to be used.
	// If isBlock is true this should be the block comment style
	// If false this should be the single line comment style
	comment string
}

var commentStyles = map[string]*languageStyle{
	"bazel":      {isBlock: false, comment: "#"},
	"c":          {isBlock: false, comment: "//"},
	"docker":     {isBlock: false, comment: "#"},
	"golang":     {isBlock: false, comment: "//"},
	"javascript": {isBlock: false, comment: "//"},
	"lua":        {isBlock: false, comment: "--"},
	"make":       {isBlock: false, comment: "#"},
	"protobuf":   {isBlock: false, comment: "//"},
	"python":     {isBlock: false, comment: "#"},
	"rust":       {isBlock: false, comment: "//"},
	"shell":      {isBlock: false, comment: "#"},
	"sql":        {isBlock: false, comment: "--"},
	"terraform":  {isBlock: false, comment: "#"},
	"yaml":       {isBlock: false, comment: "#"},
}

var commonExtensions = map[string]string{
	".c":     "c",
	".c++":   "c",
	".cc":    "c",
	".cpp":   "c",
	".go":    "golang",
	".h":     "c",
	".js":    "javascript",
	".jsx":   "javascript",
	".lua":   "lua",
	".mk":    "make",
	".patch": "shell",
	".proto": "protobuf",
	".py":    "python",
	".rs":    "rust",
	".sh":    "shell",
	".sql":   "sql",
	".tf":    "terraform",
	".ts":    "javascript",
	".tsx":   "javascript",
	".yaml":  "yaml",
	".yml":   "yaml",
}
