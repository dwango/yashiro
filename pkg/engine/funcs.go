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
	"encoding/json"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"sigs.k8s.io/yaml"
)

func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	extra := template.FuncMap{
		"toYaml":        toYAML,
		"fromYaml":      fromYAML,
		"fromYamlArray": fromYAMLArray,
		"toJson":        toJSON,
		"fromJson":      fromJSON,
		"fromJsonArray": fromJSONArray,
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

// toYAML takes an any, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func toYAML(v any) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}

// fromYAML converts a YAML document into a map[string]any.
//
// This is not a general-purpose YAML parser, and will not parse all valid
// YAML documents.
func fromYAML(str string) (map[string]any, error) {
	m := map[string]any{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		return nil, err
	}

	return m, nil
}

// fromYAMLArray converts a YAML array into a []any.
//
// This is not a general-purpose YAML parser, and will not parse all valid
// YAML documents.
func fromYAMLArray(str string) ([]any, error) {
	a := []any{}

	if err := yaml.Unmarshal([]byte(str), &a); err != nil {
		return nil, err
	}

	return a, nil
}

// toJSON takes an interface, marshals it to json, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}

// fromJSON converts a JSON document into a map[string]any.
//
// This is not a general-purpose JSON parser, and will not parse all valid
// JSON documents.
func fromJSON(str string) (map[string]any, error) {
	m := map[string]any{}

	if err := json.Unmarshal([]byte(str), &m); err != nil {
		return nil, err
	}

	return m, nil
}

// fromJSONArray converts a JSON array into a []any.
//
// This is not a general-purpose JSON parser, and will not parse all valid
// JSON documents.
func fromJSONArray(str string) ([]any, error) {
	a := []any{}

	if err := json.Unmarshal([]byte(str), &a); err != nil {
		return nil, err
	}

	return a, nil
}
