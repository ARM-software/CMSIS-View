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

type TdMember struct {
	Type   Token
	Enum   map[int16]string
	Offset uint32
	Info   string
}

// An internal "copy" of scvd.Typedef to be used by this package
type TdTypedef struct {
	Size    int
	Members map[string]TdMember
}

var Typedefs map[string]TdTypedef

type EventVals struct {
	Val1	string
	Val2	string
	Val3	string
	Val4	string
	Val5	string
	Val6	string
}
var EventV EventVals

func Eval(s *string, redefs map[string]string) (Value, error) {
	var ex Expression
	var v Value
	var err error


	ex.in = s
	ex.pos = 0
	ex.redefs = redefs;
	if ex.next, err = ex.lex(); err != nil {
		return v, err
	}
	return ex.expression()
}
