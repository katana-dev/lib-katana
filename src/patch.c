#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include "patch.h"
#include "sysex.h"

#ifndef SPARSE_BUFFER_SIZE
#define SPARSE_BUFFER_SIZE 1040
#endif

char katana_patch_data_alloc(KatanaPatch* patch, unsigned char encoding);
short katana_patch_internal_offset_sparse(unsigned short sysex_offset);

KatanaPatch* katana_patch_new(FirmwareVersion firmware_version, unsigned char encoding)
{
	if(encoding != PATCH_ENCODING_SPARSE){
		return NULL;
	}
	
	KatanaPatch* patch = malloc(sizeof(KatanaPatch));
	if(patch == NULL){
		return NULL;
	}
	
	patch->format_version = CURRENT_PATCH_FORMAT_VERSION;
	patch->firmware_version = firmware_version;
	patch->encoding = encoding;
	
	if(!katana_patch_data_alloc(patch, encoding)){
		katana_patch_free(patch);
		return NULL;
	}

	return patch;
}

char katana_patch_data_alloc(KatanaPatch* patch, unsigned char encoding)
{
	switch (encoding) {
		// Sparse encoding.
		case PATCH_ENCODING_SPARSE:
			patch->data = calloc(SPARSE_BUFFER_SIZE, sizeof(unsigned char));
			if(patch->data != NULL){
				patch->data_size = SPARSE_BUFFER_SIZE;
				return true;
			}
			break;
	}
	
	return false;
}

void katana_patch_free(KatanaPatch* patch)
{
	if(patch != NULL){
		if(patch->data != NULL){
			free(patch->data);
		}
		
		free(patch);
	}
}

short katana_patch_read_value(KatanaPatch* patch, unsigned short sysex_offset, unsigned short size)
{
	// Check the arguments are valid.
	if(patch == NULL || sysex_offset > SYSEX_SHORT_MAX ||
		sysex_offset + size - 1 > SYSEX_SHORT_MAX){
		return -1;
	}
	
	short internal_offset_msb, internal_offset_lsb;
	
	switch (size) {
		//7bit unsigned char.
		case 1:
			internal_offset_msb = katana_patch_internal_offset(sysex_offset, patch->encoding);
			if(internal_offset_msb >= 0){
				return (short)patch->data[internal_offset_msb];
			}
			else if(internal_offset_msb == INTERNAL_OFFSET_DISCARD){
				return PATCH_PADDING_VALUE;
			}
			break;
		
		//14bit unsigned short.
		case 2:
			internal_offset_msb = katana_patch_internal_offset(sysex_offset, patch->encoding);
			internal_offset_lsb = katana_patch_internal_offset(sysex_offset+1, patch->encoding);
			if(internal_offset_msb >= 0 && internal_offset_lsb >= 0){
				return
					(patch->data[internal_offset_msb] * SYSEX_BYTE_MULTIPLIER) +
					patch->data[internal_offset_lsb];
			}
			else if(internal_offset_msb == INTERNAL_OFFSET_DISCARD &&
				internal_offset_lsb == INTERNAL_OFFSET_DISCARD){
				return PATCH_PADDING_VALUE;
			}
			break;
	}
	
	return -1;
}

short katana_patch_write_value(KatanaPatch* patch, unsigned short sysex_offset,
	unsigned short value, unsigned short size)
{
	// Check the arguments are valid.
	if(patch == NULL || sysex_offset > SYSEX_SHORT_MAX ||
		sysex_offset + size - 1 > SYSEX_SHORT_MAX){
		return -1;
	}
	
	short internal_offset_msb, internal_offset_lsb;
	
	switch (size) {
		//7bit unsigned char.
		case 1:
			if(value > SYSEX_BYTE_MAX){
				return -1;
			}
			internal_offset_msb = katana_patch_internal_offset(sysex_offset, patch->encoding);
			if(internal_offset_msb >= 0){
				patch->data[internal_offset_msb] = value & SYSEX_BYTE_MAX;
				return 1;
			}
			else if(internal_offset_msb == INTERNAL_OFFSET_DISCARD){
				return 0;
			}
			break;
		
		//14bit unsigned short.
		case 2:
			if(value > SYSEX_SHORT_MAX){
				return -1;
			}
			internal_offset_msb = katana_patch_internal_offset(sysex_offset, patch->encoding);
			internal_offset_lsb = katana_patch_internal_offset(sysex_offset+1, patch->encoding);
			if(internal_offset_msb >= 0 && internal_offset_lsb >= 0){
				patch->data[internal_offset_msb] = (value / SYSEX_BYTE_MULTIPLIER) % SYSEX_BYTE_MULTIPLIER;
				patch->data[internal_offset_lsb] = value % SYSEX_BYTE_MULTIPLIER;
				return 2;
			}
			else if(internal_offset_msb == INTERNAL_OFFSET_DISCARD &&
				internal_offset_lsb == INTERNAL_OFFSET_DISCARD){
				return 0;
			}
			break;
	}
	
	return -1;
}

short katana_patch_read_sysex(KatanaPatch* patch, unsigned short sysex_offset,
	unsigned char* buffer, unsigned short length)
{
	// Check the arguments are valid.
	if(patch == NULL || buffer == NULL || length == 0 ||
		sysex_offset > SYSEX_SHORT_MAX || sysex_offset + length - 1 > SYSEX_SHORT_MAX){
		return -1;
	}
	
	short internal_offset;
	short written = 0;
	
	// For every byte of sysex data, map it to an internal location.
	for (unsigned short i = 0; i < length; i++) {
		internal_offset = katana_patch_internal_offset(sysex_offset + i, patch->encoding);
		
		// Write all areas that are encoded.
		if(internal_offset >= 0){
			buffer[i] = patch->data[internal_offset];
			written++;
		}
		
		// Pad any other area.
		else{
			buffer[i] = PATCH_PADDING_BYTE;
		}
	}
	
	return written;
}

short katana_patch_write_sysex(KatanaPatch* patch, unsigned short sysex_offset,
	unsigned char* data, unsigned short data_size)
{
	// Check the arguments are valid.
	if(patch == NULL || data == NULL || data_size == 0 ||
		sysex_offset > SYSEX_SHORT_MAX || sysex_offset + data_size - 1 > SYSEX_SHORT_MAX){
		return -1;
	}
	
	short internal_offset;
	short written = 0;
	
	// For every byte of sysex data, map it to an internal location.
	for (unsigned short i = 0; i < data_size; i++) {
		internal_offset = katana_patch_internal_offset(sysex_offset + i, patch->encoding);
		
		// Ignore any negative values, writing all areas that look good.
		if(internal_offset >= 0){
			// Apply the byte mask to ensure everything is in 7bits format.
			patch->data[internal_offset] = data[i] & SYSEX_BYTE_MAX;
			written++;
		}
	}
	
	return written;
}

short katana_patch_internal_offset(unsigned short sysex_offset, unsigned char encoding)
{
	// When outside of theoretical 7bit ushort value, invalid.
	if(sysex_offset > SYSEX_SHORT_MAX){
		return INTERNAL_OFFSET_ARGUMENT_ERROR;
	}
	
	// When outside of largest known patch parameter, discard.
	else if(sysex_offset > PATCH_SYSEX_SHORT_MAX){
		return INTERNAL_OFFSET_DISCARD;
	}
	
	// Delegate to encoding specific function.
	switch (encoding) {
		case PATCH_ENCODING_SPARSE:
			return katana_patch_internal_offset_sparse(sysex_offset);
	}
	
	// Encoding not defined, invalid.
	return INTERNAL_OFFSET_ARGUMENT_ERROR;
}

short katana_patch_internal_offset_sparse(unsigned short sysex_offset)
{
	// First used section 0-106.
	if(sysex_offset <= 106) return sysex_offset;
	
	// Discard 107-191, 85 gap.
	// Next section 192-1058.
	else if(sysex_offset < 192) return INTERNAL_OFFSET_DISCARD;
	else if(sysex_offset <= 1058) return sysex_offset - 85;
	
	// Discard 1059-2063, 1005 gap (+ 85 = 1090).
	// Next section 2064-2106.
	else if(sysex_offset < 2064) return INTERNAL_OFFSET_DISCARD;
	else if(sysex_offset <= 2106) return sysex_offset - 1090;
	
	// Discard 2107-2303, 197 gap (+ 1090 = 1287).
	// Next section 2304-2326.
	else if(sysex_offset < 2304) return INTERNAL_OFFSET_DISCARD;
	else if(sysex_offset <= 2326) return sysex_offset - 1287;
	
	// Discard beyond max value.
	else return INTERNAL_OFFSET_DISCARD;
}
