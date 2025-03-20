/*-----------------------------------------------------------------------------
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

/*-----------------------------------------------------------------------------
 * Application thread
 *----------------------------------------------------------------------------*/
static __NO_RETURN void AppThread (void *argument) {
  volatile uint32_t counter = 0U;

  (void)argument;

  for (;;) {
    counter++;
    osDelay(100U);
  }
}

/*-----------------------------------------------------------------------------
 * Fault trigger thread
 *----------------------------------------------------------------------------*/
static __NO_RETURN void FaultTriggerThread (void *argument) {
  char ch;

  (void)argument;

  // Display user interface message
  printf("\r\n--- Fault example ---\r\n\r\n");
  printf("To trigger a fault please input a corresponding number:\r\n");
  printf(" - 0: Terminate the example\r\n");
  printf(" - 1: Data access (precise) Memory Management fault\r\n");
  printf(" - 2: Data access (precise) Bus fault\r\n");
  printf(" - 3: Data access (imprecise) Bus fault\r\n");
  printf(" - 4: Instruction execution Bus fault\r\n");
  printf(" - 5: Undefined instruction Usage fault\r\n");
  printf(" - 6: Divide by 0 Usage fault\r\n\r\n");
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

/*-----------------------------------------------------------------------------
 * Application main function
 *----------------------------------------------------------------------------*/
int main (void) {

  SystemCoreClockUpdate();                      // Update SystemCoreClock variable

  SCB->SHCSR |= SCB_SHCSR_BUSFAULTENA_Msk |     // Enable BusFault
                SCB_SHCSR_USGFAULTENA_Msk;      // Enable UsageFault
  SCB->CCR   |= SCB_CCR_DIV_0_TRP_Msk;          // Enable divide by 0 trap

  /* Configure MPU
       - region 0: ROM                   - 0x00000000 .. 0x1FFFFFFF (end is extended to be able to trigger Bus fault)
       - region 1: RAM                   - 0x20000000 .. 0x3FFFFFFF (end is extended to be able to trigger Bus fault)
       - region 2: RAM (privileged only) - 0x20000000 .. 0x200000FF
       - region 3: Peripherals           - 0x40000000 .. 0x4FFFFFFF
  */
  ARM_MPU_Disable();

  ARM_MPU_SetRegion(ARM_MPU_RBAR(0U, 0x00000000), ARM_MPU_RASR_EX(0U, ARM_MPU_AP_RO,   ARM_MPU_ACCESS_NORMAL(ARM_MPU_CACHEP_NOCACHE, ARM_MPU_CACHEP_NOCACHE, 0U), 0x00U, ARM_MPU_REGION_SIZE_512MB));
  ARM_MPU_SetRegion(ARM_MPU_RBAR(1U, 0x20000000), ARM_MPU_RASR_EX(1U, ARM_MPU_AP_FULL, ARM_MPU_ACCESS_NORMAL(ARM_MPU_CACHEP_NOCACHE, ARM_MPU_CACHEP_NOCACHE, 0U), 0x00U, ARM_MPU_REGION_SIZE_512MB));
  ARM_MPU_SetRegion(ARM_MPU_RBAR(2U, 0x20000000), ARM_MPU_RASR_EX(1U, ARM_MPU_AP_PRIV, ARM_MPU_ACCESS_NORMAL(ARM_MPU_CACHEP_NOCACHE, ARM_MPU_CACHEP_NOCACHE, 0U), 0x00U, ARM_MPU_REGION_SIZE_256B));
  ARM_MPU_SetRegion(ARM_MPU_RBAR(3U, 0x40000000), ARM_MPU_RASR_EX(1U, ARM_MPU_AP_FULL, ARM_MPU_ACCESS_DEVICE(0U)                                                , 0x00U, ARM_MPU_REGION_SIZE_256MB));

  ARM_MPU_Enable(MPU_CTRL_PRIVDEFENA_Msk);      // Enable Privileged Default

  stdio_init();                                 // Initialize STDIO

  if (ARM_FaultOccurred() != 0U) {              // If fault information exists
    printf("\r\n\r\n- System Restarted -\r\n\r\n");
  }

  EventRecorderInitialize(EventRecordAll, 1U);  // Initialize and start Event Recorder

  if (ARM_FaultOccurred() != 0U) {              // If fault information exists
    ARM_FaultRecord();                          // Output decoded fault information via Event Recorder
    EventRecorderStop();                        // Stop Event Recorder
  }

  osKernelInitialize();                         // Initialize CMSIS-RTOS2
                                                // Create threads
  tid_AppThread          = osThreadNew(AppThread,          NULL, NULL);
  tid_FaultTriggerThread = osThreadNew(FaultTriggerThread, NULL, NULL);
  osKernelStart();                              // Start thread execution

  for (;;) {                                    // Loop forever
  }
}
