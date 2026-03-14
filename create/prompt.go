package create

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func PromptCreate() {
	prompt := promptui.Select{
		Label: "What do you want created?",
		Items: []string{
			"Microservice",
			"Modelo",
			"Rpc",
		},
	}

	opt, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch opt {
	case 0:
		err := CmdMicro.Execute()
		if err != nil {
			fmt.Printf("Error executing microservice command: %v", err)
			return
		}
	case 1:
		err := CmdModelo.Execute()
		if err != nil {
			fmt.Printf("Error executing modelo command: %v", err)
			return
		}
	case 2:
		err := CmdRpc.Execute()
		if err != nil {
			fmt.Printf("Error executing rpc command: %v", err)
			return
		}
	}
}
