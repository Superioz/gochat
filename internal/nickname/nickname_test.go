package nickname

import (
	"testing"
)

// tests for the random seed to be applicable.
// the test fails if the random generator generates
// more values of one kind than expected, therefore it would
// be not truly random.
func TestGetRandom(t *testing.T) {
	names := make(map[string]int)

	// get `names.length`x random names
	for i := 0; i < len(nameData.Names); i++ {
		name := GetRandom()
		size := names[name]

		names[name] = size + 1
	}

	// get greates amount of names
	var greatest int
	for _, s := range names {
		if s > greatest {
			greatest = s
		}
	}

	// check for greates generated value is bigger than
	// the different generates names
	perc := float32(len(nameData.Names)) / float32(len(names))
	if perc > float32(greatest) {
		t.Error("random value generated too many values of one kind")
	}
}
