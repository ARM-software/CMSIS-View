/*
 * Copyright (c) 2016-2020 ARM Limited. All rights reserved.
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
 * -----------------------------------------------------------------------------
 *
 * $Revision:   V1.0.2
 *
 * Project:     Event Statistic
 * Title:       main.c
 *
 * -----------------------------------------------------------------------------
 */

#include "RTE_Components.h"             // Component selection
#include CMSIS_device_header

#include "EventRecorder.h"              // Keil::Compiler:Event Recorder

#include <math.h>
#include <stdio.h>
#include <stdlib.h>

#define TABLE_SIZE 1000
float sin_table[TABLE_SIZE];

// Calculate table with sine values
void CalcSinTable (void)  {
  unsigned int i, max_i;
  float f = 0.0;

  max_i = TABLE_SIZE - (rand () % 500);
  EventStartAv (15, max_i, 0);                  // Start group A, slot 15, passing the max_i variable
  for (i = 0; i < max_i; i++)  {
    if (i == 200)  {
       EventStartAv (0, max_i, 0);              // Start group A, slot 0, passing the max_i variable
    }

    sin_table[i] = sinf(f);
    f = f + (3.141592 / TABLE_SIZE);

    if (i == 800)  {                            // Measure 800 table entries
      EventStopA (0);                           // Stop group A, slot 0
    }
  }

   EventStopA (15);                              // Stop group A, slot 15 (stops also slots 0..14)
}

// Return number of sqrt operations to exceed sum
unsigned int FindSqrtSum (float max_sum)  {
  unsigned int i;
  float sqrt_sum;

  sqrt_sum = 0.0;
  for (i = 0; i < 10000; i++) {
    sqrt_sum += sqrtf((float) i);
    if (sqrt_sum > max_sum)  {
      return (i);
    }
  }
  return (i);
}

unsigned int j, num, MaxSqrtSum;


int main (void) {

  SystemCoreClockUpdate();                      // System Initialization

  EventRecorderInitialize (EventRecordAll, 1U); // Initialize and start Event Recorder
  EventRecorderClockUpdate();
  EventStartC (0);                              // start measurement event group C, slot 0
  printf ("Started\n");
  for (j = 0; j < 1000; j++)  {
    CalcSinTable ();                            // calculate table with sinus values

    EventStartB(0);                             // start group B, slot 0
    MaxSqrtSum = rand () / 65536;               // limit for sqrt calculation
    num = FindSqrtSum ((float) MaxSqrtSum);     // return number of sqrt operations
    EventStopBv(0, MaxSqrtSum, num);            // stop group B, slot 0, output values: MaxSqrtSum, num

    if (j % 10 == 0) {
      printf("Progress: %3d%%\r", j/10+1);
    }
  }

  printf ("Finished       \n");
  EventStopC(0);                                // stop measurement event group C, slot 0
  return 0;
}
