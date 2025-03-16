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
)

type Member struct {
	Offset string
	IType  Type
	Enums  map[int64]string
}
type ITypedef struct {
	Size      uint32
	BigEndian bool
	Members   map[string]Member
}
type Typedefs map[string]ITypedef

// Eval evaluates a string expression and returns its computed Value.
// It takes a pointer to the string expression `s`, a map of type definitions `typedefs`,
// and a map `tdUsed` to track used type definitions.
// It returns the evaluated Value and an error if the evaluation fails.
//
// Parameters:
//   - s: A pointer to the string expression to be evaluated.
//   - typedefs: A map of type definitions used in the evaluation.
//   - tdUsed: A map to track used type definitions during the evaluation.
//
// Returns:
//   - Value: The computed value of the evaluated expression.
//   - error: An error if the evaluation fails.
func Eval(s *string, typedefs Typedefs, tdUsed map[string]string) (Value, error) {
	var ex Expression
	var v Value
	var err error

	ex.in = s
	ex.pos = 0
	ex.typedefs = typedefs
	ex.tdUsed = tdUsed
	if ex.next, err = ex.lex(); err != nil {
		return v, err
	}
	return ex.expression()
}

// GetValue evaluates the value of the Enum and returns it as an int64.
// If an error occurs during evaluation, it returns the error unless the error is eval.ErrEof.
//
// Returns:
//   - int64: The evaluated integer value of the Enum.
//   - error: An error if the evaluation fails, except for eval.ErrEof.
func GetValue(value string, typedefs Typedefs) (int64, error) {
	n, err := Eval(&value, typedefs, nil)
	if err != nil && !errors.Is(err, ErrEof) {
		return 0, err
	}
	return n.GetInt(), nil
}

// GetIdValue evaluates the ID and returns its value as an IDType(uint16).
// If an error occurs during evaluation, it returns 0 and the error.
// It ignores eval.ErrEof.
func GetIdValue(id string, typedefs Typedefs) (uint16, error) { //nolint:golint,revive
	n, err := Eval(&id, typedefs, nil)
	if err != nil && !errors.Is(err, ErrEof) {
		return 0, err
	}
	return uint16(n.GetInt()), nil
}
