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

import "sync"

type Variable struct {
	n string
	v Value
}

var mu sync.Mutex
var names map[string]*Variable

// ClearNames removes all entries from the 'names' map.
// It acquires a lock before modifying the map to ensure
// thread safety and releases the lock after the operation
// is complete.
func ClearNames() {
	mu.Lock()
	for k := range names {
		delete(names, k)
	}
	mu.Unlock()
}

// CountNames returns the number of elements in the 'names' slice.
// It calculates the length of the 'names' slice and returns it as an integer.
func CountNames() int {
	return len(names)
}

// GetVar retrieves a variable by its name from a global map.
// It returns a pointer to the Variable and an error if the variable name is not found.
//
// Parameters:
//   - n: The name of the variable to retrieve.
//
// Returns:
//   - *Variable: A pointer to the Variable if found, otherwise nil.
//   - error: An error if the variable name is not found, otherwise nil.
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

// SetVarI sets a variable with the given name and integer value.
// It creates a new Variable instance, assigns the provided name and value,
// and stores it in the global map of variable names. The function is thread-safe.
//
// Parameters:
//   - n: The name of the variable to set.
//   - i: The integer value to assign to the variable.
//
// Returns:
//
//	A pointer to the newly created Variable instance.
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

// SetVar creates a new Variable with the given name and value,
// stores it in the global map, and returns a pointer to the Variable.
// It ensures thread safety by locking and unlocking a mutex.
//
// Parameters:
//   - n: The name of the variable.
//   - val: The value to be assigned to the variable.
//
// Returns:
//
//	A pointer to the newly created Variable.
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

// setValue assigns a new value to the variable. It locks the mutex to ensure
// thread safety, checks if the variable name exists in the names map, and if
// it does, creates a new Variable instance with the provided value and updates
// the map. If the variable name does not exist, it returns a typeError.
//
// Parameters:
//   - val: A pointer to the Value to be assigned to the variable.
//
// Returns:
//   - error: Returns a typeError if the variable name is unknown, otherwise nil.
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

// getValue retrieves the value associated with the Variable instance.
// It locks the mutex to ensure thread safety while accessing the shared resource.
// If the variable name is not found in the names map, it returns an error indicating an unknown variable.
// Otherwise, it returns the corresponding value.
//
// Returns:
//   - Value: The value associated with the variable.
//   - error: An error if the variable name is not found.
func (v *Variable) getValue() (Value, error) {
	mu.Lock()
	defer mu.Unlock()
	val, ok := names[v.n]
	if !ok {
		return Value{}, typeError("Unknown variable", v.v.s)
	}
	return val.v, nil
}
