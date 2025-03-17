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
	"errors"
	"math"
	"reflect"
	"testing"
)

var errXx = errors.New("xx")

func TestNumError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		Func string
		Num  string
		Err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"test", fields{"f", "n", errXx}, "expression.f: parsing \"n\": xx"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &NumError{
				Func: tt.fields.Func,
				Num:  tt.fields.Num,
				Err:  tt.fields.Err,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("NumError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpression_get(t *testing.T) {
	t.Parallel()

	var s0 = "a"
	var s1 = ""

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    byte
		want1   int
		wantErr bool
	}{
		{s0, fields{&s0, 0, Value{}}, 'a', 1, false},
		{s1, fields{&s1, 0, Value{}}, 0, 0, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.get()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.get() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Expression.get() %s = %v, want %v", tt.name, got, tt.want)
			}
			if ex.pos != tt.want1 {
				t.Errorf("Expression.get() %s pos = %v, want %v", tt.name, ex.pos, tt.want1)
			}
		})
	}
}

func TestExpression_peek(t *testing.T) {
	t.Parallel()

	var s0 = "a"
	var s1 = ""

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    byte
		want1   int
		wantErr bool
	}{
		{s0, fields{&s0, 0, Value{}}, 'a', 0, false},
		{s1, fields{&s1, 0, Value{}}, 0, 0, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.peek()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.peek() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Expression.peek() %s = %v, want %v", tt.name, got, tt.want)
			}
			if ex.pos != tt.want1 {
				t.Errorf("Expression.peek() %s pos = %v, want %v", tt.name, ex.pos, tt.want1)
			}
		})
	}
}

func Test_lower(t *testing.T) {
	t.Parallel()

	var s0 = "x"
	var s1 = "X"

	type args struct {
		c byte
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{s0, args{s0[0]}, 'x'},
		{s1, args{s1[0]}, 'x'},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := lower(tt.args.c); got != tt.want {
				t.Errorf("lower() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_parseUint(t *testing.T) {
	t.Parallel()

	var s0 = "123a"
	var s1 = "0"
	var s2 = "0x"
	var s3 = "0x1234567"
	var s4 = "1234567"
	var s5 = "12345678901234567890"
	var s6 = "123456789012345678901"
	var s7 = ""
	var s8 = "0xffffffffffffffff"
	var s9 = "0x10000000000000000"

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint64
		want1   int
		wantErr bool
	}{
		{s0, fields{&s0, 0, Value{}}, 123, 3, false},
		{s1, fields{&s1, 0, Value{}}, 0, 1, false},
		{s2, fields{&s2, 0, Value{}}, 0, 2, true},
		{s3, fields{&s3, 0, Value{}}, 0x1234567, 9, false},
		{s4, fields{&s4, 0, Value{}}, 1234567, 7, false},
		{s5, fields{&s5, 0, Value{}}, 12345678901234567890, 20, false},
		{s6, fields{&s6, 0, Value{}}, 0xffffffffffffffff, 21, true},
		{s7, fields{&s7, 0, Value{}}, 0, 0, true},
		{s8, fields{&s8, 0, Value{}}, 0xffffffffffffffff, 18, false},
		{s9, fields{&s9, 0, Value{}}, 0xffffffffffffffff, 19, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.parseUint()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.parseUint() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Expression.parseUint() %s = %v, want %v", tt.name, got, tt.want)
			}
			if ex.pos != tt.want1 {
				t.Errorf("Expression.parseUint() %s pos = %v, want %v", tt.name, ex.pos, tt.want1)
			}
		})
	}
}

func TestExpression_ParseFloat(t *testing.T) {
	t.Parallel()

	var s0 = "0e"
	var s1 = "1.23"
	var s2 = "+1.23"
	var s3 = "-1.23"
	var s4 = "-1.23e5"
	var s5 = "1234567890123456789"
	var s6 = "1234567890123456789."
	var s7 = "."
	var s8 = "+1e+300"
	var s9 = ".."
	var s10 = ""
	var s11 = "+"
	var s12 = ".1"
	var s13 = "+1e+"
	var s14 = "+1ex"
	var s15 = "1"

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    float64
		want1   int
		wantErr bool
	}{
		{s6, fields{&s6, 0, Value{}}, 1234567890123456000., 20, false},
		{s0, fields{&s0, 0, Value{}}, 0, 1, true},
		{s1, fields{&s1, 0, Value{}}, 1.23, 4, false},
		{s2, fields{&s2, 0, Value{}}, 1.23, 5, false},
		{s3, fields{&s3, 0, Value{}}, -1.23, 5, false},
		{s4, fields{&s4, 0, Value{}}, -1.23e5, 7, false},
		{s5, fields{&s5, 0, Value{}}, 1234567890123456000, 19, false},
		{s6, fields{&s6, 0, Value{}}, 1234567890123456000., 20, false},
		{s7, fields{&s7, 0, Value{}}, 0, 1, true},
		{s8, fields{&s8, 0, Value{}}, +1e+300, 7, false},
		{s9, fields{&s9, 0, Value{}}, 0, 1, true},
		{s10, fields{&s10, 0, Value{}}, 0, 0, true},
		{s11, fields{&s11, 0, Value{}}, 0, 1, true},
		{s12, fields{&s12, 0, Value{}}, 0.1, 2, false},
		{s13, fields{&s13, 0, Value{}}, 0, 4, true},
		{s14, fields{&s14, 0, Value{}}, 0, 3, true},
		{s15, fields{&s15, 0, Value{}}, 1.0, 1, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.parseFloat()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.parseFloat() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Expression.parseFloat() %s = %v, want %v", tt.name, got, tt.want)
			}
			if ex.pos != tt.want1 {
				t.Errorf("Expression.parseFloat() %s pos = %v, want %v", tt.name, ex.pos, tt.want1)
			}
		})
	}
}

func TestExpression_hex(t *testing.T) {
	t.Parallel()

	var s0 = "0000x"
	var s1 = "12"
	var s2 = "1bx"
	var s3 = "9fa"
	var s4 = "fF"

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name   string
		fields fields
		wantCx byte
		wantS  string
	}{
		{s0, fields{&s0, 0, Value{}}, 0, "0000"},
		{s1, fields{&s1, 0, Value{}}, 0x12, "12"},
		{s2, fields{&s2, 0, Value{}}, 0x1b, "1b"},
		{s3, fields{&s3, 0, Value{}}, 0xfa, "9fa"},
		{s4, fields{&s4, 0, Value{}}, 0xff, "fF"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			gotCx, gotS := ex.hex()
			if gotCx != tt.wantCx {
				t.Errorf("Expression.hex() %s gotCx = %v, want %v", tt.name, gotCx, tt.wantCx)
			}
			if gotS != tt.wantS {
				t.Errorf("Expression.hex() %s gotS = %v, want %v", tt.name, gotS, tt.wantS)
			}
		})
	}
}

func TestExpression_hex4(t *testing.T) {
	t.Parallel()

	var s0 = "0000"
	var s1 = "12Ab"
	var s2 = "12xb"
	var s3 = "9f"
	var s4 = "fFfF"

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		wantCx  uint16
		wantS   string
		wantErr bool
	}{
		{s0, fields{&s0, 0, Value{}}, 0, "0000", false},
		{s1, fields{&s1, 0, Value{}}, 0x12Ab, "12Ab", false},
		{s2, fields{&s2, 0, Value{}}, 0x12, "12x", true},
		{s3, fields{&s3, 0, Value{}}, 0x9f, "9f", true},
		{s4, fields{&s4, 0, Value{}}, 0xffff, "fFfF", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			gotCx, gotS, err := ex.hex4()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.hex4() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if gotCx != tt.wantCx {
				t.Errorf("Expression.hex4() %s gotCx = %v, want %v", tt.name, gotCx, tt.wantCx)
			}
			if gotS != tt.wantS {
				t.Errorf("Expression.hex4() %s gotS = %v, want %v", tt.name, gotS, tt.wantS)
			}
		})
	}
}

func TestExpression_hex8(t *testing.T) {
	t.Parallel()

	var s0 = "00000000"
	var s1 = "12Ab34Cd"
	var s2 = "12xb"
	var s3 = "9f"
	var s4 = "fFfFfFfF"

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		wantCx  uint32
		wantS   string
		wantErr bool
	}{
		{s0, fields{&s0, 0, Value{}}, 0, "00000000", false},
		{s1, fields{&s1, 0, Value{}}, 0x12Ab34Cd, "12Ab34Cd", false},
		{s2, fields{&s2, 0, Value{}}, 0x12, "12x", true},
		{s3, fields{&s3, 0, Value{}}, 0x9f, "9f", true},
		{s4, fields{&s4, 0, Value{}}, 0xffffffff, "fFfFfFfF", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			gotCx, gotS, err := ex.hex8()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.hex8() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if gotCx != tt.wantCx {
				t.Errorf("Expression.hex8() %s gotCx = %v, want %v", tt.name, gotCx, tt.wantCx)
			}
			if gotS != tt.wantS {
				t.Errorf("Expression.hex8() %s gotS = %v, want %v", tt.name, gotS, tt.wantS)
			}
		})
	}
}

func TestExpression_lex(t *testing.T) {
	t.Parallel()

	var s0 = "+"
	var s1 = "123$"
	var s2 = "0x$"
	var s3 = "1.2$"
	var s4 = "  277e-2$"
	var s5 = "abc$"
	var s6 = "a6Z_c"
	var s7 = "iNf$"
	var s8 = "NaN$" // compare of structs with NaN does not work in test
	var s9 = "\"a\\ax\\by\\eq\\ft\\nb\\rg\\tz\\vsc\"$"
	var s10 = "'X'"
	var s11 = "\"x\\U001234afX\""
	var s12 = "\"q\\u34afQ\""
	var s13 = "'\\u4711'"
	var s14 = "'\\U001234af'"
	var s15 = "//sdjhfaskjfb"
	var s16 = ">>="
	var s17 = "1.2ex"
	var s18 = "\\"
	var s19 = "\"\\"
	var s20 = "\"\\'q\\\"v\\101w\""
	var s21 = "\""
	var s22 = "\"\\0"
	var s23 = "\"\\07"
	var s24 = "\"\\074"
	var s25 = "\"\\07x\""
	var s26 = "\"\\7y\""
	var s27 = "\"\\xay\""
	var s28 = "\"\\x2ay\""
	var s29 = "\"\\x24dy\""
	var s30 = "\"\\x2e5ay\""
	var s31 = "\"\\u471\""
	var s32 = "\"\\U47112\""
	var s33 = "'"
	var s34 = "'\\0"
	var s35 = "'\\07"
	var s36 = "'\\074"
	var s37 = "'\\07'"
	var s38 = "'\\7'"
	var s39 = "'\\xa'"
	var s40 = "'\\x2a'"
	var s41 = "'\\x24d'"
	var s42 = "'\\x2e5a'"
	var s43 = "'\\u471'"
	var s44 = "'\\U47112'"
	var s45 = "'\\"
	var s46 = ">"
	var s47 = ">>"
	var s48 = "'\\''"
	var s49 = "'\\\"'"
	var s50 = "'\\a'"
	var s51 = "'\\b'"
	var s52 = "'\\e'"
	var s53 = "'\\f'"
	var s54 = "'\\n'"
	var s55 = "'\\r'"
	var s56 = "'\\t'"
	var s57 = "'\\v'"

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		want1   int
		wantErr bool
	}{
		{s0, fields{&s0, 0, Value{}}, Value{t: Add}, 1, false},
		{s1, fields{&s1, 0, Value{}}, Value{t: Integer, i: 123}, 3, false},
		{s2, fields{&s2, 0, Value{}}, Value{}, 2, true},
		{s3, fields{&s3, 0, Value{}}, Value{t: Floating, f: 1.2}, 3, false},
		{s4, fields{&s4, 0, Value{}}, Value{t: Floating, f: 2.77}, 8, false},
		{s5, fields{&s5, 0, Value{}}, Value{t: Identifier, s: "abc"}, 3, false},
		{s6, fields{&s6, 0, Value{}}, Value{t: Identifier, s: "a6Z_c"}, 5, false},
		{s7, fields{&s7, 0, Value{}}, Value{t: Floating, f: math.Inf(0)}, 3, false},
		{s8, fields{&s8, 0, Value{}}, Value{t: Floating, f: math.NaN()}, 3, false},
		{s9, fields{&s9, 0, Value{}}, Value{t: String, s: "a\ax\by\x1bq\ft\nb\rg\tz\vsc"}, 28, false},
		{s10, fields{&s10, 0, Value{}}, Value{t: Integer, i: 'X'}, 3, false},
		{s11, fields{&s11, 0, Value{}}, Value{t: String, s: "x\xef\xbf\xbdX"}, 14, false},
		{s12, fields{&s12, 0, Value{}}, Value{t: String, s: "q\xe3\x92\xafQ"}, 10, false},
		{s13, fields{&s13, 0, Value{}}, Value{t: Integer, i: 0x4711}, 8, false},
		{s14, fields{&s14, 0, Value{}}, Value{t: Integer, i: 0x001234af}, 12, false},
		{s15, fields{&s15, 0, Value{}}, Value{t: Nix}, 13, true},
		{s16, fields{&s16, 0, Value{}}, Value{t: ShrAssign}, 3, false},
		{s17, fields{&s17, 0, Value{}}, Value{t: Nix}, 4, true},
		{s18, fields{&s18, 0, Value{}}, Value{t: Nix}, 1, true},
		{s19, fields{&s19, 0, Value{}}, Value{t: Nix}, 2, true},
		{s20, fields{&s20, 0, Value{}}, Value{t: String, s: "'q\"vAw"}, 13, false},
		{s21, fields{&s21, 0, Value{}}, Value{t: Nix}, 1, true},
		{s22, fields{&s22, 0, Value{}}, Value{t: Nix}, 3, true},
		{s23, fields{&s23, 0, Value{}}, Value{t: Nix}, 4, true},
		{s24, fields{&s24, 0, Value{}}, Value{t: Nix}, 5, true},
		{s25, fields{&s25, 0, Value{}}, Value{t: String, s: "\007x"}, 6, false},
		{s26, fields{&s26, 0, Value{}}, Value{t: String, s: "\007y"}, 5, false},
		{s27, fields{&s27, 0, Value{}}, Value{t: String, s: "\ny"}, 6, false},
		{s28, fields{&s28, 0, Value{}}, Value{t: String, s: "*y"}, 7, false},
		{s29, fields{&s29, 0, Value{}}, Value{t: String, s: "My"}, 8, false},
		{s30, fields{&s30, 0, Value{}}, Value{t: String, s: "Zy"}, 9, false},
		{s31, fields{&s31, 0, Value{}}, Value{t: Nix}, 7, true},
		{s32, fields{&s32, 0, Value{}}, Value{t: Nix}, 9, true},
		{s33, fields{&s33, 0, Value{}}, Value{t: Nix}, 1, true},
		{s34, fields{&s34, 0, Value{}}, Value{t: Nix}, 3, true},
		{s35, fields{&s35, 0, Value{}}, Value{t: Nix}, 4, true},
		{s36, fields{&s36, 0, Value{}}, Value{t: Nix}, 5, true},
		{s37, fields{&s37, 0, Value{}}, Value{t: Integer, i: 7}, 5, false},
		{s38, fields{&s38, 0, Value{}}, Value{t: Integer, i: 7}, 4, false},
		{s39, fields{&s39, 0, Value{}}, Value{t: Integer, i: 0xa}, 5, false},
		{s40, fields{&s40, 0, Value{}}, Value{t: Integer, i: 0x2a}, 6, false},
		{s41, fields{&s41, 0, Value{}}, Value{t: Integer, i: 0x4d}, 7, false},
		{s42, fields{&s42, 0, Value{}}, Value{t: Integer, i: 0x5a}, 8, false},
		{s43, fields{&s43, 0, Value{}}, Value{t: Nix}, 7, true},
		{s44, fields{&s44, 0, Value{}}, Value{t: Nix}, 9, true},
		{s45, fields{&s45, 0, Value{}}, Value{t: Nix}, 2, true},
		{s46, fields{&s46, 0, Value{}}, Value{t: Greater}, 1, false},
		{s47, fields{&s47, 0, Value{}}, Value{t: Shr}, 2, false},
		{s48, fields{&s48, 0, Value{}}, Value{t: Integer, i: 0x27}, 4, false},
		{s49, fields{&s49, 0, Value{}}, Value{t: Integer, i: 0x22}, 4, false},
		{s50, fields{&s50, 0, Value{}}, Value{t: Integer, i: 7}, 4, false},
		{s51, fields{&s51, 0, Value{}}, Value{t: Integer, i: 8}, 4, false},
		{s52, fields{&s52, 0, Value{}}, Value{t: Integer, i: 0x1b}, 4, false},
		{s53, fields{&s53, 0, Value{}}, Value{t: Integer, i: 0xc}, 4, false},
		{s54, fields{&s54, 0, Value{}}, Value{t: Integer, i: 0xa}, 4, false},
		{s55, fields{&s55, 0, Value{}}, Value{t: Integer, i: 0xd}, 4, false},
		{s56, fields{&s56, 0, Value{}}, Value{t: Integer, i: 9}, 4, false},
		{s57, fields{&s57, 0, Value{}}, Value{t: Integer, i: 11}, 4, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.lex()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.lex() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if tt.name == "NaN$" { // special case, DeepEqual does not work with NaN
				if got.t != Floating || got.t != tt.want.t ||
					!math.IsNaN(got.f) || !math.IsNaN(tt.want.f) {
					t.Errorf("Expression.lex() %s = %v, want %v", tt.name, got, tt.want)
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.lex() %s = %v, want %v", tt.name, got, tt.want)
			}
			if ex.pos != tt.want1 {
				t.Errorf("Expression.lex() %s pos = %v, want %v", tt.name, ex.pos, tt.want1)
			}
		})
	}
}

func TestExpression_primary(t *testing.T) {
	t.Parallel()

	var s0 = "+"
	var s1 = "4711)+"
	var s2 = "$"
	var s3 = "5"
	var s4 = "6)"
	var s5 = ":me+"
	var s6 = ":me:en+"
	var s7 = ":"
	var s8 = ":xx"
	var s9 = ":1"
	var s10 = ":xx+"
	var s11 = ":me:"
	var s12 = ":me:1"
	var s13 = ":me:ex"
	var s14 = ":me:en"

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantEOF bool
		wantErr bool
	}{
		{"Integer", fields{&s0, 0, Value{t: Integer, i: 0x12345}}, Value{t: Integer, i: 0x12345}, false, false},
		{"Floating", fields{&s0, 0, Value{t: Floating, f: 1.2345}}, Value{t: Floating, f: 1.2345}, false, false},
		{"Identifier", fields{&s0, 0, Value{t: Identifier, s: "vari"}}, Value{t: Identifier, s: "vari"}, false, false},
		{"String", fields{&s0, 0, Value{t: String, s: "abc"}}, Value{t: String, s: "abc"}, false, false},
		{"subExpression", fields{&s1, 0, Value{t: ParenO}}, Value{t: Integer, i: 4711}, false, false},
		{"typedef:member", fields{&s5, 0, Value{t: Identifier, s: "td"}}, Value{t: Integer, i: 123}, false, false},
		{"typedef:member:enum", fields{&s6, 0, Value{t: Identifier, s: "td"}}, Value{t: Integer, i: 4711}, false, false},
		{"Integer_fail", fields{&s2, 0, Value{t: Integer, i: 0x12345}}, Value{t: Integer, i: 0x12345}, false, true},
		{"Floating_fail", fields{&s2, 0, Value{t: Floating, f: 1.2345}}, Value{t: Floating, f: 1.2345}, false, true},
		{"Identifier_fail", fields{&s2, 0, Value{t: Identifier, s: "vari"}}, Value{t: Identifier, s: "vari"}, false, true},
		{"String_fail", fields{&s2, 0, Value{t: String, s: "abc"}}, Value{t: String, s: "abc"}, false, true},
		{"subExpression_fail1", fields{&s2, 0, Value{t: ParenO}}, Value{t: Nix}, false, true},
		{"subExpression_fail2", fields{&s0, 0, Value{t: ParenO}}, Value{t: Nix}, false, true},
		{"subExpression_fail3", fields{&s3, 0, Value{t: ParenO}}, Value{t: Integer, i: 5}, false, true},
		{"subExpression_fail4", fields{&s4, 0, Value{t: ParenO}}, Value{t: Integer, i: 6}, true, true},
		{"typedef:fail", fields{&s7, 0, Value{t: Identifier, s: "td"}}, Value{t: Identifier, s: "td"}, false, true},
		{"typedef:fail1", fields{&s8, 0, Value{t: Identifier, s: "td"}}, Value{t: Identifier, s: "xx"}, true, true},
		{"typedef:fail2", fields{&s9, 0, Value{t: Identifier, s: "td"}}, Value{t: Integer, i: 1}, false, true},
		{"typedef:fail3", fields{&s10, 0, Value{t: Identifier, s: "td"}}, Value{t: Identifier, s: "xx"}, false, true},
		{"typedef:fail4", fields{&s11, 0, Value{t: Identifier, s: "td"}}, Value{t: Identifier, s: "me"}, false, true},
		{"typedef:fail5", fields{&s12, 0, Value{t: Identifier, s: "td"}}, Value{t: Integer, i: 1}, false, true},
		{"typedef:fail6", fields{&s13, 0, Value{t: Identifier, s: "td"}}, Value{t: Identifier, s: "ex"}, false, true},
		{"typedef:fail7", fields{&s14, 0, Value{t: Identifier, s: "td"}}, Value{t: Integer, i: 4711}, true, true},
		{"fail", fields{&s0, 0, Value{t: Add}}, Value{t: Add}, false, true},
	}
	var enums Member
	enums.Offset = "123"
	enums.Enums = make(map[int64]string)
	enums.Enums[4711] = "en"
	var td ITypedef
	td.Members = make(map[string]Member)
	td.Members["me"] = enums
	tds := make(Typedefs)
	tds["td"] = td
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:       tt.fields.in,
				pos:      tt.fields.pos,
				next:     tt.fields.next,
				typedefs: tds,
			}
			got, err := ex.primary()
			if errors.Is(err, ErrEof) != tt.wantEOF {
				t.Errorf("Expression.primary() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.primary() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.primary() %s = %v, want %v", tt.name, got, tt.want)
			}
			if err == nil && ex.next.t != Add {
				t.Errorf("Expression.primary() %s %v, want %v", tt.name, ex.next.t, Add)
			}
		})
	}
}

func TestExpression_arguments(t *testing.T) {
	t.Parallel()

	var s0 = ""
	var s1 = ",123"
	var s2 = ","
	var s3 = ", +"

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantEOF bool
		wantErr bool
	}{
		{"0 arg", fields{&s0, 0, Value{}}, Value{t: Nix}, false, false},
		{"1 arg", fields{&s0, 0, Value{t: Integer, i: 1}}, Value{t: List, l: []Value{{t: Integer, i: 1}}}, false, false},
		{"2 arg", fields{&s1, 0, Value{t: Integer, i: 1}}, Value{t: List, l: []Value{{t: Integer, i: 1}, {t: Integer, i: 123}}}, false, false},
		{"arg err", fields{&s2, 0, Value{t: Integer, i: 1}}, Value{t: List, l: []Value{{t: Integer, i: 1}}}, false, true},
		{"arg err1", fields{&s3, 0, Value{t: Integer, i: 1}}, Value{t: Nix}, false, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.arguments()
			if errors.Is(err, ErrEof) != tt.wantEOF {
				t.Errorf("Expression.arguments() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !errors.Is(err, ErrEof) && (err != nil) != tt.wantErr {
				t.Errorf("Expression.arguments() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.arguments() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_postfix(t *testing.T) { //nolint:golint,paralleltest
	var s0 = "++ +"
	var s1 = "++$"
	var s2 = "++"
	var s3 = "-- +"
	var s4 = "--$"
	var s5 = "--"
	var s6 = ".abc +"
	var s6a = ".b +"
	var s7 = "."
	var s8 = ".123"
	var s9 = ".abc"
	var s10 = "->abc +"
	var s11 = "->"
	var s12 = "->123"
	var s13 = "->abc"
	var s14 = "() +"
	var s15 = "()"
	var s17 = "(\"reg\")"
	var s18 = "("
	var s19 = "(123"
	var s20 = "(+"
	var s21 = "[123] +"
	var s22 = "[123]"
	var s23 = "["
	var s24 = "[123"
	var s25 = "[+"
	var s26 = "[PostfixName]"
	var s27 = "[abc]"
	var s28 = "(1,2,3,4) +"
	var s29 = "(\"reg\") +"

	tds := make(Typedefs)
	tds["td"] = ITypedef{Size: 4, Members: map[string]Member{"b": {Offset: "2", IType: Uint8}}}
	tds["tdo"] = ITypedef{Members: map[string]Member{"b": {Offset: "xxx", IType: Uint8}}}
	tdu := make(map[string]string)
	tdu["tdname"] = "td"
	tdu["tdname1"] = "td"
	tdu["tdnameo"] = "tdo"

	type fields struct {
		in       *string
		pos      int
		next     Value
		typedefs Typedefs
		tdUsed   map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantEOF bool
		wantErr bool
	}{
		{"Postincrement", fields{&s0, 0, Value{t: Identifier, s: "PostfixName"}, nil, nil}, Value{t: Identifier, s: "PostfixName"}, false, false},
		{"Postincrement_fail", fields{&s1, 0, Value{t: Integer, i: 0x12345}, nil, nil}, Value{t: Integer, i: 0x12345}, false, true},
		{"Postincrement_eof", fields{&s2, 0, Value{t: Identifier, s: "PostfixName"}, nil, nil}, Value{t: Identifier, s: "PostfixName"}, true, false},
		{"Postincrement_fail1", fields{&s0, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, false, true},
		{"Postincrement_fail2", fields{&s0, 0, Value{t: Identifier, s: "PostfixName1"}, nil, nil}, Value{t: Nix}, false, true},
		{"Postdecrement", fields{&s3, 0, Value{t: Identifier, s: "PostfixName"}, nil, nil}, Value{t: Identifier, s: "PostfixName"}, false, false},
		{"Postdecrement_fail", fields{&s4, 0, Value{t: Integer, i: 0x12345}, nil, nil}, Value{t: Integer, i: 0x12345}, false, true},
		{"Postdecrement_eof", fields{&s5, 0, Value{t: Identifier, s: "PostfixName"}, nil, nil}, Value{t: Identifier, s: "PostfixName"}, true, false},
		{"Postdecrement_fail1", fields{&s3, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, false, true},
		{"Postdecrement_fail2", fields{&s3, 0, Value{t: Identifier, s: "PostfixName1"}, nil, nil}, Value{t: Nix}, false, true},
		{"Dot", fields{&s6, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, false, false},
		{"Dot_val", fields{&s6a, 0, Value{t: Identifier, s: "tdname"}, tds, tdu}, Value{t: Integer, i: 0}, false, false},
		{"Dot_fail", fields{&s6, 0, Value{t: Integer, i: 0x12345}, nil, nil}, Value{t: Integer, i: 0x12345}, false, true},
		{"Dot_eof_fail", fields{&s7, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, true, true},
		{"Dot_fail1", fields{&s8, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Integer, i: 123}, false, true},
		{"Dot_eof", fields{&s9, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Nix}, true, false},
		{"Dot_val_erroffset", fields{&s6a, 0, Value{t: Identifier, s: "tdnameo"}, tds, tdu}, Value{t: Identifier, s: "b"}, false, true},
		{"Dot_val_errvar", fields{&s6a, 0, Value{t: Identifier, s: "tdname1"}, tds, tdu}, Value{t: Nix}, false, true},
		{"Pointer", fields{&s10, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, false, false},
		{"Pointer_fail", fields{&s10, 0, Value{t: Integer, i: 0x12345}, nil, nil}, Value{t: Integer, i: 0x12345}, false, true},
		{"Pointer_eof_fail", fields{&s11, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, true, true},
		{"Pointer_fail1", fields{&s12, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Integer, i: 123}, false, true},
		{"Pointer_eof", fields{&s13, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Nix}, true, false},
		{"Function", fields{&s14, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, false, false},
		{"Function_eof", fields{&s15, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, true, false},
		{"Function1_eof", fields{&s17, 0, Value{t: Identifier, s: "__GetRegVal"}, nil, nil}, Value{t: Integer, i: 0}, true, false},
		{"Function_GetRegVal", fields{&s29, 0, Value{t: Identifier, s: "__GetRegVal"}, nil, nil}, Value{t: Integer, i: 0}, false, false},
		{"Function_CalcMemUsed", fields{&s28, 0, Value{t: Identifier, s: "__CalcMemUsed"}, nil, nil}, Value{t: Integer, i: 0}, false, false},
		{"Function_FcntErr", fields{&s29, 0, Value{t: Identifier, s: "xxx"}, nil, nil}, Value{t: Identifier, s: "xxx"}, false, true},
		{"Function_err", fields{&s18, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, false, true},
		{"Function_err1", fields{&s19, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, false, true},
		{"Function_err2", fields{&s20, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, false, true},
		{"Index", fields{&s21, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name", i: 123}, false, false},
		{"Index_eof", fields{&s22, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name", i: 123}, true, false},
		{"Index_err", fields{&s23, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, false, true},
		{"Index_err1", fields{&s24, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, false, true},
		{"Index_err2", fields{&s25, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, false, true},
		{"Index_name", fields{&s26, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name", i: 789}, true, false},
		{"Index_name_err", fields{&s27, 0, Value{t: Identifier, s: "name"}, nil, nil}, Value{t: Identifier, s: "name"}, false, true},
	}

	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			ex := &Expression{
				in:       tt.fields.in,
				pos:      tt.fields.pos,
				next:     tt.fields.next,
				typedefs: tt.fields.typedefs,
				tdUsed:   tt.fields.tdUsed,
			}
			ClearNames()
			vari := SetVar("PostfixName", Value{t: Integer, i: 789})
			if tt.fields.next.t == Identifier && tt.fields.next.s == "PostfixName" {
				tt.fields.next.v = vari
				tt.want.v = vari // return value should be the same as input because of postinc
			}
			vari1 := SetVar("PostfixName1", Value{t: Nix})
			if tt.fields.next.t == Identifier && tt.fields.next.s == "PostfixName1" {
				tt.fields.next.v = vari1
			}
			vari2 := SetVar("tdname", Value{t: Integer, i: 123})
			if tt.fields.next.t == Identifier && tt.fields.next.s == "tdname" {
				tt.fields.next.v = vari2
			}
			vari3 := SetVar("tdname1", Value{t: Nix})
			if tt.fields.next.t == Identifier && tt.fields.next.s == "tdname1" {
				tt.fields.next.v = vari3
			}
			got, err := ex.postfix()
			if errors.Is(err, ErrEof) != tt.wantEOF {
				t.Errorf("Expression.primary() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !errors.Is(err, ErrEof) && (err != nil) != tt.wantErr {
				t.Errorf("Expression.postfix() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.postfix() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_unary(t *testing.T) {
	t.Parallel()

	var s0 = "0x12345"
	var s1 = "12.345"
	var s2 = "v_unary"
	var s3 = "$"
	var s4 = ""
	var s5 = "++"
	var s6 = "name"
	var s7 = "\"string\""

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"+IntExpr", fields{&s0, 0, Value{t: Add}}, Value{t: Integer, i: 0x12345}, false},
		{"+IntExpr_err", fields{&s3, 0, Value{t: Add}}, Value{t: Nix}, true},
		{"+IntExpr_eof", fields{&s4, 0, Value{t: Add}}, Value{t: Nix}, true},
		{"+IntExpr_err1", fields{&s5, 0, Value{t: Add}}, Value{t: AddAdd}, true},
		{"+IntExpr_err2", fields{&s6, 0, Value{t: Add}}, Value{t: Identifier, s: "name"}, true},
		{"+IntExpr_err3", fields{&s7, 0, Value{t: Add}}, Value{t: String, s: "string"}, true},
		{"-IntExpr", fields{&s0, 0, Value{t: Sub}}, Value{t: Integer, i: -0x12345}, false},
		{"-IntExpr_err", fields{&s3, 0, Value{t: Sub}}, Value{t: Nix}, true},
		{"-IntExpr_eof", fields{&s4, 0, Value{t: Sub}}, Value{t: Nix}, true},
		{"-IntExpr_err1", fields{&s5, 0, Value{t: Sub}}, Value{t: AddAdd}, true},
		{"-IntExpr_err2", fields{&s6, 0, Value{t: Sub}}, Value{t: Identifier, s: "name"}, true},
		{"-IntExpr_err3", fields{&s7, 0, Value{t: Sub}}, Value{t: String, s: "string"}, true},
		{"~IntExpr", fields{&s0, 0, Value{t: Compl}}, Value{t: Integer, i: 0x12345 ^ -1}, false},
		{"~IntExpr_err", fields{&s3, 0, Value{t: Compl}}, Value{t: Nix}, true},
		{"~IntExpr_eof", fields{&s4, 0, Value{t: Compl}}, Value{t: Nix}, true},
		{"~IntExpr_err1", fields{&s5, 0, Value{t: Compl}}, Value{t: AddAdd}, true},
		{"~IntExpr_err2", fields{&s6, 0, Value{t: Compl}}, Value{t: Identifier, s: "name"}, true},
		{"~IntExpr_err3", fields{&s7, 0, Value{t: Compl}}, Value{t: String, s: "string"}, true},
		{"!IntExpr", fields{&s0, 0, Value{t: Not}}, Value{t: Integer, i: 0}, false},
		{"!IntExpr_err", fields{&s3, 0, Value{t: Not}}, Value{t: Nix}, true},
		{"!IntExpr_eof", fields{&s4, 0, Value{t: Not}}, Value{t: Nix}, true},
		{"!IntExpr_err1", fields{&s5, 0, Value{t: Not}}, Value{t: AddAdd}, true},
		{"!IntExpr_err2", fields{&s6, 0, Value{t: Not}}, Value{t: Identifier, s: "name"}, true},
		{"!IntExpr_err3", fields{&s7, 0, Value{t: Not}}, Value{t: String, s: "string"}, true},
		{"IntExpr", fields{&s0, 0, Value{t: Integer, i: 0x12345}}, Value{t: Integer, i: 0x12345}, false},
		{"+FloatExpr", fields{&s1, 0, Value{t: Add}}, Value{t: Floating, f: 12.345}, false},
		{"-FloatExpr", fields{&s1, 0, Value{t: Sub}}, Value{t: Floating, f: -12.345}, false},
		{"+v1", fields{&s2, 0, Value{t: Add}}, Value{t: Integer, i: 0xa2b3}, false},
		{"-v1", fields{&s2, 0, Value{t: Sub}}, Value{t: Integer, i: -0xa2b3}, false},
		{"~v1", fields{&s2, 0, Value{t: Compl}}, Value{t: Integer, i: 0xa2b3 ^ -1}, false},
		{"!v1", fields{&s2, 0, Value{t: Not}}, Value{t: Integer, i: 0}, false},
		{"v1", fields{&s2, 0, Value{t: Identifier, s: "v1"}}, Value{t: Identifier, s: "v1"}, false},
	}
	SetVar("v_unary", Value{t: Integer, i: 0xa2b3})
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.unary()
			got.v = nil // cannot compare pointers
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.unary() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.unary() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_castExpr(t *testing.T) {
	t.Parallel()

	var s0 = "(uint8_t)v_castExpr"
	var s1 = "(int8_t)0x12345"
	var s2 = "(int16_t)0x12345"
	var s3 = "(int32_t)0x123456789"
	var s4 = "(int64_t)456.789"
	var s5 = "(uint8_t)-0x12345"
	var s6 = "(uint16_t)-0x12345"
	var s7 = "(uint32_t)-0x123456789"
	var s8 = "(uint64_t)456.789"
	var s9 = "(double)12345789"
	var s10 = "(float)123456789"
	var s11 = "($"
	var s12 = "(++"
	var s13 = "(1)"
	var s14 = "(v_castExpr"
	var s15 = "(v_castExpr)"
	var s16 = "(uint8_t)"
	var s17 = "(uint8_t"
	var s18 = "(uint8_t)$"
	var s19 = "(uint8_t)++"
	var s20 = "(uint8_t)name"
	var s21 = "(uint8_t)\"string\""
	var s22 = "(uint8_t+"

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{s0, fields{&s0, 1, Value{t: ParenO}}, Value{t: Integer, i: 0xE3}, false},
		{s1, fields{&s1, 1, Value{t: ParenO}}, Value{t: Integer, i: 0x45}, false},
		{s2, fields{&s2, 1, Value{t: ParenO}}, Value{t: Integer, i: 0x2345}, false},
		{s3, fields{&s3, 1, Value{t: ParenO}}, Value{t: Integer, i: 0x23456789}, false},
		{s4, fields{&s4, 1, Value{t: ParenO}}, Value{t: Integer, i: 456}, false},
		{s5, fields{&s5, 1, Value{t: ParenO}}, Value{t: Integer, i: (-0x12345) & 0xFF}, false},
		{s6, fields{&s6, 1, Value{t: ParenO}}, Value{t: Integer, i: (-0x12345) & 0xFFFF}, false},
		{s7, fields{&s7, 1, Value{t: ParenO}}, Value{t: Integer, i: (-0x23456789) & 0xFFFFFFFF}, false},
		{s8, fields{&s8, 1, Value{t: ParenO}}, Value{t: Integer, i: 456}, false},
		{s9, fields{&s9, 1, Value{t: ParenO}}, Value{t: Floating, f: 12345789.0}, false},
		{s10, fields{&s10, 1, Value{t: ParenO}}, Value{t: Floating, f: 123456792.0}, false},
		{s11, fields{&s11, 1, Value{t: ParenO}}, Value{t: Nix}, true},
		{s12, fields{&s12, 1, Value{t: ParenO}}, Value{t: AddAdd}, true},
		{s13, fields{&s13, 1, Value{t: ParenO}}, Value{t: Integer, i: 1}, false},
		{s14, fields{&s14, 1, Value{t: ParenO}}, Value{t: Floating, f: 483.12}, true},
		{s15, fields{&s15, 1, Value{t: ParenO}}, Value{t: Floating, f: 483.12}, false},
		{s16, fields{&s16, 1, Value{t: ParenO}}, Value{t: Nix}, true},
		{s17, fields{&s17, 1, Value{t: ParenO}}, Value{t: Nix}, true},
		{s18, fields{&s18, 1, Value{t: ParenO}}, Value{t: Nix}, true},
		{s19, fields{&s19, 1, Value{t: ParenO}}, Value{t: AddAdd}, true},
		{s20, fields{&s20, 1, Value{t: ParenO}}, Value{t: Identifier, s: "name"}, true},
		{s21, fields{&s21, 1, Value{t: ParenO}}, Value{t: String, s: "string"}, true},
		{s22, fields{&s22, 1, Value{t: ParenO}}, Value{t: Add}, true},
	}
	SetVar("v_castExpr", Value{t: Floating, f: 483.12})
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.castExpr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.castExpr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.castExpr() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_mulExpr(t *testing.T) {
	t.Parallel()

	var s0 = "*v1_mulExpr"
	var s1 = "*678"
	var s2 = "*6.78"
	var s3 = "/v2_mulExpr"
	var s4 = "/15"
	var s5 = "/1.2"
	var s6 = "%v3_mulExpr"
	var s7 = "%15"
	var s8 = "*"
	var s9 = "*++"
	var s10 = "*name"
	var s11 = "*\"string\""
	var s12 = "/"
	var s13 = "/++"
	var s14 = "/name"
	var s15 = "/\"string\""
	var s16 = "%"
	var s17 = "%++"
	var s18 = "%name"
	var s19 = "%\"string\""

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{s0, fields{&s0, 0, Value{t: Integer, i: 345}}, Value{t: Floating, f: 425.73}, false},
		{"I*I", fields{&s1, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 233910}, false},
		{"I*F", fields{&s2, 0, Value{t: Integer, i: 345}}, Value{t: Floating, f: 2339.1}, false},
		{"F*I", fields{&s1, 0, Value{t: Floating, f: 3.4}}, Value{t: Floating, f: 2305.2}, false},
		{"F*F", fields{&s2, 0, Value{t: Floating, f: 3.4}}, Value{t: Floating, f: 23.052}, false},
		{s3, fields{&s3, 0, Value{t: Integer, i: 345}}, Value{t: Floating, f: 230}, false},
		{"I/I", fields{&s4, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 23}, false},
		{"I/F", fields{&s5, 0, Value{t: Integer, i: 345}}, Value{t: Floating, f: 287.5}, false},
		{"F/I", fields{&s4, 0, Value{t: Floating, f: 3.45}}, Value{t: Floating, f: 0.23}, false},
		{"F/F", fields{&s5, 0, Value{t: Floating, f: 3.6}}, Value{t: Floating, f: 3}, false},
		{s6, fields{&s6, 0, Value{t: Integer, i: 347}}, Value{t: Integer, i: 2}, false},
		{"I%I", fields{&s7, 0, Value{t: Integer, i: 347}}, Value{t: Integer, i: 2}, false},
		{s8, fields{&s8, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{s9, fields{&s9, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s0, fields{&s0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{s10, fields{&s10, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{s11, fields{&s11, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{s12, fields{&s12, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{s13, fields{&s13, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s3, fields{&s3, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{s14, fields{&s14, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{s15, fields{&s15, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{s16, fields{&s16, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{s17, fields{&s17, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s6, fields{&s6, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{s18, fields{&s18, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{s19, fields{&s19, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
	}
	SetVar("v1_mulExpr", Value{t: Floating, f: 1.234})
	SetVar("v2_mulExpr", Value{t: Floating, f: 1.5})
	SetVar("v3_mulExpr", Value{t: Integer, i: 15})
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.mulExpr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.mulExpr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.mulExpr() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_addExpr(t *testing.T) {
	t.Parallel()

	var s0 = "+v1_addExpr"
	var s1 = "+678"
	var s2 = "+6.78"
	var s3 = "-v2_addExpr"
	var s4 = "-15"
	var s5 = "-1.2"
	var s6 = "+"
	var s7 = "+ ++"
	var s8 = "+name"
	var s9 = "+\"string\""
	var s10 = "-"
	var s11 = "- ++"
	var s12 = "-name"
	var s13 = "-\"string\""

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{s0, fields{&s0, 0, Value{t: Integer, i: 345}}, Value{t: Floating, f: 346.234}, false},
		{"I+I", fields{&s1, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1023}, false},
		{"I+F", fields{&s2, 0, Value{t: Integer, i: 345}}, Value{t: Floating, f: 351.78}, false},
		{"F+I", fields{&s1, 0, Value{t: Floating, f: 3.4}}, Value{t: Floating, f: 681.4}, false},
		{"F+F", fields{&s2, 0, Value{t: Floating, f: 3.4}}, Value{t: Floating, f: 10.18}, false},
		{s3, fields{&s3, 0, Value{t: Integer, i: 345}}, Value{t: Floating, f: 343.5}, false},
		{"I-I", fields{&s4, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 330}, false},
		{"I-F", fields{&s5, 0, Value{t: Integer, i: 345}}, Value{t: Floating, f: 343.8}, false},
		{"F-I", fields{&s4, 0, Value{t: Floating, f: 3.45}}, Value{t: Floating, f: -11.55}, false},
		{"F-F", fields{&s5, 0, Value{t: Floating, f: 3.4}}, Value{t: Floating, f: 2.2}, false},
		{s6, fields{&s6, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{s7, fields{&s7, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s0, fields{&s0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{s8, fields{&s8, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{s9, fields{&s9, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{s10, fields{&s10, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{s11, fields{&s11, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s3, fields{&s3, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{s12, fields{&s12, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{s13, fields{&s13, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
	}
	SetVar("v1_addExpr", Value{t: Floating, f: 1.234})
	SetVar("v2_addExpr", Value{t: Floating, f: 1.5})
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.addExpr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.addExpr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.addExpr() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_shiftExpr(t *testing.T) {
	t.Parallel()

	var s0 = "<<v1_shiftExpr"
	var s1 = "<<7"
	var s3 = ">>v1_shiftExpr"
	var s4 = ">>1"
	var s5 = "<<"
	var s6 = "<< ++"
	var s7 = "<<name"
	var s8 = "<<\"string\""
	var s9 = ">>"
	var s10 = ">> ++"
	var s11 = ">>name"
	var s12 = ">>\"string\""

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"345" + s0, fields{&s0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 2760}, false},
		{"345" + s1, fields{&s1, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 44160}, false},
		{"345" + s3, fields{&s3, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 43}, false},
		{"345" + s4, fields{&s4, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 172}, false},
		{"345" + s5, fields{&s5, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{"345" + s6, fields{&s6, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s0, fields{&s0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s7, fields{&s7, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s8, fields{&s8, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"345" + s9, fields{&s9, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{"345" + s10, fields{&s10, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s3, fields{&s3, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s11, fields{&s11, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s12, fields{&s12, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
	}
	SetVar("v1_shiftExpr", Value{t: Integer, i: 3})
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.shiftExpr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.shiftExpr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.shiftExpr() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_relExpr(t *testing.T) {
	t.Parallel()

	var s000 = "<v1_relExpr"
	var s001 = "<7"
	var s002 = "<789"
	var s003 = "<345"
	var s004 = "<v2_relExpr"
	var s005 = "<7.1"
	var s006 = "<789.1"
	var s007 = "<345.0"
	var s008 = "<v1_relExpr"
	var s009 = "<7"
	var s010 = "<789"
	var s011 = "<345"
	var s012 = "<v2_relExpr"
	var s013 = "<7.1"
	var s014 = "<789.1"
	var s015 = "<345.0"
	var s016 = "<"
	var s017 = "< ++"
	var s018 = "<name"
	var s019 = "<\"string\""

	var s100 = "<=v1_relExpr"
	var s101 = "<=7"
	var s102 = "<=789"
	var s103 = "<=345"
	var s104 = "<=v2_relExpr"
	var s105 = "<=7.1"
	var s106 = "<=789.1"
	var s107 = "<=345.0"
	var s108 = "<=v1_relExpr"
	var s109 = "<=7"
	var s110 = "<=789"
	var s111 = "<=345"
	var s112 = "<=v2_relExpr"
	var s113 = "<=7.1"
	var s114 = "<=789.1"
	var s115 = "<=345.0"
	var s116 = "<="
	var s117 = "<= ++"
	var s118 = "<=name"
	var s119 = "<=\"string\""

	var s200 = ">v1_relExpr"
	var s201 = ">7"
	var s202 = ">789"
	var s203 = ">345"
	var s204 = ">v2_relExpr"
	var s205 = ">7.1"
	var s206 = ">789.1"
	var s207 = ">345.0"
	var s208 = ">v1_relExpr"
	var s209 = ">7"
	var s210 = ">789"
	var s211 = ">345"
	var s212 = ">v2_relExpr"
	var s213 = ">7.1"
	var s214 = ">789.1"
	var s215 = ">345.0"
	var s216 = ">"
	var s217 = "> ++"
	var s218 = ">name"
	var s219 = ">\"string\""

	var s300 = ">=v1_relExpr"
	var s301 = ">=7"
	var s302 = ">=789"
	var s303 = ">=345"
	var s304 = ">=v2_relExpr"
	var s305 = ">=7.1"
	var s306 = ">=789.1"
	var s307 = ">=345.0"
	var s308 = ">=v1_relExpr"
	var s309 = ">=7"
	var s310 = ">=789"
	var s311 = ">=345"
	var s312 = ">=v2_relExpr"
	var s313 = ">=7.1"
	var s314 = ">=789.1"
	var s315 = ">=345.0"
	var s316 = ">="
	var s317 = ">= ++"
	var s318 = ">=name"
	var s319 = ">=\"string\""

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"345" + s000, fields{&s000, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s001, fields{&s001, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s002, fields{&s002, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s003, fields{&s003, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s004, fields{&s004, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s005, fields{&s005, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s006, fields{&s006, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s007, fields{&s007, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s008, fields{&s008, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s009, fields{&s009, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s010, fields{&s010, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s011, fields{&s011, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s012, fields{&s012, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s013, fields{&s013, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s014, fields{&s014, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s015, fields{&s015, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345" + s016, fields{&s016, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{"345" + s017, fields{&s017, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s000, fields{&s000, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s018, fields{&s018, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s019, fields{&s019, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},

		{"345" + s100, fields{&s100, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s101, fields{&s101, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s102, fields{&s102, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s103, fields{&s103, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s104, fields{&s104, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s105, fields{&s105, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s106, fields{&s106, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s107, fields{&s107, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s108, fields{&s108, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s109, fields{&s109, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s110, fields{&s110, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s111, fields{&s111, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s112, fields{&s112, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s113, fields{&s113, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s114, fields{&s114, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s115, fields{&s115, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345" + s116, fields{&s116, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{"345" + s117, fields{&s117, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s100, fields{&s100, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s118, fields{&s118, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s119, fields{&s119, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},

		{"345" + s200, fields{&s200, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s201, fields{&s201, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s202, fields{&s202, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s203, fields{&s203, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s204, fields{&s204, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s205, fields{&s205, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s206, fields{&s206, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s207, fields{&s207, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s208, fields{&s208, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s209, fields{&s209, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s210, fields{&s210, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s211, fields{&s211, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s212, fields{&s212, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s213, fields{&s213, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s214, fields{&s214, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s215, fields{&s215, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345" + s216, fields{&s216, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{"345" + s217, fields{&s217, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s200, fields{&s200, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s218, fields{&s218, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s219, fields{&s219, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},

		{"345" + s300, fields{&s300, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s301, fields{&s301, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s302, fields{&s302, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s303, fields{&s303, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s304, fields{&s304, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s305, fields{&s305, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s306, fields{&s306, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s307, fields{&s307, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s308, fields{&s308, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s309, fields{&s309, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s310, fields{&s310, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s311, fields{&s311, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s312, fields{&s312, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s313, fields{&s313, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s314, fields{&s314, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s315, fields{&s315, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345" + s316, fields{&s316, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{"345" + s317, fields{&s317, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s300, fields{&s300, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s318, fields{&s318, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s319, fields{&s319, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
	}
	SetVar("v1_relExpr", Value{t: Integer, i: 3})
	SetVar("v2_relExpr", Value{t: Floating, f: 1.5})
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.relExpr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.relExpr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.relExpr() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_equExpr(t *testing.T) {
	t.Parallel()

	var s000 = "==v1_equExpr"
	var s001 = "==7"
	var s002 = "==789"
	var s003 = "==345"
	var s004 = "==v2_equExpr"
	var s005 = "==7.1"
	var s006 = "==789.1"
	var s007 = "==345.0"
	var s008 = "==v1_equExpr"
	var s009 = "==7"
	var s010 = "==789"
	var s011 = "==345"
	var s012 = "==v2_equExpr"
	var s013 = "==7.1"
	var s014 = "==789.1"
	var s015 = "==345.0"
	var s016 = "=="
	var s017 = "== ++"
	var s018 = "==name"
	var s019 = "==\"string\""

	var s100 = "!=v1_equExpr"
	var s101 = "!=7"
	var s102 = "!=789"
	var s103 = "!=345"
	var s104 = "!=v2_equExpr"
	var s105 = "!=7.1"
	var s106 = "!=789.1"
	var s107 = "!=345.0"
	var s108 = "!=v1_equExpr"
	var s109 = "!=7"
	var s110 = "!=789"
	var s111 = "!=345"
	var s112 = "!=v2_equExpr"
	var s113 = "!=7.1"
	var s114 = "!=789.1"
	var s115 = "!=345.0"
	var s116 = "!="
	var s117 = "!= ++"
	var s118 = "!=name"
	var s119 = "!=\"string\""

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"345" + s000, fields{&s000, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s001, fields{&s001, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s002, fields{&s002, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s003, fields{&s003, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s004, fields{&s004, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s005, fields{&s005, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s006, fields{&s006, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s007, fields{&s007, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s008, fields{&s008, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s009, fields{&s009, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s010, fields{&s010, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s011, fields{&s011, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s012, fields{&s012, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s013, fields{&s013, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s014, fields{&s014, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s015, fields{&s015, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345" + s016, fields{&s016, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{"345" + s017, fields{&s017, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s000, fields{&s000, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s018, fields{&s018, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s019, fields{&s019, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},

		{"345" + s100, fields{&s100, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s101, fields{&s101, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s102, fields{&s102, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s103, fields{&s103, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345" + s104, fields{&s104, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s105, fields{&s105, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s106, fields{&s106, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345" + s107, fields{&s107, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s108, fields{&s108, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s109, fields{&s109, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s110, fields{&s110, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s111, fields{&s111, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0" + s112, fields{&s112, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s113, fields{&s113, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s114, fields{&s114, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0" + s115, fields{&s115, 0, Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345" + s116, fields{&s116, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{"345" + s117, fields{&s117, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s100, fields{&s100, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s118, fields{&s118, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s119, fields{&s119, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
	}
	SetVar("v1_equExpr", Value{t: Integer, i: 3})
	SetVar("v2_equExpr", Value{t: Floating, f: 1.5})
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.equExpr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.equExpr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.equExpr() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_andExpr(t *testing.T) {
	t.Parallel()

	var s0 = "&v1_andExpr"
	var s1 = "&0xaf5f0ff0"
	var s2 = "&"
	var s3 = "& ++"
	var s4 = "&name"
	var s5 = "&\"string\""

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"0x55aa00ff" + s0, fields{&s0, 0, Value{t: Integer, i: 0x55aa00ff}}, Value{t: Integer, i: 0x050A00F0}, false},
		{"0x55aa00ff" + s1, fields{&s1, 0, Value{t: Integer, i: 0x55aa00ff}}, Value{t: Integer, i: 0x050A00F0}, false},
		{"345" + s2, fields{&s2, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{"345" + s3, fields{&s3, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s0, fields{&s0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s4, fields{&s4, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s5, fields{&s5, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
	}
	SetVar("v1_andExpr", Value{t: Integer, i: 0xaf5f0ff0})
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.andExpr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.andExpr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.andExpr() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_xorExpr(t *testing.T) {
	t.Parallel()

	var s0 = "^v1_xorExpr"
	var s1 = "^0xaf5f0ff0"
	var s2 = "^"
	var s3 = "^ ++"
	var s4 = "^name"
	var s5 = "^\"string\""

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"0x55aa00ff" + s0, fields{&s0, 0, Value{t: Integer, i: 0x55aa00ff}}, Value{t: Integer, i: 0xFAF50F0F}, false},
		{"0x55aa00ff" + s1, fields{&s1, 0, Value{t: Integer, i: 0x55aa00ff}}, Value{t: Integer, i: 0xFAF50F0F}, false},
		{"345" + s2, fields{&s2, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{"345" + s3, fields{&s3, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s0, fields{&s0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s4, fields{&s4, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s5, fields{&s5, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
	}
	SetVar("v1_xorExpr", Value{t: Integer, i: 0xaf5f0ff0})
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.xorExpr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.xorExpr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.xorExpr() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_orExpr(t *testing.T) {
	t.Parallel()

	var s0 = "|v1_orExpr"
	var s1 = "|0xaf5f0ff0"
	var s2 = "|"
	var s3 = "| ++"
	var s4 = "|name"
	var s5 = "|\"string\""

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"0x55aa00ff" + s0, fields{&s0, 0, Value{t: Integer, i: 0x55aa00ff}}, Value{t: Integer, i: 0xFFFF0FFF}, false},
		{"0x55aa00ff" + s1, fields{&s1, 0, Value{t: Integer, i: 0x55aa00ff}}, Value{t: Integer, i: 0xFFFF0FFF}, false},
		{"345" + s2, fields{&s2, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{"345" + s3, fields{&s3, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s0, fields{&s0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s4, fields{&s4, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s5, fields{&s5, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
	}
	SetVar("v1_orExpr", Value{t: Integer, i: 0xaf5f0ff0})
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.orExpr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.orExpr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.orExpr() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_logAndExpr(t *testing.T) {
	t.Parallel()

	var s0 = "&&v1_logAndExpr"
	var s1 = "&&0"
	var s2 = "&&1"
	var s3 = "&&"
	var s4 = "&& ++"
	var s5 = "&&name"
	var s6 = "&&\"string\""

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"1" + s0, fields{&s0, 0, Value{t: Integer, i: 1}}, Value{t: Integer, i: 1}, false},
		{"0" + s1, fields{&s1, 0, Value{t: Integer, i: 0}}, Value{t: Integer, i: 0}, false},
		{"0" + s2, fields{&s2, 0, Value{t: Integer, i: 0}}, Value{t: Integer, i: 0}, false},
		{"1" + s1, fields{&s1, 0, Value{t: Integer, i: 1}}, Value{t: Integer, i: 0}, false},
		{"1" + s2, fields{&s2, 0, Value{t: Integer, i: 1}}, Value{t: Integer, i: 1}, false},
		{"345" + s3, fields{&s3, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{"345" + s4, fields{&s4, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s0, fields{&s0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s5, fields{&s5, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s6, fields{&s6, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
	}
	SetVar("v1_logAndExpr", Value{t: Integer, i: 1})
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.logAndExpr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.logAndExpr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.logAndExpr() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_logOrExpr(t *testing.T) {
	t.Parallel()

	var s0 = "||v1_logOrExpr"
	var s1 = "||0"
	var s2 = "||1"
	var s3 = "||"
	var s4 = "|| ++"
	var s5 = "||name"
	var s6 = "||\"string\""

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"1" + s0, fields{&s0, 0, Value{t: Integer, i: 1}}, Value{t: Integer, i: 1}, false},
		{"0" + s1, fields{&s1, 0, Value{t: Integer, i: 0}}, Value{t: Integer, i: 0}, false},
		{"0" + s2, fields{&s2, 0, Value{t: Integer, i: 0}}, Value{t: Integer, i: 1}, false},
		{"1" + s1, fields{&s1, 0, Value{t: Integer, i: 1}}, Value{t: Integer, i: 1}, false},
		{"1" + s2, fields{&s2, 0, Value{t: Integer, i: 1}}, Value{t: Integer, i: 1}, false},
		{"345" + s3, fields{&s3, 0, Value{t: Integer, i: 345}}, Value{t: Nix}, true},
		{"345" + s4, fields{&s4, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s0, fields{&s0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s5, fields{&s5, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s6, fields{&s6, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
	}
	SetVar("v1_logOrExpr", Value{t: Integer, i: 1})
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.logOrExpr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.logOrExpr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.logOrExpr() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_condExpr(t *testing.T) {
	t.Parallel()

	var s0 = "?v1_condExpr:v2_condExpr"
	var s1 = "?2:3"
	var s2 = "?"
	var s3 = "?1"
	var s4 = "? ++"
	var s5 = "?name"
	var s6 = "?\"string\""
	var s7 = "?1:"
	var s8 = "?1: ++"
	var s9 = "?v1_condExpr:name"

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"0" + s9, fields{&s9, 0, Value{t: Integer, i: 0}}, Value{t: Identifier, s: "name"}, true},
		{"1" + s0, fields{&s0, 0, Value{t: Integer, i: 1}}, Value{t: Integer, i: 2}, false},
		{"0" + s0, fields{&s0, 0, Value{t: Integer, i: 0}}, Value{t: Integer, i: 3}, false},
		{"1" + s1, fields{&s1, 0, Value{t: Integer, i: 1}}, Value{t: Integer, i: 2}, false},
		{"0" + s1, fields{&s1, 0, Value{t: Integer, i: 0}}, Value{t: Integer, i: 3}, false},
		{"1.23" + s1, fields{&s1, 0, Value{t: Floating, f: 1.23}}, Value{t: Integer, i: 2}, false},
		{"0.0" + s1, fields{&s1, 0, Value{t: Floating, f: 0.0}}, Value{t: Integer, i: 3}, false},
		{"1" + s2, fields{&s2, 0, Value{t: Integer, i: 1}}, Value{t: Integer, i: 1}, true},
		{"345" + s3, fields{&s3, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"345" + s4, fields{&s4, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"name" + s0, fields{&s0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s5, fields{&s5, 0, Value{t: Integer, i: 345}}, Value{t: Identifier, s: "name"}, true},
		{"345" + s6, fields{&s6, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"1" + s7, fields{&s7, 0, Value{t: Integer, i: 1}}, Value{t: Integer, i: 1}, true},
		{"345" + s8, fields{&s8, 0, Value{t: Integer, i: 345}}, Value{t: AddAdd}, true},
		{"\"string\"" + s1, fields{&s1, 0, Value{t: String, s: "string"}}, Value{t: String, s: "string"}, true},
	}
	for _, tt := range tests {
		tt := tt
		SetVar("v1_condExpr", Value{t: Integer, i: 2})
		SetVar("v2_condExpr", Value{t: Integer, i: 3})
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.condExpr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.condExpr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.condExpr() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_asnExpr(t *testing.T) {
	t.Parallel()

	var shl0 = "<<=v1_asnExpr"
	var shl1 = "<<=7"
	var shl2 = "<<="
	var shl3 = "<<= ++"
	var shl4 = "<<=name"
	var shl5 = "<<=\"string\""

	var shr0 = ">>=v1_asnExpr"
	var shr1 = ">>=1"
	var shr2 = ">>="
	var shr3 = ">>= ++"
	var shr4 = ">>=name"
	var shr5 = ">>=\"string\""

	var plus0 = "+=v1_asnExpr"
	var plus1 = "+=1"
	var plus2 = "+="
	var plus3 = "+= ++"
	var plus4 = "+=name"
	var plus5 = "+=\"string\""

	var minus0 = "-=v1_asnExpr"
	var minus1 = "-=1"
	var minus2 = "-="
	var minus3 = "-= ++"
	var minus4 = "-=name"
	var minus5 = "-=\"string\""

	var or0 = "|=v2_asnExpr"
	var or1 = "|=0xaf5f0ff0"
	var or2 = "|="
	var or3 = "|= ++"
	var or4 = "|=name"
	var or5 = "|=\"string\""

	var and0 = "&=v2_asnExpr"
	var and1 = "&=0xaf5f0ff0"
	var and2 = "&="
	var and3 = "&= ++"
	var and4 = "&=name"
	var and5 = "&=\"string\""

	var xor0 = "^=v2_asnExpr"
	var xor1 = "^=0xaf5f0ff0"
	var xor2 = "^="
	var xor3 = "^= ++"
	var xor4 = "^=name"
	var xor5 = "^=\"string\""

	var mul0 = "*=v1_asnExpr"
	var mul1 = "*=7"
	var mul2 = "*="
	var mul3 = "*= ++"
	var mul4 = "*=name"
	var mul5 = "*=\"string\""

	var div0 = "/=v1_asnExpr"
	var div1 = "/=7"
	var div2 = "/="
	var div3 = "/= ++"
	var div4 = "/=name"
	var div5 = "/=\"string\""

	var mod0 = "%=v3_asnExpr"
	var mod1 = "%=14"
	var mod2 = "%="
	var mod3 = "%= ++"
	var mod4 = "%=name"
	var mod5 = "%=\"string\""

	var ass0 = "=v1_asnExpr"
	var ass1 = "=345"
	var ass2 = "="
	var ass3 = "= ++"
	var ass4 = "=name"
	var ass5 = "=\"string\""

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"345" + shl0, fields{&shl0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v00_asnExpr" + shl0, fields{&shl0, 0, Value{t: Identifier, s: "v00_asnExpr"}}, Value{t: Integer, i: 2760}, false},
		{"v01_asnExpr" + shl1, fields{&shl1, 0, Value{t: Identifier, s: "v01_asnExpr"}}, Value{t: Integer, i: 44160}, false},
		{"345" + shl0, fields{&shl0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v0_asnExpr" + shl2, fields{&shl2, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},
		{"v0_asnExpr" + shl3, fields{&shl3, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: AddAdd}, true},
		{"name" + shl0, fields{&shl0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + shl4, fields{&shl4, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + shl5, fields{&shl5, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},

		{"345" + shr0, fields{&shr0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v02_asnExpr" + shr0, fields{&shr0, 0, Value{t: Identifier, s: "v02_asnExpr"}}, Value{t: Integer, i: 43}, false},
		{"v03_asnExpr" + shr1, fields{&shr1, 0, Value{t: Identifier, s: "v03_asnExpr"}}, Value{t: Integer, i: 172}, false},
		{"345" + shr0, fields{&shr0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v0_asnExpr" + shr2, fields{&shr2, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},
		{"v0_asnExpr" + shr3, fields{&shr3, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: AddAdd}, true},
		{"name" + shr0, fields{&shr0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + shr4, fields{&shr4, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + shr5, fields{&shr5, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},

		{"345" + plus0, fields{&plus0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v04_asnExpr" + plus0, fields{&plus0, 0, Value{t: Identifier, s: "v04_asnExpr"}}, Value{t: Integer, i: 348}, false},
		{"v05_asnExpr" + plus1, fields{&plus1, 0, Value{t: Identifier, s: "v05_asnExpr"}}, Value{t: Integer, i: 346}, false},
		{"345" + plus0, fields{&plus0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v0_asnExpr" + plus2, fields{&plus2, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},
		{"v0_asnExpr" + plus3, fields{&plus3, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: AddAdd}, true},
		{"name" + plus0, fields{&plus0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + plus4, fields{&plus4, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + plus5, fields{&plus5, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},

		{"345" + minus0, fields{&minus0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v06_asnExpr" + minus0, fields{&minus0, 0, Value{t: Identifier, s: "v06_asnExpr"}}, Value{t: Integer, i: 342}, false},
		{"v07_asnExpr" + minus1, fields{&minus1, 0, Value{t: Identifier, s: "v07_asnExpr"}}, Value{t: Integer, i: 344}, false},
		{"345" + minus0, fields{&minus0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v0_asnExpr" + minus2, fields{&minus2, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},
		{"v0_asnExpr" + minus3, fields{&minus3, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: AddAdd}, true},
		{"name" + minus0, fields{&minus0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + minus4, fields{&minus4, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + minus5, fields{&minus5, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},

		{"345" + or0, fields{&or0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v08_asnExpr" + or0, fields{&or0, 0, Value{t: Identifier, s: "v08_asnExpr"}}, Value{t: Integer, i: 0xffff0fff}, false},
		{"v09_asnExpr" + or1, fields{&or1, 0, Value{t: Identifier, s: "v09_asnExpr"}}, Value{t: Integer, i: 0xffff0fff}, false},
		{"345" + or0, fields{&or0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v0_asnExpr" + or2, fields{&or2, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},
		{"v0_asnExpr" + or3, fields{&or3, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: AddAdd}, true},
		{"name" + or0, fields{&or0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + or4, fields{&or4, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + or5, fields{&or5, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},

		{"345" + and0, fields{&and0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v010_asnExpr" + and0, fields{&and0, 0, Value{t: Identifier, s: "v010_asnExpr"}}, Value{t: Integer, i: 0x050a00f0}, false},
		{"v011_asnExpr" + and1, fields{&and1, 0, Value{t: Identifier, s: "v011_asnExpr"}}, Value{t: Integer, i: 0x050a00f0}, false},
		{"345" + and0, fields{&and0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v0_asnExpr" + and2, fields{&and2, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},
		{"v0_asnExpr" + and3, fields{&and3, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: AddAdd}, true},
		{"name" + and0, fields{&and0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + and4, fields{&and4, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + and5, fields{&and5, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},

		{"345" + xor0, fields{&xor0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v012_asnExpr" + xor0, fields{&xor0, 0, Value{t: Identifier, s: "v012_asnExpr"}}, Value{t: Integer, i: 0xFAF50F0F}, false},
		{"v013_asnExpr" + xor1, fields{&xor1, 0, Value{t: Identifier, s: "v013_asnExpr"}}, Value{t: Integer, i: 0xFAF50F0F}, false},
		{"345" + xor0, fields{&xor0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v0_asnExpr" + xor2, fields{&xor2, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},
		{"v0_asnExpr" + xor3, fields{&xor3, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: AddAdd}, true},
		{"name" + xor0, fields{&xor0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + xor4, fields{&xor4, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + xor5, fields{&xor5, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},

		{"345" + mul0, fields{&mul0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v014_asnExpr" + mul0, fields{&mul0, 0, Value{t: Identifier, s: "v014_asnExpr"}}, Value{t: Integer, i: 1035}, false},
		{"v015_asnExpr" + mul1, fields{&mul1, 0, Value{t: Identifier, s: "v015_asnExpr"}}, Value{t: Integer, i: 2415}, false},
		{"345" + mul0, fields{&mul0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v0_asnExpr" + mul2, fields{&mul2, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},
		{"v0_asnExpr" + mul3, fields{&mul3, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: AddAdd}, true},
		{"name" + mul0, fields{&mul0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + mul4, fields{&mul4, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + mul5, fields{&mul5, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},

		{"345" + div0, fields{&div0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v016_asnExpr" + div0, fields{&div0, 0, Value{t: Identifier, s: "v016_asnExpr"}}, Value{t: Integer, i: 115}, false},
		{"v017_asnExpr" + div1, fields{&div1, 0, Value{t: Identifier, s: "v017_asnExpr"}}, Value{t: Integer, i: 49}, false},
		{"345" + div0, fields{&div0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v0_asnExpr" + div2, fields{&div2, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},
		{"v0_asnExpr" + div3, fields{&div3, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: AddAdd}, true},
		{"name" + div0, fields{&div0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + div4, fields{&div4, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + div5, fields{&div5, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},

		{"345" + mod0, fields{&mod0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v018_asnExpr" + mod0, fields{&mod0, 0, Value{t: Identifier, s: "v018_asnExpr"}}, Value{t: Integer, i: 9}, false},
		{"v019_asnExpr" + mod1, fields{&mod1, 0, Value{t: Identifier, s: "v019_asnExpr"}}, Value{t: Integer, i: 9}, false},
		{"345" + mod0, fields{&mod0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v0_asnExpr" + mod2, fields{&mod2, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},
		{"v0_asnExpr" + mod3, fields{&mod3, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: AddAdd}, true},
		{"name" + mod0, fields{&mod0, 0, Value{t: Identifier, s: "name"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + mod4, fields{&mod4, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "name"}, true},
		{"v0_asnExpr" + mod5, fields{&mod5, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},

		{"345" + ass0, fields{&ass0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v020_asnExpr" + ass0, fields{&ass0, 0, Value{t: Identifier, s: "v020_asnExpr"}}, Value{t: Integer, i: 3}, false},
		{"v021_asnExpr" + ass1, fields{&ass1, 0, Value{t: Identifier, s: "v021_asnExpr"}}, Value{t: Integer, i: 345}, false},
		{"345" + ass0, fields{&ass0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v0_asnExpr" + ass2, fields{&ass2, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "v0_asnExpr"}, true},
		{"v0_asnExpr" + ass3, fields{&ass3, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: AddAdd}, true},
		{"name" + ass0, fields{&ass0, 0, Value{t: Identifier, s: "name"}}, Value{t: Integer, i: 3}, true},
		{"v0_asnExpr" + ass4, fields{&ass4, 0, Value{t: Identifier, s: "v0_asnExpr"}}, Value{t: Identifier, s: "name"}, true},
		{"v022_asnExpr" + ass5, fields{&ass5, 0, Value{t: Identifier, s: "v022_asnExpr"}}, Value{t: String, s: "string"}, false},
	}
	for _, tt := range tests {
		tt := tt
		SetVar("v00_asnExpr", Value{t: Integer, i: 345})
		SetVar("v01_asnExpr", Value{t: Integer, i: 345})
		SetVar("v02_asnExpr", Value{t: Integer, i: 345})
		SetVar("v03_asnExpr", Value{t: Integer, i: 345})
		SetVar("v04_asnExpr", Value{t: Integer, i: 345})
		SetVar("v05_asnExpr", Value{t: Integer, i: 345})
		SetVar("v06_asnExpr", Value{t: Integer, i: 345})
		SetVar("v07_asnExpr", Value{t: Integer, i: 345})
		SetVar("v08_asnExpr", Value{t: Integer, i: 0x55aa00ff})
		SetVar("v09_asnExpr", Value{t: Integer, i: 0x55aa00ff})
		SetVar("v010_asnExpr", Value{t: Integer, i: 0x55aa00ff})
		SetVar("v011_asnExpr", Value{t: Integer, i: 0x55aa00ff})
		SetVar("v012_asnExpr", Value{t: Integer, i: 0x55aa00ff})
		SetVar("v013_asnExpr", Value{t: Integer, i: 0x55aa00ff})
		SetVar("v014_asnExpr", Value{t: Integer, i: 345})
		SetVar("v015_asnExpr", Value{t: Integer, i: 345})
		SetVar("v016_asnExpr", Value{t: Integer, i: 345})
		SetVar("v017_asnExpr", Value{t: Integer, i: 345})
		SetVar("v018_asnExpr", Value{t: Integer, i: 345})
		SetVar("v019_asnExpr", Value{t: Integer, i: 345})
		SetVar("v020_asnExpr", Value{t: Integer, i: 345})
		SetVar("v021_asnExpr", Value{t: Integer, i: 345})
		SetVar("v022_asnExpr", Value{t: Integer, i: 345})
		SetVar("v0_asnExpr", Value{t: Integer, i: 345})
		SetVar("v1_asnExpr", Value{t: Integer, i: 3})
		SetVar("v2_asnExpr", Value{t: Integer, i: 0xaf5f0ff0})
		SetVar("v3_asnExpr", Value{t: Integer, i: 14})
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.asnExpr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.asnExpr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			got.v = nil
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.asnExpr() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExpression_expression(t *testing.T) {
	t.Parallel()

	var s0 = ",v1_expExpr"
	var s1 = ",7"
	var s2 = ","
	var s3 = ", ++"

	type fields struct {
		in   *string
		pos  int
		next Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"345" + s0, fields{&s0, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, false},
		{"v0_expExpr" + s0, fields{&s0, 0, Value{t: Identifier, s: "v0_expExpr"}}, Value{t: Integer, i: 1}, false},
		{"345" + s1, fields{&s1, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, false},
		{"v0_expExpr" + s1, fields{&s1, 0, Value{t: Identifier, s: "v0_expExpr"}}, Value{t: Integer, i: 1}, false},
		{"345" + s2, fields{&s2, 0, Value{t: Integer, i: 345}}, Value{t: Integer, i: 345}, true},
		{"v0_expExpr" + s3, fields{&s3, 0, Value{t: Identifier, s: "v0_expExpr"}}, Value{t: Integer, i: 1}, true},
	}
	for _, tt := range tests {
		tt := tt
		SetVar("v0_expExpr", Value{t: Integer, i: 1})
		SetVar("v1_expExpr", Value{t: Integer, i: 345})
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ex := &Expression{
				in:   tt.fields.in,
				pos:  tt.fields.pos,
				next: tt.fields.next,
			}
			got, err := ex.expression()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.expression() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.expression() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
