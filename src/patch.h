#ifndef LIB_KATANA_PATCH_H
#define LIB_KATANA_PATCH_H

#include "firmware.h"

#ifndef CURRENT_PATCH_FORMAT_VERSION
#define CURRENT_PATCH_FORMAT_VERSION 1
#define PATCH_SYSEX_SHORT_MAX 2326
#define PATCH_PADDING_VALUE 0
#define PATCH_PADDING_BYTE 0x00
#endif

enum {
	PATCH_ENCODING_SPARSE = 0
};

enum {
	INTERNAL_OFFSET_ARGUMENT_ERROR = -1,
	INTERNAL_OFFSET_DISCARD = -2
};

typedef struct {
	// region or patch nr?
	unsigned char format_version;
	FirmwareVersion firmware_version;
	unsigned char encoding;
	unsigned char* data;
	unsigned short data_size;
} KatanaPatch;

/**
 * Create a new KatanaPatch.
 * @param  firmware_version Version number for the target.
 * @param  encoding         Encoding type for the internal storage.
 * @return                  Pointer to a new KatanaPatch.
 */
KatanaPatch* katana_patch_new(FirmwareVersion firmware_version, unsigned char encoding);

/**
 * Frees memory of a KatanaPatch.
 * @param patch Pointer to the KatanaPatch to be freed.
 */
void katana_patch_free(KatanaPatch* patch);

/**
 * Reads a scalar value at the given offset.
 * Note: Will return padding value for discarded address.
 * @param  patch        The KatanaPatch to read from.
 * @param  sysex_offset The sysex patch offset to read from.
 * @param  size         The size in bytes for this value (accepts 1 or 2).
 * @return The scalar value read. Or -1 if the arguments were invalid.
 */
short katana_patch_read_value(KatanaPatch* patch, unsigned short sysex_offset, unsigned short size);

/**
 * Write a scalar value at the given offset.
 * Note: Will return padding value for discarded address.
 * @param  patch        The KatanaPatch to write from.
 * @param  sysex_offset The sysex patch offset to write to.
 * @param  value        The scalar value to write.
 * @param  size         The size in bytes for this value (accepts 1 or 2).
 * @return The size of the value written. 0 when discarded, or -1 if the arguments were invalid.
 */
short katana_patch_write_value(KatanaPatch* patch, unsigned short sysex_offset,
	unsigned short value, unsigned short size);

/**
 * Reads an arbitrary block of sysex data from the patch's data structure.
 * Note: Will write padding bytes for every invalid or discarded address in the requested range.
 * @param  patch        The KatanaPatch to read from.
 * @param  sysex_offset The sysex patch offset to start reading from.
 * @param  buffer       Pointer to a buffer that will hold the read data.
 * @param  length       Amount of bytes to read (your buffer should be large enough).
 * @return The number of bytes read that were not padding bytes.
 *         Or -1 if the arguments were invalid.
 */
short katana_patch_read_sysex(KatanaPatch* patch, unsigned short sysex_offset, unsigned char* buffer,
	unsigned short length);

/**
 * Writes an arbitrary block of sysex data to the patch's data structure.
 * @param patch        The KatanaPatch to write to.
 * @param sysex_offset The sysex patch offset to start writing from.
 * @param data         Pointer to the sysex data array.
 * @param data_size    Size of the sysex data array.
 * @return Number of bytes successfully written to internal storage.
 *         Or -1 if the arguments were invalid.
 */
short katana_patch_write_sysex(KatanaPatch* patch, unsigned short sysex_offset, unsigned char* data,
	unsigned short data_size);

/**
 * Find the internal offset for a given sysex offset.
 * @param  sysex_offset The patch offset from a sysex mapping perspective.
 * @param  encoding     The encoding scheme used for internal data.
 * @return Positive values are an internal offset number, negative values are errors or special cases.
 */
short katana_patch_internal_offset(unsigned short sysex_offset, unsigned char encoding);

#endif /* end of include guard: LIB_KATANA_PATCH_H */
