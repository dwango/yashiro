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

	kms "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/dwango/yashiro/internal/values"
	"github.com/dwango/yashiro/pkg/config"
)

func Test_newAwsClient(t *testing.T) {
	type args struct {
		cfg *config.AwsConfig
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
				cfg: &config.AwsConfig{},
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

type mockKmsClient func(ctx context.Context, params *kms.GetSecretValueInput, optFns ...func(*kms.Options)) (*kms.GetSecretValueOutput, error)

func (m mockKmsClient) GetSecretValue(ctx context.Context, params *kms.GetSecretValueInput, optFns ...func(*kms.Options)) (*kms.GetSecretValueOutput, error) {
	return m(ctx, params, optFns...)
}

func Test_awsClient_GetValues(t *testing.T) {
	type fields struct {
		ssmClient           ssmClient
		kmsClient           kmsClient
		parameterStoreValue []config.AwsParameterStoreValueConfig
		secretsManagerValue []config.ValueConfig
	}
	type args struct {
		ctx            context.Context
		ignoreNotFound bool
	}
	returnStrPtr := func(s string) *string { return &s }

	textStrSsmClient := mockSsmClient(func(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
		return &ssm.GetParameterOutput{
			Parameter: &ssmTypes.Parameter{
				Value: returnStrPtr("test"),
			},
		}, nil
	})
	textStrKmsClient := mockKmsClient(func(ctx context.Context, params *kms.GetSecretValueInput, optFns ...func(*kms.Options)) (*kms.GetSecretValueOutput, error) {
		return &kms.GetSecretValueOutput{
			SecretString: returnStrPtr("test"),
		}, nil
	})

	notFoundErrSsmClient := mockSsmClient(func(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
		return nil, &ssmTypes.ParameterNotFound{}
	})
	notFoundErrKmsClient := mockKmsClient(func(ctx context.Context, params *kms.GetSecretValueInput, optFns ...func(*kms.Options)) (*kms.GetSecretValueOutput, error) {
		return nil, &kmsTypes.ResourceNotFoundException{}
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
				ssmClient: textStrSsmClient,
				kmsClient: textStrKmsClient,
				parameterStoreValue: []config.AwsParameterStoreValueConfig{
					{ValueConfig: config.ValueConfig{Name: "ssmKey"}, Decryption: nil},
				},
				secretsManagerValue: []config.ValueConfig{
					{Name: "kmsKey"},
				},
			},
			args: args{
				ctx:            context.Background(),
				ignoreNotFound: false,
			},
			want: values.Values{"ssmKey": "test", "kmsKey": "test"},
		},
		{
			name: "ok: json",
			fields: fields{
				ssmClient: mockSsmClient(func(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
					return &ssm.GetParameterOutput{
						Parameter: &ssmTypes.Parameter{
							Value: returnStrPtr(`{"key":"value"}`),
						},
					}, nil
				}),
				kmsClient: mockKmsClient(func(ctx context.Context, params *kms.GetSecretValueInput, optFns ...func(*kms.Options)) (*kms.GetSecretValueOutput, error) {
					return &kms.GetSecretValueOutput{
						SecretString: returnStrPtr(`{"key":"value"}`),
					}, nil
				}),
				parameterStoreValue: []config.AwsParameterStoreValueConfig{
					{ValueConfig: config.ValueConfig{Name: "ssmKey", IsJSON: true}},
				},
				secretsManagerValue: []config.ValueConfig{
					{Name: "kmsKey", IsJSON: true},
				},
			},
			args: args{
				ctx:            context.Background(),
				ignoreNotFound: false,
			},
			want: values.Values{"ssmKey": map[string]any{"key": "value"}, "kmsKey": map[string]any{"key": "value"}},
		},
		{
			name: "ok: ignore not found error",
			fields: fields{
				ssmClient: notFoundErrSsmClient,
				kmsClient: notFoundErrKmsClient,
				parameterStoreValue: []config.AwsParameterStoreValueConfig{
					{ValueConfig: config.ValueConfig{Name: "ssmKey"}},
				},
				secretsManagerValue: []config.ValueConfig{
					{Name: "kmsKey"},
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
				ssmClient: notFoundErrSsmClient,
				kmsClient: textStrKmsClient,
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
			name: "error: return not found from kms",
			fields: fields{
				ssmClient:           textStrSsmClient,
				kmsClient:           notFoundErrKmsClient,
				parameterStoreValue: []config.AwsParameterStoreValueConfig{},
				secretsManagerValue: []config.ValueConfig{
					{Name: "kmsKey"},
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
				kmsClient: textStrKmsClient,
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
			name: "error: return another error from kms",
			fields: fields{
				ssmClient: textStrSsmClient,
				kmsClient: mockKmsClient(func(ctx context.Context, params *kms.GetSecretValueInput, optFns ...func(*kms.Options)) (*kms.GetSecretValueOutput, error) {
					return nil, &kmsTypes.InternalServiceError{}
				}),
				parameterStoreValue: []config.AwsParameterStoreValueConfig{},
				secretsManagerValue: []config.ValueConfig{
					{Name: "kmsKey"},
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
				kmsClient:           tt.fields.kmsClient,
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
