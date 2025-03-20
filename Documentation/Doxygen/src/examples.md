# Examples {#ExampleProjects}

CMSIS-View provides several example projects that demonstrate the usage of CMSIS-View software components in embedded applications. The projects are available in the `Examples/` directory and are also part of the CMSIS-View pack.

Example Project                                      | Description
:----------------------------------------------------|:-----------------------------------------
\subpage scvd_evt_stat "EventStatistic"              | Shows the usage of start/stop events for statistical code profiling using [Arm Virtual Hardware (AVH)](https://arm-software.github.io/AVH/main/overview/html/index.html) processor simulation.
\subpage flt_example_CM7 "Fault/FVP_MPS2_Cortex-M7"  | Shows the usage of the Fault component on a Cortex-M7 [Arm Virtual Hardware (AVH)](https://arm-software.github.io/AVH/main/overview/html/index.html) processor simulation.
\subpage flt_example_CM33 "Fault/B-U585I-IOT02A"     | Shows the usage of the Fault component on a Cortex-M33 with TrustZone. This example runs on STMicroelectronics [**B-U585I-IOT02A**](https://www.st.com/en/evaluation-tools/b-u585i-iot02a.html) evaluation board.

Other examples with CMSIS-View support can be also found in various CMSIS packs with device support, BSP and software components. For example [CMSIS-RTX](https://arm-software.github.io/CMSIS-RTX/latest/index.html), [MDK Middleware](https://www.keil.com/pack/doc/mw/General/html/index.html) and others, or also in GitHub projects at [https://github.com/Arm-Examples](https://github.com/Arm-Examples).

\page scvd_evt_stat EventStatistic

The EventStatistics example project shows how to use start/stop events with the Event Recorder and obtain execution statistics.

The start/stop events allow to measure execution times with different slots (0 - 15) in four different groups (A - D). The call to 'EventStart' starts a timer slot; the call to 'EventStop' stops the related timer. A call to EventStop with slot 15 stops the timers of all slots in the specific group.

This example project is located in the `Examples/EventStatistic/` directory. It uses the [**Arm Keil Studio Pack**](https://marketplace.visualstudio.com/items?itemName=Arm.keil-studio-pack)
for project build and runs on [Arm Virtual Hardware (AVH)](https://arm-software.github.io/AVH/main/overview/html/index.html) processor simulation and does not require real hardware board.

**Build and run**

Follow the instructions in the [example's README.md](https://github.com/ARM-software/CMSIS-View/blob/main/Examples/EventStatistic/README.md) to build and run the project.

Use the \ref evntlst to analyze the outcomes.

**main.c File**

The application code in the main.c file does some time consuming calculations that are recorded.

```c
#include "RTE_Components.h"             // Component selection
#include CMSIS_device_header            // defined in the RTE_Components.h file

#include "EventRecorder.h"

#include <math.h>
#include <stdio.h>
#include <stdlib.h>

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

// Return number of sqrt operations to exceed sum
unsigned int FindSqrtSum (float max_sum)  {
  unsigned int i;
  float sqrt_sum;

  sqrt_sum = 0.0;
  for (i = 0; i < 10000; i++) {
    sqrt_sum += sqrtf((float) i);
    if (sqrt_sum > max_sum)  {
      return (i);
    }
  }
  return (i);
}

unsigned int j, num, MaxSqrtSum;


int main (void) {

  SystemCoreClockUpdate();                      // System Initialization

  EventRecorderInitialize (EventRecordAll, 1U); // Initialize and start Event Recorder
  EventRecorderClockUpdate();
  EventStartC (0);                              // start measurement event group C, slot 0
  printf ("Started\n");
  for (j = 0; j < 10000; j++)  {
    CalcSinTable ();                            // calculate table with sinus values

    EventStartB(0);                             // start group B, slot 0
    MaxSqrtSum = rand () / 65536;               // limit for sqrt calculation
    num = FindSqrtSum ((float) MaxSqrtSum);     // return number of sqrt operations
    EventStopBv(0, MaxSqrtSum, num);            // stop group B, slot 0, output values: MaxSqrtSum, num
  }

  EventStopC(0);                                // stop measurement event group C, slot 0

  while (1) {
        __NOP();
  }
}
```

\page flt_example_CM7 Fault/FVP_MPS2_Cortex-M7

This example project shows \ref fault on an Arm Cortex-M7 based device, using [Arm Virtual Hardware (AVH)](https://arm-software.github.io/AVH/main/simulation/html/using.html) processor simulation as the target platform.

**Build and run**

Follow the instructions in the [example's README.md](https://github.com/ARM-software/CMSIS-View/blob/main/Examples/Fault/FVP_MPS2_Cortex-M7/README.md) to build and run the project.

**Screenshots**

User Interface:

![](./images/Fault_CM7_UI.png)

EventRecorder.log file containing saved Memory Management fault information, processed with eventlist utility:

![](./images/Fault_CM7_eventlist_MemManage.png)


\page flt_example_CM33 Fault/B-U585I-IOT02A

This example project shows \ref fault on an Arm Cortex-M33 based device. STMicroelectronics [B-U585I-IOT02A](https://www.st.com/en/evaluation-tools/b-u585i-iot02a.html) evaluation board is used as the target platform.

**Build and run**

Follow the instructions in the [example's README.md](https://github.com/ARM-software/CMSIS-View/blob/main/Examples/Fault/B-U585I-IOT02A/README.md) to build and run the project.

**Screenshots**

User Interface:

![](./images/Fault_CM33_UI.png)

Component View when no fault was saved yet:

![](./images/Fault_CM33_CV_NoFault.png)

Component View when Memory Management fault has occurred:

![](./images/Fault_CM33_CV_MemManage.png)

Event Recorder messages when Memory Management fault was decoded and output to Event Recorder:

![](./images/Fault_CM33_EvR_MemManage.png)

Serial Terminal output when Memory Management fault was decoded and output to STDIO:

![](./images/Fault_CM33_STDIO_MemManage.png)
