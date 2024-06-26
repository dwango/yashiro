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

package yashiro

import (
	"github.com/dwango/yashiro/pkg/config"
	"github.com/dwango/yashiro/pkg/engine"
)

// Engine initializes external store client and template.
type Engine = engine.Engine

// Config is the configuration for this library.
type Config = config.Config

var (
	// NewEngine returns a new Engine.
	NewEngine = engine.New

	// IgnoreNotFound is an option to ignore missing external store values.
	IgnoreNotFound = engine.IgnoreNotFound
)
