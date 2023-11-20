# Fault example (Cortex-M33)

## STMicroelectronics B-U585I-IOT02A board

This project is a simple **Fault** component example running on **Arm Cortex-M33** microcontroller
on a STMicroelectronics [**B-U585I-IOT02A**](https://www.st.com/en/evaluation-tools/b-u585i-iot02a.html) evaluation board.

The application allows triggering of specific faults upon which the fault information is saved and system is reset.
When system restarts the fault information is output via the **Event Recorder** and via the **STDIO**.

The fault information can also be inspected with **Component Viewer** in a debug session.

## Prerequisites

### Software

 - [**CMSIS-Toolbox v2.0.0**](https://github.com/Open-CMSIS-Pack/cmsis-toolbox/releases) or newer
 - [**Keil MDK v5.38**](https://developer.arm.com/Tools%20and%20Software/Keil%20MDK) or newer containing:
   - Arm Compiler 6 (part of the MDK)
 - [**STM32CubeMX v6.8.1**](https://www.st.com/en/development-tools/stm32cubemx.html) or newer with:
   - STM32Cube MCU Package for STM32U5 Series v1.2.0
 - [**STM32CubeProgrammer**](https://www.st.com/en/development-tools/stm32cubeprog.html) utility
 - [**Arm GNU Toolchain v12.3-Rel1**](https://developer.arm.com/downloads/-/arm-gnu-toolchain-downloads)
   (only necessary when building example with GCC)

### CMSIS Packs

 - Required packs:
    - ARM::CMSIS-View
    - ARM::CMSIS v6.0.0 or newer
    - ARM::CMSIS-RTX v1.0.0 or newer
    - ARM::CMSIS-Compiler v2.0.0 or newer
    - Keil::STM32U5xx_DFP v2.1.0
    - Keil::B-U585I-IOT02A_BSP v1.0.0

### Hardware

This board has to be properly configured with **TrustZone** enabled. Please follow the steps below carefully:

Configure the following **Option bytes** with **STM32CubeProgrammer** utility:

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

### Build and Run with uVision

To try the example with uVision, do the following steps:

 1. Open the `Fault.uvmpw` in the uVision
 2. Build the `Project: Fault_S` project
 3. Build the `Project: Fault_NS` project
 4. Download the built application to the MCU's Flash
 5. Open the **Serial Terminal** application and connect to the **STMicroelectronics STLink Virtual COM Port (COMx)** (115200-8-N-1)
 6. Press **RESET** button on the board
 7. Follow the instructions in the **Serial Terminal** and observe the results

> **Note**
> - In the debug session fault information can be inspected in the **Component View** and **Event Recorder** windows

### Build with CMSIS-Toolbox

Alternatively, this example can be built with [**CMSIS-Toolbox**](https://github.com/Open-CMSIS-Pack/cmsis-toolbox).

To build the example with CMSIS-Toolbox do the following steps:

 1. Use the `csolution` command to create `.cprj` project files (for **Arm Compiler 6** toolchain):
    ```
    csolution convert -s Fault.csolution.yml
    ```
    or, for **GCC** toolchain use the following command:
    ```
    csolution convert -s Fault.csolution.yml -t GCC
    ```
 2. Use the `cbuild` command to create executable files for Secure and Non-secure applications:
    ```
    cbuild ./Secure/Fault_S.Debug+HW.cprj
    cbuild ./NonSecure/Fault_NS.Debug+HW.cprj
    ```

> **Note**
> - To run and debug executables built with CMSIS-Toolbox with uVision, it is necessary to adapt uVision settings relating to output file,
 and also adapt Debug.ini and Flash.ini scripts accordingly

## User Interface

This example uses **Serial Terminal** as User Interface.

The fault triggering is done by entering a number via **Serial Terminal** application (see possible values below).

 - 0: terminate the example
 - 1: trigger the Non-Secure fault, Non-Secure data access Memory Management fault
 - 2: trigger the Non-Secure fault, Non-Secure data access Bus fault
 - 3: trigger the Non-Secure fault, Non-Secure undefined instruction Usage fault
 - 4: trigger the Non-Secure fault, Non-Secure divide by 0 Usage fault
 - 5: trigger the Secure fault, Non-Secure data access from Secure RAM memory
 - 6: trigger the Secure fault, Non-Secure instruction execution from Secure Flash memory
 - 7: trigger the Secure fault, Secure undefined instruction Usage fault

## Example details

Clock Settings:

 - XTAL = MSIS =   4 MHz
 - Core =      **160 MHz**

The example contains 2 applications: Secure and Non-Secure.

**Secure application** (bare-metal, no RTOS):

 - setups the system (clocks, power, security and privilege rights (GTZC), caching) according to the CubeMX configuration
 - provides a function for triggering a fault on the Secure side

**Non-Secure application** (uses RTX RTOS and Standard C Library):

 - setups the peripherals used by the non-secure application (LEDs, button, UART1) according to the CubeMX configuration
 - it runs 2 threads:
    - `AppThread` thread: Blink Green LED with 1 second interval
    - `FaultTriggerThread` thread: Trigger a fault according to input from the STDIO

When a fault is triggered the fault handler saves the fault information with the `ARM_FaultSave` function.

When valid fault information exists it can be viewed with **Component Viewer** or with `ARM_FaultRecord` function the fault information
can be output to the **Event Recorder** or with the `ARM_FaultPrint` function it can be output to the STDIO.
