#include "sysex.h"

unsigned short offset_from_sysex(unsigned char sysex[2])
{
	return
		((sysex[0] & SYSEX_BYTE_MAX) * SYSEX_BYTE_MASK) +
		(sysex[1] & SYSEX_BYTE_MAX);
}

int offset_to_sysex(unsigned short offset, unsigned char* buffer)
{
	buffer[0] = (offset / 0x80) % 0x80;
	buffer[1] = offset % 0x80;
	return 2;
}

KatanaAddress address_from_sysex(unsigned char sysex[4])
{
	KatanaAddress address;
	address.region = offset_from_sysex((unsigned char[]){sysex[0], sysex[1]});
	address.offset = offset_from_sysex((unsigned char[]){sysex[2], sysex[3]});
	return address;
}

int address_to_sysex(KatanaAddress address, unsigned char* buffer)
{
	int step = offset_to_sysex(address.region, buffer);
	return step + offset_to_sysex(address.offset, buffer + step);
}
