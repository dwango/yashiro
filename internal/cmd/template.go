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
	"os"
	"path/filepath"
	"strings"

	"github.com/dwango/yashiro/pkg/config"
	"github.com/dwango/yashiro/pkg/engine"
	"github.com/spf13/cobra"
)

const example = `  # specify single file.
  ysr template example.yaml.tmpl

  # specify config file.
  ysr template -c config.yaml example.yaml.tmpl

  # specify multiple files using glob pattern.
  ysr template ./example/*.tmpl
`

var textTypeValues = []string{
	string(engine.TextTypePlane),
	string(engine.TextTypeJSON),
	string(engine.TextTypeJSONArray),
	string(engine.TextTypeYAML),
	string(engine.TextTypeYAMLArray),
	string(engine.TextTypeYAMLDocs),
}

func newTemplateCommand() *cobra.Command {
	var configFile string
	var ignoreNotFound bool
	var textType string

	cmd := cobra.Command{
		Use:     "template <file>",
		Short:   "Generate a replaced text",
		Example: example,
		Args: func(_ *cobra.Command, args []string) error {
			return checkArgsLength(len(args), "template file")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			cfg, err := config.NewFromFile(ctx, configFile)
			if err != nil {
				return err
			}

			eng, err := engine.New(cfg,
				engine.IgnoreNotFound(ignoreNotFound), engine.TextType(engine.TextTypeOpt(textType)),
			)
			if err != nil {
				return err
			}

			b, err := readAllFiles(args[0])
			if err != nil {
				return err
			}

			return eng.Render(ctx, string(b), os.Stdout)
		},
	}

	f := cmd.Flags()
	f.StringVarP(&configFile, "config", "c", config.DefaultConfigFilename, "specify config file.")
	f.StringVar(&textType, "text-type", string(engine.TextTypePlane),
		fmt.Sprintf("specify text type after rendering. available values: %s", strings.Join(textTypeValues, ", ")),
	)
	f.BoolVar(&ignoreNotFound, "ignore-not-found", false, "ignore values are not found in the external store.")

	return &cmd
}

func readAllFiles(pattern string) ([]byte, error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("file not found: '%s'", pattern)
	}

	b := make([]byte, 0, 1024)
	for _, f := range files {
		c, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		b = append(b, c...)
	}

	return b, nil
}
