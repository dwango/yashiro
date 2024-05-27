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
package cache

import (
	"context"
	"maps"

	"github.com/dwango/yashiro/internal/values"
)

type memoryCache struct {
	values values.Values
}

func newMemoryCache() (Cache, error) {
	return &memoryCache{}, nil
}

// Load implements Cache.
func (m memoryCache) Load(_ context.Context) (values.Values, bool, error) {
	expired := false
	if len(m.values) == 0 {
		expired = true
	}

	return m.values, expired, nil
}

// Save implements Cache.
func (m *memoryCache) Save(_ context.Context, val values.Values) error {
	newVal := make(values.Values, len(val))
	maps.Copy(newVal, val)
	m.values = newVal

	return nil
}
