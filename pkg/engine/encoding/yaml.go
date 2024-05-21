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
	"bufio"
	"bytes"
	"fmt"

	"sigs.k8s.io/yaml"
)

type yamlDocType int

const (
	yamlDocTypeSingle yamlDocType = iota
	yamlDocTypeArray
	yamlDocTypeMulti
)

type yamlEncodeAndDecoder struct {
	docType yamlDocType
}

const (
	yamlSeparator = "\n---"
	separator     = "---\n"
)

func (ed yamlEncodeAndDecoder) EncodeAndDecode(b []byte) ([]byte, error) {
	switch ed.docType {
	case yamlDocTypeMulti:
		buf := bytes.NewBuffer(make([]byte, 0, len(b)))

		scn := bufio.NewScanner(bytes.NewReader(b))
		scn.Split(splitYAMLDocument)
		for scn.Scan() {
			v := map[string]any{}
			if err := yaml.Unmarshal(scn.Bytes(), &v); err != nil {
				return nil, fmt.Errorf("%w: %w", ErrFailedToEncodeAndDecode, err)
			}
			if len(v) == 0 {
				continue
			}
			b, _ := yaml.Marshal(v)
			buf.WriteString(separator)
			buf.Write(b)
		}
		return buf.Bytes(), nil
	case yamlDocTypeArray:
		v := []any{}
		if err := yaml.Unmarshal(b, &v); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToEncodeAndDecode, err)
		}
		return yaml.Marshal(v)
	default:
		v := map[string]any{}
		if err := yaml.Unmarshal(b, &v); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToEncodeAndDecode, err)
		}
		return yaml.Marshal(v)
	}
}

// splitYAMLDocument is a bufio.SplitFunc for splitting YAML streams into individual documents.
func splitYAMLDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	sep := len([]byte(yamlSeparator))
	if i := bytes.Index(data, []byte(yamlSeparator)); i >= 0 {
		// We have a potential document terminator
		i += sep
		after := data[i:]
		if len(after) == 0 {
			// we can't read any more characters
			if atEOF {
				return len(data), data[:len(data)-sep], nil
			}
			return 0, nil, nil
		}
		if j := bytes.IndexByte(after, '\n'); j >= 0 {
			return i + j + 1, data[0 : i-sep], nil
		}
		return 0, nil, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
