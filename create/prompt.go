package create

import (
	"fmt"

	"github.com/cgalvisleon/et/logs"
	"github.com/charmbracelet/huh"
)

func PromptCreate() {
	var opt string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What do you want created?").
				Options(
					huh.NewOption("Microservice", "microservice"),
					huh.NewOption("Modelo", "modelo"),
					huh.NewOption("Rpc", "rpc"),
				).
				Value(&opt),
		),
	)

	if err := form.Run(); err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch opt {
	case "microservice":
		err := CmdMicro.Execute()
		if err != nil {
			logs.Errorf("Error executing microservice command: %v", err)
			return
		}
	case "modelo":
		err := CmdModelo.Execute()
		if err != nil {
			logs.Errorf("Error executing modelo command: %v", err)
			return
		}
	case "rpc":
		err := CmdRpc.Execute()
		if err != nil {
			logs.Errorf("Error executing rpc command: %v", err)
			return
		}
	}
}
