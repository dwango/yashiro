/**
 * Copyright 2023 DWANGO Co., Ltd.
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
	"io"
	"text/template"

	"github.com/dwango/yashiro/internal/client"
	"github.com/dwango/yashiro/pkg/config"
	"github.com/dwango/yashiro/pkg/engine/encoding"
)

// Engine is a template engine.
type Engine interface {
	Render(ctx context.Context, text string, dest io.Writer) error
}

type engine struct {
	client           client.Client
	encodeAndDecoder encoding.EncodeAndDecoder
	template         *template.Template
	option           *opts
}

func New(cfg *config.Config, option ...Option) (Engine, error) {
	opts := defaultOpts
	for _, o := range option {
		o(opts)
	}

	cli, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	var encAndDec encoding.EncodeAndDecoder
	if opts.TextType == TextTypePlane {
		encAndDec = &noOpEncodeAndDecoder{}
	} else {
		encAndDec, err = encoding.NewEncodeAndDecoder(opts.TextType)
		if err != nil {
			return nil, err
		}
	}

	tmpl := template.New("yashiro").Option("missingkey=error").Funcs(funcMap())

	return &engine{
		client:           cli,
		encodeAndDecoder: encAndDec,
		template:         tmpl,
		option:           opts,
	}, nil
}

func (e engine) Render(ctx context.Context, text string, dest io.Writer) error {
	values, err := e.client.GetValues(ctx, e.option.IgnoreNotFound)
	if err != nil {
		return err
	}

	return e.render(text, dest, values)
}

func (e engine) render(text string, dest io.Writer, data any) error {
	if _, err := e.template.Parse(text); err != nil {
		return err
	}

	tmp := &bytes.Buffer{}
	if err := e.template.Execute(tmp, data); err != nil {
		return err
	}

	b, err := e.encodeAndDecoder.EncodeAndDecode(tmp.Bytes())
	if err != nil {
		return err
	}

	_, err = io.Copy(dest, bytes.NewReader(b))
	if err != nil {
		return err
	}

	return nil
}

type noOpEncodeAndDecoder struct{}

func (noOpEncodeAndDecoder) EncodeAndDecode(b []byte) ([]byte, error) {
	return b, nil
}
