# Fault example (Cortex-M7)

This project is a simple **Fault** component example running on **Arm Cortex-M7** microcontroller simulated by 
[**Arm Virtual Hardware**](https://arm-software.github.io/AVH/main/simulation/html/Using.html) with the **VHT_MPS2_Cortex-M7** model simulator.

The application allows triggering of specific faults upon which the fault information is saved and system is reset.
When system restarts the fault information is output via the **Event Recorder**.

The fault information can also be inspected with **Component Viewer** in a debug session.

> **Note**
> - This example runs on the **Arm Virtual Hardware** simulator and does not require any hardware.

## Prerequisites

### Software

 - [**CMSIS-Toolbox v2.0.0**](https://github.com/Open-CMSIS-Pack/cmsis-toolbox/releases) or newer
 - [**Keil MDK v5.38**](https://www.keil.com/mdk5) or newer containing:
   - Arm Compiler 6 (part of the MDK)
   - Arm Virtual Hardware (AVH) for MPS2 platform with Cortex-M7 (part of the MDK-Professional)
 - [**eventlist v1.1.0**](https://github.com/ARM-software/CMSIS-View/releases/tag/tools%2Feventlist%2F1.1.0) or newer
 - [**Arm GNU Toolchain v12.3.Rel1**](https://developer.arm.com/downloads/-/arm-gnu-toolchain-downloads)
   (only necessary when building example with GCC)

### CMSIS Packs

 - Required packs:
    - ARM::CMSIS-View
    - ARM::CMSIS v6.0.0 or newer
    - ARM::CMSIS-RTX v1.0.0 or newer
    - ARM::CMSIS-Compiler v2.0.0 or newer
    - Keil::V2M-MPS2_CMx_BSP v1.8.0

Missing packs can be installed by executing the following `csolution` and `cpackget` commands:

```
csolution list packs -s Fault.csolution.yml -m >missing_packs_list.txt
cpackget add -f missing_packs_list.txt
```

## Build

1. Use the `csolution` command to create `.cprj` project files (for **Arm Compiler 6** toolchain):
   ```
   csolution convert -s Fault.csolution.yml
   ```
   or, for **GCC** toolchain use the following command:
   ```
   csolution convert -s Fault.csolution.yml -t GCC
   ```

2. Use the `cbuild` command to create executable files:
   ```
   cbuild Fault.Debug+VHT_MPS2_Cortex-M7.cprj
   ```

## Run

### AVH Target

Execute the following steps:
 - run the AVH model (with example built with **Arm Compiler 6** toolchain) from the command line by executing the following command:
   ```
   VHT_MPS2_Cortex-M7 -f vht_config.txt out/Fault/VHT_MPS2_Cortex-M7/Debug/Fault.axf
   ```
   or, run the AVH model (with example built with **GCC** toolchain) from the command line by executing the following command:
   ```
   VHT_MPS2_Cortex-M7 -f vht_config.txt out/Fault/VHT_MPS2_Cortex-M7/Debug/Fault.elf
   ```
   > **Note:** The Arm Virtual Hardware executables have to be in the environment path, otherwise absolute path to the 
   `VHT_MPS2_Cortex-M7.exe` (e.g. `c:\Keil\ARM\VHT\VHT_MPS2_Cortex-M7`) has to be provided instead of `VHT_MPS2_Cortex-M7`.

   The generated file `EventRecorder.log` contains the events that were generated during the example execution.
   This file is the input for the `eventlist` utility which can be used for further analysis.

 - follow the instructions in the simulator console

The fault triggering is done by entering a number via simulator console (see possible values below).

 - 0: terminate the example
 - 1: trigger the data access (precise) Memory Management fault
 - 2: trigger the data access (precise) Bus fault
 - 3: trigger the data access (imprecise) Bus fault
 - 4: trigger the instruction execution Bus fault
 - 5: trigger the undefined instruction Usage fault
 - 6: trigger the divide by 0 Usage fault

### Running the example in the uVision

 - open `Fault.Debug+VHT_MPS2_Cortex-M7.cprj` with the **uVision**
 - open **Options for Target**
 - select **Debug** tab
 - under **Use** select **Models Cortex-M Debugger** and click on **Settings** button end enter the following:
   - Command: `$K\ARM\VHT\VHT_MPS2_Cortex-M7.exe`
   - Target: `armcortexm7ct`
   - Configuration File: `.\vht_config.txt`
 - under **Initialization File:** enter the following: `.\debug.ini`
 - start the **Debug Session**
 - **Run** the code
 - you can see the events recorded by opening **View** - **Analysis Windows** - **Event Recorder**
 - you can see the state of the recorded fault by opening **View** - **Watch Windows** - **Fault**

## Analysis

To analyze the result `eventlist` utility is needed, copy the executable `eventlist` file to the same folder where `EventRecorder.log` is located.

To process `EventRecorder.log` file with the `eventlist` utility in **Windows Command Prompt** (cmd.exe) execute the following command:
```
eventlist -I %CMSIS_PACK_ROOT%/ARM/CMSIS-View/1.0.0/Fault/ARM_Fault.scvd -I %CMSIS_PACK_ROOT%/ARM/CMSIS-View/1.0.0/EventRecorder/EventRecorder.scvd -I %CMSIS_PACK_ROOT%/ARM/CMSIS-RTX/1.0.0/RTX5.scvd EventRecorder.log
```

> **Note**
> - If CMSIS-View v1.0.0 or CMSIS-RTX v1.0.0 packs are not installed, in the previous command replace corresponding path with the path of the latest installed packs
 (for example replace `%CMSIS_PACK_ROOT%/ARM/CMSIS-View/1.0.0/Fault/` with `%CMSIS_PACK_ROOT%/ARM/CMSIS-View/1.0.1/Fault/`)
