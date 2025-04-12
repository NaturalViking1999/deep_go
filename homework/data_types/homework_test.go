package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

// uint32 - имеет 4-х байтовое представление 0x.00.00.00.FF, где FF - наибольший значимый байт
// Значит поменять нужно 4 байта
// 1 байт в Hex имеет 8-ми битовое представление, поэтому 1 шаг сдвига это 8 бит.
func ToLittleEndian(number uint32) uint32 {
	// Сдвигаем на 24 влево, так как 4-ый байт меняется местами с 1-ым, то есть делает 3 шага
	num1 := (number & 0x000000FF) << 24
	// Сдвигаем на 8 влево, так как 3-ий байт меняется местами со 2-ым, то есть делает 1 шаг
	num2 := (number & 0x0000FF00) << 8
	// Сдвигаем на 8 вправо, так как 2-ой байт меняется местами со 3-ым, то есть делает 1 шаг
	num3 := (number & 0x00FF0000) >> 8
	// Сдвигаем на 24 вправо, так как 1-ый байт меняется местами с 4-ым, то есть делает 3 шага
	num4 := (number & 0xFF000000) >> 24
	return num1 | num2 | num3 | num4 // need to implement
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
