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

package client

import (
	"context"
	"errors"
	"fmt"

	kms "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/dwango/yashiro/pkg/config"
)

type awsClient struct {
	ssmClient           *ssm.Client
	kmsClient           *kms.Client
	parameterStoreValue []config.AwsParameterStoreValueConfig
	secretsManagerValue []config.ValueConfig
}

func newAwsClient(cfg *config.AwsConfig) (Client, error) {
	if cfg.SdkConfig == nil {
		return nil, fmt.Errorf("require aws sdk config")
	}

	return &awsClient{
		ssmClient:           ssm.NewFromConfig(*cfg.SdkConfig),
		kmsClient:           kms.NewFromConfig(*cfg.SdkConfig),
		parameterStoreValue: cfg.ParameterStoreValues,
		secretsManagerValue: cfg.SecretsManagerValues,
	}, nil
}

func (c awsClient) GetValues(ctx context.Context, ignoreEmpty bool) (Values, error) {
	values := make(Values, len(c.parameterStoreValue)+len(c.secretsManagerValue))

	for _, v := range c.parameterStoreValue {
		output, err := c.ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
			Name:           &v.Name,
			WithDecryption: v.Decryption,
		})

		if err != nil {
			return nil, err
		}

		if err := values.SetValue(v, output.Parameter.Value); err != nil {
			if ignoreEmpty && errors.Is(err, ErrValueIsEmpty) {
				continue
			}
			return nil, err
		}
	}

	for _, v := range c.secretsManagerValue {
		output, err := c.kmsClient.GetSecretValue(ctx, &kms.GetSecretValueInput{
			SecretId: &v.Name,
		})
		if err != nil {
			return nil, err
		}

		if err := values.SetValue(v, output.SecretString); err != nil {
			return nil, err
		}
	}

	return values, nil
}
