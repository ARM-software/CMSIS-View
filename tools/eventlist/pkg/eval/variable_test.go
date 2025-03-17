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

//nolint:golint,paralleltest
package eval

import (
	"reflect"
	"testing"
)

func TestClearNames(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetVar("v1_ClearNames", Value{t: Integer, i: 789})
			ClearNames()
			if CountNames() != 0 {
				t.Errorf("ClearNames() = %v, want %v", CountNames(), 0)
			}
		})
	}
}

func TestCountNames(t *testing.T) {
	tests := []struct {
		name  string
		clear bool
		want  int
	}{
		{"testEmpty", true, 0},
		{"testOne", false, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clear {
				ClearNames()
			} else {
				SetVar("v1_CountNames", Value{t: Integer, i: 789})
			}
			if got := CountNames(); got != tt.want {
				t.Errorf("CountNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVar(t *testing.T) {
	var vari *Variable
	type args struct {
		n string
	}
	tests := []struct {
		name    string
		clear   bool
		args    args
		want    *Variable
		wantErr bool
	}{
		{"GetVar_empty", true, args{"v1_GetVar"}, nil, true},
		{"GetVar_ok", false, args{"v1_GetVar"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clear {
				ClearNames()
			} else {
				tt.want = vari
			}
			got, err := GetVar(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVar() = %v, want %v", got, tt.want)
			}
			vari = SetVar(tt.args.n, Value{t: Integer, i: 345})
		})
	}
}

func TestSetVarI(t *testing.T) {
	type args struct {
		n string
		i int64
	}
	tests := []struct {
		name  string
		clear bool
		args  args
	}{
		{"SetVarI_empty", true, args{"v1_SetVarI", 345}},
		{"SetVarI_over", false, args{"v1_SetVarI", 123}},
		{"SetVarI_new", false, args{"v1_SetVarI1", 159}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clear {
				ClearNames()
			}
			SetVarI(tt.args.n, tt.args.i)
			v := Variable{tt.args.n, Value{}}
			vari := Value{t: Integer, i: tt.args.i}
			got, err := v.getValue()
			if err != nil {
				t.Errorf("SetVarI() = %v", err)
			} else if !reflect.DeepEqual(got, vari) {
				t.Errorf("SetVarI() = %v, want %v", got, vari)
			}
			SetVar("v1_SetVarI", Value{t: Integer, i: 789})
		})
	}
}

func TestSetVar(t *testing.T) {
	type args struct {
		n   string
		val Value
	}
	tests := []struct {
		name string
		args args
	}{
		{"SetVar_empty", args{"v1_SetVar", Value{t: Integer, i: 345}}},
		{"SetVar_over", args{"v1_SetVar", Value{t: Integer, i: 123}}},
		{"SetVar_new", args{"v1_SetVar1", Value{t: Integer, i: 159}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetVar(tt.args.n, tt.args.val)
			v := Variable{tt.args.n, Value{}}
			got, err := v.getValue()
			if err != nil {
				t.Errorf("SetVar() = %v", err)
			} else if !reflect.DeepEqual(got, tt.args.val) {
				t.Errorf("SetVar() = %v, want %v", got, tt.args.val)
			}
			SetVar("v1_SetVar", Value{t: Integer, i: 789})
		})
	}
}

func TestVariable_setValue(t *testing.T) {
	type fields struct {
		n string
		v Value
	}
	type args struct {
		val Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"setValue_ok", fields{"v1_setValue", Value{t: Integer, i: 123}}, args{Value{t: Integer, i: 345}}, false},
		{"setValue_nok", fields{"v1_xxxx", Value{t: Integer, i: 123}}, args{Value{t: Integer, i: 345}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Variable{
				n: tt.fields.n,
				v: tt.fields.v,
			}
			SetVar("v1_setValue", Value{t: Integer, i: 789})
			var err error
			if err = v.setValue(&tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Variable.setValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			var got Value
			if err == nil {
				got, err = v.getValue()
				if err != nil {
					t.Errorf("Variable.setValue() = %v", err)
				} else if !reflect.DeepEqual(got, tt.args.val) {
					t.Errorf("Variable.setValue() = %v, want %v", got, tt.args.val)
				}
			}
		})
	}
}

func TestVariable_getValue(t *testing.T) {
	t.Parallel()

	type fields struct {
		n string
		v Value
	}
	tests := []struct {
		name   string
		fields fields
		want   Value
	}{
		{"getValue", fields{"v1_getValue", Value{t: Integer, i: 123}}, Value{t: Integer, i: 789}},
	}
	SetVar("v1_getValue", Value{t: Integer, i: 789})
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := &Variable{
				n: tt.fields.n,
				v: tt.fields.v,
			}
			got, err := v.getValue()
			if err != nil {
				t.Errorf("Variable.getValue() = %v", err)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Variable.getValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
