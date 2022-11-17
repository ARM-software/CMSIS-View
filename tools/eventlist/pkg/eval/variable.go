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

import "sync"

type Variable struct {
	n string
	v Value
}

var mu sync.Mutex
var names map[string]*Variable

func ClearNames() {
	mu.Lock()
	for k := range names {
		delete(names, k)
	}
	mu.Unlock()
}

func CountNames() int {
	return len(names)
}

func GetVar(n string) (*Variable, error) {
	var v *Variable
	var ok bool
	mu.Lock()
	defer mu.Unlock()
	if v, ok = names[n]; !ok {
		return v, syntaxError("unkown variable name", "")
	}
	return v, nil
}

func SetVarI(n string, i int64) *Variable {
	val := Value{t: Integer, i: i}
	v := new(Variable)
	v.n = n
	v.v = val
	mu.Lock()
	defer mu.Unlock()
	if len(names) == 0 {
		names = make(map[string]*Variable)
	}
	names[n] = v
	return v
}

func SetVar(n string, val Value) *Variable {
	v := new(Variable)
	v.n = n
	v.v = val
	mu.Lock()
	defer mu.Unlock()
	if len(names) == 0 {
		names = make(map[string]*Variable)
	}
	names[n] = v
	return v
}

func (v *Variable) setValue(val *Value) error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := names[v.n]; !ok {
		return typeError("Unknown variable", v.v.s)
	}
	vari := new(Variable)
	vari.n = v.n
	vari.v = *val
	names[v.n] = vari
	return nil
}

func (v *Variable) getValue() (Value, error) {
	mu.Lock()
	defer mu.Unlock()
	val, ok := names[v.n]
	if !ok {
		return Value{}, typeError("Unknown variable", v.v.s)
	}
	return val.v, nil
}
