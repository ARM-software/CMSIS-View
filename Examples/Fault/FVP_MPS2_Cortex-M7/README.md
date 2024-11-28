# Fault example (Cortex-M7)

This project is a simple **Fault** component example running on an **Arm Cortex-M7** microcontroller simulated by
[**Arm Virtual Hardware**](https://www.arm.com/products/development-tools/simulation/virtual-hardware) with the **FVP_MPS2_Cortex-M7** model simulator.

The application allows triggering of specific faults upon which the fault information is saved and system is reset.
When system restarts the fault information is output via the **Event Recorder**.

The fault information can also be inspected with **Component Viewer** in a debug session.

> **Note**  
> This example runs on the [**Arm Virtual Hardware**](https://www.arm.com/products/development-tools/simulation/virtual-hardware) simulator
> and does not require any hardware.

## Prerequisites

### Software

- [**Arm Keil Studio for VS Code**](https://marketplace.visualstudio.com/items?itemName=Arm.keil-studio-pack)
- [**eventlist**](https://github.com/ARM-software/CMSIS-View/releases/tag/tools%2Feventlist%2F1.1.0) **v1.1.0** or newer

## Build and Run

To try the example with the **Arm Keil Studio**, follow the steps below:

 1. open the example in the **Visual Studio Code**.
 2. in the **Configure Solution** tab select the desired compiler (**AC6** or **GCC**), and click on the **OK** button.
 3. build the solution (in the **CMSIS** extension view click on the **Build solution** button).
 4. run the **FVP model** from the command line by executing the following command:
    - for **AC6**:
      ```shell
      FVP_MPS2_Cortex-M7 -f fvp_config.txt out/Fault/FVP_MPS2_Cortex-M7/Debug/Fault.axf
      ```
    - for **GCC**:
      ```shell
      FVP_MPS2_Cortex-M7 -f fvp_config.txt out/Fault/FVP_MPS2_Cortex-M7/Debug/Fault.elf
      ```

    > **Note**  
    > **The Arm Virtual Hardware executable files have to be in the environment path**.  
    > You can install **Arm Virtual Hardware** via **Arm Keil Studio** by following these steps:
    > - click on the **Arm Tools**.
    > - select **Add Arm Tools Configuration to Workspace**.
    > - under **Arm Virtual Hardware for Cortex®-M based on Fast Models** select the latest available version.
    > - save the **vcpkg-configuration.json** file.

 5. follow the instructions in the **simulator console**.
 6. the result of example running is an `EventRecorder.log` file that contains events that were generated during the code execution.
 7. view the events by using `eventlist` utility by executing the following command:
    ```shell
    eventlist -I %CMSIS_PACK_ROOT%/ARM/CMSIS-View/1.2.0/Fault/ARM_Fault.scvd -I %CMSIS_PACK_ROOT%/ARM/CMSIS-View/1.2.0/EventRecorder/EventRecorder.scvd -I %CMSIS_PACK_ROOT%/ARM/CMSIS-RTX/5.9.0/RTX5.scvd EventRecorder.log
    ```

    > **Note**  
    > If `CMSIS-View v1.2.0` or `CMSIS-RTX v5.9.0` packs are not installed, in the previous command replace the corresponding path with the path of the latest installed packs
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
