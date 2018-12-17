package network

import (
	"reflect"
	"testing"
)

func TestPacketTyping(t *testing.T) {
	var p Packet = &MessagePacket{}

	if reflect.TypeOf(p) != reflect.TypeOf(&MessagePacket{}) {
		t.Error("returned value is not equals expected value!")
	}
}

func TestPacketCasting(t *testing.T) {
	var p Packet = &MessagePacket{Message: "msg"}

	switch p.(type) {
	case *MessagePacket:
		h := p.(*MessagePacket)

		if h.Message != "msg" {
			t.Error("couldn't cast packet to hand shake!")
		}
		break
	default:
		t.Error("returned value is not equals expected value!")
	}
}
