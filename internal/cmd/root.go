/**
 * Copyright 2023 DWANGO Co., Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"fmt"
	"strings"

	"github.com/dwango/yashiro/pkg/config"
	"github.com/spf13/cobra"
)

var (
	configFile   string
	globalConfig = &config.Config{}
)

// New returns a new cobra.Command.
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "ysr",
		Short:         "ysh replaces template file according to config file.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	f := cmd.PersistentFlags()
	f.StringVarP(&configFile, "config", "c", config.DefaultConfigFilename, "specify config file.")

	cmd.AddCommand(newTemplateCommand())
	cmd.AddCommand(newVersionCommand())

	return cmd
}

func checkArgsLength(argsReceived int, requiredArgs ...string) error {
	expectedNum := len(requiredArgs)
	if argsReceived != expectedNum {
		arg := "arguments"
		if expectedNum == 1 {
			arg = "argument"
		}
		return fmt.Errorf("this command needs %v %s: %s", expectedNum, arg, strings.Join(requiredArgs, ", "))
	}
	return nil
}

// preLoadConfig is PreRunE function for cobra.Command. This function preloads config file.
func preLoadConfig(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	if err := globalConfig.LoadFromFile(ctx, configFile); err != nil {
		return err
	}

	return nil
}
