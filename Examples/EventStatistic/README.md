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

The following command runs the VHT simulation model from the command line:  
`> VHT_MPS3_Corstone_SSE-300 -f vht_config.txt --simlimit=60 -C cpu0.semihosting-enable=1 .\EventStatistic.Debug+AVH_OutDir\EventStatistic.Debug+AVH.axf`

When using `cpu0.semihosting-enable=1` the file `EventRecorder.log` is generated that contains the events that are generated during execution.

This file can be analysed using the `EventList` utility with the following command:
`> eventlist EventRecorder.log`
