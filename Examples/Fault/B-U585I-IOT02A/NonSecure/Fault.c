/*----------------------------------------------------------------------------
 * Name:    Fault.c
 * Purpose: Fault example program
 *----------------------------------------------------------------------------*/

#include <stdio.h>

#include "main.h"

#include "RTE_Components.h"
#include  CMSIS_device_header

#include "cmsis_os2.h"
#include "../Secure/interface.h"        // Secure functions interface header

#include "ARM_Fault.h"
#include "ARM_FaultTrigger.h"

#include "EventRecorder.h"

extern osThreadId_t tid_AppThread;
extern osThreadId_t tid_FaultTriggerThread;

/* Global Thread IDs (for debug) */
osThreadId_t tid_AppThread;
osThreadId_t tid_FaultTriggerThread;

/*---------------------------------------------------------------------------
 * Application thread
 *---------------------------------------------------------------------------*/
static __NO_RETURN void AppThread (void *argument) {

  (void)argument;

  for (;;) {
    osDelay(500U);
    HAL_GPIO_WritePin(LED_GREEN_GPIO_Port, LED_GREEN_Pin, GPIO_PIN_RESET);
    osDelay(500U);
    HAL_GPIO_WritePin(LED_GREEN_GPIO_Port, LED_GREEN_Pin, GPIO_PIN_SET);
  }
}

/*---------------------------------------------------------------------------
 * Fault trigger thread
 *---------------------------------------------------------------------------*/
static __NO_RETURN void FaultTriggerThread (void *argument) {
  char ch;

  (void)argument;

  // Display user interface message
  printf("\r\n--- Fault example ---\r\n\r\n");
  printf("To trigger a fault please input a corresponding number:\r\n");
  printf(" - 1: Non-Secure fault, Non-Secure data access Memory Management fault\r\n");
  printf(" - 2: Non-Secure fault, Non-Secure data access Bus fault\r\n");
  printf(" - 3: Non-Secure fault, Non-Secure undefined instruction Usage fault\r\n");
  printf(" - 4: Non-Secure fault, Non-Secure divide by 0 Usage fault\r\n");
  printf(" - 5: Secure fault, Non-Secure data access from Secure RAM memory\r\n");
  printf(" - 6: Secure fault, Non-Secure instruction execution from Secure Flash memory\r\n");
  printf(" - 7: Secure fault, Secure undefined instruction Usage fault\r\n");
  printf("Input>");

  for (;;) {
    ch = (char)getchar();                       // Read character from console (blocking)
    ARM_FaultTrigger((uint32_t)(ch - '0'));     // Trigger a fault
  }
}

/*---------------------------------------------------------------------------
 * Application initialization
 *---------------------------------------------------------------------------*/
void AppInitialize (void) {
  tid_AppThread          = osThreadNew(AppThread,          NULL, NULL);
  tid_FaultTriggerThread = osThreadNew(FaultTriggerThread, NULL, NULL);
}
