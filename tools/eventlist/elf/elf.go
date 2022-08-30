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

package elf

import (
	"debug/elf"
	"errors"
	"strings"
)

type elfSection struct {
	name string
	addr uint64
	data []uint8
}

type sections struct {
	sections []*elfSection
}

var Sections sections

type symbol struct {
	addr uint64
	size uint64
}

type symbols struct {
	symbols map[string]symbol
}

var Symbols symbols

func (s *sections) Readelf(name *string) error {
	file, err := elf.Open(*name)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, section := range file.Sections {
		if section.Type == elf.SHT_PROGBITS && (section.Flags&elf.SHF_ALLOC) != 0 {
			sect := new(elfSection)
			sect.name = section.Name
			sect.addr = section.Addr
			if sect.data, err = section.Data(); err != nil {
				return err
			}
			s.sections = append(s.sections, sect)
		}
	}
	var syms []elf.Symbol
	if syms, err = file.Symbols(); err != nil && !errors.Is(err, elf.ErrNoSymbols) {
		return err
	}
	if len(Symbols.symbols) == 0 {
		Symbols.symbols = make(map[string]symbol)
	}
	for _, s := range syms {
		Symbols.symbols[s.Name] = symbol{s.Value, s.Size}
	}
	return nil
}

func (s *sections) GetString(addr uint64) string {
	for _, es := range s.sections {
		if addr >= es.addr && addr < es.addr+uint64(len(es.data)) {
			l := strings.IndexByte(string(es.data[addr-es.addr:]), 0)
			return string(es.data[addr-es.addr : addr-es.addr+uint64(l)])
		}
	}
	return ""
}

func (s *symbols) Init(name string, addr uint64, size uint64) {
	s.symbols = make(map[string]symbol)
	s.symbols[name] = symbol{addr, size}
}

func (s *symbols) GetAddrSize(name string) (addr uint64, size uint64, found bool) {
	sym, found := s.symbols[name]
	if !found {
		return 0, 0, false
	}
	return sym.addr, sym.size, true
}
