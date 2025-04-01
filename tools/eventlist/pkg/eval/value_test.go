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
	"eventlist/pkg/elf"
	"reflect"
	"testing"
)

func TestValue_Compose(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		i int64
		f float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		t Token
		i int64
		f float64
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Value
	}{
		{"test", fields{t: Integer, i: 123, f: 1.23, s: "abc"}, args{t: Floating, i: 789, f: 7.89, s: "xxx"}, Value{t: Floating, i: 789, f: 7.89, s: "xxx"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.i,
				f: tt.fields.f,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			v.Compose(tt.args.t, tt.args.i, tt.args.f, tt.args.s)
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Compose() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_getValue(t *testing.T) { //nolint:golint,paralleltest
	vari := Variable{"v1_getValue", Value{t: Integer, i: 456}}

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name    string
		fields  fields
		clear   bool
		want    Value
		wantErr bool
	}{
		{"test_normal", fields{t: Integer, v: &vari}, false, Value{t: Integer, i: 789}, false},
		{"test_error", fields{}, true, Value{}, true},
		{"test_error1", fields{t: Integer, v: &vari}, true, Value{}, true},
	}

	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if tt.clear {
				ClearNames()
			} else {
				SetVar("v1_getValue", Value{t: Integer, i: 789})
			}
			got, err := v.getValue()
			if (err != nil) != tt.wantErr {
				t.Errorf("Value.getValue() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value.getValue() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestValue_setValue(t *testing.T) { //nolint:golint,paralleltest
	vari := Variable{"v1_setValue", Value{t: Integer, i: 456}}
	val1 := Value{t: Integer, i: 123}

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		clear   bool
		want    *Value
		wantErr bool
	}{
		{"test_normal", fields{t: Identifier, v: &vari}, args{&val1}, false, &val1, false},
		{"test_error", fields{t: Identifier}, args{&val1}, true, &Value{}, true},
		{"test_error1", fields{t: Identifier, v: &vari}, args{&val1}, true, &Value{}, true},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if tt.clear {
				ClearNames()
			} else {
				SetVar("v1_setValue", Value{t: Integer, i: 789})
			}
			var err error
			if err = v.setValue(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.setValue() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			var got *Variable
			if err == nil {
				if got, err = GetVar("v1_setValue"); err != nil {
					t.Errorf("Value.setValue() %s error = %v", tt.name, err)
				}
				if !reflect.DeepEqual(got.v, *tt.want) {
					t.Errorf("Value.getValue() %s = %v, want %v", tt.name, got, tt.want)
				}
			}
		})
	}
}

func TestValue_addList(t *testing.T) {
	t.Parallel()

	v1 := Value{t: Integer, i: 1}

	type fields struct {
		t Token
		i int64
		f float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"add1", fields{}, args{v1}, false},
		{"add2", fields{t: List, l: []Value{{t: Integer, i: 1}}}, args{v1}, false},
		{"err", fields{t: Integer}, args{v1}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.i,
				f: tt.fields.f,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.addList(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.addList() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestValue_GetInt(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{"test_int", fields{t: Integer, I: 123}, 123},
		{"test_-int", fields{t: Integer, I: -123}, -123},
		{"test_float", fields{t: Floating, F: 45.67}, 45},
		{"test_-float", fields{t: Floating, F: -45.67}, -45},
		{"test_nix", fields{t: String, s: "abc"}, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if got := v.GetInt(); got != tt.want {
				t.Errorf("Value.GetInt() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestValue_GetUInt(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{"test_int", fields{t: Integer, I: 123}, 123},
		{"test_float", fields{t: Floating, F: 45.67}, 45},
		{"test_nix", fields{t: String, s: "abc"}, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if got := v.GetUInt(); got != tt.want {
				t.Errorf("Value.GetUInt() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestValue_GetFloat(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{"test_int", fields{t: Integer, I: 123}, 123},
		{"test_float", fields{t: Floating, F: 45.67}, 45.67},
		{"test_nix", fields{t: String, s: "abc"}, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if got := v.GetFloat(); got != tt.want {
				t.Errorf("Value.GetFloat() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestValue_GetList(t *testing.T) {
	t.Parallel()

	v1 := []Value{{t: Integer, i: 4711}}

	type fields struct {
		t Token
		i int64
		f float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name   string
		fields fields
		want   []Value
	}{
		{"nil", fields{}, nil},
		{"4711", fields{t: List, l: []Value{{t: Integer, i: 4711}}}, v1},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.i,
				f: tt.fields.f,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if got := v.GetList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value.GetList() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestValue_IsInteger(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"N", fields{t: Nix}, false},
		{"C", fields{t: Integer}, true},
		{"F", fields{t: Floating}, false},
		{"S", fields{t: String}, false},
		{"L", fields{t: List}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if got := v.IsInteger(); got != tt.want {
				t.Errorf("Value.IsInteger() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestValue_IsFloating(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"N", fields{t: Nix}, false},
		{"C", fields{t: Integer}, false},
		{"F", fields{t: Floating}, true},
		{"S", fields{t: String}, false},
		{"L", fields{t: List}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if got := v.IsFloating(); got != tt.want {
				t.Errorf("Value.IsFloating() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestValue_IsString(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"N", fields{t: Nix}, false},
		{"C", fields{t: Integer}, false},
		{"F", fields{t: Floating}, false},
		{"S", fields{t: String}, true},
		{"L", fields{t: List}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if got := v.IsString(); got != tt.want {
				t.Errorf("Value.IsString() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestValue_IsList(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		i int64
		f float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"N", fields{t: Nix}, false},
		{"C", fields{t: Integer}, false},
		{"F", fields{t: Floating}, false},
		{"S", fields{t: String}, false},
		{"L", fields{t: List}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.i,
				f: tt.fields.f,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if got := v.IsList(); got != tt.want {
				t.Errorf("Value.IsList() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestValue_Function(t *testing.T) { //nolint:golint,paralleltest
	calcMemUsedArgs := Value{t: List, l: []Value{{t: Integer, i: 1}, {t: Integer, i: 2}, {t: Integer, i: 3}, {t: Integer, i: 4}}}
	calcMemUsedArgs1 := Value{t: List, l: []Value{{t: String}, {t: Integer, i: 2}, {t: Integer, i: 3}, {t: Integer, i: 4}}}
	getRegValArgs := Value{t: List, l: []Value{{t: String, s: "reg"}}}
	symbolExistsArgs := Value{t: List, l: []Value{{t: String, s: "LEDOn"}}}
	symbolExistsArgs1 := Value{t: List, l: []Value{{t: String, s: "xxxx"}}}

	elf.Symbols.Init("LEDOn", 0x38000178, 4)

	type fields struct {
		t Token
		i int64
		f float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"CalcMemUsed", fields{t: Identifier, s: "__CalcMemUsed"}, args{&calcMemUsedArgs}, Value{t: Integer, i: 0}, false},
		{"GetRegVal", fields{t: Identifier, s: "__GetRegVal"}, args{&getRegValArgs}, Value{t: Integer, i: 0}, false},
		{"SymbolExist", fields{t: Identifier, s: "__Symbol_exists"}, args{&symbolExistsArgs}, Value{t: Integer, i: 1}, false},
		{"SymbolExist1", fields{t: Identifier, s: "__Symbol_exists"}, args{&symbolExistsArgs1}, Value{t: Integer, i: 0}, false},
		{"FindSymbol", fields{t: Identifier, s: "__FindSymbol"}, args{&symbolExistsArgs}, Value{t: Integer, i: 1}, false},
		{"FindSymbol1", fields{t: Identifier, s: "__FindSymbol"}, args{&symbolExistsArgs1}, Value{t: Integer, i: 0}, false},
		{"offsetOf", fields{t: Identifier, s: "__Offset_of"}, args{&symbolExistsArgs}, Value{t: Integer, i: 0x38000178}, false},
		{"offsetOf1", fields{t: Identifier, s: "__Offset_of"}, args{&symbolExistsArgs1}, Value{t: Integer, i: 0}, false},
		{"sizeOf", fields{t: Identifier, s: "__size_of"}, args{&symbolExistsArgs}, Value{t: Integer, i: 4}, false},
		{"sizeOf1", fields{t: Identifier, s: "__size_of"}, args{&symbolExistsArgs1}, Value{t: Integer, i: 0}, false},
		{"NoId", fields{t: Nix}, args{&Value{}}, Value{}, true},
		{"NoList", fields{t: Identifier, s: "abc"}, args{&Value{}}, Value{t: Identifier, s: "abc"}, true},
		{"Nil", fields{t: Identifier, s: "abc"}, args{}, Value{t: Identifier, s: "abc"}, true},
		{"NoFct", fields{t: Identifier, s: "abc"}, args{&calcMemUsedArgs}, Value{t: Identifier, s: "abc"}, true},
		{"wrongCnt", fields{t: Identifier, s: "__CalcMemUsed"}, args{&getRegValArgs}, Value{t: Identifier, s: "__CalcMemUsed"}, true},
		{"wrongType", fields{t: Identifier, s: "__CalcMemUsed"}, args{&calcMemUsedArgs1}, Value{t: Identifier, s: "__CalcMemUsed"}, true},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			v := &Value{
				t: tt.fields.t,
				i: tt.fields.i,
				f: tt.fields.f,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Function(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.Function() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Function() %s = %v, want %v", tt.name, v, tt.want)
			}
		})
	}
}

func TestValue_Extract(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}

	type args struct {
		sz        uint32
		bigEndian bool
		off       uint32
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"Extract little-endian", fields{t: Integer, I: 0x1234567890ABCDEF}, args{sz: 4, bigEndian: false, off: 2}, Value{t: Integer, i: 0x90AB}, false},
		//{"Extract big-endian", fields{t: Integer, I: 0x1234567890ABCDEF}, args{sz: 4, bigEndian: true, off: 2}, Value{t: Integer, i: 0x5678}, false},
		{"Extract with offset 0", fields{t: Integer, I: 0x1234567890ABCDEF}, args{sz: 4, bigEndian: false, off: 0}, Value{t: Integer, i: 0x90ABCDEF}, false},
		//{"Extract with size 1", fields{t: Integer, I: 0x1234567890ABCDEF}, args{sz: 1, bigEndian: false, off: 1}, Value{t: Integer, i: 0xAB}, false},
		{"Extract non-integer type", fields{t: Floating, F: 123.456}, args{sz: 4, bigEndian: false, off: 2}, Value{t: Floating, f: 123.456}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			err := v.Extract(tt.args.sz, tt.args.bigEndian, tt.args.off)
			if (err != nil) != tt.wantErr {
				t.Errorf("Value.Extract() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Extract() = %v, want %v", *v, tt.want)
			}
		})
	}
}

func TestValue_Inc(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"Postincrement_I", fields{t: Integer, I: 0x12345}, Value{t: Integer, i: 0x12346}, false},
		{"Postincrement_F", fields{t: Floating, F: 123.45}, Value{t: Floating, f: 124.45}, false},
		{"Postincrement_fail", fields{t: Identifier}, Value{t: Identifier}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Inc(); (err != nil) != tt.wantErr {
				t.Errorf("Value.Inc() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Inc() %s = %v, want %v", tt.name, v, tt.want)
			}
		})
	}
}

func TestValue_Dec(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"Postdecrement_I", fields{t: Integer, I: 0x12345}, Value{t: Integer, i: 0x12344}, false},
		{"Postincrement_F", fields{t: Floating, F: 123.45}, Value{t: Floating, f: 122.45}, false},
		{"Postdecrement_fail", fields{t: Identifier}, Value{t: Identifier}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Dec(); (err != nil) != tt.wantErr {
				t.Errorf("Value.Dec() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Dec() %s = %v, want %v", tt.name, v, tt.want)
			}
		})
	}
}

func TestValue_Plus(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"+IntExpr", fields{t: Integer, I: 0x12345}, Value{t: Integer, i: 0x12345}, false},
		{"+FloatExpr", fields{t: Floating, F: 12.345}, Value{t: Floating, f: 12.345}, false},
		{"+err", fields{t: Add}, Value{t: Add}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Plus(); (err != nil) != tt.wantErr {
				t.Errorf("Value.Plus() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Plus() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Neg(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"-IntExpr", fields{t: Integer, I: 0x12345}, Value{t: Integer, i: -0x12345}, false},
		{"-FloatExpr", fields{t: Floating, F: 12.345}, Value{t: Floating, f: -12.345}, false},
		{"-err", fields{t: Sub}, Value{t: Sub}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Neg(); (err != nil) != tt.wantErr {
				t.Errorf("Value.Neg() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Neg() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Compl(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"~IntExpr", fields{t: Integer, I: 0x12345}, Value{t: Integer, i: 0x12345 ^ -1}, false},
		{"~err", fields{t: Compl}, Value{t: Compl}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Compl(); (err != nil) != tt.wantErr {
				t.Errorf("Value.Compl() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Compl() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Not(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    Value
		wantErr bool
	}{
		{"!IntExpr1", fields{t: Integer, I: 1}, Value{t: Integer, i: 0}, false},
		{"!IntExpr0", fields{t: Integer, I: 0}, Value{t: Integer, i: 1}, false},
		{"!err", fields{t: Not}, Value{t: Not}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Not(); (err != nil) != tt.wantErr {
				t.Errorf("Value.Not() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Not() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Cast(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		ty Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"(int8_t)0x12345", fields{t: Integer, I: 0x12345}, args{Int8}, Value{t: Integer, i: 0x45}, false},
		{"(int16_t)0x12345", fields{t: Integer, I: 0x12345}, args{Int16}, Value{t: Integer, i: 0x2345}, false},
		{"(int32_t)0x123456789", fields{t: Integer, I: 0x123456789}, args{Int32}, Value{t: Integer, i: 0x23456789}, false},
		{"(int64_t)0x12345678901234", fields{t: Integer, I: 0x12345678901234}, args{Int64}, Value{t: Integer, i: 0x12345678901234}, false},
		{"(uint8_t)-0x12345", fields{t: Integer, I: -0x12345}, args{Uint8}, Value{t: Integer, i: (-0x12345) & 0xFF}, false},
		{"(uint16_t)-0x12345", fields{t: Integer, I: -0x12345}, args{Uint16}, Value{t: Integer, i: (-0x12345) & 0xFFFF}, false},
		{"(uint32_t)-0x123456789", fields{t: Integer, I: -0x123456789}, args{Uint32}, Value{t: Integer, i: (-0x23456789) & 0xFFFFFFFF}, false},
		{"(uint64_t)-0x12345678901234", fields{t: Integer, I: -0x12345678901234}, args{Uint64}, Value{t: Integer, i: -0x12345678901234}, false},
		{"(int8_t)-483.12", fields{t: Floating, F: -483.12}, args{Int8}, Value{t: Integer, i: 0x1D}, false},
		{"(int16_t)-483.12", fields{t: Floating, F: -483.12}, args{Int16}, Value{t: Integer, i: -483}, false},
		{"(int32_t)-78483.12", fields{t: Floating, F: -78483.12}, args{Int32}, Value{t: Integer, i: -78483}, false},
		{"(int64_t)-9278483.12", fields{t: Floating, F: -9278483.12}, args{Int64}, Value{t: Integer, i: -9278483}, false},
		{"(uint8_t)483.12", fields{t: Floating, F: 483.12}, args{Uint8}, Value{t: Integer, i: 0xE3}, false},
		{"(uint16_t)483.12", fields{t: Floating, F: 483.12}, args{Uint16}, Value{t: Integer, i: 0x1E3}, false},
		{"(uint32_t)78483.12", fields{t: Floating, F: 78483.12}, args{Uint32}, Value{t: Integer, i: 78483}, false},
		{"(uint64_t)-9278483.12", fields{t: Floating, F: 9278483.12}, args{Uint64}, Value{t: Integer, i: 9278483}, false},
		{"(double)12345789", fields{t: Floating, F: 12345789}, args{Double}, Value{t: Floating, f: 12345789.0}, false},
		{"(float)123456789", fields{t: Floating, F: 123456789}, args{Float}, Value{t: Floating, f: 123456792.0}, false},
		{"(int8_t)err", fields{t: Nix}, args{Int8}, Value{t: Nix}, true},
		{"(int16_t)err", fields{t: Nix}, args{Int16}, Value{t: Nix}, true},
		{"(int32_t)err", fields{t: Nix}, args{Int32}, Value{t: Nix}, true},
		{"(int64_t)err", fields{t: Nix}, args{Int64}, Value{t: Nix}, true},
		{"(uint8_t)err", fields{t: Nix}, args{Uint8}, Value{t: Nix}, true},
		{"(uint16_t)err", fields{t: Nix}, args{Uint16}, Value{t: Nix}, true},
		{"(uint32_t)err", fields{t: Nix}, args{Uint32}, Value{t: Nix}, true},
		{"(uint64_t)err", fields{t: Nix}, args{Uint64}, Value{t: Nix}, true},
		{"(double)err", fields{t: Nix}, args{Double}, Value{t: Nix}, true},
		{"(float)err", fields{t: Nix}, args{Float}, Value{t: Nix}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Cast(tt.args.ty); (err != nil) != tt.wantErr {
				t.Errorf("Value.Cast() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Cast() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Mul(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"I*I", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 678}}, Value{t: Integer, i: 233910}, false},
		{"I*F", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 6.78}}, Value{t: Floating, f: 2339.1}, false},
		{"F*I", fields{t: Floating, F: 3.4}, args{&Value{t: Integer, i: 678}}, Value{t: Floating, f: 2305.2}, false},
		{"F*F", fields{t: Floating, F: 3.4}, args{&Value{t: Floating, f: 6.78}}, Value{t: Floating, f: 23.052}, false},
		{"I*X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
		{"F*X", fields{t: Floating, F: 3.4}, args{&Value{t: Nix}}, Value{t: Floating, f: 3.4}, true},
		{"X*F", fields{t: Nix}, args{&Value{t: Floating, f: 3.4}}, Value{t: Nix}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Mul(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.Mul() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Mul() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Div(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"I/I", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 15}}, Value{t: Integer, i: 23}, false},
		{"I/F", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 1.2}}, Value{t: Floating, f: 287.5}, false},
		{"F/I", fields{t: Floating, F: 3.45}, args{&Value{t: Integer, i: 15}}, Value{t: Floating, f: 0.23}, false},
		{"F/F", fields{t: Floating, F: 3.6}, args{&Value{t: Floating, f: 1.2}}, Value{t: Floating, f: 3}, false},
		{"I/0", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 0}}, Value{t: Integer, i: 345}, true},
		{"F/0", fields{t: Floating, F: 3.4}, args{&Value{t: Integer, i: 0}}, Value{t: Floating, f: 3.4}, true},
		{"I/0.0", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 0.0}}, Value{t: Integer, i: 345}, true},
		{"F/0.0", fields{t: Floating, F: 3.4}, args{&Value{t: Floating, f: 0.0}}, Value{t: Floating, f: 3.4}, true},
		{"I/X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
		{"F/X", fields{t: Floating, F: 3.4}, args{&Value{t: Nix}}, Value{t: Floating, f: 3.4}, true},
		{"X/F", fields{t: Nix}, args{&Value{t: Floating, f: 3.4}}, Value{t: Nix}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Div(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.Div() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Div() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Mod(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"I%I", fields{t: Integer, I: 347}, args{&Value{t: Integer, i: 15}}, Value{t: Integer, i: 2}, false},
		{"I%0", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 0}}, Value{t: Integer, i: 345}, true},
		{"I%F", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 1.2}}, Value{t: Integer, i: 345}, true},
		{"F%I", fields{t: Floating, F: 3.45}, args{&Value{t: Integer, i: 15}}, Value{t: Floating, f: 3.45}, true},
		{"I%X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
		{"X%I", fields{t: Nix}, args{&Value{t: Integer, i: 15}}, Value{t: Nix}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Mod(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.Mod() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Mod() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Add(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"I+I", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 678}}, Value{t: Integer, i: 1023}, false},
		{"I+F", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 6.78}}, Value{t: Floating, f: 351.78}, false},
		{"F+I", fields{t: Floating, F: 3.4}, args{&Value{t: Integer, i: 678}}, Value{t: Floating, f: 681.4}, false},
		{"F+F", fields{t: Floating, F: 3.4}, args{&Value{t: Floating, f: 6.78}}, Value{t: Floating, f: 10.18}, false},
		{"I+X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
		{"F+X", fields{t: Floating, F: 3.4}, args{&Value{t: Nix}}, Value{t: Floating, f: 3.4}, true},
		{"X+F", fields{t: Nix}, args{&Value{t: Floating, f: 3.4}}, Value{t: Nix}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Add(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.Add() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Add() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Sub(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"I-I", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 15}}, Value{t: Integer, i: 330}, false},
		{"I-F", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 1.2}}, Value{t: Floating, f: 343.8}, false},
		{"F-I", fields{t: Floating, F: 3.45}, args{&Value{t: Integer, i: 15}}, Value{t: Floating, f: -11.55}, false},
		{"F-F", fields{t: Floating, F: 3.4}, args{&Value{t: Floating, f: 1.2}}, Value{t: Floating, f: 2.2}, false},
		{"I-X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
		{"F-X", fields{t: Floating, F: 3.4}, args{&Value{t: Nix}}, Value{t: Floating, f: 3.4}, true},
		{"X-F", fields{t: Nix}, args{&Value{t: Floating, f: 3.4}}, Value{t: Nix}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Sub(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.Sub() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Sub() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Shl(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"345<<7", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 7}}, Value{t: Integer, i: 44160}, false},
		{"X<<7", fields{t: Nix}, args{&Value{t: Integer, i: 7}}, Value{t: Nix}, true},
		{"345<<X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Shl(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.Shl() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Shl() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Shr(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"345>>1", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 1}}, Value{t: Integer, i: 172}, false},
		{"X>>1", fields{t: Nix}, args{&Value{t: Integer, i: 1}}, Value{t: Nix}, true},
		{"345>>X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Shr(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.Shr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Shr() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Less(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"345<7", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 7}}, Value{t: Integer, i: 0}, false},
		{"345<789", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 789}}, Value{t: Integer, i: 1}, false},
		{"345<345", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345<7.1", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 7.1}}, Value{t: Integer, i: 0}, false},
		{"345<789.1", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 789.1}}, Value{t: Integer, i: 1}, false},
		{"345<345.0", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0<7", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 7}}, Value{t: Integer, i: 0}, false},
		{"345.0<789", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 789}}, Value{t: Integer, i: 1}, false},
		{"345.0<345", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345.0<7.1", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 7.1}}, Value{t: Integer, i: 0}, false},
		{"345.0<789.1", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 789.1}}, Value{t: Integer, i: 1}, false},
		{"345.0<345.0", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"I<X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
		{"F<X", fields{t: Floating, F: 3.4}, args{&Value{t: Nix}}, Value{t: Floating, f: 3.4}, true},
		{"X<F", fields{t: Nix}, args{&Value{t: Floating, f: 3.4}}, Value{t: Nix}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Less(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.Less() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Less() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_LessEqual(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"345<=7", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 7}}, Value{t: Integer, i: 0}, false},
		{"345<=789", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 789}}, Value{t: Integer, i: 1}, false},
		{"345<=345", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345<=7.1", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 7.1}}, Value{t: Integer, i: 0}, false},
		{"345<=789.1", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 789.1}}, Value{t: Integer, i: 1}, false},
		{"345<=345.0", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0<=7", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 7}}, Value{t: Integer, i: 0}, false},
		{"345.0<=789", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 789}}, Value{t: Integer, i: 1}, false},
		{"345.0<=345", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345.0<=7.1", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 7.1}}, Value{t: Integer, i: 0}, false},
		{"345.0<=789.1", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 789.1}}, Value{t: Integer, i: 1}, false},
		{"345.0<=345.0", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"I<=X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
		{"F<=X", fields{t: Floating, F: 3.4}, args{&Value{t: Nix}}, Value{t: Floating, f: 3.4}, true},
		{"X<=F", fields{t: Nix}, args{&Value{t: Floating, f: 3.4}}, Value{t: Nix}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.LessEqual(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.LessEqual() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.LessEqual() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Greater(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"345>7", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 7}}, Value{t: Integer, i: 1}, false},
		{"345>789", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 789}}, Value{t: Integer, i: 0}, false},
		{"345>345", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345>7.1", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 7.1}}, Value{t: Integer, i: 1}, false},
		{"345>789.1", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 789.1}}, Value{t: Integer, i: 0}, false},
		{"345>345.0", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0>7", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 7}}, Value{t: Integer, i: 1}, false},
		{"345.0>789", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 789}}, Value{t: Integer, i: 0}, false},
		{"345.0>345", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345.0>7.1", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 7.1}}, Value{t: Integer, i: 1}, false},
		{"345.0>789.1", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 789.1}}, Value{t: Integer, i: 0}, false},
		{"345.0>345.0", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"I>X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
		{"F>X", fields{t: Floating, F: 3.4}, args{&Value{t: Nix}}, Value{t: Floating, f: 3.4}, true},
		{"X>F", fields{t: Nix}, args{&Value{t: Floating, f: 3.4}}, Value{t: Nix}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Greater(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.Greater() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Greater() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_GreaterEqual(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"345>=7", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 7}}, Value{t: Integer, i: 1}, false},
		{"345>=789", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 789}}, Value{t: Integer, i: 0}, false},
		{"345>=345", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345>=7.1", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 7.1}}, Value{t: Integer, i: 1}, false},
		{"345>=789.1", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 789.1}}, Value{t: Integer, i: 0}, false},
		{"345>=345.0", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0>=7", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 7}}, Value{t: Integer, i: 1}, false},
		{"345.0>=789", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 789}}, Value{t: Integer, i: 0}, false},
		{"345.0>=345", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345.0>=7.1", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 7.1}}, Value{t: Integer, i: 1}, false},
		{"345.0>=789.1", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 789.1}}, Value{t: Integer, i: 0}, false},
		{"345.0>=345.0", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"I>=X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
		{"F>=X", fields{t: Floating, F: 3.4}, args{&Value{t: Nix}}, Value{t: Floating, f: 3.4}, true},
		{"X>=F", fields{t: Nix}, args{&Value{t: Floating, f: 3.4}}, Value{t: Nix}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.GreaterEqual(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.GreaterEqual() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.GreaterEqual() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Equal(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"345==7", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 7}}, Value{t: Integer, i: 0}, false},
		{"345==789", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 789}}, Value{t: Integer, i: 0}, false},
		{"345==345", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345==7.1", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 7.1}}, Value{t: Integer, i: 0}, false},
		{"345==789.1", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 789.1}}, Value{t: Integer, i: 0}, false},
		{"345==345.0", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"345.0==7", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 7}}, Value{t: Integer, i: 0}, false},
		{"345.0==789", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 789}}, Value{t: Integer, i: 0}, false},
		{"345.0==345", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 345}}, Value{t: Integer, i: 1}, false},
		{"345.0==7.1", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 7.1}}, Value{t: Integer, i: 0}, false},
		{"345.0==789.1", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 789.1}}, Value{t: Integer, i: 0}, false},
		{"345.0==345.0", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 1}, false},
		{"I==X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
		{"F==X", fields{t: Floating, F: 3.4}, args{&Value{t: Nix}}, Value{t: Floating, f: 3.4}, true},
		{"X==F", fields{t: Nix}, args{&Value{t: Floating, f: 3.4}}, Value{t: Nix}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Equal(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.Equal() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Equal() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_NotEqual(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"345!=7", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 7}}, Value{t: Integer, i: 1}, false},
		{"345!=789", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 789}}, Value{t: Integer, i: 1}, false},
		{"345!=345", fields{t: Integer, I: 345}, args{&Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345!=7.1", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 7.1}}, Value{t: Integer, i: 1}, false},
		{"345!=789.1", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 789.1}}, Value{t: Integer, i: 1}, false},
		{"345!=345.0", fields{t: Integer, I: 345}, args{&Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"345.0!=7", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 7}}, Value{t: Integer, i: 1}, false},
		{"345.0!=789", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 789}}, Value{t: Integer, i: 1}, false},
		{"345.0!=345", fields{t: Floating, F: 345.0}, args{&Value{t: Integer, i: 345}}, Value{t: Integer, i: 0}, false},
		{"345.0!=7.1", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 7.1}}, Value{t: Integer, i: 1}, false},
		{"345.0!=789.1", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 789.1}}, Value{t: Integer, i: 1}, false},
		{"345.0!=345.0", fields{t: Floating, F: 345.0}, args{&Value{t: Floating, f: 345.0}}, Value{t: Integer, i: 0}, false},
		{"I!=X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
		{"F!=X", fields{t: Floating, F: 3.4}, args{&Value{t: Nix}}, Value{t: Floating, f: 3.4}, true},
		{"X!=F", fields{t: Nix}, args{&Value{t: Floating, f: 3.4}}, Value{t: Nix}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.NotEqual(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.NotEqual() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.NotEqual() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_And(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"0x55aa00ff&0xaf5f0ff0", fields{t: Integer, I: 0x55aa00ff}, args{&Value{t: Integer, i: 0xaf5f0ff0}}, Value{t: Integer, i: 0x050A00F0}, false},
		{"X&7", fields{t: Nix}, args{&Value{t: Integer, i: 7}}, Value{t: Nix}, true},
		{"345&X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.And(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.And() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.And() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Xor(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"0x55aa00ff^0xaf5f0ff0", fields{t: Integer, I: 0x55aa00ff}, args{&Value{t: Integer, i: 0xaf5f0ff0}}, Value{t: Integer, i: 0xFAF50F0F}, false},
		{"X^7", fields{t: Nix}, args{&Value{t: Integer, i: 7}}, Value{t: Nix}, true},
		{"345^X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Xor(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.Xor() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Xor() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_Or(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"0x55aa00ff|&0xaf5f0ff0", fields{t: Integer, I: 0x55aa00ff}, args{&Value{t: Integer, i: 0xaf5f0ff0}}, Value{t: Integer, i: 0xFFFF0FFF}, false},
		{"X|7", fields{t: Nix}, args{&Value{t: Integer, i: 7}}, Value{t: Nix}, true},
		{"345|X", fields{t: Integer, I: 345}, args{&Value{t: Nix}}, Value{t: Integer, i: 345}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.Or(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.Or() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.Or() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_LogAnd(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"0&&0", fields{t: Integer, I: 0}, args{&Value{t: Integer, i: 0}}, Value{t: Integer, i: 0}, false},
		{"0&&1", fields{t: Integer, I: 0}, args{&Value{t: Integer, i: 1}}, Value{t: Integer, i: 0}, false},
		{"1&&0", fields{t: Integer, I: 1}, args{&Value{t: Integer, i: 0}}, Value{t: Integer, i: 0}, false},
		{"1&&1", fields{t: Integer, I: 1}, args{&Value{t: Integer, i: 1}}, Value{t: Integer, i: 1}, false},
		{"0&&0.0", fields{t: Integer, I: 0}, args{&Value{t: Floating, f: 0.0}}, Value{t: Integer, i: 0}, false},
		{"0&&1.0", fields{t: Integer, I: 0}, args{&Value{t: Floating, f: 1.0}}, Value{t: Integer, i: 0}, false},
		{"1&&0.0", fields{t: Integer, I: 1}, args{&Value{t: Floating, f: 0.0}}, Value{t: Integer, i: 0}, false},
		{"1&&1.0", fields{t: Integer, I: 1}, args{&Value{t: Floating, f: 1.0}}, Value{t: Integer, i: 1}, false},
		{"0.0&&0.0", fields{t: Floating, F: 0.0}, args{&Value{t: Integer, i: 0}}, Value{t: Integer, i: 0}, false},
		{"0.0&&1.0", fields{t: Floating, F: 1.0}, args{&Value{t: Integer, i: 0}}, Value{t: Integer, i: 0}, false},
		{"1.0&&0.0", fields{t: Floating, F: 0.0}, args{&Value{t: Integer, i: 1}}, Value{t: Integer, i: 0}, false},
		{"1.0&&1.0", fields{t: Floating, F: 1.0}, args{&Value{t: Integer, i: 1}}, Value{t: Integer, i: 1}, false},
		{"X&&1", fields{t: Nix}, args{&Value{t: Integer, i: 1}}, Value{t: Nix}, true},
		{"1&&X", fields{t: Integer, I: 1}, args{&Value{t: Nix}}, Value{t: Integer, i: 1}, true},
		{"X&&1.0", fields{t: Nix}, args{&Value{t: Floating, f: 1.0}}, Value{t: Nix}, true},
		{"1.0&&X", fields{t: Floating, F: 1.0}, args{&Value{t: Nix}}, Value{t: Floating, f: 1.0}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.LogAnd(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.LogAnd() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.LogAnd() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}

func TestValue_LogOr(t *testing.T) {
	t.Parallel()

	type fields struct {
		t Token
		I int64
		F float64
		s string
		v *Variable
		l []Value
	}
	type args struct {
		v1 *Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{"0||0", fields{t: Integer, I: 0}, args{&Value{t: Integer, i: 0}}, Value{t: Integer, i: 0}, false},
		{"0||1", fields{t: Integer, I: 0}, args{&Value{t: Integer, i: 1}}, Value{t: Integer, i: 1}, false},
		{"1||0", fields{t: Integer, I: 1}, args{&Value{t: Integer, i: 0}}, Value{t: Integer, i: 1}, false},
		{"1||1", fields{t: Integer, I: 1}, args{&Value{t: Integer, i: 1}}, Value{t: Integer, i: 1}, false},
		{"0||0.0", fields{t: Integer, I: 0}, args{&Value{t: Floating, f: 0.0}}, Value{t: Integer, i: 0}, false},
		{"0||1.0", fields{t: Integer, I: 0}, args{&Value{t: Floating, f: 1.0}}, Value{t: Integer, i: 1}, false},
		{"1||0.0", fields{t: Integer, I: 1}, args{&Value{t: Floating, f: 0.0}}, Value{t: Integer, i: 1}, false},
		{"1||1.0", fields{t: Integer, I: 1}, args{&Value{t: Floating, f: 1.0}}, Value{t: Integer, i: 1}, false},
		{"0.0||0.0", fields{t: Floating, F: 0.0}, args{&Value{t: Integer, i: 0}}, Value{t: Integer, i: 0}, false},
		{"0.0||1.0", fields{t: Floating, F: 1.0}, args{&Value{t: Integer, i: 0}}, Value{t: Integer, i: 1}, false},
		{"1.0||0.0", fields{t: Floating, F: 0.0}, args{&Value{t: Integer, i: 1}}, Value{t: Integer, i: 1}, false},
		{"1.0||1.0", fields{t: Floating, F: 1.0}, args{&Value{t: Integer, i: 1}}, Value{t: Integer, i: 1}, false},
		{"X||1", fields{t: Nix}, args{&Value{t: Integer, i: 1}}, Value{t: Nix}, true},
		{"1||X", fields{t: Integer, I: 1}, args{&Value{t: Nix}}, Value{t: Integer, i: 1}, true},
		{"X||1.0", fields{t: Nix}, args{&Value{t: Floating, f: 1.0}}, Value{t: Nix}, true},
		{"1.0||X", fields{t: Floating, F: 1.0}, args{&Value{t: Nix}}, Value{t: Floating, f: 1.0}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Value{
				t: tt.fields.t,
				i: tt.fields.I,
				f: tt.fields.F,
				s: tt.fields.s,
				v: tt.fields.v,
				l: tt.fields.l,
			}
			if err := v.LogOr(tt.args.v1); (err != nil) != tt.wantErr {
				t.Errorf("Value.LogOr() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(*v, tt.want) {
				t.Errorf("Value.LorOr() %s = %v, want %v", tt.name, *v, tt.want)
			}
		})
	}
}
