# Overview {#mainpage}

**CMSIS-View** equips software developers with software components, utilities and methodologies that provide visibility into internal operation of embedded applications and software components.

With the software components of CMSIS-View, developers can collect time-accurate event-based information, display program execution status, and analyze fault exceptions. It allows to analyze execution flows, debug potential issues, and measure execution times. The data can be observed in real-time in an IDE or can be saved as a log file during program execution. 

A **Software Component Viewer Description** (*.SCVD) file in \ref SCVD_Format (XML) defines the content that is displayed in the **Component Viewer** and **Event Recorder**.

In addition, using the Event Recorder API, you can annotate your code so that you can get statistical data on the time spent in a loop or on the energy consumption ([ULINKplus](https://developer.arm.com/Tools%20and%20Software/ULINKplus) required).

Key elements of CMSIS-View are:

 - \ref evr "Event Recorder" - an embedded software component that provides \ref Ref_EventRecorder for event annotations in the code.
 - \ref evntlst "eventlist utility" - a command line tool for processing Event Recorder log files.
 - \ref fault "Fault" - an embedded software component with infrastructure and \ref Ref_Fault to store, record, and analyze exception fault information.

> **Note**
> - CMSIS-View replaces and extends functionality previously provided as part of *Keil::ARM_Compiler* pack.
> - See [Migrating projects from CMSIS v5 to CMSIS v6](https://learn.arm.com/learning-paths/microcontrollers/project-migration-cmsis-v6) for a guidance on updating existing projects to CMSIS-View.

## Access to CMSIS-View {#view_access}

CMSIS-View is actively maintained in [**CMSIS-View GitHub repository**](https://github.com/ARM-software/CMSIS-View) and is released as a standalone [**CMSIS-View pack**](https://www.keil.arm.com/packs/cmsis-view-arm/versions/) in the [CMSIS-Pack format](https://open-cmsis-pack.github.io/Open-CMSIS-Pack-Spec/main/html/index.html).

The table below explains the content of **ARM::CMSIS-View** pack.

Directory                             | Description
:-------------------------------------|:------------------------------------------------------
📂 Documentation                      | Folder with this CMSIS-View documenation
📂 EventRecorder                      | \ref evr implementation
📂 Examples                           | \ref ExampleProjects "Examples projects" using CMSIS-View
📂 Fault                              | Implementation of the \ref fault "Fault component"
📄 ARM.CMSIS-View.pdsc                | Pack description file in CMSIS-Pack format
📄 LICENSE                            | License Agreement (Apache 2.0)

See [CMSIS Documentation](https://arm-software.github.io/CMSIS_6/) for an overview of CMSIS software components, tools and specifications.

## Documentation Structure {#doc_structure}

This user's guide contains the following chapters:

 - \ref rev_hist : lists CMSIS-View releases
 - \ref evr : explores the features and operation of the **Event Recorder** including configuration, technical data, and theory of operation.
 - \ref ev_stat : describes how to use Event Statistics to create statistical data on code execution and power consumption.
 - \ref evntlst : shows the usage of `eventlist`, a command line tool for processing Event Recorder records stored to a log file.
 - \ref cmp_viewer : explains the use of Component Viewer.
 - \ref SCVD_Format : describes the format of the Software Component View Description (*.SCVD) files that define the content that is displayed.
 - \ref fault : infrastructure and functions to store, record, and analyze the Cortex-M Exception Fault information.
 - \ref ExampleProjects are available demonstrating standard use cases.
 - [**API References**](topics.html) describes the API and the functions of the **Event Recorder** and **Fault** components in details.

## License {#license}

CMSIS-View is provided free of charge by Arm under the [Apache 2.0 License](https://raw.githubusercontent.com/ARM-software/CMSIS-View/main/LICENSE).
