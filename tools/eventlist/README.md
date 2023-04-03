# eventlist

This utility is a command line tool for processing Event Recorder records stored to a log file.

## Usage

```bash
Usage:
  eventlist [-I <scvdFile>]... [-o <outputFile>] [-a <elf/axfFile>] [-b] <logFile>

Flags:
  -a <fileName>     elf/axf file name
  -b --begin        show statistic at beginning
  -f <txt/xml/json> output format, default: txt
  -h --help         show short help
  -I <fileName>     include SCVD file name
  -o <fileName>     output file name
  -s --statistic    show statistic only
  -V --version      show version info
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

    ```txt
    ☑️ Note:
        Make sure 'git' and 'bash' paths are listed under the PATH environment
        variable and set the git bash priority higher in the path.
    ```

  - **Linux/MacOS:**
    - GNU Bash (minimum recommended version **5.0.17**)

## Clone repository

Clone GitHub repository to create a local copy on your computer to make
it easier to develop and test. Cloning of the repository can be done by following
the below git command:

```bash
git clone git@github.com:ARM-software/CMSIS-View.git
```

## Build components

The steps below demonstrate how to build and create an executable:

- Go to eventlist directory
  - cd \<**root**\>/tools/eventlist
- Run the command to build an executable under `build` directory
  - `./make.sh build` : Build and generate executable for host OS & architecture in current directory.
  - `./make.sh build -arch <ARCH> -os <OS> -outdir <OUT directory>` : Build and generate executable for provided configs.\
    for e.g.

    ```bash
    ./make.sh build -arch amd64 -os darwin -outdir "Path/to/output/dir"
    ```

## Run Tests

One can directly run the tests from the command line.

- Go to eventlist directory
  - cd \<**root**\>/tools/eventlist
- Clean existing cache test results
  - go clean -cache
- Run command
  - `./make.sh test` : Run all tests.
  - `./make.sh test <PACKAGE>` : Run test related to the specified package.\
    for e.g.

    ```bash
    ./make.sh test eventlist/pkg/event
    ```

## Code coverage

Users can get coverage and generate code coverage report in HTML format

- Go to eventlist directory
    > cd \<**root**\>/tools/eventlist
- Run command
  - `./make.sh coverage`        : Run tests and show coverage info.\
  - `./make.sh coverage -html <FILE path>` : Run tests with coverage info and generate specified HTML coverage report.\

    for e.g.

    ```bash
    ./make.sh coverage -html cov/coverage.html
    ```

```txt
☑️ Note:
   for more usable commands, Use `./make.sh -h`.
```

## License

**eventlist** is licensed under Apache 2.0.

## Log File Format

The log file is expected to use the [Common Trace Format](https://diamon.org/ctf/#specification). The binary trace
stream layout is describes using the *Trace Stream Description Language* (TSDL) in
[eventlist.tsdl](docs/eventlist.tsdl).
