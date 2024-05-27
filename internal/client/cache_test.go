/**
 * Copyright 2024 DWANGO Co., Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/dwango/yashiro/internal/client/cache"
	"github.com/dwango/yashiro/internal/values"
)

func Test_newClientWithCache(t *testing.T) {
	type args struct {
		client Client
		cache  cache.Cache
	}
	tests := []struct {
		name string
		args args
		want Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newClientWithCache(tt.args.client, tt.args.cache); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newClientWithCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockClient func(ctx context.Context, ignoreNotFound bool) (values.Values, error)

func (m mockClient) GetValues(ctx context.Context, ignoreNotFound bool) (values.Values, error) {
	return m(ctx, ignoreNotFound)
}

type mockCache func(ctx context.Context) (values.Values, bool, error)

func (m mockCache) Load(ctx context.Context) (values.Values, bool, error) {
	return m(ctx)
}
func (m mockCache) Save(ctx context.Context, val values.Values) error {
	return nil
}

func Test_clientWithCache_GetValues(t *testing.T) {
	type fields struct {
		client Client
		cache  cache.Cache
	}
	type args struct {
		ctx            context.Context
		ignoreNotFound bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    values.Values
		wantErr bool
	}{
		{
			name: "ok: get values from cache",
			fields: fields{
				client: mockClient(func(ctx context.Context, ignoreNotFound bool) (values.Values, error) {
					return values.Values{
						"key-client": "value-client",
					}, nil
				}),
				cache: mockCache(func(ctx context.Context) (values.Values, bool, error) {
					return values.Values{
						"key-cache": "value-cache",
					}, false, nil
				}),
			},
			args: args{
				ctx:            context.Background(),
				ignoreNotFound: false,
			},
			want: values.Values{
				"key-cache": "value-cache",
			},
		},
		{
			name: "ok: get values from client(no cache)",
			fields: fields{
				client: mockClient(func(ctx context.Context, ignoreNotFound bool) (values.Values, error) {
					return values.Values{
						"key-client": "value-client",
					}, nil
				}),
				cache: mockCache(func(ctx context.Context) (values.Values, bool, error) {
					return values.Values{}, true, nil
				}),
			},
			args: args{
				ctx:            context.Background(),
				ignoreNotFound: false,
			},
			want: values.Values{
				"key-client": "value-client",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &clientWithCache{
				client: tt.fields.client,
				cache:  tt.fields.cache,
			}
			got, err := c.GetValues(tt.args.ctx, tt.args.ignoreNotFound)
			if (err != nil) != tt.wantErr {
				t.Errorf("clientWithCache.GetValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("clientWithCache.GetValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
