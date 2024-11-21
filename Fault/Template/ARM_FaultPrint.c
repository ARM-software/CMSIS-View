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

#include "ARM_Fault.h"

#include <stdio.h>

// General defines
#ifndef EXC_RETURN_SPSEL
#define EXC_RETURN_SPSEL       (1UL << 2)
#endif

// Armv8/8.1-M Mainline architecture related defines
#if    (ARM_FAULT_ARCH_ARMV8x_M_MAIN != 0)
#ifndef SAU_SFSR_LSERR_Msk
#define SAU_SFSR_LSERR_Msk     (1UL << 7)               // SAU SFSR: LSERR Mask
#endif
#ifndef SAU_SFSR_SFARVALID_Msk
#define SAU_SFSR_SFARVALID_Msk (1UL << 6)               // SAU SFSR: SFARVALID Mask
#endif
#ifndef SAU_SFSR_LSPERR_Msk
#define SAU_SFSR_LSPERR_Msk    (1UL << 5)               // SAU SFSR: LSPERR Mask
#endif
#ifndef SAU_SFSR_INVTRAN_Msk
#define SAU_SFSR_INVTRAN_Msk   (1UL << 4)               // SAU SFSR: INVTRAN Mask
#endif
#ifndef SAU_SFSR_AUVIOL_Msk
#define SAU_SFSR_AUVIOL_Msk    (1UL << 3)               // SAU SFSR: AUVIOL Mask
#endif
#ifndef SAU_SFSR_INVER_Msk
#define SAU_SFSR_INVER_Msk     (1UL << 2)               // SAU SFSR: INVER Mask
#endif
#ifndef SAU_SFSR_INVIS_Msk
#define SAU_SFSR_INVIS_Msk     (1UL << 1)               // SAU SFSR: INVIS Mask
#endif
#ifndef SAU_SFSR_INVEP_Msk
#define SAU_SFSR_INVEP_Msk     (1UL)                    // SAU SFSR: INVEP Mask
#endif
#endif

// ARM_FaultPrint function -----------------------------------------------------

/**
  Output decoded fault information via STDIO.
  Should be called when system is running in normal operating mode with
  standard input/output fully functional.
*/
void ARM_FaultPrint (void) {
  int8_t fault_info_valid;

  /* Check if there is available valid fault information */
  fault_info_valid = (int8_t)ARM_FaultOccurred();

  // Output: Header and version information
  printf("\n --- Fault (v%s) ---\n\n", (const char *)ARM_FaultVersion);

  // Output: Message if fault info is invalid
  if (fault_info_valid == 0) {
    printf("\n  No fault saved yet or fault information is invalid!\n\n");
    return;
  }

  // Output: Fault count
  printf("  Fault count:         %u\n\n", (unsigned int)ARM_FaultInfo.Count);

  // Output: Exception which saved the fault information
  printf("  Exception Handler:   ");

  if (ARM_FaultInfo.Content.TZ_Enabled != 0U) {
    if (ARM_FaultInfo.Content.TZ_SaveMode != 0U) {
      printf("Secure - ");
    } else {
      printf("Non-Secure - ");
    }
  }

  switch (ARM_FaultInfo.ExceptionState.xPSR & IPSR_ISR_Msk) {
    case 3:
      printf("HardFault");
      break;
    case 4:
      printf("MemManage fault");
      break;
    case 5:
      printf("BusFault");
      break;
    case 6:
      printf("UsageFault");
      break;
    case 7:
      printf("SecureFault");
      break;
    default:
      printf("unknown, exception number = %u", (unsigned int)(ARM_FaultInfo.ExceptionState.xPSR & IPSR_ISR_Msk));
      break;
  }
  printf("\n");

#if (ARM_FAULT_ARCH_ARMV8x_M != 0)
  // Output: state in which the fault occurred
  if (ARM_FaultInfo.Content.TZ_Enabled != 0U) {
    printf("  State:               ");

    if (ARM_FaultInfo.Content.TZ_FaultMode != 0U) {
      printf("Secure");
    } else {
      printf("Non-Secure");
    }
    printf("\n");
  }
#endif

  // Output: Mode in which the fault occurred
  printf("  Mode:                ");

  if ((ARM_FaultInfo.ExceptionState.EXC_RETURN & EXC_RETURN_SPSEL) == 0U) {
    printf("Handler");
  } else {
    printf("Thread");
  }
  printf("\n");

#if (ARM_FAULT_FAULT_REGS_EXIST != 0)   // If fault registers exist
  /* Output: Decoded HardFault information */
  if (ARM_FaultInfo.Content.FaultRegs != 0U) {
    uint32_t scb_hfsr = ARM_FaultInfo.FaultRegisters.HFSR;

    if ((scb_hfsr & (SCB_HFSR_VECTTBL_Msk   |
                     SCB_HFSR_FORCED_Msk    |
                     SCB_HFSR_DEBUGEVT_Msk  )) != 0U) {

      printf("  Fault:               HardFault - ");

      if ((scb_hfsr & SCB_HFSR_VECTTBL_Msk) != 0U) {
        printf("Bus error on vector read");
      }
      if ((scb_hfsr & SCB_HFSR_FORCED_Msk) != 0U) {
        printf("Escalated fault (original fault was disabled or it caused another lower priority fault)");
      }
      if ((scb_hfsr & SCB_HFSR_DEBUGEVT_Msk) != 0U) {
        printf("Breakpoint hit with Debug Monitor disabled");
      }
      printf("\n");
    }
  }

  /* Output: Decoded MemManage fault information */
  if (ARM_FaultInfo.Content.FaultRegs != 0U) {
    uint32_t scb_cfsr   = ARM_FaultInfo.FaultRegisters.CFSR;
    uint32_t scb_mmfar  = ARM_FaultInfo.FaultRegisters.MMFAR;
    uint8_t  faults_cnt = 0U;

    if ((scb_cfsr & (SCB_CFSR_IACCVIOL_Msk  |
                     SCB_CFSR_DACCVIOL_Msk  |
                     SCB_CFSR_MUNSTKERR_Msk |
#ifdef SCB_CFSR_MLSPERR_Msk
                     SCB_CFSR_MLSPERR_Msk   |
#endif
                     SCB_CFSR_MSTKERR_Msk   )) != 0U) {

      printf("  Fault:               MemManage - ");

      if ((scb_cfsr & SCB_CFSR_IACCVIOL_Msk) != 0U) {
        faults_cnt++;
        printf("Instruction execution failure due to MPU violation or fault");
      }
      if ((scb_cfsr & SCB_CFSR_DACCVIOL_Msk) != 0U) {
        faults_cnt++;
        if (faults_cnt > 1U) {
          printf("\n  Fault:               MemManage - ");
        }
        printf("Data access failure due to MPU violation or fault");
      }
      if ((scb_cfsr & SCB_CFSR_MUNSTKERR_Msk) != 0U) {
        faults_cnt++;
        if (faults_cnt > 1U) {
          printf("\n  Fault:               MemManage - ");
        }
        printf("Exception exit unstacking failure due to MPU access violation");
      }
      if ((scb_cfsr & SCB_CFSR_MSTKERR_Msk) != 0U) {
        faults_cnt++;
        if (faults_cnt > 1U) {
          printf("\n  Fault:               MemManage - ");
        }
        printf("Exception entry stacking failure due to MPU access violation");
      }
#ifdef SCB_CFSR_MLSPERR_Msk
      if ((scb_cfsr & SCB_CFSR_MLSPERR_Msk) != 0U) {
        faults_cnt++;
        if (faults_cnt > 1U) {
          printf("\n  Fault:               MemManage - ");
        }
        printf("Floating-point lazy stacking failure due to MPU access violation");
      }
#endif
      if ((scb_cfsr & SCB_CFSR_MMARVALID_Msk) != 0U) {
        printf(", fault address 0x%08X", (unsigned int)scb_mmfar);
      }
      printf("\n");
    }
  }

  /* Output: Decoded BusFault information */
  if (ARM_FaultInfo.Content.FaultRegs != 0U) {
    uint32_t scb_cfsr   = ARM_FaultInfo.FaultRegisters.CFSR;
    uint32_t scb_bfar   = ARM_FaultInfo.FaultRegisters.BFAR;
    uint8_t  faults_cnt = 0U;

    if ((scb_cfsr & (SCB_CFSR_IBUSERR_Msk     |
                     SCB_CFSR_PRECISERR_Msk   |
                     SCB_CFSR_IMPRECISERR_Msk |
                     SCB_CFSR_UNSTKERR_Msk    |
#ifdef SCB_CFSR_LSPERR_Msk
                     SCB_CFSR_LSPERR_Msk      |
#endif
                     SCB_CFSR_STKERR_Msk      )) != 0U) {

      printf("  Fault:               BusFault - ");

      if ((scb_cfsr & SCB_CFSR_IBUSERR_Msk) != 0U) {
        faults_cnt++;
        printf("Instruction prefetch failure due to bus fault");
      }
      if ((scb_cfsr & SCB_CFSR_PRECISERR_Msk) != 0U) {
        faults_cnt++;
        if (faults_cnt > 1U) {
          printf("\n  Fault:               BusFault - ");
        }
        printf("Data access failure due to bus fault (precise)");
      }
      if ((scb_cfsr & SCB_CFSR_IMPRECISERR_Msk) != 0U) {
        faults_cnt++;
        if (faults_cnt > 1U) {
          printf("\n  Fault:               BusFault - ");
        }
        printf("Data access failure due to bus fault (imprecise)");
      }
      if ((scb_cfsr & SCB_CFSR_UNSTKERR_Msk) != 0U) {
        faults_cnt++;
        if (faults_cnt > 1U) {
          printf("\n  Fault:               BusFault - ");
        }
        printf("Exception exit unstacking failure due to bus fault");
      }
      if ((scb_cfsr & SCB_CFSR_STKERR_Msk) != 0U) {
        faults_cnt++;
        if (faults_cnt > 1U) {
          printf("\n  Fault:               BusFault - ");
        }
        printf("Exception entry stacking failure due to bus fault");
      }
#ifdef SCB_CFSR_LSPERR_Msk
      if ((scb_cfsr & SCB_CFSR_LSPERR_Msk) != 0U) {
        faults_cnt++;
        if (faults_cnt > 1U) {
          printf("\n  Fault:               BusFault - ");
        }
        printf("Floating-point lazy stacking failure due to bus fault");
      }
#endif
      if ((scb_cfsr & SCB_CFSR_BFARVALID_Msk) != 0U) {
        printf(", fault address 0x%08X", (unsigned int)scb_bfar);
      }
      printf("\n");
    }
  }

  /* Output Decoded UsageFault information */
  if (ARM_FaultInfo.Content.FaultRegs != 0U) {
    uint32_t scb_cfsr = ARM_FaultInfo.FaultRegisters.CFSR;

    if ((scb_cfsr & (SCB_CFSR_UNDEFINSTR_Msk |
                     SCB_CFSR_INVSTATE_Msk   |
                     SCB_CFSR_INVPC_Msk      |
                     SCB_CFSR_NOCP_Msk       |
#ifdef SCB_CFSR_STKOF_Msk
                     SCB_CFSR_STKOF_Msk      |
#endif
                     SCB_CFSR_UNALIGNED_Msk  |
                     SCB_CFSR_DIVBYZERO_Msk  )) != 0U) {

      printf("  Fault:               UsageFault - ");

      if ((scb_cfsr & SCB_CFSR_UNDEFINSTR_Msk) != 0U) {
        printf("Execution of undefined instruction");
      }
      if ((scb_cfsr & SCB_CFSR_INVSTATE_Msk) != 0U) {
        printf("Execution of Thumb instruction with Thumb mode turned off");
      }
      if ((scb_cfsr & SCB_CFSR_INVPC_Msk) != 0U) {
        printf("Invalid exception return value");
      }
      if ((scb_cfsr & SCB_CFSR_NOCP_Msk) != 0U) {
        printf("Coprocessor instruction with coprocessor disabled or non-existent");
      }
#ifdef SCB_CFSR_STKOF_Msk
      if ((scb_cfsr & SCB_CFSR_STKOF_Msk) != 0U) {
        printf("Stack overflow");
      }
#endif
      if ((scb_cfsr & SCB_CFSR_UNALIGNED_Msk) != 0U) {
        printf("Unaligned load/store");
      }
      if ((scb_cfsr & SCB_CFSR_DIVBYZERO_Msk) != 0U) {
        printf("Divide by 0");
      }
      printf("\n");
    }
  }

#if (ARM_FAULT_ARCH_ARMV8x_M_MAIN != 0)
  /* Output: Decoded SecureFault information */
  if (ARM_FaultInfo.Content.SecureFaultRegs != 0U) {
    uint32_t scb_sfsr = ARM_FaultInfo.FaultRegisters.SFSR;
    uint32_t scb_sfar = ARM_FaultInfo.FaultRegisters.SFAR;

    if ((scb_sfsr & (SAU_SFSR_INVEP_Msk   |
                     SAU_SFSR_INVIS_Msk   |
                     SAU_SFSR_INVER_Msk   |
                     SAU_SFSR_AUVIOL_Msk  |
                     SAU_SFSR_INVTRAN_Msk |
                     SAU_SFSR_LSPERR_Msk  |
                     SAU_SFSR_LSERR_Msk   )) != 0U) {

      printf("  Fault:               SecureFault - ");

      if ((scb_sfsr & SAU_SFSR_INVEP_Msk) != 0U) {
        printf("Invalid entry point due to invalid attempt to enter Secure state");
      }
      if ((scb_sfsr & SAU_SFSR_INVIS_Msk) != 0U) {
        printf("Invalid integrity signature in exception stack frame found on unstacking");
      }
      if ((scb_sfsr & SAU_SFSR_INVER_Msk) != 0U) {
        printf("Invalid exception return due to mismatch on EXC_RETURN.DCRS or EXC_RETURN.ES");
      }
      if ((scb_sfsr & SAU_SFSR_AUVIOL_Msk) != 0U) {
        printf("Attribution unit violation due to Non-secure access to Secure address space");
      }
      if ((scb_sfsr & SAU_SFSR_INVTRAN_Msk) != 0U) {
        printf("Invalid transaction caused by domain crossing branch not flagged as such");
      }
      if ((scb_sfsr & SAU_SFSR_LSPERR_Msk) != 0U) {
        printf("Lazy stacking preservation failure due to SAU or IDAU violation");
      }
      if ((scb_sfsr & SAU_SFSR_LSERR_Msk) != 0U) {
        printf("Lazy stacking activation or deactivation failure");
      }
      if ((scb_sfsr & SAU_SFSR_SFARVALID_Msk) != 0U) {
        printf(", fault address 0x%08X", (unsigned int)scb_sfar);
      }
      printf("\n");
    }
  }
#endif
#endif

  /* Output: Program Counter */
  /* Output here is named PC (Program Counter) since in most situations stacked Return Address will be
     the address of the instruction which caused the fault, there are some exceptions (asynchronous faults)
     but these are for the sake of simplicity not taken into account here */
  printf("  Program Counter:     ");

  if (ARM_FaultInfo.Content.StateContext != 0U) {
    printf("0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.ReturnAddress);
  } else {
    printf("unknown (was not stacked)\n");
  }

  /* Output: Registers */
  /* Registers R4 .. R11 values might be either: stacked (if additional state context (TrustZone only)
     was stacked) or values as they were when fault handler started execution */
  printf("\n  Registers:\n");
  if (ARM_FaultInfo.Content.StateContext != 0U) {
    printf("   - R0:               0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.R0);
    printf("   - R1:               0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.R1);
    printf("   - R2:               0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.R2);
    printf("   - R3:               0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.R3);
  } else {
    printf("   - R0 .. R3:         unknown (were not stacked)\n");
  }

  /* Output: R4 .. R11 */
  printf("   - R4:               0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.R4);
  printf("   - R5:               0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.R5);
  printf("   - R6:               0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.R6);
  printf("   - R7:               0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.R7);
  printf("   - R8:               0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.R8);
  printf("   - R9:               0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.R9);
  printf("   - R10:              0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.R10);
  printf("   - R11:              0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.R11);

  if (ARM_FaultInfo.Content.StateContext != 0U) {
    printf("   - R12:              0x%08X\n",   (unsigned int)ARM_FaultInfo.Registers.R12);
    printf("   - LR:               0x%08X\n",   (unsigned int)ARM_FaultInfo.Registers.LR);
    printf("   - Return Address:   0x%08X\n",   (unsigned int)ARM_FaultInfo.Registers.ReturnAddress);
    printf("   - xPSR:             0x%08X\n\n", (unsigned int)ARM_FaultInfo.Registers.xPSR);
  } else {
    printf("   - R12:              unknown (was not stacked)\n");
    printf("   - LR:               unknown (was not stacked)\n");
    printf("   - Return Address:   unknown (was not stacked)\n");
    printf("   - xPSR:             unknown (was not stacked)\n");
  }

  printf("   - MSP:              0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.MSP);
  if (ARM_FaultInfo.Content.LimitRegs != 0U) {
    printf("   - MSPLIM:           0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.MSPLIM);
  }
  printf("   - PSP:              0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.PSP);
  if (ARM_FaultInfo.Content.LimitRegs != 0U) {
    printf("   - PSPLIM:           0x%08X\n", (unsigned int)ARM_FaultInfo.Registers.PSPLIM);
  }

  /* Output: Exception State */
  printf("\n  Exception State:\n");
  printf("   - xPSR:             0x%08X\n", (unsigned int)ARM_FaultInfo.ExceptionState.xPSR);
  printf("   - Exception Return: 0x%08X\n", (unsigned int)ARM_FaultInfo.ExceptionState.EXC_RETURN);
  printf("\n");

#if (ARM_FAULT_FAULT_REGS_EXIST != 0)
  /* Output: Fault Registers (if they exist) */
  if (ARM_FaultInfo.Content.FaultRegs != 0U) {
    printf("  Fault Registers:\n");

    printf("   - CFSR:             0x%08X\n", (unsigned int)ARM_FaultInfo.FaultRegisters.CFSR);
    printf("   - HFSR:             0x%08X\n", (unsigned int)ARM_FaultInfo.FaultRegisters.HFSR);
    printf("   - DFSR:             0x%08X\n", (unsigned int)ARM_FaultInfo.FaultRegisters.DFSR);
    printf("   - MMFAR:            0x%08X\n", (unsigned int)ARM_FaultInfo.FaultRegisters.MMFAR);
    printf("   - BFAR:             0x%08X\n", (unsigned int)ARM_FaultInfo.FaultRegisters.BFAR);
    printf("   - AFSR:             0x%08X\n", (unsigned int)ARM_FaultInfo.FaultRegisters.AFSR);

    if (ARM_FaultInfo.Content.SecureFaultRegs != 0U) {
      printf("   - SFSR:             0x%08X\n", (unsigned int)ARM_FaultInfo.FaultRegisters.SFSR);
      printf("   - SFAR:             0x%08X\n", (unsigned int)ARM_FaultInfo.FaultRegisters.SFAR);
    }

#if (ARM_FAULT_ARCH_ARMV8_1M_MAIN != 0)
    if (ARM_FaultInfo.Content.RAS_FaultReg != 0U) {
      printf("   - RFSR:             0x%08X\n", (unsigned int)ARM_FaultInfo.FaultRegisters.RFSR);
    }
#endif

    printf("\n");
  }
#else
  /* Output: Message if fault registers do not exist */
  printf("  Fault Registers do not exist!\n\n");
#endif
}
