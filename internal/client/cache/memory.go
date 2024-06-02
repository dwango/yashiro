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
	"strings"
	"time"
)

type memoryCache struct {
	caches         map[string]*cacheData
	expireDuration time.Duration
	keyPrefix      string
}

func newMemoryCache(expireDuration time.Duration, options ...Option) (Cache, error) {
	opts := defaultOpts
	for _, o := range options {
		o(opts)
	}

	keyPrefix := keyToHex(strings.Join(opts.CacheKeys, "_")) + "_"

	return &memoryCache{
		caches:         make(map[string]*cacheData),
		expireDuration: expireDuration,
		keyPrefix:      keyPrefix,
	}, nil
}

type cacheData struct {
	value    string
	saveTime time.Time
}

// Load implements Cache.
func (m memoryCache) Load(_ context.Context, key string, _ bool) (*string, bool, error) {
	data, ok := m.caches[m.keyPrefix+key]
	if !ok {
		return nil, true, nil
	}

	if time.Since(data.saveTime) > m.expireDuration {
		return &data.value, true, nil
	}

	return &data.value, false, nil
}

// Save implements Cache.
func (m *memoryCache) Save(_ context.Context, key string, value *string, _ bool) error {
	if value == nil {
		return nil
	}

	data := &cacheData{
		value:    *value,
		saveTime: time.Now(),
	}
	m.caches[m.keyPrefix+key] = data

	return nil
}
