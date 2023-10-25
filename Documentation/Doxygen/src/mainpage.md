# Overview {#mainpage}

**CMSIS-View** offers software developers methodologies, software components, and utilities that provide visibility into internal operation of embedded applications and software components.

With the software components of CMSIS-View, developers can collect time-accurate event-based information, display program execution status, and analyze fault exceptions. It allows to analyze execution flows, debug potential issues, and measure execution times. The data can be observed in real-time in an IDE or can be saved as a log file during program execution. 

A **Software Component Viewer Description** (*.SCVD) file in \ref SCVD_Format (XML) defines the content that is displayed in the **Component Viewer** and **Event Recorder**.

In addition, using the Event Recorder API, you can annotate your
code so that you can get statistical data on the time spent in a loop or on the energy consumption
([ULINKplus](https://developer.arm.com/Tools%20and%20Software/ULINKplus) required).

Key elements of CMSIS-View are:

- \ref er_use "Event Recorder" - is an embedded software component that provides an [API (function calls)](modules.html) for event annotations in the code.
- \ref evntlst, a command line tool for processing Event Recorder log files.
- \ref fault with infrastructure and functions to store, record, and analyze exception fault information.

## Access to CMSIS-View {#view_access}

CMSIS-View is actively maintained in [**CMSIS-View GitHub repository**](https://github.com/ARM-software/CMSIS-View) that contains the full source of CMSIS-View firmware, implementation of eventlist utility, examples, as well as this documentation.

CMSIS-View software components are released in [CMSIS-Pack format](https://open-cmsis-pack.github.io/Open-CMSIS-Pack-Spec/main/html/index.html). An overview of the pack and downloads are available on [CMSIS-View Pack page](https://www.keil.arm.com/packs/cmsis-view-arm/versions/).

The table below explains the content of **ARM::CMSIS-View** pack.

Directory                             | Description
:-------------------------------------|:------------------------------------------------------
📂 Documentation                      | Folder with this CMSIS-View documenation
📂 EventRecorder                      | \ref evr implementation
📂 Examples                           | \ref ExampleProjects "Examples projects" using CMSIS-View
📂 Fault                              | Implementation of the \ref fault "Fault component"
📄 LICENSE                            | License Agreement (Apache 2.0)
📄 ARM.CMSIS-View.pdsc                | Pack description file in CMSIS-Pack format

See [CMSIS Documentation](https://arm-software.github.io/CMSIS_6/) for an overview of CMSIS software components, tools and specifications.

## Documentation Content {#content}

This user's guide contains the following chapters:

- \ref rev_hist : lists CMSIS-View releases
- \ref evr : explores the features and operation of the **Event Recorder** including
  configuration, technical data, and theory of operation.
- \ref ev_stat : describes how to use Event Statistics to create statistical data on code execution and power consumption.
- \ref evntlst : shows the usage of `eventlist`, a command line tool for processing Event Recorder records stored to a log file.
- \ref cmp_viewer : explains the use of Component Viewer.
- \ref SCVD_Format : describes the format of the Software Component View Description (*.SCVD) files that define the content that is displayed.
- \ref fault : infrastructure and functions to store, record, and analyze the Cortex-M Exception Fault information.
- \ref ExampleProjects are available demonstrating standard use cases.
- **[API References](modules.html)** describes the API and the functions of the **Event Recorder** and **Fault components** in details.

## License {#license}

CMSIS-View is provided free of charge by Arm under the [Apache 2.0 License](https://raw.githubusercontent.com/ARM-software/CMSIS-View/main/LICENSE).
