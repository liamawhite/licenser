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

import "io"

// Handler is the interface required to implement a license
type Handler interface {
	// Reader returns an io.Reader of the license bytes that will be prepended to files
	Reader() io.Reader

	// IsPresent returns true if it can find the license in the passed reader
	IsPresent(in io.Reader) bool
}
