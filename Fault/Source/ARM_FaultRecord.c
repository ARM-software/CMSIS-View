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

//lint -e46  "Suppress: field type should be _Bool, unsigned int or signed int [MISRA 2012 Rule 6.1, required]"
//lint -e750 "Suppress: local macro not referenced [MISRA 2012 Rule 2.5, advisory]"
//lint -e835 "Suppress: A zero has been given as left argument to operator '+'"

#include "ARM_Fault.h"
#include "EventRecorder.h"

// ARM Fault component number
#define EvtFault_No     0xEEU           // Component number for ARM Fault Record

// ARM Fault Event IDs
#define EvtFault_FaultInfo_Empty        EventID(EventLevelOp,    EvtFault_No, 0x00U)
#define EvtFault_FaultInfo_Invalid      EventID(EventLevelError, EvtFault_No, 0x01U)
#define EvtFault_HardFault_VECTTBL      EventID(EventLevelError, EvtFault_No, 0x02U)
#define EvtFault_HardFault_FORCED       EventID(EventLevelError, EvtFault_No, 0x04U)
#define EvtFault_HardFault_DEBUGEVT     EventID(EventLevelError, EvtFault_No, 0x06U)
#define EvtFault_MemManage_IACCVIOL     EventID(EventLevelError, EvtFault_No, 0x08U)
#define EvtFault_MemManage_DACCVIOL     EventID(EventLevelError, EvtFault_No, 0x0CU)
#define EvtFault_MemManage_MUNSTKERR    EventID(EventLevelError, EvtFault_No, 0x10U)
#define EvtFault_MemManage_MSTKERR      EventID(EventLevelError, EvtFault_No, 0x14U)
#define EvtFault_MemManage_MLSPERR      EventID(EventLevelError, EvtFault_No, 0x16U)
#define EvtFault_BusFault_IBUSERR       EventID(EventLevelError, EvtFault_No, 0x1AU)
#define EvtFault_BusFault_PRECISERR     EventID(EventLevelError, EvtFault_No, 0x1EU)
#define EvtFault_BusFault_IMPRECISERR   EventID(EventLevelError, EvtFault_No, 0x22U)
#define EvtFault_BusFault_UNSTKERR      EventID(EventLevelError, EvtFault_No, 0x26U)
#define EvtFault_BusFault_STKERR        EventID(EventLevelError, EvtFault_No, 0x2AU)
#define EvtFault_BusFault_LSPERR        EventID(EventLevelError, EvtFault_No, 0x2CU)
#define EvtFault_UsageFault_UNDEFINSTR  EventID(EventLevelError, EvtFault_No, 0x30U)
#define EvtFault_UsageFault_INVSTATE    EventID(EventLevelError, EvtFault_No, 0x32U)
#define EvtFault_UsageFault_INVPC       EventID(EventLevelError, EvtFault_No, 0x34U)
#define EvtFault_UsageFault_NOCP        EventID(EventLevelError, EvtFault_No, 0x36U)
#define EvtFault_UsageFault_STKOF       EventID(EventLevelError, EvtFault_No, 0x38U)
#define EvtFault_UsageFault_UNALIGNED   EventID(EventLevelError, EvtFault_No, 0x39U)
#define EvtFault_UsageFault_DIVBYZERO   EventID(EventLevelError, EvtFault_No, 0x3BU)
#define EvtFault_SecureFault_INVEP      EventID(EventLevelError, EvtFault_No, 0x3DU)
#define EvtFault_SecureFault_INVIS      EventID(EventLevelError, EvtFault_No, 0x41U)
#define EvtFault_SecureFault_INVER      EventID(EventLevelError, EvtFault_No, 0x45U)
#define EvtFault_SecureFault_AUVIOL     EventID(EventLevelError, EvtFault_No, 0x49U)
#define EvtFault_SecureFault_INVTRAN    EventID(EventLevelError, EvtFault_No, 0x4DU)
#define EvtFault_SecureFault_LSPERR     EventID(EventLevelError, EvtFault_No, 0x51U)
#define EvtFault_SecureFault_LSERR      EventID(EventLevelError, EvtFault_No, 0x55U)

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

// ARM_FaultRecord function ----------------------------------------------------

/**
  Output decoded fault information via EventRecorder.
  Should be called when system is running in normal operating mode with
  EventRecorder fully functional.
*/
void ARM_FaultRecord (void) {
  int8_t   fault_info_valid    = 1;
  int8_t   fault_info_magic_ok = 1;
  int8_t   fault_info_crc_ok   = 1;
#if (ARM_FAULT_FAULT_REGS_EXIST != 0)
  int8_t   state_context_valid = 1;
  uint32_t return_address      = 0U;
  uint32_t evt_id_inc          = 0U;
#endif

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
  } else {
    return_address = ARM_FaultInfo.ReturnAddress;
  }
#endif

  // Output: Error message if magic number or CRC is invalid
  if (fault_info_magic_ok == 0) {
    //lint -e845 "Suppress: The right argument to operator '|' is certain to be 0"
    (void)EventRecord2(EvtFault_FaultInfo_Empty,   0U, 0U);
  } else if (fault_info_crc_ok == 0) {
    (void)EventRecord2(EvtFault_FaultInfo_Invalid, 0U, 0U);
  } else {
    // Fault information is valid
  }

#if (ARM_FAULT_FAULT_REGS_EXIST != 0)
  /* Output: Decoded HardFault information */
  if ((fault_info_valid != 0) && (ARM_FaultInfo.type.fault_regs != 0U)) {
    uint32_t scb_hfsr = ARM_FaultInfo.SCB_HFSR;

    if ((scb_hfsr & (SCB_HFSR_VECTTBL_Msk   |
                     SCB_HFSR_FORCED_Msk    |
                     SCB_HFSR_DEBUGEVT_Msk  )) != 0U) {

      if (state_context_valid == 0) {
        evt_id_inc = 0U;
      } else {
        evt_id_inc = 1U;
      }

      if ((scb_hfsr & SCB_HFSR_VECTTBL_Msk) != 0U) {
        (void)EventRecord2(EvtFault_HardFault_VECTTBL  + evt_id_inc, return_address, 0U);
      }
      if ((scb_hfsr & SCB_HFSR_FORCED_Msk) != 0U) {
        (void)EventRecord2(EvtFault_HardFault_FORCED   + evt_id_inc, return_address, 0U);
      }
      if ((scb_hfsr & SCB_HFSR_DEBUGEVT_Msk) != 0U) {
        (void)EventRecord2(EvtFault_HardFault_DEBUGEVT + evt_id_inc, return_address, 0U);
      }
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

      evt_id_inc = 0U;
      if (state_context_valid != 0) {
        evt_id_inc += 1U;
      }
      if ((scb_cfsr & SCB_CFSR_MMARVALID_Msk) != 0U) {
        evt_id_inc += 2U;
      }

      if ((scb_cfsr & SCB_CFSR_IACCVIOL_Msk) != 0U) {
        (void)EventRecord2(EvtFault_MemManage_IACCVIOL  + evt_id_inc, return_address, scb_mmfar);
      }
      if ((scb_cfsr & SCB_CFSR_DACCVIOL_Msk) != 0U) {
        (void)EventRecord2(EvtFault_MemManage_DACCVIOL  + evt_id_inc, return_address, scb_mmfar);
      }
      if ((scb_cfsr & SCB_CFSR_MUNSTKERR_Msk) != 0U) {
        (void)EventRecord2(EvtFault_MemManage_MUNSTKERR + evt_id_inc, return_address, scb_mmfar);
      }
      if ((scb_cfsr & SCB_CFSR_MSTKERR_Msk) != 0U) {
        (void)EventRecord2(EvtFault_MemManage_MSTKERR   + evt_id_inc, return_address, scb_mmfar);
      }
#ifdef SCB_CFSR_MLSPERR_Msk
      if ((scb_cfsr & SCB_CFSR_MLSPERR_Msk) != 0U) {
        (void)EventRecord2(EvtFault_MemManage_MLSPERR   + evt_id_inc, return_address, scb_mmfar);
      }
#endif
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

      evt_id_inc = 0U;
      if (state_context_valid != 0) {
        evt_id_inc += 1U;
      }
      if ((scb_cfsr & SCB_CFSR_BFARVALID_Msk) != 0U) {
        evt_id_inc += 2U;
      }

      if ((scb_cfsr & SCB_CFSR_IBUSERR_Msk) != 0U) {
        (void)EventRecord2(EvtFault_BusFault_IBUSERR     + evt_id_inc, return_address, scb_bfar);
      }
      if ((scb_cfsr & SCB_CFSR_PRECISERR_Msk) != 0U) {
        (void)EventRecord2(EvtFault_BusFault_PRECISERR   + evt_id_inc, return_address, scb_bfar);
      }
      if ((scb_cfsr & SCB_CFSR_IMPRECISERR_Msk) != 0U) {
        (void)EventRecord2(EvtFault_BusFault_IMPRECISERR + evt_id_inc, return_address, scb_bfar);
      }
      if ((scb_cfsr & SCB_CFSR_UNSTKERR_Msk) != 0U) {
        (void)EventRecord2(EvtFault_BusFault_UNSTKERR    + evt_id_inc, return_address, scb_bfar);
      }
      if ((scb_cfsr & SCB_CFSR_STKERR_Msk) != 0U) {
        (void)EventRecord2(EvtFault_BusFault_STKERR      + evt_id_inc, return_address, scb_bfar);
      }
#ifdef SCB_CFSR_LSPERR_Msk
      if ((scb_cfsr & SCB_CFSR_LSPERR_Msk) != 0U) {
        (void)EventRecord2(EvtFault_BusFault_LSPERR   + evt_id_inc, return_address, scb_bfar);
      }
#endif
      if ((scb_cfsr & SCB_CFSR_BFARVALID_Msk) != 0U) {
        (void)EventRecord2(EvtFault_BusFault_IBUSERR   + evt_id_inc, return_address, scb_bfar);
      }
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

      if (state_context_valid == 0) {
        evt_id_inc = 0U;
      } else {
        evt_id_inc = 1U;
      }

      if ((scb_cfsr & SCB_CFSR_UNDEFINSTR_Msk) != 0U) {
        (void)EventRecord2(EvtFault_UsageFault_UNDEFINSTR + evt_id_inc, return_address, 0U);
      }
      if ((scb_cfsr & SCB_CFSR_INVSTATE_Msk) != 0U) {
        (void)EventRecord2(EvtFault_UsageFault_INVSTATE   + evt_id_inc, return_address, 0U);
      }
      if ((scb_cfsr & SCB_CFSR_INVPC_Msk) != 0U) {
        (void)EventRecord2(EvtFault_UsageFault_INVPC      + evt_id_inc, return_address, 0U);
      }
      if ((scb_cfsr & SCB_CFSR_NOCP_Msk) != 0U) {
        (void)EventRecord2(EvtFault_UsageFault_NOCP       + evt_id_inc, return_address, 0U);
      }
#ifdef SCB_CFSR_STKOF_Msk
      if ((scb_cfsr & SCB_CFSR_STKOF_Msk) != 0U) {
        (void)EventRecord2(EvtFault_UsageFault_STKOF      + evt_id_inc, return_address, 0U);
      }
#endif
      if ((scb_cfsr & SCB_CFSR_UNALIGNED_Msk) != 0U) {
        (void)EventRecord2(EvtFault_UsageFault_UNALIGNED  + evt_id_inc, return_address, 0U);
      }
      if ((scb_cfsr & SCB_CFSR_DIVBYZERO_Msk) != 0U) {
        (void)EventRecord2(EvtFault_UsageFault_DIVBYZERO  + evt_id_inc, return_address, 0U);
      }
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

      evt_id_inc = 0U;
      if (state_context_valid != 0) {
        evt_id_inc += 1U;
      }
      if ((scb_sfsr & SAU_SFSR_SFARVALID_Msk) != 0U) {
        evt_id_inc += 2U;
      }

      if ((scb_sfsr & SAU_SFSR_INVEP_Msk) != 0U) {
        (void)EventRecord2(EvtFault_SecureFault_INVEP   + evt_id_inc, return_address, scb_sfar);
      }
      if ((scb_sfsr & SAU_SFSR_INVIS_Msk) != 0U) {
        (void)EventRecord2(EvtFault_SecureFault_INVIS   + evt_id_inc, return_address, scb_sfar);
      }
      if ((scb_sfsr & SAU_SFSR_INVER_Msk) != 0U) {
        (void)EventRecord2(EvtFault_SecureFault_INVER   + evt_id_inc, return_address, scb_sfar);
      }
      if ((scb_sfsr & SAU_SFSR_AUVIOL_Msk) != 0U) {
        (void)EventRecord2(EvtFault_SecureFault_AUVIOL  + evt_id_inc, return_address, scb_sfar);
      }
      if ((scb_sfsr & SAU_SFSR_INVTRAN_Msk) != 0U) {
        (void)EventRecord2(EvtFault_SecureFault_INVTRAN + evt_id_inc, return_address, scb_sfar);
      }
      if ((scb_sfsr & SAU_SFSR_LSPERR_Msk) != 0U) {
        (void)EventRecord2(EvtFault_SecureFault_LSPERR  + evt_id_inc, return_address, scb_sfar);
      }
      if ((scb_sfsr & SAU_SFSR_LSERR_Msk) != 0U) {
        (void)EventRecord2(EvtFault_SecureFault_LSERR   + evt_id_inc, return_address, scb_sfar);
      }
    }
  }
#endif
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
