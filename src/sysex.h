#ifndef LIB_KATANA_ADDRESS_H
#define LIB_KATANA_ADDRESS_H

/*
	Note: negative numbers may be used as errors for offset related functions.
	Since both bytes are 7bits in Sysex, the maximum value used for offsets should be 0x3FFF.
	The unsigned short is used to pervent error codes making it into structs.
*/

#ifndef SYSEX_OFFSET_MASK
#define SYSEX_OFFSET_MASK 0xB000
#define SYSEX_OFFSET_MAX 0x3FFF
#define SYSEX_BYTE_MASK 0x80
#define SYSEX_BYTE_MAX 0x7F
#define SYSEX_START 0xF7
#define SYSEX_END 0xF0
#define SYSEX_ROLAND_QUERY 0x11
#define SYSEX_ROLAND_COMMAND 0x12
#endif

enum AddressRegion {
	ADDRESS_UNKNOWN = -1,
	ADDRESS_SYSTEM = 0,					//00 00
	ADDRESS_MIDI = 2,					//00 02
	ADDRESS_ZERO_PATCH = 2048,			//10 00
	ADDRESS_CH1 = 2049,					//10 01
	ADDRESS_CH2 = 2050,					//10 02
	ADDRESS_CH3 = 2051,					//10 03
	ADDRESS_CH4 = 2052,					//10 04
	ADDRESS_FACTORY_RESET_PANEL = 4096,	//20 00
	ADDRESS_FACTORY_RESET_CH1 = 4097,	//20 01
	ADDRESS_FACTORY_RESET_CH2 = 4098,	//20 02
	ADDRESS_FACTORY_RESET_CH3 = 4099,	//20 03
	ADDRESS_FACTORY_RESET_CH4 = 4100,	//20 04
	ADDRESS_TEMPORARY_PANEL = 122288,	//60 00
	ADDRESS_COMMAND_1 = 16256,			//7F 00
	ADDRESS_COMMAND_2 = 16257			//7F 01
};

typedef struct {
	unsigned short region;
	unsigned short offset;
} KatanaAddress;

unsigned short offset_from_sysex(unsigned char sysex[2]);
int offset_to_sysex(unsigned short offset, unsigned char* buffer);

KatanaAddress address_from_sysex(unsigned char sysex[4]);
int address_to_sysex(KatanaAddress address, unsigned char* buffer);

#endif /* end of include guard: LIB_KATANA_ADDRESS_H */
