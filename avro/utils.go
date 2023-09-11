// Description: Avro utils
// Author: Pixie79
// ============================================================================
// package avro

package avro

import (
	"encoding/binary"
)

// intToBytes converts an integer to a byte array.
//
// It takes an integer as a parameter and returns a byte array.
func intToBytes(n int) []byte {
	byteArray := make([]byte, 4)
	binary.BigEndian.PutUint32(byteArray, uint32(n))
	return byteArray
}

// addZeroToStart adds a zero byte at the start of the given byte array.
//
// byteArray: the input byte array
// Returns: the modified byte array with a zero byte added at the start
func addZeroToStart(byteArray []byte) []byte {
	return append([]byte{0}, byteArray...)
}

// EncodedBuffer returns the encoded buffer of an integer.
//
// It takes an integer as input and returns a byte slice that represents the encoded buffer.
func EncodedBuffer(i int) []byte {
	return addZeroToStart(intToBytes(i))
}
