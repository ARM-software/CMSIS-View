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

//nolint:golint,paralleltest
package scvd

import (
	"testing"
	"eventlist/pkg/eval"
)

func TestComponentViewer_getFromFile(t *testing.T) {
	var name = "../../../testdata/test.xml"
	var wrongName = "../../../testdata/xxxxx"

	type fields struct {
		Component Component
		Typedefs  Typedefs
		Events    Events
	}
	type args struct {
		name *string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"getFromFile", fields{}, args{&name}, false},
		{"getFromFile err", fields{}, args{&wrongName}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viewer := &ComponentViewer{
				Component: tt.fields.Component,
				Typedefs:  tt.fields.Typedefs,
				Events:    tt.fields.Events,
			}
			if err := viewer.getFromFile(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("ComponentViewer.getFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnum_getInfo(t *testing.T) {
	type fields struct {
		Name  string
		Value string
		Info  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int16
		wantErr bool
	}{
		{"getInfo", fields{"Name", "1+1", "Info"}, 2, false},
		{"getInfo err", fields{"Name", "??", "Info"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enum := &Enum{
				Name:  tt.fields.Name,
				Value: tt.fields.Value,
				Info:  tt.fields.Info,
			}
			got, err := enum.getInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("Enum.getInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Enum.getInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestID_getIdValue(t *testing.T) {
	var id1 ID = "2+3"
	var id2 ID = "=="

	tests := []struct {
		name    string
		id      *ID
		want    uint16
		wantErr bool
	}{
		{"getIdValue", &id1, 5, false},
		{"getIdValue err", &id2, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.id.getIdValue()
			if (err != nil) != tt.wantErr {
				t.Errorf("ID.getIdValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ID.getIdValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOne(t *testing.T) {
	var name = "../../../testdata/test.xml"
	var name1 = "../../../testdata/testBytes.xml"
	var wrongName = "../../../testdata/xxxxx"
	var nameErr1 = "../../../testdata/test_err1.xml"
	var nameErr2 = "../../../testdata/test_err2.xml"
	var nameErr3 = "../../../testdata/test_err3.xml"
	var sc ScvdData

	type args struct {
		filename *string
	}
	tests := []struct {
		name    string
		args    args
		ev      uint16
		evWant  string
		td      string
		member  string
		enum    int16
		tdWant  string
		wantErr bool
	}{
		{"getOne Bytes", args{&name1}, 0x2003, "%x[val1.B0] %x[val1.B1] %x[val1.B2] %x[val1.B3]", "attr", "B0", 1, "ready", false},
		{"getOne", args{&name}, 0xEF00, "File=fff", "attr", "member", 1, "ready", false},
		{"getOne err", args{&wrongName}, 0, "", "", "", 0, "", true},
		{"getOne err1", args{&nameErr1}, 0, "", "", "", 0, "", true},
		{"getOne err2", args{&nameErr2}, 0, "", "", "", 0, "", true},
		{"getOne err3", args{&nameErr3}, 0, "", "", "", 0, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc.Events = make(map[uint16]Event)
			for k := range eval.Typedefs {
				delete(eval.Typedefs, k)
			}
			eval.Typedefs = make(map[string]map[string]eval.TdMember)
			if err := sc.GetOne(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("getOne() error = %v, wantErr %v", err, tt.wantErr)
			}
			if string(sc.Events[tt.ev].Value) != tt.evWant {
				t.Errorf("getOne() event = %v, want %v", string(sc.Events[tt.ev].Value), tt.evWant)
			}
			if eval.Typedefs[tt.td][tt.member].Enum != nil && eval.Typedefs[tt.td][tt.member].Enum[tt.enum] != tt.tdWant {
				t.Errorf("getOne() enum = %v, want %v", eval.Typedefs[tt.td][tt.member].Enum[tt.enum], tt.tdWant)
			}
		})
	}
}

func TestGet(t *testing.T) {
	var files = []string{"../../../testdata/test.xml"}
	var files1 = []string{"../../../testdata/xxxxx"}
	var sc ScvdData

	type args struct {
		scvdFiles *[]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Get", args{&files}, false},
		{"Get err", args{&files1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := sc.Get(tt.args.scvdFiles); (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
