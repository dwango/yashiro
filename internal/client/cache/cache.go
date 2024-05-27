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
	"errors"
	"fmt"

	"github.com/dwango/yashiro/internal/values"
	"github.com/dwango/yashiro/pkg/config"
)

var (
	ErrInvalidCacheType = errors.New("invalid cache type")
)

type Cache interface {
	// Load returns values from cache and whether or not cache is expired. If cache is empty,
	// returned values is empty and expired=true.
	Load(ctx context.Context) (values.Values, bool, error)

	// Save saves values to cache.
	Save(ctx context.Context, val values.Values) error
}

func New(cfg config.CacheConfig) (Cache, error) {
	switch cfg.Type {
	case config.CacheTypeUnspecified, config.CacheTypeMemory:
		return newMemoryCache()
	case config.CacheTypeFile:
		return newFileCache(cfg.File)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidCacheType, cfg.Type)
	}
}
