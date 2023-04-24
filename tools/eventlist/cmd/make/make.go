/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
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
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/josephspurrier/goversioninfo"
)

const program = "eventlist"
const mainPath = "./cmd/" + program
const resourceFileName = "resource.syso"
const unknownVersion = "0.0.0.0"
const unknownYear = "2023"

var legalCopyright = "Arm Ltd. and Contributors"

// Errors
var ErrGitTag = errors.New("git tag error")
var ErrVersion = errors.New("version error")
var ErrCommand = errors.New("command error")

func reportError(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

type Options struct {
	targetOs   string
	targetArch string
	outDir     string
	covReport  string
}

type runner struct {
	options Options
	args    []string
}

func (r runner) run(command string) {
	switch {
	case command == "build":
		versionInfo, err := createResourceInfoFile(r.options.targetArch)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		info := versionInfo.StringFileInfo.FileVersion + " " + versionInfo.StringFileInfo.LegalCopyright
		if err = r.build(r.options, info); err != nil {
			fmt.Println(err.Error())
		}
	case command == "test":
		if err := r.test(); err != nil {
			fmt.Println(err.Error())
			return
		}
	case command == "coverage":
		if r.options.covReport == "" {
			if err := r.coverage(); err != nil {
				fmt.Println(err.Error())
				return
			}
		} else {
			if err := r.coverageReport(r.options.covReport); err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	case command == "lint":
		r.lint()
	case command == "format":
		r.format()
	}
}

func (r runner) executeCommand(command string) (err error) {
	var stdout, stderr bytes.Buffer
	fmt.Println(command)
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	stdoutStr := stdout.String()
	stderrStr := stderr.String()
	if stdoutStr != "" {
		fmt.Println(stdoutStr)
	}
	if stderrStr != "" {
		fmt.Println(stderrStr)
	}
	return err
}

func (r runner) build(options Options, versionInfo string) (err error) {
	var extn string

	if options.targetOs == "windows" {
		extn = ".exe"
	}
	cmd := "GOOS=" + options.targetOs + " GOARCH=" + options.targetArch +
		" go build -v -ldflags '-X \"main.versionInfo=" + versionInfo +
		"\"' -o " + options.outDir + "/" + program + extn + " " + mainPath

	if err = r.executeCommand(cmd); err == nil {
		fmt.Println("build finished successfully!")
	}
	return err
}

func (r runner) test() (err error) {
	args := "./..."
	if len(r.args) != 0 {
		args = strings.Join(r.args[:], " ")
	}
	return r.executeCommand("go test " + args)
}

func (r runner) coverage() (err error) {
	args := "./..."
	if len(r.args) != 0 {
		args = strings.Join(r.args[:], " ")
	}
	return r.executeCommand("go test -cover " + args)
}

func (r runner) coverageReport(covReport string) (err error) {
	covDir := path.Dir(covReport)
	if covReport == covDir {
		return reportError(ErrCommand, "invalid file path '"+covReport+"'")
	}

	if _, err = os.Stat(covDir); os.IsNotExist(err) {
		if err = os.Mkdir(covDir, os.ModePerm); err != nil {
			return
		}
	}
	err = r.executeCommand("go test ./... -coverprofile " + covDir + "/cover.out")
	if err != nil {
		return
	}
	err = r.executeCommand("go tool cover -html=" + covDir + "/cover.out -o " + covReport)
	if err == nil {
		fmt.Println("info: HTML coverage output written to " + covReport)
	}
	return
}

func (r runner) lint() {
	_ = r.executeCommand("golangci-lint run --config=./.golangci.yaml")
}

func (r runner) format() {
	_ = r.executeCommand("gofmt -s -w .")
}

func fetchVersionInfoFromGit() (version version, err error) {
	out, err := exec.Command("git", "describe", "--tags", "--match", "tools/eventlist/*").Output()
	if len(out) == 0 && err != nil {
		fmt.Println("warning: no release tag found, setting version to default \"" + unknownVersion + "\"")
		return newVersion(unknownVersion)
	}
	if err != nil {
		return
	}
	tag := strings.TrimSpace(string(out))
	if tag == "" {
		return version, reportError(ErrGitTag, "no git release tag found")
	}
	tokens := strings.Split(tag, "/")
	if len(tokens) != 3 {
		return version, reportError(ErrGitTag, "invalid release tag")
	}
	return newVersion(tokens[2])
}

func fetchChangeYearFromGit() (year string) {
	out, err := exec.Command("git", "log", "-n", "1", "--format=%ad", "--date=format:%Y").Output()
	if len(out) == 0 || err != nil {
		fmt.Println("warning: no change log found, setting year to default \"" + unknownYear + "\"")
		return unknownYear
	}
	return strings.TrimSpace(string(out))
}

func createResourceInfoFile(arch string) (version goversioninfo.VersionInfo, err error) {
	gitVersion, err := fetchVersionInfoFromGit()
	if err != nil {
		return
	}
	gitYear := fetchChangeYearFromGit()

	verInfo := goversioninfo.VersionInfo{}

	verInfo.FixedFileInfo.FileVersion = goversioninfo.FileVersion{
		Major: gitVersion.major,
		Minor: gitVersion.minor,
		Patch: gitVersion.patch,
		Build: gitVersion.numCommit,
	}

	verInfo.FixedFileInfo.ProductVersion = verInfo.FixedFileInfo.FileVersion
	verInfo.StringFileInfo = goversioninfo.StringFileInfo{
		FileDescription:  program,
		InternalName:     program,
		ProductName:      program,
		OriginalFilename: program + ".exe",
		FileVersion:      gitVersion.String(),
		ProductVersion:   gitVersion.String(),
		LegalCopyright:   "Copyright (c) 2022-" + gitYear + " " + legalCopyright,
	}
	verInfo.VarFileInfo.Translation = goversioninfo.Translation{
		LangID:    1033,
		CharsetID: 1200,
	}

	// Fill the structures with config data
	verInfo.Build()
	// Write the data to a buffer
	verInfo.Walk()

	return verInfo, verInfo.WriteSyso(mainPath+"/"+resourceFileName, arch)
}

func isCommandValid(command string) (result bool) {
	for _, cmd := range []string{
		"build", "coverage", "coverage-report",
		"format", "help", "lint", "test",
	} {
		if cmd == command {
			return true
		}
	}
	fmt.Println(reportError(ErrCommand, "invalid command").Error())
	return false
}

type version struct {
	major, minor, patch int
	numCommit           int
	shaCommit           string
}

func (v version) String() string {
	if v.Empty() {
		return unknownVersion
	}
	if v.shaCommit == "" && v.numCommit == 0 {
		return fmt.Sprintf("%d.%d.%d.%d", v.major, v.minor, v.patch, v.numCommit)
	}
	return fmt.Sprintf("%d.%d.%d-dev%d+%s", v.major, v.minor, v.patch, v.numCommit, v.shaCommit)
}

func (v version) Empty() bool {
	if v.major == 0 && v.minor == 0 && v.patch == 0 && v.shaCommit == "" && v.numCommit == 0 {
		return true
	}
	return false
}

func newVersion(verStr string) (ver version, err error) {
	if verStr == "" || verStr == unknownVersion {
		return
	}

	versionStr := strings.TrimSpace(verStr)
	tokens := strings.Split(versionStr, "-")
	numTokens := len(tokens)

	if !(numTokens == 1 || numTokens == 3) {
		return ver, reportError(ErrVersion, "invalid version string")
	}
	verParts := strings.Split(tokens[0], ".")
	if len(verParts) != 3 {
		return ver, reportError(ErrVersion, "invalid version string")
	}

	// Major
	ver.major, err = strconv.Atoi(verParts[0])
	if err != nil {
		return version{}, err
	}
	// Minor
	ver.minor, err = strconv.Atoi(verParts[1])
	if err != nil {
		return version{}, err
	}
	// Patch
	ver.patch, err = strconv.Atoi(verParts[2])
	if err != nil {
		return version{}, err
	}

	if numTokens == 3 {
		// Number of commits
		ver.numCommit, err = strconv.Atoi(tokens[1])
		if err != nil {
			return version{}, err
		}
		// SHA of commit
		ver.shaCommit = tokens[2]
	}
	return ver, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(reportError(ErrCommand, "invalid command").Error())
		os.Exit(1)
	}

	command := os.Args[1]
	if !isCommandValid(command) {
		os.Exit(1)
	}

	commFlag := flag.CommandLine
	targetOs := commFlag.String("os", runtime.GOOS, "Target Operating System")
	targetArch := commFlag.String("arch", runtime.GOARCH, "Target architecture")
	outDir := commFlag.String("outdir", ".", "Output directory")
	covReport := commFlag.String("html", "", "Coverage report")
	_ = commFlag.Parse(os.Args[2:])
	arguments := commFlag.Args()

	runner := runner{
		options: Options{
			targetOs:   *targetOs,
			targetArch: *targetArch,
			outDir:     *outDir,
			covReport:  *covReport,
		},
		args: arguments,
	}
	runner.run(command)
}
