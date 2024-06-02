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
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/dwango/yashiro/pkg/config"
)

var (
	ErrInvalidCacheType = errors.New("invalid cache type")
)

type Cache interface {
	// Load returns cached string by using the key, and whether or not cache is expired. If a cache is empty,
	// returned a string is nil and expired is true.
	Load(ctx context.Context, key string, decrypt bool) (*string, bool, error)

	// Save saves value to cache. If encrypt is true, value is encrypted before saving.
	Save(ctx context.Context, key string, value *string, encrypt bool) error
}

func New(cfg config.CacheConfig, options ...Option) (Cache, error) {
	expireDuration := config.DefaultExpireDuration
	if cfg.ExpireDuration != 0 {
		expireDuration = time.Duration(cfg.ExpireDuration)
	}

	switch cfg.Type {
	case config.CacheTypeUnspecified, config.CacheTypeMemory:
		return newMemoryCache(expireDuration)
	case config.CacheTypeFile:
		return newFileCache(cfg.File, expireDuration, options...)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidCacheType, cfg.Type)
	}
}

func keyToHex(key string) string {
	return hex.EncodeToString([]byte(key))
}
