/*
 * Copyright (c) 2022-2025 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the License); you may
 * not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an AS IS BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package eval

import (
	"reflect"
	"testing"
)

func TestEval(t *testing.T) {
	t.Parallel()

	var s0 = "1+1"
	var s1 = "1+0.23"
	var s2 = "1+"
	var s3 = ""
	tds := make(Typedefs)

	type args struct {
		s        *string
		typedefs Typedefs
		tdUsed   map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    Value
		wantErr bool
	}{
		{"test " + s0, args{&s0, tds, nil}, Value{t: Integer, i: 2}, false},
		{"test " + s1, args{&s1, tds, nil}, Value{t: Floating, f: 1.23}, false},
		{"test " + s2, args{&s2, tds, nil}, Value{t: Nix}, true},
		{"test eof", args{&s3, tds, nil}, Value{t: Nix}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Eval(tt.args.s, tt.args.typedefs, tt.args.tdUsed)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetValue(t *testing.T) { //nolint:golint,paralleltest
	var tds = make(Typedefs)

	type args struct {
		value    string
		typedefs Typedefs
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"GetInfo", args{"1+1", tds}, 2, false},
		{"GetInfo err", args{"??", tds}, 0, true},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetValue(tt.args.value, tt.args.typedefs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIdValue(t *testing.T) { //nolint:golint,paralleltest
	id1 := "2+3"
	id2 := "=="
	var tds = make(Typedefs)

	type args struct {
		id       string
		typedefs Typedefs
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
	}{
		{"getIdValue", args{id1, tds}, 5, false},
		{"getIdValue err", args{id2, tds}, 0, true},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIdValue(tt.args.id, tt.args.typedefs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIdValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetIdValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
