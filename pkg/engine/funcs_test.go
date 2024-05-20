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

package engine

import (
	"reflect"
	"testing"
)

func Test_funcMap(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		isNil bool
	}{
		{
			name:  "does not exist env function",
			key:   "env",
			isNil: true,
		},
		{
			name:  "does not exist expanddev function",
			key:   "expanddev",
			isNil: true,
		},
		{
			name:  "exists toYaml function",
			key:   "toYaml",
			isNil: false,
		},
		{
			name:  "exists fromYaml function",
			key:   "fromYaml",
			isNil: false,
		},
		{
			name:  "exists fromYamlArray function",
			key:   "fromYamlArray",
			isNil: false,
		},
		{
			name:  "exists toJson function",
			key:   "toJson",
			isNil: false,
		},
		{
			name:  "exists fromJson function",
			key:   "fromJson",
			isNil: false,
		},
		{
			name:  "exists fromJsonArray function",
			key:   "fromJsonArray",
			isNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := funcMap(); (got[tt.key] == nil) != tt.isNil {
				t.Errorf("funcMap()[%v] = %v", tt.key, got[tt.key])
			}
		})
	}
}

func Test_toYAML(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				v: map[string]any{"key": "value"},
			},
			want: `key: value`,
		},
		{
			name: "ok: empty",
			args: args{
				v: map[string]any{},
			},
			want: "{}",
		},
		{
			name: "ok: nil",
			args: args{
				v: nil,
			},
			want: "null",
		},
		{
			name: "ok: invalid type",
			args: args{
				v: make(chan int),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toYAML(tt.args.v); got != tt.want {
				t.Errorf("toYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromYAML(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				str: `key: value`,
			},
			want: map[string]any{"key": "value"},
		},
		{
			name: "ok: empty",
			args: args{
				str: "{}",
			},
			want: map[string]any{},
		},
		{
			name: "ok: null",
			args: args{
				str: "null",
			},
			want: map[string]any{},
		},
		{
			name: "error: invalid yaml",
			args: args{
				str: "invalid yaml",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromYAML(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromYAMLArray(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    []any
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				str: "- key: value\n- key2: value2",
			},
			want: []any{map[string]any{"key": "value"}, map[string]any{"key2": "value2"}},
		},
		{
			name: "ok: empty",
			args: args{
				str: "[]",
			},
			want: []any{},
		},
		{
			name: "ok: null",
			args: args{
				str: "null",
			},
			want: []any{},
		},
		{
			name: "error: invalid yaml",
			args: args{
				str: "invalid yaml",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromYAMLArray(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromYAMLArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromYAMLArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toJSON(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				v: map[string]any{"key": "value"},
			},
			want: `{"key":"value"}`,
		},
		{
			name: "ok: empty",
			args: args{
				v: map[string]any{},
			},
			want: "{}",
		},
		{
			name: "ok: nil",
			args: args{
				v: nil,
			},
			want: "null",
		},
		{
			name: "ok: invalid type",
			args: args{
				v: make(chan int),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toJSON(tt.args.v); got != tt.want {
				t.Errorf("toJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromJSON(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				str: `{"key":"value"}`,
			},
			want: map[string]any{"key": "value"},
		},
		{
			name: "ok: empty",
			args: args{
				str: "{}",
			},
			want: map[string]any{},
		},
		// FIXME: This test case is not deep equal.
		// {
		// 	name: "ok: nil",
		// 	args: args{
		// 		str: "null",
		// 	},
		// 	want: map[string]any{},
		// },
		{
			name: "error: invalid json",
			args: args{
				str: "invalid json",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromJSON(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromJSONArray(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    []any
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				str: `[{"key":"value"},{"key2":"value2"}]`,
			},
			want: []any{map[string]any{"key": "value"}, map[string]any{"key2": "value2"}},
		},
		{
			name: "ok: empty",
			args: args{
				str: "[]",
			},
			want: []any{},
		},
		// FIXME: This test case is not deep equal.
		// {
		// 	name: "ok: null",
		// 	args: args{
		// 		str: "null",
		// 	},
		// 	want: []any{},
		// },
		{
			name: "error: invalid json",
			args: args{
				str: "invalid json",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromJSONArray(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromJSONArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromJSONArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
