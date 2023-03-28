\mainpage

**CMSIS-View** equips software developers with methodologies, software components and utilities that provide visibility into internal operation of embedded applications and software components.

With the software components of CMSIS-View, developers can collect time-accurate event-based information, display program execution status, and analyze fault exceptions. The data can be observed in real-time in an IDE or can be saved as a log file during program execution. It allows to analyze execution flows, debug potential issues, and measure execution times.

A **Software Component Viewer Description** (*.SCVD) file in \ref SCVD_Format (XML) defines the content that is displayed in the **Component Viewer** and **Event Recorder**.

In addition, using the Event Recorder API, you can annotate your
code so that you can get statistical data on the time spent in a loop or on the energy consumption
(<a href="https://developer.arm.com/Tools%20and%20Software/ULINKplus">ULINKplus</a> required).

Key elements of CMSIS-View are:

- \ref er_use "Event Recorder" - is an embedded software component that provides an [API (function calls)](modules.html) for event annotations in the code.
- \ref SCVD_Format "SCVD file specification" defines the content that is displayed.
- \ref evntlst, a command line tool for processing Event Recorder log files.
- \ref fault analysis with infrastructure and functions to store, record, and analyze exception fault information.

# Content {#content}

This user's guide contains the following chapters:

- \subpage rev_hist : lists CMSIS-View releases
- \subpage evr : explores the features and operation of the **Event Recorder** including
  configuration, technical data, and theory of operation.
- \subpage ev_stat : describes how to use Event Statistics to create statistical data on code execution and power consumption.
- \subpage evntlst : shows the usage of `eventlist`, a command line tool for processing Event Recorder records stored to a log file.
- \subpage cmp_viewer : explains the use of Component Viewer.
- \subpage SCVD_Format : describes the format of the Software Component View Description (*.SCVD) files that define the output of the MDK debugger views.
- \subpage fault : infrastructure and functions to store, record, and analyze the Cortex-M Exception Fault information.
- \subpage ExampleProjects are available demonstrating standard use cases.
- **[API References](modules.html)** describes the API and the functions of the **Event Recorder** and **Fault components** in details.
