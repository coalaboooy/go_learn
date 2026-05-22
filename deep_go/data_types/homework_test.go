package main

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

// go test -v homework_test.go

func ToLittleEndian(number uint32) uint32 {
	var bytes [4]byte
	pointer := unsafe.Pointer(&number)
	for i := range 4 {
		currentByte := *(*uint8)(unsafe.Add(pointer, i))
		bytes[i] = currentByte
		//bytes[3-i] = currentByte  		     // Альтернатива - записывать в обратном порядке, а потом
	}
	new_number := binary.BigEndian.Uint32(bytes[:])
	//new_number := binary.LittleEndian.Uint32(bytes[:]) // Формировать число в LE
	// Работает, потому что по умолчанию BE.Uint32 сам переворачивает порядок
	return new_number
}

func TestСonversion(t *testing.T) {
	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}
