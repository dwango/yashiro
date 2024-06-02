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

	"github.com/dwango/yashiro/pkg/config"
)

func TestNew(t *testing.T) {
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
			name: "error: not found value config",
			args: args{
				cfg: &config.Config{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

var (
	mockLoadFunc = func(_ context.Context, key string, decrypt bool) (*string, bool, error) {
		return stringPtr("value"), false, nil
	}
	mockLoadFuncNotFound = func(_ context.Context, key string, decrypt bool) (*string, bool, error) {
		return nil, true, nil
	}
	mockLoadFuncExpired = func(_ context.Context, key string, decrypt bool) (*string, bool, error) {
		return stringPtr("value"), true, nil
	}

	mockSaveFunc = func(_ context.Context, key string, value *string, encrypt bool) error {
		return nil
	}
)

type mockCache struct {
	load func(ctx context.Context, key string, decrypt bool) (*string, bool, error)
	save func(ctx context.Context, key string, value *string, encrypt bool) error
}

func (m mockCache) Load(ctx context.Context, key string, decrypt bool) (*string, bool, error) {
	return m.load(ctx, key, decrypt)
}

func (m mockCache) Save(ctx context.Context, key string, value *string, encrypt bool) error {
	return m.save(ctx, key, value, encrypt)
}

func stringPtr(s string) *string {
	return &s
}
