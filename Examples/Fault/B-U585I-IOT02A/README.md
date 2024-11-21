# Fault example (Cortex-M33)

This project is a simple **Fault** component example running on an **Arm Cortex-M33** microcontroller
on a STMicroelectronics [**B-U585I-IOT02A**](https://www.st.com/en/evaluation-tools/b-u585i-iot02a.html) evaluation board.

The application allows triggering of specific faults upon which the fault information is saved and system is reset.
When system restarts the fault information is output via the **Event Recorder** and via the **STDIO (STLink Virtual COM Port)**.

The fault information can also be inspected with **Component Viewer** in a debug session.

## Prerequisites

### Software

- [**Arm Keil Studio Pack**](https://marketplace.visualstudio.com/items?itemName=Arm.keil-studio-pack)
- [**CMSIS-Toolbox**](https://github.com/Open-CMSIS-Pack/cmsis-toolbox/releases) **v2.6.0** or newer
- [**STM32CubeMX**](https://www.st.com/en/development-tools/stm32cubemx.html) **v6.12.1** or newer
- [**STM32CubeProgrammer**](https://www.st.com/en/development-tools/stm32cubeprog.html) utility
- [**Keil MDK**](https://developer.arm.com/Tools%20and%20Software/Keil%20MDK) **v5.41** or newer

### CMSIS Packs

- Required packs:
  - [ARM::CMSIS-View](https://www.keil.arm.com/packs/cmsis-view-arm/versions/) **v1.2.0** or newer
  - [ARM::CMSIS](https://www.keil.arm.com/packs/cmsis-arm/overview/) **v6.1.0** or newer
  - [ARM::CMSIS-RTX](https://www.keil.arm.com/packs/cmsis-rtx-arm/versions/) **v5.9.0** or newer
  - [ARM::CMSIS-Compiler](https://www.keil.arm.com/packs/cmsis-compiler-arm/versions/) **v2.1.0** or newer
  - [ARM::CMSIS-Driver_STM32](https://www.keil.arm.com/packs/cmsis-driver_stm32-arm/overview/) **v1.0.0** or newer
  - [Keil::STM32U5xx_DFP](https://www.keil.arm.com/packs/stm32u5xx_dfp-keil/overview/) **v3.0.0** or newer
  - [Keil::B-U585I-IOT02A_BSP](https://www.keil.arm.com/packs/b-u585i-iot02a_bsp-keil/overview/) **v2.0.0** or newer

### Hardware

This board has to be properly configured with **TrustZone** enabled. Please follow the steps below carefully:

Configure the following **Option bytes** with the **STM32CubeProgrammer** utility:

- **User Configuration**:
  - `TZEN`: checked
  - `DBANK`: checked
- **Boot Configuration**:
  - `SECBOOTADD0`:  Value = 0x1800 Address = 0x0c000000
- **Secure Area 1**:
  - `SECWM1_PSTRT`: Value = 0x0    Address = 0x08000000
  - `SECWM1_PEND`:  Value = 0x7f   Address = 0x080fe000
- **Write Protection 1**:
  - `WRP1A_PSTRT`:  Value = 0x7f   Address = 0x080fe000
  - `WRP1A_PEND`:   Value = 0x0    Address = 0x08000000
- **Secure Area 1**:
  - `SECWM2_PSTRT`: Value = 0x7f   Address = 0x08000000
  - `SECWM2_PEND`:  Value = 0x0    Address = 0x08100000
- **Write Protection 1**:
  - `WRP2A_PSTRT`:  Value = 0x7f   Address = 0x081fe000
  - `WRP2A_PEND`:   Value = 0x0    Address = 0x08100000

## Build and Run

### Arm Keil Studio

#### Compiler: Arm Compiler 6

To try the example with the **Arm Keil Studio**, do the following steps:

 1. open the **Visual Studio Code**.
 2. click on the **CMSIS** extension, click on the **Create a New Solution** button, then under **Create new solution** for
    **Target Board (Optional)** select **B-U585I-IOT02A (Rev.C)**, under **Templates, Reference Applications, and Examples**
    look for and select the **Fault example**, choose the desired **Solution Location** and click on the **Create** button.
 3. in the **Configure Solution** tab select **AC6** compiler and click on the **OK** button.
 4. build the solution (click on the **hammer** button).
 5. download the built applications (Fault_S/Fault_NS) to the MCU's Flash (click on the **Run** button).  
    In case of issues use **uVision** for downloading to flash (just open **Fault.csolution.yml** in the **uVision**, click on the **Rebuild** button,
    in **Options for Target** dialog under **Debug** select **ST-Link Debugger** and click on the **OK** button; now click on the **Download** button).
 6. open the **Serial Terminal** application and connect to the **STMicroelectronics STLink Virtual COM Port (COMx)** (115200-8-N-1).
 7. press the **RESET** button on the board.
 8. follow the instructions in the **Serial Terminal** and observe the results.

#### Compiler: GCC

To try the example with the **Arm Keil Studio**, do the following steps:

 1. open the **Visual Studio Code**.
 2. click on the **CMSIS** extension, click on the **Create a New Solution** button, then under **Create new solution** for
    **Target Board (Optional)** select **B-U585I-IOT02A (Rev.C)**, under **Templates, Reference Applications, and Examples**
    look for and select the **Fault example**, choose the desired **Solution Location** and click on the **Create** button.
 3. in the **Configure Solution** tab select **AC6** compiler.
 4. in the **Fault.csolution.yml** csolution file replace `compiler: AC6` with `compiler: GCC`.
 5. launch the **STM32CubeMX generator** by selecting **Device:CubeMX** component in the project view and clicking on the **cog** button.
 6. in the **STM32CubeMX**: open from the menu **Project Manager - Project - Toolchain/IDE**,
    select **STM32CubeIDE** and clear **Generate Under Root** check box, and click on the **GENERATE CODE** button.
 7. clean the solution (click on the **...** button and select **Clean all out and tmp directories**).
 8. build the solution (click on the **hammer** button).
 9. download the built applications (Fault_S/Fault_NS) to the MCU's Flash (click on the **Run** button).  
    In case of issues use **uVision** for downloading to flash (just open **Fault.csolution.yml** in the **uVision**,
    in **Options for Target** dialog under **Debug** select **ST-Link Debugger** and click on the **OK** button; now click on the **Download** button).
 10. open the **Serial Terminal** application and connect to the **STMicroelectronics STLink Virtual COM Port (COMx)** (115200-8-N-1).
 11. press the **RESET** button on the board.
 12. follow the instructions in the **Serial Terminal** and observe the results.

### uVision

To try the example with the **uVision**, do the following steps:

 1. open the **uVision**.
 2. start the **Pack Installer** and under **Boards** tab select **B-U585I-IOT02A (Rev.C)**, then under **Examples** tab
    for **Fault example (B-U585I-IOT02A)** example click on the **Copy** button, choose the desired **Destination Folder**, and
    select **Use Pack Folder Structure** check box, select also **Launch uVision** check box and click on the **OK** button.
 3. build the **Project: Fault_S** project (click on the **Rebuild** button).
 4. build the **Project: Fault_NS** project (set **Fault_NS** as active project and click on the **Rebuild** button).
 5. download the built applications (Fault_S/Fault_NS) to the MCU's Flash (click on the **Download** button).
 6. open the **Serial Terminal** application and connect to the **STMicroelectronics STLink Virtual COM Port (COMx)** (115200-8-N-1).
 7. press the **RESET** button on the board.
 8. follow the instructions in the **Serial Terminal** and observe the results.

> **Note**  
> In the debug session fault information can be inspected in the **Component View** and **Event Recorder** windows

## User Interface

This example uses **Serial Terminal** for User Interface.

The fault triggering is done by entering a number via **Serial Terminal** application, see possible values below:

- `1` : trigger the Non-Secure fault, Non-Secure data access Memory Management fault
- `2` : trigger the Non-Secure fault, Non-Secure data access Bus fault
- `3` : trigger the Non-Secure fault, Non-Secure undefined instruction Usage fault
- `4` : trigger the Non-Secure fault, Non-Secure divide by 0 Usage fault
- `5` : trigger the Secure fault, Non-Secure data access from Secure RAM memory
- `6` : trigger the Secure fault, Non-Secure instruction execution from Secure Flash memory
- `7` : trigger the Secure fault, Secure undefined instruction Usage fault

## Example details

Clock Settings:

- **System Core Clock**: **160 MHz**

The example contains 2 applications: **Secure** and **Non-Secure**.

**Secure application** (bare-metal, no RTOS):

- setups the system (clocks, power, security and privilege rights (GTZC), caching) according to the STM32CubeMX configuration
- provides a function for triggering a fault on the Secure side

**Non-Secure application** (uses RTX RTOS and Standard C Library):

- setups the peripherals used by the non-secure application (LED, USART1) according to the STM32CubeMX configuration
- it runs 2 threads:
  - `AppThread` thread: Blink Green LED with 1 second interval
  - `FaultTriggerThread` thread: Trigger a fault according to input from the STDIO

When a fault is triggered the fault handler saves the fault information with the `ARM_FaultSave` function.

When valid fault information exists it can be viewed with **Component Viewer** or with `ARM_FaultRecord` function the fault information
can be output to the **Event Recorder** or with the `ARM_FaultPrint` function it can be output to the STDIO.
