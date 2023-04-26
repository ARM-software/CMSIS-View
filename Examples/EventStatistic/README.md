# Event Statistic Example

This project shows how to use start/stop events with the Event Recorder that allow to measure execution times with:
-  different slots (0 - 15)
-  different groups (A - D)

The following API calls control this time recording:
- `EventStart` starts a timer slot.
- `EventStop` stops the related timer.
- `EventStop` with slot 15 stops the timers of all slots for the specified group.

Refer to [Using Event Statistics](https://arm-software.github.io/CMSIS-View/main/ev_stat.html#es_use) for more information.

This demo application does some time consuming calculations that are recorded
and can be displayed in the Event Statistics window.

>Note:
This example runs on Arm Virtual Hardware on the [VHT_MPS3_Corstone_SSE-300 model](https://arm-software.github.io/AVH/main/simulation/html/Using.html) and does not require any hardware.

## Prerequisites

Tools:

- [**CMSIS-Toolbox 1.5.0 or higher**](https://github.com/Open-CMSIS-Pack/cmsis-toolbox)
- [**Keil MDK 5.38 or higher**](https://www.keil.com/mdk5)
  - Arm Compiler 6 (part of MDK)
  - Arm Virtual Hardware for MPS3 platform with Corstone-300 (part of MDK-Professional)
- [**eventlist**](https://github.com/ARM-software/CMSIS-View/releases/latest) utility from this repository

As an alternative the example runs also on [**AMI Arm Virtual Hardware**](https://aws.amazon.com/marketplace/search/results?searchTerms=Arm+Virtual+Hardware) available via the AWS Marketplace as this image contains all relevant tools.

## Compile Project

The following commands convert and build the project with build type `Debug` and target type `AVH`:

```sh
EventStatistic $ cbuild .\EventStatistic.csolution.yml --configuration .Debug+AVH
```

## Execute

The following command runs the example for 60 seconds (parameter *--simlimit*) on the VHT simulation model:

```sh
EventStatistic $ VHT_MPS3_Corstone_SSE-300 -f vht_config.txt --simlimit=60 ./out/EventStatistic/Debug/AVH/Debug+AVH.axf
```

## Analyze Events

This file can be analyzed using the `eventlist` utility with the following command:

```sh
EventStatistic $ eventlist -s EventRecorder.log

   Start/Stop event statistic
   --------------------------

Event count      total       min         max         average     first       last
----- -----      -----       ---         ---         -------     -----       ----
A(0)  10000    31.44509s    1.69997ms   3.80041ms   3.14451ms   3.29962ms   3.59964ms
      Min: Start: 31.94980000 val1=0x000001f5, val2=0x00000000 Stop: 31.95149997 val1=0x10004d43, val2=0x0000003c
      Max: Start: 84.70757283 val1=0x000003a5, val2=0x00000000 Stop: 84.71137324 val1=0x10004d43, val2=0x00000038

A(15) 10000   169.75100s    2.49964ms  42.78648s   16.97510ms   3.99995ms   4.30004ms
      Min: Start: 81.87697318 val1=0x000001f8, val2=0x00000000 Stop: 81.87947282 val1=0x10004d43, val2=0x0000003c
      Max: Start: 37.41299154 val1=0x0000032f, val2=0x00000000 Stop: 80.19947314 val1=0x10004d43, val2=0x0000003c

B(0)  10000    10.83677s    0.00000s  169.29161ms   1.08368ms   1.60016ms   1.00010ms
      Min: Start: 1.76679986 val1=0x10004d43, val2=0x0000005c Stop: 1.76679986 val1=0x0000018e, val2=0x00000047
      Max: Start: 37.24369993 val1=0x10004d43, val2=0x0000005c Stop: 37.41299154 val1=0x000066bf, val2=0x00000487

C(0)      1   180.67372s  180.67372s  180.67372s  180.67372s  180.67372s  180.67372s
      Min: Start: 0.00000000 val1=0x10004d43, val2=0x00000057 Stop: 180.67371888 val1=0x10004d43, val2=0x00000062
      Max: Start: 0.00000000 val1=0x10004d43, val2=0x00000057 Stop: 180.67371888 val1=0x10004d43, val2=0x00000062
```

When adding the AXF file and the [SCVD file](https://arm-software.github.io/CMSIS-View/main/SCVD_Format.html) to the `eventlist` command the context of the program is shown

```sh
EventStatistic $ eventlist -a ./out/EventStatistic/Debug/AVH/Debug+AVH.axf -I $CMSIS_PACK_ROOT/ARM/CMSIS-View/1.2.0/EventRecorder/EventRecorder.scvd EventRecorder.log

  :

53947 180.66841874 EvCtrl    StartAv(15)             v1=776 v2=0
53948 180.66911914 EvCtrl    StartAv(0)              v1=776 v2=0
53949 180.67271878 EvCtrl    StopA(15)               File=./EventStatistic/main.c(60)
53950 180.67271878 EvCtrl    StartB(0)               File=./EventStatistic/main.c(92)
53951 180.67371888 EvCtrl    StopBv(0)               v1=15150 v2=802
53952 180.67371888 EvCtrl    StopC(0)                File=./EventStatistic/main.c(98)

   Start/Stop event statistic
   --------------------------

Event count      total       min         max         average     first       last
----- -----      -----       ---         ---         -------     -----       ----
A(0)  10000    31.44509s    1.69997ms   3.80041ms   3.14451ms   3.29962ms   3.59964ms
      Min: Start: 31.94980000 v1=501 v2=0 Stop: 31.95149997 File=./EventStatistic/main.c(60)
      Max: Start: 84.70757283 v1=933 v2=0 Stop: 84.71137324 File=./EventStatistic/main.c(56)

A(15) 10000   169.75100s    2.49964ms  42.78648s   16.97510ms   3.99995ms   4.30004ms
      Min: Start: 81.87697318 v1=504 v2=0 Stop: 81.87947282 File=./EventStatistic/main.c(60)
      Max: Start: 37.41299154 v1=815 v2=0 Stop: 80.19947314 File=./EventStatistic/main.c(60)

B(0)  10000    10.83677s    0.00000s  169.29161ms   1.08368ms   1.60016ms   1.00010ms
      Min: Start: 1.76679986 File=./EventStatistic/main.c(92) Stop: 1.76679986 v1=398 v2=71
      Max: Start: 37.24369993 File=./EventStatistic/main.c(92) Stop: 37.41299154 v1=26303 v2=1159

C(0)      1   180.67372s  180.67372s  180.67372s  180.67372s  180.67372s  180.67372s
      Min: Start: 0.00000000 File=./EventStatistic/main.c(87) Stop: 180.67371888 File=./EventStatistic/main.c(98)
      Max: Start: 0.00000000 File=./EventStatistic/main.c(87) Stop: 180.67371888 File=./EventStatistic/main.c(98)
```
