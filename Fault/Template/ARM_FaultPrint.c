/*------------------------------------------------------------------------------
 * MDK - Component ARM::CMSIS-View:Fault:Storage
 * Copyright (c) 2022-2023 ARM Germany GmbH. All rights reserved.
 *------------------------------------------------------------------------------
 * Name:    ARM_FaultPrint.c
 * Purpose: Output decoded fault information via STDIO
 * Rev.:    V0.4.0
 *----------------------------------------------------------------------------*/

#include "ARM_Fault.h"

#include <stdio.h>

#if    (ARM_FAULT_FAULT_REGS_EXIST != 0)
// Define CFSR mask for detecting state context stacking failure
#ifndef SCB_CFSR_Stack_Err_Msk
#ifdef  SCB_CFSR_STKOF_Msk
#define SCB_CFSR_Stack_Err_Msk (SCB_CFSR_STKERR_Msk | SCB_CFSR_MSTKERR_Msk | SCB_CFSR_STKOF_Msk)
#else
#define SCB_CFSR_Stack_Err_Msk (SCB_CFSR_STKERR_Msk | SCB_CFSR_MSTKERR_Msk)
#endif
#endif
#endif

// General defines
#ifndef EXC_RETURN_SPSEL
#define EXC_RETURN_SPSEL       (1UL << 2)
#endif

// Armv8/8.1-M architecture related defines
#if    (ARM_FAULT_ARCH_ARMV8x_M != 0)
#define ARM_FAULT_ASC_INTEGRITY_SIG    (0xFEFA125AU)    // Additional State Context Integrity Signature
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

// Local functions prototypes
static uint32_t CalcCRC32 (      uint32_t init_val,
                           const uint8_t *data_ptr,
                                 uint32_t data_len,
                                 uint32_t polynom);

// ARM_FaultPrint function -----------------------------------------------------

/**
  Output decoded fault information via STDIO.
  Should be called when system is running in normal operating mode with
  standard input/output fully functional.
*/
void ARM_FaultPrint (void) {
  int8_t fault_info_valid    = 1;
  int8_t fault_info_magic_ok = 1;
  int8_t fault_info_crc_ok   = 1;
  int8_t state_context_valid = 1;

  // Check if magic number is valid
  if (ARM_FaultInfo.magic_number != ARM_FAULT_MAGIC_NUMBER) {
    fault_info_valid    = 0;
    fault_info_magic_ok = 0;
  }

  // Check if CRC of the ARM_FaultInfo is correct
  if (fault_info_valid != 0) {
    if (ARM_FaultInfo.crc32 != CalcCRC32(ARM_FAULT_CRC32_INIT_VAL,
                                        (const uint8_t *)&ARM_FaultInfo.type,
                                        (sizeof(ARM_FaultInfo) - (sizeof(ARM_FaultInfo.magic_number) + sizeof(ARM_FaultInfo.crc32))),
                                         ARM_FAULT_CRC32_POLYNOM)) {
      fault_info_valid  = 0;
      fault_info_crc_ok = 0;
    }
  }

#if (ARM_FAULT_FAULT_REGS_EXIST != 0)
  // Check if the state context was stacked properly (if CFSR is available)
  if ((ARM_FaultInfo.SCB_CFSR & (SCB_CFSR_Stack_Err_Msk)) != 0U) {
    state_context_valid = 0;
  }
#endif

  // Output: Header and version information
  printf("\n --- FaultRecorder (v%s) ---\n\n", (const char *)ARM_FaultVersion);

  // Output: Error message if magic number or CRC is invalid
  if (fault_info_magic_ok == 0) {
    printf("\n  No fault saved yet!\n\n");
  } else if (fault_info_crc_ok == 0) {
    printf("\n  Invalid CRC of the saved fault information!\n\n");
  } else {
    // Fault information is valid
  }

  // Output: Fault count
  if (fault_info_valid != 0) {
    printf("  Fault count:       %u\n\n", ARM_FaultInfo.count);
  }

  // Output: Exception which recorded the fault information
  if (fault_info_valid != 0) {
    uint32_t exc_num = ARM_FaultInfo.xPSR_in_handler & IPSR_ISR_Msk;

    printf("  Exception Handler: ");

#if (ARM_FAULT_ARCH_ARMV8x_M != 0)
    if (ARM_FaultInfo.type.tz_secure != 0U) {
      printf("Secure - ");
    } else {
      printf("Non-Secure - ");
    }
#endif

    switch (exc_num) {
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
        printf("unknown, exception number = %u", exc_num);
        break;
    }

    printf("\n");
  }

#if (ARM_FAULT_ARCH_ARMV8x_M != 0)
  // Output: state in which the fault occurred
  if (fault_info_valid != 0) {
    uint32_t exc_return = ARM_FaultInfo.EXC_RETURN;

    printf("  State:             ");

    if ((exc_return & EXC_RETURN_S) != 0U) {
      printf("Secure");
    } else {
      printf("Non-Secure");
    }

    printf("\n");
  }
#endif

  // Output: Mode in which the fault occurred
  if (fault_info_valid != 0) {
    uint32_t exc_return = ARM_FaultInfo.EXC_RETURN;

    printf("  Mode:              ");

    if ((exc_return & EXC_RETURN_SPSEL) == 0U) {
      printf("Handler");
    } else {
      printf("Thread");
    }

    printf("\n");
  }

#if (ARM_FAULT_FAULT_REGS_EXIST != 0)
  /* Output: Decoded HardFault information */
  if ((fault_info_valid != 0) && (ARM_FaultInfo.type.fault_regs != 0U)) {
    uint32_t scb_hfsr = ARM_FaultInfo.SCB_HFSR;

    if ((scb_hfsr & (SCB_HFSR_VECTTBL_Msk   |
                     SCB_HFSR_FORCED_Msk    |
                     SCB_HFSR_DEBUGEVT_Msk  )) != 0U) {

      printf("  Fault:             HardFault - ");

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
  if ((fault_info_valid != 0) && (ARM_FaultInfo.type.fault_regs != 0U)) {
    uint32_t scb_cfsr  = ARM_FaultInfo.SCB_CFSR;
    uint32_t scb_mmfar = ARM_FaultInfo.SCB_MMFAR;

    if ((scb_cfsr & (SCB_CFSR_IACCVIOL_Msk  |
                     SCB_CFSR_DACCVIOL_Msk  |
                     SCB_CFSR_MUNSTKERR_Msk |
#ifdef SCB_CFSR_MLSPERR_Msk
                     SCB_CFSR_MLSPERR_Msk   |
#endif
                     SCB_CFSR_MSTKERR_Msk   )) != 0U) {

      printf("  Fault:             MemManage - ");

      if ((scb_cfsr & SCB_CFSR_IACCVIOL_Msk) != 0U) {
        printf("Instruction execution failure due to MPU violation or fault");
      }
      if ((scb_cfsr & SCB_CFSR_DACCVIOL_Msk) != 0U) {
        printf("Data access failure due to MPU violation or fault");
      }
      if ((scb_cfsr & SCB_CFSR_MUNSTKERR_Msk) != 0U) {
        printf("Exception exit unstacking failure due to MPU access violation");
      }
      if ((scb_cfsr & SCB_CFSR_MSTKERR_Msk) != 0U) {
        printf("Exception entry stacking failure due to MPU access violation");
      }
#ifdef SCB_CFSR_MLSPERR_Msk
      if ((scb_cfsr & SCB_CFSR_MLSPERR_Msk) != 0U) {
        printf("Floating-point lazy stacking failure due to MPU access violation");
      }
#endif
      if ((scb_cfsr & SCB_CFSR_MMARVALID_Msk) != 0U) {
        printf(", fault address 0x%08X", scb_mmfar);
      }

      printf("\n");
    }
  }

  /* Output: Decoded BusFault information */
  if ((fault_info_valid != 0) && (ARM_FaultInfo.type.fault_regs != 0U)) {
    uint32_t scb_cfsr = ARM_FaultInfo.SCB_CFSR;
    uint32_t scb_bfar = ARM_FaultInfo.SCB_BFAR;

    if ((scb_cfsr & (SCB_CFSR_IBUSERR_Msk     |
                     SCB_CFSR_PRECISERR_Msk   |
                     SCB_CFSR_IMPRECISERR_Msk |
                     SCB_CFSR_UNSTKERR_Msk    |
#ifdef SCB_CFSR_LSPERR_Msk
                     SCB_CFSR_LSPERR_Msk      |
#endif
                     SCB_CFSR_STKERR_Msk      )) != 0U) {

      printf("  Fault:             BusFault - ");

      if ((scb_cfsr & SCB_CFSR_IBUSERR_Msk) != 0U) {
        printf("Instruction prefetch failure due to bus fault");
      }
      if ((scb_cfsr & SCB_CFSR_PRECISERR_Msk) != 0U) {
        printf("Data access failure due to bus fault (precise)");
      }
      if ((scb_cfsr & SCB_CFSR_IMPRECISERR_Msk) != 0U) {
        printf("Data access failure due to bus fault (imprecise)");
      }
      if ((scb_cfsr & SCB_CFSR_UNSTKERR_Msk) != 0U) {
        printf("Exception exit unstacking failure due to bus fault");
      }
      if ((scb_cfsr & SCB_CFSR_STKERR_Msk) != 0U) {
        printf("Exception entry stacking failure due to bus fault");
      }
#ifdef SCB_CFSR_LSPERR_Msk
      if ((scb_cfsr & SCB_CFSR_LSPERR_Msk) != 0U) {
        printf("Floating-point lazy stacking failure due to bus fault");
      }
#endif
      if ((scb_cfsr & SCB_CFSR_BFARVALID_Msk) != 0U) {
        printf(", fault address 0x%08X", scb_bfar);
      }

      printf("\n");
    }
  }

  /* Output Decoded UsageFault information */
  if ((fault_info_valid != 0) && (ARM_FaultInfo.type.fault_regs != 0U)) {
    uint32_t scb_cfsr = ARM_FaultInfo.SCB_CFSR;

    if ((scb_cfsr & (SCB_CFSR_UNDEFINSTR_Msk |
                     SCB_CFSR_INVSTATE_Msk   |
                     SCB_CFSR_INVPC_Msk      |
                     SCB_CFSR_NOCP_Msk       |
#ifdef SCB_CFSR_STKOF_Msk
                     SCB_CFSR_STKOF_Msk      |
#endif
                     SCB_CFSR_UNALIGNED_Msk  |
                     SCB_CFSR_DIVBYZERO_Msk  )) != 0U) {

      printf("  Fault:             UsageFault - ");

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
  if ((fault_info_valid != 0) && (ARM_FaultInfo.type.fault_regs != 0U)) {
    uint32_t scb_sfsr = ARM_FaultInfo.SCB_SFSR;
    uint32_t scb_sfar = ARM_FaultInfo.SCB_SFAR;

    if ((scb_sfsr & (SAU_SFSR_INVEP_Msk   |
                     SAU_SFSR_INVIS_Msk   |
                     SAU_SFSR_INVER_Msk   |
                     SAU_SFSR_AUVIOL_Msk  |
                     SAU_SFSR_INVTRAN_Msk |
                     SAU_SFSR_LSPERR_Msk  |
                     SAU_SFSR_LSERR_Msk   )) != 0U) {

      printf("  Fault:             SecureFault - ");

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
        printf(", fault address 0x%08X", scb_sfar);
      }

      printf("\n");
    }
  }
#endif
#endif

  // Output: Program Counter, MSP (if TrustZone also MSPLIM), PSP (if TrustZone also PSPLIM)
  if (fault_info_valid != 0) {

    printf("\n");

#if (ARM_FAULT_FAULT_REGS_EXIST != 0)
    printf("   - Return Address: ");

    if (state_context_valid != 0) {
      printf("0x%08X\n", ARM_FaultInfo.ReturnAddress);
    } else {
      printf("unknown\n");
    }
#else
    printf("   - Return Address: 0x%08X\n", ARM_FaultInfo.ReturnAddress);
#endif
    printf("   - MSP:            0x%08X\n", ARM_FaultInfo.MSP);
#if (ARM_FAULT_ARCH_ARMV8x_M     != 0)
#if (ARM_FAULT_ARCH_ARMV8_M_BASE != 0)
    if ((ARM_FaultInfo.EXC_RETURN & EXC_RETURN_S) != 0) {
      printf("   - MSPLIM:         0x%08X\n", ARM_FaultInfo.MSPLIM);
    }
#else
    printf("   - MSPLIM:         0x%08X\n", ARM_FaultInfo.MSPLIM);
#endif
#endif
    printf("   - PSP:            0x%08X\n", ARM_FaultInfo.PSP);
#if (ARM_FAULT_ARCH_ARMV8x_M     != 0)
#if (ARM_FAULT_ARCH_ARMV8_M_BASE != 0)
    if ((ARM_FaultInfo.EXC_RETURN & EXC_RETURN_S) != 0) {
      printf("   - PSPLIM:         0x%08X\n", ARM_FaultInfo.PSPLIM);
    }
#else
    printf("   - PSPLIM:         0x%08X\n", ARM_FaultInfo.PSPLIM);
#endif
#endif

    printf("\n");
  }

  /* Output: state context information */
  if ((fault_info_valid != 0) && (state_context_valid != 0))  {
    printf("  Exception stacked State Context:\n");

    printf("   - R0:             0x%08X\n", ARM_FaultInfo.R0);
    printf("   - R1:             0x%08X\n", ARM_FaultInfo.R1);
    printf("   - R2:             0x%08X\n", ARM_FaultInfo.R2);
    printf("   - R3:             0x%08X\n", ARM_FaultInfo.R3);
  }

#if (ARM_FAULT_ARCH_ARMV8x_M != 0)
  if ((fault_info_valid != 0) && (state_context_valid != 0) && (ARM_FaultInfo.type.armv8m != 0U))  {
    /* Output: additional state context (if it exists) */
    if ((ARM_FaultInfo.IntegritySignature & 0xFFFFFFFEU) == ARM_FAULT_ASC_INTEGRITY_SIG) {
      printf("   - R4:             0x%08X\n", ARM_FaultInfo.R4);
      printf("   - R5:             0x%08X\n", ARM_FaultInfo.R5);
      printf("   - R6:             0x%08X\n", ARM_FaultInfo.R6);
      printf("   - R7:             0x%08X\n", ARM_FaultInfo.R7);
      printf("   - R8:             0x%08X\n", ARM_FaultInfo.R8);
      printf("   - R9:             0x%08X\n", ARM_FaultInfo.R9);
      printf("   - R10:            0x%08X\n", ARM_FaultInfo.R10);
      printf("   - R11:            0x%08X\n", ARM_FaultInfo.R11);
    }
  }
#endif

  if ((fault_info_valid != 0) && (state_context_valid != 0))  {
    printf("   - R12:            0x%08X\n",   ARM_FaultInfo.R12);
    printf("   - LR:             0x%08X\n",   ARM_FaultInfo.LR);
    printf("   - Return Address: 0x%08X\n",   ARM_FaultInfo.ReturnAddress);
    printf("   - xPSR:           0x%08X\n\n", ARM_FaultInfo.xPSR);
  }

#if (ARM_FAULT_FAULT_REGS_EXIST  != 0)
  /* Output: Fault registers (if they exist) */
  if (fault_info_valid != 0) {
    printf("  Fault Registers:\n");

    printf("   - CFSR:           0x%08X\n", ARM_FaultInfo.SCB_CFSR);
    printf("   - HFSR:           0x%08X\n", ARM_FaultInfo.SCB_HFSR);
    printf("   - DFSR:           0x%08X\n", ARM_FaultInfo.SCB_DFSR);
    printf("   - MMFAR:          0x%08X\n", ARM_FaultInfo.SCB_MMFAR);
    printf("   - BFAR:           0x%08X\n", ARM_FaultInfo.SCB_BFAR);
    printf("   - AFSR:           0x%08X\n", ARM_FaultInfo.SCB_AFSR);

#if (ARM_FAULT_ARCH_ARMV8x_M_MAIN != 0)
    if (ARM_FaultInfo.type.tz_secure != 0U) {
      printf("   - SFSR:           0x%08X\n", ARM_FaultInfo.SCB_SFSR);
      printf("   - SFAR:           0x%08X\n", ARM_FaultInfo.SCB_SFAR);
    }
#endif

    printf("\n");
  }
#endif
}


// Helper functions

/**
  Calculate CRC-32 on data block in memory
  \param[in]    init_val        initial CRC value
  \param[in]    data_ptr        pointer to data
  \param[in]    data_len        data length (in bytes)
  \param[in]    polynom         CRC polynom
  \return       CRC-32 value (32-bit)
*/
static uint32_t CalcCRC32 (      uint32_t init_val,
                           const uint8_t *data_ptr,
                                 uint32_t data_len,
                                 uint32_t polynom) {
  uint32_t crc32, i;

  crc32 = init_val;
  while (data_len != 0U) {
    crc32 ^= ((uint32_t)*data_ptr) << 24;
    for (i = 8U; i != 0U; i--) {
      if ((crc32 & 0x80000000U) != 0U) {
        crc32 <<= 1;
        crc32  ^= polynom;
      } else {
        crc32 <<= 1;
      }
    }
    data_ptr++;
    data_len--;
  }

  return crc32;
}
