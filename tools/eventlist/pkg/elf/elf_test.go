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

package elf

import (
	"reflect"
	"testing"
)

func Test_sections_Readelf(t *testing.T) { //nolint:golint,paralleltest
	fileTest := "../../testdata/elftest.elf"
	fileNix := "../../testdata/nix.elf"
	fileSym := "../../testdata/elfsym.elf"

	type args struct {
		name   *string
		symbol string
	}
	tests := []struct {
		name    string
		s       *sections
		args    args
		want    uint64
		wantErr bool
	}{
		{"Sym", &sections{}, args{&fileSym, "LEDOn"}, 0x38000178, false},
		{"test", &sections{}, args{name: &fileTest}, 0, false},
		{"errName", &sections{}, args{name: &fileNix}, 0, true},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = tt.s.Readelf(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("sections.Readelf() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if err == nil && tt.args.name == &fileTest && tt.s.GetString(0x4010) != "def" {
				t.Errorf("sections.Readelf() %s data not found", tt.name)
			}
			if err == nil && tt.args.name == &fileSym {
				a, _, f := Symbols.GetAddrSize(tt.args.symbol)
				if !f {
					t.Errorf("sections.Readelf() %s symbol not found", tt.name)
				} else if a != tt.want {
					t.Errorf("sections.Readelf() %s = %v, want %v", tt.name, a, tt.want)
				}
			}
		})
	}
}

func TestGetString(t *testing.T) {
	t.Parallel()

	sect1 := &elfSection{"", 123, []uint8{'a', 'b', 'c', 0}}
	sect2 := &elfSection{"", 100, []uint8{0, 1, 2, 'd', 'e', 'f', 0}}

	type args struct {
		addr uint64
	}
	tests := []struct {
		name string
		e    *elfSection
		s    *sections
		args args
		want string
	}{
		{"test_1", sect1, &sections{}, args{123}, "abc"},
		{"test_2", sect2, &sections{}, args{103}, "def"},
		{"test_err1", sect2, &sections{}, args{99}, ""},
		{"test_err2", sect1, &sections{}, args{127}, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.s.sections = append(tt.s.sections, tt.e)
			if got := tt.s.GetString(tt.args.addr); got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_symbols_Init(t *testing.T) {
	t.Parallel()

	type fields struct {
		symbols map[string]symbol
	}
	type args struct {
		name string
		addr uint64
		size uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   symbol
	}{
		{"test", fields{}, args{"n", 123, 456}, symbol{123, 456}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &symbols{
				symbols: tt.fields.symbols,
			}
			s.Init(tt.args.name, tt.args.addr, tt.args.size)
			if !reflect.DeepEqual(s.symbols[tt.args.name], tt.want) {
				t.Errorf("Test_symbols.Init() %s = %v, want %v", tt.name, s, tt.want)
			}
		})
	}
}

func Test_symbols_GetAddrSize(t *testing.T) {
	t.Parallel()

	var syms = make(map[string]symbol)

	syms["symbol"] = symbol{123, 456}

	type fields struct {
		symbols map[string]symbol
	}
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantAddr  uint64
		wantSize  uint64
		wantFound bool
	}{
		{"sym_ok", fields{syms}, args{"symbol"}, 123, 456, true},
		{"sym_fail", fields{syms}, args{"s"}, 0, 0, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &symbols{
				symbols: tt.fields.symbols,
			}
			gotAddr, gotSize, gotFound := s.GetAddrSize(tt.args.name)
			if gotAddr != tt.wantAddr {
				t.Errorf("symbols.GetAddrSize() gotAddr = %v, want %v", gotAddr, tt.wantAddr)
			}
			if gotSize != tt.wantSize {
				t.Errorf("symbols.GetAddrSize() gotSize = %v, want %v", gotSize, tt.wantSize)
			}
			if gotFound != tt.wantFound {
				t.Errorf("symbols.GetAddrSize() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}
