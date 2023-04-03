/*
 * Copyright (c) 2022-2023 Arm Limited. All rights reserved.
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

// Check if Arm Architecture is supported
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

// Determine if architecture is Armv8/8.1-M architecture
#if   ((defined(__ARM_ARCH_8M_BASE__)   && (__ARM_ARCH_8M_BASE__   != 0)) || \
       (defined(__ARM_ARCH_8M_MAIN__)   && (__ARM_ARCH_8M_MAIN__   != 0)) || \
       (defined(__ARM_ARCH_8_1M_MAIN__) && (__ARM_ARCH_8_1M_MAIN__ != 0))    )
#define ARM_FAULT_ARCH_ARMV8x_M        (1)
#else
#define ARM_FAULT_ARCH_ARMV8x_M        (0)
#endif

// Determine if architecture is Armv8-M Baseline architecture
#if    (defined(__ARM_ARCH_8M_BASE__)   && (__ARM_ARCH_8M_BASE__   != 0))
#define ARM_FAULT_ARCH_ARMV8_M_BASE    (1)
#else
#define ARM_FAULT_ARCH_ARMV8_M_BASE    (0)
#endif

// Determine if architecture is Armv8-M Mainline or Armv8.1 architecture
#if   ((defined(__ARM_ARCH_8M_MAIN__)   && (__ARM_ARCH_8M_MAIN__   != 0)) || \
       (defined(__ARM_ARCH_8_1M_MAIN__) && (__ARM_ARCH_8_1M_MAIN__ != 0))    )
#define ARM_FAULT_ARCH_ARMV8x_M_MAIN   (1)
#else
#define ARM_FAULT_ARCH_ARMV8x_M_MAIN   (0)
#endif

// Determine if the code is compiled for Secure World
#if    (defined (__ARM_FEATURE_CMSE) && (__ARM_FEATURE_CMSE == 3))
#define ARM_FAULT_TZ_SECURE            (1)
#else
#define ARM_FAULT_TZ_SECURE            (0)
#endif

// Fault Storage definitions
#define ARM_FAULT_FAULT_INFO_VER_MAJOR (0U)             // ARM Fault ARM_FaultInfo type version.major
#define ARM_FAULT_FAULT_INFO_VER_MINOR (1U)             // ARM Fault ARM_FaultInfo type version.minor
#define ARM_FAULT_MAGIC_NUMBER         (0x52746C46U)    // ARM Fault Magic number (ASCII "FltR")
#define ARM_FAULT_CRC32_INIT_VAL       (0xFFFFFFFFU)    // ARM Fault CRC-32 initial value
#define ARM_FAULT_CRC32_POLYNOM        (0x04C11DB7U)    // ARM Fault CRC-32 polynom


#define ARM_FAULT_VERSION              "0.4.0"

#ifdef __cplusplus
extern "C" {
#endif

// Fault information type definition
typedef struct {
  uint32_t magic_number;
  uint32_t crc32;
  uint32_t count;

  // Fault information structure
  struct {
    struct {
      uint8_t minor;                    // Fault information structure version: minor
      uint8_t major;                    // Fault information structure version: major
    } version;
    uint16_t fault_regs    :  1;        // == 1 - contains fault registers
    uint16_t armv8m        :  1;        // == 1 - contains Armv8/8.1-M related information
    uint16_t tz_secure     :  1;        // == 1 - recording was done running in Secure World
    uint16_t reserved      : 13;        // Reserved (0)
  } type;

  // State context (same as Basic Stack Frame)
  uint32_t R0;                          // R0  register value before exception
  uint32_t R1;                          // R1  register value before exception
  uint32_t R2;                          // R2  register value before exception
  uint32_t R3;                          // R3  register value before exception
  uint32_t R12;                         // R12 register value before exception
  uint32_t LR;                          // Link Register (R14) value before exception
  uint32_t ReturnAddress;               // Return address from exception
  uint32_t xPSR;                        // Program Status Register value before exception

  // Common Registers
  uint32_t xPSR_in_handler;             // Program Status Register value, in exception handler
  uint32_t EXC_RETURN;                  // Exception Return code (LR), in exception handler
  uint32_t MSP;                         // Main    Stack Pointer value
  uint32_t PSP;                         // Process Stack Pointer value

#if (ARM_FAULT_FAULT_REGS_EXIST != 0)
  // Fault Registers
  uint32_t SCB_CFSR;                    // System Control Block - Configurable Fault Status Register value
  uint32_t SCB_HFSR;                    // System Control Block - HardFault          Status Register value
  uint32_t SCB_DFSR;                    // System Control Block - Debug Fault        Status Register value
  uint32_t SCB_MMFAR;                   // System Control Block - MemManage Fault    Status Register value
  uint32_t SCB_BFAR;                    // System Control Block - BusFault           Status Register value
  uint32_t SCB_AFSR;                    // System Control Block - Auxiliary Fault    Status Register value
#endif

#if (ARM_FAULT_ARCH_ARMV8x_M != 0)
  // Additional state context (only for Armv8/8.1-M arch)
  uint32_t IntegritySignature;          // Integrity Signature
  uint32_t Reserved;                    // Reserved
  uint32_t R4;                          // R4  register value before exception
  uint32_t R5;                          // R5  register value before exception
  uint32_t R6;                          // R6  register value before exception
  uint32_t R7;                          // R7  register value before exception
  uint32_t R8;                          // R8  register value before exception
  uint32_t R9;                          // R9  register value before exception
  uint32_t R10;                         // R10 register value before exception
  uint32_t R11;                         // R11 register value before exception

  // Additional Armv8/8.1-M arch specific Registers
  uint32_t MSPLIM;                      // Main    Stack Pointer Limit Register value
  uint32_t PSPLIM;                      // Process Stack Pointer Limit Register value
#endif

#if (ARM_FAULT_ARCH_ARMV8x_M_MAIN != 0)
  // Additional Armv8/8.1-M Mainline arch specific Fault Registers
  uint32_t SCB_SFSR;                    // System Control Block - Secure Fault Status  Register value
  uint32_t SCB_SFAR;                    // System Control Block - Secure Fault Address Register value
#endif
} ARM_FaultInfo_t;

// ARM Fault variables ---------------------------------------------------------

// Fault component version information
extern const char ARM_FaultVersion[];

// Fault Information
extern ARM_FaultInfo_t ARM_FaultInfo;

// ARM Fault Storage functions -------------------------------------------------

/// Clear the saved fault information.
extern void ARM_FaultClear (void);

/// Check if the fault occurred and if the fault information was saved properly.
/// \return       status (1=Fault occurred and valid fault information exists,
///                       0=no fault information saved yet or is invalid)
extern uint32_t ARM_FaultOccurred (void);

/// Save the fault information.
extern void ARM_FaultSave (void);

/// Callback function called after fault information was saved.
extern void ARM_FaultExit (void);

// ARM Fault User Code template ------------------------

/// Output decoded fault information via STDIO.
extern void ARM_FaultPrint (void);

// ARM Fault Record function ---------------------------------------------------

/// Output decoded fault information via EventRecorder.
extern void ARM_FaultRecord (void);

#ifdef __cplusplus
}
#endif

#endif /* ARM_FAULT_H__ */
