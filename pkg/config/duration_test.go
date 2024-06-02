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
package config

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		d       Duration
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			d:    Duration(24 * time.Hour),
			want: `"24h0m0s"`,
		},
	}
	for _, tt := range tests {
		got, err := json.Marshal(tt.d)
		if (err != nil) != tt.wantErr {
			t.Errorf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(string(got), tt.want) {
			t.Errorf("json.Marshal() = %v, want %v", string(got), tt.want)
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {
	type testStruct struct {
		Duration Duration `json:"duration"`
	}
	tests := []struct {
		name    string
		json    string
		want    testStruct
		wantErr bool
	}{
		{
			name: "ok: string",
			json: `{"duration": "1s"}`,
			want: testStruct{
				Duration: Duration(time.Second),
			},
		},
		{
			name: "ok: number",
			json: `{"duration": 1000000000}`, // 1s
			want: testStruct{
				Duration: Duration(time.Second),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got testStruct
			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonUnmarshal() = %v, want %v", got, tt.want)
			}
		})
	}

}
