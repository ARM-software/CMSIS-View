# Event Statistics {#ev_stat}

## Overview {#about_event_statistics}

The \ref Event_Execution_Statistic functions allow you to collect and statistical data about the code execution. Any 
debug adapter can be used to record execution timing and number of calls for annotated code sections:

\image html EventStatistics_wo_Energy.png "Event Statistics for user code"

Energy profiling is of annotated code sections is possible using <a href="https://www2.keil.com/mdk5/ulink/ulinkplus">ULINKplus</a>. 
When combined with power measurement, the Event Statistics window displays the energy consumption of the code section with min/man/average values:

\image html EventStatistics_w_Energy.png "User code energy profiling"

For more information, refer to \ref es_use.

**Benefits of Event Statistics:**
 - Collect statistical data about the code execution (time and energy).
 - Log files enable comparisons between different build runs in continuous integration (CI) environments.
 - Improve overall code quality and energy profile (especially relevant for battery driven applications).

## Using Event Statistics{#es_use}

The following steps enable the MDK debugger views for \estatistics on timing, number of calls, and current consumption.

To use \estatistics in the application code:
  -# Follow the first two steps in \ref er_use.
  -# Annotate the C source with \ref Event_Execution_Statistic.

\ref Event_Execution_Statistic functions may be placed throughout the application source code to measure execution performance between 
corresponding start and stop events:

- <b>EventStart<i>G</i> (<i>slot</i>)</b> or <b>EventStart<i>G</i>v (<i>slot</i>, <i>val1</i>, <i>val2</i>)</b> functions define the start point of an execution slot.
- <b>EventStop<i>G</i> (<i>slot</i>)</b> or <b>EventStop<i>G</i>v (<i>slot</i>, <i>val1</i>, <i>val2</i>)</b> functions define the stop point of an execution slot.

The \estatistics window shows collected data about execution time, number of calls, and (when using <a href="https://www2.keil.com/mdk5/ulink/ulinkplus">ULINKplus</a>) the current consumption
for each execution slot. 

For the minimum and maximum time or current consumption it also shows for start and stop events the:
- C source file name and line number of event calls via <b>EventStart<i>G</i> (<i>slot</i>)</b> or <b>EventStop<i>G</i> (<i>slot</i>)</b>.
- Integer values <i>val1</i> and <i>val2</i> of event calls via <b>EventStart<i>G</i>v (<i>slot</i>, <i>val1</i>, <i>val2</i>)</b> or <b>EventStop<i>G</i>v (<i>slot</i>, <i>val1</i>, <i>val2</i>)</b>.

Each execution slot is identified by the function name group letter <i>G</i> = {A, B, C, D} and a <i>slot</i> number (0 to 15). \ref er_filtering may be used to control the recording of each group.
A call to <b>EventStop<i>G</i></b> or <b>EventStop<i>G</i>v</b>  with <i>slot</i>=15 stops measurement for all slots in a group and may be used at global exits of an execution block.

The following code is  from the \ref scvd_evt_stat example project that is part of the <b>Keil::ARM_Compiler</b> pack:

**Code example**

```
#define TABLE_SIZE 1000
float sin_table[TABLE_SIZE];
 
// Calculate table with sine values 
void CalcSinTable (void)  {
  unsigned int i, max_i;
  float f = 0.0;
  
  max_i = TABLE_SIZE - (rand () % 500);
  EventStartAv (15, max_i, 0);                  // Start group A, slot 15, passing the max_i variable
  for (i = 0; i < max_i; i++)  {
    if (i == 200)  {
       EventStartAv (0, max_i, 0);              // Start group A, slot 0, passing the max_i variable
    }
 
    sin_table[i] = sinf(f);
    f = f + (3.141592 / TABLE_SIZE);
  
    if (i == 800)  {                            // Measure 800 table entries
      EventStopA (0);                           // Stop group A, slot 0
    }
  }
 
  EventStopA (15);                              // Stop group A, slot 15 (stops also slots 0..14)
}
 
...
 
int main (void) {
 
  SystemCoreClockUpdate();                      // System Initialization
  
  EventRecorderInitialize(EventRecordAll, 1U);  // Initialize and start Event Recorder
 
  EventStartC (0);                              // start measurement event group C, slot 0
 
  for (j = 0; j < 1000; j++)  {
    CalcSinTable ();                            // calculate table with sinus values
  
    EventStartB(0);                             // start group B, slot 0
    MaxSqrtSum = rand () / 65536;               // limit for sqrt calculation
    num = FindSqrtSum ((float) MaxSqrtSum);     // return number of sqrt operations
    EventStopBv(0, MaxSqrtSum, num);            // stop group B, slot 0, output values: MaxSqrtSum, num
  }
 
  EventStopC(0);                                // stop measurement event group C, slot 0
  
  for (;;) {}
}
```

Build and run the example project which uses the µVision simulator (available in all versions of MDK). In a debug session,
\erecorder, displays the following output:

\image html er_with_statistics_annotated.png "Event Recorder with Start/Stop events and values"

The \estatistics window shows the statistical data about the code execution:

\image html es_start_stop_wo_energy_annotated.png "Event Statistics"

For more information on the usage of the functions, refer to the 
<a href="group__Event__Execution__Statistic.html"><b>Event Execution Statistics API</b></a>.

## Display current consumption{#es_display_energy}

Using a <a href="https://www2.keil.com/mdk5/ulink/ulinkplus">ULINKplus</a> debug adapter, you can also record and analyze the
energy that has been consumed in each execution slot. Using the above example on a hardware target with a ULINKplus, you get the
following displays in the \estatistics window (the \erecorder window does not change):

\image html es_start_stop_w_energy.png "Event Statistics displaying the energy consumption"
