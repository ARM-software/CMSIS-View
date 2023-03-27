# CMSIS-View

The [**Keil MDK Event Recorder**](https://developer.arm.com/documentation/101407/latest/Debugging/Debug-Windows-and-Dialogs/Event-Recorder) is now available as open source component along with tooling.
This repository contains the source code of:
  - [**ARM::CMSIS-View**](https://arm-software.github.io/CMSIS-View/main/index.html) software pack that provides the event recorder software component.
  - [**EventList**](./tools/eventlist) command line utility that allows to dump the events on command line.
  - [**Example Projects**](./Examples) that show the usage of the Event Recorder.

## Repository toplevel structure

```
    📦
    ┣ 📂 .github          GitHub Action workflow and configuration
    ┣ 📂 Documentation    Target directory for generated documentation
    ┣ 📂 Doxygen          Source directory for documentation
    ┣ 📂 EventRecorder    Source code of EventRecorder software component
    ┣ 📂 Examples         Usage examples of the EventRecorder component
    ┣ 📂 tools            EventList command line tool source code
    ┗ 📂 Scripts          Helper scripts
```

## Generating Software Pack

Some helper scripts are provided to generate the release artifacts from this repository.

### Doxygen Documentation

Generating the HTML-formatted documentation from its Doxygen-based source is done via

```bash
CMSIS-View $ ./Doxygen/gen_doc.sh
``` 

Prerequisites for this script to succeed are:
 - Doxygen 1.9.2

### CMSIS-Pack Bundle

The CMSIS-Pack bundle can be generated with

```bash
CMSIS-View $ ./gen_pack.sh
``` 

Prerequisites for this script to succeed are:
 - Generated documentation (see above)
 - 7z
 - packchk
 - xmllint (optional)

### Version and Changelog Inference

The version and changelog embedded into the documentation and pack are inferred from the
local Git history. In order to get the full changelog one needs to have a full clone (not
a shallow one) including all release tags.

One can check the inferred version by calling the helper script:

```bash
CMSIS-View $ ./Scripts/git_describe.sh
Git version: '1.2.4-dev5+g6a00f35'
1.2.4-dev5+g6a00f35
^ ^ ^    ^ ^
| | |    | ┗ Commit SHA of the commit used for generation
| | |    ┗ Number of commits added since latest tag
| | ┗ Patch version from the latest tag incremented by one
| ┗ Minor version as in the latest tag
┗ Major version as in the latest tag
``` 

The full changelog can be inspected by calling the helper script:

```bash
CMSIS-View $ ./Scripts/gen_changelog.sh [-h] [-f <format>]
Generating changelog ...
Git version: '1.2.4-dev5+g6a00f35'
1.2.3 (2022-07-28)
- EventRecorder
  - Change 1
  - Change 2
- EventList
  - Change 1
  - Change 2
1.1.1 (2022-03-31)
Added clock frequency to internal Event Recorder Initialization message
1.1.0 (2022-02-25)
Added Event Recorder logging via Semihosting
1.0.0 (2021-01-28)
Add EventRecorder component
```

The version numbers are taken from the available tags. The shown release dates and
changelogs are one of:

1. For annotated tags the tagger date and the associated message is used.
2. For simple tags the committer date and message of the pointed-to commit is used.


## License

CMSIS-View is licensed under Apache 2.0.

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
