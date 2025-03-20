/* USER CODE BEGIN Header */
/**
  ******************************************************************************
  * @file           : main.c
  * @brief          : Main program body
  ******************************************************************************
  * @attention
  *
  * Copyright (c) 2024 STMicroelectronics.
  * All rights reserved.
  *
  * This software is licensed under terms that can be found in the LICENSE file
  * in the root directory of this software component.
  * If no LICENSE file comes with this software, it is provided AS-IS.
  *
  ******************************************************************************
  */
/* USER CODE END Header */
/* Includes ------------------------------------------------------------------*/
#include "main.h"

/* Private includes ----------------------------------------------------------*/
/* USER CODE BEGIN Includes */
#include <stdio.h>

#include "RTE_Components.h"
#include  CMSIS_device_header

#include "cmsis_os2.h"
#include "cmsis_vio.h"
#include "EventRecorder.h"

#include "ARM_Fault.h"

/* USER CODE END Includes */

/* Private typedef -----------------------------------------------------------*/
/* USER CODE BEGIN PTD */

/* USER CODE END PTD */

/* Private define ------------------------------------------------------------*/
/* USER CODE BEGIN PD */

/* USER CODE END PD */

/* Private macro -------------------------------------------------------------*/
/* USER CODE BEGIN PM */

/* USER CODE END PM */

/* Private variables ---------------------------------------------------------*/

UART_HandleTypeDef huart1;

/* USER CODE BEGIN PV */

/* USER CODE END PV */

/* Private function prototypes -----------------------------------------------*/
static void MX_GPIO_Init(void);
static void MX_GTZC_NS_Init(void);
static void MX_USART1_UART_Init(void);
/* USER CODE BEGIN PFP */

/* USER CODE END PFP */

/* Private user code ---------------------------------------------------------*/
/* USER CODE BEGIN 0 */

/* USER CODE END 0 */

/**
  * @brief  The application entry point.
  * @retval int
  */
int main(void)
{

  /* USER CODE BEGIN 1 */

  /* USER CODE END 1 */

  /* MCU Configuration--------------------------------------------------------*/

  /* Reset of all peripherals, Initializes the Flash interface and the Systick. */
  HAL_Init();

  /* USER CODE BEGIN Init */

/**
  * CubeMX cannot generate required MPU configuration!
  *
  * Configure MPU
  *   - region 0: ROM                   - 0x08100000 .. 0x081FFFFF
  *   - region 1: RAM (privileged only) - 0x20040000 .. 0x200400FF
  *   - region 2: RAM                   - 0x20040100 .. 0x200CFFFF
  *   - region 3: Peripherals           - 0x40000000 .. 0x4FFFFFFF
  */
  ARM_MPU_Disable();

  // Memory attributes for Flash (index 0) = Outer/Inner: Normal cacheable memory, Non-Transient, no Write-Back, Read Allocate, no Write Allocate
  ARM_MPU_SetMemAttr(0UL, ARM_MPU_ATTR(ARM_MPU_ATTR_MEMORY_(1UL, 0UL, 1UL, 0UL), ARM_MPU_ATTR_MEMORY_(1UL, 0UL, 1UL, 0UL)));

  // Memory attributes for RAM (index 1) = Outer/Inner: Normal non-cacheable memory
  ARM_MPU_SetMemAttr(1UL, ARM_MPU_ATTR(ARM_MPU_ATTR_NON_CACHEABLE, ARM_MPU_ATTR_NON_CACHEABLE));

  // Memory attributes for Peripherals (index 2) = Device memory: nG (non-Gathering), nR (non-Reordering), nE (no Early Write Acknowledgment)
  ARM_MPU_SetMemAttr(2UL, ARM_MPU_ATTR(ARM_MPU_ATTR_DEVICE, ARM_MPU_ATTR_DEVICE_nGnRnE));

  /* Configure regions
                region,             (BASE       , Shareability  , RO , NP , XN ),             (LIMIT     , ATTR IDX) */
  ARM_MPU_SetRegion(0U, ARM_MPU_RBAR(0x08100000 , ARM_MPU_SH_NON, 1UL, 1UL, 0UL), ARM_MPU_RLAR(0x081FFFFF, 0UL));
  ARM_MPU_SetRegion(1U, ARM_MPU_RBAR(0x20040000 , ARM_MPU_SH_NON, 0UL, 0UL, 1UL), ARM_MPU_RLAR(0x200400FF, 1UL));
  ARM_MPU_SetRegion(2U, ARM_MPU_RBAR(0x20040100 , ARM_MPU_SH_NON, 0UL, 1UL, 1UL), ARM_MPU_RLAR(0x200CFFFF, 1UL));
  ARM_MPU_SetRegion(3U, ARM_MPU_RBAR(0x40000000 , ARM_MPU_SH_NON, 0UL, 1UL, 1UL), ARM_MPU_RLAR(0x4FFFFFFF, 2UL));

  ARM_MPU_Enable(MPU_CTRL_PRIVDEFENA_Msk);      // Enable Privileged Default

  /* USER CODE END Init */

  /* GTZC initialisation */
  MX_GTZC_NS_Init();

  /* USER CODE BEGIN SysInit */

  SCB->SHCSR |= SCB_SHCSR_BUSFAULTENA_Msk |     // Enable BusFault
                SCB_SHCSR_USGFAULTENA_Msk;      // Enable UsageFault
  SCB->CCR   |= SCB_CCR_DIV_0_TRP_Msk;          // Enable divide by 0 trap
  /* USER CODE END SysInit */

  /* Initialize all configured peripherals */
  MX_GPIO_Init();
  MX_USART1_UART_Init();
  /* USER CODE BEGIN 2 */

  stdio_init();                                 // Initialize STDIO

  vioInit();                                    // Initialize Virtual I/O

  if (ARM_FaultOccurred() != 0U) {              // If fault information exists
    printf("\r\n\r\n- System Restarted -\r\n\r\n");
  }

  EventRecorderInitialize(EventRecordAll, 1U);  // Initialize and start Event Recorder

  if (ARM_FaultOccurred() != 0U) {              // If fault information exists
    ARM_FaultPrint();                           // Output decoded fault information via STDIO
    ARM_FaultRecord();                          // Output decoded fault information via Event Recorder
    EventRecorderStop();                        // Stop Event Recorder
  }

  app_main();                                   // Execute Application main

  /* USER CODE END 2 */

  /* Infinite loop */
  /* USER CODE BEGIN WHILE */
  while (1)
  {
    /* USER CODE END WHILE */

    /* USER CODE BEGIN 3 */
  }
  /* USER CODE END 3 */
}

/**
  * @brief GTZC_NS Initialization Function
  * @param None
  * @retval None
  */
static void MX_GTZC_NS_Init(void)
{

  /* USER CODE BEGIN GTZC_NS_Init 0 */

  /* USER CODE END GTZC_NS_Init 0 */

  /* USER CODE BEGIN GTZC_NS_Init 1 */

  /* USER CODE END GTZC_NS_Init 1 */
  /* USER CODE BEGIN GTZC_NS_Init 2 */

  /* USER CODE END GTZC_NS_Init 2 */

}

/**
  * @brief USART1 Initialization Function
  * @param None
  * @retval None
  */
static void MX_USART1_UART_Init(void)
{

  /* USER CODE BEGIN USART1_Init 0 */

  /* USER CODE END USART1_Init 0 */

  /* USER CODE BEGIN USART1_Init 1 */

  /* USER CODE END USART1_Init 1 */
  huart1.Instance = USART1;
  huart1.Init.BaudRate = 115200;
  huart1.Init.WordLength = UART_WORDLENGTH_8B;
  huart1.Init.StopBits = UART_STOPBITS_1;
  huart1.Init.Parity = UART_PARITY_NONE;
  huart1.Init.Mode = UART_MODE_TX_RX;
  huart1.Init.HwFlowCtl = UART_HWCONTROL_NONE;
  huart1.Init.OverSampling = UART_OVERSAMPLING_16;
  huart1.Init.OneBitSampling = UART_ONE_BIT_SAMPLE_DISABLE;
  huart1.Init.ClockPrescaler = UART_PRESCALER_DIV1;
  huart1.AdvancedInit.AdvFeatureInit = UART_ADVFEATURE_NO_INIT;
  if (HAL_UART_Init(&huart1) != HAL_OK)
  {
    Error_Handler();
  }
  if (HAL_UARTEx_SetTxFifoThreshold(&huart1, UART_TXFIFO_THRESHOLD_1_8) != HAL_OK)
  {
    Error_Handler();
  }
  if (HAL_UARTEx_SetRxFifoThreshold(&huart1, UART_RXFIFO_THRESHOLD_1_8) != HAL_OK)
  {
    Error_Handler();
  }
  if (HAL_UARTEx_DisableFifoMode(&huart1) != HAL_OK)
  {
    Error_Handler();
  }
  /* USER CODE BEGIN USART1_Init 2 */

  /* USER CODE END USART1_Init 2 */

}

/**
  * @brief GPIO Initialization Function
  * @param None
  * @retval None
  */
static void MX_GPIO_Init(void)
{
  GPIO_InitTypeDef GPIO_InitStruct = {0};
/* USER CODE BEGIN MX_GPIO_Init_1 */
/* USER CODE END MX_GPIO_Init_1 */

  /* GPIO Ports Clock Enable */
  __HAL_RCC_GPIOH_CLK_ENABLE();
  __HAL_RCC_GPIOA_CLK_ENABLE();

  /*Configure GPIO pin Output Level */
  HAL_GPIO_WritePin(GPIOH, LED_RED_Pin|LED_GREEN_Pin, GPIO_PIN_RESET);

  /*Configure GPIO pins : LED_RED_Pin LED_GREEN_Pin */
  GPIO_InitStruct.Pin = LED_RED_Pin|LED_GREEN_Pin;
  GPIO_InitStruct.Mode = GPIO_MODE_OUTPUT_PP;
  GPIO_InitStruct.Pull = GPIO_PULLUP;
  GPIO_InitStruct.Speed = GPIO_SPEED_FREQ_HIGH;
  HAL_GPIO_Init(GPIOH, &GPIO_InitStruct);

/* USER CODE BEGIN MX_GPIO_Init_2 */
/* USER CODE END MX_GPIO_Init_2 */
}

/* USER CODE BEGIN 4 */

/* USER CODE END 4 */

/**
  * @brief  Period elapsed callback in non blocking mode
  * @note   This function is called  when TIM16 interrupt took place, inside
  * HAL_TIM_IRQHandler(). It makes a direct call to HAL_IncTick() to increment
  * a global variable "uwTick" used as application time base.
  * @param  htim : TIM handle
  * @retval None
  */
void HAL_TIM_PeriodElapsedCallback(TIM_HandleTypeDef *htim)
{
  /* USER CODE BEGIN Callback 0 */

  /* USER CODE END Callback 0 */
  if (htim->Instance == TIM16) {
    HAL_IncTick();
  }
  /* USER CODE BEGIN Callback 1 */

  /* USER CODE END Callback 1 */
}

/**
  * @brief  This function is executed in case of error occurrence.
  * @retval None
  */
void Error_Handler(void)
{
  /* USER CODE BEGIN Error_Handler_Debug */
  /* User can add his own implementation to report the HAL error return state */
  __disable_irq();
  while (1)
  {
  }
  /* USER CODE END Error_Handler_Debug */
}

#ifdef  USE_FULL_ASSERT
/**
  * @brief  Reports the name of the source file and the source line number
  *         where the assert_param error has occurred.
  * @param  file: pointer to the source file name
  * @param  line: assert_param error line source number
  * @retval None
  */
void assert_failed(uint8_t *file, uint32_t line)
{
  /* USER CODE BEGIN 6 */
  /* User can add his own implementation to report the file name and line number,
     ex: printf("Wrong parameters value: file %s on line %d\r\n", file, line) */
  /* USER CODE END 6 */
}
#endif /* USE_FULL_ASSERT */
