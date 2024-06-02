/**
 * Copyright 2024 DWANGO Co., Ltd.
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
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	secs "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	secsTypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/dwango/yashiro/internal/client/cache"
	"github.com/dwango/yashiro/internal/values"
	"github.com/dwango/yashiro/pkg/config"
)

func Test_newAwsClient(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    Client
		wantErr bool
	}{
		{
			name: "error: aws sdk config is nil",
			args: args{
				cfg: &config.Config{Aws: &config.AwsConfig{}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newAwsClient(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("newAwsClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newAwsClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockSsmClient func(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)

func (m mockSsmClient) GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	return m(ctx, params, optFns...)
}

type mockSecsClient func(ctx context.Context, params *secs.GetSecretValueInput, optFns ...func(*secs.Options)) (*secs.GetSecretValueOutput, error)

func (m mockSecsClient) GetSecretValue(ctx context.Context, params *secs.GetSecretValueInput, optFns ...func(*secs.Options)) (*secs.GetSecretValueOutput, error) {
	return m(ctx, params, optFns...)
}

var (
	textStrSsmClient = mockSsmClient(func(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
		return &ssm.GetParameterOutput{
			Parameter: &ssmTypes.Parameter{
				Value: stringPtr("test"),
			},
		}, nil
	})

	textStrSecsClient = mockSecsClient(func(ctx context.Context, params *secs.GetSecretValueInput, optFns ...func(*secs.Options)) (*secs.GetSecretValueOutput, error) {
		return &secs.GetSecretValueOutput{
			SecretString: stringPtr("test"),
		}, nil
	})
)

func Test_awsClient_GetValues(t *testing.T) {
	type fields struct {
		ssmClient           ssmClient
		secsClient          secsClient
		parameterStoreValue []config.AwsParameterStoreValueConfig
		secretsManagerValue []config.ValueConfig
	}
	type args struct {
		ctx            context.Context
		ignoreNotFound bool
	}

	notFoundErrSsmClient := mockSsmClient(func(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
		return nil, &ssmTypes.ParameterNotFound{}
	})
	notFoundErrSecsClient := mockSecsClient(func(ctx context.Context, params *secs.GetSecretValueInput, optFns ...func(*secs.Options)) (*secs.GetSecretValueOutput, error) {
		return nil, &secsTypes.ResourceNotFoundException{}
	})

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    values.Values
		wantErr bool
	}{
		{
			name: "ok: text",
			fields: fields{
				ssmClient:  textStrSsmClient,
				secsClient: textStrSecsClient,
				parameterStoreValue: []config.AwsParameterStoreValueConfig{
					{ValueConfig: config.ValueConfig{Name: "ssmKey"}, Decryption: nil},
				},
				secretsManagerValue: []config.ValueConfig{
					{Name: "secsKey"},
				},
			},
			args: args{
				ctx:            context.Background(),
				ignoreNotFound: false,
			},
			want: values.Values{"ssmKey": "test", "secsKey": "test"},
		},
		{
			name: "ok: json",
			fields: fields{
				ssmClient: mockSsmClient(func(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
					return &ssm.GetParameterOutput{
						Parameter: &ssmTypes.Parameter{
							Value: stringPtr(`{"key":"value"}`),
						},
					}, nil
				}),
				secsClient: mockSecsClient(func(ctx context.Context, params *secs.GetSecretValueInput, optFns ...func(*secs.Options)) (*secs.GetSecretValueOutput, error) {
					return &secs.GetSecretValueOutput{
						SecretString: stringPtr(`{"key":"value"}`),
					}, nil
				}),
				parameterStoreValue: []config.AwsParameterStoreValueConfig{
					{ValueConfig: config.ValueConfig{Name: "ssmKey", IsJSON: true}},
				},
				secretsManagerValue: []config.ValueConfig{
					{Name: "secsKey", IsJSON: true},
				},
			},
			args: args{
				ctx:            context.Background(),
				ignoreNotFound: false,
			},
			want: values.Values{"ssmKey": map[string]any{"key": "value"}, "secsKey": map[string]any{"key": "value"}},
		},
		{
			name: "ok: ignore not found error",
			fields: fields{
				ssmClient:  notFoundErrSsmClient,
				secsClient: notFoundErrSecsClient,
				parameterStoreValue: []config.AwsParameterStoreValueConfig{
					{ValueConfig: config.ValueConfig{Name: "ssmKey"}},
				},
				secretsManagerValue: []config.ValueConfig{
					{Name: "secsKey"},
				},
			},
			args: args{
				ctx:            context.Background(),
				ignoreNotFound: true,
			},
			want: values.Values{},
		},
		{
			name: "error: return not found from ssm",
			fields: fields{
				ssmClient:  notFoundErrSsmClient,
				secsClient: textStrSecsClient,
				parameterStoreValue: []config.AwsParameterStoreValueConfig{
					{ValueConfig: config.ValueConfig{Name: "ssmKey"}},
				},
				secretsManagerValue: []config.ValueConfig{},
			},
			args: args{
				ctx:            context.Background(),
				ignoreNotFound: false,
			},
			wantErr: true,
		},
		{
			name: "error: return not found from secs",
			fields: fields{
				ssmClient:           textStrSsmClient,
				secsClient:          notFoundErrSecsClient,
				parameterStoreValue: []config.AwsParameterStoreValueConfig{},
				secretsManagerValue: []config.ValueConfig{
					{Name: "secsKey"},
				},
			},
			args: args{
				ctx:            context.Background(),
				ignoreNotFound: false,
			},
			wantErr: true,
		},
		{
			name: "error: return another error from ssm",
			fields: fields{
				ssmClient: mockSsmClient(func(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
					return nil, &ssmTypes.InternalServerError{}
				}),
				secsClient: textStrSecsClient,
				parameterStoreValue: []config.AwsParameterStoreValueConfig{
					{ValueConfig: config.ValueConfig{Name: "ssmKey"}},
				},
				secretsManagerValue: []config.ValueConfig{},
			},
			args: args{
				ctx:            context.Background(),
				ignoreNotFound: true,
			},
			wantErr: true,
		},
		{
			name: "error: return another error from secs",
			fields: fields{
				ssmClient: textStrSsmClient,
				secsClient: mockSecsClient(func(ctx context.Context, params *secs.GetSecretValueInput, optFns ...func(*secs.Options)) (*secs.GetSecretValueOutput, error) {
					return nil, &secsTypes.InternalServiceError{}
				}),
				parameterStoreValue: []config.AwsParameterStoreValueConfig{},
				secretsManagerValue: []config.ValueConfig{
					{Name: "secsKey"},
				},
			},
			args: args{
				ctx:            context.Background(),
				ignoreNotFound: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := awsClient{
				ssmClient:           tt.fields.ssmClient,
				secsClient:          tt.fields.secsClient,
				parameterStoreValue: tt.fields.parameterStoreValue,
				secretsManagerValue: tt.fields.secretsManagerValue,
			}
			got, err := c.GetValues(tt.args.ctx, tt.args.ignoreNotFound)
			if (err != nil) != tt.wantErr {
				t.Errorf("awsClient.GetValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("awsClient.GetValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ssmClientWithCache_GetParameter(t *testing.T) {
	var params = &ssm.GetParameterInput{Name: stringPtr("any")}
	var textStrSsmClientWant = &ssm.GetParameterOutput{Parameter: &ssmTypes.Parameter{Value: stringPtr("test")}}

	type fields struct {
		client ssmClient
		cache  cache.Cache
	}
	type args struct {
		ctx    context.Context
		params *ssm.GetParameterInput
		optFns []func(*ssm.Options)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ssm.GetParameterOutput
		wantErr bool
	}{
		{
			name: "ok: get from cache",
			fields: fields{
				client: nil,
				cache:  mockCache{load: mockLoadFunc, save: mockSaveFunc},
			},
			args: args{
				ctx:    context.Background(),
				params: params,
			},
			want: &ssm.GetParameterOutput{Parameter: &ssmTypes.Parameter{Value: stringPtr("value")}},
		},
		{
			name: "ok: get from cache(cache disabled)",
			fields: fields{
				client: textStrSsmClient,
				cache:  nil,
			},
			args: args{
				ctx:    context.Background(),
				params: &ssm.GetParameterInput{Name: stringPtr("key")},
			},
			want: textStrSsmClientWant,
		},
		{
			name: "ok: get from client(no cache)",
			fields: fields{
				client: textStrSsmClient,
				cache:  mockCache{load: mockLoadFuncNotFound, save: mockSaveFunc},
			},
			args: args{
				ctx:    context.Background(),
				params: params,
			},
			want: textStrSsmClientWant,
		},
		{
			name: "ok: get from client(cache expired)",
			fields: fields{
				client: textStrSsmClient,
				cache:  mockCache{load: mockLoadFuncExpired, save: mockSaveFunc},
			},
			args: args{
				ctx:    context.Background(),
				params: params,
			},
			want: textStrSsmClientWant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ssmClientWithCache{
				client: tt.fields.client,
				cache:  tt.fields.cache,
			}
			got, err := c.GetParameter(tt.args.ctx, tt.args.params, tt.args.optFns...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ssmClientWithCache.GetParameter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ssmClientWithCache.GetParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_secsClientWithCache_GetSecretValue(t *testing.T) {
	var params = &secs.GetSecretValueInput{SecretId: stringPtr("any")}
	var textStrSecsClientWant = &secs.GetSecretValueOutput{SecretString: stringPtr("test")}

	type fields struct {
		client secsClient
		cache  cache.Cache
	}
	type args struct {
		ctx    context.Context
		params *secs.GetSecretValueInput
		optFns []func(*secs.Options)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *secs.GetSecretValueOutput
		wantErr bool
	}{
		{
			name: "ok: get from cache",
			fields: fields{
				client: nil,
				cache:  mockCache{load: mockLoadFunc, save: mockSaveFunc},
			},
			args: args{
				ctx:    context.Background(),
				params: params,
			},
			want: &secs.GetSecretValueOutput{SecretString: stringPtr("value")},
		},
		{
			name: "ok: get from cache(cache disabled)",
			fields: fields{
				client: textStrSecsClient,
				cache:  nil,
			},
			args: args{
				ctx:    context.Background(),
				params: params,
			},
			want: textStrSecsClientWant,
		},
		{
			name: "ok: get from client(no cache)",
			fields: fields{
				client: textStrSecsClient,
				cache:  mockCache{load: mockLoadFuncNotFound, save: mockSaveFunc},
			},
			args: args{
				ctx:    context.Background(),
				params: params,
			},
			want: textStrSecsClientWant,
		},
		{
			name: "ok: get from client(cache expired)",
			fields: fields{
				client: textStrSecsClient,
				cache:  mockCache{load: mockLoadFuncExpired, save: mockSaveFunc},
			},
			args: args{
				ctx:    context.Background(),
				params: params,
			},
			want: textStrSecsClientWant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := secsClientWithCache{
				client: tt.fields.client,
				cache:  tt.fields.cache,
			}
			got, err := c.GetSecretValue(tt.args.ctx, tt.args.params, tt.args.optFns...)
			if (err != nil) != tt.wantErr {
				t.Errorf("secsClientWithCache.GetSecretValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("secsClientWithCache.GetSecretValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAwsAccountId(t *testing.T) {
	type args struct {
		sdkConfig *aws.Config
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAwsAccountId(tt.args.sdkConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAwsAccountId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getAwsAccountId() = %v, want %v", got, tt.want)
			}
		})
	}
}
