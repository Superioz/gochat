package input

import (
	"github.com/manifoldco/promptui"
	"github.com/superioz/gochat/internal/env"
)

// Prompts to choose a protocol or use
// the given `force` if it is not empty
func PromptChooseProtocol(force string) error {
	var result string
	if len(force) == 0 {
		prompt := promptui.Select{
			Label: "What protocol do you want to use?",
			Items: []string{"tcp", "amqp"},
		}

		_, r, err := prompt.Run()
		if err != nil {
			return err
		}
		result = r
	} else {
		result = force
	}

	env.SetDefaults(result)
	return nil
}
