package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// Writes an unsigned int with byte length 2 to a byte buffer
// Uses byte order of BigEndian
func WriteUint16(buf *bytes.Buffer, i uint16) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)

	err := binary.Write(buf, binary.BigEndian, b)
	if err != nil {
		fmt.Println(err)
	}
}

// Writes a string with its byte length to a byte buffer
// Uses byte order of BigEndian
func WriteString(buf *bytes.Buffer, s string) {
	WriteUint16(buf, uint16(len(s)))

	err := binary.Write(buf, binary.BigEndian, []byte(s))
	if err != nil {
		fmt.Println(err)
	}
}

// Writes a boolean to a byte buffer
// Uses byte order of BigEndian
func WriteBool(buf *bytes.Buffer, b bool) {
	var bit byte
	if b {
		bit = 1
	} else {
		bit = 0
	}

	err := binary.Write(buf, binary.BigEndian, bit)
	if err != nil {
		fmt.Println(err)
	}
}

// Reads an unsigned int with byte length 2 from byte buffer
// Uses byte order of BigEndian
func ReadUint16(buf *bytes.Buffer) (uint16, error) {
	if buf.Len() < 2 {
		return 0, errors.New("buffer size lower than 2")
	}
	return binary.BigEndian.Uint16(buf.Next(2)), nil
}

// Reads a string with its length from byte buffer
// Uses byte order of BigEndian
func ReadString(buf *bytes.Buffer) (string, error) {
	l, err := ReadUint16(buf)
	if err != nil {
		return "", err
	}
	return string(buf.Next(int(l))), nil
}

// Reads a boolean from byte buffer
// Uses byte order of BigEndian
func ReadBool(buf *bytes.Buffer) (bool, error) {
	i, err := buf.ReadByte()

	if err != nil {
		return false, err
	}
	var b bool
	if i == 1 {
		b = true
	} else {
		b = false
	}
	return b, nil
}
