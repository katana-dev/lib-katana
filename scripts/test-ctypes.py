import os
from ctypes import *

# Route to the library.
script_dir = os.path.dirname(os.path.abspath(__file__))
so_path = os.path.join(script_dir, "../src/.libs/libkatana.so")

# Load and run with ctypes.
libkatana = CDLL(so_path)

def offset(a, b):
	return libkatana.offset_from_sysex((c_ubyte * 2)(a, b))

print offset(0x00, 0x00)
print offset(0x00, 0x02)
print offset(0x10, 0x00)
print offset(0x10, 0x01)
print offset(0x10, 0x02)
print offset(0x10, 0x03)
print offset(0x10, 0x04)
print offset(0x20, 0x00)
print offset(0x20, 0x01)
print offset(0x20, 0x02)
print offset(0x20, 0x03)
print offset(0x20, 0x04)
print offset(0x60, 0x00)
print offset(0x7F, 0x00)
print offset(0x7F, 0x01)

def formatSysex(address):
	return '#[' + str.join(' ', ["%02X" % v for v in address]) + ']'

class KatanaAddress(Structure):
	_fields_ = [("region", c_ushort),
				("offset", c_ushort)]

libkatana.address_from_sysex.restype = KatanaAddress

def printAddress(a, b, c, d):
	address = libkatana.address_from_sysex((c_ubyte * 4)(a, b, c, d))
	print address.region
	print address.offset
	buff = (c_ubyte * 4)(0, 0, 0, 0)
	libkatana.address_to_sysex(address, buff)
	print formatSysex(buff)

printAddress(0x00, 0x00, 0x00, 0x94)
printAddress(0x10, 0x00, 0x00, 0x07)

print libkatana.katana_patch_internal_offset(2106, 0)
print libkatana.katana_patch_internal_offset(2107, 0)
print libkatana.katana_patch_internal_offset(2303, 0)
print libkatana.katana_patch_internal_offset(2304, 0)

# Write to FX1 acsim block starting from a few bytes too early.
# Should be a silently lossy operation.
patch = libkatana.katana_patch_new(1, 0)
print patch

buffer = (c_ubyte * 8)(1, 2, 3, 4, 5, 6, 0x2A, 0xFF)
print [v for v in buffer]
print libkatana.katana_patch_write_sysex(patch, 2061, buffer, 8)
print libkatana.katana_patch_read_sysex(patch, 2061, buffer, 8)
print [v for v in buffer]

print libkatana.katana_patch_read_value(patch, 2067, 1)
print libkatana.katana_patch_write_value(patch, 2067, 123, 1)
print libkatana.katana_patch_read_value(patch, 2067, 1)
print libkatana.katana_patch_write_value(patch, 2065, 4012, 2)
print libkatana.katana_patch_read_value(patch, 2065, 2)

print libkatana.katana_patch_read_sysex(patch, 2061, buffer, 8)
print [v for v in buffer]

libkatana.katana_patch_free(patch)
