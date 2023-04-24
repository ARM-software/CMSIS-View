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

func (v *Value) Compose(t Token, i int64, f float64, s string) {
	*v = Value{t, i, f, s, nil, nil}
}

func (v *Value) getValue() (Value, error) {
	if v.v == nil {
		return *v, typeError("not a variable", "")
	}
	return v.v.getValue()
}

func (v *Value) setValue(v1 *Value) error {
	if v.v == nil {
		return typeError("not a variable", "")
	}
	err := v.v.setValue(v1) // do not change v yet
	return err
}

func (v *Value) addList(v1 Value) error {
	if v.t == Nix {
		v.t = List
	} else if v.t != List {
		return typeError("not a list", "")
	}
	v.l = append(v.l, v1)
	return nil
}

func (v *Value) GetInt() int64 {
	switch v.t {
	case Integer:
		return v.i
	case Floating:
		return int64(v.f)
	}
	return 0
}

func (v *Value) GetUInt() uint64 {
	switch v.t {
	case Integer:
		return uint64(v.i)
	case Floating:
		return uint64(v.f)
	}
	return 0
}

func (v *Value) GetFloat() float64 {
	switch v.t {
	case Integer:
		return float64(v.i)
	case Floating:
		return v.f
	}
	return 0.0
}

func (v *Value) GetList() []Value {
	if v.IsList() {
		return v.l
	}
	return nil
}

func (v *Value) IsInteger() bool {
	return v.t == Integer
}

func (v *Value) IsFloating() bool {
	return v.t == Floating
}

func (v *Value) IsString() bool {
	return v.t == String
}

func (v *Value) IsIdentifier() bool {
	return v.t == Identifier
}

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

func (v *Value) Plus() error {
	switch v.t {
	case Integer:
	case Floating:
	default:
		return typeError("Plus", "")
	}
	return nil
}

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

func (v *Value) Compl() error {
	switch v.t {
	case Integer:
		v.i = -1 - v.i
	default:
		return typeError("Compl", "")
	}
	return nil
}

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

func (v *Value) Shl(v1 *Value) error {
	if v.t != Integer || v1.t != Integer {
		return typeError("shl", "")
	}
	v.i <<= v1.i
	return nil
}

func (v *Value) Shr(v1 *Value) error {
	if v.t != Integer || v1.t != Integer {
		return typeError("shr", "")
	}
	v.i >>= v1.i
	return nil
}

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

func (v *Value) And(v1 *Value) error {
	if v.t != Integer || v1.t != Integer {
		return typeError("And", "")
	}
	v.i &= v1.i
	return nil
}

func (v *Value) Xor(v1 *Value) error {
	if v.t != Integer || v1.t != Integer {
		return typeError("Xor", "")
	}
	v.i ^= v1.i
	return nil
}

func (v *Value) Or(v1 *Value) error {
	if v.t != Integer || v1.t != Integer {
		return typeError("Or", "")
	}
	v.i |= v1.i
	return nil
}

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
