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

package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dwango/yashiro/pkg/config"
)

// Defines errors
var (
	ErrNotfoundValueConfig = errors.New("not found value config")
	ErrValueIsEmpty        = errors.New("value is empty")
	ErrInvalidJSON         = errors.New("invalid json string")
)

// Client is the external stores client.
type Client interface {
	GetValues(ctx context.Context, ignoreNotFound bool) (Values, error)
}

// New returns a new Client.
func New(cfg *config.Config) (Client, error) {
	if cfg.Aws != nil {
		return newAwsClient(cfg.Aws)
	}

	return nil, ErrNotfoundValueConfig
}

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
