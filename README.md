# CMSIS-View

Repository of CMSIS Software Pack for software event recorder and input/output redirection.

# Repository toplevel structure

    📦
    ┣ 📂Doc            scripts to generate documentation for Eventrecorder
    ┣ 📂EventRecorder  source code for Eventrecorder
    ┗ 📂tools          command line tool source code

## Command line tools

For processing of event recorder records, the following command tool is provided:

- [**eventlist**](./tools/eventlist) - Process event recoder records.

**Refer to:**
  - [README.md](./tools/eventlist/README.md) for eventlist usage and build instructions.

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
