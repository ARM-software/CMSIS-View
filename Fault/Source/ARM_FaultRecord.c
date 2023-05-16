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
//lint -esym(9058, cmse_address_info) "Suppress: tag 'cmse_address_info' unused outside of typedefs [MISRA 2012 Rule 2.4, advisory]"

#include "ARM_Fault.h"
#include "EventRecorder.h"

// ARM Fault component number
#define EvtFault_No     0xEEU           // Component number for ARM Fault

// ARM Fault Event IDs
#define EvtFault_FaultInfo_Invalid      EventID(EventLevelOp,    EvtFault_No, 0x00U)
#define EvtFault_FaultInfo_NoFaultRegs  EventID(EventLevelError, EvtFault_No, 0x01U)
#define EvtFault_HardFault_VECTTBL      EventID(EventLevelError, EvtFault_No, 0x03U)
#define EvtFault_HardFault_FORCED       EventID(EventLevelError, EvtFault_No, 0x05U)
#define EvtFault_HardFault_DEBUGEVT     EventID(EventLevelError, EvtFault_No, 0x07U)
#define EvtFault_MemManage_IACCVIOL     EventID(EventLevelError, EvtFault_No, 0x09U)
#define EvtFault_MemManage_DACCVIOL     EventID(EventLevelError, EvtFault_No, 0x0DU)
#define EvtFault_MemManage_MUNSTKERR    EventID(EventLevelError, EvtFault_No, 0x11U)
#define EvtFault_MemManage_MSTKERR      EventID(EventLevelError, EvtFault_No, 0x15U)
#define EvtFault_MemManage_MLSPERR      EventID(EventLevelError, EvtFault_No, 0x17U)
#define EvtFault_BusFault_IBUSERR       EventID(EventLevelError, EvtFault_No, 0x1BU)
#define EvtFault_BusFault_PRECISERR     EventID(EventLevelError, EvtFault_No, 0x1FU)
#define EvtFault_BusFault_IMPRECISERR   EventID(EventLevelError, EvtFault_No, 0x23U)
#define EvtFault_BusFault_UNSTKERR      EventID(EventLevelError, EvtFault_No, 0x27U)
#define EvtFault_BusFault_STKERR        EventID(EventLevelError, EvtFault_No, 0x2BU)
#define EvtFault_BusFault_LSPERR        EventID(EventLevelError, EvtFault_No, 0x2DU)
#define EvtFault_UsageFault_UNDEFINSTR  EventID(EventLevelError, EvtFault_No, 0x31U)
#define EvtFault_UsageFault_INVSTATE    EventID(EventLevelError, EvtFault_No, 0x33U)
#define EvtFault_UsageFault_INVPC       EventID(EventLevelError, EvtFault_No, 0x35U)
#define EvtFault_UsageFault_NOCP        EventID(EventLevelError, EvtFault_No, 0x37U)
#define EvtFault_UsageFault_STKOF       EventID(EventLevelError, EvtFault_No, 0x39U)
#define EvtFault_UsageFault_UNALIGNED   EventID(EventLevelError, EvtFault_No, 0x3AU)
#define EvtFault_UsageFault_DIVBYZERO   EventID(EventLevelError, EvtFault_No, 0x3CU)
#define EvtFault_SecureFault_INVEP      EventID(EventLevelError, EvtFault_No, 0x3EU)
#define EvtFault_SecureFault_INVIS      EventID(EventLevelError, EvtFault_No, 0x42U)
#define EvtFault_SecureFault_INVER      EventID(EventLevelError, EvtFault_No, 0x46U)
#define EvtFault_SecureFault_AUVIOL     EventID(EventLevelError, EvtFault_No, 0x4AU)
#define EvtFault_SecureFault_INVTRAN    EventID(EventLevelError, EvtFault_No, 0x4EU)
#define EvtFault_SecureFault_LSPERR     EventID(EventLevelError, EvtFault_No, 0x52U)
#define EvtFault_SecureFault_LSERR      EventID(EventLevelError, EvtFault_No, 0x56U)
#define EvtFault_RAS_Fault              EventID(EventLevelError, EvtFault_No, 0x5AU)
#define EvtFault_TZ_Info                EventID(EventLevelOp,    EvtFault_No, 0x5CU)

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

// ARM_FaultRecord function ----------------------------------------------------

/**
  Output decoded fault information via EventRecorder.
  Should be called when system is running in normal operating mode with
  EventRecorder fully functional.
*/
void ARM_FaultRecord (void) {
  int8_t fault_info_valid;
  uint32_t return_address = 0U;
  uint32_t evt_id_inc     = 0U;

  /* Check if there is available valid fault information */
  fault_info_valid = (int8_t)ARM_FaultOccurred();

  // Check if state context is valid
  if (ARM_FaultInfo.Content.StateContext != 0U) {
    return_address = ARM_FaultInfo.Registers.ReturnAddress;
  }

  // Output: Message if fault info is invalid
  if (fault_info_valid == 0) {
    //lint -e845 "Suppress: The right argument to operator '|' is certain to be 0"
    (void)EventRecord2(EvtFault_FaultInfo_Invalid, 0U, 0U);
  } else {
    // Fault information is valid
  }

#if (ARM_FAULT_TZ_ENABLED != 0)         // If TrustZone is enabled
  if ((fault_info_valid != 0) && (ARM_FaultInfo.Content.TZ_Enabled != 0U)) {
    (void)EventRecord2(EvtFault_TZ_Info, ARM_FaultInfo.Content.TZ_FaultMode, ARM_FaultInfo.Content.TZ_SaveMode);
  }
#endif

#if (ARM_FAULT_FAULT_REGS_EXIST == 0)   // If fault registers do not exist
  /* Output: Message if fault registers do not exist */
  if (fault_info_valid != 0) {
    if (ARM_FaultInfo.Content.StateContext == 0U) {
      evt_id_inc = 0U;
    } else {
      evt_id_inc = 1U;
    }

    (void)EventRecord2(EvtFault_FaultInfo_NoFaultRegs + evt_id_inc, return_address, 0U);
  }
#endif

#if (ARM_FAULT_FAULT_REGS_EXIST != 0)   // If fault registers exist
  /* Output: Decoded HardFault information */
  if ((fault_info_valid != 0) && (ARM_FaultInfo.Content.FaultRegs != 0U)) {
    uint32_t scb_hfsr = ARM_FaultInfo.FaultRegisters.HFSR;

    if ((scb_hfsr & (SCB_HFSR_VECTTBL_Msk   |
                     SCB_HFSR_FORCED_Msk    |
                     SCB_HFSR_DEBUGEVT_Msk  )) != 0U) {

      if (ARM_FaultInfo.Content.StateContext == 0U) {
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
  if ((fault_info_valid != 0) && (ARM_FaultInfo.Content.FaultRegs != 0U)) {
    uint32_t scb_cfsr  = ARM_FaultInfo.FaultRegisters.CFSR;
    uint32_t scb_mmfar = ARM_FaultInfo.FaultRegisters.MMFAR;

    if ((scb_cfsr & (SCB_CFSR_IACCVIOL_Msk  |
                     SCB_CFSR_DACCVIOL_Msk  |
                     SCB_CFSR_MUNSTKERR_Msk |
#ifdef SCB_CFSR_MLSPERR_Msk
                     SCB_CFSR_MLSPERR_Msk   |
#endif
                     SCB_CFSR_MSTKERR_Msk   )) != 0U) {

      evt_id_inc = 0U;
      if (ARM_FaultInfo.Content.StateContext != 0U) {
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
  if ((fault_info_valid != 0) && (ARM_FaultInfo.Content.FaultRegs != 0U)) {
    uint32_t scb_cfsr = ARM_FaultInfo.FaultRegisters.CFSR;
    uint32_t scb_bfar = ARM_FaultInfo.FaultRegisters.BFAR;

    if ((scb_cfsr & (SCB_CFSR_IBUSERR_Msk     |
                     SCB_CFSR_PRECISERR_Msk   |
                     SCB_CFSR_IMPRECISERR_Msk |
                     SCB_CFSR_UNSTKERR_Msk    |
#ifdef SCB_CFSR_LSPERR_Msk
                     SCB_CFSR_LSPERR_Msk      |
#endif
                     SCB_CFSR_STKERR_Msk      )) != 0U) {

      evt_id_inc = 0U;
      if (ARM_FaultInfo.Content.StateContext != 0U) {
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
        (void)EventRecord2(EvtFault_BusFault_LSPERR      + evt_id_inc, return_address, scb_bfar);
      }
#endif
    }
  }

  /* Output Decoded UsageFault information */
  if ((fault_info_valid != 0) && (ARM_FaultInfo.Content.FaultRegs != 0U)) {
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

      if (ARM_FaultInfo.Content.StateContext == 0U) {
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
  if ((fault_info_valid != 0) && (ARM_FaultInfo.Content.SecureFaultRegs != 0U)) {
    uint32_t scb_sfsr = ARM_FaultInfo.FaultRegisters.SFSR;
    uint32_t scb_sfar = ARM_FaultInfo.FaultRegisters.SFAR;

    if ((scb_sfsr & (SAU_SFSR_INVEP_Msk   |
                     SAU_SFSR_INVIS_Msk   |
                     SAU_SFSR_INVER_Msk   |
                     SAU_SFSR_AUVIOL_Msk  |
                     SAU_SFSR_INVTRAN_Msk |
                     SAU_SFSR_LSPERR_Msk  |
                     SAU_SFSR_LSERR_Msk   )) != 0U) {

      evt_id_inc = 0U;
      if (ARM_FaultInfo.Content.StateContext != 0U) {
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

#if (ARM_FAULT_ARCH_ARMV8_1M_MAIN != 0)
  /* Output: RAS Fault information */
  if ((fault_info_valid != 0) && (ARM_FaultInfo.Content.RAS_FaultReg != 0U)) {
    uint32_t scb_rfsr = ARM_FaultInfo.FaultRegisters.RFSR;

    if ((scb_rfsr & SCB_RFSR_V_Msk) != 0U) {

      evt_id_inc = 0U;
      if (ARM_FaultInfo.Content.StateContext != 0U) {
        evt_id_inc += 1U;
      }

      (void)EventRecord2(EvtFault_RAS_Fault + evt_id_inc, return_address, scb_rfsr);
    }
  }
#endif
#endif
}
