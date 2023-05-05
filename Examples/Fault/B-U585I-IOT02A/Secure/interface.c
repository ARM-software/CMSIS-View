/*----------------------------------------------------------------------------
 * Name:    interface.c
 * Purpose: Non-Secure to Secure function implementation
 *----------------------------------------------------------------------------*/

#include <arm_cmse.h>
#include "interface.h"

#include "RTE_Components.h"
#include  CMSIS_device_header

#include "../NonSecure/ARM_FaultTrigger.h"

/* Non-secure callable (entry) functions */

/* This function is used to trigger Secure fault */
__attribute__((cmse_nonsecure_entry)) void Secure_TriggerFault (uint32_t fault_id) {

  switch (fault_id) {
    case ARM_FAULT_ID_SEC_USG_UNDEFINED_INSTRUCTION:  // Trigger Secure - UsageFault - undefined instruction
      __ASM volatile (
        ".syntax unified\n"
        ".inst.w 0xF1234567\n"                        // Execute undefined 32-bit instruction encoded as 0xF1234567
      );
      break;
  }
}
