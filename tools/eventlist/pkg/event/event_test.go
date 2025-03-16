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

package event

import (
	"bufio"
	"errors"
	"eventlist/pkg/elf"
	"eventlist/pkg/eval"
	"eventlist/pkg/xml/scvd"
	"reflect"
	"testing"
)

func Test_getEnum(t *testing.T) { //nolint:golint,paralleltest
	var vals eval.Member
	vals.Enums = make(map[int64]string)
	var td eval.ITypedef
	td.Members = make(map[string]eval.Member)
	var tds = make(eval.Typedefs)

	vals.Enums[4711] = "enum"
	td.Members["enumName"] = vals
	tds["typName"] = td

	var i int

	type args struct {
		typedefs eval.Typedefs
		val      int64
		value    string
		i        *int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantI   int
		wantErr bool
	}{
		{"enum E", args{tds, 4711, "typName:enumName]", &i}, "enum", 17, false},
		{"enum err1", args{tds, 4711, "typName", &i}, "", 0, true},
		{"enum err2", args{tds, 4711, "typ:", &i}, "", 0, true},
		{"enum err3", args{tds, 4711, "typName:", &i}, "", 8, true},
		{"enum err4", args{tds, 4711, "typName:eee", &i}, "", 8, true},
		{"enum err5", args{tds, 4711, "typName:eee]", &i}, "", 8, true},
		{"enum err6", args{tds, 47, "typName:enumName]", &i}, "", 8, true},
		{"enum err7", args{tds, 4711, "typName]", &i}, "enum", 8, false},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			i = 0
			got, err := getEnum(tt.args.typedefs, tt.args.val, tt.args.value, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getEnum() got = %v, want %v", got, tt.want)
			}
			if i != tt.wantI {
				t.Errorf("getEnum() idx = %v, want %v", i, tt.wantI)
			}
		})
	}
}

func TestInfo_getInfoFromBytes(t *testing.T) {
	t.Parallel()

	type fields struct {
		ID     scvd.IDType
		length uint16
		irq    bool
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Info
	}{
		{"normal", fields{}, args{[]byte{0x34, 0x12, 0x78, 0xd6}}, Info{0x1234, 0x5678, true}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info := &Info{
				ID:     tt.fields.ID,
				length: tt.fields.length,
				irq:    tt.fields.irq,
			}
			info.getInfoFromBytes(tt.args.data)
			//			if info.ID != tt.want.ID || info.length != tt.want.length || info.irq != tt.want.irq {
			if !reflect.DeepEqual(info, &tt.want) { // this does not work
				t.Errorf("getInfoFramBytes() = %v, want %v", info, &tt.want)
			}
		})
	}
}

func TestInfo_SplitID(t *testing.T) {
	t.Parallel()

	type fields struct {
		ID     scvd.IDType
		length uint16
		irq    bool
	}
	tests := []struct {
		name      string
		fields    fields
		wantClass uint16
		wantGroup uint16
		wantIdx   uint16
		wantStart bool
	}{
		{"0x0000", fields{ID: 0x0000}, 0x00, 0, 0, true},
		{"0xEF00", fields{ID: 0xEF00}, 0xEF, 0, 0, true},
		{"0xEF0A", fields{ID: 0xEF0A}, 0xEF, 0, 10, true},
		{"0xEF35", fields{ID: 0xEF35}, 0xEF, 0, 5, false},
		{"0xEF91", fields{ID: 0xEF91}, 0xEF, 2, 1, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info := &Info{
				ID:     tt.fields.ID,
				length: tt.fields.length,
				irq:    tt.fields.irq,
			}
			gotClass, gotGroup, gotIdx, gotStart := info.SplitID()
			if gotClass != tt.wantClass {
				t.Errorf("Info.SplitID() gotClass = %v, want %v", gotClass, tt.wantClass)
			}
			if gotGroup != tt.wantGroup {
				t.Errorf("Info.SplitID() gotGroup = %v, want %v", gotGroup, tt.wantGroup)
			}
			if gotIdx != tt.wantIdx {
				t.Errorf("Info.SplitID() gotIdx = %v, want %v", gotIdx, tt.wantIdx)
			}
			if gotStart != tt.wantStart {
				t.Errorf("Info.SplitID() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
		})
	}
}

func TestEventData_calculateExpression(t *testing.T) { //nolint:golint,paralleltest
	var i int

	fileTest := "../../testdata/elftest.elf"

	type fields struct {
		Time   uint64
		Value1 int32
		Value2 int32
		Value3 int32
		Value4 int32
		Data   *[]uint8
		Info   Info
	}

	var ed1 = fields{Time: 306, Value1: 257, Value2: -24, Value3: 625478261, Value4: 0x4010, Data: nil, Info: Info{}}
	var tds = make(eval.Typedefs)

	type args struct {
		typedefs eval.Typedefs
		value    string
		i        *int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantI   int
		wantErr bool
	}{
		{"expr empty", ed1, args{tds, "", &i}, "", 0, true},
		{"expr T", ed1, args{tds, "T[val2]", &i}, "-24", 7, false},
		{"expr d", ed1, args{tds, "d[val2]", &i}, "-24", 7, false},
		{"expr u", ed1, args{tds, "u[val1]", &i}, "257", 7, false},
		{"expr t", ed1, args{tds, "t[val4]", &i}, "def", 7, false},
		{"expr x", ed1, args{tds, "x[val1]", &i}, "0x101", 7, false},
		{"expr F", ed1, args{tds, "F[val4]", &i}, "def", 7, false},
		{"expr F", ed1, args{tds, "F[val1]", &i}, "0x00000101", 7, false},
		{"expr C", ed1, args{tds, "C[val2]", &i}, "", 7, true},
		{"expr I", ed1, args{tds, "I[val3]", &i}, "37.72.10.117", 7, false},
		{"expr J", ed1, args{tds, "J[val3]", &i}, "0:0:2548:a75:", 7, false},
		{"expr N", ed1, args{tds, "N[val4]", &i}, "def", 7, false},
		{"expr N", ed1, args{tds, "N[val1]", &i}, "0x00000101", 7, false},
		{"expr M", ed1, args{tds, "M[val3]", &i}, "00-00-25-48-0a-75", 7, false},
		{"expr S", ed1, args{tds, "S[val3]", &i}, "25480a75", 7, false},
		{"expr T", ed1, args{tds, "T[val1+0.234]", &i}, "257.234000", 13, false},
		{"expr ?", ed1, args{tds, "?[val3]", &i}, "?", 7, false},
		{"expr err1", ed1, args{tds, "S[", &i}, "", 2, true},
		{"expr err2", ed1, args{tds, "S[val3,", &i}, "", 6, true},
	}
	if err := elf.Sections.Readelf(&fileTest); err != nil {
		t.Errorf("Data.calculateExpression() cannot open %s", fileTest)
		return
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			e := &Data{
				Time:   tt.fields.Time,
				Value1: tt.fields.Value1,
				Value2: tt.fields.Value2,
				Value3: tt.fields.Value3,
				Value4: tt.fields.Value4,
				Data:   tt.fields.Data,
				Info:   tt.fields.Info,
			}
			i = 0
			got, err := e.calculateExpression(tt.args.typedefs, nil, tt.args.value, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Data.calculateExpression() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Data.calculateExpression() %s = %v, want %v", tt.name, got, tt.want)
			}
			if i != tt.wantI {
				t.Errorf("Data.calculateExpression() %s idx = %v, want %v", tt.name, i, tt.wantI)
			}
		})
	}
}

func TestEventData_calculateEnumExpression(t *testing.T) { //nolint:golint,paralleltest
	var vals eval.Member
	vals.Enums = make(map[int64]string)
	var td eval.ITypedef
	td.Members = make(map[string]eval.Member)
	var tds = make(eval.Typedefs)
	var tdm eval.ITypedef
	tdm.Members = make(map[string]eval.Member)
	var tdsm = make(eval.Typedefs)

	vals.Enums[4711] = "enum"
	td.Members["enumName"] = vals
	tds["typName"] = td
	tdm.Members["sub"] = vals
	tdsm["typ"] = tdm

	var i int

	type fields struct {
		Time   uint64
		Value1 int32
		Value2 int32
		Value3 int32
		Value4 int32
		Data   *[]uint8
		Info   Info
	}

	var ed1 = fields{Time: 306, Value1: 257, Value2: 4711, Value3: 625478261, Value4: 0, Data: nil, Info: Info{}}

	type args struct {
		typedefs eval.Typedefs
		value    string
		i        *int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantI   int
		wantErr bool
	}{
		{"enumExpr empty", ed1, args{tds, "", &i}, "", 0, true},
		{"enumExpr E", ed1, args{tds, "E[val2, typName]", &i}, "enum", 16, false},
		{"enumExpr Esub", ed1, args{tdsm, "E[val2, typ:sub]", &i}, "enum", 16, false},
		{"enumExpr err1", ed1, args{tds, "S[", &i}, "", 2, true},
		{"enumExpr err2", ed1, args{tds, "S[val3]", &i}, "", 6, true},
		{"enumExpr err3", ed1, args{tds, "E[val3, xxx]", &i}, "", 12, true},
		{"enumExpr err4", ed1, args{tds, "S[val3, xxx]", &i}, "", 7, true},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			e := &Data{
				Time:   tt.fields.Time,
				Value1: tt.fields.Value1,
				Value2: tt.fields.Value2,
				Value3: tt.fields.Value3,
				Value4: tt.fields.Value4,
				Data:   tt.fields.Data,
				Info:   tt.fields.Info,
			}
			i = 0
			got, err := e.calculateEnumExpression(tt.args.typedefs, tt.args.value, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Data.calculateEnumExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Data.calculateEnumExpression() = %v, want %v", got, tt.want)
			}
			if i != tt.wantI {
				t.Errorf("Data.calculateEnumExpression() idx = %v, want %v", i, tt.wantI)
			}
		})
	}
}

func TestEventData_EvalLine(t *testing.T) {
	t.Parallel()

	var ev1 scvd.EventType = scvd.EventType{ID: "id1", Value: "x%%%d[val1]y%u[val2]z"}
	var ev2 scvd.EventType = scvd.EventType{ID: "id2", Value: "x%T[val1]y%x[val2]z"}
	var ev3 scvd.EventType = scvd.EventType{ID: "id3", Value: "x%I[val3]y%J[val3]z"}
	var ev4 scvd.EventType = scvd.EventType{ID: "id4", Value: "x%M[val3]y%S[val3]z"}
	var evE1 scvd.EventType = scvd.EventType{ID: "idE1", Value: "x%E[val2, typName]y"}
	var evTD scvd.EventType = scvd.EventType{ID: "idTD", Val1: "v1", Val2: "v2", Val3: "4BY", Val4: "v4", Val5: "v5", Val6: "v6", Value: "x%x[val3.B2]y"}
	var everr1 scvd.EventType = scvd.EventType{ID: "iderr1", Value: "x%d[;]y"}
	var everr2 scvd.EventType = scvd.EventType{ID: "iderr2", Value: "x%E[;]y"}

	var vals eval.Member
	vals.Enums = make(map[int64]string)
	var td eval.ITypedef
	td.Members = make(map[string]eval.Member)

	var td1 eval.ITypedef
	td1.Size = 4
	td1.Members = make(map[string]eval.Member)
	td1.Members["B2"] = eval.Member{Offset: "2", IType: eval.Uint8}
	var tds = make(eval.Typedefs)

	vals.Enums[4711] = "enum"
	td.Members["enumName"] = vals
	tds["typName"] = td
	tds["4BY"] = td1

	type fields struct {
		Time   uint64
		Value1 int32
		Value2 int32
		Value3 int32
		Value4 int32
		Data   *[]uint8
		Info   Info
	}

	var ed1 = fields{Time: 306, Value1: 257, Value2: 4711, Value3: 625478261, Value4: 0, Data: nil, Info: Info{}}

	type args struct {
		scvdevent scvd.EventType
		typedefs  eval.Typedefs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"EvalLine ev1", ed1, args{ev1, tds}, "x%257y4711z", false},
		{"EvalLine ev2", ed1, args{ev2, tds}, "x257y0x1267z", false},
		{"EvalLine ev3", ed1, args{ev3, tds}, "x37.72.10.117y0:0:2548:a75:z", false},
		{"EvalLine ev4", ed1, args{ev4, tds}, "x00-00-25-48-0a-75y25480a75z", false},
		{"EvalLine evE1", ed1, args{evE1, tds}, "xenumy", false},
		{"EvalLine evTD", ed1, args{evTD, tds}, "x0x48y", false},
		{"EvalLine err1", ed1, args{everr1, tds}, "", true},
		{"EvalLine err2", ed1, args{everr2, tds}, "", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &Data{
				Time:   tt.fields.Time,
				Value1: tt.fields.Value1,
				Value2: tt.fields.Value2,
				Value3: tt.fields.Value3,
				Value4: tt.fields.Value4,
				Data:   tt.fields.Data,
				Info:   tt.fields.Info,
			}
			got, err := e.EvalLine(tt.args.scvdevent, tt.args.typedefs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Data.EvalLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Data.EvalLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData_GetValuesAsString(t *testing.T) {
	t.Parallel()

	type fields struct {
		Time   uint64
		Value1 int32
		Value2 int32
		Value3 int32
		Value4 int32
		Data   *[]uint8
		Typ    uint16
		Info   Info
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"GetValuesAsString 0", fields{Typ: 0, Data: &[]uint8{1, 2, 0x80}, Value1: -0x8000, Value2: 0x12345678, Value3: 0xABCDEF, Value4: 0x76543210}, ""},
		{"GetValuesAsString 1", fields{Typ: 1, Data: &[]uint8{1, 2, 0x80}}, "data=0x010280"},
		{"GetValuesAsString 2", fields{Typ: 2, Value1: -0x8000, Value2: 0x12345678}, "val1=0xffff8000, val2=0x12345678"},
		{"GetValuesAsString 3", fields{Typ: 3, Value1: -0x8000, Value2: 0x12345678, Value3: 0xABCDEF, Value4: 0x76543210}, "val1=0xffff8000, val2=0x12345678, val3=0x00abcdef, val4=0x76543210"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &Data{
				Time:   tt.fields.Time,
				Value1: tt.fields.Value1,
				Value2: tt.fields.Value2,
				Value3: tt.fields.Value3,
				Value4: tt.fields.Value4,
				Data:   tt.fields.Data,
				Typ:    tt.fields.Typ,
				Info:   tt.fields.Info,
			}
			if got := e.GetValuesAsString(); got != tt.want {
				t.Errorf("Data.GetValuesAsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convert16(t *testing.T) {
	t.Parallel()

	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want uint16
	}{
		{"normal", args{[]byte{0x55, 0xAA}}, 0xAA55},
		{"fail", args{[]byte{0x55, 0xAA, 0x55}}, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := convert16(tt.args.data); got != tt.want {
				t.Errorf("convert16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convert32(t *testing.T) {
	t.Parallel()

	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{"normal", args{[]byte{0x55, 0xAA, 0x00, 0xFF}}, 0xFF00AA55},
		{"fail", args{[]byte{0x55, 0xAA, 0x55}}, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := convert32(tt.args.data); got != tt.want {
				t.Errorf("convert32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convert64(t *testing.T) {
	t.Parallel()

	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"normal", args{[]byte{0x55, 0xAA, 0x00, 0xFF, 0x01, 0x20, 0x13, 0x41}}, 0x41132001FF00AA55},
		{"fail", args{[]byte{0x55, 0xAA, 0x55}}, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := convert64(tt.args.data); got != tt.want {
				t.Errorf("convert64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventData_Read(t *testing.T) {
	t.Parallel()

	var s0 = "../../testdata/test0.binary"
	var s1 = "../../testdata/test1.binary"
	var s2 = "../../testdata/test2.binary"
	var sNix = "../../testdata/xxxx"
	var s3 = "../../testdata/test3.binary"
	var s4 = "../../testdata/test4.binary"
	var s5 = "../../testdata/test5.binary"
	var s8 = "../../testdata/test8.binary"
	var s9 = "../../testdata/test9.binary"
	var s12 = "../../testdata/test12.binary"
	var s13 = "../../testdata/test13.binary"

	var b0 = []uint8("hello wo")

	type fields struct {
		Time   uint64
		Value1 int32
		Value2 int32
		Value3 int32
		Value4 int32
		Data   *[]uint8
		Typ    uint16
		Info   Info
	}
	type args struct {
		in *bufio.Reader
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		file       *string
		wantNoOpen bool
		want       Data
		wantEOF    bool
		wantErr    bool
	}{
		{"read fail0", fields{}, args{}, &s0, false, Data{}, true, false},
		{"read fail1", fields{}, args{}, &s1, false, Data{}, false, true},
		{"read fail2", fields{}, args{}, &s2, false, Data{}, false, true},
		{"read fail3", fields{}, args{}, &s8, false, Data{}, true, false},
		{"read fail4", fields{}, args{}, &s9, false, Data{Typ: 1, Time: 1410, Info: Info{0xfe00, 8, false}}, true, false},
		{"read fail5", fields{}, args{}, &s12, false, Data{Typ: 2, Time: 31, Info: Info{0xff00, 0, false}}, true, false},
		{"read fail6", fields{}, args{}, &s13, false, Data{Typ: 3, Time: 306, Info: Info{0xf000, 0, true}}, true, false},
		{"read failOpen", fields{}, args{}, &sNix, true, Data{}, true, false},
		{"read ok1", fields{}, args{}, &s3, false, Data{Typ: 1, Data: &b0, Time: 1410, Info: Info{0xfe00, 8, false}}, false, false},
		{"read ok2", fields{}, args{}, &s4, false, Data{Typ: 2, Value1: 1, Value2: 2, Time: 31, Info: Info{0xff00, 0, false}}, false, false},
		{"read ok3", fields{}, args{}, &s5, false, Data{Typ: 3, Value1: 805332648, Value2: 24000, Value3: 1, Value4: -65536, Time: 306, Info: Info{0xf000, 0, true}}, false, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &Data{
				Time:   tt.fields.Time,
				Value1: tt.fields.Value1,
				Value2: tt.fields.Value2,
				Value3: tt.fields.Value3,
				Value4: tt.fields.Value4,
				Data:   tt.fields.Data,
				Typ:    tt.fields.Typ,
				Info:   tt.fields.Info,
			}
			var b Binary
			tt.args.in = b.Open(tt.file)
			if (tt.args.in == nil) != tt.wantNoOpen {
				t.Errorf("Data.Read() cannot open file %v", tt.file)
			}
			err := e.Read(tt.args.in)
			if errors.Is(err, eval.ErrEof) != tt.wantEOF {
				t.Errorf("Data.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !errors.Is(err, eval.ErrEof) && (err != nil) != tt.wantErr {
				t.Errorf("Data.Read() error = %v, wantErr %v", err, tt.wantErr)
			}
			b.Close()
			if !reflect.DeepEqual(e, &tt.want) {
				t.Errorf("Data.Read() = %v, want %v", e, &tt.want)
			}
		})
	}
}

func TestData_GetValue(t *testing.T) { //nolint:golint,paralleltest
	type fields struct {
		Time   uint64
		Value1 int32
		Value2 int32
		Value3 int32
		Value4 int32
		Data   *[]uint8
		Typ    uint16
		Info   Info
	}

	var ed1 = fields{Time: 306, Value1: 0x300066a8, Value2: 24000, Value3: 1, Value4: 0, Data: nil, Info: Info{}}
	var hello = []uint8("Hello")
	var ed2 = fields{Time: 306, Value1: 0, Value2: 0, Value3: 0, Value4: 0, Data: &hello, Info: Info{}}

	var vals eval.Member
	vals.Enums = make(map[int64]string)
	var td eval.ITypedef
	td.Members = make(map[string]eval.Member)
	var tds = make(eval.Typedefs)

	vals.Enums[4711] = "enum"
	td.Members["b0"] = vals
	tds["by4"] = td

	var vals1 eval.Member
	var td1 eval.ITypedef
	td1.Members = make(map[string]eval.Member)
	var tds1 = make(eval.Typedefs)

	vals1.IType = eval.Uint8
	vals1.Offset = "2"
	td1.Members["b2"] = vals1
	td1.Size = 4
	tds1["by4"] = td1

	tdu := make(map[string]string)
	tdu["val1"] = "by4"

	var i int

	type args struct {
		value    string
		i        *int
		typedefs eval.Typedefs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		gen     int
		want    eval.Value
		wantErr bool
	}{
		{"val1.b2", ed1, args{"[val1.b2]", &i, tds1}, 3, eval.Value{}, false},
		{"val1", ed1, args{"[val1]", &i, tds}, 1, eval.Value{}, false},
		{"data", ed2, args{"[val1]", &i, tds}, 2, eval.Value{}, false},
		{"nixvar", ed2, args{"xx", &i, tds}, 42, eval.Value{}, true},
		{"valxxx", ed1, args{"[valxxx]", &i, tds}, 42, eval.Value{}, true},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			i = 0
			e := &Data{
				Time:   tt.fields.Time,
				Value1: tt.fields.Value1,
				Value2: tt.fields.Value2,
				Value3: tt.fields.Value3,
				Value4: tt.fields.Value4,
				Data:   tt.fields.Data,
				Typ:    tt.fields.Typ,
				Info:   tt.fields.Info,
			}
			switch tt.gen {
			case 1:
				tt.want.Compose(eval.Integer, 0x300066a8, 0.0, "")
			case 2:
				tt.want.Compose(eval.Integer, 0x48656C6C, 0.0, "")
			case 3:
				tt.want.Compose(eval.Integer, 0x00, 0.0, "")
			}
			got, err := e.GetValue(tt.args.value, tt.args.i, tt.args.typedefs, tdu)
			if (err != nil) != tt.wantErr {
				t.Errorf("Data.GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data.GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
