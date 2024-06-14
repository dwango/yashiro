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
	"errors"
	"fmt"

	"github.com/dwango/yashiro/internal/values"
	"github.com/dwango/yashiro/pkg/config"
)

// Define errors
var (
	ErrNotfoundValueConfig = errors.New("not found value config")
	ErrGettingValue        = errors.New("failed to get value")
)

// Client is the external stores client.
type Client interface {
	GetValues(ctx context.Context, ignoreNotFound bool) (values.Values, error)
}

// New returns a new Client.
func New(cfg *config.Config) (Client, error) {
	var client Client
	var err error

	if cfg.Aws != nil {
		client, err = newAwsClient(cfg)
	}
	if err != nil {
		return nil, err
	}

	if client == nil {
		return nil, ErrNotfoundValueConfig
	}

	return client, nil
}

func gettingValueError(name string, err error) error {
	return fmt.Errorf("%w: name='%s': %w", ErrGettingValue, name, err)
}
