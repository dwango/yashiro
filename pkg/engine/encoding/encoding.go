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
	"errors"
	"fmt"
)

type TextType string

// Define text types
const (
	TextTypeJSON      TextType = "json"
	TextTypeJSONArray TextType = "json-array"
	TextTypeYAML      TextType = "yaml"
	TextTypeYAMLArray TextType = "yaml-array"
	TextTypeYAMLDocs  TextType = "yaml-docs"
)

// Define errors
var (
	ErrUnsupportedTextType     = errors.New("unsupported text type")
	ErrFailedToEncodeAndDecode = errors.New("failed to encode and decode")
)

// EncodeAndDecoder is an interface that provides encoding and decoding functionality.
type EncodeAndDecoder interface {
	EncodeAndDecode(b []byte) ([]byte, error)
}

func NewEncodeAndDecoder(t TextType) (EncodeAndDecoder, error) {
	switch t {
	case TextTypeJSON:
		return &jsonEncodeAndDecoder{}, nil
	case TextTypeJSONArray:
		return &jsonEncodeAndDecoder{isArray: true}, nil
	case TextTypeYAML:
		return &yamlEncodeAndDecoder{docType: yamlDocTypeSingle}, nil
	case TextTypeYAMLArray:
		return &yamlEncodeAndDecoder{docType: yamlDocTypeArray}, nil
	case TextTypeYAMLDocs:
		return &yamlEncodeAndDecoder{docType: yamlDocTypeMulti}, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedTextType, t)
	}
}
