/*----------------------------------------------------------------------------
 * Name:    interface.h
 * Purpose: Non-Secure to Secure interface
 *----------------------------------------------------------------------------*/

#include <stdint.h>

/* Non-secure callable functions */
extern void Secure_TriggerFault (uint32_t fault_id);
