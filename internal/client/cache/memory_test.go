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

	"github.com/dwango/yashiro/internal/values"
)

func Test_newMemoryCache(t *testing.T) {
	tests := []struct {
		name    string
		want    Cache
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newMemoryCache()
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

func Test_memoryCache_Load(t *testing.T) {
	type fields struct {
		values values.Values
	}
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    values.Values
		want1   bool
		wantErr bool
	}{
		{
			name: "ok: get values",
			fields: fields{
				values: values.Values{
					"key": "value",
				},
			},
			args: args{
				in0: context.Background(),
			},
			want: values.Values{
				"key": "value",
			},
		},
		{
			name: "ok: no values(return expired=true)",
			fields: fields{
				values: nil,
			},
			args: args{
				in0: context.Background(),
			},
			want:  nil,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := memoryCache{
				values: tt.fields.values,
			}
			got, got1, err := m.Load(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("memoryCache.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("memoryCache.Load() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("memoryCache.Load() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_memoryCache_Save(t *testing.T) {
	type fields struct {
		values values.Values
	}
	type args struct {
		in0 context.Context
		val values.Values
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok: save values",
			fields: fields{
				values: values.Values{},
			},
			args: args{
				in0: context.Background(),
				val: values.Values{
					"key": "value",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &memoryCache{
				values: tt.fields.values,
			}
			if err := m.Save(tt.args.in0, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("memoryCache.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(m.values, tt.args.val) {
				t.Errorf("memoryCache.Save() got = %v, want %v", m.values, tt.args.val)
			}
		})
	}
}
