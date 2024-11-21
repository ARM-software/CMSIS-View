/*
 * Copyright (c) 2022-2024 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the License); you may
 * not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an AS IS BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#ifndef ARM_FAULT_H__
#define ARM_FAULT_H__

#include <stdint.h>

#include "RTE_Components.h"
#include  CMSIS_device_header

// Check if Arm architecture is supported
#if  ((!defined(__ARM_ARCH_6M__)        || (__ARM_ARCH_6M__        == 0)) && \
      (!defined(__ARM_ARCH_7M__)        || (__ARM_ARCH_7M__        == 0)) && \
      (!defined(__ARM_ARCH_7EM__)       || (__ARM_ARCH_7EM__       == 0)) && \
      (!defined(__ARM_ARCH_8M_BASE__)   || (__ARM_ARCH_8M_BASE__   == 0)) && \
      (!defined(__ARM_ARCH_8M_MAIN__)   || (__ARM_ARCH_8M_MAIN__   == 0)) && \
      (!defined(__ARM_ARCH_8_1M_MAIN__) || (__ARM_ARCH_8_1M_MAIN__ == 0))    )
#error "Unknown or unsupported Arm Architecture!"
#endif

// Determine if fault registers are available
#if   ((defined(__ARM_ARCH_7M__)        && (__ARM_ARCH_7M__        != 0)) || \
       (defined(__ARM_ARCH_7EM__)       && (__ARM_ARCH_7EM__       != 0)) || \
       (defined(__ARM_ARCH_8M_MAIN__)   && (__ARM_ARCH_8M_MAIN__   != 0)) || \
       (defined(__ARM_ARCH_8_1M_MAIN__) && (__ARM_ARCH_8_1M_MAIN__ != 0))    )
#define ARM_FAULT_FAULT_REGS_EXIST     (1)
#else
#define ARM_FAULT_FAULT_REGS_EXIST     (0)
#endif

// Determine if architecture is Armv8/8.1-M
#if   ((defined(__ARM_ARCH_8M_BASE__)   && (__ARM_ARCH_8M_BASE__   != 0)) || \
       (defined(__ARM_ARCH_8M_MAIN__)   && (__ARM_ARCH_8M_MAIN__   != 0)) || \
       (defined(__ARM_ARCH_8_1M_MAIN__) && (__ARM_ARCH_8_1M_MAIN__ != 0))    )
#define ARM_FAULT_ARCH_ARMV8x_M        (1)
#else
#define ARM_FAULT_ARCH_ARMV8x_M        (0)
#endif

// Determine if architecture is Armv8-M Baseline
#if    (defined(__ARM_ARCH_8M_BASE__)   && (__ARM_ARCH_8M_BASE__   != 0))
#define ARM_FAULT_ARCH_ARMV8_M_BASE    (1)
#else
#define ARM_FAULT_ARCH_ARMV8_M_BASE    (0)
#endif

// Determine if architecture is Armv8-M Mainline or Armv8.1 Mainline
#if   ((defined(__ARM_ARCH_8M_MAIN__)   && (__ARM_ARCH_8M_MAIN__   != 0)) || \
       (defined(__ARM_ARCH_8_1M_MAIN__) && (__ARM_ARCH_8_1M_MAIN__ != 0))    )
#define ARM_FAULT_ARCH_ARMV8x_M_MAIN   (1)
#else
#define ARM_FAULT_ARCH_ARMV8x_M_MAIN   (0)
#endif

// Determine if architecture is Armv8.1-M Mainline
#if    (defined(__ARM_ARCH_8_1M_MAIN__) && (__ARM_ARCH_8_1M_MAIN__ != 0))
#define ARM_FAULT_ARCH_ARMV8_1M_MAIN   (1)
#else
#define ARM_FAULT_ARCH_ARMV8_1M_MAIN   (0)
#endif

// Determine if the code is compiled with Cortex-M Security Extensions enabled
#if     defined (__ARM_FEATURE_CMSE)
#define ARM_FAULT_TZ_ENABLED           (1)
#else
#define ARM_FAULT_TZ_ENABLED           (0)
#endif

// Determine if the code is compiled for Secure World
#if    (defined (__ARM_FEATURE_CMSE) && (__ARM_FEATURE_CMSE == 3))
#define ARM_FAULT_TZ_SECURE            (1)
#else
#define ARM_FAULT_TZ_SECURE            (0)
#endif

// Fault component version
#define ARM_FAULT_VERSION              "1.1.0"

// Fault Information structure type version
#define ARM_FAULT_FAULT_INFO_VER_MAJOR (1U)             // ARM_FaultInfo type Version.Major
#define ARM_FAULT_FAULT_INFO_VER_MINOR (1U)             // ARM_FaultInfo type Version.Minor

#ifdef __cplusplus
extern "C" {
#endif

/// Fault information structure type definition
typedef struct {
  uint32_t MagicNumber;                 //!< Magic number (ASCII "FltR")
  uint32_t CRC32;                       //!< CRC32 of the structure content (excluding MagicNumber and CRC32 fields)
  uint32_t Count;                       //!< Saved faults counter

  struct {                              // Version
    uint8_t Minor;                      //!< Fault information structure version: Minor, see \ref ARM_FAULT_FAULT_INFO_VER_MINOR
    uint8_t Major;                      //!< Fault information structure version: Major, see \ref ARM_FAULT_FAULT_INFO_VER_MAJOR
  } Version;

  struct {                              // Content
                                        // Compile-time information
    uint16_t FaultRegsExist    :  1;    //!< Fault registers: 0 - absent; 1 - available
    uint16_t Armv8xM_Main      :  1;    //!< Armv8/8.1-M Mainline information: 0 - absent; 1 - available
    uint16_t TZ_Enabled        :  1;    //!< TrustZone (Cortex-M security extensions): 0 - not enabled; 1 - enabled
    uint16_t TZ_SaveMode       :  1;    //!< Fault information was saved in: 0 - TrustZone-disabled or non-secure mode; 1 - secure mode

                                        // Runtime-time information
    uint16_t TZ_FaultMode      :  1;    //!< Fault happened in: 0 - TrustZone-disabled or non-secure mode; 1 - secure mode
    uint16_t StateContext      :  1;    //!< State Context: 0 - was not saved; 1 - was saved
    uint16_t AdditionalContext :  1;    //!< Additional State Context: 0 - was not saved; 1 - was saved
    uint16_t LimitRegs         :  1;    //!< MSPLIM and PSPLIM: 0 - were not saved; 1 - were saved
    uint16_t FaultRegs         :  1;    //!< Fault registers: 0 - were not saved; 1 - were saved
    uint16_t SecureFaultRegs   :  1;    //!< Secure Fault registers: 0 - were not saved; 1 - were saved
    uint16_t RAS_FaultReg      :  1;    //!< RAS Fault register: 0 - was not saved; 1 - was saved

    uint16_t Reserved          :  5;    //!< Reserved (0)
  } Content;

  struct {                              // Registers
    uint32_t R0;                        //!< R0  Register value
    uint32_t R1;                        //!< R1  Register value
    uint32_t R2;                        //!< R2  Register value
    uint32_t R3;                        //!< R3  Register value
    uint32_t R4;                        //!< R4  Register value
    uint32_t R5;                        //!< R5  Register value
    uint32_t R6;                        //!< R6  Register value
    uint32_t R7;                        //!< R7  Register value
    uint32_t R8;                        //!< R8  Register value
    uint32_t R9;                        //!< R9  Register value
    uint32_t R10;                       //!< R10 Register value
    uint32_t R11;                       //!< R11 Register value
    uint32_t R12;                       //!< R12 Register value
    uint32_t LR;                        //!< Link Register (R14) value
    uint32_t ReturnAddress;             //!< Return address from exception
    uint32_t xPSR;                      //!< Program Status Register value
    uint32_t MSP;                       //!< Main    Stack Pointer value
    uint32_t PSP;                       //!< Process Stack Pointer value
    uint32_t MSPLIM;                    //!< Main    Stack Pointer Limit Register value (only for Armv8/8.1-M arch)
    uint32_t PSPLIM;                    //!< Process Stack Pointer Limit Register value (only for Armv8/8.1-M arch)
  } Registers;

  struct {                              // Exception State
    uint32_t xPSR;                      //!< Program Status Register value, in exception handler
    uint32_t EXC_RETURN;                //!< Exception Return code (LR), in exception handler
  } ExceptionState;

#if (ARM_FAULT_FAULT_REGS_EXIST != 0)
  struct {                              // Fault Registers
    uint32_t CFSR;                      //!< System Control Block - Configurable Fault Status  Register value
    uint32_t HFSR;                      //!< System Control Block - HardFault          Status  Register value
    uint32_t DFSR;                      //!< System Control Block - Debug Fault        Status  Register value
    uint32_t MMFAR;                     //!< System Control Block - MemManage Fault    Address Register value
    uint32_t BFAR;                      //!< System Control Block - BusFault           Address Register value
    uint32_t AFSR;                      //!< System Control Block - Auxiliary Fault    Status  Register value

    // Additional Armv8/8.1-M Mainline arch specific fault registers
    uint32_t SFSR;                      //!< System Control Block - Secure Fault       Status  Register value
    uint32_t SFAR;                      //!< System Control Block - Secure Fault       Address Register value

    // Additional Armv8.1-M Mainline arch specific fault register
    uint32_t RFSR;                      //!< System Control Block - RAS Fault          Status  Register value
  } FaultRegisters;
#endif
} ARM_FaultInfo_t;

// ARM Fault variables ---------------------------------------------------------

//! Fault component version information
extern const char ARM_FaultVersion[];

//! Fault Information
extern ARM_FaultInfo_t ARM_FaultInfo;

// ARM Fault Storage functions -------------------------------------------------

/// \brief Clear the saved fault information.
extern void ARM_FaultClear (void);

/// \brief Check if the fault occurred and if the fault information was saved properly.
/// \return       status (1 = fault occurred and valid fault information exists,
///                       0 = no fault information saved yet or fault information is invalid)
extern uint32_t ARM_FaultOccurred (void);

/// \brief Save the fault information.
extern void ARM_FaultSave (void);

/// \brief Callback function called after fault information was saved.
extern void ARM_FaultExit (void);

// ARM Fault User Code template ------------------------

/// \brief Output decoded fault information via STDIO.
extern void ARM_FaultPrint (void);

// ARM Fault Record function ---------------------------------------------------

/// \brief Output decoded fault information via Event Recorder.
extern void ARM_FaultRecord (void);

#ifdef __cplusplus
}
#endif

#endif /* ARM_FAULT_H__ */
