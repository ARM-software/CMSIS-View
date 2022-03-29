/*
 * Copyright (c) 2016-2022 Arm Limited. All rights reserved.
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

//lint -emacro((923,9078),CoreDebug,DWT,SysTick) "cast from unsigned long to pointer"
//lint -ecall(534,__disable_irq) "Ignoring return value"

#include "RTE_Components.h"
#include CMSIS_device_header

#include <string.h>
#include "EventRecorder.h"
#include "EventRecorderConf.h"

#if (EVENT_TIMESTAMP_SOURCE == 2)
#include "cmsis_os2.h"
#endif

#if !defined(__USED)
  #define __USED __attribute__((used))
#endif
#if !defined(__WEAK)
  #define __WEAK __attribute__((weak))
#endif
#if !defined(__ALIGNED)
  #define __ALIGNED(x) __attribute__((aligned(x)))
#endif
#if !defined(__NO_INIT)
 //lint -esym(9071, __NO_INIT) "defined macro is reserved to the compiler"
 #if   defined (__CC_ARM)                                           /* ARM Compiler 4/5 */
  #define __NO_INIT __attribute__ ((section (".bss.noinit"), zero_init))
 #elif defined (__ARMCC_VERSION) && (__ARMCC_VERSION >= 6010050)    /* ARM Compiler 6 */
  #define __NO_INIT __attribute__ ((section (".bss.noinit")))
 #elif defined (__GNUC__)                                           /* GNU Compiler */
  #define __NO_INIT __attribute__ ((section (".bss.noinit")))
 #else
  #warning "No compiler specific solution for __NO_INIT. __NO_INIT is ignored."
  #define __NO_INIT
 #endif
#endif

//lint -e(9026) "Function-like macro"
#define INT_LOG2(n) \
 (((n) == (1UL<<31)) ? 31 : ((n) == (1UL<<30)) ? 30 : ((n) == (1UL<<29)) ? 29 : ((n) == (1UL<<28)) ? 28 : \
  ((n) == (1UL<<27)) ? 27 : ((n) == (1UL<<26)) ? 26 : ((n) == (1UL<<25)) ? 25 : ((n) == (1UL<<24)) ? 24 : \
  ((n) == (1UL<<23)) ? 23 : ((n) == (1UL<<22)) ? 22 : ((n) == (1UL<<21)) ? 21 : ((n) == (1UL<<20)) ? 20 : \
  ((n) == (1UL<<19)) ? 19 : ((n) == (1UL<<18)) ? 18 : ((n) == (1UL<<17)) ? 17 : ((n) == (1UL<<16)) ? 16 : \
  ((n) == (1UL<<15)) ? 15 : ((n) == (1UL<<14)) ? 14 : ((n) == (1UL<<13)) ? 13 : ((n) == (1UL<<12)) ? 12 : \
  ((n) == (1UL<<11)) ? 11 : ((n) == (1UL<<10)) ? 10 : ((n) == (1UL<< 9)) ?  9 : ((n) == (1UL<< 8)) ?  8 : \
  ((n) == (1UL<< 7)) ?  7 : ((n) == (1UL<< 6)) ?  6 : ((n) == (1UL<< 5)) ?  5 : ((n) == (1UL<< 4)) ?  4 : \
  ((n) == (1UL<< 3)) ?  3 : ((n) == (1UL<< 2)) ?  2 : ((n) == (1UL<< 1)) ?  1 : ((n) == (1UL<< 0)) ?  0 : -1)

/* Event Recorder Component */
#define CID_EVENT               0xFFU

/* Event Recorder Messages */
#define MID_EVENT_INIT          0x00U   // Initialization
#define MID_EVENT_START         0x01U   // Start Recorder
#define MID_EVENT_STOP          0x02U   // Stop Recorder
#define MID_EVENT_CLOCK         0x03U   // Clock changed

//lint -emacro((835),ID_EVENT_*) "A zero has been given as argument to operator '|'"
#define ID_EVENT_INIT   (((uint32_t)CID_EVENT << 8) | MID_EVENT_INIT  | EVENT_RECORD_FIRST | EVENT_RECORD_LAST)
#define ID_EVENT_START  (((uint32_t)CID_EVENT << 8) | MID_EVENT_START | EVENT_RECORD_FIRST | EVENT_RECORD_LAST)
#define ID_EVENT_STOP   (((uint32_t)CID_EVENT << 8) | MID_EVENT_STOP  | EVENT_RECORD_FIRST | EVENT_RECORD_LAST)
#define ID_EVENT_CLOCK  (((uint32_t)CID_EVENT << 8) | MID_EVENT_CLOCK | EVENT_RECORD_FIRST | EVENT_RECORD_LAST)

/* Event Recorder Signature */
#define SIGNATURE               0xE1A5276BU

/* Number of Records in Event Buffer */
#define EVENT_RECORD_COUNT_BITS INT_LOG2(EVENT_RECORD_COUNT)
#if ((EVENT_RECORD_COUNT_BITS > 20) || (EVENT_RECORD_COUNT_BITS < 3))
#error "Invalid number of Event Buffer Records!"
#endif

/* Maximum number of Locked Records */
#define EVENT_RECORD_MAX_LOCKED 7U

/* Event Data Maximum Length */
#if (((EVENT_RECORD_COUNT / 4U) * 8U) > 256U)
#define EVENT_DATA_MAX_LENGTH   256U
#else
#define EVENT_DATA_MAX_LENGTH   ((EVENT_RECORD_COUNT / 4U) * 8U)
#endif

/* Event Record Information */
#define EVENT_RECORD_ID_MASK    0x0000FFFFU
#define EVENT_RECORD_DLEN_POS   16
#define EVENT_RECORD_DLEN_MASK  0x00070000U
#define EVENT_RECORD_CTX_POS    16
#define EVENT_RECORD_CTX_MASK   0x00070000U
#define EVENT_RECORD_IRQ        0x00080000U
#define EVENT_RECORD_SEQ_POS    20
#define EVENT_RECORD_SEQ_MASK   0x00F00000U
#define EVENT_RECORD_FIRST      0x01000000U
#define EVENT_RECORD_LAST       0x02000000U
#define EVENT_RECORD_LOCKED     0x04000000U
#define EVENT_RECORD_VALID      0x08000000U
#define EVENT_RECORD_MSB_TS     0x10000000U
#define EVENT_RECORD_MSB_VAL1   0x20000000U
#define EVENT_RECORD_MSB_VAL2   0x40000000U
#define EVENT_RECORD_TBIT       0x80000000U

/* Event Record */
typedef struct {
  uint32_t ts;                  // Timestamp (32-bit, Toggle bit instead of MSB)
  uint32_t val1;                // Value 1   (32-bit, Toggle bit instead of MSB)
  uint32_t val2;                // Value 2   (32-bit, Toggle bit instead of MSB)
  uint32_t info;                // Record Information
                                //  [ 7.. 0]: Message ID (8-bit)
                                //  [15.. 8]: Component ID (8-bit)
                                //  [18..16]: Data Length (1..8) / Event Context
                                //      [19]: IRQ Flag
                                //  [23..20]: Sequence Number
                                //      [24]: First Record
                                //      [25]: Last Record
                                //      [26]: Locked Record
                                //      [27]: Valid Record
                                //      [28]: Timestamp MSB
                                //      [29]: Value 1 MSB
                                //      [30]: Value 2 MSB
                                //      [31]: Toggle bit
} EventRecord_t;

#ifdef RTE_Compiler_EventRecorder_Semihosting

/* Event Record Types (Log) */
#define EVENT_TYPE_DATA         0x0001U  // EventRecordData
#define EVENT_TYPE_VAL2         0x0002U  // EventRecord2
#define EVENT_TYPE_VAL4         0x0003U  // EventRecord4

/* Event Record Header (Log) */
typedef struct __PACKED {
  uint16_t      type;           // Type (EVENT_TYPE_xxx)
  uint16_t    length;           // Total length (in bytes)
} EventRecordHead_t;

/* Event Record for EventRecord2 (Log) */
typedef struct __PACKED {
  uint64_t        ts;           // Timestamp
  struct __PACKED {             // Record Information
    uint32_t      id : 16;      //  [ 7.. 0]: Message ID (8-bit)
                                //  [15.. 8]: Component ID (8-bit)
    uint32_t   rsrvd : 15;      // Reserved
    uint32_t     irq :  1;      // IRQ Flag
  } info;
  uint32_t val1;                // Value 1
  uint32_t val2;                // Value 2
} EventRecord2_t;

/* Event Record for EventRecord4 (Log) */
typedef struct __PACKED {
  uint64_t        ts;           // Timestamp
  struct __PACKED {             // Record Information
    uint32_t      id : 16;      //  [ 7.. 0]: Message ID (8-bit)
                                //  [15.. 8]: Component ID (8-bit)
    uint32_t   rsrvd : 15;      // Reserved
    uint32_t     irq :  1;      // IRQ Flag
  } info;
  uint32_t val1;                // Value 1
  uint32_t val2;                // Value 2
  uint32_t val3;                // Value 3
  uint32_t val4;                // Value 4
} EventRecord4_t;

/* Event Record for EventRecordData (Log) */
typedef struct __PACKED {
  uint64_t        ts;           // Timestamp
  struct __PACKED{              // Record Information
    uint32_t      id : 16;      //  [ 7.. 0]: Message ID (8-bit)
                                //  [15.. 8]: Component ID (8-bit)
    uint32_t  length : 15;      //  Data Length
    uint32_t     irq :  1;      // IRQ Flag
  } info;
//uint8_t data[info.length];    // Data
} EventRecordData_t;

#endif

/* Event Buffer */
static EventRecord_t EventBuffer[EVENT_RECORD_COUNT] __NO_INIT __ALIGNED(16);

/* Event Filter: 1024 enable bits for 8-bit Component ID with 2-bit Level */
static uint8_t EventFilter[128] __NO_INIT __ALIGNED(4);

/* Event Recorder Status */
typedef struct {
  uint8_t  state;               // Recorder State: 0 - Inactive, 1 - Running
  uint8_t  context;             // Current Event Context
  uint16_t info_crc;            // EventRecorderInfo CRC16-CCITT
  uint32_t record_index;        // Current Record Index
  uint32_t records_written;     // Number of records written
  uint32_t records_dumped;      // Number of records dumped
  uint32_t ts_overflow;         // Timestamp overflow counter
  uint32_t ts_freq;             // Timestamp frequency
  uint32_t ts_last;             // Timestamp last value
  uint32_t init_count;          // Initialization counter
  uint32_t signature;           // Initialization signature
} EventStatus_t;

static EventStatus_t EventStatus __NO_INIT __ALIGNED(64);

/* Global Event Recorder Information */
typedef struct {
  uint8_t    protocol_type;     // Protocol Type: 1 - DAP
  uint8_t    reserved;          // Reserved (must be zero)
  uint16_t   protocol_version;  // Protocol Version: [15..8]=major, [7..0]=minor
  // Version 1.x specific information
  uint32_t       record_count;  // Number of Records in Event Buffer
  EventRecord_t *event_buffer;  // Pointer to Event Buffer
  uint8_t       *event_filter;  // Pointer to Event Filter
  EventStatus_t *event_status;  // Pointer to Event Status
  uint8_t        ts_source;     // Timestamp source
  uint8_t        reserved3[3];  // Reserved (must be zero)
} EventRecorderInfo_t;

//lint -esym(754, EventRecorderInfo*) "Referenced   (used by debugger)"
//lint -esym(765, EventRecorderInfo)  "Global scope (used by debugger)"
//lint -esym(9003,EventRecorderInfo)  "Global scope (used by debugger)"
//lint ++fem
extern const EventRecorderInfo_t EventRecorderInfo;
__USED const EventRecorderInfo_t EventRecorderInfo =
{
  1U, 0U,
  0x0101U,                      // Protocol Version 1.1
  EVENT_RECORD_COUNT,
  &EventBuffer[0],
  &EventFilter[0],
  &EventStatus,
  EVENT_TIMESTAMP_SOURCE,
  { 0U, 0U, 0U }
};
//lint --fem


/* Atomic operation helper functions */

#if (__CORTEX_M < 3U)

__STATIC_INLINE uint8_t atomic_inc8 (uint8_t *mem) {
  uint32_t primask = __get_PRIMASK();
  uint8_t  ret;

  __disable_irq();
  ret = *mem;
  *mem = ret + 1U;
  if (primask == 0U) {
    __enable_irq();
  }

  return ret;
}

__STATIC_INLINE uint32_t atomic_inc32 (uint32_t *mem) {
  uint32_t primask = __get_PRIMASK();
  uint32_t ret;

  __disable_irq();
  ret = *mem;
  *mem = ret + 1U;
  if (primask == 0U) {
    __enable_irq();
  }

  return ret;
}

__STATIC_INLINE uint8_t atomic_wr8 (uint8_t *mem, uint8_t val) {
  uint32_t primask = __get_PRIMASK();
  uint8_t  ret;

  __disable_irq();
  ret = *mem;
  *mem = val;
  if (primask == 0U) {
    __enable_irq();
  }

  return ret;
}

__STATIC_INLINE uint32_t atomic_wr32 (uint32_t *mem, uint32_t val) {
  uint32_t primask = __get_PRIMASK();
  uint32_t ret;

  __disable_irq();
  ret = *mem;
  *mem = val;
  if (primask == 0U) {
    __enable_irq();
  }

  return ret;
}

#else /* (__CORTEX_M >= 3U) */

//lint ++flb

#if defined(__CC_ARM)
static __asm    uint8_t atomic_inc8 (uint8_t *mem) {
  mov    r2,r0
1
  ldrexb r0,[r2]
  adds   r1,r0,#1
  strexb r3,r1,[r2]
  cbz    r3,%F2
  b      %B1
2
  bx     lr
}
#else
__STATIC_INLINE uint8_t atomic_inc8 (uint8_t *mem) {
  register uint32_t val, res;
  register uint8_t  ret;

  __ASM volatile (
  ".syntax unified\n"
  "1:\n\t"
    "ldrexb %[ret],[%[mem]]\n\t"
    "adds   %[val],%[ret],#1\n\t"
    "strexb %[res],%[val],[%[mem]]\n\t"
    "cbz    %[res],2f\n\t"
    "b      1b\n"
  "2:"
  : [ret] "=&l" (ret),
    [val] "=&l" (val),
    [res] "=&l" (res)
  : [mem] "l"   (mem)
  : "cc", "memory"
  );

  return ret;
}
#endif

#if defined(__CC_ARM)
static __asm    uint32_t atomic_inc32 (uint32_t *mem) {
  mov   r2,r0
1
  ldrex r0,[r2]
  adds  r1,r0,#1
  strex r3,r1,[r2]
  cbz   r3,%F2
  b     %B1
2
  bx    lr
}
#else
__STATIC_INLINE uint32_t atomic_inc32 (uint32_t *mem) {
  register uint32_t val, res;
  register uint32_t ret;

  __ASM volatile (
  ".syntax unified\n"
  "1:\n\t"
    "ldrex %[ret],[%[mem]]\n\t"
    "adds  %[val],%[ret],#1\n\t"
    "strex %[res],%[val],[%[mem]]\n\t"
    "cbz   %[res],2f\n\t"
    "b     1b\n"
  "2:"
  : [ret] "=&l" (ret),
    [val] "=&l" (val),
    [res] "=&l" (res)
  : [mem] "l"   (mem)
  : "cc", "memory"
  );

  return ret;
}
#endif

#if defined(__CC_ARM)
static __asm    uint8_t atomic_wr8 (uint8_t *mem, uint8_t val) {
  mov    r2,r0
1
  ldrexb r0,[r2]
  strexb r3,r1,[r2]
  cbz    r3,%F2
  b      %B1
2
  bx     lr
}
#else
__STATIC_INLINE uint8_t atomic_wr8 (uint8_t *mem, uint8_t val) {
  register uint32_t res;
  register uint8_t  ret;

  __ASM volatile (
  ".syntax unified\n"
  "1:\n\t"
    "ldrexb %[ret],[%[mem]]\n\t"
    "strexb %[res],%[val],[%[mem]]\n\t"
    "cbz    %[res],2f\n\t"
    "b      1b\n"
  "2:"
  : [ret] "=&l" (ret),
    [res] "=&l" (res)
  : [mem] "l"   (mem),
    [val] "l"   (val)
  : "memory"
  );

  return ret;
}
#endif

#if defined(__CC_ARM)
static __asm    uint32_t atomic_wr32 (uint32_t *mem, uint32_t val) {
  mov   r2,r0
1
  ldrex r0,[r2]
  strex r3,r1,[r2]
  cbz   r3,%F2
  b     %B1
2
  bx    lr
}
#else
__STATIC_INLINE uint32_t atomic_wr32 (uint32_t *mem, uint32_t val) {
  register uint32_t res;
  register uint32_t ret;

  __ASM volatile (
  ".syntax unified\n"
  "1:\n\t"
    "ldrex %[ret],[%[mem]]\n\t"
    "strex %[res],%[val],[%[mem]]\n\t"
    "cbz   %[res],2f\n\t"
    "b     1b\n"
  "2:"
  : [ret] "=&l" (ret),
    [res] "=&l" (res)
  : [mem] "l"   (mem),
    [val] "l"   (val)
  : "memory"
  );

  return ret;
}
#endif

//lint --flb

#endif


__STATIC_INLINE uint32_t GetContext (void) {
  return ((uint32_t)atomic_inc8(&EventStatus.context));
}

__STATIC_INLINE uint32_t GetRecordIndex (void) {
  return (atomic_inc32(&EventStatus.record_index));
}

__STATIC_INLINE uint32_t UpdateTS (uint32_t ts) {
  return (atomic_wr32(&EventStatus.ts_last, ts));
}

static uint8_t TS_OverflowLock;

__STATIC_INLINE uint8_t LockTS_Overflow (void) {
  return (atomic_wr8(&TS_OverflowLock, 1U));
}

__STATIC_INLINE void UnlockTS_Overflow (void) {
  (void) (atomic_wr8(&TS_OverflowLock, 0U));
}

__STATIC_INLINE void IncrementRecordsWritten (void) {
  (void)atomic_inc32(&EventStatus.records_written);
}

__STATIC_INLINE void IncrementRecordsDumped (void) {
  (void)atomic_inc32(&EventStatus.records_dumped);
}


#if (__CORTEX_M < 3U)

__STATIC_INLINE uint32_t LockRecord (uint32_t *mem, uint32_t info) {
  uint32_t primask = __get_PRIMASK();
  uint32_t val;

  __disable_irq();
  val = *mem;
  if ((val & EVENT_RECORD_LOCKED) == 0U) {
     val = (val & EVENT_RECORD_TBIT) | info;
    *mem = val;
  } else {
     val = 0U;
  }
  if (primask == 0U) {
    __enable_irq();
  }

  return val;
}

__STATIC_INLINE uint32_t UnlockRecord (uint32_t *mem, uint32_t info) {
  uint32_t primask = __get_PRIMASK();
  uint32_t val;
  uint32_t ret;

  __disable_irq();
  val = *mem;
  if ((val & EVENT_RECORD_LOCKED) != 0U) {
    *mem = info;
     ret = 1U;
  } else {
     ret = 0U;
  }
  if (primask == 0U) {
    __enable_irq();
  }

  return ret;
}

#else /* (__CORTEX_M >= 3U) */

//lint ++flb

#if defined(__CC_ARM)
static __asm    uint32_t LockRecord (uint32_t *mem, uint32_t info) {
  mov   r2,r0
1
  ldrex r0,[r2]
  lsls  r3,r0,#__cpp(32 - INT_LOG2(EVENT_RECORD_LOCKED))
  bcc   %F2
  clrex
  movs  r0,#0
  bx    lr
2
  and   r0,r0,#__cpp(EVENT_RECORD_TBIT)
  orrs  r0,r1
  strex r3,r0,[r2]
  cbz   r3,%F3
  b     %B1
3
  bx    lr
}
#else
__STATIC_INLINE uint32_t LockRecord (uint32_t *mem, uint32_t info) {
  register uint32_t val, res, tmp;

  __ASM volatile (
  ".syntax unified\n"
  "1:\n\t"
    "ldrex %[val],[%[mem]]\n\t"
    "lsls  %[tmp],%[val],%[Ln]\n\t"
    "bcc   2f\n\t"
    "clrex\n\t"
    "movs  %[val],#0\n\t"
    "b     3f\n"
  "2:\n\t"
#if (defined(__ARM_ARCH_8M_BASE__) && (__ARM_ARCH_8M_BASE__ != 0))
    "lsrs  %[val],%[val],%[Tp]\n\t"
    "lsls  %[val],%[val],%[Tp]\n\t"
#else
    "and   %[val],%[val],%[Tbit]\n\t"
#endif
    "orrs  %[val],%[info]\n\t"
    "strex %[res],%[val],[%[mem]]\n\t"
    "cbz   %[res],3f\n\t"
    "b     1b\n"
  "3:"
  : [val]  "=&l" (val),
    [res]  "=&l" (res),
    [tmp]  "=&l" (tmp)
  : [mem]  "l"   (mem),
    [info] "l"   (info),
    [Ln]   "I"   (32 - INT_LOG2(EVENT_RECORD_LOCKED)),
#if (defined(__ARM_ARCH_8M_BASE__) && (__ARM_ARCH_8M_BASE__ != 0))
    [Tp]   "I"   (INT_LOG2(EVENT_RECORD_TBIT))
#else
    [Tbit] "I"   (EVENT_RECORD_TBIT)
#endif
  : "cc", "memory"
  );

  return val;
}
#endif

#if defined(__CC_ARM)
static __asm    uint32_t UnlockRecord (uint32_t *mem, uint32_t info) {
  mov   r2,r0
1
  ldrex r0,[r2]
  lsls  r0,r0,#__cpp(32 - INT_LOG2(EVENT_RECORD_LOCKED))
  bcs   %F2
  clrex
  movs  r0,#0
  bx    lr
2
  strex r3,r1,[r2]
  cbz   r3,%F3
  b     %B1
3
  movs  r0,#1
4
  bx    lr
}
#else
__STATIC_INLINE uint32_t UnlockRecord (uint32_t *mem, uint32_t info) {
  register uint32_t val, res, ret;

  __ASM volatile (
  ".syntax unified\n"
  "1:\n\t"
    "ldrex %[val],[%[mem]]\n\t"
    "lsls  %[val],%[val],%[Ln]\n\t"
    "bcs   2f\n\t"
    "clrex\n\t"
    "movs  %[ret],#0\n\t"
    "b     4f\n"
  "2:\n\t"
    "strex %[res],%[info],[%[mem]]\n\t"
    "cbz   %[res],3f\n\t"
    "b     1b\n"
  "3:\n\t"
    "movs  %[ret],#1\n"
  "4:"
  : [ret]  "=&l" (ret),
    [val]  "=&l" (val),
    [res]  "=&l" (res)
  : [mem]  "l"   (mem),
    [info] "l"   (info),
    [Ln]   "I"   (32 - INT_LOG2(EVENT_RECORD_LOCKED))
  : "cc", "memory"
  );

  return ret;
}
#endif

//lint --flb

#endif


/* Semihosting */

#ifdef RTE_Compiler_EventRecorder_Semihosting

#ifndef EVENT_LOG_FILENAME
#define EVENT_LOG_FILENAME      "EventRecorder.log"
#endif

#define SYS_OPEN                0x01U
#define SYS_CLOSE               0x02U
#define SYS_WRITE               0x05U

#define MODE_wb                 5U

#if defined(__CC_ARM)
static __asm    int32_t semihosting_call (uint32_t operation, void *args) {
  bkpt  0xab
  bx    lr
}
#else
__STATIC_INLINE int32_t semihosting_call (uint32_t operation, void *args) {
  //lint --e{438} "Last value assigned to variable not used"
  //lint --e{529} "Symbol not subsequently referenced"
  //lint --e{923} "cast from pointer to unsigned int"
  register uint32_t __r0 __ASM("r0") = operation;
  register uint32_t __r1 __ASM("r1") = (uint32_t)args;

  __ASM volatile (
    "bkpt 0xab" : "=r"(__r0) : "r"(__r0), "r"(__r1) :
  );
 
  return (int32_t)__r0;
}
#endif

typedef int32_t FILEHANDLE;

static FILEHANDLE FileHandle = -1;

static FILEHANDLE sys_open (const char *name, uint32_t openmode) {
  //lint --e{446}  "side effect in initializer"
  //lint --e{934}  "Taking address of near auto variable"
  struct {
    const char    *name;
    uint32_t       openmode;
    size_t         len;
  } args = { name, openmode, strlen(name) };
  (void)args.name;
  (void)args.openmode;
  (void)args.len;
  return semihosting_call(SYS_OPEN, &args);
}

/*
static int32_t sys_close (FILEHANDLE fh) {
  //lint --e{934}  "Taking address of near auto variable"
  struct {
    FILEHANDLE     fh;
  } args = { fh };
  (void)args.fh;
  return semihosting_call(SYS_CLOSE, &args);
}
*/

static int32_t sys_write (FILEHANDLE fh, const uint8_t *buf, uint32_t len) {
  //lint --e{934}  "Taking address of near auto variable"
  struct {
    FILEHANDLE     fh;
    const uint8_t *buf;
    uint32_t       len;
  } args = { fh, buf, len };
  (void)args.fh;
  (void)args.buf;
  (void)args.len;
  return semihosting_call(SYS_WRITE, &args);
}

#endif


/**
  Record a single item
  \param[in]    id     event identifier (component, message with context & first/last flags)
  \param[in]    ts     timestamp
  \param[in]    val1   first data value
  \param[in]    val2   second data value
  \return       status (1=Success, 0=Failure)
*/
static uint32_t EventRecordItem (uint32_t id, uint32_t ts, uint32_t val1, uint32_t val2) {
  EventRecord_t *record;
  uint32_t cnt, i;
  uint32_t info;
  uint32_t tbit;
  uint32_t seq;

  for (cnt = EVENT_RECORD_MAX_LOCKED; cnt != 0U; cnt--) {
    i = GetRecordIndex();
    record = &EventBuffer[i & (EVENT_RECORD_COUNT - 1U)];
    seq  = ((i / EVENT_RECORD_COUNT) << EVENT_RECORD_SEQ_POS) & EVENT_RECORD_SEQ_MASK;
    info = id                                    | 
           seq                                   |
           ((ts   >> 3) & EVENT_RECORD_MSB_TS)   |
           ((val1 >> 2) & EVENT_RECORD_MSB_VAL1) |
           ((val2 >> 1) & EVENT_RECORD_MSB_VAL2) |
           EVENT_RECORD_VALID                    |
           EVENT_RECORD_LOCKED;
    info = LockRecord(&record->info, info);
    if ((info & EVENT_RECORD_LOCKED) != 0U) {
      info ^= EVENT_RECORD_LOCKED;
      info ^= EVENT_RECORD_TBIT;
      tbit  = info & EVENT_RECORD_TBIT;
      record->ts   = (ts   & ~EVENT_RECORD_TBIT) | tbit;
      record->val1 = (val1 & ~EVENT_RECORD_TBIT) | tbit;
      record->val2 = (val2 & ~EVENT_RECORD_TBIT) | tbit;
      if ((UnlockRecord(&record->info, info)) != 0U) {
        IncrementRecordsWritten();
        //lint -e{904} "Return statement before end of function"
        return 1U;
      } else {
        break;
      }
    }
  }

  IncrementRecordsDumped();
  return 0U;
}


#ifdef RTE_Compiler_EventRecorder_Semihosting
 
/**
  Record an event with variable data size to a log file
  \param[in]    id     event identifier (component number, message number)
  \param[in]    data   event data buffer
  \param[in]    len    event data length
  \param[in]    ts     timestamp
*/
static void EventRecordData_Log (uint32_t id,
                                 const void *data, uint32_t len,
                                 uint64_t ts) {
  //lint --e{934}  "Taking address of near auto variable"
  struct {
    EventRecordHead_t head;
    EventRecordData_t record;
  } event;

  event.head.type          = EVENT_TYPE_DATA;
  event.head.length        = (uint16_t)(sizeof(event.record) + len); 
  event.record.ts          = ts;
  event.record.info.id     = (uint16_t)id;
  //lint -e{9034} "Expression assigned to a narrower or different essential type"
  event.record.info.length = len;
  event.record.info.irq    = (__get_IPSR() != 0U) ? 1U : 0U;

  (void)sys_write(FileHandle, (uint8_t *)&event,  sizeof(event));
  //lint -e{9079} "conversion from pointer to void to pointer to other type"
  (void)sys_write(FileHandle,             data,   len);
}

/**
  Record an event with two 32-bit data values to a log file
  \param[in]    id     event identifier (component number, message number)
  \param[in]    val1   first data value
  \param[in]    val2   second data value
  \param[in]    ts     timestamp
*/
static void EventRecord2_Log (uint32_t id,
                              uint32_t val1, uint32_t val2,
                              uint64_t ts) {
  //lint --e{934}  "Taking address of near auto variable"
  struct {
    EventRecordHead_t head;
    EventRecord2_t    record;
  } event;

  event.head.type          = EVENT_TYPE_VAL2;
  event.head.length        = (uint16_t)sizeof(event.record);
  event.record.ts          = ts;
  event.record.info.id     = (uint16_t)id;
  event.record.info.rsrvd  = 0U;
  event.record.info.irq    = (__get_IPSR() != 0U) ? 1U : 0U;
  event.record.val1        = val1;
  event.record.val2        = val2;

  (void)sys_write(FileHandle, (uint8_t *)&event,  sizeof(event));
}
 
/**
  Record an event with four 32-bit data values a log file
  \param[in]    id     event identifier (component number, message number)
  \param[in]    val1   first data value
  \param[in]    val2   second data value
  \param[in]    val3   third data value
  \param[in]    val4   fourth data value
  \param[in]    ts     timestamp
*/
static void EventRecord4_Log (uint32_t id,
                              uint32_t val1, uint32_t val2, uint32_t val3, uint32_t val4,
                              uint64_t ts) {
  //lint --e{934}  "Taking address of near auto variable"
  struct {
    EventRecordHead_t head;
    EventRecord4_t    record;
  } event;

  event.head.type          = EVENT_TYPE_VAL4;
  event.head.length        = (uint16_t)sizeof(event.record);
  event.record.ts          = ts;
  event.record.info.id     = (uint16_t)id;
  event.record.info.rsrvd  = 0U;
  event.record.info.irq    = (__get_IPSR() != 0U) ? 1U : 0U;
  event.record.val1        = val1;
  event.record.val2        = val2;
  event.record.val3        = val3;
  event.record.val4        = val4;

  (void)sys_write(FileHandle, (uint8_t *)&event,  sizeof(event));
}

#endif


#ifdef RTE_Compiler_EventRecorder_Semihosting

/**
  Get timestamp and handle overflow

  \return       timestamp (64-bit)
*/
static uint64_t EventGetTS64 (void) {
  uint32_t ts;
  uint32_t ts_last;
  uint32_t ts_last_prev;
  uint32_t ts_overflow;

  do {
    ts_overflow = *((volatile uint32_t *)&EventStatus.ts_overflow);
    ts_last     = *((volatile uint32_t *)&EventStatus.ts_last);
    ts = EventRecorderTimerGetCount();
    if (ts < ts_last) {
      if (LockTS_Overflow() == 0U) {
        EventStatus.ts_overflow++;
        UnlockTS_Overflow();
      }
      ts_overflow++;
    } else {
      if (TS_OverflowLock != 0U) {
        ts_overflow++;
      }
    }
    ts_last_prev = UpdateTS(ts);
  } while (ts_last != ts_last_prev);

  return (ts | ((uint64_t)ts_overflow << 32));
}

#else

/**
  Get timestamp and handle overflow

  \return       timestamp (32-bit)
*/
static uint32_t EventGetTS (void) {
  uint32_t ts;
  uint32_t ts_last;
  uint32_t ts_last_prev;

  do {
    ts_last = *((volatile uint32_t *)&EventStatus.ts_last);
    ts = EventRecorderTimerGetCount();
    if (ts < ts_last) {
      if (LockTS_Overflow() == 0U) {
        EventStatus.ts_overflow++;
        UnlockTS_Overflow();
      }
    }
    ts_last_prev = UpdateTS(ts);
  } while (ts_last != ts_last_prev);

  return (ts);
}

#endif


/**
  Check event filter based on specified level and component
  \param[in]    id     event identifier (level, component number, message number)
  \return              1=Enabled, 0=Disabled
*/
__STATIC_INLINE uint32_t EventCheckFilter (uint32_t id) {
  uint32_t ret;

  if (EventStatus.state == 0U) {
    ret = 0U;
  } else {
    ret = ((uint32_t)EventFilter[(id >> (8 + 3)) & 0x7FU] >> ((id >> 8) & 0x7U)) & 1U;
  }
  return (ret);
}


/**
  Calculate CRC16-CCITT (16-bit, polynom=0x1021, init_value=0xFFFF)
  \param[in]    data  pointer to data
  \param[in]    len   data length (in bytes)
  \return             CRC16-CCITT value
*/
static uint16_t crc16_ccitt (const uint8_t *data, uint32_t len) {
  uint16_t crc;
  uint32_t n;

  crc = 0xFFFFU;
  while (len != 0U) {
    //lint -e{9049} "increment/decrement operation combined with other operation"
    crc ^= ((uint16_t)*data++ << 8);
    for (n = 8U; n != 0U; n--) {
      if ((crc & 0x8000U) != 0U) {
        crc <<= 1;
        crc  ^= 0x1021U;
      } else {
        crc <<= 1;
      }
    }
    len--;
  }

  return (crc);
}


#if   (EVENT_TIMESTAMP_SOURCE == 0)
#if ((__CORTEX_M < 3U) || (__CORTEX_M == 23U))
#ifndef _lint
#warning "Invalid Time Stamp Source selected in EventRecorderConf.h!"
#endif
static uint32_t TimeStamp __NO_INIT;
#endif

#elif (EVENT_TIMESTAMP_SOURCE == 1)

/* SysTick period in cycles */
#ifndef SYSTICK_PERIOD
#define SYSTICK_PERIOD  0x01000000U
#endif

/* SysTick variables */
static volatile uint32_t SysTickCount;
static volatile uint8_t  SysTickUpdated;

/* SysTick IRQ handler */
void SysTick_Handler (void);
void SysTick_Handler (void) {
  SysTickCount  += SYSTICK_PERIOD;
  SysTickUpdated = 1U;
}

/* Setup SysTick */
__STATIC_INLINE uint32_t SysTickSetup (void) {
  SysTickCount  = 0U;
  SysTick->LOAD = SYSTICK_PERIOD - 1U;
  SysTick->VAL  = 0U;
  SysTick->CTRL = SysTick_CTRL_ENABLE_Msk     |
                  SysTick_CTRL_TICKINT_Msk    |
                  SysTick_CTRL_CLKSOURCE_Msk;
  return 1U;
}

/* Get SysTick frequency */
__STATIC_INLINE uint32_t SysTickGetFreq (void) {
  return (SystemCoreClock);
}

/* Get SysTick count */
__STATIC_INLINE uint32_t SysTickGetCount (void) {
  uint32_t val;

  do {
    SysTickUpdated = 0U;
    val  = SysTickCount;
    val += (SYSTICK_PERIOD - 1U) - SysTick->VAL;
  } while (SysTickUpdated != 0U);

  return (val);
}

#elif (EVENT_TIMESTAMP_SOURCE == 2)

static uint8_t SysTimerState;

/* Check if SysTimer is running  */
__STATIC_INLINE uint32_t SysTimerIsRunning (void) {
  if (SysTimerState == 0U) {
    if (osKernelGetState() >= osKernelRunning) {
      SysTimerState = 1U;
    }
  }
  return (SysTimerState);
}

#endif


/**
  Setup timer hardware
  \return       status (1=Success, 0=Failure)
*/
#if (EVENT_TIMESTAMP_SOURCE < 3)
__WEAK uint32_t EventRecorderTimerSetup (void) {
#if   (EVENT_TIMESTAMP_SOURCE == 0)
  #if ((__CORTEX_M >= 3U) && (__CORTEX_M != 23U))
    CoreDebug->DEMCR |= CoreDebug_DEMCR_TRCENA_Msk;
    DWT->CTRL |= DWT_CTRL_CYCCNTENA_Msk;
    return 1U;
  #else
    TimeStamp = 0U;
    return 1U;
  #endif
#elif (EVENT_TIMESTAMP_SOURCE == 1)
  return (SysTickSetup());
#elif (EVENT_TIMESTAMP_SOURCE == 2)
  SysTimerState = 0U;
  return 1U;
#endif
}
#endif

/**
  Get timer frequency
  \return       timer frequency in Hz
*/
#if (EVENT_TIMESTAMP_SOURCE < 3)
__WEAK uint32_t EventRecorderTimerGetFreq (void) {
#if   (EVENT_TIMESTAMP_SOURCE == 0)
  #if ((__CORTEX_M >= 3U) && (__CORTEX_M != 23U))
    return (SystemCoreClock);
  #else
    return 0U;
  #endif
#elif (EVENT_TIMESTAMP_SOURCE == 1)
  return (SysTickGetFreq());
#elif (EVENT_TIMESTAMP_SOURCE == 2)
  uint32_t freq;

  if (SysTimerIsRunning() != 0U) {
    freq = osKernelGetSysTimerFreq();
  } else {
    freq = 0U;
  }
  return (freq);
#endif
}
#endif

/**
  Get timer count
  \return       timer count (32-bit)
*/
#if (EVENT_TIMESTAMP_SOURCE < 3)
__WEAK uint32_t EventRecorderTimerGetCount (void) {
#if   (EVENT_TIMESTAMP_SOURCE == 0)
  #if ((__CORTEX_M >= 3U) && (__CORTEX_M != 23U))
    return (DWT->CYCCNT);
  #else
    return (TimeStamp++);
  #endif
#elif (EVENT_TIMESTAMP_SOURCE == 1)
  return (SysTickGetCount());
#elif (EVENT_TIMESTAMP_SOURCE == 2)
  uint32_t count;

  if (SysTimerIsRunning() != 0U) {
    count = osKernelGetSysTimerCount();
  } else {
    count = 0U;
  }
  return (count);
#endif
}
#endif


/**
  Initialize Event Recorder
  \param[in]    recording   initial level mask for event record filter
  \param[in]    start       initial recording setup (1=start, 0=stop)
  \return       status (1=Success, 0=Failure)
*/
uint32_t EventRecorderInitialize (uint32_t recording, uint32_t start) {
  EventRecord_t *record;
  uint16_t crc;
  uint32_t freq;
  uint32_t ret;
  uint32_t ts;
  uint32_t n;

  EventStatus.state = 0U;
  memset(&EventFilter[0], 0, sizeof(EventFilter));

  crc = crc16_ccitt((const uint8_t *)&EventRecorderInfo, sizeof(EventRecorderInfo));

  if (EventStatus.signature != SIGNATURE) {
    EventStatus.signature  = SIGNATURE;
    EventStatus.info_crc   = crc;
    EventStatus.init_count = 1U;
  } else {
    if (EventStatus.info_crc != crc) {
      EventStatus.info_crc   = crc;
      EventStatus.init_count = 1U;
    } else {
      EventStatus.init_count++;
    }
  }

  if (EventStatus.init_count == 1U) {
    EventStatus.context         = 0U;
    EventStatus.record_index    = 0U;
    EventStatus.records_written = 0U;
    EventStatus.records_dumped  = 0U;
    memset(&EventBuffer[0], 0, sizeof(EventBuffer));
#ifdef RTE_Compiler_EventRecorder_Semihosting
    FileHandle = sys_open(EVENT_LOG_FILENAME, MODE_wb);
#endif
  } else {
    for (n = 0U; n < EVENT_RECORD_COUNT; n++) {
      record = &EventBuffer[n];
      if ((record->info & EVENT_RECORD_LOCKED) != 0U) {
        record->info &= ~(EVENT_RECORD_LOCKED | EVENT_RECORD_VALID);
      }
    }
  }

  if (EventStatus.init_count == 1U) {
    ret = EventRecorderTimerSetup();
    if (ret != 0U) {
      #if (defined(EVENT_TIMESTAMP_FREQ) && (EVENT_TIMESTAMP_FREQ != 0U))
        freq = EVENT_TIMESTAMP_FREQ;
      #else
        freq = EventRecorderTimerGetFreq();
      #endif
    } else {
      freq = 0U;
    }
    EventStatus.ts_freq     = freq;
    EventStatus.ts_last     = 0U;
    EventStatus.ts_overflow = 0U;
    TS_OverflowLock         = 0U;
  } else {
#if    (EVENT_TIMESTAMP_SOURCE == 0)
  #if ((__CORTEX_M >= 3U) && (__CORTEX_M != 23U))
    ret = EventRecorderTimerSetup();
  #else
    ret = 1U;
  #endif
#elif ((EVENT_TIMESTAMP_SOURCE >= 1) && (EVENT_TIMESTAMP_SOURCE <= 3))
    ret = EventRecorderTimerSetup();
    if (ret != 0U) {
      #if (defined(EVENT_TIMESTAMP_FREQ) && (EVENT_TIMESTAMP_FREQ != 0U))
        freq = EVENT_TIMESTAMP_FREQ;
      #else
        freq = EventRecorderTimerGetFreq();
      #endif
    } else {
      freq = 0U;
    }
    EventStatus.ts_freq     = freq;
    EventStatus.ts_last     = 0U;
    EventStatus.ts_overflow = 0U;
    TS_OverflowLock         = 0U;
#else
    ret = 1U;
#endif
  }

  if (ret != 0U) {

    (void)EventRecorderEnable(recording,      0x00U,            0xFEU);
    (void)EventRecorderEnable(EventRecordAll, EvtStatistics_No, EvtStatistics_No);
    (void)EventRecorderEnable(EventRecordOp,  EvtPrintf_No,     EvtPrintf_No);

#ifdef RTE_Compiler_EventRecorder_Semihosting
    uint64_t ts64 = EventGetTS64();
    EventRecord2_Log(ID_EVENT_INIT & EVENT_RECORD_ID_MASK, EventStatus.init_count, EventStatus.ts_freq, ts64);
    ts = (uint32_t)ts64;
#else
    ts = EventGetTS();
#endif

    (void)EventRecordItem(ID_EVENT_INIT, ts, EventStatus.init_count, EventStatus.ts_freq);

    if (start != 0U) {
      (void)EventRecorderStart();
    }
  }

  return (ret);
}

/**
  Enable recording of events with specified level and component range
  \param[in]    recording   level mask for event record filter
  \param[in]    comp_start  first component number of range
  \param[in]    comp_end    last Component number of range
  \return       status (1=Success, 0=Failure)
*/
uint32_t EventRecorderEnable (uint32_t recording, uint32_t comp_start, uint32_t comp_end) {
  uint32_t ofs;
  uint32_t i, j;

  if ((comp_start >= 0xFFU) || (comp_end >= 0xFFU)) {
    //lint -e{904} "Return statement before end of function"
    return 0U;
  }

  ofs = 0U;
  for (i = 0U; i < 4U; i++) {
    if ((recording & (1UL << i)) != 0U) {
      for (j = comp_start; j <= comp_end; j++) {
        EventFilter[ofs + (j >> 3)] |= (1U << (j & 0x7U));
      }
    }
    ofs += 32U;
  }

  return 1U;
}

/**
  Disable recording of events with specified level and component range
  \param[in]    recording   level mask for event record filter
  \param[in]    comp_start  first component number of range
  \param[in]    comp_end    last Component number of range
  \return       status (1=Success, 0=Failure)
*/
uint32_t EventRecorderDisable (uint32_t recording, uint32_t comp_start, uint32_t comp_end) {
  uint32_t ofs;
  uint32_t i, j;

  if ((comp_start >= 0xFFU) || (comp_end >= 0xFFU)) {
    //lint -e{904} "Return statement before end of function"
    return 0U;
  }

  ofs = 0U;
  for (i = 0U; i < 4U; i++) {
    if ((recording & (1UL << i)) != 0U) {
      for (j = comp_start; j <= comp_end; j++) {
        EventFilter[ofs + (j >> 3)] &= ~(1U << (j & 0x7U));
      }
    }
    ofs += 32U;
  }

  return 1U;
}

/**
  Start event recording
  \return       status (1=Success, 0=Failure)
*/
uint32_t EventRecorderStart (void) {
  uint32_t ts;

  if (EventStatus.state != 0U) {
    //lint -e{904} "Return statement before end of function"
    return 1U;
  }
  EventStatus.state = 1U;

#ifdef RTE_Compiler_EventRecorder_Semihosting
  uint64_t ts64 = EventGetTS64();
  EventRecord2_Log(ID_EVENT_START & EVENT_RECORD_ID_MASK, 0U, 0U, ts64);
  ts = (uint32_t)ts64;
#else
  ts = EventGetTS();
#endif

  (void)EventRecordItem(ID_EVENT_START, ts, 0U, 0U);

  return 1U;
}

/**
  Stop event recording
  \return       status (1=Success, 0=Failure)
*/
uint32_t EventRecorderStop (void) {
  uint32_t ts;

  if (EventStatus.state == 0U) {
    //lint -e{904} "Return statement before end of function"
    return 1U;
  }
  EventStatus.state = 0U;

#ifdef RTE_Compiler_EventRecorder_Semihosting
  uint64_t ts64 = EventGetTS64();
  EventRecord2_Log(ID_EVENT_STOP & EVENT_RECORD_ID_MASK, 0U, 0U, ts64);
  ts = (uint32_t)ts64;
#else
  ts = EventGetTS();
#endif

  (void)EventRecordItem(ID_EVENT_STOP, ts, 0U, 0U);

  return 1U;
}

/**
  Update Event Recorder timestamp clock
  \return       status (1=Success, 0=Failure)
*/
uint32_t EventRecorderClockUpdate (void) {
  uint32_t ts;

  EventStatus.ts_freq = EventRecorderTimerGetFreq();

#ifdef RTE_Compiler_EventRecorder_Semihosting
  uint64_t ts64 = EventGetTS64();
  EventRecord2_Log(ID_EVENT_CLOCK & EVENT_RECORD_ID_MASK, EventStatus.ts_freq, 0U, ts64);
  ts = (uint32_t)ts64;
#else
  ts = EventGetTS();
#endif

  (void)EventRecordItem(ID_EVENT_CLOCK, ts, EventStatus.ts_freq, 0U);

  return 1U;
}

/**
  Record an event with variable data size
  \param[in]    id     event identifier (level, component number, message number)
  \param[in]    data   event data buffer
  \param[in]    len    event data length
  \return       status (1=Success, 0=Failure)
*/
uint32_t EventRecordData (uint32_t id, const void *data, uint32_t len) {
  //lint --e{934}  "Taking address of near auto variable"
  //lint --e{9016} "pointer arithmetic other than array indexing used"
  const uint8_t *dptr;
  uint32_t ts;
  uint32_t ctx;
  uint32_t val[2];
  uint32_t ret;

  if ((data == NULL) || (len > EVENT_DATA_MAX_LENGTH)) {
    //lint -e{904} "Return statement before end of function"
    return 0U;
  }

  if (EventCheckFilter(id) == 0U) {
    //lint -e{904} "Return statement before end of function"
    return 1U;
  }

#ifdef RTE_Compiler_EventRecorder_Semihosting
  uint64_t ts64 = EventGetTS64();
  EventRecordData_Log(id, data, len, ts64);
  ts = (uint32_t)ts64;
#else
  ts = EventGetTS();
#endif

  id &= EVENT_RECORD_ID_MASK;
  id |= (__get_IPSR() != 0U) ? EVENT_RECORD_IRQ : 0U;
  //lint -e{9079} -e{9087} "conversion from pointer to void to pointer to other type"
  dptr = (const uint8_t *)data;

  if (len == 0U) {
    ret = EventRecordItem(id, ts, 0U, 0U);
    //lint -e{904} "Return statement before end of function"
    return (ret);
  }

  if (len <= 8U) {
    val[0] = 0U;
    val[1] = 0U;
    memcpy(val, dptr, len);
    id |= (len << EVENT_RECORD_DLEN_POS) & EVENT_RECORD_DLEN_MASK;
    ret = EventRecordItem(id | EVENT_RECORD_FIRST | EVENT_RECORD_LAST, ts, val[0], val[1]);
    //lint -e{904} "Return statement before end of function"
    return (ret);
  }

  ctx = (GetContext() << EVENT_RECORD_CTX_POS) & EVENT_RECORD_CTX_MASK;

  memcpy(val, dptr, 8U);
  dptr += 8U;
  len  -= 8U;
  id |= ctx;
  ret = EventRecordItem(id | EVENT_RECORD_FIRST, ts, val[0], val[1]);
  if (ret == 0U) {
    //lint -e{904} "Return statement before end of function"
    return 0U;
  }

  //lint -e{9044} "function parameter modified"
  id = 0xFF01U | ctx;

  while (len > 8U) {
    memcpy(val, dptr, 8U);
    dptr += 8U;
    len  -= 8U;
    ret = EventRecordItem(id, ts, val[0], val[1]);
    id++;
    if (ret == 0U) {
      //lint -e{904} "Return statement before end of function"
      return 0U;
    }
  }

  val[0] = 0U;
  val[1] = 0U;
  memcpy(val, dptr, len);
  id &= ~0xFF00U;
  id |= len << 8;
  ret = EventRecordItem(id | EVENT_RECORD_LAST, ts, val[0], val[1]);

  return (ret);
}

/**
  Record an event with two 32-bit data values
  \param[in]    id     event identifier (level, component number, message number)
  \param[in]    val1   first data value
  \param[in]    val2   second data value
  \return       status (1=Success, 0=Failure)
*/
uint32_t EventRecord2 (uint32_t id, uint32_t val1, uint32_t val2) {
  uint32_t ts;
  uint32_t ret;

  if (EventCheckFilter(id) == 0U) {
    //lint -e{904} "Return statement before end of function"
    return 1U;
  }

#ifdef RTE_Compiler_EventRecorder_Semihosting
  uint64_t ts64 = EventGetTS64();
  EventRecord2_Log(id, val1, val2, ts64);
  ts = (uint32_t)ts64;
#else
  ts = EventGetTS();
#endif

  id &= EVENT_RECORD_ID_MASK;
  id |= (__get_IPSR() != 0U) ? EVENT_RECORD_IRQ : 0U;

  ret = EventRecordItem(id | EVENT_RECORD_FIRST | EVENT_RECORD_LAST, ts, val1, val2);

  return (ret);
}

/**
  Record an event with four 32-bit data values
  \param[in]    id     event identifier (level, component number, message number)
  \param[in]    val1   first data value
  \param[in]    val2   second data value
  \param[in]    val3   third data value
  \param[in]    val4   fourth data value
  \return       status (1=Success, 0=Failure)
*/
uint32_t EventRecord4 (uint32_t id,
                       uint32_t val1, uint32_t val2, uint32_t val3, uint32_t val4) {
  uint32_t ts;
  uint32_t ctx;
  uint32_t ret;

  if (EventCheckFilter(id) == 0U) {
    //lint -e{904} "Return statement before end of function"
    return 1U;
  }

#ifdef RTE_Compiler_EventRecorder_Semihosting
  uint64_t ts64 = EventGetTS64();
  EventRecord4_Log(id, val1, val2, val3, val4, ts64);
  ts = (uint32_t)ts64;
#else
  ts = EventGetTS();
#endif

  id &= EVENT_RECORD_ID_MASK;
  id |= (__get_IPSR() != 0U) ? EVENT_RECORD_IRQ : 0U;
  ctx = (GetContext() << EVENT_RECORD_CTX_POS) & EVENT_RECORD_CTX_MASK;

  ret = EventRecordItem(id | ctx | EVENT_RECORD_FIRST, ts, val1, val2);
  if (ret == 0U) {
    //lint -e{904} "Return statement before end of function"
    return 0U;
  }
  ret = EventRecordItem(1U | ctx | EVENT_RECORD_LAST,  ts, val3, val4);

  return (ret);
}
