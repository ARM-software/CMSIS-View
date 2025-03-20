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

// enumError creates and returns a pointer to an eval.NumError struct.
// The function takes two string parameters: fn and str, which represent
// the function name and the string that caused the error, respectively.
// It returns a pointer to an eval.NumError containing the provided
// function name, the erroneous string, and a predefined error value.
func enumError(fn, str string) *eval.NumError {
	return &eval.NumError{Func: fn, Num: str, Err: errEnum}
}

// formatError creates and returns a pointer to an eval.NumError struct.
// It takes two string parameters: fn, which represents the function name,
// and str, which represents the erroneous string value.
// The returned eval.NumError contains the provided function name, string value,
// and a predefined error indicating a formatting issue.
func formatError(fn, str string) *eval.NumError {
	return &eval.NumError{Func: fn, Num: str, Err: errFormat}
}

// getEnum retrieves the enumeration name corresponding to a given integer value from a set of typedefs.
// It parses the input string to determine the appropriate typedef and member to search within.
//
// Parameters:
// - typedefs: A map of Typedefs containing the enumeration definitions.
// - val: The integer value to look up in the enumeration.
// - value: The input string containing the typedef and member information.
// - i: A pointer to an integer representing the current position in the input string.
//
// Returns:
// - A string representing the name of the enumeration corresponding to the given integer value.
// - An error if the enumeration name could not be found or if there was a parsing error.
func getEnum(typedefs eval.Typedefs, val int64, value string, i *int) (string, error) {
	j := strings.IndexAny(value[*i:], ":]")
	if j == -1 {
		return "", formatError("getEnum", value[*i:])
	}
	td := strings.TrimSpace(value[*i : *i+j])
	if value[*i+j] == ':' {
		itypedef := typedefs[td] // select members of typedef
		if itypedef.Members == nil {
			return "", enumError("getEnum", td)
		}
		*i += j + 1
		j = strings.IndexAny(value[*i:], "]")
		if j == -1 {
			return "", formatError("getEnum", value[*i:])
		}
		md := strings.TrimSpace(value[*i : *i+j])
		entry := itypedef.Members[md].Enums
		if entry == nil {
			return "", enumError("getEnum", md)
		}
		name, ok := entry[val]
		if !ok {
			return "", enumError("getEnum", strconv.Itoa(int(val)))
		}
		*i += j + 1
		return name, nil
	}
	*i += j + 1 // only enum name, no member
	for _, mm := range typedefs[td].Members {
		if name, ok := mm.Enums[val]; ok {
			return name, nil
		}
	}
	return "", formatError("getEnum", value[*i:])
}

type Info struct {
	ID     scvd.IDType
	length uint16
	irq    bool
}

// getInfoFromBytes populates the Info struct fields from the given byte slice.
// The byte slice is expected to be at least 4 bytes long.
// The first 2 bytes are converted to the ID field.
// The next 2 bytes are converted to the length field.
// The irq field is set based on the most significant bit of the length field.
// The length field is then masked to remove the most significant bit.
func (info *Info) getInfoFromBytes(data []byte) {
	info.ID = scvd.IDType(convert16(data[0:2]))
	info.length = convert16(data[2:4])
	info.irq = (info.length & 0x8000) != 0
	info.length &= 0x7FFF
}

// SplitID splits the ID field of the Info struct into its constituent parts:
// class, group, idx, and start. The ID is expected to be a 16-bit value with
// the following structure:
// - The upper 8 bits represent the class (0xEF).
// - The next 2 bits represent the group (0..3 corresponding to A..D).
// - The next bit represents the start flag (0 for Start, 1 for Stop).
// - The lower 4 bits represent the index (0..15).
//
// Returns:
// - class: The class part of the ID.
// - group: The group part of the ID.
// - idx: The index part of the ID.
// - start: The start flag derived from the ID.
func (info *Info) SplitID() (class uint16, group uint16, idx uint16, start bool) {
	class = uint16(info.ID >> 8)     // should be 0xEF
	group = uint16(info.ID >> 6 & 3) // 0..3 are A..D
	idx = uint16(info.ID & 0xF)      // 0..15
	start = (info.ID >> 5 & 1) == 0  // 0 is Start, 1 is Stop
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

// calculateExpression evaluates a given expression based on the provided typedefs and value.
// It processes the expression character by character and returns the evaluated result as a string.
//
// Parameters:
// - typedefs: A map of type definitions used for evaluation.
// - tdUsed: A map to track used type definitions.
// - value: The expression string to be evaluated.
// - i: A pointer to the current index in the expression string.
//
// Returns:
// - A string representing the evaluated result of the expression.
// - An error if the expression contains syntax errors or evaluation fails.
//
// Supported expression characters:
// - 'd': Signed decimal
// - 'u': Unsigned decimal
// - 't': Text
// - 'x': Hexadecimal
// - 'F': File
// - 'C': Address with file (currently returns syntax error)
// - 'I': IPV4 address
// - 'J': IPV6 address (partial support)
// - 'N': String address
// - 'M': MAC address
// - 'S': Address
// - 'T': Type dependent (floating point or integer)
// - 'U': USB descriptor (currently not implemented)
// - Default: Returns the character itself as a string
func (e *Data) calculateExpression(typedefs eval.Typedefs, tdUsed map[string]string, value string, i *int) (string, error) {
	var val eval.Value
	var out string
	var err error

	if *i >= len(value) {
		return "", eval.ErrSyntax
	}
	c := value[*i]
	if *i+1 < len(value) && value[*i+1] == '[' {
		*i++
		val, err = e.GetValue(value, i, typedefs, tdUsed)
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
		case val.IsFloating():
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

// calculateEnumExpression evaluates an enum expression from the given string value starting at the position indicated by i.
// It uses the provided typedefs to resolve the enum type and value.
//
// Parameters:
//
//	typedefs - a collection of type definitions used for evaluation.
//	value - the string containing the enum expression to be evaluated.
//	i - a pointer to the current position in the value string.
//
// Returns:
//
//	A string representing the evaluated enum expression, or an error if the evaluation fails.
//
// Errors:
//
//	Returns eval.ErrSyntax if the syntax of the value string is incorrect.
//	Returns an error if the value cannot be resolved or if the enum cannot be found.
func (e *Data) calculateEnumExpression(typedefs eval.Typedefs, value string, i *int) (string, error) {
	var val eval.Value
	var out string
	var err error

	if *i >= len(value) {
		return "", eval.ErrSyntax
	}
	c := value[*i]
	if *i+1 < len(value) && value[*i+1] == '[' {
		*i++
		val, err = e.GetValue(value, i, typedefs, nil)
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

// EvalLine evaluates a line of event data based on the provided event type and typedefs.
// It processes the event's value string, replacing placeholders with corresponding values
// from the event or calculated expressions.
//
// Parameters:
//   - scvdevent: The event type containing values and the value string to be evaluated.
//   - typedefs: A collection of type definitions used for evaluating expressions.
//
// Returns:
//   - A string with the evaluated event data.
//   - An error if any issues occur during the evaluation process.
func (e *Data) EvalLine(scvdevent scvd.EventType, typedefs eval.Typedefs) (string, error) {
	var tdUsed = make(map[string]string)
	if scvdevent.Val1 != "" {
		tdUsed["val1"] = scvdevent.Val1
	}
	if scvdevent.Val2 != "" {
		tdUsed["val2"] = scvdevent.Val2
	}
	if scvdevent.Val3 != "" {
		tdUsed["val3"] = scvdevent.Val3
	}
	if scvdevent.Val4 != "" {
		tdUsed["val4"] = scvdevent.Val4
	}
	if scvdevent.Val5 != "" {
		tdUsed["val5"] = scvdevent.Val5
	}
	if scvdevent.Val6 != "" {
		tdUsed["val6"] = scvdevent.Val6
	}
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
					out, err := e.calculateExpression(typedefs, tdUsed, string(scvdevent.Value), &i)
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

// GetValuesAsString returns a string representation of the Data object
// based on its type. The format of the returned string varies depending
// on the value of the Typ field:
//
//   - If Typ is 1 (EventrecordData), the returned string is a hexadecimal
//     representation of the Data field.
//   - If Typ is 2 (Eventrecord2), the returned string includes two hexadecimal
//     values corresponding to Value1 and Value2.
//   - If Typ is 3 (Eventrecord4), the returned string includes four hexadecimal
//     values corresponding to Value1, Value2, Value3, and Value4.
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

// convert16 converts a byte slice of length 2 to a uint16 using little-endian encoding.
// If the length of the input slice is not 2, it returns 0.
//
// Parameters:
//   - data: A byte slice expected to have a length of 2.
//
// Returns:
//   - A uint16 value representing the little-endian encoded bytes from the input slice.
//   - 0 if the input slice length is not 2.
func convert16(data []byte) uint16 {
	if len(data) != 2 {
		return 0
	}
	return binary.LittleEndian.Uint16([]byte{data[0], data[1]})
}

// convert32 converts a slice of 4 bytes into a uint32 value using little-endian byte order.
// If the length of the input slice is not 4, it returns 0.
//
// Parameters:
//
//	data []byte - A slice of bytes to be converted.
//
// Returns:
//
//	uint32 - The converted 32-bit unsigned integer, or 0 if the input slice length is not 4.
func convert32(data []byte) uint32 {
	if len(data) != 4 {
		return 0
	}
	return binary.LittleEndian.Uint32([]byte{data[0], data[1], data[2], data[3]})
}

// convert64 converts a byte slice of length 8 to a uint64 using little-endian encoding.
// If the length of the input slice is not 8, it returns 0.
//
// Parameters:
//   - data: A byte slice that should be exactly 8 bytes long.
//
// Returns:
//   - A uint64 value represented by the byte slice in little-endian order.
//   - 0 if the input slice length is not 8.
func convert64(data []byte) uint64 {
	if len(data) != 8 {
		return 0
	}
	return binary.LittleEndian.Uint64([]byte{data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7]})
}

// Read reads data from the provided bufio.Reader and populates the Data struct.
// It expects the input data to be in a specific binary format and processes it accordingly.
//
// Parameters:
//   - in: A pointer to a bufio.Reader from which the data will be read.
//
// Returns:
//   - error: Returns an error if the input reader is nil, if there is an issue reading from the reader,
//     or if the data format is invalid or incomplete.
//
// The function reads the following from the input reader:
//   - A 2-byte type identifier.
//   - A 2-byte length field indicating the length of the subsequent data.
//   - A data segment of the specified length.
//
// Depending on the type identifier, the function processes the data segment and populates the appropriate
// fields in the Data struct. The supported types and their corresponding data formats are:
//   - Type 1 (EventrecordData): Expects at least 12 bytes of data plus the length specified in the Info field.
//   - Type 2 (Eventrecord2): Expects at least 20 bytes of data.
//   - Type 3 (Eventrecord4): Expects at least 28 bytes of data.
//
// If the data is successfully read and processed, the function returns nil. Otherwise, it returns an error.
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

// GetValue evaluates a value expression within a given context and returns the result.
// It supports evaluating expressions that are enclosed in square brackets and uses
// predefined variables (val1, val2, val3, val4) for evaluation.
//
// Parameters:
//   - value: The string containing the expression to be evaluated.
//   - i: A pointer to an integer representing the current position in the value string.
//   - typedefs: A collection of type definitions used during evaluation.
//   - tdUsed: A map to track which type definitions are used during evaluation.
//
// Returns:
//   - eval.Value: The result of the evaluated expression.
//   - error: An error if the evaluation fails or if there is a syntax error in the expression.
func (e *Data) GetValue(value string, i *int, typedefs eval.Typedefs, tdUsed map[string]string) (eval.Value, error) {
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
		n, err = eval.Eval(&sid, typedefs, tdUsed) // evaluate the expression
		if err != nil {
			return eval.Value{}, err
		}
		*i += j // stay on ending char
		return n, nil
	}
	return eval.Value{}, eval.ErrSyntax
}

// Open opens a file specified by the filename and returns a bufio.Reader to read from it.
// If there is an error opening the file, it returns nil.
//
// Parameters:
//   - filename: A pointer to a string containing the name of the file to open.
//
// Returns:
//   - A pointer to a bufio.Reader if the file is successfully opened, or nil if there is an error.
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
