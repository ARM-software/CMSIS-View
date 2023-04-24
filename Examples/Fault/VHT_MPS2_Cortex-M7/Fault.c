/*----------------------------------------------------------------------------
 * Name:    Fault.c
 * Purpose: Fault example program
 *----------------------------------------------------------------------------*/

#include <stdio.h>

#include "RTE_Components.h"
#include  CMSIS_device_header

#include "cmsis_os2.h"

#include "ARM_Fault.h"
#include "ARM_FaultTrigger.h"

#include "EventRecorder.h"

extern osThreadId_t tid_AppThread;
extern osThreadId_t tid_FaultTriggerThread;

/* STDIO initialize function */
extern int stdio_init (void);

/* Global Thread IDs (for debug) */
osThreadId_t tid_AppThread;
osThreadId_t tid_FaultTriggerThread;

/*---------------------------------------------------------------------------
 * Application thread
 *---------------------------------------------------------------------------*/
static __NO_RETURN void AppThread (void *argument) {
  volatile uint32_t counter = 0U;

  (void)argument;

  for (;;) {
    counter++;
    osDelay(100U);
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
  printf(" - 0: terminate the example\r\n");
  printf(" - 1: trigger the escalated hard fault\r\n");
  printf(" - 2: trigger the data access (precise) bus fault\r\n");
  printf(" - 3: trigger the data access (imprecise) bus fault\r\n");
  printf(" - 4: trigger the instruction execution bus fault\r\n");
  printf(" - 5: trigger the no coprocessor usage fault\r\n");
  printf(" - 6: trigger the undefined instruction usage fault\r\n");
  printf(" - 7: trigger the divide by 0 usage fault\r\n\r\n");
  printf("Input>");

  for (;;) {
    ch = (char)getchar();                       // Read character from console (blocking)
    if (ch == '0') {
      putchar(0x04);                            // Shutdown the simulator
    } else {
      ARM_FaultTrigger((uint32_t)(ch - '0'));   // Trigger a fault
    }
  }
}

/*---------------------------------------------------------------------------
 * Application main function
 *---------------------------------------------------------------------------*/
int main (void) {

  SystemCoreClockUpdate();                      // Update SystemCoreClock variable

  stdio_init();                                 // Initialize STDIO
  if (ARM_FaultOccurred() != 0U) {              // If fault information exists
    printf("\r\n\r\n- System Restarted -\r\n\r\n");
  }

  EventRecorderInitialize(EventRecordAll, 1U);  // Initialize and start EventRecorder

  if (ARM_FaultOccurred() != 0U) {              // If fault information exists
    ARM_FaultRecord();                          // Output decoded fault information via EventRecorder
    EventRecorderStop();                        // Stop EventRecorder
  } else {                                      // If fault information does not exist
    ARM_FaultClear();                           // Clear (initialize) fault information
  }

  osKernelInitialize();                         // Initialize CMSIS-RTOS2
                                                // Create threads
  tid_AppThread          = osThreadNew(AppThread,          NULL, NULL);
  tid_FaultTriggerThread = osThreadNew(FaultTriggerThread, NULL, NULL);
  osKernelStart();                              // Start thread execution

  for (;;);
}
