define memory mem with size = 4G;

define symbol use_ROM1 = (__ROM1_SIZE != 0x0);
define symbol use_ROM2 = (__ROM2_SIZE != 0x0);
define symbol use_ROM3 = (__ROM3_SIZE != 0x0);
define symbol use_ROM4 = (__ROM4_SIZE != 0x0);

define symbol use_RAM1 = (__RAM1_SIZE != 0x0);
define symbol use_RAM2 = (__RAM2_SIZE != 0x0);
define symbol use_RAM3 = (__RAM3_SIZE != 0x0);
define symbol use_RAM4 = (__RAM4_SIZE != 0x0);

if (use_ROM1)
{
  define region ROM1_region = mem:[from __ROM1_BASE to (__ROM1_BASE+__ROM1_SIZE-1)];
}
else
{
  define region ROM1_region = [];
}

if (use_ROM2)
{
  define region ROM2_region = mem:[from __ROM2_BASE to (__ROM1_BASE+__ROM2_SIZE-1)];
}
else
{
  define region ROM2_region = [];
}
if (use_ROM3)
{
  define region ROM3_region = mem:[from __ROM3_BASE to (__ROM3_BASE+__ROM3_SIZE-1)];
}
else
{
  define region ROM3_region = [];
}
if (use_ROM4)
{
  define region ROM4_region = mem:[from __ROM4_BASE to (__ROM4_BASE+__ROM4_SIZE-1)];
}
else
{
  define region ROM4_region = [];
}

define region ROM_region = ROM1_region |  ROM2_region | ROM3_region | ROM4_region;

if (use_RAM1)
{
  define region RAM1_region = mem:[from __RAM1_BASE to (__RAM1_BASE+__RAM1_SIZE-1)];
}
else
{
  define region RAM1_region = [];
}
if (use_RAM2)
{
  define region RAM2_region = mem:[from __RAM2_BASE to (__RAM1_BASE+__RAM2_SIZE-1)];
}
else
{
  define region RAM2_region = [];
}
if (use_RAM3)
{
  define region RAM3_region = mem:[from __RAM3_BASE to (__RAM3_BASE+__RAM3_SIZE-1)];
}
else
{
  define region RAM3_region = [];
}
if (use_RAM4)
{
  define region RAM4_region = mem:[from __RAM4_BASE to (__RAM4_BASE+__RAM4_SIZE-1)];
}
else
{
  define region RAM4_region = [];
}
define region RAM_region = RAM1_region | RAM2_region | RAM3_region | RAM4_region;

do not initialize  { section .noinit };
initialize by copy { readwrite };
if (isdefinedsymbol(__USE_DLIB_PERTHREAD))
{
  // Required in a multi-threaded application
  initialize by copy with packing = none { section __DLIB_PERTHREAD };
}

place at address mem:__ROM1_BASE { readonly section .intvec };

if (!isempty(ROM_region))
{
  place in ROM_region  { readonly };
}

if (!isempty(RAM_region))
{
  define block CSTACK     with alignment = 8, size = __STACK_SIZE     { };
  define block PROC_STACK with alignment = 8, size = 0 { };
  define block HEAP       with alignment = 8, size = __HEAP_SIZE       { };
  place in RAM_region  { readwrite, block CSTACK, block PROC_STACK, block HEAP };
}
