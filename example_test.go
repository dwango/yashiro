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

package yashiro_test

import (
	"context"
	"log"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/dwango/yashiro"
	"github.com/dwango/yashiro/pkg/config"
)

func Example() {
	ctx := context.TODO()

	sdkConfig, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	refName := "example"
	cfg := &config.Config{
		Aws: &config.AwsConfig{
			ParameterStoreValues: []config.AwsParameterStoreValueConfig{
				{
					ValueConfig: config.ValueConfig{
						Name:   "/example",
						Ref:    &refName,
						IsJSON: true,
					},
				},
			},
			SdkConfig: &sdkConfig,
		},
	}

	eng, err := yashiro.NewEngine(cfg)
	if err != nil {
		log.Fatalf("failed to create engine: %s", err)
	}

	text := `This is example code. The message is {{ .example.message }}.`

	if err := eng.Render(ctx, text, os.Stdout); err != nil {
		log.Fatalf("failed to render text: %s", err)
	}
}
