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
	"fmt"
	"os"
	"time"

	"github.com/liamawhite/licenser/pkg/license"
	"github.com/liamawhite/licenser/pkg/processor"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify <directory>",
	Short: "Verify licenses are present in files in your directory",

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("not enough arguments passed")
		}
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		fInfo, err := f.Stat()
		if err != nil {
			return err
		}
		if !fInfo.IsDir() {
			return fmt.Errorf("%q is not a directory", args[0])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		license := license.NewApache20(time.Now().Year(), "")
		l := processor.New(args[0], license)
		if ok := l.Verify(recurseDirectories); !ok {
			os.Exit(1)
		}
		fmt.Println("verification successful!")
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
