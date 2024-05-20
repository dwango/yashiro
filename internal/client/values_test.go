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

package client

import (
	"errors"
	"testing"

	"github.com/dwango/yashiro/pkg/config"
)

type mockConfigValue struct {
	isJSON bool
}

func (m mockConfigValue) GetReferenceName() string {
	return "test"
}

func (m mockConfigValue) GetIsJSON() bool {
	return m.isJSON
}

func TestValues_SetValue(t *testing.T) {
	type args struct {
		cfg   config.Value
		value *string
	}
	returnStrPtr := func(s string) *string { return &s }
	tests := []struct {
		name    string
		v       Values
		args    args
		wantErr error
	}{
		{
			name: "ok",
			v:    make(Values),
			args: args{
				cfg:   mockConfigValue{isJSON: false},
				value: returnStrPtr("test"),
			},
		},
		{
			name: "error: value is nil",
			v:    make(Values),
			args: args{
				cfg:   mockConfigValue{isJSON: false},
				value: nil,
			},
			wantErr: ErrValueIsEmpty,
		},
		{
			name: "error: value is empty",
			v:    make(Values),
			args: args{
				cfg:   mockConfigValue{isJSON: false},
				value: returnStrPtr(""),
			},
			wantErr: ErrValueIsEmpty,
		},
		{
			name: "error: value is invalid json",
			v:    make(Values),
			args: args{
				cfg:   mockConfigValue{isJSON: true},
				value: returnStrPtr("INVALID JSON"),
			},
			wantErr: ErrInvalidJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.SetValue(tt.args.cfg, tt.args.value); err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Values.SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
