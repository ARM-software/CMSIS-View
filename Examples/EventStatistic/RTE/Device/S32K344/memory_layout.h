/*---------------------------------------------------------------------------
 * Name:    S32K3xx_memmap.h
 * Purpose: S32K3xx Memory Mapping
 * Rev.:    1.0.0
 *---------------------------------------------------------------------------*/
/*
 * Copyright (c) 2021 Arm Limited or its affiliates. All rights reserved.
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

#ifndef S32K3xx_MEMMAP_H
#define S32K3xx_MEMMAP_H

/* default memory mapping used for S32K344 cores
 *  BOOT_HEADER    (r)  :    Start: 0x10000000,   Size: 0x00000100
 *
 *  CM7_0_FLASH   (rx) :     Start: 0x00400000,   Size: 0x003D0000   0x00400000
 *  CM7_0_RAM     (rwx):     Start: 0x20400000,   Size: 0x00050000   (2 x 160kB)
 *  CM7_0_DTCM    (rwx):     Start: 0x20000000,   Size: 0x00020000
 */

#define BOOT_HEADER_START   0x10000000
#define BOOT_HEADER_SIZE    0x00000100

/* CM7 0 */
#define CM7_0_FLASH_START   0x00400000
#define CM7_0_FLASH_SIZE    0x003D0000

#define CM7_0_RAM_START     0x20400000
#define CM7_0_RAM_SIZE      0x00050000

#define CM7_0_DTCM_START    0x20000000
#define CM7_0_DTCM_SIZE     0x00020000

#define  FLASH_START__          CM7_0_FLASH_START
#define FLASH_SIZE__           CM7_0_FLASH_SIZE
#define RAM_START__            CM7_0_RAM_START
#define RAM_SIZE__             CM7_0_RAM_SIZE
#define DTCM_START__           CM7_0_DTCM_START
#define DTCM_SIZE__            CM7_0_DTCM_SIZE

#define __ROM1_BASE            FLASH_START__
#define __ROM1_SIZE            FLASH_SIZE__
#define __RAM1_BASE            RAM_START__
#define __RAM1_SIZE            RAM_SIZE__
#define __RAM2_BASE            DTCM_START__
#define __RAM2_SIZE            DTCM_SIZE__

#endif /* S32K3xx_MEMMAP_H */
