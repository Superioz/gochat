package input

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Input string

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

type Command struct {
	RawInput  Input
	Label     string
	Arguments []string
}

func (c Command) IsValid() bool {
	return len(c.Label) > 0
}

func (c Command) Argument(index int) (string, error) {
	if len(c.Arguments) > index {
		return c.Arguments[index], nil
	}
	return "", errors.New("arguments out of bounds")
}

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
