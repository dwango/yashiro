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
package encoding

import (
	"reflect"
	"testing"
)

func Test_jsonEncodeAndDecoder_EncodeAndDecode(t *testing.T) {
	type fields struct {
		isArray bool
	}
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantStr string
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				isArray: false,
			},
			args: args{
				str: `{"key":"value"}`,
			},
			wantStr: `{"key":"value"}`,
		},
		{
			name: "ok: array",
			fields: fields{
				isArray: true,
			},
			args: args{
				str: `[{"key":"value"},{"key2":"value2"}]`,
			},
			wantStr: `[{"key":"value"},{"key2":"value2"}]`,
		},
		{
			name: "error: invalid json",
			fields: fields{
				isArray: false,
			},
			args: args{
				str: "invalid json",
			},
			wantErr: true,
		},
		{
			name: "error: invalid json array",
			fields: fields{
				isArray: true,
			},
			args: args{
				str: "invalid json",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := jsonEncodeAndDecoder{
				isArray: tt.fields.isArray,
			}
			got, err := ed.EncodeAndDecode([]byte(tt.args.str))
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonEncodeAndDecoder.EncodeAndDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(got), tt.wantStr) {
				t.Errorf("jsonEncodeAndDecoder.EncodeAndDecode() = %v, want %v", string(got), tt.wantStr)
			}
		})
	}
}
