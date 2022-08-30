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

package main

//go:generate goversioninfo -gofile=versioninfo.go

import (
	"eventlist/elf"
	"eventlist/output"
	"eventlist/xml/scvd"
	"flag"
	"fmt"
	"os"
	"strings"
)

var Progname string

type includes []string

func (s *includes) String() string {
	if s == nil || len(*s) == 0 {
		return ""
	}
	return (*s)[0]
}

func (s *includes) Set(v string) error {
	*s = append(*s, v)
	return nil
}

var paths includes

func infoOpt(flags *flag.FlagSet, sopt string, lopt string, opt string) {
	fmt.Print("\t")
	if sopt != "" {
		fmt.Printf("-%s", sopt)
	}
	if lopt != "" {
		if sopt != "" {
			fmt.Print(" ")
		}
		fmt.Printf("--%s", lopt)
	}
	if opt != "" {
		fmt.Printf(" %s", opt)
	}
	if lopt == "help" {
		fmt.Printf("\t%s\n", "show short help")
	} else {
		f := flags.Lookup(sopt)
		if f == nil {
			fmt.Printf("\t%s\n", "unknown option")
		} else {
			if opt == "" {
				fmt.Printf("\t%s\n", f.Usage)
			} else {
				fmt.Printf(" %s\t%s\n", f.Value.String(), f.Usage)
			}
		}
	}
}

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
	fmt.Println(os.Args)

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
		fmt.Printf("Usage: %s [-I <scvdFile>]... [-o <outputFile>] [-a <elf/axfFile>] [-b] <logFile>\n",
			Progname)
		infoOpt(commFlag, "a", "", "<fileName>")
		infoOpt(commFlag, "b", "begin", "")
		infoOpt(commFlag, "h", "help", "")
		infoOpt(commFlag, "I", "", "<fileName>")
		infoOpt(commFlag, "o", "", "<fileName>")
		infoOpt(commFlag, "s", "statistic", "")
		infoOpt(commFlag, "V", "version", "")
		usage = true
	}
	// parse command line
	commFlag.Var(&paths, "I", "include SCVD file name")
	outputFile := commFlag.String("o", "", "output file name")
	elfFile := commFlag.String("a", "", "elf/axf file name")
	var statBegin bool
	commFlag.BoolVar(&statBegin, "b", false, "show statistic at beginning")
	commFlag.BoolVar(&statBegin, "begin", false, "show statistic at beginning")
	var showVersion bool
	commFlag.BoolVar(&showVersion, "V", false, "show version info")
	commFlag.BoolVar(&showVersion, "version", false, "show version info")
	var showStatistic bool
	commFlag.BoolVar(&showStatistic, "s", false, "show statistic only")
	commFlag.BoolVar(&showStatistic, "statistic", false, "show statistic only")
	err = commFlag.Parse(os.Args[1:])

	if usage || err != nil {
		return
	}

	if showVersion {
		version := versionInfo.StringFileInfo.ProductVersion
		i := strings.LastIndex(version, ".")
		if i > 0 {
			version = version[:i]
		}
		fmt.Printf("%s: Version %s\n", Progname, version)
		fmt.Printf("%s\n", versionInfo.StringFileInfo.LegalCopyright)
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
	evdefs := make(map[uint16]scvd.Event)
	typedefs := make(map[string]map[string]map[int16]string)

	var p []string = paths
	if err = scvd.Get(&p, evdefs, typedefs); err != nil {
		fmt.Print(Progname + ": ")
		fmt.Println(err)
		return
	}

	if err := output.Print(outputFile, &eventFile[0], evdefs, typedefs, statBegin, showStatistic); err != nil {
		fmt.Print(Progname + ": ")
		fmt.Println(err)
	}
}
