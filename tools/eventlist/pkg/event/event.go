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
	"encoding/binary"
	"errors"
	"eventlist/pkg/elf"
	"eventlist/pkg/eval"
	"eventlist/pkg/xml/scvd"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var errEnum = errors.New("invalid enum")

var errFormat = errors.New("invalid format expression")

func enumError(fn, str string) *eval.NumError {
	return &eval.NumError{Func: fn, Num: str, Err: errEnum}
}

func formatError(fn, str string) *eval.NumError {
	return &eval.NumError{Func: fn, Num: str, Err: errFormat}
}

// get the enum value as string
// count closing ]
func getEnum(typedefs map[string]map[string]map[int16]string, val int64, value string, i *int) (string, error) {
	j := strings.IndexAny(value[*i:], ":]")
	if j == -1 {
		return "", formatError("getEnum", value[*i:])
	}
	td := strings.TrimSpace(value[*i : *i+j])
	if value[*i+j] == ':' {
		members := typedefs[td] // select members of typedef
		if members == nil {
			return "", enumError("getEnum", td)
		}
		*i += j + 1
		j = strings.IndexAny(value[*i:], "]")
		if j == -1 {
			return "", formatError("getEnum", value[*i:])
		}
		md := strings.TrimSpace(value[*i : *i+j])
		entry := members[md]
		if entry == nil {
			return "", enumError("getEnum", md)
		}
		name, ok := entry[int16(val)]
		if !ok {
			return "", enumError("getEnum", strconv.Itoa(int(val)))
		}
		*i += j + 1
		return name, nil
	}
	*i += j + 1 // only enum name, no member
	for _, mm := range typedefs[td] {
		if name, ok := mm[int16(val)]; ok {
			return name, nil
		}
	}
	return "", formatError("getEnum", value[*i:])
}

type Info struct {
	ID     uint16
	length uint16
	irq    bool
}

// get the info fields from byte stream
func (info *Info) getInfoFromBytes(data []byte) {
	info.ID = convert16(data[0:2])
	info.length = convert16(data[2:4])
	info.irq = (info.length & 0x8000) != 0
	info.length &= 0x7FFF
}

func (info *Info) SplitID() (class uint16, group uint16, idx uint16, start bool) {
	class = info.ID >> 8            // should be 0xEF
	group = info.ID >> 6 & 3        // 0..3 are A..D
	idx = info.ID & 0xF             // 0..15
	start = (info.ID >> 5 & 1) == 0 // 0 is Start, 1 is Stop
	return
}

type Data struct {
	Time   uint64
	Value1 int32 // val1
	Value2 int32 // val2
	Value3 int32 // val3
	Value4 int32 // val4
	Data   *[]uint8
	Typ    uint16
	Info   Info
}

// calculate a format expression and return the result
// if unknown code then return the code only
func (e *Data) calculateExpression(value string, i *int) (string, error) {
	var val eval.Value
	var out string
	var err error

	if *i >= len(value) {
		return "", eval.ErrSyntax
	}
	c := value[*i]
	if *i+1 < len(value) && value[*i+1] == '[' {
		*i++
		val, err = e.GetValue(value, i)
		if err != nil {
			return "", err
		}
		if value[*i] != ']' {
			return "", eval.ErrSyntax
		}
		*i++
	}
	switch c {
	case 'd': // signed decimal
		out = fmt.Sprintf("%d", val.GetInt())
	case 'u': // unsigned decimal
		out = fmt.Sprintf("%d", val.GetUInt())
	case 't': // text
		out = elf.Sections.GetString(val.GetUInt())
	case 'x': // hexadecimal
		out = fmt.Sprintf("0x%02x", val.GetUInt())
	case 'F': // File
		out = elf.Sections.GetString(val.GetUInt())
		if len(out) == 0 {
			out = fmt.Sprintf("0x%08x", val.GetUInt())
		}
	case 'C': // address with file
		return "", eval.ErrSyntax
	case 'I': // IPV4
		out = fmt.Sprintf("%d.%d.%d.%d", val.GetUInt()>>24&0xFF, val.GetUInt()>>16&0xFF,
			val.GetUInt()>>8&0xFF, val.GetUInt()&0xFF)
	case 'J': // IPV6			TODO: only part of IPV6
		out = fmt.Sprintf("%x:%x:%x:%x:", val.GetUInt()>>48&0xFFFF, val.GetUInt()>>32&0xFFFF,
			val.GetUInt()>>16&0xFFFF, val.GetUInt()&0xFFFF)
	case 'N': // string address
		out = elf.Sections.GetString(val.GetUInt())
		if len(out) == 0 {
			out = fmt.Sprintf("0x%08x", val.GetUInt())
		}
	case 'M': // MAC address
		out = fmt.Sprintf("%02x-%02x-%02x-%02x-%02x-%02x", val.GetUInt()>>40&0xFF, val.GetUInt()>>32&0xFF,
			val.GetUInt()>>24&0xFF, val.GetUInt()>>16&0xFF, val.GetUInt()>>8&0xFF, val.GetUInt()&0xFF)
	case 'S': // address
		out = fmt.Sprintf("%08x", val.GetUInt())
	case 'T': // type dependant
		switch {
		case val.IsFloating(): // TODO: Float not yet possible because of event record format
			out = fmt.Sprintf("%f", val.GetFloat())
		case val.IsInteger():
			out = fmt.Sprintf("%d", val.GetInt())
		}
	case 'U': // USB descriptor
	default:
		out = string(c)
	}
	return out, nil
}

func (e *Data) calculateEnumExpression(typedefs map[string]map[string]map[int16]string,
	value string, i *int) (string, error) {
	var val eval.Value
	var out string
	var err error

	if *i >= len(value) {
		return "", eval.ErrSyntax
	}
	c := value[*i]
	if *i+1 < len(value) && value[*i+1] == '[' {
		*i++
		val, err = e.GetValue(value, i)
		if err != nil {
			return "", err
		}
		if value[*i] != ',' {
			return "", eval.ErrSyntax
		}
		*i++
	}
	if c == 'E' {
		out, err = getEnum(typedefs, val.GetInt(), value, i)
		if err != nil {
			return "", err
		}
	} else {
		return "", eval.ErrSyntax
	}
	return out, nil
}

func (e *Data) EvalLine(scvdevent scvd.Event, typedefs map[string]map[string]map[int16]string) (string, error) {
	var s string
	for i := 0; i < len(scvdevent.Value); i++ {
		c := scvdevent.Value[i]
		if c == '%' {
			if i+1 < len(scvdevent.Value) {
				i++
				c := scvdevent.Value[i]
				switch c {
				case '%':
					s += string(c)
					continue
				case 'd': // signed decimal
					fallthrough
				case 'u': // unsigned decimal
					fallthrough
				case 't': // text
					fallthrough
				case 'x': // hexadecimal
					fallthrough
				case 'F': // File
					fallthrough
				case 'C': // address with file
					fallthrough
				case 'I': // IPV4
					fallthrough
				case 'J': // IPV6
					fallthrough
				case 'N': // string address
					fallthrough
				case 'M': // MAC address
					fallthrough
				case 'S': // address
					fallthrough
				case 'T': // type dependant
					fallthrough
				case 'U': // USB descriptor
					out, err := e.calculateExpression(string(scvdevent.Value), &i)
					if err != nil {
						return "", err
					}
					s += out
					i--
				case 'E': // enum
					out, err := e.calculateEnumExpression(typedefs, string(scvdevent.Value), &i)
					if err != nil {
						return "", err
					}
					s += out
					i--
				}
			}
		} else {
			s += string(c)
		}
	}
	return s, nil
}

func (e *Data) GetValuesAsString() string {
	value := ""
	switch e.Typ {
	case 1: // EventrecordData
		value = "data=0x"
		for _, d := range *e.Data {
			value += fmt.Sprintf("%02x", d)
		}
	case 2: // Eventrecord2
		value = fmt.Sprintf("val1=0x%08x, val2=0x%08x", uint32(e.Value1), uint32(e.Value2))
	case 3: // Eventrecord4
		value = fmt.Sprintf("val1=0x%08x, val2=0x%08x, val3=0x%08x, val4=0x%08x",
			uint32(e.Value1), uint32(e.Value2), uint32(e.Value3), uint32(e.Value4))
	}
	return value
}

type Binary struct {
	file *os.File
}

func convert16(data []byte) uint16 {
	if len(data) != 2 {
		return 0
	}
	return binary.LittleEndian.Uint16([]byte{data[0], data[1]})
}

func convert32(data []byte) uint32 {
	if len(data) != 4 {
		return 0
	}
	return binary.LittleEndian.Uint32([]byte{data[0], data[1], data[2], data[3]})
}

func convert64(data []byte) uint64 {
	if len(data) != 8 {
		return 0
	}
	return binary.LittleEndian.Uint64([]byte{data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7]})
}

// get one data record
func (e *Data) Read(in *bufio.Reader) error {
	if in == nil {
		return eval.ErrEof
	}
	a2 := make([]byte, 2)
	_, err := io.ReadFull(in, a2)
	if err != nil {
		return eval.ErrEof
	}
	typ := convert16(a2)
	_, err = io.ReadFull(in, a2)
	if err != nil {
		return err
	}
	length := int(convert16(a2))
	data := make([]byte, length)
	_, err = io.ReadFull(in, data)
	if err != nil {
		return err
	}
	if len(data) < 12 {
		return eval.ErrEof
	}
	e.Time = convert64(data[:8])
	e.Info.getInfoFromBytes(data[8:12])
	e.Typ = typ
	switch typ {
	case 1: // EventrecordData
		if len(data) < 12+int(e.Info.length) {
			return eval.ErrEof
		}
		e.Data = new([]uint8)
		*e.Data = data[12 : 12+int(e.Info.length)]
	case 2: // Eventrecord2
		if len(data) < 20 {
			return eval.ErrEof
		}
		e.Value1 = int32(convert32(data[12:16]))
		e.Value2 = int32(convert32(data[16:20]))
	case 3: // Eventrecord4
		if len(data) < 28 {
			return eval.ErrEof
		}
		e.Value1 = int32(convert32(data[12:16]))
		e.Value2 = int32(convert32(data[16:20]))
		e.Value3 = int32(convert32(data[20:24]))
		e.Value4 = int32(convert32(data[24:28]))
	}
	return nil
}

func (e *Data) GetValue(value string, i *int) (eval.Value, error) {
	if *i < len(value) && value[*i] == '[' {
		if e.Data == nil {
			eval.SetVarI("val1", int64(e.Value1))
			eval.SetVarI("val2", int64(e.Value2))
			eval.SetVarI("val3", int64(e.Value3))
			eval.SetVarI("val4", int64(e.Value4))
		} else {
			ed := *e.Data
			var ed8 [8]uint8
			copy(ed8[:8], ed)
			v := uint32(ed8[0])<<24 | uint32(ed8[1])<<16 | uint32(ed8[2])<<8 | uint32(ed8[3])
			eval.SetVarI("val1", int64(v))
			v = uint32(ed8[4])<<24 | uint32(ed8[5])<<16 | uint32(ed8[6])<<8 | uint32(ed8[7])
			eval.SetVarI("val2", int64(v))
			eval.SetVarI("val3", 0)
			eval.SetVarI("val4", 0)
		}
		*i++ // skip [
		j := strings.IndexAny(value[*i:], ",]")
		var n eval.Value
		var err error
		if j == -1 {
			return eval.Value{}, eval.ErrSyntax
		}
		sid := value[*i : *i+j]
		n, err = eval.Eval(&sid)
		if err != nil {
			return eval.Value{}, err
		}
		*i += j // stay on ending char
		return n, nil
	}
	return eval.Value{}, eval.ErrSyntax
}

func (b *Binary) Open(filename *string) *bufio.Reader {
	var err error
	b.file, err = os.Open(*filename)

	if err != nil {
		return nil
	}
	return bufio.NewReader(b.file)
}

func (b *Binary) Close() error {
	return b.file.Close()
}
