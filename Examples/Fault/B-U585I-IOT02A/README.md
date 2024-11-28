# Fault example for B-U585I-IOT02A board (Cortex-M33)

This project is a simple **Fault** component example running on an **Arm Cortex-M33** microcontroller
on a STMicroelectronics [**B-U585I-IOT02A**](https://www.st.com/en/evaluation-tools/b-u585i-iot02a.html) evaluation board.

The application allows triggering of specific faults upon which the fault information is saved and system is reset.
When system restarts the fault information is output via the **Event Recorder** and via the **STDIO (STLink Virtual COM Port)**.

The fault information can also be inspected with **Component Viewer** in a debug session.

## Prerequisites

### Software

- [**Arm Keil Studio for VS Code**](https://marketplace.visualstudio.com/items?itemName=Arm.keil-studio-pack)
- [**STM32CubeMX**](https://www.st.com/en/development-tools/stm32cubemx.html) **v6.12.1** or newer
- [**STM32CubeProgrammer**](https://www.st.com/en/development-tools/stm32cubeprog.html) utility
- [**Keil MDK**](https://developer.arm.com/Tools%20and%20Software/Keil%20MDK) **v5.41** or newer

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

To try the example with the **Arm Keil Studio**, follow the steps below:

 1. open the example in the **Visual Studio Code**.
 2. in the **Configure Solution** tab select the **AC6** compiler and click on the **OK** button.

    > **Note**  
    > To use **GCC** instead of AC6 compiler do the following:
    > - in the **Fault.csolution.yml** csolution file replace `compiler: AC6` with `compiler: GCC`.
    > - launch the **STM32CubeMX generator**: in the **CMSIS** extension project view click on the **Run Generator** button for **Device:CubeMX** component.
    > - in the **STM32CubeMX**: open from the menu **Project Manager - Project - Toolchain/IDE**,
    >   select **STM32CubeIDE** and clear **Generate Under Root** check box, and click on the **GENERATE CODE** button.

 3. build the solution (in the **CMSIS** extension view click on the **Build solution** button).
 4. download the built applications (Fault_S/Fault_NS) to the MCU's Flash (in the **CMSIS** extension view click on the **Run** button).  

    > **Note**  
    > In case of issues with download you can use **STM32CubeProgrammer** or **uVision**.  
    > Procedure using **uVision**: open **Fault.csolution.yml** in the **uVision**, click on the **Rebuild** button,
    > in **Options for Target** dialog under **Debug** select **ST-Link Debugger** and click on the **OK** button; now click on the **Download** button.

 5. open the **Serial Terminal** application and connect to the **STMicroelectronics STLink Virtual COM Port (COMx)** (115200-8-N-1).
 6. press the **RESET** button on the board.
 7. follow the instructions in the **Serial Terminal** and observe the results.

### uVision

To try the example with the **uVision**, follow the steps below:

 1. open the example in the **uVision**.
 2. build the projects (select **Project - Batch Build**).
 3. download the built applications (Fault_S/Fault_NS) to the MCU's Flash (click on the **Download** button).
 4. open the **Serial Terminal** application and connect to the **STMicroelectronics STLink Virtual COM Port (COMx)** (115200-8-N-1).
 5. press the **RESET** button on the board.
 6. follow the instructions in the **Serial Terminal** and observe the results.

> **Note**  
> In the debug session fault information can be inspected in the **Component View** and **Event Recorder** windows.

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
