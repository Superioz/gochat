package network

import (
	"bytes"
	"testing"
)

func TestWriteUint16(t *testing.T) {
	b := bytes.NewBuffer([]byte{})

	val := uint16(20)
	WriteUint16(b, val)
	val2, _ := ReadUint16(b)

	if val != val2 {
		t.Errorf("written value is not equals expected value. %d!=%d", val, val2)
	}
}

func TestWriteString(t *testing.T) {
	b := bytes.NewBuffer([]byte{})

	val := "test"
	WriteString(b, val)
	val2, _ := ReadString(b)

	if val != val2 {
		t.Errorf("written value is not equals expected value. %s!=%s", val, val2)
	}
}

func TestWriteBool(t *testing.T) {
	b := bytes.NewBuffer([]byte{})

	WriteBool(b, true)
	val2, _ := ReadBool(b)

	if true != val2 {
		t.Error("written value is not equals expected value.")
	}
}
