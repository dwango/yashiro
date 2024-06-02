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
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/dwango/yashiro/pkg/config"
)

func Test_newFileCache(t *testing.T) {
	type args struct {
		cfg            config.FileCacheConfig
		expireDuration time.Duration
		options        []Option
	}
	tests := []struct {
		name      string
		args      args
		wantFiles []string
		wantErr   bool
	}{
		{
			name: "ok",
			args: args{
				cfg: config.FileCacheConfig{
					CachePath: "testdata/constructor",
				},
			},
			wantFiles: []string{"testdata/constructor/_key", "testdata/constructor/._keyHash"},
		},
		{
			name: "ok with cache keys option",
			args: args{
				cfg: config.FileCacheConfig{
					CachePath: "testdata/constructor",
				},
				options: []Option{WithCacheKeys("key1")},
			},
			wantFiles: []string{"testdata/constructor/6b657931_key", "testdata/constructor/.6b657931_keyHash"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newFileCache(tt.args.cfg, tt.args.expireDuration, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("newFileCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, file := range tt.wantFiles {
				if _, err := os.Stat(file); err != nil {
					t.Errorf("file is not found = %v", file)
				}
			}
		})
	}
}

func Test_fileCache_SaveAndLoad(t *testing.T) {
	const cachePath = "testdata/save-and-load"
	block, _ := aes.NewCipher([]byte("0123456789abcdef0123456789abcdef"))

	type fields struct {
		cachePath      string
		cipherBlock    cipher.Block
		expireDuration time.Duration
		filenamePrefix string
	}
	type args struct {
		in0     context.Context
		key     string
		value   *string
		encrypt bool
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantExpired bool
		wantErr     bool
	}{
		{
			name: "ok: save and load plain value",
			fields: fields{
				cachePath:      cachePath,
				cipherBlock:    block,
				expireDuration: notExpireDuration,
			},
			args: args{
				key:   "plain-key",
				value: stringPtr("plain-value"),
			},
		},
		{
			name: "ok: save and load encrypted value",
			fields: fields{
				cachePath:      cachePath,
				cipherBlock:    block,
				expireDuration: notExpireDuration,
			},
			args: args{
				key:     "encrypted-key",
				value:   stringPtr("encrypted-value"),
				encrypt: true,
			},
		},
		{
			name: "ok: load expired value",
			fields: fields{
				cachePath:      cachePath,
				cipherBlock:    block,
				expireDuration: 0,
			},
			args: args{
				key:   "expired-key",
				value: stringPtr("expired-value"),
			},
			wantExpired: true,
		},
		{
			name: "ok: with prefix",
			fields: fields{
				cachePath:      cachePath,
				cipherBlock:    block,
				expireDuration: notExpireDuration,
				filenamePrefix: "test_",
			},
			args: args{
				key:   "prefix-key",
				value: stringPtr("prefix-value"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fileCache{
				cachePath:      tt.fields.cachePath,
				cipherBlock:    tt.fields.cipherBlock,
				expireDuration: tt.fields.expireDuration,
				filenamePrefix: tt.fields.filenamePrefix,
			}
			// Save
			if err := f.Save(tt.args.in0, tt.args.key, tt.args.value, tt.args.encrypt); (err != nil) != tt.wantErr {
				t.Errorf("fileCache.Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Load
			gotValue, gotExpired, err := f.Load(tt.args.in0, tt.args.key, tt.args.encrypt)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileCache.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExpired != tt.wantExpired {
				t.Errorf("fileCache.Load() expired = %v, want %v", gotExpired, tt.wantExpired)
			}
			if !reflect.DeepEqual(gotValue, tt.args.value) {
				t.Errorf("fileCache.Load() got = %v, want %v", gotValue, tt.args.value)
			}
		})
	}
}
