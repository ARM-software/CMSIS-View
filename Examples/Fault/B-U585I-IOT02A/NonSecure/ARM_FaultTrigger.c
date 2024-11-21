/*
 * Copyright (c) 2024 Arm Limited. All rights reserved.
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
#include "../STM32CubeMX/B-U585I-IOT02A/STM32CubeMX/Secure_nsclib/secure_nsc.h"

// ARM_FaultTrigger function ---------------------------------------------------

/**
  Trigger a fault.
  \param[in]    fault_id    Fault ID of the fault to be triggered
*/
void ARM_FaultTrigger (uint32_t fault_id) {
  volatile uint32_t val;
  void (*ptr_func) (void);

  switch (fault_id) {
    case ARM_FAULT_ID_MEM_DATA:                     // Trigger Non-Secure MemManage fault - data access
      val = *((uint32_t *)0x20040000);              // Read from address not allowed by the MPU (non-privileged access not allowed)
      break;

    case ARM_FAULT_ID_BUS_DATA:                     // Trigger Non-Secure BusFault - data access
      val = *((uint32_t *)0x200C0000);              // Read from invalid RAM address
      break;

    case ARM_FAULT_ID_USG_UNDEFINED_INSTRUCTION:    // Trigger Non-Secure UsageFault - undefined instruction
      __ASM volatile (
        ".syntax unified\n"
        ".inst.w 0xF1234567\n"                      // Execute undefined 32-bit instruction encoded as 0xF1234567
      );
      break;

    case ARM_FAULT_ID_USG_DIV_0:                    // Trigger Non-Secure UsageFault - divide by 0
      val = 0U;
      val = 123/val;
      break;

    case ARM_FAULT_ID_SEC_DATA:                     // Trigger Secure BusFault - data access
      val = *((uint32_t *)0x30000000);              // Read from Secure RAM address
      break;

    case ARM_FAULT_ID_SEC_INSTRUCTION:              // Trigger Secure BusFault - instruction execution
      ptr_func = (void (*) (void))(0xC000000);
      ptr_func();                                   // Call function from Secure Flash address
      break;

    case ARM_FAULT_ID_SEC_USG_UNDEFINED_INSTRUCTION:  // Trigger Secure - UsageFault - undefined instruction
      Secure_TriggerFault(fault_id);                // Call Secure function that will trigger a fault
      break;

    default:
      break;
  }
}
