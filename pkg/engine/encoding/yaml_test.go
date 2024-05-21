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

func Test_yamlEncodeAndDecoder_EncodeAndDecode(t *testing.T) {
	type fields struct {
		docType yamlDocType
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
			name: "ok: single type",
			fields: fields{
				docType: yamlDocTypeSingle,
			},
			args: args{
				str: "---\nkey: value\n",
			},
			wantStr: "key: value\n",
		},
		{
			name: "ok: array type",
			fields: fields{
				docType: yamlDocTypeArray,
			},
			args: args{
				str: "- key: value\n- key2: value2",
			},
			wantStr: "- key: value\n- key2: value2\n",
		},
		{
			name: "ok: multi type",
			fields: fields{
				docType: yamlDocTypeMulti,
			},
			args: args{
				str: "---\nkey: value\n---\nkey2: value2",
			},
			wantStr: "---\nkey: value\n---\nkey2: value2\n",
		},
		{
			name: "ok: single type with comment",
			fields: fields{
				docType: yamlDocTypeSingle,
			},
			args: args{
				str: "# comment\n---\nkey: value\n",
			},
			wantStr: "key: value\n",
		},
		{
			name: "ok: multi type with comment",
			fields: fields{
				docType: yamlDocTypeMulti,
			},
			args: args{
				str: "# comment\n---\nkey: value\n---\nkey2: value2",
			},
			wantStr: "---\nkey: value\n---\nkey2: value2\n",
		},
		{
			name: "error: invalid yaml with single type",
			fields: fields{
				docType: yamlDocTypeSingle,
			},
			args: args{
				str: "invalid yaml",
			},
			wantErr: true,
		},
		{
			name: "error: invalid yaml with array type",
			fields: fields{
				docType: yamlDocTypeArray,
			},
			args: args{
				str: "invalid yaml",
			},
			wantErr: true,
		},
		{
			name: "error: invalid yaml with multi type",
			fields: fields{
				docType: yamlDocTypeMulti,
			},
			args: args{
				str: "invalid yaml",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := yamlEncodeAndDecoder{
				docType: tt.fields.docType,
			}
			got, err := ed.EncodeAndDecode([]byte(tt.args.str))
			if (err != nil) != tt.wantErr {
				t.Errorf("yamlEncodeAndDecoder.EncodeAndDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(got), tt.wantStr) {
				t.Errorf("yamlEncodeAndDecoder.EncodeAndDecode() = %v, want %v", string(got), tt.wantStr)
			}
		})
	}
}

func Test_splitYAMLDocument(t *testing.T) {
	type args struct {
		dataStr string
		atEOF   bool
	}
	tests := []struct {
		name         string
		args         args
		wantAdvance  int
		wantTokenStr string
		wantErr      bool
	}{
		{
			name: "ok: at EOF separated",
			args: args{
				dataStr: "abc\n---\ndef",
				atEOF:   true,
			},
			wantAdvance:  8,
			wantTokenStr: "abc",
		},
		{
			name: "ok: empty",
			args: args{
				dataStr: "",
				atEOF:   true,
			},
			wantAdvance:  0,
			wantTokenStr: "",
		},
		{
			name: "ok: at EOF",
			args: args{
				dataStr: "test",
				atEOF:   true,
			},
			wantAdvance:  4,
			wantTokenStr: "test",
		},
		{
			name: "ok: not at EOF",
			args: args{
				dataStr: "test",
				atEOF:   false,
			},
			wantAdvance:  0,
			wantTokenStr: "",
		},
		{
			name: "ok: at EOF separator without newline",
			args: args{
				dataStr: "---",
				atEOF:   true,
			},
			wantAdvance:  3,
			wantTokenStr: "---",
		},
		{
			name: "ok: at EOF separator",
			args: args{
				dataStr: "---\n",
				atEOF:   true,
			},
			wantAdvance:  4,
			wantTokenStr: "---\n",
		},
		{
			name: "ok: not at EOF separator",
			args: args{
				dataStr: "---\n",
				atEOF:   false,
			},
			wantAdvance:  0,
			wantTokenStr: "",
		},
		{
			name: "ok: not at EOF separator after newline",
			args: args{
				dataStr: "\n---\n",
				atEOF:   false,
			},
			wantAdvance:  5,
			wantTokenStr: "",
		},
		{
			name: "ok: at EOF separator after newline",
			args: args{
				dataStr: "\n---\n",
				atEOF:   true,
			},
			wantAdvance:  5,
			wantTokenStr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAdvance, gotToken, err := splitYAMLDocument([]byte(tt.args.dataStr), tt.args.atEOF)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitYAMLDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAdvance != tt.wantAdvance {
				t.Errorf("splitYAMLDocument() gotAdvance = %v, want %v", gotAdvance, tt.wantAdvance)
			}
			if !reflect.DeepEqual(string(gotToken), tt.wantTokenStr) {
				t.Errorf("splitYAMLDocument() gotToken = %v, want %v", string(gotToken), tt.wantTokenStr)
			}
		})
	}
}
