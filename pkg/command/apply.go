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
	"errors"
	"time"

	"github.com/liamawhite/licenser/pkg/license"
	"github.com/liamawhite/licenser/pkg/processor"
	"github.com/spf13/cobra"
)

var (
	isDryRun     bool
	templatePath string
	markerString string
)

var applyCmd = &cobra.Command{
	Use:   "apply [-t <template file> -m <license-mark>] <copyright-owner>",
	Short: "Apply licenses to files in your directory",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("not enough arguments passed")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		handler, err := newHandler(templatePath, markerString, args[0])
		if err != nil {
			return err
		}

		l := processor.New(".", handler)
		if ok := l.Apply(recurseDirectories, isDryRun); !ok {
			return errors.New("error applying license")
		}

		return nil
	},
}

func init() {
	applyCmd.Flags().BoolVarP(&isDryRun, "dry-run", "d", false, "output result to stdout")
	applyCmd.Flags().StringVarP(&templatePath, "license-template", "t", "", "license template file to use. By default Apache 2.0 license template is used")
	applyCmd.Flags().StringVarP(&markerString, "license-mark", "m", "", "mark to check against the file when checking if the license header is present")
	rootCmd.AddCommand(applyCmd)
}

func newHandler(template, marker, owner string) (license.Handler, error) {
	var h license.Handler
	if template == "" {
		h = license.NewApache20(time.Now().Year(), owner)
	} else {
		if marker == "" {
			return nil, errors.New("--license-mark is required when using --license-template")
		}
		h = license.FromTemplateFile(template, marker, time.Now().Year(), owner)
	}
	return h, nil
}
