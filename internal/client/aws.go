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
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/dwango/yashiro/pkg/config"
)

type ssmClient interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

type kmsClient interface {
	GetSecretValue(ctx context.Context, params *kms.GetSecretValueInput, optFns ...func(*kms.Options)) (*kms.GetSecretValueOutput, error)
}

type awsClient struct {
	ssmClient           ssmClient
	kmsClient           kmsClient
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

func (c awsClient) GetValues(ctx context.Context, ignoreNotFound bool) (Values, error) {
	values := make(Values, len(c.parameterStoreValue)+len(c.secretsManagerValue))

	for _, v := range c.parameterStoreValue {
		output, err := c.ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
			Name:           &v.Name,
			WithDecryption: v.Decryption,
		})

		if err != nil {
			var notFoundErr *ssmTypes.ParameterNotFound
			if ignoreNotFound && errors.As(err, &notFoundErr) {
				continue
			}
			return nil, err
		}

		if err := values.SetValue(v, output.Parameter.Value); err != nil {
			return nil, err
		}
	}

	for _, v := range c.secretsManagerValue {
		output, err := c.kmsClient.GetSecretValue(ctx, &kms.GetSecretValueInput{
			SecretId: &v.Name,
		})

		if err != nil {
			var notFoundErr *kmsTypes.ResourceNotFoundException
			if ignoreNotFound && errors.As(err, &notFoundErr) {
				continue
			}
			return nil, err
		}

		if err := values.SetValue(v, output.SecretString); err != nil {
			return nil, err
		}
	}

	return values, nil
}
