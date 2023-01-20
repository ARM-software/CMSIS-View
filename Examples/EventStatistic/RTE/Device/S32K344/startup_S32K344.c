/*---------------------------------------------------------------------------
 * Name:    startup_S32K344.c
 * Purpose: S32K344 CMSIS Core Device Startup File
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

#include "RTE_Components.h"             // Component selection
#include CMSIS_device_header
#include "memory_layout.h"

/*----------------------------------------------------------------------------
  External References
 *----------------------------------------------------------------------------*/
extern uint32_t __INITIAL_SP;

extern __NO_RETURN void __PROGRAM_START(void);

/*----------------------------------------------------------------------------
  Internal References
 *----------------------------------------------------------------------------*/
            void Reset_Handler  (void) __attribute__ ((weak));
__NO_RETURN void Reset_Handler_C(void) __attribute__ ((weak));
            void Default_Handler(void) __attribute__ ((weak));

/*----------------------------------------------------------------------------
  Exception / Interrupt Handler
 *----------------------------------------------------------------------------*/
void RESERVED_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));

/* Exceptions */
void NMI_Handler                                  (void) __attribute__ ((weak, alias("Default_Handler")));
void HardFault_Handler                            (void) __attribute__ ((weak));
void MemManage_Handler                            (void) __attribute__ ((weak, alias("Default_Handler")));
void BusFault_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void UsageFault_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void SVC_Handler                                  (void) __attribute__ ((weak, alias("Default_Handler")));
void DebugMon_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void PendSV_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void SysTick_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));

  /* Interrupts */
void INT0_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void INT1_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void INT2_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void INT3_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD0_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD1_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD2_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD3_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD4_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD5_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD6_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD7_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD8_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD9_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD10_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD11_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD12_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD13_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD14_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD15_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD16_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD17_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD18_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD19_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD20_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD21_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD22_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD23_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD24_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD25_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD26_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD27_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD28_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD29_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD30_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void DMATCD31_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void ERM_0_Handler                                (void) __attribute__ ((weak, alias("Default_Handler")));
void ERM_1_Handler                                (void) __attribute__ ((weak, alias("Default_Handler")));
void MCM_Handler                                  (void) __attribute__ ((weak, alias("Default_Handler")));
void STM0_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void STM1_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void SWT0_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void CTI0_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void FLASH_0_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void FLASH_1_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void FLASH_2_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void RGM_Handler                                  (void) __attribute__ ((weak, alias("Default_Handler")));
void PMC_Handler                                  (void) __attribute__ ((weak, alias("Default_Handler")));
void SIUL_0_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void SIUL_1_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void SIUL_2_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void SIUL_3_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS0_0_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS0_1_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS0_2_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS0_3_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS0_4_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS0_5_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS1_0_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS1_1_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS1_2_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS1_3_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS1_4_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS1_5_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS2_0_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS2_1_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS2_2_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS2_3_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS2_4_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void EMIOS2_5_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void WKPU_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void CMU0_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void CMU1_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void CMU2_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void BCTU_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void LCU0_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void LCU1_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void PIT0_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void PIT1_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void PIT2_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void RTC_Handler                                  (void) __attribute__ ((weak, alias("Default_Handler")));
void EMAC_0_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void EMAC_1_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void EMAC_2_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void EMAC_3_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN0_0_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN0_1_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN0_2_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN0_3_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN1_0_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN1_1_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN1_2_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN2_0_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN2_1_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN2_2_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN3_0_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN3_1_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN4_0_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN4_1_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN5_0_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FlexCAN5_1_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void FLEXIO_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART0_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART1_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART2_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART3_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART4_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART5_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART6_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART7_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART8_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART9_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART10_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART11_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART12_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART13_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART14_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART15_Handler                             (void) __attribute__ ((weak, alias("Default_Handler")));
void LPI2C0_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void LPI2C1_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void LPSPI0_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void LPSPI1_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void LPSPI2_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void LPSPI3_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void LPSPI4_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void LPSPI5_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void QSPI_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void SAI0_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void SAI1_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void JDC_Handler                                  (void) __attribute__ ((weak, alias("Default_Handler")));
void ADC0_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void ADC1_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void ADC2_Handler                                 (void) __attribute__ ((weak, alias("Default_Handler")));
void LPCMP0_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void LPCMP1_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void LPCMP2_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void FCCU_0_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void FCCU_1_Handler                               (void) __attribute__ ((weak, alias("Default_Handler")));
void STCU_LBIST_MBIST_Handler                     (void) __attribute__ ((weak, alias("Default_Handler")));
void HSE_MU0_TX_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void HSE_MU0_RX_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void HSE_MU0_ORED_Handler                         (void) __attribute__ ((weak, alias("Default_Handler")));
void HSE_MU1_TX_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void HSE_MU1_RX_Handler                           (void) __attribute__ ((weak, alias("Default_Handler")));
void HSE_MU1_ORED_Handler                         (void) __attribute__ ((weak, alias("Default_Handler")));
void SoC_PLL_Handler                              (void) __attribute__ ((weak, alias("Default_Handler")));



/*----------------------------------------------------------------------------
  Exception / Interrupt Vector table
 *----------------------------------------------------------------------------*/

#if defined ( __GNUC__ )
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wpedantic"
#endif

extern const VECTOR_TABLE_Type __VECTOR_TABLE[NUMBER_OF_INT_VECTORS];
       const VECTOR_TABLE_Type __VECTOR_TABLE[NUMBER_OF_INT_VECTORS] __VECTOR_TABLE_ATTRIBUTE = {
  (VECTOR_TABLE_Type)(&__INITIAL_SP),            /*     Initial Stack Pointer */
  Reset_Handler,                                 /*     Reset Handler */
  NMI_Handler,                                   /* -14 NMI Handler */
  HardFault_Handler,                             /* -13 Hard Fault Handler */
  MemManage_Handler,                             /* -12 MPU Fault Handler */
  BusFault_Handler,                              /* -11 Bus Fault Handler */
  UsageFault_Handler,                            /* -10 Usage Fault Handler */
  RESERVED_Handler,                              /*     Reserved */
  RESERVED_Handler,                              /*     Reserved */
  RESERVED_Handler,                              /*     Reserved */
  RESERVED_Handler,                              /*     Reserved */
  SVC_Handler,                                   /*  -5 SVCall Handler */
  DebugMon_Handler,                              /*  -4 Debug Monitor Handler */
  RESERVED_Handler,                              /*     Reserved */
  PendSV_Handler,                                /*  -2 PendSV Handler */
  SysTick_Handler,                               /*  -1 SysTick Handler */

  /* Interrupts */
  INT0_Handler,                                  /*   0 CPU to CPU int0 */
  INT1_Handler,                                  /*   1 CPU to CPU int1 */
  INT2_Handler,                                  /*   2 CPU to CPU int2 */
  INT3_Handler,                                  /*   3 CPU to CPU int3 */
  DMATCD0_Handler,                               /*   4 DMA transfer complete and error CH0 */
  DMATCD1_Handler,                               /*   5 DMA transfer complete and error CH1 */
  DMATCD2_Handler,                               /*   6 DMA transfer complete and error CH2 */
  DMATCD3_Handler,                               /*   7 DMA transfer complete and error CH3 */
  DMATCD4_Handler,                               /*   8 DMA transfer complete and error CH4 */
  DMATCD5_Handler,                               /*   9 DMA transfer complete and error CH5 */
  DMATCD6_Handler,                               /*  10 DMA transfer complete and error CH6 */
  DMATCD7_Handler,                               /*  11 DMA transfer complete and error CH7 */
  DMATCD8_Handler,                               /*  12 DMA transfer complete and error CH8 */
  DMATCD9_Handler,                               /*  13 DMA transfer complete and error CH9 */
  DMATCD10_Handler,                              /*  14 DMA transfer complete and error CH10 */
  DMATCD11_Handler,                              /*  15 DMA transfer complete and error CH11 */
  DMATCD12_Handler,                              /*  16 DMA transfer complete and error CH12 */
  DMATCD13_Handler,                              /*  17 DMA transfer complete and error CH13 */
  DMATCD14_Handler,                              /*  18 DMA transfer complete and error CH14 */
  DMATCD15_Handler,                              /*  19 DMA transfer complete and error CH15 */
  DMATCD16_Handler,                              /*  20 DMA transfer complete and error CH16 */
  DMATCD17_Handler,                              /*  21 DMA transfer complete and error CH17 */
  DMATCD18_Handler,                              /*  22 DMA transfer complete and error CH18 */
  DMATCD19_Handler,                              /*  23 DMA transfer complete and error CH19 */
  DMATCD20_Handler,                              /*  24 DMA transfer complete and error CH20 */
  DMATCD21_Handler,                              /*  25 DMA transfer complete and error CH21 */
  DMATCD22_Handler,                              /*  26 DMA transfer complete and error CH22 */
  DMATCD23_Handler,                              /*  27 DMA transfer complete and error CH23 */
  DMATCD24_Handler,                              /*  28 DMA transfer complete and error CH24 */
  DMATCD25_Handler,                              /*  29 DMA transfer complete and error CH25 */
  DMATCD26_Handler,                              /*  30 DMA transfer complete and error CH26 */
  DMATCD27_Handler,                              /*  31 DMA transfer complete and error CH27 */
  DMATCD28_Handler,                              /*  32 DMA transfer complete and error CH28 */
  DMATCD29_Handler,                              /*  33 DMA transfer complete and error CH29 */
  DMATCD30_Handler,                              /*  34 DMA transfer complete and error CH30 */
  DMATCD31_Handler,                              /*  35 DMA transfer complete and error CH31 */
  ERM_0_Handler,                                 /*  36 Single bit ECC error */
  ERM_1_Handler,                                 /*  37 Multi bit ECC error */
  MCM_Handler,                                   /*  38 Multi bit ECC error */
  STM0_Handler,                                  /*  39 Single interrupt vector for all four channels */
  STM1_Handler,                                  /*  40 Single interrupt vector for all four channels */
  RESERVED_Handler,                              /*  41 Reserved */
  SWT0_Handler,                                  /*  42 Platform watchdog initial time-out */
  RESERVED_Handler,                              /*  43 Reserved */
  RESERVED_Handler,                              /*  44 Reserved */
  CTI0_Handler,                                  /*  45 CTI Interrupt 0 */
  RESERVED_Handler,                              /*  46 Reserved */
  RESERVED_Handler,                              /*  47 Reserved */
  FLASH_0_Handler,                               /*  48 Program or erase operation is completed */
  FLASH_1_Handler,                               /*  49 Main watchdog timeout interrupt */
  FLASH_2_Handler,                               /*  50 Alternate watchdog timeout interrupt */
  RGM_Handler,                                   /*  51 Interrupt request to the system */
  PMC_Handler,                                   /*  52 One interrupt for all LVDs, One interrupt for all HVDs */
  SIUL_0_Handler,                                /*  53 External Interrupt Vector 0 */
  SIUL_1_Handler,                                /*  54 External Interrupt Vector 1 */
  SIUL_2_Handler,                                /*  55 External Interrupt Vector 2 */
  SIUL_3_Handler,                                /*  56 External Interrupt Vector 3 */
  RESERVED_Handler,                              /*  57 Reserved */
  RESERVED_Handler,                              /*  58 Reserved */
  RESERVED_Handler,                              /*  59 Reserved */
  RESERVED_Handler,                              /*  60 Reserved */
  EMIOS0_0_Handler,                              /*  61 Interrupt request 23,22,21,20 */
  EMIOS0_1_Handler,                              /*  62 Interrupt request 19,18,17,16 */
  EMIOS0_2_Handler,                              /*  63 Interrupt request 15,14,13,12 */
  EMIOS0_3_Handler,                              /*  64 Interrupt request 11,10,9,8 */
  EMIOS0_4_Handler,                              /*  65 Interrupt request 7,6,5,4 */
  EMIOS0_5_Handler,                              /*  66 Interrupt request 3,2,1,0 */
  RESERVED_Handler,                              /*  67 Reserved */
  RESERVED_Handler,                              /*  68 Reserved */
  EMIOS1_0_Handler,                              /*  69 Interrupt request 23,22,21,20 */
  EMIOS1_1_Handler,                              /*  70 Interrupt request 19,18,17,16 */
  EMIOS1_2_Handler,                              /*  71 Interrupt request 15,14,13,12 */
  EMIOS1_3_Handler,                              /*  72 Interrupt request 11,10,9,8 */
  EMIOS1_4_Handler,                              /*  73 Interrupt request 7,6,5,4 */
  EMIOS1_5_Handler,                              /*  74 Interrupt request 3,2,1,0 */
  RESERVED_Handler,                              /*  75 Reserved */
  RESERVED_Handler,                              /*  76 Reserved */
  EMIOS2_0_Handler,                              /*  77 Interrupt request 23,22,21,20 */
  EMIOS2_1_Handler,                              /*  78 Interrupt request 19,18,17,16 */
  EMIOS2_2_Handler,                              /*  79 Interrupt request 15,14,13,12 */
  EMIOS2_3_Handler,                              /*  80 Interrupt request 11,10,9,8 */
  EMIOS2_4_Handler,                              /*  81 Interrupt request 7,6,5,4 */
  EMIOS2_5_Handler,                              /*  82 Interrupt request 3,2,1,0 */
  WKPU_Handler,                                  /*  83 Interrupts from pad group 0,1,2,3, Interrupts from pad group 0_64, Interrupts from pad group 1_64, Interrupts from pad group 2_64, Interrupts from pad group 3_64 */
  CMU0_Handler,                                  /*  84 CMU0 interrupt */
  CMU1_Handler,                                  /*  85 CMU1 interrupt */
  CMU2_Handler,                                  /*  86 CMU2 interrupt */
  BCTU_Handler,                                  /*  87 An interrupt is requested when a conversion is issued to the ADC, An interrupt is requested when new data is available from ADC0 conversion, An interrupt is requested when new data is available from ADC1 conversion, An interrupt is requested when new data is available from ADC2 conversion, An interrupt is requested when the last command of a list is issued to the ADC,An Interrupt output for FIFO1,An Interrupt output for FIFO2 */
  RESERVED_Handler,                              /*  88 Reserved */
  RESERVED_Handler,                              /*  89 Reserved */
  RESERVED_Handler,                              /*  90 Reserved */
  RESERVED_Handler,                              /*  91 Reserved */
  LCU0_Handler,                                  /*  92 Interrupt 0, Interrupt 1 Interrupt 2 */
  LCU1_Handler,                                  /*  93 Interrupt 0, Interrupt 1 Interrupt 2 */
  RESERVED_Handler,                              /*  94 Reserved */
  RESERVED_Handler,                              /*  95 Reserved */
  PIT0_Handler,                                  /*  96 Interrupt for Channel0,Interrupt for Channel1,Interrupt for Channel2,Interrupt for Channel3,Interrupt for Channel4 */
  PIT1_Handler,                                  /*  97 Interrupt for Channel0,Interrupt for Channel1,Interrupt for Channel2,Interrupt for Channel3 */
  PIT2_Handler,                                  /*  98 Interrupt for Channel0,Interrupt for Channel1,Interrupt for Channel2,Interrupt for Channel3 */
  RESERVED_Handler,                              /*  99 Reserved */
  RESERVED_Handler,                              /* 100 Reserved */
  RESERVED_Handler,                              /* 101 Reserved */
  RTC_Handler,                                   /* 102 RTCF or ROVRF interrupt to be serviced by the system controller, APIF interrupt to be serviced by the system controller */
  RESERVED_Handler,                              /* 103 Reserved */
  RESERVED_Handler,                              /* 104 Reserved */
  EMAC_0_Handler,                                /* 105 Common interrupt */
  EMAC_1_Handler,                                /* 106 Tx interrupt 0, Tx interrupt 1 */
  EMAC_2_Handler,                                /* 107 Rx interrupt 0, Rx interrupt 1 */
  EMAC_3_Handler,                                /* 108 Safety interrupt correctable, Safety interrupt un-correctable */
  FlexCAN0_0_Handler,                            /* 109 Interrupt indicating that the CAN bus went to Bus Off state */
  FlexCAN0_1_Handler,                            /* 110 Message Buffer Interrupt line 0-31,ORed Interrupt for Message Buffers */
  FlexCAN0_2_Handler,                            /* 111 Message Buffer Interrupt line 32-63 */
  FlexCAN0_3_Handler,                            /* 112 Message Buffer Interrupt line 64-95 */
  FlexCAN1_0_Handler,                            /* 113 Interrupt indicating that the CAN bus went to Bus Off state */
  FlexCAN1_1_Handler,                            /* 114 Message Buffer Interrupt line 0-31 */
  FlexCAN1_2_Handler,                            /* 115 Message Buffer Interrupt line 32-63 */
  FlexCAN2_0_Handler,                            /* 116 Interrupt indicating that the CAN bus went to Bus Off state */
  FlexCAN2_1_Handler,                            /* 117 Message Buffer Interrupt line 0-31 */
  FlexCAN2_2_Handler,                            /* 118 Message Buffer Interrupt line 32-63 */
  FlexCAN3_0_Handler,                            /* 119 Interrupt indicating that the CAN bus went to Bus Off state */
  FlexCAN3_1_Handler,                            /* 120 Message Buffer Interrupt line 0-31 */
  FlexCAN4_0_Handler,                            /* 121 Interrupt indicating that the CAN bus went to Bus Off state */
  FlexCAN4_1_Handler,                            /* 122 Message Buffer Interrupt line 0-31 */
  FlexCAN5_0_Handler,                            /* 123 Interrupt indicating that the CAN bus went to Bus Off state */
  FlexCAN5_1_Handler,                            /* 124 Message Buffer Interrupt line 0-31 */
  RESERVED_Handler,                              /* 125 Reserved */
  RESERVED_Handler,                              /* 126 Reserved */
  RESERVED_Handler,                              /* 127 Reserved */
  RESERVED_Handler,                              /* 128 Reserved */
  RESERVED_Handler,                              /* 129 Reserved */
  RESERVED_Handler,                              /* 130 Reserved */
  RESERVED_Handler,                              /* 131 Reserved */
  RESERVED_Handler,                              /* 132 Reserved */
  RESERVED_Handler,                              /* 133 Reserved */
  RESERVED_Handler,                              /* 134 Reserved */
  RESERVED_Handler,                              /* 135 Reserved */
  RESERVED_Handler,                              /* 136 Reserved */
  RESERVED_Handler,                              /* 137 Reserved */
  RESERVED_Handler,                              /* 138 Reserved */
  FLEXIO_Handler,                                /* 139 FlexIO Interrupt */
  RESERVED_Handler,                              /* 140 Reserved */
  LPUART0_Handler,                               /* 141 TX and RX interrupt */
  LPUART1_Handler,                               /* 142 TX and RX interrupt */
  LPUART2_Handler,                               /* 143 TX and RX interrupt */
  LPUART3_Handler,                               /* 144 TX and RX interrupt */
  LPUART4_Handler,                               /* 145 TX and RX interrupt */
  LPUART5_Handler,                               /* 146 TX and RX interrupt */
  LPUART6_Handler,                               /* 147 TX and RX interrupt */
  LPUART7_Handler,                               /* 148 TX and RX interrupt */
  LPUART8_Handler,                               /* 149 TX and RX interrupt */
  LPUART9_Handler,                               /* 150 TX and RX interrupt */
  LPUART10_Handler,                              /* 151 TX and RX interrupt */
  LPUART11_Handler,                              /* 152 TX and RX interrupt */
  LPUART12_Handler,                              /* 153 TX and RX interrupt */
  LPUART13_Handler,                              /* 154 TX and RX interrupt */
  LPUART14_Handler,                              /* 155 TX and RX interrupt */
  LPUART15_Handler,                              /* 156 TX and RX interrupt */
  RESERVED_Handler,                              /* 157 Reserved */
  RESERVED_Handler,                              /* 158 Reserved */
  RESERVED_Handler,                              /* 159 Reserved */
  RESERVED_Handler,                              /* 160 Reserved */
  LPI2C0_Handler,                                /* 161 LPI2C Master Interrupt, LPI2C Interrupt */
  LPI2C1_Handler,                                /* 162 LPI2C Master Interrupt, LPI2C Interrupt */
  RESERVED_Handler,                              /* 163 Reserved */
  RESERVED_Handler,                              /* 164 Reserved */
  LPSPI0_Handler,                                /* 165 LPSPI Interrupt */
  LPSPI1_Handler,                                /* 166 LPSPI Interrupt */
  LPSPI2_Handler,                                /* 167 LPSPI Interrupt */
  LPSPI3_Handler,                                /* 168 LPSPI Interrupt */
  LPSPI4_Handler,                                /* 169 LPSPI Interrupt */
  LPSPI5_Handler,                                /* 170 LPSPI Interrupt */
  RESERVED_Handler,                              /* 171 Reserved */
  RESERVED_Handler,                              /* 172 Reserved */
  QSPI_Handler,                                  /* 173 TX Buffer Fill interrupt, Transfer Complete / Transaction Finished, RX Buffer Drain interrupt, Buffer Overflow / Underrun interrupt, Serial Flash Communication Error interrupt, All interrupts ORed output */
  SAI0_Handler,                                  /* 174 RX interrupt,TX interrupt */
  SAI1_Handler,                                  /* 175 RX interrupt,TX interrupt */
  RESERVED_Handler,                              /* 176 Reserved */
  RESERVED_Handler,                              /* 177 Reserved */
  JDC_Handler,                                   /* 178 Indicates new data to be read from JIN_IPS register or can be written to JOUT_IPS register */
  RESERVED_Handler,                              /* 179 Reserved */
  ADC0_Handler,                                  /* 180 End of conversion, Error interrupt, Watchdog interrupt */
  ADC1_Handler,                                  /* 181 End of conversion, Error interrupt, Watchdog interrupt */
  ADC2_Handler,                                  /* 182 End of conversion, Error interrupt, Watchdog interrupt */
  LPCMP0_Handler,                                /* 183 Async interrupt */
  LPCMP1_Handler,                                /* 184 Async interrupt */
  LPCMP2_Handler,                                /* 185 Async interrupt */
  RESERVED_Handler,                              /* 186 Reserved */
  RESERVED_Handler,                              /* 187 Reserved */
  RESERVED_Handler,                              /* 188 Reserved */
  FCCU_0_Handler,                                /* 189 Interrupt request(ALARM state) */
  FCCU_1_Handler,                                /* 190 Interrupt request(miscellaneous conditions) */
  STCU_LBIST_MBIST_Handler,                      /* 191 Interrupt request(miscellaneous conditions) */
  HSE_MU0_TX_Handler,                            /* 192 ORed TX interrupt to MU-0 */
  HSE_MU0_RX_Handler,                            /* 193 ORed RX interrupt to MU-0 */
  HSE_MU0_ORED_Handler,                          /* 194 ORed general purpose interrupt request to MU-0 */
  HSE_MU1_TX_Handler,                            /* 195 ORed TX interrupt to MU-1 */
  HSE_MU1_RX_Handler,                            /* 196 ORed RX interrupt to MU-1 */
  HSE_MU1_ORED_Handler,                          /* 197 ORed general purpose interrupt request to MU-1 */
  RESERVED_Handler,                              /* 198 Reserved */
  RESERVED_Handler,                              /* 199 Reserved */
  RESERVED_Handler,                              /* 200 Reserved */
  RESERVED_Handler,                              /* 201 Reserved */
  RESERVED_Handler,                              /* 202 Reserved */
  RESERVED_Handler,                              /* 203 Reserved */
  RESERVED_Handler,                              /* 204 Reserved */
  RESERVED_Handler,                              /* 205 Reserved */
  RESERVED_Handler,                              /* 206 Reserved */
  RESERVED_Handler,                              /* 207 Reserved */
  RESERVED_Handler,                              /* 208 Reserved */
  RESERVED_Handler,                              /* 209 Reserved */
  RESERVED_Handler,                              /* 210 Reserved */
  RESERVED_Handler,                              /* 211 Reserved */
  SoC_PLL_Handler                                /* 212 PLL LOL interrupt */
};

#if defined ( __GNUC__ )
#pragma GCC diagnostic pop
#endif

/*----------------------------------------------------------------------------
  Reset Handler called on controller reset
 *----------------------------------------------------------------------------*/

__attribute__((naked))
void Reset_Handler(void)
{
  #define xstr(s) str(s)
  #define str(s)  #s

  __ASM volatile (
#ifndef __ICCARM__
    ".syntax unified\n"
#endif

    // Initialize SRAM
    "        LDR      R2,=0\n"
    "        LDR      R3,=0\n"
    "        LDR      R4,=" xstr(RAM_START__) "\n"
    "        LDR      R5,=" xstr(RAM_SIZE__) "\n"

    "        LSRS     R5,R5,#3\n"
    "SRAM_Loop:\n"
    "        STRD     r2,r3,[R4,#0]\n"
    "        ADDS     R4,R4,#0x08\n"
    "        SUBS     R5,R5,#1\n"
    "        CMP      R5,#0x00\n"
    "        BNE      SRAM_Loop\n"

    // Initialize DTCM
    "        LDR      R2,=0\n"
    "        LDR      R3,=0\n"
    "        LDR      R4,=" xstr(DTCM_START__) "\n"
    "        LDR      R5,=" xstr(DTCM_SIZE__) "\n"

    "        LSRS     R5,R5,#3\n"
    "DTCM_Loop:\n"
    "        STRD     r2,r3,[R4,#0]\n"
    "        ADDS     R4,R4,#0x08\n"
    "        SUBS     R5,R5,#1\n"
    "        CMP      R5,#0x00\n"
    "        BNE      DTCM_Loop\n"

    "        B        Reset_Handler_C\n"
  );

  #undef str
  #undef xstr
}

void Reset_Handler_C(void)
{
  SystemInit();
  __PROGRAM_START();
}


#if defined(__ARMCC_VERSION) && (__ARMCC_VERSION >= 6010050)
  #pragma clang diagnostic push
  #pragma clang diagnostic ignored "-Wmissing-noreturn"
#endif

/*----------------------------------------------------------------------------
  Hard Fault Handler
 *----------------------------------------------------------------------------*/
void HardFault_Handler(void)
{
  while(1);
}

/*----------------------------------------------------------------------------
  Default Handler for Exceptions / Interrupts
 *----------------------------------------------------------------------------*/
void Default_Handler(void)
{
  while(1);
}

/*----------------------------------------------------------------------------
  Reserved Handler for unimplemented Exceptions / Interrupts
 *----------------------------------------------------------------------------*/
void Reserved_Handler(void) {
  while(1);
}

#if defined(__ARMCC_VERSION) && (__ARMCC_VERSION >= 6010050)
  #pragma clang diagnostic pop
#endif

