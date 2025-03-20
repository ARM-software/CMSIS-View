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

//nolint:golint,paralleltest
package scvd

import (
	"eventlist/pkg/eval"
	"testing"
)

func TestComponentViewer_getFromFile(t *testing.T) {
	var name = "../../../testdata/test.xml"
	var wrongName = "../../../testdata/xxxxx"

	type fields struct {
		Component ComponentsType
		Typedefs  TypedefsType
		Events    EventsType
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

func Test_getOne(t *testing.T) {
	var name = "../../../testdata/test.xml"
	var name1 = "../../../testdata/test1.xml"
	var wrongName = "../../../testdata/xxxxx"
	var nameErr1 = "../../../testdata/test_err1.xml"
	var nameErr2 = "../../../testdata/test_err2.xml"
	var nameErr3 = "../../../testdata/test_err3.xml"
	var evs = make(Events)
	var tds = make(eval.Typedefs)

	type args struct {
		filename *string
		events   Events
		typedefs eval.Typedefs
	}
	tests := []struct {
		name    string
		args    args
		ev      IDType
		evWant  string
		td      string
		member  string
		enum    int64
		tdWant  string
		wantErr bool
	}{
		{"getOne", args{&name, evs, tds}, 0xEF00, "File=fff", "attr", "member", 1, "ready", false},
		{"getOne1", args{&name1, evs, tds}, 0x2003, "%x[val1.B0] %x[val1.B1] %x[val1.B2] %x[val1.B3]", "BY4", "B2", -1, "2", false},
		{"getOne err", args{&wrongName, evs, tds}, 0, "", "", "", 0, "", true},
		{"getOne err1", args{&nameErr1, evs, tds}, 0, "", "", "", 0, "", true},
		{"getOne err2", args{&nameErr2, evs, tds}, 0, "", "", "", 0, "", true},
		{"getOne err3", args{&nameErr3, evs, tds}, 0, "", "", "", 0, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getOne(tt.args.filename, tt.args.events, tt.args.typedefs); (err != nil) != tt.wantErr {
				t.Errorf("getOne() error = %v, wantErr %v", err, tt.wantErr)
			}
			if string(evs[tt.ev].Value) != tt.evWant {
				t.Errorf("getOne() event = %v, want %v", string(evs[tt.ev].Value), tt.evWant)
			}
			if tt.enum == -1 {
				if tt.args.events[tt.ev].Val1 != tt.td {
					t.Errorf("getOne() val1 = %v, want %v", tt.args.events[tt.ev].Val1, tt.td)
				}
				if tds[tt.td].Members[tt.member].Offset != tt.tdWant {
					t.Errorf("getOne() enum = %v, want %v", tds[tt.td].Members[tt.member].Offset, tt.tdWant)
				}
			} else {
				if tds[tt.td].Members[tt.member].Enums[tt.enum] != tt.tdWant {
					t.Errorf("getOne() enum = %v, want %v", tds[tt.td].Members[tt.member].Enums[tt.enum], tt.tdWant)
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	var files = []string{"../../../testdata/test.xml"}
	var files1 = []string{"../../../testdata/xxxxx"}
	var evs = make(Events)
	var tds = make(eval.Typedefs)

	type args struct {
		scvdFiles *[]string
		events    Events
		typedefs  eval.Typedefs
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Get", args{&files, evs, tds}, false},
		{"Get err", args{&files1, evs, tds}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Get(tt.args.scvdFiles, tt.args.events, tt.args.typedefs); (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
