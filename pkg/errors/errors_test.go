/**
 * Copyright 2024 DWANGO Co., Ltd.
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

package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dwango/yashiro/internal/client/cache"
	"github.com/dwango/yashiro/pkg/engine"
)

func TestIsCacheProcessingError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true",
			args: args{
				err: fmt.Errorf("test: %w", cache.ErrCacheProcessing),
			},
			want: true,
		},
		{
			name: "false",
			args: args{
				err: fmt.Errorf("test: %w", errors.New("different error")),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCacheProcessingError(tt.args.err); got != tt.want {
				t.Errorf("IsCacheProcessingError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRenderingError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true",
			args: args{
				err: fmt.Errorf("test: %w", engine.ErrRendering),
			},
			want: true,
		},
		{
			name: "false",
			args: args{
				err: fmt.Errorf("test: %w", errors.New("different error")),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRenderingError(tt.args.err); got != tt.want {
				t.Errorf("IsRenderingError() = %v, want %v", got, tt.want)
			}
		})
	}
}
