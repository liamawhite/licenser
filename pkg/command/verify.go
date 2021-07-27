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

package command

import (
	"fmt"
	"os"

	"github.com/liamawhite/licenser/pkg/processor"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify [-t <template file> -m <license-mark>]",
	Short: "Verify licenses are present in files in your directory",
	Long: `Verify licenses are present in files in your directory.
	
Verify will ignore the following files:
  - *.md, *.golden
  - .gitignore
  - Files that should be ignored according to .gitignore (experimental)
  - .licenserignore
  - Files that should be ignored according to .licenserignore (experimental)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler, err := newHandler(templatePath, markerString, "")
		if err != nil {
			return err
		}

		l := processor.New(".", handler)
		if ok := l.Verify(recurseDirectories); !ok {
			os.Exit(1)
		}
		fmt.Println("verification successful!")

		return nil
	},
}

func init() {
	verifyCmd.Flags().StringVarP(&templatePath, "license-template", "t", "", "license template file to use. By default Apache 2.0 license template is used")
	verifyCmd.Flags().StringVarP(&markerString, "license-mark", "m", "", "mark to check against the file when checking if the license header is present")
	rootCmd.AddCommand(verifyCmd)
}
