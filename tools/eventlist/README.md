# eventlist

This utility is a command line tool for processing Event Recorder records stored to a log file.

## Usage

```bash
Usage:
  eventlist [-I <scvdFile>]... [-o <outputFile>] [-a <elf/axfFile>] [-b] <logFile>

Flags:
  -a <fileName>   elf/axf file name
  -b --begin      show statistic at beginning
  -h --help       show short help
  -I <fileName>   include SCVD file name
  -o <fileName>   output file name
  -s --statistic  show statistic only
  -V --version    show version info
```

## Building the tool locally

This section contains a complete guide to get you the project build on
your local machine for development and testing purposes.

## Prerequisites

The following applications are required to be installed on your
machine to allow **eventlist** utility to be built and run.

Note that some of the required tools are platform dependent:

- [Git](https://git-scm.com/)
- [golang](https://go.dev/doc/install) (minimum recommended version **1.17**)
- Platform specific command line terminal
  - **Windows:**
    - [GIT Bash](https://gitforwindows.org/)
    ```
    ☑️ Note:
        Make sure 'git' and 'bash' paths are listed under the PATH environment
        variable and set the git bash priority higher in the path.
    ```
  - **Linux/MacOS:**
    - GNU Bash (minimum recommended version **5.0.17**)

## Clone repository

Clone GitHub repository to create a local copy on your computer to make
it easier to develop and test. Cloning of repository can be done by following
the below git command:

```bash
git clone git@github.com:ARM-software/CMSIS-View.git
```

## Build components

The commands below demonstrate how to build and create executable:

  - Go to eventlist directory
    > cd \<**root**\>/tools/eventlist
  - Create and switch to build directory
    ```bash
    mkdir build
    cd build
    ```
  - Run go command to build an executable
    > go build ./..

## Run Tests

One can directly run the tests from command line.
  - Go to eventlist directory
    > cd \<**root**\>/tools/eventlist
  - Clean existing cache test results
    > go clean -cache
  - Run the executable
    > go test ./...

## License

**eventlist** is licensed under Apache 2.0.

## Log File Format

The log file is expected to use the [Common Trace Format](https://diamon.org/ctf/#specification). The binary trace
stream layout is describes using the *Trace Stream Description Language* (TSDL) in
[eventlist.tsdl](docs/eventlist.tsdl).
