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
package values

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dwango/yashiro/pkg/config"
)

// Define errors
var (
	ErrValueIsEmpty = errors.New("value is empty")
	ErrInvalidJSON  = errors.New("invalid json string")
)

// Values are stored values from external stores.
type Values map[string]any

// SetValue sets the getting value from external stores. If value is json string, is set
// as map[string]any.
func (v Values) SetValue(cfg config.Value, value *string) error {
	if value == nil || len(*value) == 0 {
		return ErrValueIsEmpty
	}

	if v == nil {
		v = make(Values)
	}

	var val any
	if cfg.GetIsJSON() {
		val = make(map[string]any)
		if err := json.Unmarshal([]byte(*value), &val); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidJSON, err)
		}
	} else {
		val = *value
	}

	v[cfg.GetReferenceName()] = val

	return nil
}
