/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
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

#include <stdint.h>

// Fault IDs for fault triggering
#define ARM_FAULT_ID_MEM_DATA                          (1U)
#define ARM_FAULT_ID_BUS_DATA_PRECISE                  (2U)
#define ARM_FAULT_ID_BUS_DATA_IMPRECISE                (3U)
#define ARM_FAULT_ID_BUS_INSTRUCTION                   (4U)
#define ARM_FAULT_ID_USG_UNDEFINED_INSTRUCTION         (5U)
#define ARM_FAULT_ID_USG_DIV_0                         (6U)

// ARM_FaultTrigger function ---------------------------------------------------

/// Trigger a fault.
/// \param[in]    fault_id    Fault Id of the fault to be triggered
extern void ARM_FaultTrigger (uint32_t fault_id);
