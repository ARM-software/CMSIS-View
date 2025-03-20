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
 *
 *      Name:    retarget_stdio.c
 *      Purpose: Retarget STDIO to USART0
 *
 */

#include "Driver_USART.h"

extern int stdio_init     (void);
extern int stderr_putchar (int ch);
extern int stdout_putchar (int ch);
extern int stdin_getchar  (void);

#define USART_DRV_NUM           0
#define USART_BAUDRATE          115200

#define _USART_Driver_(n)  Driver_USART##n
#define  USART_Driver_(n) _USART_Driver_(n)

extern ARM_DRIVER_USART  USART_Driver_(USART_DRV_NUM);
#define ptrUSART       (&USART_Driver_(USART_DRV_NUM))

/**
  Initialize stdio

  \return          0 on success, or -1 on error.
*/
int stdio_init (void) {

  if (ptrUSART->Initialize(NULL) != ARM_DRIVER_OK) {
    return -1;
  }

  if (ptrUSART->PowerControl(ARM_POWER_FULL) != ARM_DRIVER_OK) {
    return -1;
  }

  if (ptrUSART->Control(ARM_USART_MODE_ASYNCHRONOUS |
                        ARM_USART_DATA_BITS_8       |
                        ARM_USART_PARITY_NONE       |
                        ARM_USART_STOP_BITS_1       |
                        ARM_USART_FLOW_CONTROL_NONE,
                        USART_BAUDRATE) != ARM_DRIVER_OK) {
    return -1;
  }

  if (ptrUSART->Control(ARM_USART_CONTROL_RX, 1U) != ARM_DRIVER_OK) {
    return -1;
  }

  return 0;
}

/**
  Put a character to the stderr

  \param[in]   ch  Character to output
  \return          The character written, or -1 on error.
*/
int stderr_putchar (int ch) {
  uint8_t buf[1];

  buf[0] = (uint8_t)ch;

  if (ptrUSART->Send(buf, 1U) != ARM_DRIVER_OK) {
    return -1;
  }

  while (ptrUSART->GetTxCount() != 1U);

  return ch;
}

/**
  Put a character to the stdout

  \param[in]   ch  Character to output
  \return          The character written, or -1 on write error.
*/
int stdout_putchar (int ch) {
  uint8_t buf[1];

  buf[0] = (uint8_t)ch;

  if (ptrUSART->Send(buf, 1U) != ARM_DRIVER_OK) {
    return -1;
  }

  while (ptrUSART->GetTxCount() != 1U);

  return ch;
}

/**
  Get a character from the stdio

  \return     The next character from the input, or -1 on error.
*/
int stdin_getchar (void) {
  uint8_t buf[1];

  if (ptrUSART->Receive(buf, 1U) != ARM_DRIVER_OK) {
    return -1;
  }

  while (ptrUSART->GetRxCount() != 1U);

  return (int)buf[0];
}
