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

package config

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"sigs.k8s.io/yaml"
)

const DefaultConfigFilename = "./yashiro.yaml"

// Config is Yashiro configuration.
type Config struct {
	Aws *AwsConfig `json:"aws,omitempty"`
}

// AwsConfig is AWS service configuration.
type AwsConfig struct {
	ParameterStoreValues []AwsParameterStoreValueConfig `json:"parameter_store,omitempty"`
	SecretsManagerValues []ValueConfig                  `json:"secrets_manager,omitempty"`
	SdkConfig            *aws.Config                    `json:"-"`
}

// ValueConfig is a value of external store configuration.
type ValueConfig struct {
	Name   string  `json:"name"`
	Ref    *string `json:"ref,omitempty"`
	IsJSON bool    `json:"is_json"`
}

// AwsParameterStoreValueConfig is a AWS Systems Manager Parameter Store configuration. This
// is extended ValueConfig for parameter decryption.
type AwsParameterStoreValueConfig struct {
	ValueConfig
	Decryption *bool `json:"decryption,omitempty"`
}

// NewFromFile returns a new Config according to a file. The configuration file is assumed to
// be in YAML format.
func NewFromFile(ctx context.Context, filename string) (*Config, error) {
	b, err := getConfigFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, err
	}

	if cfg.Aws != nil {
		awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, err
		}
		cfg.Aws.SdkConfig = &awsCfg
	}

	return cfg, nil
}

// Value is interface of external store value.
type Value interface {
	GetReferenceName() string
	GetIsJSON() bool
}

// GetReferenceName returns name of variable reference. If Ref is not set, returns Name.
func (c ValueConfig) GetReferenceName() string {
	if c.Ref != nil && len(*c.Ref) != 0 {
		return *c.Ref
	}

	return c.Name
}

func (c ValueConfig) GetIsJSON() bool {
	return c.IsJSON
}

func getConfigFile(filename string) ([]byte, error) {
	if len(filename) == 0 {
		filename = DefaultConfigFilename
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}
