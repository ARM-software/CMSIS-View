# Fault example (Cortex-M7)

This project is a simple **Fault** component example running on an **Arm Cortex-M7** microcontroller simulated by
[**Arm Virtual Hardware**](https://www.arm.com/products/development-tools/simulation/virtual-hardware) with the **FVP_MPS2_Cortex-M7** model simulator.

The application allows triggering of specific faults upon which the fault information is saved and system is reset.
When system restarts the fault information is output via the **Event Recorder**.

The fault information can also be inspected with **Component Viewer** in a debug session.

> **Note**  
> This example runs on the **Arm Virtual Hardware** simulator and does not require any hardware.

## Prerequisites

### Software

- [**Arm Keil Studio Pack**](https://marketplace.visualstudio.com/items?itemName=Arm.keil-studio-pack)
- [**CMSIS-Toolbox**](https://github.com/Open-CMSIS-Pack/cmsis-toolbox/releases) **v2.6.0** or newer
- [**Keil MDK**](https://developer.arm.com/Tools%20and%20Software/Keil%20MDK) **v5.41** or newer
- [**eventlist**](https://github.com/ARM-software/CMSIS-View/releases/tag/tools%2Feventlist%2F1.1.0) **v1.1.0** or newer

### CMSIS Packs

- Required packs:
  - [ARM::CMSIS-View](https://www.keil.arm.com/packs/cmsis-view-arm/versions/) **v1.2.0** or newer
  - [ARM::CMSIS](https://www.keil.arm.com/packs/cmsis-arm/overview/) **v6.1.0** or newer
  - [ARM::CMSIS-RTX](https://www.keil.arm.com/packs/cmsis-rtx-arm/versions/) **v5.9.0** or newer
  - [ARM::CMSIS-Compiler](https://www.keil.arm.com/packs/cmsis-compiler-arm/versions/) **v2.1.0** or newer
  - [Keil::V2M-MPS2_CMx_BSP](https://www.keil.arm.com/packs/v2m-mps2_cmx_bsp-keil/versions/) **v1.8.0** or newer

## Build and Run

### Arm Keil Studio

#### Compiler: Arm Compiler 6

To try the example with the **Arm Keil Studio**, do the following steps:

 1. open the **Visual Studio Code**.
 2. click on the **CMSIS** extension, click on the **Create a New Solution** button, then under **Create new solution** for
    **Target Board (Optional)** select **V2M-MPS2 (B)**, under **Templates, Reference Applications, and Examples**
    look for and select the **Fault example**, choose the desired **Solution Location** and click on the **Create** button.
 3. in the **Configure Solution** tab select **AC6** compiler and click on the **OK** button.
 4. build the solution (click on the **hammer** button).
 5. run the AVH model from the command line by executing the following command:

    ```shell
    FVP_MPS2_Cortex-M7 -f fvp_config.txt out/Fault/FVP_MPS2_Cortex-M7/Debug/Fault.axf
    ```

    > **Note**  
    > The Arm Virtual Hardware executable files have to be in the environment path, otherwise executable file has to be started from
    > absolute path e.g. `C:\Keil_v5\ARM\avh-fvp\bin\models\FVP_MPS2_Cortex-M7.exe` has to be used instead of `FVP_MPS2_Cortex-M7`.
 6. follow the instructions in the **simulator console**.
 7. the result of example running is an `EventRecorder.log` file that contains events that were generated during the code execution.
 8. view the events by using `eventlist` utility by executing the following command:

    ```shell
    eventlist -I %CMSIS_PACK_ROOT%/ARM/CMSIS-View/1.2.0/Fault/ARM_Fault.scvd -I %CMSIS_PACK_ROOT%/ARM/CMSIS-View/1.2.0/EventRecorder/EventRecorder.scvd -I %CMSIS_PACK_ROOT%/ARM/CMSIS-RTX/5.9.0/RTX5.scvd EventRecorder.log
    ```

    > **Note**  
    > If `CMSIS-View v1.2.0` or `CMSIS-RTX v5.9.0` packs are not installed, in the previous command replace corresponding path with the path of the latest installed packs
    > (for example replace `%CMSIS_PACK_ROOT%/ARM/CMSIS-View/1.2.0/Fault/` with `%CMSIS_PACK_ROOT%/ARM/CMSIS-View/1.2.1/Fault/`)

#### Compiler: GCC

To try the example with the **Arm Keil Studio**, do the following steps:

 1. open the **Visual Studio Code**.
 2. click on the **CMSIS** extension, click on the **Create a New Solution** button, then under **Create new solution** for
    **Target Board (Optional)** select **V2M-MPS2 (B)**, under **Templates, Reference Applications, and Examples**
    look for and select the **Fault example**, choose the desired **Solution Location** and click on the **Create** button.
 3. in the **Configure Solution** tab select **GCC** compiler and click on the **OK** button.
 4. build the solution (click on the **hammer** button).
 5. run the AVH model from the command line by executing the following command:

    ```shell
    FVP_MPS2_Cortex-M7 -f fvp_config.txt out/Fault/FVP_MPS2_Cortex-M7/Debug/Fault.elf
    ```

    > **Note**  
    > The Arm Virtual Hardware executable files have to be in the environment path, otherwise executable file has to be started from
    > absolute path e.g. `C:\Keil_v5\ARM\avh-fvp\bin\models\FVP_MPS2_Cortex-M7.exe` has to be used instead of `FVP_MPS2_Cortex-M7`.
 6. follow the instructions in the **simulator console**.
 7. the result of example running is an `EventRecorder.log` file that contains events that were generated during the code execution.
 8. view the events by using `eventlist` utility by executing the following command:

    ```shell
    eventlist -I %CMSIS_PACK_ROOT%/ARM/CMSIS-View/1.2.0/Fault/ARM_Fault.scvd -I %CMSIS_PACK_ROOT%/ARM/CMSIS-View/1.2.0/EventRecorder/EventRecorder.scvd -I %CMSIS_PACK_ROOT%/ARM/CMSIS-RTX/5.9.0/RTX5.scvd EventRecorder.log
    ```

    > **Note**  
    > If `CMSIS-View v1.2.0` or `CMSIS-RTX v5.9.0` packs are not installed, in the previous command replace corresponding path with the path of the latest installed packs
    > (for example replace `%CMSIS_PACK_ROOT%/ARM/CMSIS-View/1.2.0/Fault/` with `%CMSIS_PACK_ROOT%/ARM/CMSIS-View/1.2.1/Fault/`)

## User Interface

This example uses **simulator console** for User Interface.

The fault triggering is done by entering a number via **simulator console**, see possible values below:

- `0` : Terminate the example
- `1` : trigger the Data access (precise) Memory Management fault
- `2` : trigger the Data access (precise) Bus fault
- `3` : trigger the Data access (imprecise) Bus fault
- `4` : trigger the Instruction execution Bus fault
- `5` : trigger the Undefined instruction Usage fault
- `6` : trigger the Divide by 0 Usage fault
