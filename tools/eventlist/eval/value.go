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

package eval

import (
	"eventlist/elf"
)

type Value struct {
	t Token
	i int64
	f float64
	s string
	v *Variable
	l []Value
}

var TResult = [...][11]Token {
//        none,   I8,   U8,  I16,  U16,  I32,  U32,  I64,  U64,  F32,  F64
//------------------------------------------------------------------------
/*none*/ { Nix,  Nix,  Nix,  Nix,  Nix,  Nix,  Nix,  Nix,  Nix,  Nix,  Nix },
/* i8 */ { Nix,   I8,  I32,  U16,  U16,  I32,  U32,  I64,  U64,  F32,  F64 },
/* u8 */ { Nix,  I32,   U8,  I16,  U16,  I32,  U32,  I64,  U64,  F32,  F64 },
/* i16*/ { Nix,  I32,  I16,  I16,  U16,  I32,  U32,  I64,  U64,  F32,  F64 },
/* u16*/ { Nix,  I32,  U32,  U32,  U16,  I32,  U32,  I64,  U64,  F32,  F64 },
/* i32*/ { Nix,  I32,  I32,  I32,  I32,  I32,  U32,  I64,  U64,  F32,  F64 },
/* u32*/ { Nix,  U32,  U32,  U32,  U32,  U32,  U32,  I64,  U64,  F32,  F64 },
/* i64*/ { Nix,  I64,  U64,  I64,  I64,  I64,  I64,  I64,  U64,  F64,  F64 },
/* u64*/ { Nix,  U64,  U64,  U64,  U64,  U64,  U64,  U64,  U64,  F64,  F64 },
/* f32*/ { Nix,  F32,  F32,  F32,  F32,  F32,  F32,  F64,  F64,  F32,  F64 },
/* f64*/ { Nix,  F64,  F64,  F64,  F64,  F64,  F64,  F64,  F64,  F64,  F64 },
}

func calcResult(t Token, t1 Token) (Token, error) {
	if int(t) >= len(TResult[0]) || t == Nix || int(t1) >= len(TResult[0]) || t1 == Nix {
		return Nix, typeError("calcResult", "")
	}
	return TResult[t][t1], nil
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

func (v *Value) GetInt64() int64 {
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		return v.i
	case F32, F64:
		return int64(v.f)
	}
	return 0
}

func (v *Value) GetUInt64() uint64 {
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		return uint64(v.i)
	case F32, F64:
		return uint64(v.f)
	}
	return 0
}

func (v *Value) GetFloat64() float64 {
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		return float64(v.i)
	case F32, F64:
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
	return v.t >= I8 && v.t <= U64
}

func (v *Value) IsFloating() bool {
	return v.t == F64 || v.t == F32
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
	"__CalcMemUsed":   {CALCMEMUSED, 4, I64, I32},
	"__GetRegVal":     {GETREGVAL, 1, String, I32},
	"__Symbol_exists": {SYMBOLEXIST, 1, String, U8},
	"__FindSymbol":    {FINDSYMBOL, 1, String, U8},
	"__Offset_of":     {OFFSETOF, 1, String, I64},
	"__size_of":       {SIZEOF, 1, String, I64},
}

func (v *Value) Function(v1 *Value) error {
	const fnFunction = "Function"

	if v1 == nil {
		return typeError(fnFunction, "")
	}
	if !v.IsIdentifier() {
		return typeError(fnFunction, "")
	}
	if !v1.IsList() {
		return typeError(fnFunction, "")
	}
	var f Function
	var found bool
	if f, found = fctMap[v.s]; !found {
		return typeError(fnFunction, "")
	}
	if f.params != len(v1.GetList()) {
		return typeError(fnFunction, "")
	}
	for _, par := range v1.GetList() {
		if f.parType == String && par.t != String {
			return typeError(fnFunction, "")
		}
		if f.parType >= I8 && f.parType <= U64 && !(par.t >= I8 && par.t <= U64) {
			return typeError(fnFunction, "")
		}
		if f.parType >= F32 && f.parType <= F64 && !(par.t >= F32 && par.t <=F64) {
			return typeError(fnFunction, "")
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
	case I64, U64, I32, U32, I16, U16, I8, U8:
		v.i++
	case F32, F64:
		v.f++
	default:
		return typeError("Inc", "")
	}
	return nil
}

func (v *Value) Dec() error {
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		v.i--
	case F32, F64:
		v.f--
	default:
		return typeError("Dec", "")
	}
	return nil
}

func (v *Value) Plus() error {
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
	case F32, F64:
	default:
		return typeError("Plus", "")
	}
	return nil
}

func (v *Value) Neg() error {
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		v.i = -v.i
	case F32, F64:
		v.f = -v.f
	default:
		return typeError("Neg", "")
	}
	return nil
}

func (v *Value) Compl() error {
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		v.i = -1 - v.i
	default:
		return typeError("Compl", "")
	}
	return nil
}

func (v *Value) Not() error {
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		if v.i == 0 {
			v.i = 1
		} else {
			v.i = 0
		}
		v.t = U8
	default:
		return typeError("Compl", "")
	}
	return nil
}

func (v *Value) Cast(t Token) error {
	const fnCast = "Cast"

	switch t {
	case U8:
		switch v.t {
		case I64, U64, I32, U32, I16, U16, I8, U8:
			v.i = int64(uint8(v.i))
			v.t = U8
		case F32:
			v.i = int64(uint8(float32(v.f)))
			v.t = U8
			v.f = 0
		case F64:
			v.i = int64(uint8(v.f))
			v.t = U8
			v.f = 0
		default:
			return typeError(fnCast, "")
		}
	case I8:
		switch v.t {
		case I64, U64, I32, U32, I16, U16, I8, U8:
			v.i = int64(int8(v.i))
			v.t = I8
		case F32:
			v.i = int64(int8(float32(v.f)))
			v.t = I8
			v.f = 0
		case F64:
			v.i = int64(int8(v.f))
			v.t = I8
			v.f = 0
		default:
			return typeError(fnCast, "")
		}
	case U16:
		switch v.t {
		case I64, U64, I32, U32, I16, U16, I8, U8:
			v.i = int64(uint16(v.i))
			v.t = U16
		case F32:
			v.i = int64(uint16(float32(v.f)))
			v.t = U16
			v.f = 0
		case F64:
			v.i = int64(uint16(v.f))
			v.t = U16
			v.f = 0
		default:
			return typeError(fnCast, "")
		}
	case I16:
		switch v.t {
		case I64, U64, I32, U32, I16, U16, I8, U8:
			v.i = int64(int16(v.i))
			v.t = I16
		case F32:
			v.i = int64(int16(float32(v.f)))
			v.t = I16
			v.f = 0
		case F64:
			v.i = int64(int16(v.f))
			v.t = I16
			v.f = 0
		default:
			return typeError(fnCast, "")
		}
	case U32:
		switch v.t {
		case I64, U64, I32, U32, I16, U16, I8, U8:
			v.i = int64(uint32(v.i))
			v.t = U32
		case F32:
			v.i = int64(uint32(float32(v.f)))
			v.t = U32
			v.f = 0
		case F64:
			v.i = int64(uint32(v.f))
			v.t = U32
			v.f = 0
		default:
			return typeError(fnCast, "")
		}
	case I32:
		switch v.t {
		case I64, U64, I32, U32, I16, U16, I8, U8:
			v.i = int64(int32(v.i))
			v.t = I32
		case F32:
			v.i = int64(int32(float32(v.f)))
			v.t = I32
			v.f = 0
		case F64:
			v.i = int64(int32(v.f))
			v.t = I32
			v.f = 0
		default:
			return typeError(fnCast, "")
		}
	case U64:
		switch v.t {
		case I64, U64, I32, U32, I16, U16, I8, U8:
			v.i = int64(uint64(v.i))
			v.t = U64
		case F32:
			v.i = int64(float32(v.f))
			v.t = U64
			v.f = 0
		case F64:
			v.i = int64(uint64(v.f))
			v.t = U64
			v.f = 0
		default:
			return typeError(fnCast, "")
		}
	case I64:
		switch v.t {
		case I64, U64, I32, U32, I16, U16, I8, U8:
			v.t = I64
		case F32:
			v.i = int64(float32(v.f))
			v.t = I64
			v.f = 0
		case F64:
			v.i = int64(v.f)
			v.t = I64
			v.f = 0
		default:
			return typeError(fnCast, "")
		}
	case F32:
		switch v.t {
		case I64:
			v.f = float64(float32(v.i))
			v.t = F32
			v.i = 0
		case U64:
			v.f = float64(float32(uint64(v.i)))
			v.t = F32
			v.i = 0
		case I32:
			v.f = float64(float32(int32(v.i)))
			v.t = F32
			v.i = 0
		case U32:
			v.f = float64(float32(uint32(v.i)))
			v.t = F32
			v.i = 0
		case I16:
			v.f = float64(float32(int16(v.i)))
			v.t = F32
			v.i = 0
		case U16:
			v.f = float64(float32(uint16(v.i)))
			v.t = F32
			v.i = 0
		case I8:
			v.f = float64(float32(int8(v.i)))
			v.t = F32
			v.i = 0
		case U8:
			v.f = float64(float32(uint8(v.i)))
			v.t = F32
			v.i = 0
		case F32:
			v.f = float64(float32(v.f))
			v.t = F32
		case F64:
			v.f = float64(float32(v.f))
			v.t = F32
		default:
			return typeError(fnCast, "")
		}
	case F64:
		switch v.t {
		case I64:
			v.f = float64(v.i)
			v.t = F64
			v.i = 0
		case U64:
			v.f = float64(uint64(v.i))
			v.t = F64
			v.i = 0
		case I32:
			v.f = float64(int32(v.i))
			v.t = F64
			v.i = 0
		case U32:
			v.f = float64(uint32(v.i))
			v.t = F64
			v.i = 0
		case I16:
			v.f = float64(int16(v.i))
			v.t = F64
			v.i = 0
		case U16:
			v.f = float64(uint16(v.i))
			v.t = F64
			v.i = 0
		case I8:
			v.f = float64(int8(v.i))
			v.t = F64
			v.i = 0
		case U8:
			v.f = float64(uint8(v.i))
			v.t = F64
		case F32:
			v.f = float64(float32(v.f))
			v.t = F64
		case F64:
		default:
			return typeError(fnCast, "")
		}
	default:
		return typeError(fnCast, "")
	}
	return nil
}

func (v *Value) Mul(v1 *Value) error {
	const fnMul = "Mul"
	var resultT	Token
	var err	error
	if resultT, err = calcResult(v.t, v1.t); err != nil {
		return typeError(fnMul, "")
	}
	if err = v.Cast(resultT); err != nil {
		return typeError(fnMul, "")
	}
	vy := v1
	if err = vy.Cast(resultT); err != nil {
		return typeError(fnMul, "")
	}
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		v.i *= vy.i
	case F32, F64:
		v.f *= vy.f
	default:
		return typeError(fnMul, "")
	}
	return nil
}

func (v *Value) Div(v1 *Value) error {
	const fnDiv = "Div"
	var resultT	Token
	var err	error
	if resultT, err = calcResult(v.t, v1.t); err != nil {
		return typeError(fnDiv, "")
	}
	if err = v.Cast(resultT); err != nil {
		return typeError(fnDiv, "")
	}
	vy := v1
	if err = vy.Cast(resultT); err != nil {
		return typeError(fnDiv, "")
	}
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		if v1.i == 0 {
			return typeError("division by 0", "")
		}
		v.i /= vy.i
	case F32, F64:
		if v1.f == 0.0 {
			return typeError("division by 0", "")
		}
		v.f /= vy.f
	default:
		return typeError(fnDiv, "")
	}
	return nil
}

func (v *Value) Mod(v1 *Value) error {
	const fnMod = "Mod"
	var resultT	Token
	var err	error
	if resultT, err = calcResult(v.t, v1.t); err != nil {
		return typeError(fnMod, "")
	}
	if err = v.Cast(resultT); err != nil {
		return typeError(fnMod, "")
	}
	vy := v1
	if err = vy.Cast(resultT); err != nil {
		return typeError(fnMod, "")
	}
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		if v1.i == 0 {
			return typeError("modular by 0", "")
		}
		v.i %= vy.i
	case F32, F64:
		return typeError("mod with floatings", "")
	default:
		return typeError(fnMod, "")
	}
	return nil
}

func (v *Value) Add(v1 *Value) error {
	const fnAdd = "Add"
	var resultT	Token
	var err	error
	if resultT, err = calcResult(v.t, v1.t); err != nil {
		return typeError(fnAdd, "")
	}
	if err = v.Cast(resultT); err != nil {
		return typeError(fnAdd, "")
	}
	vy := v1
	if err = vy.Cast(resultT); err != nil {
		return typeError(fnAdd, "")
	}
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		v.i += vy.i
	case F32, F64:
		v.f += vy.f
	default:
		return typeError(fnAdd, "")
	}
	return nil
}

func (v *Value) Sub(v1 *Value) error {
	const fnSub = "Sub"
	var resultT	Token
	var err	error
	if resultT, err = calcResult(v.t, v1.t); err != nil {
		return typeError(fnSub, "")
	}
	if err = v.Cast(resultT); err != nil {
		return typeError(fnSub, "")
	}
	vy := v1
	if err = vy.Cast(resultT); err != nil {
		return typeError(fnSub, "")
	}
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		v.i -= vy.i
	case F32, F64:
		v.f -= vy.f
	default:
		return typeError(fnSub, "")
	}
	return nil
}

func (v *Value) Shl(v1 *Value) error {
	if !v.IsInteger() || !v1.IsInteger() {
		return typeError("shl", "")
	}
	v.i <<= v1.i
	return nil
}

func (v *Value) Shr(v1 *Value) error {
	if !v.IsInteger() || !v1.IsInteger() {
		return typeError("shr", "")
	}
	v.i >>= v1.i
	return nil
}

func (v *Value) Less(v1 *Value) error {
	const fnLess = "Less"
	var resultT	Token
	var err	error
	if resultT, err = calcResult(v.t, v1.t); err != nil {
		return typeError(fnLess, "")
	}
	if err = v.Cast(resultT); err != nil {
		return typeError(fnLess, "")
	}
	vy := v1
	if err = vy.Cast(resultT); err != nil {
		return typeError(fnLess, "")
	}
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		if v.i < vy.i {
			v.i = 1
		} else {
			v.i = 0
		}
		v.t = U8
	case F32, F64:
		if v.f < vy.f {
			v.i = 1
		} else {
			v.i = 0
		}
		v.t = U8
		v.f = 0.0
	default:
		return typeError(fnLess, "")
	}
	return nil
}

func (v *Value) LessEqual(v1 *Value) error {
	const fnLessEqual = "LessEqual"
	var resultT	Token
	var err	error
	if resultT, err = calcResult(v.t, v1.t); err != nil {
		return typeError(fnLessEqual, "")
	}
	if err = v.Cast(resultT); err != nil {
		return typeError(fnLessEqual, "")
	}
	vy := v1
	if err = vy.Cast(resultT); err != nil {
		return typeError(fnLessEqual, "")
	}
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		if v.i <= vy.i {
			v.i = 1
		} else {
			v.i = 0
		}
		v.t = U8
	case F32, F64:
		if v.f <= vy.f {
			v.i = 1
		} else {
			v.i = 0
		}
		v.t = U8
		v.f = 0.0
	default:
		return typeError(fnLessEqual, "")
	}
	return nil
}

func (v *Value) Greater(v1 *Value) error {
	const fnGreater = "Greater"
	var resultT	Token
	var err	error
	if resultT, err = calcResult(v.t, v1.t); err != nil {
		return typeError(fnGreater, "")
	}
	if err = v.Cast(resultT); err != nil {
		return typeError(fnGreater, "")
	}
	vy := v1
	if err = vy.Cast(resultT); err != nil {
		return typeError(fnGreater, "")
	}
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		if v.i > vy.i {
			v.i = 1
		} else {
			v.i = 0
		}
		v.t = U8
	case F32, F64:
		if v.f > vy.f {
			v.i = 1
		} else {
			v.i = 0
		}
		v.t = U8
		v.f = 0.0
	default:
		return typeError(fnGreater, "")
	}
	return nil
}

func (v *Value) GreaterEqual(v1 *Value) error {
	const fnGreaterEqual = "GreaterEqual"
	var resultT	Token
	var err	error
	if resultT, err = calcResult(v.t, v1.t); err != nil {
		return typeError(fnGreaterEqual, "")
	}
	if err = v.Cast(resultT); err != nil {
		return typeError(fnGreaterEqual, "")
	}
	vy := v1
	if err = vy.Cast(resultT); err != nil {
		return typeError(fnGreaterEqual, "")
	}
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		if v.i >= vy.i {
			v.i = 1
		} else {
			v.i = 0
		}
		v.t = U8
	case F32, F64:
		if v.f >= vy.f {
			v.i = 1
		} else {
			v.i = 0
		}
		v.t = U8
		v.f = 0.0
	default:
		return typeError(fnGreaterEqual, "")
	}
	return nil
}

func (v *Value) Equal(v1 *Value) error {
	const fnEqual = "Equal"
	var resultT	Token
	var err	error
	if resultT, err = calcResult(v.t, v1.t); err != nil {
		return typeError(fnEqual, "")
	}
	if err = v.Cast(resultT); err != nil {
		return typeError(fnEqual, "")
	}
	vy := v1
	if err = vy.Cast(resultT); err != nil {
		return typeError(fnEqual, "")
	}
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		if v.i == vy.i {
			v.i = 1
		} else {
			v.i = 0
		}
		v.t = U8
	case F32, F64:
		if v.f == vy.f {
			v.i = 1
		} else {
			v.i = 0
		}
		v.t = U8
		v.f = 0.0
	default:
		return typeError(fnEqual, "")
	}
	return nil
}

func (v *Value) NotEqual(v1 *Value) error {
	const fnNotEqual = "NotEqual"
	var resultT	Token
	var err	error
	if resultT, err = calcResult(v.t, v1.t); err != nil {
		return typeError(fnNotEqual, "")
	}
	if err = v.Cast(resultT); err != nil {
		return typeError(fnNotEqual, "")
	}
	vy := v1
	if err = vy.Cast(resultT); err != nil {
		return typeError(fnNotEqual, "")
	}
	switch v.t {
	case I64, U64, I32, U32, I16, U16, I8, U8:
		if v.i != vy.i {
			v.i = 1
		} else {
			v.i = 0
		}
		v.t = U8
	case F32, F64:
		if v.f != vy.f {
			v.i = 1
		} else {
			v.i = 0
		}
		v.t = U8
		v.f = 0.0
	default:
		return typeError(fnNotEqual, "")
	}
	return nil
}

func (v *Value) And(v1 *Value) error {
	if !v.IsInteger() || !v1.IsInteger() {
		return typeError("And", "")
	}
	v.i &= v1.i
	if v1.t > v.t {
		v.t = v1.t
	}
	return nil
}

func (v *Value) Xor(v1 *Value) error {
	if !v.IsInteger() || !v1.IsInteger() {
		return typeError("Xor", "")
	}
	v.i ^= v1.i
	if v1.t > v.t {
		v.t = v1.t
	}
	return nil
}

func (v *Value) Or(v1 *Value) error {
	if !v.IsInteger() || !v1.IsInteger() {
		return typeError("Or", "")
	}
	v.i |= v1.i
	if v1.t > v.t {
		v.t = v1.t
	}
	return nil
}

func (v *Value) LogAnd(v1 *Value) error {
	if !v.IsInteger() && !v.IsFloating() || !v1.IsInteger() && !v1.IsFloating() {
		return typeError("LogAnd", "")
	}
	if ((v.t == F32 || v.t == F64) && v.f != 0.0 || v.i != 0) &&
		((v1.t == F32 || v1.t == F64) && v1.f != 0.0 || v1.i != 0) {
		v.i = 1
	} else {
		v.i = 0
	}
	v.t = U8
	v.f = 0
	return nil
}

func (v *Value) LogOr(v1 *Value) error {
	if !v.IsInteger() && !v.IsFloating() || !v1.IsInteger() && !v1.IsFloating() {
		return typeError("LogOr", "")
	}
	if ((v.t == F32 || v.t == F64) && v.f != 0.0 || v.i != 0) ||
		((v1.t == F32 || v1.t == F64) && v1.f != 0.0 || v1.i != 0) {
		v.i = 1
	} else {
		v.i = 0
	}
	v.t = U8
	v.f = 0
	return nil
}
