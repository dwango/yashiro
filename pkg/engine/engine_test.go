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
	"bytes"
	"context"
	"reflect"
	"testing"
	"text/template"

	"github.com/dwango/yashiro/internal/client"
	"github.com/dwango/yashiro/internal/values"
	"github.com/dwango/yashiro/pkg/config"
	"github.com/dwango/yashiro/pkg/engine/encoding"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg    *config.Config
		option []Option
	}
	tests := []struct {
		name    string
		args    args
		want    Engine
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.cfg, tt.args.option...)
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

type mockClient func(ctx context.Context, ignoreNotFound bool) (values.Values, error)

func (m mockClient) GetValues(ctx context.Context, ignoreNotFound bool) (values.Values, error) {
	return m(ctx, ignoreNotFound)
}

func Test_engine_Render(t *testing.T) {
	type fields struct {
		client           client.Client
		encodeAndDecoder encoding.EncodeAndDecoder
		template         *template.Template
		option           *opts
	}
	type args struct {
		ctx  context.Context
		text string
	}

	createEncodeAndDecoder := func(docType encoding.TextType) encoding.EncodeAndDecoder {
		encAndDec, _ := encoding.NewEncodeAndDecoder(docType)
		return encAndDec
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantDest string
		wantErr  bool
	}{
		{
			name: "ok: render",
			fields: fields{
				client: mockClient(func(ctx context.Context, ignoreNotFound bool) (values.Values, error) {
					return map[string]any{"key": "value"}, nil
				}),
				encodeAndDecoder: &noOpEncodeAndDecoder{},
				template:         template.New("test"),
				option:           &opts{},
			},
			args: args{
				ctx:  context.Background(),
				text: "{{ .key }}",
			},
			wantDest: "value",
		},
		{
			name: "ok: deep render",
			fields: fields{
				client: mockClient(func(ctx context.Context, ignoreNotFound bool) (values.Values, error) {
					return map[string]any{"Values": map[string]any{"key": "value"}}, nil
				}),
				encodeAndDecoder: &noOpEncodeAndDecoder{},
				template:         template.New("test"),
				option:           &opts{},
			},
			args: args{
				ctx:  context.Background(),
				text: "{{ .Values.key }}",
			},
			wantDest: "value",
		},
		{
			name: "ok: render with function",
			fields: fields{
				client: mockClient(func(ctx context.Context, ignoreNotFound bool) (values.Values, error) {
					return map[string]any{"key": "value"}, nil
				}),
				encodeAndDecoder: &noOpEncodeAndDecoder{},
				template:         template.New("test").Funcs(funcMap()),
				option:           &opts{},
			},
			args: args{
				ctx:  context.Background(),
				text: "{{ .key | upper }}",
			},
			wantDest: "VALUE",
		},
		{
			name: "ok: encode and decode as yaml-docs after rendering",
			fields: fields{
				client: mockClient(func(ctx context.Context, ignoreNotFound bool) (values.Values, error) {
					return map[string]any{"key": "value"}, nil
				}),
				encodeAndDecoder: createEncodeAndDecoder(encoding.TextTypeYAMLDocs),
				template:         template.New("test"),
				option:           &opts{},
			},
			args: args{
				ctx:  context.Background(),
				text: "---\nkey: {{ .key }}\n---\n# comment\nkey2: value2\n",
			},
			wantDest: "---\nkey: value\n---\nkey2: value2\n",
		},
		{
			name: "error: failed to get values",
			fields: fields{
				client: mockClient(func(ctx context.Context, ignoreNotFound bool) (values.Values, error) {
					return nil, values.ErrValueIsEmpty
				}),
				encodeAndDecoder: &noOpEncodeAndDecoder{},
				template:         template.New("test"),
				option:           &opts{},
			},
			args: args{
				ctx:  context.Background(),
				text: "{{ .key }}",
			},
			wantErr: true,
		},
		{
			name: "error: failed to parse template",
			fields: fields{
				client: mockClient(func(ctx context.Context, ignoreNotFound bool) (values.Values, error) {
					return map[string]any{"key": "value"}, nil
				}),
				encodeAndDecoder: &noOpEncodeAndDecoder{},
				template:         template.New("test"),
				option:           &opts{},
			},
			args: args{
				ctx:  context.Background(),
				text: "{{ .key",
			},
			wantErr: true,
		},
		{
			name: "error: failed to execute template",
			fields: fields{
				client: mockClient(func(ctx context.Context, ignoreNotFound bool) (values.Values, error) {
					return map[string]any{"key": "value"}, nil
				}),
				encodeAndDecoder: &noOpEncodeAndDecoder{},
				template:         template.New("test").Option("missingkey=error"),
				option:           &opts{},
			},
			args: args{
				ctx:  context.Background(),
				text: "{{ .notfound }}",
			},
			wantErr: true,
		},
		{
			name: "error: failed to encode and decode",
			fields: fields{
				client: mockClient(func(ctx context.Context, ignoreNotFound bool) (values.Values, error) {
					return map[string]any{"key": "value"}, nil
				}),
				encodeAndDecoder: createEncodeAndDecoder(encoding.TextTypeJSON),
				template:         template.New("test"),
				option:           &opts{},
			},
			args: args{
				ctx:  context.Background(),
				text: "{{ .key }}",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := engine{
				client:           tt.fields.client,
				encodeAndDecoder: tt.fields.encodeAndDecoder,
				template:         tt.fields.template,
				option:           tt.fields.option,
			}
			dest := &bytes.Buffer{}
			if err := e.Render(tt.args.ctx, tt.args.text, dest); (err != nil) != tt.wantErr {
				t.Errorf("engine.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDest := dest.String(); gotDest != tt.wantDest {
				t.Errorf("engine.Render() = %v, want %v", gotDest, tt.wantDest)
			}
		})
	}
}

func Test_engine_render(t *testing.T) {
	type fields struct {
		client   client.Client
		template *template.Template
		option   *opts
	}
	type args struct {
		text string
		data any
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantDest string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := engine{
				client:   tt.fields.client,
				template: tt.fields.template,
				option:   tt.fields.option,
			}
			dest := &bytes.Buffer{}
			if err := e.render(tt.args.text, dest, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("engine.render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDest := dest.String(); gotDest != tt.wantDest {
				t.Errorf("engine.render() = %v, want %v", gotDest, tt.wantDest)
			}
		})
	}
}
