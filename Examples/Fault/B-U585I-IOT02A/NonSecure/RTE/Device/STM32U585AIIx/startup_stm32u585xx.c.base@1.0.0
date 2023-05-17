/******************************************************************************
 * @file     startup_stm32u585xx.c
 * @brief    CMSIS-Core Device Startup File for STM32U585xx Device
 * @version  V1.0.0
 * @date     23. February 2023
 ******************************************************************************/
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

#include "stm32u5xx.h"

/*----------------------------------------------------------------------------
  Exception / Interrupt Handler Function Prototype
 *----------------------------------------------------------------------------*/
typedef void(*VECTOR_TABLE_Type)(void);

/*----------------------------------------------------------------------------
  External References
 *----------------------------------------------------------------------------*/
extern uint32_t __INITIAL_SP;
extern uint32_t __STACK_LIMIT;
#if defined (__ARM_FEATURE_CMSE) && (__ARM_FEATURE_CMSE == 3U)
extern uint32_t __STACK_SEAL;
#endif

extern __NO_RETURN void __PROGRAM_START(void);

/*----------------------------------------------------------------------------
  Internal References
 *----------------------------------------------------------------------------*/
__NO_RETURN void Reset_Handler  (void);
            void Default_Handler(void);

/*----------------------------------------------------------------------------
  Exception / Interrupt Handler
 *----------------------------------------------------------------------------*/
/* Exceptions */
void NMI_Handler                  (void) __attribute__ ((weak, alias("Default_Handler")));
void HardFault_Handler            (void) __attribute__ ((weak));
void MemManage_Handler            (void) __attribute__ ((weak, alias("Default_Handler")));
void BusFault_Handler             (void) __attribute__ ((weak, alias("Default_Handler")));
void UsageFault_Handler           (void) __attribute__ ((weak, alias("Default_Handler")));
void SecureFault_Handler          (void) __attribute__ ((weak, alias("Default_Handler")));
void SVC_Handler                  (void) __attribute__ ((weak, alias("Default_Handler")));
void DebugMon_Handler             (void) __attribute__ ((weak, alias("Default_Handler")));
void PendSV_Handler               (void) __attribute__ ((weak, alias("Default_Handler")));
void SysTick_Handler              (void) __attribute__ ((weak, alias("Default_Handler")));

void WWDG_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void PVD_PVM_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void RTC_IRQHandler               (void) __attribute__ ((weak, alias("Default_Handler")));
void RTC_S_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void TAMP_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void RAMCFG_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void FLASH_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void FLASH_S_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void GTZC_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void RCC_IRQHandler               (void) __attribute__ ((weak, alias("Default_Handler")));
void RCC_S_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI0_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI1_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI2_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI3_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI4_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI5_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI6_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI7_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI8_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI9_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI10_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI11_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI12_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI13_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI14_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void EXTI15_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void IWDG_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void SAES_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel0_IRQHandler   (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel1_IRQHandler   (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel2_IRQHandler   (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel3_IRQHandler   (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel4_IRQHandler   (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel5_IRQHandler   (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel6_IRQHandler   (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel7_IRQHandler   (void) __attribute__ ((weak, alias("Default_Handler")));
void ADC1_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void DAC1_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void FDCAN1_IT0_IRQHandler        (void) __attribute__ ((weak, alias("Default_Handler")));
void FDCAN1_IT1_IRQHandler        (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM1_BRK_IRQHandler          (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM1_UP_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM1_TRG_COM_IRQHandler      (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM1_CC_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM2_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM3_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM4_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM5_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM6_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM7_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM8_BRK_IRQHandler          (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM8_UP_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM8_TRG_COM_IRQHandler      (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM8_CC_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void I2C1_EV_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void I2C1_ER_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void I2C2_EV_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void I2C2_ER_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void SPI1_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void SPI2_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void USART1_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void USART2_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void USART3_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void UART4_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void UART5_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void LPUART1_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void LPTIM1_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void LPTIM2_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM15_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM16_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void TIM17_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void COMP_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void OTG_FS_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void CRS_IRQHandler               (void) __attribute__ ((weak, alias("Default_Handler")));
void FMC_IRQHandler               (void) __attribute__ ((weak, alias("Default_Handler")));
void OCTOSPI1_IRQHandler          (void) __attribute__ ((weak, alias("Default_Handler")));
void PWR_S3WU_IRQHandler          (void) __attribute__ ((weak, alias("Default_Handler")));
void SDMMC1_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void SDMMC2_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel8_IRQHandler   (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel9_IRQHandler   (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel10_IRQHandler  (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel11_IRQHandler  (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel12_IRQHandler  (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel13_IRQHandler  (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel14_IRQHandler  (void) __attribute__ ((weak, alias("Default_Handler")));
void GPDMA1_Channel15_IRQHandler  (void) __attribute__ ((weak, alias("Default_Handler")));
void I2C3_EV_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void I2C3_ER_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void SAI1_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void SAI2_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void TSC_IRQHandler               (void) __attribute__ ((weak, alias("Default_Handler")));
void AES_IRQHandler               (void) __attribute__ ((weak, alias("Default_Handler")));
void RNG_IRQHandler               (void) __attribute__ ((weak, alias("Default_Handler")));
void FPU_IRQHandler               (void) __attribute__ ((weak, alias("Default_Handler")));
void HASH_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void PKA_IRQHandler               (void) __attribute__ ((weak, alias("Default_Handler")));
void LPTIM3_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void SPI3_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void I2C4_ER_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void I2C4_EV_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void MDF1_FLT0_IRQHandler         (void) __attribute__ ((weak, alias("Default_Handler")));
void MDF1_FLT1_IRQHandler         (void) __attribute__ ((weak, alias("Default_Handler")));
void MDF1_FLT2_IRQHandler         (void) __attribute__ ((weak, alias("Default_Handler")));
void MDF1_FLT3_IRQHandler         (void) __attribute__ ((weak, alias("Default_Handler")));
void UCPD1_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void ICACHE_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void OTFDEC1_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void OTFDEC2_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void LPTIM4_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void DCACHE1_IRQHandler           (void) __attribute__ ((weak, alias("Default_Handler")));
void ADF1_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void ADC4_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));
void LPDMA1_Channel0_IRQHandler   (void) __attribute__ ((weak, alias("Default_Handler")));
void LPDMA1_Channel1_IRQHandler   (void) __attribute__ ((weak, alias("Default_Handler")));
void LPDMA1_Channel2_IRQHandler   (void) __attribute__ ((weak, alias("Default_Handler")));
void LPDMA1_Channel3_IRQHandler   (void) __attribute__ ((weak, alias("Default_Handler")));
void DMA2D_IRQHandler             (void) __attribute__ ((weak, alias("Default_Handler")));
void DCMI_PSSI_IRQHandler         (void) __attribute__ ((weak, alias("Default_Handler")));
void OCTOSPI2_IRQHandler          (void) __attribute__ ((weak, alias("Default_Handler")));
void MDF1_FLT4_IRQHandler         (void) __attribute__ ((weak, alias("Default_Handler")));
void MDF1_FLT5_IRQHandler         (void) __attribute__ ((weak, alias("Default_Handler")));
void CORDIC_IRQHandler            (void) __attribute__ ((weak, alias("Default_Handler")));
void FMAC_IRQHandler              (void) __attribute__ ((weak, alias("Default_Handler")));


/*----------------------------------------------------------------------------
  Exception / Interrupt Vector table
 *----------------------------------------------------------------------------*/

#if defined ( __GNUC__ )
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wpedantic"
#endif

extern const VECTOR_TABLE_Type __VECTOR_TABLE[];
       const VECTOR_TABLE_Type __VECTOR_TABLE[] __VECTOR_TABLE_ATTRIBUTE = {
  (VECTOR_TABLE_Type)(&__INITIAL_SP),       /*     Initial Stack Pointer */
  Reset_Handler,                            /*     Reset Handler */
  NMI_Handler,                              /* -14 NMI Handler */
  HardFault_Handler,                        /* -13 Hard Fault Handler */
  MemManage_Handler,                        /* -12 MPU Fault Handler */
  BusFault_Handler,                         /* -11 Bus Fault Handler */
  UsageFault_Handler,                       /* -10 Usage Fault Handler */
  SecureFault_Handler,                      /*  -9 Secure Fault Handler */
  0,                                        /*     Reserved */
  0,                                        /*     Reserved */
  0,                                        /*     Reserved */
  SVC_Handler,                              /*  -5 SVCall Handler */
  DebugMon_Handler,                         /*  -4 Debug Monitor Handler */
  0,                                        /*     Reserved */
  PendSV_Handler,                           /*  -2 PendSV Handler */
  SysTick_Handler,                          /*  -1 SysTick Handler */

  /* Interrupts */
  WWDG_IRQHandler,                          /*     Window WatchDog */
  PVD_PVM_IRQHandler,                       /*     PVD/PVM through EXTI Line detection interrupt */
  RTC_IRQHandler,                           /*     RTC non-secure interrupt */
  RTC_S_IRQHandler,                         /*     RTC secure interrupt */
  TAMP_IRQHandler,                          /*     Tamper non-secure interrupt */
  RAMCFG_IRQHandler,                        /*     RAMCFG global interrupt */
  FLASH_IRQHandler,                         /*     FLASH non-secure global interrupt */
  FLASH_S_IRQHandler,                       /*     FLASH secure global interrupt */
  GTZC_IRQHandler,                          /*     Global TrustZone Controller interrupt */
  RCC_IRQHandler,                           /*     RCC non-secure global interrupt */
  RCC_S_IRQHandler,                         /*     RCC secure global interrupt */
  EXTI0_IRQHandler,                         /*     EXTI Line0 interrupt */
  EXTI1_IRQHandler,                         /*     EXTI Line1 interrupt */
  EXTI2_IRQHandler,                         /*     EXTI Line2 interrupt */
  EXTI3_IRQHandler,                         /*     EXTI Line3 interrupt */
  EXTI4_IRQHandler,                         /*     EXTI Line4 interrupt */
  EXTI5_IRQHandler,                         /*     EXTI Line5 interrupt */
  EXTI6_IRQHandler,                         /*     EXTI Line6 interrupt */
  EXTI7_IRQHandler,                         /*     EXTI Line7 interrupt */
  EXTI8_IRQHandler,                         /*     EXTI Line8 interrupt */
  EXTI9_IRQHandler,                         /*     EXTI Line9 interrupt */
  EXTI10_IRQHandler,                        /*     EXTI Line10 interrupt */
  EXTI11_IRQHandler,                        /*     EXTI Line11 interrupt */
  EXTI12_IRQHandler,                        /*     EXTI Line12 interrupt */
  EXTI13_IRQHandler,                        /*     EXTI Line13 interrupt */
  EXTI14_IRQHandler,                        /*     EXTI Line14 interrupt */
  EXTI15_IRQHandler,                        /*     EXTI Line15 interrupt */
  IWDG_IRQHandler,                          /*     IWDG global interrupt */
  SAES_IRQHandler,                          /*     Secure AES global interrupt */
  GPDMA1_Channel0_IRQHandler,               /*     GPDMA1 Channel 0 global interrupt */
  GPDMA1_Channel1_IRQHandler,               /*     GPDMA1 Channel 1 global interrupt */
  GPDMA1_Channel2_IRQHandler,               /*     GPDMA1 Channel 2 global interrupt */
  GPDMA1_Channel3_IRQHandler,               /*     GPDMA1 Channel 3 global interrupt */
  GPDMA1_Channel4_IRQHandler,               /*     GPDMA1 Channel 4 global interrupt */
  GPDMA1_Channel5_IRQHandler,               /*     GPDMA1 Channel 5 global interrupt */
  GPDMA1_Channel6_IRQHandler,               /*     GPDMA1 Channel 6 global interrupt */
  GPDMA1_Channel7_IRQHandler,               /*     GPDMA1 Channel 7 global interrupt */
  ADC1_IRQHandler,                          /*     ADC1 global interrupt */
  DAC1_IRQHandler,                          /*     DAC1 global interrupt */
  FDCAN1_IT0_IRQHandler,                    /*     FDCAN1 interrupt 0 */
  FDCAN1_IT1_IRQHandler,                    /*     FDCAN1 interrupt 1 */
  TIM1_BRK_IRQHandler,                      /*     TIM1 Break interrupt */
  TIM1_UP_IRQHandler,                       /*     TIM1 Update interrupt */
  TIM1_TRG_COM_IRQHandler,                  /*     TIM1 Trigger and Commutation interrupt */
  TIM1_CC_IRQHandler,                       /*     TIM1 Capture Compare interrupt */
  TIM2_IRQHandler,                          /*     TIM2 global interrupt */
  TIM3_IRQHandler,                          /*     TIM3 global interrupt */
  TIM4_IRQHandler,                          /*     TIM4 global interrupt */
  TIM5_IRQHandler,                          /*     TIM5 global interrupt */
  TIM6_IRQHandler,                          /*     TIM6 global interrupt */
  TIM7_IRQHandler,                          /*     TIM7 global interrupt */
  TIM8_BRK_IRQHandler,                      /*     TIM8 Break interrupt */
  TIM8_UP_IRQHandler,                       /*     TIM8 Update interrupt */
  TIM8_TRG_COM_IRQHandler,                  /*     TIM8 Trigger and Commutation interrupt */
  TIM8_CC_IRQHandler,                       /*     TIM8 Capture Compare interrupt */
  I2C1_EV_IRQHandler,                       /*     I2C1 Event interrupt */
  I2C1_ER_IRQHandler,                       /*     I2C1 Error interrupt */
  I2C2_EV_IRQHandler,                       /*     I2C2 Event interrupt */
  I2C2_ER_IRQHandler,                       /*     I2C2 Error interrupt */
  SPI1_IRQHandler,                          /*     SPI1 global interrupt */
  SPI2_IRQHandler,                          /*     SPI2 global interrupt */
  USART1_IRQHandler,                        /*     USART1 global interrupt */
  USART2_IRQHandler,                        /*     USART2 global interrupt */
  USART3_IRQHandler,                        /*     USART3 global interrupt */
  UART4_IRQHandler,                         /*     UART4 global interrupt */
  UART5_IRQHandler,                         /*     UART5 global interrupt */
  LPUART1_IRQHandler,                       /*     LPUART1 global interrupt */
  LPTIM1_IRQHandler,                        /*     LPTIM1 global interrupt */
  LPTIM2_IRQHandler,                        /*     LPTIM2 global interrupt */
  TIM15_IRQHandler,                         /*     TIM15 global interrupt */
  TIM16_IRQHandler,                         /*     TIM16 global interrupt */
  TIM17_IRQHandler,                         /*     TIM17 global interrupt */
  COMP_IRQHandler,                          /*     COMP1 and COMP2 through EXTI Lines interrupt */
  OTG_FS_IRQHandler,                        /*     USB OTG FS global interrupt */
  CRS_IRQHandler,                           /*     CRS global interrupt */
  FMC_IRQHandler,                           /*     FMC global interrupt */
  OCTOSPI1_IRQHandler,                      /*     OctoSPI1 global interrupt */
  PWR_S3WU_IRQHandler,                      /*     PWR wake up from Stop3 interrupt */
  SDMMC1_IRQHandler,                        /*     SDMMC1 global interrupt */
  SDMMC2_IRQHandler,                        /*     SDMMC2 global interrupt */
  GPDMA1_Channel8_IRQHandler,               /*     GPDMA1 Channel 8 global interrupt */
  GPDMA1_Channel9_IRQHandler,               /*     GPDMA1 Channel 9 global interrupt */
  GPDMA1_Channel10_IRQHandler,              /*     GPDMA1 Channel 10 global interrupt */
  GPDMA1_Channel11_IRQHandler,              /*     GPDMA1 Channel 11 global interrupt */
  GPDMA1_Channel12_IRQHandler,              /*     GPDMA1 Channel 12 global interrupt */
  GPDMA1_Channel13_IRQHandler,              /*     GPDMA1 Channel 13 global interrupt */
  GPDMA1_Channel14_IRQHandler,              /*     GPDMA1 Channel 14 global interrupt */
  GPDMA1_Channel15_IRQHandler,              /*     GPDMA1 Channel 15 global interrupt */
  I2C3_EV_IRQHandler,                       /*     I2C3 event interrupt */
  I2C3_ER_IRQHandler,                       /*     I2C3 error interrupt */
  SAI1_IRQHandler,                          /*     Serial Audio Interface 1 global interrupt */
  SAI2_IRQHandler,                          /*     Serial Audio Interface 2 global interrupt */
  TSC_IRQHandler,                           /*     Touch Sense Controller global interrupt */
  AES_IRQHandler,                           /*     AES global interrupt */
  RNG_IRQHandler,                           /*     RNG global interrupt */
  FPU_IRQHandler,                           /*     FPU global interrupt */
  HASH_IRQHandler,                          /*     HASH global interrupt */
  PKA_IRQHandler,                           /*     PKA global interrupt */
  LPTIM3_IRQHandler,                        /*     LPTIM3 global interrupt */
  SPI3_IRQHandler,                          /*     SPI3 global interrupt */
  I2C4_ER_IRQHandler,                       /*     I2C4 Error interrupt */
  I2C4_EV_IRQHandler,                       /*     I2C4 Event interrupt */
  MDF1_FLT0_IRQHandler,                     /*     MDF1 Filter 0 global interrupt */
  MDF1_FLT1_IRQHandler,                     /*     MDF1 Filter 1 global interrupt */
  MDF1_FLT2_IRQHandler,                     /*     MDF1 Filter 2 global interrupt */
  MDF1_FLT3_IRQHandler,                     /*     MDF1 Filter 3 global interrupt */
  UCPD1_IRQHandler,                         /*     UCPD1 global interrupt */
  ICACHE_IRQHandler,                        /*     Instruction cache global interrupt */
  OTFDEC1_IRQHandler,                       /*     OTFDEC1 global interrupt */
  OTFDEC2_IRQHandler,                       /*     OTFDEC2 global interrupt */
  LPTIM4_IRQHandler,                        /*     LPTIM4 global interrupt */
  DCACHE1_IRQHandler,                       /*     Data cache global interrupt */
  ADF1_IRQHandler,                          /*     ADF interrupt */
  ADC4_IRQHandler,                          /*     ADC4 (12bits) global interrupt */
  LPDMA1_Channel0_IRQHandler,               /*     LPDMA1 SmartRun Channel 0 global interrupt */
  LPDMA1_Channel1_IRQHandler,               /*     LPDMA1 SmartRun Channel 1 global interrupt */
  LPDMA1_Channel2_IRQHandler,               /*     LPDMA1 SmartRun Channel 2 global interrupt */
  LPDMA1_Channel3_IRQHandler,               /*     LPDMA1 SmartRun Channel 3 global interrupt */
  DMA2D_IRQHandler,                         /*     DMA2D global interrupt */
  DCMI_PSSI_IRQHandler,                     /*     DCMI/PSSI global interrupt */
  OCTOSPI2_IRQHandler,                      /*     OCTOSPI2 global interrupt */
  MDF1_FLT4_IRQHandler,                     /*     MDF1 Filter 4 global interrupt */
  MDF1_FLT5_IRQHandler,                     /*     MDF1 Filter 5 global interrupt */
  CORDIC_IRQHandler,                        /*     CORDIC global interrupt */
  FMAC_IRQHandler                           /*     FMAC global interrupt */
};

#if defined ( __GNUC__ )
#pragma GCC diagnostic pop
#endif

/*----------------------------------------------------------------------------
  Reset Handler called on controller reset
 *----------------------------------------------------------------------------*/
__NO_RETURN void Reset_Handler(void)
{
  __set_PSP((uint32_t)(&__INITIAL_SP));

  __set_MSPLIM((uint32_t)(&__STACK_LIMIT));
  __set_PSPLIM((uint32_t)(&__STACK_LIMIT));

#if defined (__ARM_FEATURE_CMSE) && (__ARM_FEATURE_CMSE == 3U)
  __TZ_set_STACKSEAL_S((uint32_t *)(&__STACK_SEAL));
#endif

  SystemInit();                             /* CMSIS System Initialization */
  __PROGRAM_START();                        /* Enter PreMain (C library entry point) */
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

#if defined(__ARMCC_VERSION) && (__ARMCC_VERSION >= 6010050)
  #pragma clang diagnostic pop
#endif

