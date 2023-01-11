#ifndef MEMORY_LAYOUT_H
#define MEMORY_LAYOUT_H

/*
;-------- <<< Use Configuration Wizard in Context Menu >>> -------------------
*/

/*--------------------- Flash Configuration ----------------------------------
; <h> Flash Configuration
;   <o0> Flash Base Address <0x0-0xFFFFFFFF:8>
;   <o1> Flash Size (in Bytes) <0x0-0xFFFFFFFF:8>
; </h>
 *----------------------------------------------------------------------------*/
#define __ROM_BASE          0x00000000
#define __ROM_SIZE          0x00080000

/*--------------------- Embedded RAM Configuration ---------------------------
; <h> RAM Configuration
;   <o0> RAM Base Address    <0x0-0xFFFFFFFF:8>
;   <o1> RAM Size (in Bytes) <0x0-0xFFFFFFFF:8>
;   <o2> Uninit RAM Size (in Bytes) <0x0-0xFFFFFFFF:8>
; </h>
 *----------------------------------------------------------------------------*/
#define __RAM_BASE          0x20000000
#define __RAM_SIZE          0x00040000
#define __RAM_NOINIT_SIZE   0x00001000

/*--------------------- Stack / Heap Configuration ---------------------------
; <h> Stack / Heap Configuration
;   <o0> Stack Size (in Bytes) <0x0-0xFFFFFFFF:8>
;   <o1> Heap Size (in Bytes) <0x0-0xFFFFFFFF:8>
; </h>
 *----------------------------------------------------------------------------*/
#define __STACK_SIZE        0x00000200
#define __HEAP_SIZE         0x00000C00

/*--------------------- CMSE Veneer Configuration ---------------------------
; <h> CMSE Veneer Configuration
;   <o0>  CMSE Veneer Size (in Bytes) <0x0-0xFFFFFFFF:32>
; </h>
 *----------------------------------------------------------------------------*/
#define __CMSEVENEER_SIZE    0x200

/*
;------------- <<< end of configuration section >>> ---------------------------
*/

#endif
