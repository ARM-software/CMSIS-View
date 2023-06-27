/*
 * Copyright (c) 2022-2023 Arm Limited. All rights reserved.
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
	var vals = make(map[int16]string)
	var enms = make(map[string]eval.TdMember)
	var tds = make(map[string]map[string]eval.TdMember)

	vals[4711] = "enum"
	var e = enms["enumName"]
	e.Enum = vals
	enms["enumName"] = e
	tds["typName"] = enms

	var i int

	type args struct {
		typedefs map[string]map[string]eval.TdMember
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
		ID     uint16
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
		ID     uint16
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
	var member = make(map[int16]string)
	var members = make(map[string]eval.TdMember)
	var tds = make(map[string]map[string]eval.TdMember)

	var m = members["B2"]
	m.Enum = member
	members["B2"] = m
	tds["4BY"] = members
	event := scvd.Event{Val1: "4BY"}
	var sc scvd.ScvdData
	sc.Events = make(map[uint16]scvd.Event)
	sc.Events[1] = event
	eval.Typedefs = tds

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

	type args struct {
		sc    *scvd.ScvdData
		event scvd.Event
		value string
		i     *int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantI   int
		wantErr bool
	}{
		{"expr x member", ed1, args{&sc, event, "x[val1.B2]", &i}, "-24", 7, false},
		{"expr empty", ed1, args{&sc, event, "", &i}, "", 0, true},
		{"expr T", ed1, args{&sc, event, "T[val2]", &i}, "-24", 7, false},
		{"expr d", ed1, args{&sc, event, "d[val2]", &i}, "-24", 7, false},
		{"expr u", ed1, args{&sc, event, "u[val1]", &i}, "257", 7, false},
		{"expr t", ed1, args{&sc, event, "t[val4]", &i}, "def", 7, false},
		{"expr x", ed1, args{&sc, event, "x[val1]", &i}, "0x101", 7, false},
		{"expr F", ed1, args{&sc, event, "F[val4]", &i}, "def", 7, false},
		{"expr F", ed1, args{&sc, event, "F[val1]", &i}, "0x00000101", 7, false},
		{"expr C", ed1, args{&sc, event, "C[val2]", &i}, "", 7, true},
		{"expr I", ed1, args{&sc, event, "I[val3]", &i}, "37.72.10.117", 7, false},
		{"expr J", ed1, args{&sc, event, "J[val3]", &i}, "0:0:2548:a75:", 7, false},
		{"expr N", ed1, args{&sc, event, "N[val4]", &i}, "def", 7, false},
		{"expr N", ed1, args{&sc, event, "N[val1]", &i}, "0x00000101", 7, false},
		{"expr M", ed1, args{&sc, event, "M[val3]", &i}, "00-00-25-48-0a-75", 7, false},
		{"expr S", ed1, args{&sc, event, "S[val3]", &i}, "25480a75", 7, false},
		{"expr ?", ed1, args{&sc, event, "?[val3]", &i}, "?", 7, false},
		{"expr err1", ed1, args{&sc, event, "S[", &i}, "", 2, true},
		{"expr err2", ed1, args{&sc, event, "S[val3,", &i}, "", 6, true},
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
			got, err := e.calculateExpression(tt.args.sc, tt.args.event, tt.args.value, tt.args.i)
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
	event := scvd.Event{Val1: "4BY"}
//	var sc scvdm; //                                          6.ScvdData

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
		sc    *scvd.ScvdData
		event scvd.Event
		value string
		i     *int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantI   int
		wantErr bool
	}{
		{"enumExpr empty", ed1, args{&sc, event, "", &i}, "", 0, true},
		{"enumExpr E", ed1, args{&sc, event, "E[val2, typName]", &i}, "enum", 16, false},
		{"enumExpr err1", ed1, args{&sc, event, "S[", &i}, "", 2, true},
		{"enumExpr err2", ed1, args{&sc, event, "S[val3]", &i}, "", 6, true},
		{"enumExpr err3", ed1, args{&sc, event, "E[val3, xxx]", &i}, "", 12, true},
		{"enumExpr err4", ed1, args{&sc, event, "S[val3, xxx]", &i}, "", 7, true},
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
			got, err := e.calculateEnumExpression(tt.args.sc, tt.args.event, tt.args.value, tt.args.i)
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

	var ev1 scvd.Event = scvd.Event{ID: "id1", Value: "x%%%d[val1]y%u[val2]z"}
	var ev2 scvd.Event = scvd.Event{ID: "id2", Value: "x%T[val1]y%x[val2]z"}
	var ev3 scvd.Event = scvd.Event{ID: "id3", Value: "x%I[val3]y%J[val3]z"}
	var ev4 scvd.Event = scvd.Event{ID: "id4", Value: "x%M[val3]y%S[val3]z"}
	var evE1 scvd.Event = scvd.Event{ID: "idE1", Value: "x%E[val2, typName]y"}
	var everr1 scvd.Event = scvd.Event{ID: "iderr1", Value: "x%d[;]y"}
	var everr2 scvd.Event = scvd.Event{ID: "iderr2", Value: "x%E[;]y"}

	var vals = make(map[int16]string)

	vals[4711] = "enum"

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
		scvdevent scvd.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"EvalLine ev1", ed1, args{ev1}, "x%257y4711z", false},
		{"EvalLine ev2", ed1, args{ev2}, "x257y0x1267z", false},
		{"EvalLine ev3", ed1, args{ev3}, "x37.72.10.117y0:0:2548:a75:z", false},
		{"EvalLine ev4", ed1, args{ev4}, "x00-00-25-48-0a-75y25480a75z", false},
		{"EvalLine evE1", ed1, args{evE1}, "xenumy", false},
		{"EvalLine err1", ed1, args{everr1}, "", true},
		{"EvalLine err2", ed1, args{everr2}, "", true},
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
			got, err := e.EvalLine(nil, tt.args.scvdevent)
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
		{"len0", args{[]byte{}}, 0},
		{"len1", args{[]byte{0x55}}, 0x55},
		{"len3", args{[]byte{0x55, 0xAA, 0x55}}, 0xAA55},
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
		{"len0", args{[]byte{}}, 0},
		{"len1", args{[]byte{0x55}}, 0x55},
		{"len2", args{[]byte{0x55, 0xAA}}, 0xAA55},
		{"len3", args{[]byte{0x55, 0xAA, 0xFF}}, 0xFFAA55},
		{"len5", args{[]byte{0x55, 0xAA, 0x00, 0xFF, 0x11}}, 0xFF00AA55},
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
		{"len0", args{[]byte{}}, 0},
		{"len1", args{[]byte{0x55}}, 0x55},
		{"len2", args{[]byte{0x55, 0xAA}}, 0xAA55},
		{"len3", args{[]byte{0x55, 0xAA, 0xFF}}, 0xFFAA55},
		{"len4", args{[]byte{0x55, 0xAA, 0x00, 0xFF}}, 0xFF00AA55},
		{"len5", args{[]byte{0x55, 0xAA, 0x00, 0xFF, 0x01}}, 0x01FF00AA55},
		{"len6", args{[]byte{0x55, 0xAA, 0x00, 0xFF, 0x01, 0x20}}, 0x2001FF00AA55},
		{"len7", args{[]byte{0x55, 0xAA, 0x00, 0xFF, 0x01, 0x20, 0x13}}, 0x132001FF00AA55},
		{"len9", args{[]byte{0x55, 0xAA, 0x00, 0xFF, 0x01, 0x20, 0x13, 0x41, 0x77}}, 0x41132001FF00AA55},
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
	var s14 = "../../testdata/test14.binary"

	var b0 = []uint8("hello wo")
	b1 := append(b0, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0x0f)

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
		{"read ok1", fields{}, args{}, &s3, false, Data{Typ: 1, Value1: 0x6c6c6568, Value2: 0x6f77206f, Data: &b0, Time: 1410, Info: Info{0xfe00, 8, false}}, false, false},
		{"read ok2", fields{}, args{}, &s4, false, Data{Typ: 2, Value1: 1, Value2: 2, Time: 31, Info: Info{0xff00, 0, false}}, false, false},
		{"read ok3", fields{}, args{}, &s5, false, Data{Typ: 3, Value1: 805332648, Value2: 24000, Value3: 1, Value4: -65536, Time: 306, Info: Info{0xf000, 0, true}}, false, false},
		{"read ok4", fields{}, args{}, &s14, false, Data{Typ: 1, Value1: 0x6c6c6568, Value2: 0x6f77206f, Value3: 0x78563412, Value4: 0x0fdebc9a, Data: &b1, Time: 1410, Info: Info{0xfe00, 16, false}}, false, false},
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
				t.Errorf("Data.Read() %s cannot open file %v", tt.name, tt.file)
			}
			err := e.Read(tt.args.in)
			if errors.Is(err, eval.ErrEof) != tt.wantEOF {
				t.Errorf("Data.Read() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !errors.Is(err, eval.ErrEof) && (err != nil) != tt.wantErr {
				t.Errorf("Data.Read() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			b.Close()
			if !reflect.DeepEqual(e, &tt.want) {
				t.Errorf("Data.Read() %s = %v, want %v", tt.name, e, &tt.want)
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
	event := scvd.Event{Val1: "4BY"}
	var sc scvd.ScvdData
	sc.Events = make(map[uint16]scvd.Event)
	sc.Events[1] = event

	var i int

	type args struct {
		value string
		sc    *scvd.ScvdData
		event scvd.Event
		i     *int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		gen     int
		want    eval.Value
		wantErr bool
	}{
		{"val1", ed1, args{"[val1]", &sc, event, &i}, 1, eval.Value{}, false},
		//		{"data", ed2, args{"[val1]", &sc, event, &i}, 2, eval.Value{}, false},
		{"nixvar", ed2, args{"xx", &sc, event, &i}, 3, eval.Value{}, true},
		{"valxxx", ed1, args{"[valxxx]", &sc, event, &i}, 4, eval.Value{}, true},
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
				tt.want.Compose(eval.I32, 0x300066a8, 0.0, "")
			case 2:
				tt.want.Compose(eval.I32, 0x48656C6C, 0.0, "")
			}
			got, err := e.GetValue(tt.args.sc, tt.args.event, tt.args.value, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Data.GetValue() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data.GetValue() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
