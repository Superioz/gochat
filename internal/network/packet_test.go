package network

import (
	"reflect"
	"testing"
)

func TestPacketTyping(t *testing.T) {
	var p Packet = &HandshakePacket{}

	if reflect.TypeOf(p) != reflect.TypeOf(&HandshakePacket{}) {
		t.Error("returned value is not equals expected value!")
	}
}

func TestPacketCasting(t *testing.T) {
	var p Packet = &HandshakePacket{Passed: true}

	switch p.(type) {
	case *HandshakePacket:
		h := p.(*HandshakePacket)

		if !h.Passed {
			t.Error("couldn't cast packet to hand shake!")
		}
		break
	default:
		t.Error("returned value is not equals expected value!")
	}
}
