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

package cache

import (
	"reflect"
	"testing"
	"time"

	"github.com/dwango/yashiro/pkg/config"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg     config.CacheConfig
		options []Option
	}
	tests := []struct {
		name    string
		args    args
		want    Cache
		wantErr bool
	}{
		{
			name: "error: unknown cache type",
			args: args{
				cfg: config.CacheConfig{
					Type: "unknown",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.cfg, tt.args.options...)
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

const notExpireDuration = 100 * 365 * 24 * time.Hour // 100 years

func stringPtr(s string) *string {
	return &s
}
