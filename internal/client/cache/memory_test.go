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
	"context"
	"reflect"
	"testing"
	"time"
)

func Test_newMemoryCache(t *testing.T) {
	type args struct {
		expireDuration time.Duration
		options        []Option
	}
	tests := []struct {
		name    string
		args    args
		want    Cache
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newMemoryCache(tt.args.expireDuration, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("newMemoryCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newMemoryCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memoryCache_SaveAndLoad(t *testing.T) {
	type fields struct {
		caches         map[string]*cacheData
		expireDuration time.Duration
		keyPrefix      string
	}
	type args struct {
		in0   context.Context
		key   string
		value *string
		in3   bool
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantExpired bool
		wantErr     bool
	}{
		{
			name: "ok: save and load",
			fields: fields{
				caches:         make(map[string]*cacheData),
				expireDuration: notExpireDuration,
			},
			args: args{
				key:   "key",
				value: stringPtr("value"),
			},
		},
		{
			name: "ok: load expired value",
			fields: fields{
				caches:         make(map[string]*cacheData),
				expireDuration: 0,
			},
			args: args{
				key:   "expired-key",
				value: stringPtr("expired-value"),
			},
			wantExpired: true,
		},
		{
			name: "ok: with key prefix",
			fields: fields{
				caches:         make(map[string]*cacheData),
				expireDuration: notExpireDuration,
				keyPrefix:      "prefix_",
			},
			args: args{
				key:   "prefix-key",
				value: stringPtr("prefix-value"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &memoryCache{
				caches:         tt.fields.caches,
				expireDuration: tt.fields.expireDuration,
				keyPrefix:      tt.fields.keyPrefix,
			}
			// Save
			if err := m.Save(tt.args.in0, tt.args.key, tt.args.value, tt.args.in3); (err != nil) != tt.wantErr {
				t.Errorf("memoryCache.Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Load
			gotValue, gotExpired, err := m.Load(tt.args.in0, tt.args.key, tt.args.in3)
			if (err != nil) != tt.wantErr {
				t.Errorf("memoryCache.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExpired != tt.wantExpired {
				t.Errorf("memoryCache.Load() expired = %v, want %v", gotExpired, tt.wantExpired)
			}
			if !reflect.DeepEqual(gotValue, tt.args.value) {
				t.Errorf("memoryCache.Load() got = %v, want %v", gotValue, tt.args.value)
			}
		})
	}
}
