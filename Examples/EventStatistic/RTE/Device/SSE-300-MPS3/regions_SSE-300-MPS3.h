#ifndef REGIONS_SSE-300-MPS3_H
#define REGIONS_SSE-300-MPS3_H


//-------- <<< Use Configuration Wizard in Context Menu >>> --------------------

// <h>ROM Configuration
// =======================
// </h>

// <h>RAM Configuration
// =======================
// <h> ITCM_NS
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x00000000
#define __RAM0_BASE 0x00000000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00080000
#define __RAM0_SIZE 0x00080000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM0_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM0_NOINIT 0
// </h>

// <h> SRAM_NS
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x01000000
#define __RAM1_BASE 0x01000000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00100000
#define __RAM1_SIZE 0x00100000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM1_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM1_NOINIT 0
// </h>

// <h> DTCM0_NS
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x20000000
#define __RAM2_BASE 0x20000000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00020000
#define __RAM2_SIZE 0x00020000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM2_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM2_NOINIT 0
// </h>

// <h> DTCM1_NS
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x20020000
#define __RAM3_BASE 0x20020000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00020000
#define __RAM3_SIZE 0x00020000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM3_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM3_NOINIT 0
// </h>

// <h> DTCM2_NS
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x20040000
#define __RAM4_BASE 0x20040000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00020000
#define __RAM4_SIZE 0x00020000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM4_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM4_NOINIT 0
// </h>

// <h> DTCM3_NS
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x20060000
#define __RAM5_BASE 0x20060000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00020000
#define __RAM5_SIZE 0x00020000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM5_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM5_NOINIT 0
// </h>

// <h> ISRAM0_NS
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x21000000
#define __RAM6_BASE 0x21000000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00100000
#define __RAM6_SIZE 0x00100000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM6_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM6_NOINIT 0
// </h>

// <h> ISRAM1_NS
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x21100000
#define __RAM7_BASE 0x21100000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00100000
#define __RAM7_SIZE 0x00100000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM7_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM7_NOINIT 0
// </h>

// <h> QSPI_SRAM_NS
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x28000000
#define __RAM8_BASE 0x28000000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00800000
#define __RAM8_SIZE 0x00800000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM8_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM8_NOINIT 0
// </h>

// <h> ITCM_S
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x10000000
#define __RAM9_BASE 0x10000000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00080000
#define __RAM9_SIZE 0x00080000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM9_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM9_NOINIT 0
// </h>

// <h> SRAM_S
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x11000000
#define __RAM10_BASE 0x11000000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00100000
#define __RAM10_SIZE 0x00100000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM10_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM10_NOINIT 0
// </h>

// <h> DTCM0_S
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x30000000
#define __RAM11_BASE 0x30000000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00020000
#define __RAM11_SIZE 0x00020000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM11_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM11_NOINIT 0
// </h>

// <h> DTCM1_S
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x30020000
#define __RAM12_BASE 0x30020000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00020000
#define __RAM12_SIZE 0x00020000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM12_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM12_NOINIT 0
// </h>

// <h> DTCM2_S
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x30040000
#define __RAM13_BASE 0x30040000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00020000
#define __RAM13_SIZE 0x00020000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM13_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM13_NOINIT 0
// </h>

// <h> DTCM3_S
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x30060000
#define __RAM14_BASE 0x30060000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00020000
#define __RAM14_SIZE 0x00020000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM14_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM14_NOINIT 0
// </h>

// <h> ISRAM0_S
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x31000000
#define __RAM15_BASE 0x31000000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00100000
#define __RAM15_SIZE 0x00100000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM15_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM15_NOINIT 0
// </h>

// <h> ISRAM1_S
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x31100000
#define __RAM16_BASE 0x31100000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00100000
#define __RAM16_SIZE 0x00100000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM16_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM16_NOINIT 0
// </h>

// <h> QSPI_SRAM_S
//   <o> Base address <0x0-0xFFFFFFFF:8>
//   <i> Defines base address of memory region.
//   <i> Default: 0x38000000
#define __RAM17_BASE 0x38000000
//   <o> Region size [bytes] <0x0-0xFFFFFFFF:8>
//   <i> Defines size of memory region.
//   <i> Default: 0x00800000
#define __RAM17_SIZE 0x00800000
//   <q>Default region
//   <i> Enables memory region globally for the application.
#define __RAM17_DEFAULT 1
//   <q>No zero initialize
//   <i> Excludes region from zero initialization.
#define __RAM17_NOINIT 0
// </h>

// </h>


#endif /* REGIONS_SSE-300-MPS3_H */
