package create

import (
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
)

func PrompStr(label string, require bool) (string, error) {
	validate := func(input string) error {
		if len(input) == 0 && require {
			return fmt.Errorf("invalid %s", label)
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func PrompInt(label string, require bool) (int, error) {
	validate := func(input string) error {
		if len(input) == 0 && require {
			return fmt.Errorf("invalid %s", label)
		}

		_, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("invalid %s", label)
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	value, err := prompt.Run()

	if err != nil {
		return 0, err
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return result, nil
}
