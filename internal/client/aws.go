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

	"github.com/aws/aws-sdk-go-v2/aws"
	secs "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	secsTypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/dwango/yashiro/internal/client/cache"
	"github.com/dwango/yashiro/internal/values"
	"github.com/dwango/yashiro/pkg/config"
)

type awsClient struct {
	ssmClient           ssmClient
	secsClient          secsClient
	parameterStoreValue []config.AwsParameterStoreValueConfig
	secretsManagerValue []config.ValueConfig
}

func newAwsClient(cfg *config.Config) (Client, error) {
	if cfg.Aws.SdkConfig == nil {
		return nil, fmt.Errorf("require aws sdk config")
	}

	var cc cache.Cache
	if cfg.Global.EnableCache {
		// get AWS account ID
		accountID, err := getAwsAccountId(cfg.Aws.SdkConfig)
		if err != nil {
			return nil, err
		}
		cc, err = cache.New(cfg.Global.Cache, cache.WithCacheKeys("aws", cfg.Aws.SdkConfig.Region, accountID))
		if err != nil {
			return nil, err
		}
	}

	return &awsClient{
		ssmClient: &ssmClientWithCache{
			client: ssm.NewFromConfig(*cfg.Aws.SdkConfig),
			cache:  cc,
		},
		secsClient: &secsClientWithCache{
			client: secs.NewFromConfig(*cfg.Aws.SdkConfig),
			cache:  cc,
		},
		parameterStoreValue: cfg.Aws.ParameterStoreValues,
		secretsManagerValue: cfg.Aws.SecretsManagerValues,
	}, nil
}

func (c awsClient) GetValues(ctx context.Context, ignoreNotFound bool) (values.Values, error) {
	values := make(values.Values, len(c.parameterStoreValue)+len(c.secretsManagerValue))

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
			return nil, gettingValueError(v.Name, err)
		}

		if err := values.SetValue(v, output.Parameter.Value); err != nil {
			return nil, err
		}
	}

	for _, v := range c.secretsManagerValue {
		output, err := c.secsClient.GetSecretValue(ctx, &secs.GetSecretValueInput{
			SecretId: &v.Name,
		})

		if err != nil {
			var notFoundErr *secsTypes.ResourceNotFoundException
			if ignoreNotFound && errors.As(err, &notFoundErr) {
				continue
			}
			return nil, gettingValueError(v.Name, err)
		}

		if err := values.SetValue(v, output.SecretString); err != nil {
			return nil, err
		}
	}

	return values, nil
}

type ssmClient interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

type ssmClientWithCache struct {
	client ssmClient
	cache  cache.Cache
}

func (c ssmClientWithCache) GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	if c.cache == nil {
		return c.getParameter(ctx, params, optFns...)
	}

	key := *params.Name // Name is required, so do not check nil
	isSensitive := params.WithDecryption != nil && *params.WithDecryption

	// Load from cache.
	value, expired, err := c.cache.Load(ctx, key, isSensitive)
	if err != nil {
		return nil, err
	}

	// If a cache value is expired or not found, get a value from the external store.
	if value == nil || expired {
		output, err := c.getParameter(ctx, params, optFns...)
		if err != nil {
			return nil, err
		}

		// Create or update cache.
		if err := c.cache.Save(ctx, key, output.Parameter.Value, isSensitive); err != nil {
			return nil, err
		}

		return output, nil
	}

	return &ssm.GetParameterOutput{Parameter: &ssmTypes.Parameter{Value: value}}, nil
}

func (c ssmClientWithCache) getParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	output, err := c.client.GetParameter(ctx, params, optFns...)
	if err != nil {
		return nil, err
	}
	return output, nil
}

type secsClient interface {
	GetSecretValue(ctx context.Context, params *secs.GetSecretValueInput, optFns ...func(*secs.Options)) (*secs.GetSecretValueOutput, error)
}

type secsClientWithCache struct {
	client secsClient
	cache  cache.Cache
}

func (c secsClientWithCache) GetSecretValue(ctx context.Context, params *secs.GetSecretValueInput, optFns ...func(*secs.Options)) (*secs.GetSecretValueOutput, error) {
	if c.cache == nil {
		return c.getSecretValue(ctx, params, optFns...)
	}

	key := *params.SecretId // SecretId is required, so do not check nil

	// Load from cache. Secret is always sensitive.
	value, expired, err := c.cache.Load(ctx, key, true)
	if err != nil {
		return nil, err
	}

	// If a cache value is expired or not found, get a value from the external store.
	if value == nil || expired {
		output, err := c.getSecretValue(ctx, params, optFns...)
		if err != nil {
			return nil, err
		}

		// Create or update cache.
		if err := c.cache.Save(ctx, key, output.SecretString, true); err != nil {
			return nil, err
		}

		return output, nil
	}

	return &secs.GetSecretValueOutput{SecretString: value}, nil
}

func (c secsClientWithCache) getSecretValue(ctx context.Context, params *secs.GetSecretValueInput, optFns ...func(*secs.Options)) (*secs.GetSecretValueOutput, error) {
	output, err := c.client.GetSecretValue(ctx, params, optFns...)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func getAwsAccountId(sdkConfig *aws.Config) (string, error) {
	stsClient := sts.NewFromConfig(*sdkConfig)
	output, err := stsClient.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return *output.Account, nil
}
