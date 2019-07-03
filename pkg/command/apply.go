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
	"os"
	"time"

	"github.com/liamawhite/licenser/pkg/license"
	"github.com/liamawhite/licenser/pkg/processor"
	"github.com/spf13/cobra"
)

var (
	isDryRun bool
)

var applyCmd = &cobra.Command{
	Use:   "apply <copyright-owner>",
	Short: "Apply licenses to files in your directory",

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) <= 1 {
			return errors.New("not enough arguments passed")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		license := license.NewApache20(time.Now().Year(), args[1])
		l := processor.New(".", license)
		if ok := l.Apply(recurseDirectories, isDryRun); !ok {
			os.Exit(1)
		}
	},
}

func init() {
	applyCmd.Flags().BoolVarP(&isDryRun, "dry-run", "d", false, "output result to stdout")
	rootCmd.AddCommand(applyCmd)
}
