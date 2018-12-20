package console

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Input string

// gets the `input` as command format
// which is <label> <args:>
func (i Input) AsCommand() Command {
	f := strings.Fields(string(i))
	if len(f) == 0 {
		return Command{RawInput: i}
	}

	var a []string
	if len(f) > 1 {
		a = f[1:]
	}
	return Command{RawInput: i, Label: f[0], Arguments: a}
}

// represents a command consisting of a label and arguments
type Command struct {
	RawInput  Input
	Label     string
	Arguments []string
}

// checks if the length of the label is greater than zero
func (c Command) IsValid() bool {
	return len(c.Label) > 0
}

// gets an argument with `index` safely
// returns an `error` if the `index` does not exist
func (c Command) Argument(index int) (string, error) {
	if len(c.Arguments) > index {
		return c.Arguments[index], nil
	}
	return "", errors.New("arguments out of bounds")
}

// listens to the input of the console and returns
// a channel where new messages get channeled to
func ListenToConsole() chan Input {
	i := make(chan Input)
	go func(chan Input) {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			i <- Input(s.Text())
		}
	}(i)
	return i
}
