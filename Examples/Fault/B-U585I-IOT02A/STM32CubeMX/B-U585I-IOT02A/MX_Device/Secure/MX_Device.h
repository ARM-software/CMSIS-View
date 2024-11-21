/******************************************************************************
 * File Name   : MX_Device.h
 * Date        : 15/11/2024 10:26:27
 * Description : STM32Cube MX parameter definitions
 * Note        : This file is generated with a generator out of the
 *               STM32CubeMX project and its generated files (DO NOT EDIT!)
 ******************************************************************************/

#ifndef MX_DEVICE_H__
#define MX_DEVICE_H__

/* MX_Device.h version */
#define MX_DEVICE_VERSION                       0x01000000


/*------------------------------ I2C1           -----------------------------*/
#define MX_I2C1                                 1

/* Filter Settings */
#define MX_I2C1_ANF_ENABLE                      1
#define MX_I2C1_DNF                             0

/* Peripheral Clock Frequency */
#define MX_I2C1_PERIPH_CLOCK_FREQ               160000000

/* Pins */

/* I2C1_SCL */
#define MX_I2C1_SCL_Pin                         PB8
#define MX_I2C1_SCL_GPIO_Pin                    GPIO_PIN_8
#define MX_I2C1_SCL_GPIOx                       GPIOB
#define MX_I2C1_SCL_GPIO_Mode                   GPIO_MODE_AF_OD
#define MX_I2C1_SCL_GPIO_PuPd                   GPIO_NOPULL
#define MX_I2C1_SCL_GPIO_Speed                  GPIO_SPEED_FREQ_LOW
#define MX_I2C1_SCL_GPIO_AF                     GPIO_AF4_I2C1

/* I2C1_SDA */
#define MX_I2C1_SDA_Pin                         PB9
#define MX_I2C1_SDA_GPIO_Pin                    GPIO_PIN_9
#define MX_I2C1_SDA_GPIOx                       GPIOB
#define MX_I2C1_SDA_GPIO_Mode                   GPIO_MODE_AF_OD
#define MX_I2C1_SDA_GPIO_PuPd                   GPIO_NOPULL
#define MX_I2C1_SDA_GPIO_Speed                  GPIO_SPEED_FREQ_LOW
#define MX_I2C1_SDA_GPIO_AF                     GPIO_AF4_I2C1

/*------------------------------ I2C2           -----------------------------*/
#define MX_I2C2                                 1

/* Filter Settings */
#define MX_I2C2_ANF_ENABLE                      1
#define MX_I2C2_DNF                             0

/* Peripheral Clock Frequency */
#define MX_I2C2_PERIPH_CLOCK_FREQ               160000000

/* Pins */

/* I2C2_SCL */
#define MX_I2C2_SCL_Pin                         PH4
#define MX_I2C2_SCL_GPIO_Pin                    GPIO_PIN_4
#define MX_I2C2_SCL_GPIOx                       GPIOH
#define MX_I2C2_SCL_GPIO_Mode                   GPIO_MODE_AF_OD
#define MX_I2C2_SCL_GPIO_PuPd                   GPIO_NOPULL
#define MX_I2C2_SCL_GPIO_Speed                  GPIO_SPEED_FREQ_LOW
#define MX_I2C2_SCL_GPIO_AF                     GPIO_AF4_I2C2

/* I2C2_SDA */
#define MX_I2C2_SDA_Pin                         PH5
#define MX_I2C2_SDA_GPIO_Pin                    GPIO_PIN_5
#define MX_I2C2_SDA_GPIOx                       GPIOH
#define MX_I2C2_SDA_GPIO_Mode                   GPIO_MODE_AF_OD
#define MX_I2C2_SDA_GPIO_PuPd                   GPIO_NOPULL
#define MX_I2C2_SDA_GPIO_Speed                  GPIO_SPEED_FREQ_LOW
#define MX_I2C2_SDA_GPIO_AF                     GPIO_AF4_I2C2

/*------------------------------ SPI2           -----------------------------*/
#define MX_SPI2                                 1

/* Peripheral Clock Frequency */
#define MX_SPI2_PERIPH_CLOCK_FREQ               160000000

/* Pins */

/* SPI2_MISO */
#define MX_SPI2_MISO_Pin                        PD3
#define MX_SPI2_MISO_GPIO_Pin                   GPIO_PIN_3
#define MX_SPI2_MISO_GPIOx                      GPIOD
#define MX_SPI2_MISO_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_SPI2_MISO_GPIO_PuPd                  GPIO_NOPULL
#define MX_SPI2_MISO_GPIO_Speed                 GPIO_SPEED_FREQ_LOW
#define MX_SPI2_MISO_GPIO_AF                    GPIO_AF5_SPI2

/* SPI2_MOSI */
#define MX_SPI2_MOSI_Pin                        PD4
#define MX_SPI2_MOSI_GPIO_Pin                   GPIO_PIN_4
#define MX_SPI2_MOSI_GPIOx                      GPIOD
#define MX_SPI2_MOSI_GPIO_Mode                  GPIO_MODE_AF_PP
#define MX_SPI2_MOSI_GPIO_PuPd                  GPIO_NOPULL
#define MX_SPI2_MOSI_GPIO_Speed                 GPIO_SPEED_FREQ_LOW
#define MX_SPI2_MOSI_GPIO_AF                    GPIO_AF5_SPI2

/* SPI2_NSS */
#define MX_SPI2_NSS_Pin                         PB12
#define MX_SPI2_NSS_GPIO_Pin                    GPIO_PIN_12
#define MX_SPI2_NSS_GPIOx                       GPIOB
#define MX_SPI2_NSS_GPIO_Mode                   GPIO_MODE_AF_PP
#define MX_SPI2_NSS_GPIO_PuPd                   GPIO_NOPULL
#define MX_SPI2_NSS_GPIO_Speed                  GPIO_SPEED_FREQ_LOW
#define MX_SPI2_NSS_GPIO_AF                     GPIO_AF5_SPI2

/* SPI2_SCK */
#define MX_SPI2_SCK_Pin                         PD1
#define MX_SPI2_SCK_GPIO_Pin                    GPIO_PIN_1
#define MX_SPI2_SCK_GPIOx                       GPIOD
#define MX_SPI2_SCK_GPIO_Mode                   GPIO_MODE_AF_PP
#define MX_SPI2_SCK_GPIO_PuPd                   GPIO_NOPULL
#define MX_SPI2_SCK_GPIO_Speed                  GPIO_SPEED_FREQ_LOW
#define MX_SPI2_SCK_GPIO_AF                     GPIO_AF5_SPI2

/*------------------------------ UART4          -----------------------------*/
#define MX_UART4                                1

/* Pins */

/* UART4_RX */
#define MX_UART4_RX_Pin                         PC11
#define MX_UART4_RX_GPIO_Pin                    GPIO_PIN_11
#define MX_UART4_RX_GPIOx                       GPIOC
#define MX_UART4_RX_GPIO_Mode                   GPIO_MODE_AF_PP
#define MX_UART4_RX_GPIO_PuPd                   GPIO_NOPULL
#define MX_UART4_RX_GPIO_Speed                  GPIO_SPEED_FREQ_LOW
#define MX_UART4_RX_GPIO_AF                     GPIO_AF8_UART4

/* UART4_TX */
#define MX_UART4_TX_Pin                         PC10
#define MX_UART4_TX_GPIO_Pin                    GPIO_PIN_10
#define MX_UART4_TX_GPIOx                       GPIOC
#define MX_UART4_TX_GPIO_Mode                   GPIO_MODE_AF_PP
#define MX_UART4_TX_GPIO_PuPd                   GPIO_NOPULL
#define MX_UART4_TX_GPIO_Speed                  GPIO_SPEED_FREQ_LOW
#define MX_UART4_TX_GPIO_AF                     GPIO_AF8_UART4

/*------------------------------ USB_OTG_FS     -----------------------------*/
#define MX_USB_OTG_FS                           1

/* Handle */
#define MX_USB_OTG_FS_HANDLE                    hpcd_USB_OTG_FS

/* Virtual mode */
#define MX_USB_OTG_FS_VM                        Device_Only
#define MX_USB_OTG_FS_Device_Only               1

/* Pins */

/* USB_OTG_FS_DM */
#define MX_USB_OTG_FS_DM_Pin                    PA11
#define MX_USB_OTG_FS_DM_GPIO_Pin               GPIO_PIN_11
#define MX_USB_OTG_FS_DM_GPIOx                  GPIOA
#define MX_USB_OTG_FS_DM_GPIO_Mode              GPIO_MODE_AF_PP
#define MX_USB_OTG_FS_DM_GPIO_PuPd              GPIO_NOPULL
#define MX_USB_OTG_FS_DM_GPIO_Speed             GPIO_SPEED_FREQ_LOW
#define MX_USB_OTG_FS_DM_GPIO_AF                GPIO_AF10_USB

/* USB_OTG_FS_DP */
#define MX_USB_OTG_FS_DP_Pin                    PA12
#define MX_USB_OTG_FS_DP_GPIO_Pin               GPIO_PIN_12
#define MX_USB_OTG_FS_DP_GPIOx                  GPIOA
#define MX_USB_OTG_FS_DP_GPIO_Mode              GPIO_MODE_AF_PP
#define MX_USB_OTG_FS_DP_GPIO_PuPd              GPIO_NOPULL
#define MX_USB_OTG_FS_DP_GPIO_Speed             GPIO_SPEED_FREQ_LOW
#define MX_USB_OTG_FS_DP_GPIO_AF                GPIO_AF10_USB

#endif  /* MX_DEVICE_H__ */
