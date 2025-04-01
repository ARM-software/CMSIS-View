/*
 * Copyright (c) 2023-2025 Arm Limited. All rights reserved.
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

#include "ARM_FaultTrigger.h"

#include "RTE_Components.h"
#include  CMSIS_device_header

// ARM_FaultTrigger function ---------------------------------------------------

/**
  Trigger a fault.
  \param[in]    fault_id    Fault ID of the fault to be triggered
*/
void ARM_FaultTrigger (uint32_t fault_id) {
  volatile uint32_t val;
  void (*ptr_func) (void);

  switch (fault_id) {
    case ARM_FAULT_ID_MEM_DATA:                     // Trigger MemManage fault - data access
      val = *((uint32_t *)0x20000000);              // Read from address not allowed by the MPU (non-privileged access not allowed)
      break;

    case ARM_FAULT_ID_BUS_DATA_PRECISE:             // Trigger BusFault - data access (precise)
      val = *((uint32_t *)0x3FFFFFFC);              // Read from invalid address
      break;

    case ARM_FAULT_ID_BUS_DATA_IMPRECISE:           // Trigger BusFault - data access (imprecise)
      *((uint32_t *)0x3FFFFFFC) = 1U;               // Write to invalid address
      break;

    case ARM_FAULT_ID_BUS_INSTRUCTION:              // Trigger BusFault - instruction execution
      ptr_func = (void (*) (void))(0x1FFFFFFC);
      ptr_func();                                   // Call function from invalid address
      break;

    case ARM_FAULT_ID_USG_UNDEFINED_INSTRUCTION:    // Trigger UsageFault - undefined instruction
      __ASM volatile (
        ".syntax unified\n"
        ".inst.w 0xF1234567\n"                      // Execute undefined 32-bit instruction encoded as 0xF1234567
      );
      break;

    case ARM_FAULT_ID_USG_DIV_0:                    // Trigger UsageFault - divide by 0
      val = 0U;
      val = 123/val;
      break;

    default:
      break;
  }
}
