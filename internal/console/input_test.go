package console

import "testing"

func TestCommand_IsValid(t *testing.T) {
	var raw Input = "label some arguments"
	var raw2 Input = "label"
	var raw3 Input = ""

	cmd := raw.AsCommand()
	if !cmd.IsValid() || len(cmd.Arguments) != 2 {
		t.Error("command not valid")
	}

	cmd2 := raw2.AsCommand()
	if !cmd2.IsValid() {
		t.Error("command not valid")
	}

	cmd3 := raw3.AsCommand()
	if cmd3.IsValid() {
		t.Error("command valid despite invalidation expected")
	}
}
