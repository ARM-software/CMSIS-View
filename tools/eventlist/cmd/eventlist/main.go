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

package main

import (
	"eventlist/pkg/elf"
	"eventlist/pkg/eval"
	"eventlist/pkg/output"
	"eventlist/pkg/xml/scvd"
	"flag"
	"fmt"
	"os"
	"strings"
)

var Progname string
var versionInfo string

type includes []string

// String returns the first element of the includes slice as a string.
// If the slice is nil or empty, it returns an empty string.
func (s *includes) String() string {
	if s == nil || len(*s) == 0 {
		return ""
	}
	return (*s)[0]
}

// Set appends the given string value to the includes slice.
// It implements the flag.Value interface, allowing the includes
// slice to be used as a command-line flag.
//
// Parameters:
//
//	v - the string value to be appended to the includes slice.
//
// Returns:
//
//	An error, which is always nil in this implementation.
func (s *includes) Set(v string) error {
	*s = append(*s, v)
	return nil
}

var paths includes

// infoOpt prints information about a command-line option.
//
// Parameters:
//   - flags: A FlagSet containing the defined command-line flags.
//   - sopt: The short option name (e.g., "h" for "-h").
//   - lopt: The long option name (e.g., "help" for "--help").
//   - arg: A boolean indicating whether the option requires an argument.
//
// Returns:
//   - error: An error if any of the printing operations fail, otherwise nil.
//
// The function prints the short and/or long option names, followed by
// "arg" if the option requires an argument. It then prints the usage
// information for the option, or "unknown option" if the option is not
// found in the FlagSet.
func infoOpt(flags *flag.FlagSet, sopt string, lopt string, arg bool) error {
	pos, err := fmt.Print("  ")
	if err != nil {
		return err
	}
	var n int
	if sopt != "" {
		if n, err = fmt.Printf("-%s", sopt); err != nil {
			return err
		}
		pos += n
	}
	if lopt != "" {
		if sopt == "" {
			if n, err = fmt.Printf("    "); err != nil {
				return err
			}
		} else {
			if n, err = fmt.Printf(", "); err != nil {
				return err
			}
		}
		pos += n
		if n, err = fmt.Printf("--%s", lopt); err != nil {
			return err
		}
		pos += n
	}
	if arg {
		if sopt == "" && lopt == "" {
			if n, err = fmt.Printf("  "); err != nil {
				return err
			}
			pos += n
		}
		if n, err = fmt.Printf(" arg"); err != nil {
			return err
		}
		pos += n
	}
	fmt.Printf("%*s", 22-pos, " ")
	if lopt == "help" {
		fmt.Printf("%s\n", "Print usage")
	} else {
		f := flags.Lookup(sopt)
		if f == nil {
			fmt.Printf("%s\n", "unknown option")
		} else {
			fmt.Printf("%s\n", f.Usage)
		}
	}
	return nil
}

// main is the entry point of the event listing tool. It parses command-line
// arguments, sets up the necessary configurations, and processes the event
// log file. The tool supports various options such as specifying an output
// file, format type, and level of detail. It also includes options for
// displaying version information and statistics.
//
// Usage:
//
//	eventlist [options] <logFile>
//
// Options:
//
//	-a <file>        Application file: elf/axf file name
//	-b, --begin      Output order: show statistic before events
//	-h, --help       Show help message
//	-I <file>        Include SCVD file name(s)
//	-o <file>        Output file
//	-s, --statistic  Output: show statistic but no events
//	-V, --version    Show version info
//	-f <format>      Output format: txt, json, xml
//	-l <level>       Level: Error|API|Op|Detail
func main() {
	var err error
	Progname = os.Args[0]
	idx := strings.LastIndexByte(Progname, '/')
	if idx == -1 {
		idx = strings.LastIndexByte(Progname, '\\')
	}
	if idx >= 0 {
		Progname = Progname[idx+1:]
	}
	if idx = strings.LastIndexByte(Progname, '.'); idx > 0 {
		Progname = Progname[:idx]
	}

	commFlag := flag.CommandLine

	// --- this is only for unit tests of main()
	testRun := flag.Lookup("test.run")
	if testRun != nil {
		commFlag = flag.NewFlagSet("test", flag.ContinueOnError)
		flag.CommandLine.VisitAll(func(flag *flag.Flag) {
			commFlag.Var(flag.Value, flag.Name, flag.Usage)
		})
	}
	// ---

	usage := false

	commFlag.Usage = func() {
		fmt.Printf("%s: Event Listing %s\n\n", Progname, versionInfo)
		fmt.Printf("Usage:\n  %s [options] <logFile>\n\n", Progname)
		fmt.Printf("Options:\n")
		_ = infoOpt(commFlag, "a", "", true)
		_ = infoOpt(commFlag, "b", "begin", false)
		_ = infoOpt(commFlag, "h", "help", false)
		_ = infoOpt(commFlag, "I", "", true)
		_ = infoOpt(commFlag, "o", "", true)
		_ = infoOpt(commFlag, "s", "statistic", false)
		_ = infoOpt(commFlag, "V", "version", false)
		_ = infoOpt(commFlag, "f", "format", true)
		_ = infoOpt(commFlag, "l", "level", true)
		usage = true
	}
	// parse command line
	commFlag.Var(&paths, "I", "[...] Include SCVD file name(s)")
	outputFile := commFlag.String("o", "", "Output file")
	elfFile := commFlag.String("a", "", "Application file: elf/axf file name")
	formatType := commFlag.String("f", "", "Output format: txt, json, xml")
	level := commFlag.String("l", "", "Level: Error|API|Op|Detail")
	var statBegin bool
	commFlag.BoolVar(&statBegin, "b", false, "Output order: show statistic before events")
	commFlag.BoolVar(&statBegin, "begin", false, "Output order: show statistic before events")
	var showVersion bool
	commFlag.BoolVar(&showVersion, "V", false, "Show version info")
	commFlag.BoolVar(&showVersion, "version", false, "Show version info")
	var showStatistic bool
	commFlag.BoolVar(&showStatistic, "s", false, "Output: show statistic but no events")
	commFlag.BoolVar(&showStatistic, "statistic", false, "Output: show statistic but no events")
	err = commFlag.Parse(os.Args[1:])

	if usage || err != nil {
		return
	}

	if len(os.Args) == 1 {
		commFlag.Usage()
		return
	}

	if showVersion {
		fmt.Printf("%s: Event Listing %s\n", Progname, versionInfo)
		return
	}

	eventFile := commFlag.Args()

	if len(eventFile) == 0 {
		fmt.Println(Progname + ": missing input file")
		return
	}
	if len(eventFile) > 1 {
		fmt.Println(Progname + ": only one binary input file allowed")
		return
	}

	if elfFile != nil && len(*elfFile) != 0 {
		if err = elf.Sections.Readelf(elfFile); err != nil {
			fmt.Print(Progname + ": ")
			fmt.Println(err)
			return
		}
	}
	evdefs := make(scvd.Events)
	typedefs := make(eval.Typedefs)

	var p []string = paths
	if err = scvd.Get(&p, evdefs, typedefs); err != nil {
		fmt.Print(Progname + ": ")
		fmt.Println(err)
		return
	}

	if err := output.Print(outputFile, formatType, level, &eventFile[0], evdefs, typedefs, statBegin, showStatistic); err != nil {
		fmt.Print(Progname + ": ")
		fmt.Println(err)
	}
}
