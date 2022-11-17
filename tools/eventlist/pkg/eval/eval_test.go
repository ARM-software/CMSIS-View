/*
 * Copyright (c) 2022 Arm Limited. All rights reserved.
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

	type args struct {
		s *string
	}
	tests := []struct {
		name    string
		args    args
		want    Value
		wantErr bool
	}{
		{"test " + s0, args{&s0}, Value{t: Integer, i: 2}, false},
		{"test " + s1, args{&s1}, Value{t: Floating, f: 1.23}, false},
		{"test " + s2, args{&s2}, Value{t: Nix}, true},
		{"test eof", args{&s3}, Value{t: Nix}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Eval(tt.args.s)
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
