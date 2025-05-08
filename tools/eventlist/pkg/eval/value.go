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

import "eventlist/pkg/elf"

type Value struct {
	t Token
	i int64
	f float64
	s string
	v *Variable
	l []Value
}

// Compose sets the fields of the Value struct with the provided parameters.
//
// Parameters:
//   - t: A Token representing the type of the value.
//   - i: An int64 representing an integer value.
//   - f: A float64 representing a floating-point value.
//   - s: A string representing a string value.
func (v *Value) Compose(t Token, i int64, f float64, s string) {
	*v = Value{t, i, f, s, nil, nil}
}

// getValue retrieves the value stored in the Value object.
// If the value is nil, it returns an error indicating that the value is not a variable.
// Otherwise, it delegates the call to the underlying value's getValue method.
//
// Returns:
//   - Value: The retrieved value.
//   - error: An error if the value is not a variable.
func (v *Value) getValue() (Value, error) {
	if v.v == nil {
		return *v, typeError("not a variable", "")
	}
	return v.v.getValue()
}

// setValue assigns the value of v1 to the receiver Value v.
// If the receiver's value is nil, it returns a type error indicating that
// the receiver is not a variable. Otherwise, it delegates the assignment
// to the receiver's internal value and returns any error encountered during
// the assignment process.
//
// Parameters:
//
//	v1 - The Value to be assigned to the receiver.
//
// Returns:
//
//	An error if the assignment fails or if the receiver is not a variable.
func (v *Value) setValue(v1 *Value) error {
	if v.v == nil {
		return typeError("not a variable", "")
	}
	err := v.v.setValue(v1) // do not change v yet
	return err
}

// addList adds a Value to the list contained within the receiver Value.
// If the receiver Value is of type Nix, it initializes it as a List.
// If the receiver Value is not of type List, it returns a type error.
//
// Parameters:
//
//	v1 - The Value to be added to the list.
//
// Returns:
//
//	An error if the receiver Value is not of type List, otherwise nil.
func (v *Value) addList(v1 Value) error {
	if v.t == Nix {
		v.t = List
	} else if v.t != List {
		return typeError("not a list", "")
	}
	v.l = append(v.l, v1)
	return nil
}

// GetInt returns the integer representation of the Value.
// If the Value is of type Integer, it returns the integer value directly.
// If the Value is of type Floating, it converts the floating-point value to an integer.
// If the Value is of any other type, it returns 0.
func (v *Value) GetInt() int64 {
	switch v.t {
	case Integer:
		return v.i
	case Floating:
		return int64(v.f)
	}
	return 0
}

// GetUInt returns the unsigned integer representation of the Value.
// If the Value is of type Integer, it converts the integer to uint64.
// If the Value is of type Floating, it converts the floating-point number to uint64.
// If the Value is of any other type, it returns 0.
func (v *Value) GetUInt() uint64 {
	switch v.t {
	case Integer:
		return uint64(v.i)
	case Floating:
		return uint64(v.f)
	}
	return 0
}

// GetFloat returns the float64 representation of the Value.
// If the Value is of type Integer, it converts the integer to a float64.
// If the Value is of type Floating, it returns the float64 value directly.
// If the Value is of any other type, it returns 0.0.
func (v *Value) GetFloat() float64 {
	switch v.t {
	case Integer:
		return float64(v.i)
	case Floating:
		return v.f
	}
	return 0.0
}

// GetList returns the list of Value objects if the current Value is a list.
// If the current Value is not a list, it returns nil.
func (v *Value) GetList() []Value {
	if v.IsList() {
		return v.l
	}
	return nil
}

// IsInteger checks if the Value type is an integer.
// It returns true if the Value type is Integer, otherwise false.
func (v *Value) IsInteger() bool {
	return v.t == Integer
}

// IsFloating checks if the Value type is Floating.
// It returns true if the Value type is Floating, otherwise false.
func (v *Value) IsFloating() bool {
	return v.t == Floating
}

// IsString checks if the Value type is a string.
// It returns true if the Value type is String, otherwise false.
func (v *Value) IsString() bool {
	return v.t == String
}

// IsIdentifier checks if the Value instance represents an identifier.
// It returns true if the type of the Value (v.t) is Identifier, otherwise false.
func (v *Value) IsIdentifier() bool {
	return v.t == Identifier
}

// IsList checks if the Value instance is of type List.
// It returns true if the type is List, otherwise false.
func (v *Value) IsList() bool {
	return v.t == List
}

type FuncNo int

type Function struct {
	fno     FuncNo
	params  int
	parType Token
	ret     Token
}

const (
	CALCMEMUSED FuncNo = iota
	GETREGVAL
	SYMBOLEXIST
	FINDSYMBOL
	OFFSETOF
	SIZEOF
)

var fctMap = map[string]Function{
	"__CalcMemUsed":   {CALCMEMUSED, 4, Integer, Integer},
	"__GetRegVal":     {GETREGVAL, 1, String, Integer},
	"__Symbol_exists": {SYMBOLEXIST, 1, String, Integer},
	"__FindSymbol":    {FINDSYMBOL, 1, String, Integer},
	"__Offset_of":     {OFFSETOF, 1, String, Integer},
	"__size_of":       {SIZEOF, 1, String, Integer},
}

// Function evaluates the function represented by the Value receiver (v) using the provided argument (v1).
// It performs several checks to ensure that the function and its arguments are valid:
// - v1 must not be nil.
// - v must be an identifier.
// - v1 must be a list.
// - The function must exist in the function map (fctMap).
// - The number of parameters in the function must match the length of the list in v1.
// - Each parameter type in the list must match the expected parameter type of the function.
//
// Depending on the function number (fno), it performs specific operations and sets the result in the receiver (v).
// The possible operations include calculating memory usage, getting register values, checking symbol existence,
// finding symbols, getting the offset of a symbol, and getting the size of a symbol.
//
// Returns an error if any of the checks fail or if the function cannot be evaluated.
func (v *Value) Function(v1 *Value) error {
	if v1 == nil {
		return typeError("Function", "")
	}
	if !v.IsIdentifier() {
		return typeError("Function", "")
	}
	if !v1.IsList() {
		return typeError("Function", "")
	}
	var f Function
	var found bool
	if f, found = fctMap[v.s]; !found {
		return typeError("Function", "")
	}
	if f.params != len(v1.GetList()) {
		return typeError("Function", "")
	}
	for _, par := range v1.GetList() {
		if f.parType != par.t {
			return typeError("Function", "")
		}
	}
	switch f.fno {
	case CALCMEMUSED:
		*v = Value{t: f.ret, i: 0}
	case GETREGVAL:
		*v = Value{t: f.ret, i: 0}
	case SYMBOLEXIST:
		_, _, flag := elf.Symbols.GetAddrSize(v1.GetList()[0].s)
		if flag {
			*v = Value{t: f.ret, i: 1}
		} else {
			*v = Value{t: f.ret, i: 0}
		}
	case FINDSYMBOL:
		_, _, flag := elf.Symbols.GetAddrSize(v1.GetList()[0].s)
		if flag {
			*v = Value{t: f.ret, i: 1}
		} else {
			*v = Value{t: f.ret, i: 0}
		}
	case OFFSETOF:
		a, _, flag := elf.Symbols.GetAddrSize(v1.GetList()[0].s)
		if flag {
			*v = Value{t: f.ret, i: int64(a)}
		} else {
			*v = Value{t: f.ret, i: 0}
		}
	case SIZEOF:
		_, s, flag := elf.Symbols.GetAddrSize(v1.GetList()[0].s)
		if flag {
			*v = Value{t: f.ret, i: int64(s)}
		} else {
			*v = Value{t: f.ret, i: 0}
		}
	}
	return nil
}

// Extract extracts a portion of the integer value stored in the Value struct.
// The size (sz) and offset (off) parameters specify the number of bytes to extract
// and the starting byte position, respectively. The extracted portion is then
// right-shifted by the offset and stored back in the Value struct.
//
// Parameters:
//
//	sz  - The number of bytes to extract.
//	off - The starting byte position for extraction.
//
// Returns:
//
//	An error if the type of the value is not Integer.
func (v *Value) Extract(sz uint32, bigEndian bool, off uint32) error {
	if v.t != Integer {
		return typeError("Extract", "")
	}
	tmp := uint64(v.i)
	if bigEndian {
		tmp = (tmp&0xFF)<<56 | (tmp&0xFF00)<<40 | (tmp&0xFF0000)<<24 | (tmp&0xFF000000)<<8 |
			(tmp&0xFF00000000)>>8 | (tmp&0xFF0000000000)>>24 | (tmp&0xFF000000000000)>>40 | (tmp&0xFF00000000000000)>>56
	}
	tmp &= (1 << (sz * 8)) - 1
	v.i = int64(tmp >> (off * 8))
	return nil
}

// Inc increments the value of the Value receiver based on its type.
// If the type is Integer, it increments the integer value by 1.
// If the type is Floating, it increments the floating-point value by 1.
// If the type is neither Integer nor Floating, it returns a type error.
func (v *Value) Inc() error {
	switch v.t {
	case Integer:
		v.i++
	case Floating:
		v.f++
	default:
		return typeError("Inc", "")
	}
	return nil
}

// Dec decrements the value of the Value receiver by 1.
// If the type of the value is Integer, it decrements the integer value.
// If the type of the value is Floating, it decrements the floating-point value.
// If the type is neither Integer nor Floating, it returns a type error.
func (v *Value) Dec() error {
	switch v.t {
	case Integer:
		v.i--
	case Floating:
		v.f--
	default:
		return typeError("Dec", "")
	}
	return nil
}

// Plus performs an operation based on the type of the Value receiver.
// It currently supports Integer and Floating types.
// If the type is not supported, it returns a typeError.
//
// Returns:
//   - error: typeError if the type is not supported, otherwise nil.
func (v *Value) Plus() error {
	switch v.t {
	case Integer:
	case Floating:
	default:
		return typeError("Plus", "")
	}
	return nil
}

// Neg negates the value of the Value receiver. If the type of the value is
// Integer, it negates the integer value. If the type is Floating, it negates
// the floating-point value. If the type is neither Integer nor Floating, it
// returns a typeError indicating that negation is not supported for the type.
func (v *Value) Neg() error {
	switch v.t {
	case Integer:
		v.i = -v.i
	case Floating:
		v.f = -v.f
	default:
		return typeError("Neg", "")
	}
	return nil
}

// Compl performs a bitwise complement operation on the value if it is of type Integer.
// If the value is not an Integer, it returns a type error.
// Returns an error if the operation is not applicable to the value type.
func (v *Value) Compl() error {
	switch v.t {
	case Integer:
		v.i = ^v.i
	default:
		return typeError("Compl", "")
	}
	return nil
}

// Not negates the value of the Value receiver if it is of type Integer.
// If the integer value is 0, it sets it to 1. Otherwise, it sets it to 0.
// Returns an error if the type of the Value receiver is not Integer.
func (v *Value) Not() error {
	switch v.t {
	case Integer:
		if v.i == 0 {
			v.i = 1
		} else {
			v.i = 0
		}
	default:
		return typeError("Compl", "")
	}
	return nil
}

// Cast converts the value to the specified type.
// It supports casting between integer and floating-point types.
// If the conversion is not possible, it returns a typeError.
//
// Supported types for casting:
// - Uint8
// - Int8
// - Uint16
// - Int16
// - Uint32
// - Int32
// - Uint64
// - Int64
// - Float
// - Double
//
// Parameters:
// - ty: The target type to cast the value to.
//
// Returns:
// - error: Returns a typeError if the cast is not possible.
func (v *Value) Cast(ty Type) error {
	switch ty {
	case Uint8:
		switch v.t {
		case Integer:
			v.i = int64(uint8(v.i))
		case Floating:
			v.i = int64(uint8(v.f))
			v.t = Integer
			v.f = 0
		default:
			return typeError("Cast", "")
		}
	case Int8:
		switch v.t {
		case Integer:
			v.i = int64(int8(v.i))
		case Floating:
			v.i = int64(int8(v.f))
			v.t = Integer
			v.f = 0
		default:
			return typeError("Cast", "")
		}
	case Uint16:
		switch v.t {
		case Integer:
			v.i = int64(uint16(v.i))
		case Floating:
			v.i = int64(uint16(v.f))
			v.t = Integer
			v.f = 0
		default:
			return typeError("Cast", "")
		}
	case Int16:
		switch v.t {
		case Integer:
			v.i = int64(int16(v.i))
		case Floating:
			v.i = int64(int16(v.f))
			v.t = Integer
			v.f = 0
		default:
			return typeError("Cast", "")
		}
	case Uint32:
		switch v.t {
		case Integer:
			v.i = int64(uint32(v.i))
		case Floating:
			v.i = int64(uint32(v.f))
			v.t = Integer
			v.f = 0
		default:
			return typeError("Cast", "")
		}
	case Int32:
		switch v.t {
		case Integer:
			v.i = int64(int32(v.i))
		case Floating:
			v.i = int64(int32(v.f))
			v.t = Integer
			v.f = 0
		default:
			return typeError("Cast", "")
		}
	case Uint64:
		switch v.t {
		case Integer:
			v.i = int64(uint64(v.i))
		case Floating:
			v.i = int64(uint64(v.f))
			v.t = Integer
			v.f = 0
		default:
			return typeError("Cast", "")
		}
	case Int64:
		switch v.t {
		case Integer:
		case Floating:
			v.i = int64(v.f)
			v.t = Integer
			v.f = 0
		default:
			return typeError("Cast", "")
		}
	case Float:
		switch v.t {
		case Integer:
			v.f = float64(float32(v.i))
			v.t = Floating
			v.i = 0
		case Floating:
			v.f = float64(float32(v.f))
		default:
			return typeError("Cast", "")
		}
	case Double:
		switch v.t {
		case Integer:
			v.f = float64(v.i)
			v.t = Floating
			v.i = 0
		case Floating:
		default:
			return typeError("Cast", "")
		}
	}
	return nil
}

// Mul multiplies the value of the receiver with the value of v1.
// The result is stored in the receiver. The types of the values
// are considered during the multiplication:
// - If both values are integers, the result is an integer.
// - If one value is a floating point, the result is a floating point.
// If the types are incompatible, an error is returned.
//
// Parameters:
// - v1: The value to multiply with the receiver.
//
// Returns:
// - error: An error if the types are incompatible, otherwise nil.
func (v *Value) Mul(v1 *Value) error {
	switch v.t {
	case Integer:
		switch v1.t {
		case Integer:
			v.i *= v1.i
		case Floating:
			v.f = float64(v.i) * v1.f
			v.t = Floating
			v.i = 0
		default:
			return typeError("Mul", "")
		}
	case Floating:
		switch v1.t {
		case Integer:
			v.f = v.f * float64(v1.i)
			v.t = Floating
		case Floating:
			v.f *= v1.f
		default:
			return typeError("Mul", "")
		}
	default:
		return typeError("Mul", "")
	}
	return nil
}

// Div performs division of the current Value (v) by another Value (v1).
// It supports division for both Integer and Floating types.
// If the division is by zero, it returns an error indicating "division by 0".
// If the types of v and v1 are incompatible for division, it returns a type error.
//
// Parameters:
// - v1: The Value to divide the current Value by.
//
// Returns:
// - error: An error if the division is by zero or if the types are incompatible, otherwise nil.
func (v *Value) Div(v1 *Value) error {
	switch v.t {
	case Integer:
		switch v1.t {
		case Integer:
			if v1.i == 0 {
				return typeError("division by 0", "")
			}
			v.i /= v1.i
		case Floating:
			if v1.f == 0.0 {
				return typeError("division by 0", "")
			}
			v.f = float64(v.i) / v1.f
			v.t = Floating
			v.i = 0
		default:
			return typeError("Div", "")
		}
	case Floating:
		switch v1.t {
		case Integer:
			if v1.i == 0 {
				return typeError("division by 0", "")
			}
			v.f = v.f / float64(v1.i)
			v.t = Floating
		case Floating:
			if v1.f == 0.0 {
				return typeError("division by 0", "")
			}
			v.f /= v1.f
		default:
			return typeError("Div", "")
		}
	default:
		return typeError("Div", "")
	}
	return nil
}

// Mod performs the modulus operation on the current Value (v) with another Value (v1).
// It supports only integer types and returns an error if the operation is not valid.
//
// If v or v1 is not an integer, or if v1 is zero, an appropriate error is returned.
//
// Parameters:
// - v1: The Value to perform the modulus operation with.
//
// Returns:
// - error: An error if the modulus operation is invalid, otherwise nil.
func (v *Value) Mod(v1 *Value) error {
	switch v.t {
	case Integer:
		switch v1.t {
		case Integer:
			if v1.i == 0 {
				return typeError("modular by 0", "")
			}
			v.i %= v1.i
		case Floating:
			return typeError("mod with floatings", "")
		default:
			return typeError("Mod", "")
		}
	case Floating:
		return typeError("mod with floatings", "")
	default:
		return typeError("Mod", "")
	}
	return nil
}

// Add adds the value of v1 to the receiver v. The addition is performed
// based on the type of the values. If both values are integers, their
// integer values are added. If one of the values is a floating-point
// number, the result is a floating-point number. If the types are not
// compatible, a typeError is returned.
//
// Parameters:
//
//	v1 - The value to be added to the receiver.
//
// Returns:
//
//	An error if the types are not compatible for addition, otherwise nil.
func (v *Value) Add(v1 *Value) error {
	switch v.t {
	case Integer:
		switch v1.t {
		case Integer:
			v.i += v1.i
		case Floating:
			v.f = float64(v.i) + v1.f
			v.i = 0
			v.t = Floating
		default:
			return typeError("Add", "")
		}
	case Floating:
		switch v1.t {
		case Integer:
			v.f += float64(v1.i)
		case Floating:
			v.f += v1.f
		default:
			return typeError("Add", "")
		}
	default:
		return typeError("Add", "")
	}
	return nil
}

// Sub subtracts the value of v1 from the receiver v. The result is stored in the receiver v.
// It supports subtraction for Integer and Floating types. If the types of v and v1 do not match,
// it will convert the Integer to Floating if necessary. If the types are not supported, it returns an error.
//
// Parameters:
//   - v1: The value to be subtracted from the receiver v.
//
// Returns:
//   - error: An error if the types are not supported for subtraction.
func (v *Value) Sub(v1 *Value) error {
	switch v.t {
	case Integer:
		switch v1.t {
		case Integer:
			v.i -= v1.i
		case Floating:
			v.f = float64(v.i) - v1.f
			v.i = 0
			v.t = Floating
		default:
			return typeError("Sub", "")
		}
	case Floating:
		switch v1.t {
		case Integer:
			v.f -= float64(v1.i)
		case Floating:
			v.f -= v1.f
		default:
			return typeError("Sub", "")
		}
	default:
		return typeError("Sub", "")
	}
	return nil
}

// Shl performs a left shift operation on the integer value of the receiver
// by the integer value of the provided Value v1. If either the receiver or
// v1 is not of type Integer, it returns a typeError.
//
// Parameters:
// - v1: A pointer to a Value that provides the shift amount.
//
// Returns:
// - error: An error if the types of the receiver or v1 are not Integer, otherwise nil.
func (v *Value) Shl(v1 *Value) error {
	if v.t != Integer || v1.t != Integer {
		return typeError("shl", "")
	}
	v.i <<= v1.i
	return nil
}

// Shr performs a bitwise right shift operation on the integer value of the receiver
// by the integer value of the provided Value v1. If either value is not of type Integer,
// it returns a typeError.
//
// Parameters:
//
//	v1 - The Value containing the integer by which the receiver's integer value will be shifted.
//
// Returns:
//
//	An error if either value is not of type Integer, otherwise nil.
func (v *Value) Shr(v1 *Value) error {
	if v.t != Integer || v1.t != Integer {
		return typeError("shr", "")
	}
	v.i >>= v1.i
	return nil
}

// Less compares the current Value (v) with another Value (v1) and sets the
// current Value (v) to 1 if it is less than v1, otherwise sets it to 0.
// It supports comparison between Integer and Floating types. If the types
// are incompatible, it returns a typeError.
//
// The resulting type always is forced to be Integer.
//
// Parameters:
// - v1: A pointer to the Value to compare with.
//
// Returns:
// - error: An error if the types are incompatible, otherwise nil.
func (v *Value) Less(v1 *Value) error {
	switch v.t {
	case Integer:
		switch v1.t {
		case Integer:
			if v.i < v1.i {
				v.i = 1
			} else {
				v.i = 0
			}
		case Floating:
			if float64(v.i) < v1.f {
				v.i = 1
			} else {
				v.i = 0
			}
		default:
			return typeError("Less", "")
		}
	case Floating:
		switch v1.t {
		case Integer:
			if v.f < float64(v1.i) {
				v.i = 1
			} else {
				v.i = 0
			}
			v.t = Integer
			v.f = 0.0
		case Floating:
			if v.f < v1.f {
				v.i = 1
			} else {
				v.i = 0
			}
			v.t = Integer
			v.f = 0.0
		default:
			return typeError("Less", "")
		}
	default:
		return typeError("Less", "")
	}
	return nil
}

// LessEqual compares the value of the current Value object with another Value object (v1).
// It supports comparison between Integer and Floating types. If the current Value is less than
// or equal to v1, it sets the current Value to 1 (true), otherwise sets it to 0 (false).
// If the types are incompatible for comparison, it returns a typeError.
//
// The resulting type always is forced to be Integer.
//
// Parameters:
// - v1: A pointer to another Value object to compare with.
//
// Returns:
// - error: An error if the types are incompatible for comparison, otherwise nil.
func (v *Value) LessEqual(v1 *Value) error {
	switch v.t {
	case Integer:
		switch v1.t {
		case Integer:
			if v.i <= v1.i {
				v.i = 1
			} else {
				v.i = 0
			}
		case Floating:
			if float64(v.i) <= v1.f {
				v.i = 1
			} else {
				v.i = 0
			}
		default:
			return typeError("LessEqual", "")
		}
	case Floating:
		switch v1.t {
		case Integer:
			if v.f <= float64(v1.i) {
				v.i = 1
			} else {
				v.i = 0
			}
			v.t = Integer
			v.f = 0.0
		case Floating:
			if v.f <= v1.f {
				v.i = 1
			} else {
				v.i = 0
			}
			v.t = Integer
			v.f = 0.0
		default:
			return typeError("LessEqual", "")
		}
	default:
		return typeError("LessEqual", "")
	}
	return nil
}

// Greater compares the current Value with another Value v1 and sets the current Value to 1 if it is greater than v1,
// or 0 if it is not. The comparison is based on the type of the Values (Integer or Floating).
// If the types are incompatible for comparison, it returns a type error.
//
// The resulting type always is forced to be Integer.
//
// Parameters:
//
//	v1 *Value - The Value to compare with the current Value.
//
// Returns:
//
//	error - Returns a type error if the types of the Values are incompatible for comparison, otherwise returns nil.
func (v *Value) Greater(v1 *Value) error {
	switch v.t {
	case Integer:
		switch v1.t {
		case Integer:
			if v.i > v1.i {
				v.i = 1
			} else {
				v.i = 0
			}
		case Floating:
			if float64(v.i) > v1.f {
				v.i = 1
			} else {
				v.i = 0
			}
		default:
			return typeError("Greater", "")
		}
	case Floating:
		switch v1.t {
		case Integer:
			if v.f > float64(v1.i) {
				v.i = 1
			} else {
				v.i = 0
			}
			v.t = Integer
			v.f = 0.0
		case Floating:
			if v.f > v1.f {
				v.i = 1
			} else {
				v.i = 0
			}
			v.t = Integer
			v.f = 0.0
		default:
			return typeError("Greater", "")
		}
	default:
		return typeError("Greater", "")
	}
	return nil
}

// GreaterEqual compares the value of the current Value object (v) with another Value object (v1)
// and sets the current Value object to 1 if it is greater than or equal to v1, otherwise sets it to 0.
// It supports comparisons between Integer and Floating types. If the types are incompatible, it returns an error.
//
// The resulting type always is forced to be Integer.
//
// Parameters:
// - v1: A pointer to the Value object to compare with.
//
// Returns:
// - error: An error if the types of the Value objects are incompatible for comparison.
func (v *Value) GreaterEqual(v1 *Value) error {
	switch v.t {
	case Integer:
		switch v1.t {
		case Integer:
			if v.i >= v1.i {
				v.i = 1
			} else {
				v.i = 0
			}
		case Floating:
			if float64(v.i) >= v1.f {
				v.i = 1
			} else {
				v.i = 0
			}
		default:
			return typeError("GreaterEqual", "")
		}
	case Floating:
		switch v1.t {
		case Integer:
			if v.f >= float64(v1.i) {
				v.i = 1
			} else {
				v.i = 0
			}
			v.t = Integer
			v.f = 0.0
		case Floating:
			if v.f >= v1.f {
				v.i = 1
			} else {
				v.i = 0
			}
			v.t = Integer
			v.f = 0.0
		default:
			return typeError("GreaterEqual", "")
		}
	default:
		return typeError("GreaterEqual", "")
	}
	return nil
}

// Equal compares the value of the current Value object with another Value object (v1).
// It supports comparison between Integer and Floating types.
// If the values are equal, it sets the current Value's integer field (v.i) to 1, otherwise to 0.
//
// The resulting type always is forced to be Integer.
//
// Parameters:
//
//	v1 - The Value object to compare with the current Value object.
//
// Returns:
//
//	error - An error if the types are not supported for comparison, otherwise nil.
func (v *Value) Equal(v1 *Value) error {
	switch v.t {
	case Integer:
		switch v1.t {
		case Integer:
			if v.i == v1.i {
				v.i = 1
			} else {
				v.i = 0
			}
		case Floating:
			if float64(v.i) == v1.f {
				v.i = 1
			} else {
				v.i = 0
			}
		default:
			return typeError("Equal", "")
		}
	case Floating:
		switch v1.t {
		case Integer:
			if v.f == float64(v1.i) {
				v.i = 1
			} else {
				v.i = 0
			}
			v.t = Integer
			v.f = 0.0
		case Floating:
			if v.f == v1.f {
				v.i = 1
			} else {
				v.i = 0
			}
			v.t = Integer
			v.f = 0.0
		default:
			return typeError("Equal", "")
		}
	default:
		return typeError("Equal", "")
	}
	return nil
}

// NotEqual compares the value of the current Value object with another Value object (v1).
// It sets the current Value object to 1 if they are not equal, and 0 if they are equal.
// The comparison is based on the type of the Value objects (Integer or Floating).
// If the types are incompatible, it returns a type error.
//
// The resulting type always is forced to be Integer.
//
// Parameters:
//
//	v1 (*Value): The Value object to compare with.
//
// Returns:
//
//	error: An error if the types are incompatible, otherwise nil.
func (v *Value) NotEqual(v1 *Value) error {
	switch v.t {
	case Integer:
		switch v1.t {
		case Integer:
			if v.i != v1.i {
				v.i = 1
			} else {
				v.i = 0
			}
		case Floating:
			if float64(v.i) != v1.f {
				v.i = 1
			} else {
				v.i = 0
			}
		default:
			return typeError("NotEqual", "")
		}
	case Floating:
		switch v1.t {
		case Integer:
			if v.f != float64(v1.i) {
				v.i = 1
			} else {
				v.i = 0
			}
			v.t = Integer
			v.f = 0.0
		case Floating:
			if v.f != v1.f {
				v.i = 1
			} else {
				v.i = 0
			}
			v.t = Integer
			v.f = 0.0
		default:
			return typeError("NotEqual", "")
		}
	default:
		return typeError("NotEqual", "")
	}
	return nil
}

// And performs a bitwise AND operation between the current Value and another Value (v1).
// Both Values must be of Integer type. If either Value is not an Integer, a typeError is returned.
// The result of the operation is stored in the current Value.
//
// Parameters:
//
//	v1 - The Value to perform the bitwise AND operation with.
//
// Returns:
//
//	An error if either Value is not of Integer type, otherwise nil.
func (v *Value) And(v1 *Value) error {
	if v.t != Integer || v1.t != Integer {
		return typeError("And", "")
	}
	v.i &= v1.i
	return nil
}

// Xor performs a bitwise XOR operation between the integer values of the
// receiver and the provided Value. If either Value is not of type Integer,
// it returns a typeError. The result of the XOR operation is stored in the
// receiver's integer field.
//
// Parameters:
//
//	v1 - The Value to XOR with the receiver.
//
// Returns:
//
//	An error if either Value is not of type Integer, otherwise nil.
func (v *Value) Xor(v1 *Value) error {
	if v.t != Integer || v1.t != Integer {
		return typeError("Xor", "")
	}
	v.i ^= v1.i
	return nil
}

// Or performs a bitwise OR operation between the current Value and another Value (v1).
// It returns an error if either of the Values is not of type Integer.
// The result of the operation is stored in the current Value.
//
// Parameters:
// - v1: A pointer to another Value to perform the OR operation with.
//
// Returns:
// - error: An error if the types of the Values are not Integer, otherwise nil.
func (v *Value) Or(v1 *Value) error {
	if v.t != Integer || v1.t != Integer {
		return typeError("Or", "")
	}
	v.i |= v1.i
	return nil
}

// LogAnd performs a logical AND operation between the current Value (v) and another Value (v1).
// Both values must be of type Integer or Floating. If either value is not of the correct type,
// a typeError is returned. The result of the logical AND operation is stored in the current Value (v)
// as an Integer (1 for true, 0 for false).
//
// Parameters:
//
//	v1 - The Value to perform the logical AND operation with.
//
// Returns:
//
//	error - Returns a typeError if either Value is not of type Integer or Floating.
func (v *Value) LogAnd(v1 *Value) error {
	if v.t != Integer && v.t != Floating || v1.t != Integer && v1.t != Floating {
		return typeError("LogAnd", "")
	}
	if (v.t == Floating && v.f != 0.0 || v.i != 0) &&
		(v1.t == Floating && v1.f != 0.0 || v1.i != 0) {
		v.i = 1
	} else {
		v.i = 0
	}
	v.t = Integer
	v.f = 0
	return nil
}

// LogOr performs a logical OR operation between the current Value (v) and another Value (v1).
// It supports both Integer and Floating types. If either value is non-zero, the result is 1 (true),
// otherwise, the result is 0 (false). The result is stored in the current Value (v) as an Integer type.
//
// Parameters:
//
//	v1 - The Value to perform the logical OR operation with.
//
// Returns:
//
//	error - Returns a typeError if either Value is not of type Integer or Floating.
func (v *Value) LogOr(v1 *Value) error {
	if v.t != Integer && v.t != Floating || v1.t != Integer && v1.t != Floating {
		return typeError("LogOr", "")
	}
	if (v.t == Floating && v.f != 0.0 || v.i != 0) ||
		(v1.t == Floating && v1.f != 0.0 || v1.i != 0) {
		v.i = 1
	} else {
		v.i = 0
	}
	v.t = Integer
	v.f = 0
	return nil
}
