package nickname

import (
	"testing"
)

func TestGetRandom(t *testing.T) {
	names := make(map[string]int)

	for i := 0; i < len(nameData.Names); i++ {
		name := GetRandom()
		size := names[name]

		names[name] = size + 1
	}

	var greatest int
	for _, s := range names {
		if s > greatest {
			greatest = s
		}
	}

	perc := float32(len(nameData.Names)) / float32(len(names))
	if perc > float32(greatest) {
		t.Error("random value generated too many values of one kind")
	}
}
