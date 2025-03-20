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

import (
	"errors"
	"math"
	"strings"
)

const maxUint64 = 1<<64 - 1

type Token int

const (
	Nix Token = iota
	Integer
	Floating
	String
	Identifier
	List

	Not
	Compl
	Assign
	OrAssign
	XorAssign
	AndAssign
	ShlAssign
	ShrAssign
	PlusAssign
	MinusAssign
	MulAssign
	DivAssign
	ModAssign
	Quest
	Colon
	LogOr
	LogAnd
	Or
	Xor
	And
	Equal
	NotEqual
	Less
	LessEqual
	Greater
	GreaterEqual
	Shl
	Shr
	Add
	Sub
	Mul
	Div
	Mod
	AddAdd
	SubSub
	ParenO
	ParenC
	BracketO
	BracketC
	Comma
	Semi
	Dot
	Pointer
)

var ITokens = map[string]Token{
	"!":   Not,
	"~":   Compl,
	"=":   Assign,
	"|=":  OrAssign,
	"^=":  XorAssign,
	"&=":  AndAssign,
	"<<=": ShlAssign,
	">>=": ShrAssign,
	"+=":  PlusAssign,
	"-=":  MinusAssign,
	"*=":  MulAssign,
	"/=":  DivAssign,
	"%=":  ModAssign,
	"?":   Quest,
	":":   Colon,
	"||":  LogOr,
	"&&":  LogAnd,
	"|":   Or,
	"^":   Xor,
	"&":   And,
	"==":  Equal,
	"!=":  NotEqual,
	"<":   Less,
	"<=":  LessEqual,
	">":   Greater,
	">=":  GreaterEqual,
	"<<":  Shl,
	">>":  Shr,
	"+":   Add,
	"-":   Sub,
	"*":   Mul,
	"/":   Div,
	"%":   Mod,
	"++":  AddAdd,
	"--":  SubSub,
	"(":   ParenO,
	")":   ParenC,
	"[":   BracketO,
	"]":   BracketC,
	",":   Comma,
	";":   Semi,
	".":   Dot,
	"->":  Pointer,
}

type Type int

const (
	NoType Type = iota
	Uint8
	Int8
	Uint16
	Int16
	Uint32
	Int32
	Uint64
	Int64
	Float
	Double
)

var ITypes = map[string]Type{
	"uint8_t":  Uint8,
	"int8_t":   Int8,
	"uint16_t": Uint16,
	"int16_t":  Int16,
	"uint32_t": Uint32,
	"int32_t":  Int32,
	"uint64_t": Uint64,
	"int64_t":  Int64,
	"float":    Float,
	"double":   Double,
}

type Expression struct {
	in       *string
	pos      int
	next     Value
	typedefs Typedefs
	tdUsed   map[string]string
}

var ErrRange = errors.New("value out of range")

var ErrSyntax = errors.New("syntax error")

var ErrType = errors.New("value type error")

var ErrEof = errors.New("eof") //nolint:golint,revive

type NumError struct {
	Func string // failing function
	Num  string // input value
	Err  error  // error reason
}

// Error returns a formatted error string for the NumError type.
// The returned string includes the function name, the number that caused the error,
// and the original error message.
func (e *NumError) Error() string {
	return "expression." + e.Func + ": " + "parsing \"" + e.Num + "\": " + e.Err.Error()
}

func (e *NumError) Unwrap() error { return e.Err }

func syntaxError(fn, str string) *NumError {
	return &NumError{fn, str, ErrSyntax}
}

func rangeError(fn, str string) *NumError {
	return &NumError{fn, str, ErrRange}
}

func typeError(fn, str string) *NumError {
	return &NumError{fn, str, ErrType}
}

// get retrieves the current byte from the input and advances the position.
// It returns the byte and an error if the end of the input is reached.
//
// Returns:
//   - byte: The current byte from the input.
//   - error: ErrEof if the end of the input is reached.
func (ex *Expression) get() (byte, error) {
	if ex.pos >= len(*ex.in) {
		return 0, ErrEof
	}
	c := (*ex.in)[ex.pos]
	ex.pos++
	return c, nil
}

// peek returns the next byte in the input without advancing the position.
// If the end of the input is reached, it returns an error ErrEof.
//
// Returns:
//
//	byte: The next byte in the input.
//	error: ErrEof if the end of the input is reached.
func (ex *Expression) peek() (byte, error) {
	if ex.pos >= len(*ex.in) {
		return 0, ErrEof
	}
	c := (*ex.in)[ex.pos]
	return c, nil
}

// back moves the position pointer one step back if it is not already at the beginning.
func (ex *Expression) back() {
	if ex.pos > 0 {
		ex.pos--
	}
}

// skipToEnd sets the position of the Expression to the end of the input.
// This effectively skips all remaining characters in the input.
func (ex *Expression) skipToEnd() {
	ex.pos = len(*ex.in)
}

// getPos returns the position of the expression.
// It retrieves the value of the pos field from the Expression struct.
// Used in conjunction with setPos to save and restore the position.
func (ex *Expression) getPos() int {
	return ex.pos
}

// setPos sets the position of the expression to the given value.
//
// Parameters:
//
//	pos - the new position to set for the expression.
func (ex *Expression) setPos(pos int) {
	ex.pos = pos
}

// lower converts an uppercase ASCII letter to its lowercase equivalent.
func lower(c byte) byte {
	return c | ('x' - 'X')
}

// parseUint parses an unsigned integer from the expression.
// It supports decimal, octal, and hexadecimal formats.
// The function returns the parsed uint64 value and an error if any occurred during parsing.
// If the input is empty or contains invalid characters, appropriate errors are returned.
func (ex *Expression) parseUint() (uint64, error) {
	const fnParseUint = "parseUint"

	c, err := ex.get()
	if err != nil {
		return 0, err // eof
	}

	s0 := string(c)
	base := 10
	if c == '0' { // Look for octal, hex prefix.
		if c, err = ex.get(); err != nil {
			// err could only be an EOF which is corrrect here
			// only a '0'
			return 0, nil //nolint:golint,nilerr
		}
		if lower(c) == 'x' {
			s0 += string(c)
			base = 16
			if c, err = ex.get(); err != nil {
				return 0, syntaxError(fnParseUint, "")
			}
			s0 += string(c)
		} else {
			ex.back()
			base = 8
			c = '0'
		}
	}

	// Cutoff is the smallest number such that cutoff*base > maxUint64.
	cutoff := maxUint64/uint64(base) + 1

	first := true
	var n uint64
loop:
	for {
		var d byte
		switch {
		case base == 8 && '0' <= c && c <= '7':
			d = c - '0'
		case base >= 10 && '0' <= c && c <= '9':
			d = c - '0'
		case base == 16 && 'a' <= lower(c) && lower(c) <= 'f':
			d = lower(c) - 'a' + 10
		default:
			ex.back() // back to breaking char
			break loop
		}

		if n >= cutoff { // n*base overflows
			return maxUint64, rangeError(fnParseUint, s0)
		}
		n *= uint64(base)

		n1 := n + uint64(d)
		if n1 < n || n1 > maxUint64 { // n+d overflows
			return maxUint64, rangeError(fnParseUint, s0) // cannot happen because of n*base test
		}
		n = n1
		if c, err = ex.get(); err != nil {
			break loop // end of number
		}
		first = false
		s0 += string(c)
	}

	if first && base == 16 {
		return 0, syntaxError(fnParseUint, s0)
	}
	return n, nil
}

// parseFloat parses a floating-point number from the expression.
// It handles optional leading signs, decimal points, and exponents.
// It returns the parsed float64 value and an error if the parsing fails.
// The function assumes that the input is well-formed and does not perform
// extensive validation beyond basic syntax checks.
//
// The function supports numbers with up to 19 significant digits and
// handles both positive and negative exponents. It uses a lookup table
// for efficient exponentiation.
//
// Returns:
// - float64: The parsed floating-point number.
// - error: An error if the parsing fails, otherwise nil.
func (ex *Expression) parseFloat() (float64, error) {
	var pow10Tab = []float64{
		1e1, 1e2, 1e4, 1e8, 1e16, 1e32, 1e64, 1e128, 1e256, 0,
	}
	const fnParseFloat = "parseFloat"

	var c byte
	var err error
	if c, err = ex.get(); err != nil {
		return 0, syntaxError(fnParseFloat, "")
	}
	neg := false
	switch c {
	case '-':
		neg = true
		fallthrough
	case '+':
		if c, err = ex.get(); err != nil {
			return 0, syntaxError(fnParseFloat, "")
		}
	}

	// digits
	var mantissa uint64
	const maxMantDigits = 19 // 10^19 fits in uint64
	sawdot := false
	sawdigits := false
	nd := 0
	ndMant := 0
	dp := 0

	for {
		if c == '.' {
			if sawdot {
				ex.back()
				break
			}
			sawdot = true
			dp = nd
			if c, err = ex.get(); err != nil {
				break // '.' at eof
			}
		} else if '0' <= c && c <= '9' {
			sawdigits = true
			if c == '0' && nd == 0 { // ignore leading zeros
				dp--
			} else {
				nd++
				if ndMant < maxMantDigits {
					mantissa *= 10
					mantissa += uint64(c - '0')
					ndMant++
				}
			}
			if c, err = ex.get(); err != nil {
				break // digit at eof
			}
		} else {
			break // something behind the number
		}
	}
	for mantissa > 1<<52 {
		mantissa /= 10
		ndMant--
	}
	if !sawdigits {
		return 0, rangeError("no digits", "")
	}
	if !sawdot {
		dp = nd
	}

	// optional exponent moves decimal point.
	// if we read a very large, very long number,
	// just be sure to move the decimal point by
	// a lot (say, 100000).  it doesn't matter if it's
	// not the exact number.
	if err == nil && lower(c) == 'e' {
		if c, err = ex.get(); err != nil {
			ex.back()     // back onto 'e'
			return 0, err // e without digits
		}
		esign := 1
		switch c {
		case '-':
			esign = -1
			fallthrough
		case '+':
			if c, err = ex.get(); err != nil {
				return 0, syntaxError(fnParseFloat, "") // e+- without digits
			}
		}
		if c < '0' || c > '9' {
			ex.back()
			return 0, syntaxError(fnParseFloat, "") // e without digits
		}
		e := 0
		for c >= '0' && c <= '9' {
			if e < 10000 {
				e = e*10 + int(c-'0')
			}
			if c, err = ex.get(); err != nil {
				break // end of exponent
			}
		}
		dp += e * esign
	}
	if !errors.Is(err, ErrEof) {
		ex.back()
	}

	exp := 0
	if mantissa != 0 {
		exp = dp - ndMant
	}

	var f float64
	f = float64(mantissa)
	if neg {
		f = -f
	}
	if exp == 0 { // an integer.
		return f, nil
	}
	i := 0
	for exp > 0 {
		if (exp & 1) != 0 {
			f *= pow10Tab[i]
		}
		exp >>= 1
		i++
	}
	exp = -exp
	for exp > 0 {
		if (exp & 1) != 0 {
			f /= pow10Tab[i]
		}
		exp >>= 1
		i++
	}

	return f, nil
}

// hex parses a hexadecimal number from the expression and returns the
// corresponding byte value and the string representation of the parsed
// hexadecimal number. It reads characters from the expression until it
// encounters a non-hexadecimal character, at which point it stops and
// returns the results. If a non-hexadecimal character is encountered,
// it is pushed back to the expression for further processing.
func (ex *Expression) hex() (cx byte, s string) {
	for {
		var c byte
		var err error
		if c, err = ex.get(); err != nil {
			break
		}
		if '0' <= c && c <= '9' {
			s += string(c)
			cx = cx<<4 | (c - '0')
		} else if 'A' <= c && c <= 'F' || 'a' <= c && c <= 'f' {
			s += string(c)
			cx = cx<<4 | (lower(c) - 'a' + 0xa)
		} else {
			ex.back()
			break
		}
	}
	return cx, s
}

// hex4 reads the next four characters from the expression, interprets them as a hexadecimal number,
// and returns the corresponding uint16 value along with the string representation of the characters.
// If any character is not a valid hexadecimal digit, it returns a rangeError.
//
// Returns:
//
//	cx (uint16): The uint16 value of the hexadecimal number.
//	s (string): The string representation of the four characters read.
//	err (error): An error if any character is not a valid hexadecimal digit, otherwise nil.
func (ex *Expression) hex4() (cx uint16, s string, err error) {
	for i := 0; i < 4; i++ {
		var c byte
		if c, err = ex.get(); err != nil {
			return
		}
		s += string(c)
		if '0' <= c && c <= '9' {
			cx = cx<<4 | uint16(c-'0')
		} else if 'A' <= c && c <= 'F' || 'a' <= c && c <= 'f' {
			cx = cx<<4 | uint16(lower(c)-'a'+0xa)
		} else {
			return cx, s, rangeError("hex4", s)
		}
	}
	return cx, s, nil
}

// hex8 reads the next 8 characters from the expression, interprets them as a hexadecimal number,
// and returns the corresponding uint32 value along with the string representation of the characters.
// If any character is not a valid hexadecimal digit, it returns a rangeError.
//
// Returns:
//
//	cx (uint32): The uint32 value of the hexadecimal number.
//	s (string): The string representation of the four characters read.
//	err (error): An error if any character is not a valid hexadecimal digit, otherwise nil.
func (ex *Expression) hex8() (cx uint32, s string, err error) {
	for i := 0; i < 8; i++ {
		var c byte
		if c, err = ex.get(); err != nil {
			return
		}
		s += string(c)
		if '0' <= c && c <= '9' {
			cx = cx<<4 | uint32(c-'0')
		} else if 'A' <= c && c <= 'F' || 'a' <= c && c <= 'f' {
			cx = cx<<4 | uint32(lower(c)-'a'+0xa)
		} else {
			return cx, s, rangeError("hex8", s)
		}
	}
	return cx, s, nil
}

// lex performs lexical analysis on the input expression.
// It reads characters from the input, identifies tokens, and returns a Value representing the token
// and an error if any issues are encountered during the process.
//
// The method handles various types of tokens including:
// - Whitespace characters (spaces, tabs, form feeds)
// - Digits (0-9) which can be parsed as integers or floating-point numbers
// - Identifiers (a-z, A-Z, and underscores)
// - String literals enclosed in double quotes (")
// - Character literals enclosed in single quotes (')
// - Special tokens and comments
//
// The method uses helper functions such as get, back, peek, parseUint, parseFloat, hex, hex4, hex8,
// and skipToEnd to perform specific tasks during the lexical analysis.
//
// Returns:
// - Value: A struct representing the identified token.
// - error: An error if any issues are encountered during the lexical analysis.
func (ex *Expression) lex() (Value, error) {
	const fnLex = "Lex"

	var v Value
	var c byte
	var err error

	for {
		c, err = ex.get()
		if err != nil {
			return v, err
		}
		if c != ' ' && c != '\t' && c != '\f' {
			break
		}
	}

	s0 := string(c)
	if '0' <= c && c <= '9' { // a digit
		ex.back()
		begin := ex.getPos()
		ui, err := ex.parseUint()
		if err != nil {
			return v, err
		}
		c, err = ex.peek()
		if err == nil && (c == '.' || lower(c) == 'e') {
			ex.setPos(begin)
			f, err := ex.parseFloat()
			if err != nil {
				return v, err
			}
			return Value{t: Floating, f: f}, nil
		}
		return Value{t: Integer, i: int64(ui)}, nil

	} else if 'a' <= lower(c) && lower(c) <= 'z' {
	loop:
		for {
			if c, err = ex.get(); err != nil {
				break
			}
			switch {
			case '0' <= c && c <= '9' || 'a' <= lower(c) && lower(c) <= 'z' || c == '_':
				s0 += string(c)
			default:
				ex.back()
				break loop
			}
		}
		if strings.ToLower(s0) == "inf" {
			return Value{t: Floating, f: math.Inf(0)}, nil
		} else if strings.ToLower(s0) == "nan" {
			return Value{t: Floating, f: math.NaN()}, nil
		}
		return Value{t: Identifier, s: s0}, nil

	} else if c == '"' {
		for {
			if c, err = ex.get(); err != nil {
				break
			}
			s0 += string(c)
			done := false
			if c == '\\' {
				var cx byte
				if c, err = ex.get(); err != nil {
					return v, syntaxError(fnLex, s0)
				}
				s0 += string(c)
				switch c {
				case '\'':
					c = '\''
				case '"':
					c = '"'
				case 'a':
					c = '\a'
				case 'b':
					c = '\b'
				case 'e':
					c = '\x1b' // GCC extension
				case 'f':
					c = '\f'
				case 'n':
					c = '\n'
				case 'r':
					c = '\r'
				case 't':
					c = '\t'
				case 'v':
					c = '\v'
				case '0', '1', '2', '3', '4', '5', '6', '7':
					cx = c - '0'
					if c, err = ex.get(); err != nil {
						return v, syntaxError(fnLex, s0)
					}
					if c >= '0' && c <= '7' {
						s0 += string(c)
						cx = cx<<3 | (c - '0')
						if c, err = ex.get(); err != nil {
							return v, syntaxError(fnLex, s0)
						}
						if c >= '0' && c <= '7' {
							s0 += string(c)
							cx = cx<<3 | (c - '0')
						} else {
							ex.back()
						}
					} else {
						ex.back()
					}
					c = cx
				case 'x':
					var s string
					c, s = ex.hex()
					s0 += s
				case 'u':
					var s string
					var i uint16
					if i, s, err = ex.hex4(); err != nil {
						return v, syntaxError(fnLex, s0)
					}
					s0 += s
					v.s += string(rune(i))
					done = true
				case 'U':
					var s string
					var i uint32
					if i, s, err = ex.hex8(); err != nil {
						return v, syntaxError(fnLex, s0)
					}
					s0 += s
					v.s += string(rune(i))
					done = true
				}
			} else if c == '"' {
				v.t = String
				return v, nil
			}
			if !done {
				v.s += string(c)
			}
		}
	} else if c == '\'' {
		if c, err = ex.get(); err != nil {
			return v, syntaxError(fnLex, s0)
		}
		s0 += string(c)
		done := false
		if c == '\\' {
			var cx byte
			if c, err = ex.get(); err != nil {
				return v, syntaxError(fnLex, s0)
			}
			s0 += string(c)
			switch c {
			case '\'':
				c = '\''
			case '"':
				c = '"'
			case 'a':
				c = '\a'
			case 'b':
				c = '\b'
			case 'e':
				c = '\x1b' // GCC extension
			case 'f':
				c = '\f'
			case 'n':
				c = '\n'
			case 'r':
				c = '\r'
			case 't':
				c = '\t'
			case 'v':
				c = '\v'
			case '0', '1', '2', '3', '4', '5', '6', '7':
				cx = c - '0'
				if c, err = ex.get(); err != nil {
					return v, syntaxError(fnLex, s0)
				}
				if c >= '0' && c <= '7' {
					s0 += string(c)
					cx = cx<<3 | (c - '0')
					if c, err = ex.get(); err != nil {
						return v, syntaxError(fnLex, s0)
					}
					if c >= '0' && c <= '7' {
						s0 += string(c)
						cx = cx<<3 | (c - '0')
					} else {
						ex.back()
					}
				} else {
					ex.back()
				}
				c = cx
			case 'x':
				var s string
				c, s = ex.hex()
				s0 += s
			case 'u':
				var s string
				var i uint16
				if i, s, err = ex.hex4(); err != nil {
					return v, syntaxError(fnLex, s0)
				}
				s0 += s
				v.i = int64(i)
				done = true
			case 'U':
				var s string
				var i uint32
				if i, s, err = ex.hex8(); err != nil {
					return v, syntaxError(fnLex, s0)
				}
				s0 += s
				v.i = int64(i)
				done = true
			}
		}
		if !done {
			v.i = int64(c)
		}
		if c, err = ex.get(); err != nil {
			return Value{}, syntaxError(fnLex, s0)
		}
		if c == '\'' {
			v.t = Integer
			return v, nil
		}
	} else {
		var t Token
		lastc := c
		c, err = ex.get() // 2nd char
		if !errors.Is(err, ErrEof) {
			if lastc == '/' && c == '/' { // comment till end
				ex.skipToEnd()
				return v, ErrEof
			}
			s0 += string(c)
			c, err = ex.get() // 3rd char
			if !errors.Is(err, ErrEof) {
				s0 += string(c)
				t = ITokens[s0] // try 3 chars
				if t != Nix {
					return Value{t: t}, nil
				}
				s0 = s0[:len(s0)-1]
				ex.back()
			}
			t = ITokens[s0] // try 2 chars
			if t != Nix {
				return Value{t: t}, nil
			}
			s0 = s0[:len(s0)-1]
			ex.back()
		}
		t = ITokens[s0] // try 1 char
		if t != Nix {
			return Value{t: t}, nil
		}
	}
	return Value{}, syntaxError(fnLex, s0)
}

// Integer
// Identifier
// TypedefIdentifier : Identifier
// TypedefIdentifier : Identifier : Identifier
// String
// ( expression )
// primary parses the primary expressions in the given context.
// It handles different types of tokens such as Integer, Floating, Identifier, String, and Parentheses.
// For Identifiers, it also handles typedefs and their members, including enums.
// It returns the parsed Value and an error if any issues are encountered during parsing.
func (ex *Expression) primary() (Value, error) {
	var v Value
	var err error

	v = ex.next
	switch ex.next.t {
	case Integer:
		if ex.next, err = ex.lex(); err != nil {
			return v, err
		}
	case Floating:
		if ex.next, err = ex.lex(); err != nil {
			return v, err
		}
	case Identifier:
		mu.Lock()
		v.v = names[v.s]
		mu.Unlock()
		if ex.next, err = ex.lex(); err != nil {
			return v, err
		}
		itypedef, ok := ex.typedefs[v.s]
		if ok && ex.next.t == Colon {
			if ex.next, err = ex.lex(); err != nil {
				return v, syntaxError("member name expected", "")
			}
			v1 := ex.next // typedef_name:member
			if ex.next.t != Identifier {
				return v1, syntaxError("member name expected", "")
			}
			if ex.next, err = ex.lex(); err != nil {
				return v1, err
			}
			member, ok := itypedef.Members[v1.s]
			if !ok {
				return v1, syntaxError(v1.s+" unknown in "+v.s, "")
			}
			if ok && ex.next.t == Colon {
				if ex.next, err = ex.lex(); err != nil {
					return v1, syntaxError("enum name expected", "")
				}
				v2 := ex.next
				if ex.next.t != Identifier {
					return v2, syntaxError("enum name expected", "")
				}
				for n, s := range member.Enums {
					if s == v2.s {
						v = Value{t: Integer, i: n}
						if ex.next, err = ex.lex(); err != nil {
							return v, err
						}
						return v, nil
					}
				}
				return v2, syntaxError("enum "+v2.s+"unknown", "")
			}
			return Eval(&member.Offset, ex.typedefs, ex.tdUsed)
			// return Value{t: Integer, i: int64(member.Offset)}, nil
		}
	case String:
		if ex.next, err = ex.lex(); err != nil {
			return v, err
		}
	case ParenO:
		if ex.next, err = ex.lex(); err != nil {
			return ex.next, err
		}
		if v, err = ex.expression(); err != nil {
			return v, err
		}
		if ex.next.t != ParenC {
			return v, syntaxError("expected \")\"", "")
		}
		if ex.next, err = ex.lex(); err != nil {
			return v, err
		}
	default:
		return ex.next, syntaxError("primary", "")
	}
	return v, nil
}

// asnExpr
// arguments , asnExpr
// arguments parses and returns a list of argument expressions.
// It starts by checking if the next token is Nix, in which case it returns an empty value.
// Otherwise, it parses the first argument expression and adds it to the list.
// It then enters a loop to parse additional arguments separated by commas.
// If a syntax error or end-of-file error occurs, it returns the error.
// The function returns the list of parsed argument expressions or an error if any occurred.
func (ex *Expression) arguments() (Value, error) {
	var left Value
	var right Value
	var err error

	if ex.next.t == Nix {
		return left, nil
	}
	if left, err = ex.asnExpr(); err != nil {
		return left, err
	}
	if err = right.addList(left); err != nil {
		return left, err // cannot happen because right is Nix
	}
	left = right
loop:
	for {
		switch ex.next.t {
		case Comma:
			if ex.next, err = ex.lex(); err != nil {
				return left, syntaxError("expected expression", "")
			}
			if right, err = ex.asnExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if err = left.addList(right); err != nil {
				return left, err // cannot happen because left already is a list
			}
		default:
			break loop
		}
	}
	return left, nil
}

// primary
// primary ++
// primary --
// primary . identifier
// primary -> identifier
// primary ( )
// primary ( arguments )
// primary [ asnExpr ]
// postfix evaluates postfix expressions for the given Expression.
// It handles various postfix operators such as increment (AddAdd), decrement (SubSub),
// member access (Dot), pointer access (Pointer), function calls (ParenO), and array indexing (BracketO).
// The function returns the evaluated Value and an error if any occurred during the evaluation.
//
// Supported postfix operators:
// - AddAdd: Post-increment operator (e.g., x++)
// - SubSub: Post-decrement operator (e.g., x--)
// - Dot: Member access operator (e.g., obj.field)
// - Pointer: Pointer access operator (e.g., ptr->field)
// - ParenO: Function call operator (e.g., func())
// - BracketO: Array indexing operator (e.g., arr[index])
//
// Returns:
// - Value: The evaluated value after applying the postfix operator.
// - error: An error if any occurred during the evaluation.
func (ex *Expression) postfix() (Value, error) { // TODO: not finished yet
	var left Value
	var right Value
	var v Value
	var err error

	if left, err = ex.primary(); err != nil {
		return left, err
	}
	switch ex.next.t {
	case AddAdd:
		if !left.IsIdentifier() {
			return left, syntaxError("identifier expected", "")
		}
		if v, err = left.getValue(); err != nil {
			return left, err
		}
		if err = v.Inc(); err != nil {
			return v, err
		}
		if err = left.setValue(&v); err != nil { // do not change left, it is postincrement
			return left, err // cannot happen because of working getValue
		}
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
	case SubSub:
		if !left.IsIdentifier() {
			return left, syntaxError("identifier expected", "")
		}
		if v, err = left.getValue(); err != nil {
			return left, err
		}
		if err = v.Dec(); err != nil {
			return v, err
		}
		if err = left.setValue(&v); err != nil { // do not change left, it is postdecrement
			return left, err // cannot happen because of working getValue
		}
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
	case Dot:
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
		if !left.IsIdentifier() {
			return left, syntaxError("identifier expected", "")
		}
		if !ex.next.IsIdentifier() {
			return ex.next, syntaxError("identifier expected", "")
		}
		members, ok := ex.typedefs[ex.tdUsed[left.s]]
		if ok { // it is a val* typedef
			member, ok := members.Members[ex.next.s]
			if ok { // field found
				if right, err = Eval(&member.Offset, ex.typedefs, ex.tdUsed); err != nil {
					return ex.next, err
				}
				if v, err = left.getValue(); err != nil {
					return left, err
				}
				if !right.IsInteger() {
					return right, syntaxError("integer offset expected", "")
				}
				if err = v.Extract(members.Size, members.BigEndian, uint32(right.i)); err != nil {
					return v, err
				}
				if err = v.Cast(member.IType); err != nil {
					return v, err
				}
				return v, nil
			}
		}
		if ex.next, err = ex.lex(); err != nil {
			return ex.next, err
		}
	case Pointer:
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
		if !left.IsIdentifier() {
			return left, syntaxError("identifier expected", "")
		}
		if !ex.next.IsIdentifier() {
			return ex.next, syntaxError("identifier expected", "")
		} // TODO: noch nicht implementiert
		if ex.next, err = ex.lex(); err != nil {
			return ex.next, err
		}
	case ParenO:
		if ex.next, err = ex.lex(); err != nil {
			return left, syntaxError("expected \")\"", "")
		}
		if ex.next.t != ParenC {
			if right, err = ex.arguments(); err != nil {
				return left, err
			}
			if ex.next.t != ParenC {
				return left, syntaxError("expected \")\"", "")
			}
			if err = left.Function(&right); err != nil {
				return left, err
			}
		}
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
	case BracketO:
		if ex.next, err = ex.lex(); err != nil {
			return left, syntaxError("expected expression", "")
		}
		if right, err = ex.asnExpr(); err != nil {
			return left, err
		}
		if ex.next.t != BracketC {
			return left, syntaxError("expected \"]\"", "")
		}
		v = right
		if v.IsIdentifier() {
			if v, err = v.getValue(); err != nil {
				return left, err
			}
		}
		left.i = v.GetInt() // TODO: noch nicht implementiert
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
	}
	return left, nil
}

// + postfix
// - postfix
// ~ postfix
// ! postfix
// postfix
// unary parses and evaluates a unary expression. It handles the following unary operators:
// - Add (+)
// - Sub (-)
// - Compl (~)
// - Not (!)
//
// For each operator, it advances the lexer, processes the postfix expression, and applies the corresponding operation.
// If the operand is an identifier, it retrieves its value before applying the operation.
//
// Returns the evaluated Value or an error if the expression is invalid or an operation fails.
func (ex *Expression) unary() (Value, error) {
	var v Value
	var right Value
	var err error

	switch ex.next.t {
	case Add:
		if ex.next, err = ex.lex(); err != nil {
			if errors.Is(err, ErrEof) {
				return ex.next, syntaxError("expected expression", "")
			}
			return ex.next, err
		}
		if right, err = ex.postfix(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		v = right
		if v.IsIdentifier() {
			if v, err = v.getValue(); err != nil {
				return v, err
			}
		}
		if err = v.Plus(); err != nil {
			return v, err
		}
	case Sub:
		if ex.next, err = ex.lex(); err != nil {
			if errors.Is(err, ErrEof) {
				return ex.next, syntaxError("expected expression", "")
			}
			return ex.next, err
		}
		if right, err = ex.postfix(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		v = right
		if v.IsIdentifier() {
			if v, err = v.getValue(); err != nil {
				return v, err
			}
		}
		if err = v.Neg(); err != nil {
			return v, err
		}
	case Compl:
		if ex.next, err = ex.lex(); err != nil {
			if errors.Is(err, ErrEof) {
				return ex.next, syntaxError("expected expression", "")
			}
			return ex.next, err
		}
		if right, err = ex.postfix(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		v = right
		if v.IsIdentifier() {
			if v, err = v.getValue(); err != nil {
				return v, err
			}
		}
		if err = v.Compl(); err != nil {
			return v, err
		}
	case Not:
		if ex.next, err = ex.lex(); err != nil {
			if errors.Is(err, ErrEof) {
				return ex.next, syntaxError("expected expression", "")
			}
			return ex.next, err
		}
		if right, err = ex.postfix(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		v = right
		if v.IsIdentifier() {
			if v, err = v.getValue(); err != nil {
				return v, err
			}
		}
		if err = v.Not(); err != nil {
			return v, err
		}
	default:
		if v, err = ex.postfix(); err != nil {
			return v, err
		}
	}
	return v, nil
}

// unary
// ( type ) castExpr
// castExpr parses and evaluates a cast expression in the form of (Type) expression.
// It handles the following scenarios:
// 1. If the next token is an opening parenthesis, it checks if the following token is an identifier.
//   - If not, it resets the position and treats it as a unary expression.
//   - If it is an identifier, it checks if it is a valid type.
//   - If not, it resets the position and treats it as a unary expression.
//   - If it is a valid type, it expects a closing parenthesis and then recursively calls castExpr to evaluate the expression.
//   - If the evaluated value is an identifier, it retrieves its value.
//   - It then attempts to cast the value to the specified type.
//
// 2. If the next token is not an opening parenthesis, it treats the expression as a unary expression.
//
// Returns the evaluated value and any error encountered during parsing or evaluation.
func (ex *Expression) castExpr() (Value, error) {
	var v Value
	var err error

	start := ex.getPos()
	if ex.next.t == ParenO {
		if ex.next, err = ex.lex(); err != nil {
			return ex.next, err
		}
		if !ex.next.IsIdentifier() {
			ex.setPos(start)
			ex.next.t = ParenO
			if v, err = ex.unary(); err != nil && !errors.Is(err, ErrEof) {
				return v, err
			}
			return v, nil
		}
		var ty Type
		if ty = ITypes[ex.next.s]; ty == NoType {
			ex.setPos(start)
			ex.next.t = ParenO
			if v, err = ex.unary(); err != nil && !errors.Is(err, ErrEof) {
				return v, err
			}
			return v, nil
		}
		if ex.next, err = ex.lex(); err != nil {
			return ex.next, err
		}
		if ex.next.t != ParenC {
			return ex.next, syntaxError("expected \")\"", "")
		}
		if ex.next, err = ex.lex(); err != nil {
			return ex.next, err
		}
		if v, err = ex.castExpr(); err != nil && !errors.Is(err, ErrEof) {
			return v, err
		}
		if v.IsIdentifier() {
			if v, err = v.getValue(); err != nil {
				return v, err
			}
		}
		if err = v.Cast(ty); err != nil {
			return v, err
		}
	} else {
		if v, err = ex.unary(); err != nil && !errors.Is(err, ErrEof) {
			return v, err
		}
	}
	return v, nil
}

// castExpr
// mulExpr * castExpr
// mulExpr / castExpr
// mulExpr % castExpr
// mulExpr parses and evaluates a multiplication, division, or modulus expression.
// It handles the following operators: *, /, %.
// The function first parses the left-hand side expression and then iterates through
// the operators, parsing the right-hand side expression and performing the corresponding
// operation. If an identifier is encountered, it retrieves its value before performing
// the operation. The function returns the resulting value or an error if any occurs
// during parsing or evaluation.
func (ex *Expression) mulExpr() (Value, error) {
	var left Value
	var right Value
	var err error

	if left, err = ex.castExpr(); err != nil {
		return left, err
	}
loop:
	for {
		switch ex.next.t {
		case Mul:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.castExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.Mul(&right); err != nil {
				return left, err
			}
		case Div:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.castExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.Div(&right); err != nil {
				return left, err
			}
		case Mod:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.castExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.Mod(&right); err != nil {
				return left, err
			}
		default:
			break loop
		}
	}
	return left, nil
}

// mulExpr
// addExpr + mulExpr
// addExpr - mulExpr
// addExpr parses and evaluates an addition or subtraction expression.
// It handles the following operators: +, -.
// The function first parses the left-hand side expression and then iterates through
// the operators, parsing the right-hand side expression and performing the corresponding
// operation. If an identifier is encountered, it retrieves its value before performing
// the operation. The function returns the resulting value or an error if any occurs
// during parsing or evaluation.
func (ex *Expression) addExpr() (Value, error) {
	var left Value
	var right Value
	var err error

	if left, err = ex.mulExpr(); err != nil {
		return left, err
	}
loop:
	for {
		switch ex.next.t {
		case Add:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.mulExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.Add(&right); err != nil {
				return left, err
			}
		case Sub:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.mulExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.Sub(&right); err != nil {
				return left, err
			}
		default:
			break loop
		}
	}
	return left, nil
}

// addExpr
// shiftExpr << addExpr
// shiftExpr >> addExpr
// shiftExpr parses and evaluates shift expressions (both left and right shifts).
// It starts by evaluating the left-hand side of the shift expression using addExpr.
// Then, it enters a loop to handle multiple shift operations in sequence.
// For each shift operation, it checks the type of shift (Shl for left shift, Shr for right shift),
// consumes the next token, and evaluates the right-hand side of the shift expression using addExpr.
// If either the left or right value is an identifier, it resolves the identifier to its value.
// Finally, it performs the shift operation (Shl or Shr) and updates the left value accordingly.
// The function returns the final evaluated value or an error if any occurs during parsing or evaluation.
func (ex *Expression) shiftExpr() (Value, error) {
	var left Value
	var right Value
	var err error

	if left, err = ex.addExpr(); err != nil {
		return left, err
	}
loop:
	for {
		switch ex.next.t {
		case Shl:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.addExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.Shl(&right); err != nil {
				return left, err
			}
		case Shr:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.addExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.Shr(&right); err != nil {
				return left, err
			}
		default:
			break loop
		}
	}
	return left, nil
}

// shiftExpr
// relExpr < shiftExpr
// relExpr <= shiftExpr
// relExpr > shiftExpr
// relExpr >= shiftExpr
// relExpr parses and evaluates relational expressions in the form of
// "left < right", "left <= right", "left > right", and "left >= right".
// It first evaluates the left-hand side expression using shiftExpr, then
// iteratively processes relational operators and their corresponding
// right-hand side expressions. If the left or right values are identifiers,
// it resolves their actual values before performing the comparison.
// The function returns the final evaluated Value or an error if any occurs
// during parsing or evaluation.
func (ex *Expression) relExpr() (Value, error) {
	var left Value
	var right Value
	var err error

	if left, err = ex.shiftExpr(); err != nil {
		return left, err
	}
loop:
	for {
		switch ex.next.t {
		case Less:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.shiftExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.Less(&right); err != nil {
				return left, err
			}
		case LessEqual:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.shiftExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.LessEqual(&right); err != nil {
				return left, err
			}
		case Greater:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.shiftExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.Greater(&right); err != nil {
				return left, err
			}
		case GreaterEqual:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.shiftExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.GreaterEqual(&right); err != nil {
				return left, err
			}
		default:
			break loop
		}
	}
	return left, nil
}

// relExpr
// equExpr == relExpr
// equExpr != relExpr
// equExpr evaluates equality and inequality expressions in the Expression.
// It first evaluates a relational expression and then checks for equality (==) or inequality (!=) operators.
// If an equality or inequality operator is found, it evaluates the right-hand side expression and compares it with the left-hand side.
// The function handles identifiers by resolving their values before comparison.
// It returns the final evaluated Value and an error if any occurs during the evaluation process.
func (ex *Expression) equExpr() (Value, error) {
	var left Value
	var right Value
	var err error

	if left, err = ex.relExpr(); err != nil {
		return left, err
	}
loop:
	for {
		switch ex.next.t {
		case Equal:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.relExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.Equal(&right); err != nil {
				return left, err
			}
		case NotEqual:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.relExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.NotEqual(&right); err != nil {
				return left, err
			}
		default:
			break loop
		}
	}
	return left, nil
}

// equExpr
// andExpr & equExpr
// andExpr evaluates a logical AND expression in the context of the Expression struct.
// It first evaluates the left-hand side of the expression using equExpr().
// Then, it enters a loop to process any subsequent AND operations.
// For each AND operation, it evaluates the right-hand side of the expression using equExpr().
// If either side of the expression is an identifier, it resolves the identifier to its value.
// Finally, it performs the logical AND operation on the left and right values.
// The function returns the result of the AND operation or an error if any step fails.
func (ex *Expression) andExpr() (Value, error) {
	var left Value
	var right Value
	var err error

	if left, err = ex.equExpr(); err != nil {
		return left, err
	}
loop:
	for {
		switch ex.next.t {
		case And:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.equExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.And(&right); err != nil {
				return left, err
			}
		default:
			break loop
		}
	}
	return left, nil
}

// andExpr
// xorExpr ^ andExpr
// xorExpr evaluates an XOR expression in the context of the Expression struct.
// It first evaluates the left-hand side using the andExpr method. Then, it enters
// a loop to process any subsequent XOR operations. For each XOR operation, it
// evaluates the right-hand side using the andExpr method. If either side is an
// identifier, it retrieves its value. Finally, it performs the XOR operation on
// the left and right values. The method returns the final evaluated value or an
// error if any occurs during the evaluation process.
func (ex *Expression) xorExpr() (Value, error) {
	var left Value
	var right Value
	var err error

	if left, err = ex.andExpr(); err != nil {
		return left, err
	}
loop:
	for {
		switch ex.next.t {
		case Xor:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.andExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.Xor(&right); err != nil {
				return left, err
			}
		default:
			break loop
		}
	}
	return left, nil
}

// xorExpr
// orExpr | xorExpr
// orExpr evaluates an OR expression in the context of the Expression struct.
// It first evaluates the left-hand side of the expression using xorExpr.
// Then, it enters a loop to process any subsequent OR operations.
// For each OR operation, it evaluates the right-hand side using xorExpr,
// resolves any identifiers to their values, and performs the OR operation.
// If an error occurs during any of these steps, it returns the error.
// Otherwise, it returns the final evaluated value of the OR expression.
func (ex *Expression) orExpr() (Value, error) {
	var left Value
	var right Value
	var err error

	if left, err = ex.xorExpr(); err != nil {
		return left, err
	}
loop:
	for {
		switch ex.next.t {
		case Or:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.xorExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.Or(&right); err != nil {
				return left, err
			}
		default:
			break loop
		}
	}
	return left, nil
}

// orExpr
// logAndExpr && orExpr
// logAndExpr evaluates a logical AND expression.
// It first evaluates the left operand using the orExpr method.
// Then, it enters a loop to process any subsequent logical AND operations.
// For each logical AND operation, it evaluates the right operand using the orExpr method.
// If either operand is an identifier, it retrieves its value.
// Finally, it performs the logical AND operation on the left and right operands.
// The method returns the result of the logical AND operation or an error if any occurs during evaluation.
func (ex *Expression) logAndExpr() (Value, error) {
	var left Value
	var right Value
	var err error

	if left, err = ex.orExpr(); err != nil {
		return left, err
	}
loop:
	for {
		switch ex.next.t {
		case LogAnd:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.orExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.LogAnd(&right); err != nil {
				return left, err
			}
		default:
			break loop
		}
	}
	return left, nil
}

// logAndExpr
// logOrExpr || logAndExpr
// logOrExpr evaluates a logical OR expression.
// It first evaluates the left-hand side of the expression using logAndExpr.
// Then, it enters a loop to process any subsequent logical OR operators.
// For each logical OR operator, it evaluates the right-hand side of the expression using logAndExpr.
// If either side of the expression is an identifier, it retrieves its value.
// Finally, it performs the logical OR operation on the left and right values.
// The function returns the result of the logical OR expression or an error if any occurs during evaluation.
func (ex *Expression) logOrExpr() (Value, error) {
	var left Value
	var right Value
	var err error

	if left, err = ex.logAndExpr(); err != nil {
		return left, err
	}
loop:
	for {
		switch ex.next.t {
		case LogOr:
			if ex.next, err = ex.lex(); err != nil {
				return ex.next, err
			}
			if right, err = ex.logAndExpr(); err != nil && !errors.Is(err, ErrEof) {
				return right, err
			}
			if left.IsIdentifier() {
				if left, err = left.getValue(); err != nil {
					return left, err
				}
			}
			if right.IsIdentifier() {
				if right, err = right.getValue(); err != nil {
					return right, err
				}
			}
			if err = left.LogOr(&right); err != nil {
				return left, err
			}
		default:
			break loop
		}
	}
	return left, nil
}

// logOrExpr
// logOrExpr ? expression : asnExpr
// condExpr evaluates a conditional expression (ternary operator) in the form of
// "condition ? trueExpression : falseExpression". It first evaluates the condition
// using logOrExpr, then checks for the presence of the '?' token. If the token is
// present, it evaluates the trueExpression and checks for the ':' token. If the ':'
// token is present, it evaluates the falseExpression. Depending on the value of the
// condition, it returns either the trueExpression or the falseExpression. If the
// condition is an identifier, it resolves its value before making the comparison.
// It returns the evaluated Value or an error if any part of the expression fails to
// parse or evaluate correctly.
func (ex *Expression) condExpr() (Value, error) {
	var left Value
	var mid Value
	var right Value
	var err error

	if left, err = ex.logOrExpr(); err != nil {
		return left, err
	}
	if ex.next.t != Quest {
		return left, nil
	}
	if ex.next, err = ex.lex(); err != nil {
		return left, err
	}
	if mid, err = ex.expression(); err != nil {
		return mid, err
	}
	if ex.next.t != Colon {
		return left, syntaxError("missing : in conditional expression", "")
	}
	if ex.next, err = ex.lex(); err != nil {
		return mid, err
	}
	if right, err = ex.asnExpr(); err != nil && !errors.Is(err, ErrEof) {
		return right, err
	}
	if left.IsIdentifier() {
		if left, err = left.getValue(); err != nil {
			return left, err
		}
	}
	switch left.t {
	case Integer:
		if left.i != 0 {
			left = mid
		} else {
			left = right
		}
	case Floating:
		if left.f != 0.0 {
			left = mid
		} else {
			left = right
		}
	default:
		return left, typeError("equExpr", "")
	}
	if left.IsIdentifier() {
		if left, err = left.getValue(); err != nil {
			return left, err
		}
	}
	return left, nil
}

// condExpr
// condExpr ?= asnExpr
// asnExpr parses and evaluates an assignment expression.
// It supports various assignment operators such as ShlAssign, ShrAssign, PlusAssign, MinusAssign, OrAssign, AndAssign, XorAssign, MulAssign, DivAssign, ModAssign, and Assign.
// The function returns the resulting Value and an error if any occurs during parsing or evaluation.
// It ensures that the left-hand side of the assignment is an identifier and performs the appropriate operation based on the assignment operator.
func (ex *Expression) asnExpr() (Value, error) {
	var v Value
	var left Value
	var right Value
	var err error

	if left, err = ex.condExpr(); err != nil {
		return left, err
	}
	switch ex.next.t {
	case ShlAssign:
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
		if !left.IsIdentifier() {
			return left, syntaxError("assignment not to a variable", "")
		}
		if right, err = ex.asnExpr(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		if right.IsIdentifier() {
			if right, err = right.getValue(); err != nil {
				return right, err
			}
		}
		if v, err = left.getValue(); err != nil {
			return left, err
		}
		if err = v.Shl(&right); err != nil {
			return left, err
		}
		if err = left.setValue(&v); err != nil {
			return v, err // cannot happen because left was checked before
		}
	case ShrAssign:
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
		if !left.IsIdentifier() {
			return left, syntaxError("assignment not to a variable", "")
		}
		if right, err = ex.asnExpr(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		if right.IsIdentifier() {
			if right, err = right.getValue(); err != nil {
				return right, err
			}
		}
		if v, err = left.getValue(); err != nil {
			return left, err
		}
		if err = v.Shr(&right); err != nil {
			return left, err
		}
		if err = left.setValue(&v); err != nil {
			return v, err // cannot happen because left was checked before
		}
	case PlusAssign:
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
		if !left.IsIdentifier() {
			return left, syntaxError("assignment not to a variable", "")
		}
		if right, err = ex.asnExpr(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		if right.IsIdentifier() {
			if right, err = right.getValue(); err != nil {
				return right, err
			}
		}
		if v, err = left.getValue(); err != nil {
			return left, err
		}
		if err = v.Add(&right); err != nil {
			return left, err
		}
		if err = left.setValue(&v); err != nil {
			return v, err // cannot happen because left was checked before
		}
	case MinusAssign:
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
		if !left.IsIdentifier() {
			return left, syntaxError("assignment not to a variable", "")
		}
		if right, err = ex.asnExpr(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		if right.IsIdentifier() {
			if right, err = right.getValue(); err != nil {
				return right, err
			}
		}
		if v, err = left.getValue(); err != nil {
			return left, err
		}
		if err = v.Sub(&right); err != nil {
			return left, err
		}
		if err = left.setValue(&v); err != nil {
			return v, err // cannot happen because left was checked before
		}
	case OrAssign:
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
		if !left.IsIdentifier() {
			return left, syntaxError("assignment not to a variable", "")
		}
		if right, err = ex.asnExpr(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		if right.IsIdentifier() {
			if right, err = right.getValue(); err != nil {
				return right, err
			}
		}
		if v, err = left.getValue(); err != nil {
			return left, err
		}
		if err = v.Or(&right); err != nil {
			return left, err
		}
		if err = left.setValue(&v); err != nil {
			return v, err // cannot happen because left was checked before
		}
	case AndAssign:
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
		if !left.IsIdentifier() {
			return left, syntaxError("assignment not to a variable", "")
		}
		if right, err = ex.asnExpr(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		if right.IsIdentifier() {
			if right, err = right.getValue(); err != nil {
				return right, err
			}
		}
		if v, err = left.getValue(); err != nil {
			return left, err
		}
		if err = v.And(&right); err != nil {
			return left, err
		}
		if err = left.setValue(&v); err != nil {
			return v, err // cannot happen because left was checked before
		}
	case XorAssign:
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
		if !left.IsIdentifier() {
			return left, syntaxError("assignment not to a variable", "")
		}
		if right, err = ex.asnExpr(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		if right.IsIdentifier() {
			if right, err = right.getValue(); err != nil {
				return right, err
			}
		}
		if v, err = left.getValue(); err != nil {
			return left, err
		}
		if err = v.Xor(&right); err != nil {
			return left, err
		}
		if err = left.setValue(&v); err != nil {
			return v, err // cannot happen because left was checked before
		}
	case MulAssign:
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
		if !left.IsIdentifier() {
			return left, syntaxError("assignment not to a variable", "")
		}
		if right, err = ex.asnExpr(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		if right.IsIdentifier() {
			if right, err = right.getValue(); err != nil {
				return right, err
			}
		}
		if v, err = left.getValue(); err != nil {
			return left, err
		}
		if err = v.Mul(&right); err != nil {
			return left, err
		}
		if err = left.setValue(&v); err != nil {
			return v, err // cannot happen because left was checked before
		}
	case DivAssign:
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
		if !left.IsIdentifier() {
			return left, syntaxError("assignment not to a variable", "")
		}
		if right, err = ex.asnExpr(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		if right.IsIdentifier() {
			if right, err = right.getValue(); err != nil {
				return right, err
			}
		}
		if v, err = left.getValue(); err != nil {
			return left, err
		}
		if err = v.Div(&right); err != nil {
			return left, err
		}
		if err = left.setValue(&v); err != nil {
			return v, err // cannot happen because left was checked before
		}
	case ModAssign:
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
		if !left.IsIdentifier() {
			return left, syntaxError("assignment not to a variable", "")
		}
		if right, err = ex.asnExpr(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		if right.IsIdentifier() {
			if right, err = right.getValue(); err != nil {
				return right, err
			}
		}
		if v, err = left.getValue(); err != nil {
			return left, err
		}
		if err = v.Mod(&right); err != nil {
			return left, err
		}
		if err = left.setValue(&v); err != nil {
			return v, err // cannot happen because left was checked before
		}
	case Assign:
		if ex.next, err = ex.lex(); err != nil {
			return left, err
		}
		if !left.IsIdentifier() {
			return left, syntaxError("assignment not to a variable", "")
		}
		if right, err = ex.asnExpr(); err != nil && !errors.Is(err, ErrEof) {
			return right, err
		}
		if right.IsIdentifier() {
			if right, err = right.getValue(); err != nil {
				return right, err
			}
		}
		v = right
		if err = left.setValue(&v); err != nil {
			return v, err
		}
	default:
		v = left
	}
	return v, err
}

// asnExpr
// expression "," asnExpr
// expression ";" asnExpr
// expression evaluates the current expression and returns its value.
// It first attempts to parse an assignment expression. If the resulting
// value is an identifier, it retrieves the actual value associated with
// that identifier. The function then processes any subsequent comma or
// semicolon tokens, continuing to parse assignment expressions until
// the end of the expression is reached or an error occurs.
//
// Returns:
//   - Value: The evaluated value of the expression.
//   - error: An error if the expression could not be evaluated.
func (ex *Expression) expression() (Value, error) {
	var v Value
	var err error

	if v, err = ex.asnExpr(); err != nil {
		return v, err
	}
	if v.IsIdentifier() {
		if v, err = v.getValue(); err != nil {
			return v, err
		}
	}
	for ex.next.t == Comma || ex.next.t == Semi {
		if ex.next, err = ex.lex(); err != nil {
			return v, err
		}
		if _, err = ex.asnExpr(); err != nil && !errors.Is(err, ErrEof) {
			return v, err
		}
	}
	return v, nil
}
