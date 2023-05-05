# Fault example (Cortex-M7)

This project is a simple **Fault** component example running on **Arm Cortex-M7** microcontroller
simulated by **Arm Virtual Hardware** simulator.

The application allows triggering of specific faults upon which the fault information is saved and system is reset.
When system restarts the fault information is output via the **Event Recorder**.

The fault information can also be inspected with **Component Viewer** in a debug session.

>Note: This example runs on the [Arm Virtual Hardware VHT_MPS2_Cortex-M7 model](https://arm-software.github.io/AVH/main/simulation/html/Using.html) and does not require any hardware.

## Prerequisites

### Software:
 - [**CMSIS-Toolbox v1.6.0**](https://github.com/Open-CMSIS-Pack/cmsis-toolbox/releases/tag/1.6.0) or newer
 - [**Keil MDK v5.38**](https://www.keil.com/mdk5) or newer containing:
   - Arm Compiler 6 (part of the MDK)
   - Arm Virtual Hardware (AVH) for MPS2 platform with Cortex-M7 (part of the MDK-Professional)
 - [**eventlist v1.1.0**](https://github.com/ARM-software/CMSIS-View/releases/tag/tools%2Feventlist%2F1.1.0) or newer

### CMSIS Packs:
 - Required packs:
    - ARM::CMSIS-View v1.2.0 or newer
    - ARM::CMSIS-Compiler v1.0.0 or newer
    - ARM::CMSIS v5.9.0
    - Keil::V2M-MPS2_CMx_BSP v1.8.0

   Missing packs can be installed by executing the following `csolution` and `cpackget` commands:
   ```
   csolution list packs -s Fault.csolution.yml -m >missing_packs_list.txt
   cpackget add -f missing_packs_list.txt
   ```
## Build

1. Use the `csolution` command to create `.cprj` project files:
   ```
   csolution convert -s Fault.csolution.yml
   ```

2. Use the `cbuild` command to create executable files:
   ```
   cbuild Fault.Debug+VHT_MPS2_Cortex-M7.cprj
   ```
## Run

### AVH Target

Execute the following steps:
 - run the AVH model from the command line by executing the following command:
   ```
   VHT_MPS2_Cortex-M7 -f vht_config.txt out/Fault/VHT_MPS2_Cortex-M7/Debug/Fault.axf
   ```
   >Note: The Arm Virtual Hardware executables have to be in the environment path, otherwise absolute path to the
          `VHT_MPS2_Cortex-M7.exe` (e.g. c:\Keil\ARM\VHT\VHT_MPS2_Cortex-M7) has to be provided instead of `VHT_MPS2_Cortex-M7`.

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

To show the events in a user friendly way copy `EventRecorder.scvd`, `RTX5.scvd` and `ARM_Fault.scvd` files from respective local Pack repository to the same folder where `EventRecorder.log` is located.

To process `EventRecorder.log` file with the `eventlist` utility execute the following command:
   ```
   eventlist -I EventRecorder.scvd -I RTX5.scvd -I ARM_Fault.scvd EventRecorder.log
   ```
