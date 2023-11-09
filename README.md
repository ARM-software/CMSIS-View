# CMSIS-View

**CMSIS-View** provides software components and utilities that allow embedded software developers to analyze program execution flows, debug potential issues, and measure code execution times. The data can be observed in real-time in an IDE or can be saved as a log file during program execution.

This repository contains the source code of:

- [**ARM::CMSIS-View pack**](https://www.keil.arm.com/packs/cmsis-view-arm) that provides the Event Recorder and Fault software components.
- [**EventList**](./tools/eventlist) command line utility that allows to dump the events on command line.
- [**Example Projects**](./Examples) that show the usage of the Event Recorder and Fault components.

[CMSIS-View documentation](https://arm-software.github.io/CMSIS-View) explains available functionality and APIs.

> **Note**
> - CMSIS-View replaces and extends functionality previously provided as part of *Keil::ARM_Compiler* pack.
> - See [Migrating projects from CMSIS v5 to CMSIS v6](https://learn.arm.com/learning-paths/microcontrollers/project-migration-cmsis-v6) for a guidance on updating existing projects to CMSIS-View.

## Repository toplevel structure

```txt
    📦
    ┣ 📂 .github          GitHub Action workflow and configuration
    ┣ 📂 Documentation    Documentation directory
    ┣ 📂 EventRecorder    Source code of EventRecorder software component
    ┣ 📂 Examples         Usage examples
    ┣ 📂 Fault            Source code of Fault software component
    ┣ 📂 Schema           Schema files
    ┗ 📂 tools            EventList command line tool source code
```

## Generating Software Pack

Some helper scripts are provided to generate the release artifacts from this repository.

### Doxygen Documentation

Generating the HTML-formatted documentation from its Doxygen-based source is done via

```sh
CMSIS-View $ ./Documentation/Doxygen/gen_doc.sh
```

Prerequisites for this script to succeed are:

- Doxygen 1.9.6

Also see [Documentation README](./documentation/README.md).

### CMSIS-Pack Bundle

The CMSIS-Pack bundle can be generated with

```sh
CMSIS-View $ ./gen_pack.sh
```

Prerequisites for this script to succeed are:

- Generated documentation (see above)
- 7z/GNU Zip
- packchk (e.g., via CMSIS-Toolbox)
- xmllint (optional)

### Version and Changelog Inference

The version and changelog embedded into the documentation and pack are inferred from the
local Git history. In order to get the full changelog one needs to have a full clone (not
a shallow one) including all release tags.

The version numbers and change logs are taken from the available annotated tags.

### Release Pack

A release is simply done via the GitHub Web UI. The newly created tag needs to have
the pattern `pack/<version>` where `<version>` shall be the SemVer `<major>.<minor>.<patch>`
version string for the release. The release description is used as the change log
message for the release.

When using an auto-generated tag (via Web UI) the release description is used as the
annotation message for the generated tag. Alternatively, one can prepare the release
tag in the local clone and add the annotation message independently from creating the
release.

Once the release is published via the GitHub Web UI the release workflow generates the
documentation and the pack (see above) and attaches the resulting pack archive as an
additional asset to the release.

## EventList Utility

The command line utility to decode EventRecorder log files written in Go.

### Compile and Test

To build and EventList run `make.sh` script.

```sh
CMSIS-View/tools/eventlist $ ./make.sh build
GOOS=windows GOARCH=amd64 go build -v -ldflags '-X "..."' -o ./eventlist.exe ./cmd/eventlist

build finished successfully!

CMSIS-View/tools/eventlist $ ./make.sh test
go test ./...
?       eventlist/cmd/make      [no test files]
ok      eventlist/cmd/eventlist 7.584s
ok      eventlist/pkg/elf       6.802s
ok      eventlist/pkg/eval      7.458s
ok      eventlist/pkg/event     7.471s
ok      eventlist/pkg/output    7.645s
ok      eventlist/pkg/xml/scvd  6.808s
```

One can run cross-builds for other than the own host platform by specifying `-arch <arch>`
and/or `-os <os>` on the `make.sh` command line, see `--help` for details.

### Release

A release for EventList utility is done independently from the CMSIS-View pack via
the GitHub Web UI. The release tag must match the pattern `tools/eventlist/<version>`
where `<version>` is a SemVer string.

The GitHub Action release workflow is triggered once a release is published. The
workflow builds release bundles for all supported target platforms and attaches
them as assets to the release.

## License Terms

CMSIS-View is licensed under [Apache License 2.0](LICENSE).

### Note

Individual files contain the following tag instead of the full license text.

SPDX-License-Identifier: Apache-2.0

This enables machine processing of license information based on the SPDX License Identifiers that are here available: http://spdx.org/licenses/

### External Dependencies

The components listed below are not redistributed with the project but are used internally for building, development, or testing purposes.

| Component | Version | License | Origin | Usage |
| --------- | ------- | ------- | ------ | ----- |
|goversioninfo|n/a|[MIT](https://opensource.org/licenses/MIT)|https://github.com/josephspurrier/goversioninfo| Used in [eventlist](./tools/eventlist) to generate MS Windows version info |

## Contributions and Pull Requests

Contributions are accepted under Apache 2.0. Only submit contributions where you have authored all of the code.

### Issues, Labels

Please feel free to raise an issue on GitHub
to report misbehavior (i.e. bugs)

Issues are your best way to interact directly with the maintenance team and the community.
We encourage you to append implementation suggestions as this helps to decrease the
workload of the very limited maintenance team.

We shall be monitoring and responding to issues as best we can.
Please attempt to avoid filing duplicates of open or closed items when possible.
In the spirit of openness we shall be tagging issues with the following:

- **bug** – We consider this issue to be a bug that shall be investigated.

- **wontfix** - We appreciate this issue but decided not to change the current behavior.

- **out-of-scope** - We consider this issue loosely related to CMSIS. It might be implemented outside of CMSIS. Let us know about your work.

- **question** – We have further questions about this issue. Please review and provide feedback.

- **documentation** - This issue is a documentation flaw that shall be improved in the future.

- **DONE** - We consider this issue as resolved - please review and close it. In case of no further activity, these issues shall be closed after a week.

- **duplicate** - This issue is already addressed elsewhere, see a comment with provided references.
