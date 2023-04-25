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

package main

import (
	"flag"
	"io"
	"os"
	"reflect"
	"regexp"
	"testing"
)

func Test_includes_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    *includes
		want string
	}{
		{"nil", nil, ""},
		{"one", &includes{"ab"}, "ab"},
		{"empty", &includes{}, ""},
		{"two", &includes{"cd", "ab"}, "cd"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.s.String(); got != tt.want {
				t.Errorf("includes.String() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_includes_Set(t *testing.T) {
	t.Parallel()

	type args struct {
		v string
	}
	tests := []struct {
		name    string
		s       *includes
		args    args
		want    *includes
		wantErr bool
	}{
		{"to_one", &includes{"ab"}, args{"cd"}, &includes{"ab", "cd"}, false},
		{"to_empty", &includes{}, args{"ab"}, &includes{"ab"}, false},
		{"empty", &includes{}, args{}, &includes{""}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.s.Set(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("includes.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.s, tt.want) {
				t.Errorf("includes.Set() %s = %v, want %v", tt.name, tt.s, tt.want)
			}
		})
	}
}

func Test_infoOpt(t *testing.T) { //nolint:golint,paralleltest
	type args struct {
		sopt string
		lopt string
		opt  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test.run opt", args{"test.run", "", "ef"}, "\t-test.run ef yy\trun only tests and examples matching `regexp`\n"},
		{"test.run", args{"test.run", "", ""}, "\t-test.run\trun only tests and examples matching `regexp`\n"},
		{"test help", args{"", "help", ""}, "\t--help\tshow short help\n"},
		{"test", args{"", "", ""}, "\t\tunknown option\n"},
		{"test s", args{"a", "", ""}, "\t-a\tunknown option\n"},
		{"test l", args{"", "cd", ""}, "\t--cd\tunknown option\n"},
		{"test s l", args{"a", "cd", ""}, "\t-a --cd\tunknown option\n"},
		{"test opt", args{"", "", "ef"}, "\t ef\tunknown option\n"},
		{"test s opt", args{"a", "", "ef"}, "\t-a ef\tunknown option\n"},
		{"test l opt", args{"", "cd", "ef"}, "\t--cd ef\tunknown option\n"},
		{"test s l opt", args{"a", "cd", "ef"}, "\t-a --cd ef\tunknown option\n"},
	}
	_ = flag.Set("test.run", "yy")
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			oldOut := os.Stdout
			restore := func() {
				os.Stdout = oldOut
			}
			defer restore()
			r, w, _ := os.Pipe()
			os.Stdout = w
			infoOpt(flag.CommandLine, tt.args.sopt, tt.args.lopt, tt.args.opt)
			w.Close()
			buf, _ := io.ReadAll(r)
			if string(buf) != tt.want {
				t.Errorf("infoOpt() %s = %v, want %v", tt.name, string(buf), tt.want)
			}
		})
	}
}

func Test_main(t *testing.T) { //nolint:golint,paralleltest
	outFile := "out.out"

	lines1 :=
		"   Detailed event list\\n" +
			"   -------------------\\n" +
			"\\n" +
			"Index Time \\(s\\)   Component Event Property Value\\n" +
			"----- --------   --------- -------------- -----\\n" +
			"    0 7\\.75000000 0xFF      0xFF03         val1=0x00000004, val2=0x00000002\\n" +
			"    1 7\\.75000000 0xFE      0xFE00         \"hello wo\"\\n" +
			"\\n" +
			"   Start/Stop event statistic\\n" +
			"   --------------------------\\n" +
			"\\n" +
			"Event count      total       min         max         average     first       last\\n" +
			"----- -----      -----       ---         ---         -------     -----       ----\\n"

	lines2 :=
		"   Start/Stop event statistic\\n" +
			"   --------------------------\\n" +
			"\\n" +
			"Event count      total       min         max         average     first       last\\n" +
			"----- -----      -----       ---         ---         -------     -----       ----\\n"

	help :=
		"Usage: [^ ]+ \\[-I <scvdFile>\\]\\.\\.\\. \\[-o <outputFile>\\] \\[-a <elf/axfFile>\\] \\[-b\\] <logFile>\\n" +
			"\\t-a <fileName> \\telf/axf file name\\n" +
			"\\t-b --begin\\tshow statistic at beginning\\n" +
			"\\t-h --help\\tshow short help\\n" +
			"\\t-I <fileName> \\tinclude SCVD file name\\n" +
			"\\t-o <fileName> \\toutput file name\\n" +
			"\\t-s --statistic\\tshow statistic only\\n" +
			"\\t-V --version\\tshow version info\\n"

	versionInfo = "1.2.3 (C) 2022 Arm Ltd. and Contributors"
	tests := []struct {
		name       string
		args       []string
		want       string
		removefile string
	}{
		{"-a", []string{"-a", "../../testdata/nix", "xxx"}, ".*: open ../../testdata/nix: (no such file or directory|The system cannot find the file specified.)\\n", ""},
		{"-s stdout", []string{"-s", "../../testdata/test10.binary"}, lines2, ""},
		{"-s", []string{"-s", "-o", outFile, "../../testdata/test10.binary"}, "", outFile},
		{"-statistic", []string{"-statistic", "-o", outFile, "../../testdata/test10.binary"}, "", outFile},
		{"-help", []string{"-help"}, help, ""},
		{"stdout", []string{"../../testdata/test10.binary"}, lines1, ""},
		{"-o -begin", []string{"-begin", "-o", outFile, "../../testdata/test10.binary"}, "", outFile},
		{"-o -b", []string{"-b", "-o", outFile, "../../testdata/test10.binary"}, "", outFile},
		{"-o", []string{"-o", outFile, "../../testdata/test10.binary"}, "", outFile},
		{"-o", []string{"-o", outFile, "../../testdata/nix"}, ".*: cannot open event file\\n", outFile},
		{"-V", []string{"-V"}, ".* [0-9]+\\.[0-9]+\\.[0-9]+ \\(C\\) [0-9]+ Arm Ltd. and Contributors\\n", ""},
		{"-version", []string{"-version"}, ".* [0-9]+\\.[0-9]+\\.[0-9]+ \\(C\\) [0-9]+ Arm Ltd. and Contributors\\n", ""},
		{"err", []string{"xxx", "yyy"}, ".*: only one binary input file allowed\n", ""},
		{"missing", nil, ".*: missing input file\n", ""},
		// -I must be the last test
		{"-I", []string{"-I", "../../testdata/nix", "xxx"}, ".*: open ../../testdata/nix: (no such file or directory|The system cannot find the file specified.)\\n", ""},
	}
	savedArgs := os.Args
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			oldOut := os.Stdout
			restore := func() {
				os.Stdout = oldOut
			}
			defer restore()
			defer os.Remove(tt.removefile)
			r, w, _ := os.Pipe()
			os.Stdout = w
			os.Args = append(savedArgs, tt.args...)
			main()
			w.Close()
			buf, _ := io.ReadAll(r)
			match, err := regexp.Match(tt.want, buf)
			if err != nil {
				t.Errorf("main() %s regexp match error %v, want %v", tt.name, err, tt.want)
			}
			if !match {
				t.Errorf("main() %s = %v, want %v", tt.name, string(buf), tt.want)
			}
		})
	}
}
