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

// Option is configurable Engine behavior.
type Option func(*opts)

// IgnoreNotFound ignores values are not found in the external store.
func IgnoreNotFound(b bool) Option {
	return func(o *opts) {
		o.IgnoreNotFound = b
	}
}

type opts struct {
	IgnoreNotFound bool
}

var defaultOpts = &opts{
	IgnoreNotFound: false,
}
