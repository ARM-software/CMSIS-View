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
//lint -e451 "Suppress: repeatedly included but does not have a standard include guard [MISRA 2012 Directive 4.10, required]"
//lint -e537 "Suppress: Repeated include file 'stddef.h'"

#include "ARM_Fault.h"

#include <stddef.h>
#include <string.h>

// Compiler-specific defines
#if !defined(__NAKED)
  //lint -esym(9071, __NAKED) "Suppress: defined macro is reserved to the compiler"
  #define __NAKED __attribute__((naked))
#endif
#if !defined(__WEAK)
  //lint -esym(9071, __WEAK) "Suppress: defined macro is reserved to the compiler"
  #define __WEAK __attribute__((weak))
#endif
#if !defined(__NO_INIT)
  //lint -esym(9071, __NO_INIT) "Suppress: defined macro is reserved to the compiler"
  #if   defined (__CC_ARM)                                           /* ARM Compiler 4/5 */
    #define __NO_INIT __attribute__ ((section (".bss.noinit"), zero_init))
  #elif defined (__ARMCC_VERSION) && (__ARMCC_VERSION >= 6010050)    /* ARM Compiler 6 */
    #define __NO_INIT __attribute__ ((section (".bss.noinit")))
  #elif defined (__GNUC__)                                           /* GNU Compiler */
    #define __NO_INIT __attribute__ ((section (".noinit")))
  #else
    #warning "No compiler specific solution for __NO_INIT. __NO_INIT is ignored."
    #define __NO_INIT
  #endif
#endif

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

// Fault component version information
const char ARM_FaultVersion[] = ARM_FAULT_VERSION;

// Fault Information
ARM_FaultInfo_t ARM_FaultInfo __NO_INIT;

// Local function prototype
static uint32_t CalcCRC32 (      uint32_t init_val,
                           const uint8_t *data_ptr,
                                 uint32_t data_len,
                                 uint32_t polynom);

// ARM Fault Storage functions -------------------------------------------------

/**
  Clear the saved fault information.
*/
void ARM_FaultClear (void) {
  memset(&ARM_FaultInfo, 0, sizeof(ARM_FaultInfo));
}

/**
  Check if the fault occurred and if the fault information was saved properly.
  \return       status (1=Fault occurred and valid fault information exists,
                        0=No fault information saved yet or is invalid)
*/
uint32_t ARM_FaultOccurred (void) {
  uint32_t fault_info_valid = 1;

  // Check if magic number is valid
  if (ARM_FaultInfo.magic_number != ARM_FAULT_MAGIC_NUMBER) {
    fault_info_valid = 0U;
  }

  // Check if CRC of the ARM_FaultInfo structure is valid
  if (fault_info_valid != 0U) {
    if (ARM_FaultInfo.crc32 != CalcCRC32(ARM_FAULT_CRC32_INIT_VAL,
                                        (const uint8_t *)&ARM_FaultInfo.type,
                                        (sizeof(ARM_FaultInfo) - (sizeof(ARM_FaultInfo.magic_number) + sizeof(ARM_FaultInfo.crc32))),
                                         ARM_FAULT_CRC32_POLYNOM)) {
      fault_info_valid = 0U;
    }
  }

  return fault_info_valid;  
}

/**
  Save the fault information.
  Must be called from fault handler with preserved Link Register value, 
  typically by branching to this function.
*/
__NAKED void ARM_FaultSave (void) {
  //lint -efunc(10,  ARM_FaultSave) "Suppress: expecting ';'"
  //lint -efunc(522, ARM_FaultSave) "Suppress: Warning 522: Highest operation, a 'constant', lacks side-effects [MISRA 2012 Rule 2.2, required]"
  __ASM volatile (
#ifndef __ICCARM__
    ".syntax unified\n\t"
#endif

    "mov   r12, r4\n"                   // Store R4 to R12 (to use R4 in this function)
    "movs  r4,  #0\n"                   // Clear R4 for further usage

 /* --- Read current count value --- */
    "ldr   r1,  =%c[ARM_FaultInfo_count_addr]\n"
    "ldr   r3,  [r1]\n"

 /* --- Clear ARM_FaultInfo --- */
    "movs  r0,  #0\n"                           // R0 = 0
    "ldr   r1,  =%c[ARM_FaultInfo_addr]\n"      // R1 = &ARM_FaultInfo
    "movs  r2,  %[ARM_FaultInfo_size]\n"        // R2 = sizeof(ARM_FaultInfo)/4
    "b     is_clear_done\n"
  "clear_uint32:\n"
    "stm   r1!, {r0}\n"
    "subs  r2,  r2, #1\n"
  "is_clear_done:\n"
    "bne   clear_uint32\n"

 /* --- Increment and store new count value --- */
    "ldr   r1,  =%c[ARM_FaultInfo_count_addr]\n"
    "adds  r3,  r3, #1\n"
    "str   r3,  [r1]\n"

 /* Determine the beginning of the state context or the additional state context
    (for device with TruztZone) that was stacked upon exception entry and put that
    address into R3.
    For device with TrustZone, also determine if state context was pushed from
    Non-secure World but the exception handling is happening in the Secure World
    and if so, mark it by setting bit [0] of the R4 to value 1, thus indicating usage
    of Non-secure aliases.

    after this section:
      R3          == start of state context or additional state context if that was pushed also
      R4 bit [0]: == 0 - no access to Non-secure aliases or device without TrustZone
                  == 1 -    access to Non-secure aliases

    Determine by analyzing EXC_RETURN (Link Register):
    EXC_RETURN:
      - bit [6] (S):            only on device with TrustZone
                         == 0 - Non-secure stack was used
                         == 1 - Secure     stack was used
      - bit [5] (DCRS):         only on device with TrustZone
                         == 0 - additional state context was also stacked
                         == 1 - only       state context was stacked
      - bit [2] (SPSEL): == 0 - Main    Stack Pointer (MSP) was used for stacking on exception entry
                         == 1 - Process Stack Pointer (PSP) was used for stacking on exception entry */
    "mov   r0,  lr\n"                   // R0 = LR (EXC_RETURN)
    "lsrs  r0,  r0, #3\n"               // Shift bit [2] (SPSEL) into Carry flag
    "bcc   msp_used\n"                  // If    bit [2] (SPSEL) == 0, MSP or MSP_NS was used
                                        // If    bit [2] (SPSEL) == 1, PSP or PSP_NS was used
  "psp_used:\n"
#if (ARM_FAULT_TZ_SECURE != 0)          // If code was compiled for and is running in Secure World
    "mov   r0,  lr\n"                   // R0 = LR (EXC_RETURN)
    "lsrs  r0,  r0, #7\n"               // Shift   bit [6] (S) into Carry flag
    "bcs   load_psp\n"                  // If      bit [6] (S) == 1, jump to load PSP
  "load_psp_ns:\n"                      // else if bit [6] (S) == 0, load PSP_NS
    "mrs   r3,  psp_ns\n"               // R3 = PSP_NS
    "movs  r4,  #1\n"                   // R4 = 1
    "b     r3_points_to_stack\n"        // PSP_NS loaded to R3, exit section
  "load_psp:\n"
#endif
    "mrs   r3,  psp\n"                  // R3 = PSP
    "b     r3_points_to_stack\n"        // PSP loaded to R3, exit section

  "msp_used:\n"
#if (ARM_FAULT_TZ_SECURE != 0)          // If code was compiled for and is running in Secure World
    "mov   r0,  lr\n"                   // R0 = LR (EXC_RETURN)
    "lsrs  r0,  r0, #7\n"               // Shift   bit [6] (S) into Carry flag
    "bcs   load_msp\n"                  // If      bit [6] (S) == 1, jump to load MSP
  "load_msp_ns:\n"                      // else if bit [6] (S) == 0, load MSP_NS
    "mrs   r3,  msp_ns\n"               // R3 = MSP_NS
    "movs  r4,  #1\n"                   // R4 = 1
    "b     r3_points_to_stack\n"        // MSP_NS loaded to R3, exit section
  "load_msp:\n"
#endif
    "mrs   r3,  msp\n"                  // R3 = MSP
    "b     r3_points_to_stack\n"        // MSP loaded to R3, exit section

  "r3_points_to_stack:\n"

 /* Determine if stack contains valid state context (if fault was not a stacking fault).
    If stack information is not valid mark it by setting bit [1] of the R4 to value 1.
    Note: for Armv6-M and Armv8-M Baseline CFSR register is not available, so stack is 
          considered valid although it might not always be so. */
#if (ARM_FAULT_FAULT_REGS_EXIST != 0)   // If fault registers exist
    "ldr   r1,  =%c[cfsr_err_msk]\n"    // R1 = (SCB_CFSR_Stack_Err_Msk)
#if (ARM_FAULT_TZ_SECURE != 0)          // If code was compiled for and is running in Secure World
    "lsrs  r0,  r4, #1\n"               // Shift   bit [0] of R4 into Carry flag
    "bcc   load_cfsr_addr\n"            // If      bit [0] of R4 == 0, jump to load CFSR register address
  "load_cfsr_ns_addr:\n"                // else if bit [0] of R4 == 1, load CFSR_NS register address
    "ldr   r2,  =%c[cfsr_ns_addr]\n"    // R2 = CFSR_NS address
    "b     load_cfsr\n"
  "load_cfsr_addr:\n"
#endif
    "ldr   r2,  =%c[cfsr_addr]\n"       // R2 = CFSR address
  "load_cfsr:\n"
    "ldr   r0,  [r2]\n"                 // R0 = CFSR (or CFSR_NS) register value
    "ands  r0,  r1\n"                   // Mask CFSR value with stacking error bits
    "beq   stack_check_end\n"           // If   no stacking error, jump to stack_check_end
  "stack_check_failed:\n"               // else if stacking error, stack information is invalid
    "adds  r4,  #2\n"                   // R4 |= (1 << 1)
  "stack_check_end:\n"
#endif

 /* --- Type information --- */
    "ldr   r2,  =%c[ARM_FaultInfo_type_addr]\n"
    "ldr   r0,  =%c[ARM_FaultInfo_type_val]\n"
    "str   r0,  [r2]\n"

 /* --- State Context --- */
 /* Check if state context (also additional state context if it exists) is valid and
    if it is then copy it, otherwise skip copying */
    "lsrs  r0,  r4, #2\n"               // Shift bit [1] of R4 into Carry flag
    "bcs   state_context_end\n"         // If stack is not valid (bit == 1), skip copying information from stack

#if (ARM_FAULT_ARCH_ARMV8x_M != 0)      // If arch is Armv8/8.1-M
 /* If additional state context was stacked upon exception entry, copy it into ARM_FaultInfo */
    "mov   r0,  lr\n"                   // R0 = LR (EXC_RETURN)
    "lsrs  r0,  r0, #6\n"               // Shift   bit [5] (DCRS) into Carry flag
    "bcs   additional_context_end\n"    // If      bit [5] (DCRS) == 1, skip additional state context
                                        // else if bit [5] (DCRS) == 0, copy additional state context
    "ldr   r2,  =%c[ARM_FaultInfo_additonal_ctx_addr]\n"
    "ldm   r3!, {r0, r1}\n"             // Stacked IntegritySignature, Reserved
    "stm   r2!, {r0, r1}\n"
    "ldm   r3!, {r0, r1}\n"             // Stacked R4, R5
    "stm   r2!, {r0, r1}\n"
    "ldm   r3!, {r0, r1}\n"             // Stacked R6, R7
    "stm   r2!, {r0, r1}\n"
    "ldm   r3!, {r0, r1}\n"             // Stacked R8, R9
    "stm   r2!, {r0, r1}\n"
    "ldm   r3!, {r0, r1}\n"             // Stacked R10, R11
    "stm   r2!, {r0, r1}\n"

  "additional_context_end:\n"
#endif

 /* Copy state context stacked on exception entry into ARM_FaultInfo */
    "ldr   r2,  =%c[ARM_FaultInfo_state_ctx_addr]\n"
    "ldm   r3!, {r0, r1}\n"             // Stacked R0, R1
    "stm   r2!, {r0, r1}\n"
    "ldm   r3!, {r0, r1}\n"             // Stacked R2, R3
    "stm   r2!, {r0, r1}\n"
    "ldm   r3!, {r0, r1}\n"             // Stacked R12, LR
    "stm   r2!, {r0, r1}\n"
    "ldm   r3!, {r0, r1}\n"             // Stacked ReturnAddress, xPSR
    "stm   r2!, {r0, r1}\n"

  "state_context_end:\n"

 /* --- Common Registers --- */
 /* Store values of Common Registers into ARM_FaultInfo */
    "ldr   r2,  =%c[ARM_FaultInfo_common_regs_addr]\n"
    "mrs   r0,  xpsr\n"                 // R0 = current xPSR
    "mov   r1,  lr\n"                   // R1 = current LR (exception return code)
    "stm   r2!, {r0, r1}\n"
#if (ARM_FAULT_TZ_SECURE != 0)          // If code was compiled for and is running in Secure World
    "lsrs  r0,  r4, #1\n"               // Shift   bit [0] of R4 into Carry flag
    "bcc   load_sps\n"                  // If      bit [0] of R4 == 0, jump to load MSP and PSP
  "load_sps_ns:\n"                      // else if bit [0] of R4 == 1, load MSP_NS and PSP_NS
    "mrs   r0,  msp_ns\n"               // R0 = current MSP_NS
    "mrs   r1,  psp_ns\n"               // R1 = current PSP_NS
    "b     store_sps\n"
#endif
  "load_sps:\n"
    "mrs   r0,  msp\n"                  // R0 = current MSP
    "mrs   r1,  psp\n"                  // R1 = current PSP
  "store_sps:\n"
    "stm   r2!, {r0, r1}\n"             // Store MSP, PSP

 /* --- Armv8/8.1-M specific Registers --- */
 /* Store values of Armv8/8.1-M specific Registers (if they exist) into ARM_FaultInfo */
#if (ARM_FAULT_ARCH_ARMV8x_M != 0)      // If arch is Armv8/8.1-M
    "ldr   r2,  =%c[ARM_FaultInfo_armv8_m_regs_addr]\n"
#if (ARM_FAULT_TZ_SECURE != 0)          // If code was compiled for and is running in Secure World
    "lsrs  r0,  r4, #1\n"               // Shift   bit [0] of R4 into Carry flag
    "bcc   load_splims\n"               // If      bit [0] of R4 == 0, jump to load MSPLIM and PSPLIM
#if (ARM_FAULT_ARCH_ARMV8_M_BASE !=0)   // If arch is Armv8-M Baseline
    "b     splims_end\n"                // MSPLIM_NS and PSPLIM_NS do not exist, skip loading and storing them
#else                                   // Else if arch is Armv8/8.1-M Mainline
  "load_splims_ns:\n"                   // else if bit [0] of R4 == 1, load MSPLIM_NS and PSPLIM_NS
    "mrs   r0,  msplim_ns\n"            // R0 = current MSPLIM_NS
    "mrs   r1,  psplim_ns\n"            // R1 = current PSPLIM_NS
    "b     store_splims\n"
#endif
#endif
  "load_splims:\n"
    "mrs   r0,  msplim\n"               // R0 = current MSP
    "mrs   r1,  psplim\n"               // R1 = current PSP
  "store_splims:\n"
    "stm   r2!, {r0, r1}\n"
  "splims_end:\n"
#endif

 /* --- Fault Registers --- */
 /* Store values of Fault Registers (if they exist) into ARM_FaultInfo */
#if (ARM_FAULT_FAULT_REGS_EXIST != 0)   // If fault registers exist
    "ldr   r2,  =%c[ARM_FaultInfo_fault_regs_addr]\n"
#if (ARM_FAULT_TZ_SECURE != 0)          // If code was compiled for and is running in Secure World
    "lsrs  r0,  r4, #1\n"               // Shift   bit [0] of R4 into Carry flag
    "bcc   load_scb_addr\n"             // If      bit [0] of R4 == 0, jump to load SCB address
  "load_scb_ns_addr:\n"                 // else if bit [0] of R4 == 1, load SCB_NS address
    "ldr   r3,  =%c[scb_ns_addr]\n"
    "b     load_fault_regs\n"
  "load_scb_addr:\n"
#endif
    "ldr   r3,  =%c[scb_addr]\n"
  "load_fault_regs:\n"
    "ldr   r0,  [r3, %[cfsr_ofs]]\n"    // R0 = CFSR
    "ldr   r1,  [r3, %[hfsr_ofs]]\n"    // R1 = HFSR
    "stm   r2!, {r0, r1}\n"
    "ldr   r0,  [r3, %[dfsr_ofs]]\n"    // R0 = DFSR
    "ldr   r1,  [r3, %[mmfar_ofs]]\n"   // R1 = MMFAR
    "stm   r2!, {r0, r1}\n"
    "ldr   r0,  [r3, %[bfar_ofs]]\n"    // R0 = BFSR
    "ldr   r1,  [r3, %[afsr_ofs]]\n"    // R1 = AFSR
    "stm   r2!, {r0, r1}\n"

 /* --- Armv8/8.1-M Fault Registers --- */
 /* Store values of Armv8/8.1-M Fault Registers (if they exist) and if code is running in Secure World into ARM_FaultInfo */
#if (ARM_FAULT_TZ_SECURE != 0)          // If code was compiled for and is running in Secure World
    "ldr   r2,  =%c[ARM_FaultInfo_armv8_m_fault_regs_addr]\n"
    "ldr   r3,  =%c[scb_addr]\n"
    "ldr   r0,  [r3, %[sfsr_ofs]]\n"    // R0 = SFSR
    "ldr   r1,  [r3, %[sfar_ofs]]\n"    // R1 = SFAR
    "stm   r2!, {r0, r1}\n"
#endif
#endif

 /* Calculate CRC-32 on ARM_FaultInfo structure (excluding magic_number and crc32 fields) and
    store it into ARM_FaultInfo.crc32 */
    "ldr   r0,  =%c[crc_init_val]\n"    // R0 = init_val parameter
    "ldr   r1,  =%c[crc_data_ptr]\n"    // R1 = data_ptr parameter
    "ldr   r2,  =%c[crc_data_len]\n"    // R2 = data_len parameter
    "ldr   r3,  =%c[crc_polynom]\n"     // R3 = polynom  parameter
    "bl    CalcCRC32\n"                 // Call CalcCRC32 function
    "ldr   r2,  =%c[ARM_FaultInfo_crc32_addr]\n"
    "str   r0,  [r2]\n"                 // Store CRC-32

 /* Store magic number into ARM_FaultInfo.magic_number */
    "ldr   r2,  =%c[ARM_FaultInfo_magic_number_addr]\n"
    "ldr   r0,  =%c[ARM_FaultInfo_magic_number_val]\n"
    "str   r0,  [r2]\n"

    "mov   r4,  r12\n"                  // Restore R4 from R12

    "bl    ARM_FaultExit\n"             // Call ARM_FaultExit function

 /* Inline assembly template operands */
 :  /* no outputs */
 :  /* inputs */
    [ARM_FaultInfo_addr]                    "i" (&ARM_FaultInfo)
  , [ARM_FaultInfo_size]                    "i" (sizeof(ARM_FaultInfo)/4)
  , [ARM_FaultInfo_magic_number_addr]       "i" (&ARM_FaultInfo.magic_number)
  , [ARM_FaultInfo_magic_number_val]        "i" (ARM_FAULT_MAGIC_NUMBER)
  , [ARM_FaultInfo_crc32_addr]              "i" (&ARM_FaultInfo.crc32)
  , [ARM_FaultInfo_count_addr]              "i" (&ARM_FaultInfo.count)
  , [ARM_FaultInfo_type_addr]               "i" (&ARM_FaultInfo.type)
  , [ARM_FaultInfo_type_val]                "i" (ARM_FAULT_FAULT_INFO_VER_MINOR
                                             |  (ARM_FAULT_FAULT_INFO_VER_MAJOR <<  8)
                                             |  (ARM_FAULT_FAULT_REGS_EXIST     << 16)
                                             |  (ARM_FAULT_ARCH_ARMV8x_M        << 17)
                                             |  (ARM_FAULT_TZ_SECURE            << 18))
  , [ARM_FaultInfo_state_ctx_addr]          "i" (&ARM_FaultInfo.R0)
  , [ARM_FaultInfo_common_regs_addr]        "i" (&ARM_FaultInfo.xPSR_in_handler)
#if (ARM_FAULT_FAULT_REGS_EXIST != 0)
  , [ARM_FaultInfo_fault_regs_addr]         "i" (&ARM_FaultInfo.SCB_CFSR)
  , [cfsr_err_msk]                          "i" (SCB_CFSR_Stack_Err_Msk)
  , [scb_addr]                              "i" (SCB_BASE)
  , [cfsr_addr]                             "i" (SCS_BASE + offsetof(SCB_Type, CFSR))
  , [cfsr_ofs]                              "i" (offsetof(SCB_Type, CFSR ))
  , [hfsr_ofs]                              "i" (offsetof(SCB_Type, HFSR ))
  , [dfsr_ofs]                              "i" (offsetof(SCB_Type, DFSR ))
  , [mmfar_ofs]                             "i" (offsetof(SCB_Type, MMFAR))
  , [bfar_ofs]                              "i" (offsetof(SCB_Type, BFAR ))
  , [afsr_ofs]                              "i" (offsetof(SCB_Type, AFSR ))
#if (ARM_FAULT_TZ_SECURE != 0)
  , [scb_ns_addr]                           "i" (SCB_BASE_NS)
  , [cfsr_ns_addr]                          "i" (SCS_BASE_NS + offsetof(SCB_Type, CFSR))
#endif
#endif
#if (ARM_FAULT_ARCH_ARMV8x_M != 0)
  , [ARM_FaultInfo_additonal_ctx_addr]      "i" (&ARM_FaultInfo.IntegritySignature)
  , [ARM_FaultInfo_armv8_m_regs_addr]       "i" (&ARM_FaultInfo.MSPLIM)
#endif
#if (ARM_FAULT_ARCH_ARMV8x_M_MAIN !=0)
  , [ARM_FaultInfo_armv8_m_fault_regs_addr] "i" (&ARM_FaultInfo.SCB_SFSR)
  , [sfsr_ofs]                              "i" (offsetof(SCB_Type, SFSR ))
  , [sfar_ofs]                              "i" (offsetof(SCB_Type, SFAR ))
#endif
  , [crc_init_val]                          "i" (ARM_FAULT_CRC32_INIT_VAL)
  , [crc_data_ptr]                          "i" (&ARM_FaultInfo.type)
  , [crc_data_len]                          "i" (sizeof(ARM_FaultInfo) - (sizeof(ARM_FaultInfo.magic_number) + sizeof(ARM_FaultInfo.crc32)))
  , [crc_polynom]                           "i" (ARM_FAULT_CRC32_POLYNOM)
 :  /* clobber list */
    "r0", "r1", "r2", "r3", "r4", "r12", "lr" , "cc", "memory");
}

/**
  Callback function called after fault information was saved.
  Used to provide a specific reaction to fault after it was saved.
  The default implementation will reset the system via the CMSIS NVIC_SystemReset function.
  User can override this function to provide required reaction.
*/
__WEAK __NO_RETURN void ARM_FaultExit (void) {
  NVIC_SystemReset();                   // Reset the system
}


// Helper function

#ifdef __ICCARM__
#pragma diag_suppress=Pe940
#endif

/**
  Calculate CRC-32 on data block in memory
  \param[in]    init_val        initial CRC value
  \param[in]    data_ptr        pointer to data
  \param[in]    data_len        data length (in bytes)
  \param[in]    polynom         CRC polynom
  \return       CRC-32 value (32-bit)
*/
static __NAKED __USED uint32_t CalcCRC32 (      uint32_t init_val,
                                          const uint8_t *data_ptr,
                                                uint32_t data_len,
                                                uint32_t polynom) {
  //lint -esym(528,  CalcCRC32) "Suppress: Warning 528: Symbol 'CalcCRC32(uint32_t, const uint8_t *, uint32_t, uint32_t)' not referenced"
  //lint -efunc(10,  CalcCRC32) "Suppress: expecting ';'"
  //lint -efunc(522, CalcCRC32) "Suppress: Warning 522: Highest operation, a 'constant', lacks side-effects [MISRA 2012 Rule 2.2, required]"
  //lint -efunc(533, CalcCRC32) "Warning 533: function 'CalcCRC32(uint32_t, const uint8_t *, uint32_t, uint32_t)' should return a value [MISRA 2012 Rule 17.4, mandatory]"
  //lint -efunc(715, CalcCRC32) "Info 715: Symbol not referenced [MISRA 2012 Rule 2.7, advisory]"
  __ASM volatile (
#ifndef __ICCARM__
    ".syntax unified\n"
#endif
    "mov   r12, r4\n"
    "b     check\n"
  "loop:\n"
    "ldrb  r4,  [r1]\n"
    "lsls  r4,  r4, #24\n"
    "eors  r0,  r0, r4\n"
    "lsls  r0,  r0, #1\n"
    "bcc   skip_xor_7\n"
    "eors  r0,  r0, r3\n"
  "skip_xor_7:\n"
    "lsls  r0,  r0, #1\n"
    "bcc   skip_xor_6\n"
    "eors  r0,  r0, r3\n"
  "skip_xor_6:\n"
    "lsls  r0,  r0, #1\n"
    "bcc   skip_xor_5\n"
    "eors  r0,  r0, r3\n"
  "skip_xor_5:\n"
    "lsls  r0,  r0, #1\n"
    "bcc   skip_xor_4\n"
    "eors  r0,  r0, r3\n"
  "skip_xor_4:\n"
    "lsls  r0,  r0, #1\n"
    "bcc   skip_xor_3\n"
    "eors  r0,  r0, r3\n"
  "skip_xor_3:\n"
    "lsls  r0,  r0, #1\n"
    "bcc   skip_xor_2\n"
    "eors  r0,  r0, r3\n"
  "skip_xor_2:\n"
    "lsls  r0,  r0, #1\n"
    "bcc   skip_xor_1\n"
    "eors  r0,  r0, r3\n"
  "skip_xor_1:\n"
    "lsls  r0,  r0, #1\n"
    "bcc   skip_xor_0\n"
    "eors  r0,  r0, r3\n"
  "skip_xor_0:\n"
    "adds  r1,  r1, #1\n"
    "subs  r2,  r2, #1\n"
  "check:\n"
    "cmp   r2,  #0\n"
    "bne   loop\n"
    "mov   r4,  r12\n"
    "bx    lr\n"
 :::
    "r0", "r1", "r2", "r3", "r4", "r12", "cc");
}

#ifdef __ICCARM__
#pragma diag_warning=Pe940
#endif
