package nickname

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

// just some random names taken from:
// https://pairedlife.com/friendship/badass-nicknames
var data = []byte(`{
	"names": [
		"Aspect", "Bender", "Big Papa", "Bowser", "Bruise", "Cannon", "Clink", "Cobra", "Colt", "Crank",
		"Creep", "Daemon", "Decay", "Diablo", "Doom", "Dracula", "Dragon", "Fender", "Fester", "Fisheye",
		"Flack", "Gargoyle", "Grave", "Gunner", "Hash", "Hashtag", "Indominus", "Ironclas", "Killer", "Knuckles",
		"Kraken", "Lynch", "Mad Dog", "O'Doyle", "Psycho", "Ranger", "Ratchet", "Reaper", "Rigs", "Ripley",
		"Roadkill", "Ronin", "Rubble", "Sasquatch", "Scar", "Shiver", "Skinner", "Skull Crusher", "Slasher", "Steelshot",
		"Surge", "Scythe", "Trip", "Trooper", "Tweek", "Vein", "Void", "Wadon", "Wraith", "Zero"
	]
}`)

// defines a struct for the `data`
// json fetches the names with the field label `json:"names"`
type NameData struct {
	Names []string `json:"names"`
}

var random *rand.Rand
var nameData = NameData{}

// Initializes the `nickname` package and creates a new random instance
// with a custom seed. Also, fetches all names from `jsondata.go`
func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))

	// unmarshals the data from `data`
	err := json.Unmarshal(data, &nameData)
	if err != nil {
		log.Fatal(err)
	}
}

// Returns a random name from the list above
// uses the already initialized random instance
func GetRandom() string {
	return nameData.Names[random.Intn(len(nameData.Names))]
}
