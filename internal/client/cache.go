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
package client

import (
	"context"

	"github.com/dwango/yashiro/internal/client/cache"
	"github.com/dwango/yashiro/internal/values"
)

type clientWithCache struct {
	client Client
	cache  cache.Cache
}

func newClientWithCache(client Client, cache cache.Cache) Client {
	return &clientWithCache{
		client: client,
		cache:  cache,
	}
}

// GetValues implements Client.
func (c *clientWithCache) GetValues(ctx context.Context, ignoreNotFound bool) (values.Values, error) {
	val, expired, err := c.cache.Load(ctx)
	if err != nil {
		return nil, err
	}

	// if cache is empty, get values from external store.
	if len(val) == 0 || expired {
		val, err = c.client.GetValues(ctx, ignoreNotFound)
		if err != nil {
			return nil, err
		}
	}

	// save values to cache
	if expired {
		if err := c.cache.Save(ctx, val); err != nil {
			return nil, err
		}
	}

	return val, nil
}
