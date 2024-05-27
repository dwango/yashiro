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
	"crypto/aes"
	"crypto/cipher"
	"reflect"
	"testing"

	"github.com/dwango/yashiro/internal/values"
	"github.com/dwango/yashiro/pkg/config"
)

func Test_newFileCache(t *testing.T) {
	type args struct {
		cfg config.FileCacheConfig
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
			got, err := newFileCache(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("newFileCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newFileCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileCache_SaveAndLoad(t *testing.T) {
	block, _ := aes.NewCipher([]byte("0123456789abcdef0123456789abcdef"))

	type fields struct {
		cacheBasePath string
		cipherBlock   cipher.Block
		expired       bool
	}
	type args struct {
		in0 context.Context
		val values.Values
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    values.Values
		wantErr bool
	}{
		{
			name: "ok: save and load",
			fields: fields{
				cacheBasePath: "testdata/save-and-load",
				cipherBlock:   block,
				expired:       true,
			},
			args: args{
				in0: context.Background(),
				val: values.Values{
					"key": "value",
				},
			},
			want: values.Values{
				"key": "value",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fileCache{
				cacheBasePath: tt.fields.cacheBasePath,
				cipherBlock:   tt.fields.cipherBlock,
				expired:       tt.fields.expired,
			}
			if err := f.Save(tt.args.in0, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("fileCache.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, _, err := f.Load(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileCache.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileCache.Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileCache_readOrCreateKey(t *testing.T) {
	type fields struct {
		cacheBasePath string
		cipherBlock   cipher.Block
		expired       bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			// This test case is executed only once. Therefore, if you want to retest, delete the file
			// before executing it again.
			name: "ok: create key",
			fields: fields{
				cacheBasePath: "testdata/read-or-create-key",
			},
			wantErr: false,
		},
		{
			name: "ok: read key",
			fields: fields{
				cacheBasePath: "testdata/read-or-create-key",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fileCache{
				cacheBasePath: tt.fields.cacheBasePath,
				cipherBlock:   tt.fields.cipherBlock,
				expired:       tt.fields.expired,
			}
			if _, err := f.readOrCreateKey(); (err != nil) != tt.wantErr {
				t.Errorf("fileCache.readOrCreateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_fileCache_encryptAndDecryptCache(t *testing.T) {
	block, _ := aes.NewCipher([]byte("0123456789abcdef0123456789abcdef"))

	type fields struct {
		cacheBasePath string
		cipherBlock   cipher.Block
		expired       bool
	}
	type args struct {
		values values.Values
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok: encrypt values",
			fields: fields{
				cipherBlock: block,
			},
			args: args{
				values: values.Values{
					"key": "value",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fileCache{
				cacheBasePath: tt.fields.cacheBasePath,
				cipherBlock:   tt.fields.cipherBlock,
				expired:       tt.fields.expired,
			}
			gotEncrypt, err := f.encryptCache(tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileCache.encryptCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotDecrypt, err := f.decryptCache(gotEncrypt)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileCache.decryptCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDecrypt, tt.args.values) {
				t.Errorf("fileCache.decryptCache() = %v, want %v", gotDecrypt, tt.args.values)
			}
		})
	}
}
