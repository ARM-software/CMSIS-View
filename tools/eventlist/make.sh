#!/bin/bash

# -------------------------------------------------------
# Copyright (c) 2023 Arm Limited. All rights reserved.
#
# SPDX-License-Identifier: Apache-2.0
# -------------------------------------------------------

# usage
usage() {
  echo ""
  echo "Usage:"
  echo "  make.sh <command> [OPTIONS...]"
  echo ""
  echo "commands:"
  echo "  build           : Build executable"
  echo "  coverage        : Run tests with coverage info"
  echo "  format          : Align indentation and format code"
  echo "  lint            : Run linter"
  echo "  test            : Run all tests"
  echo ""
  echo "build options:"
  echo "  -arch arg       : Optional target architecture for e.g amd64 etc [default: host arch]"
  echo "  -os arg         : Optional target operating system for e.g windows, linux, darwin etc [default: host OS]"
  echo "  -outdir arg     : Optional output directory for executable generation [default: current directory]"
  echo ""
  echo "coverage options:"
  echo "  -html arg       : Coverage file path"
}

if [ $# -eq 0 ]
then
  usage
  exit 0
fi

for cmdline in "$@"
do
  if [[ "${cmdline}" == "help" || "${cmdline}" == "-h" || "${cmdline}" == "--help" ]]; then
    usage
    exit 0
  fi
  arg="${cmdline}"
  args+=("${arg}")
done

go run cmd/make/make.go "${args[@]}"

RESULT=$?
if [ $RESULT -ne 0 ]; then
  usage
  exit 1
fi
exit 0