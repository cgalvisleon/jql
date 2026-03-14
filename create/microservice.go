package create

import (
	"fmt"

	"github.com/cgalvisleon/et/utility"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var CmdMicro = &cobra.Command{
	Use:   "micro [name schema]",
	Short: "Create project base type microservice.",
	Long:  "Template project to microservice include folder cmd, deployments, pkg, rest, test and web, with files .go required for making a microservice.",
	Run: func(cmd *cobra.Command, args []string) {
		packageName, err := utility.GoMod("module")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		var name string
		var schema string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Name").
					Prompt("? ").
					Validate(func(s string) error {
						if len(s) == 0 {
							return fmt.Errorf("name is required")
						}
						return nil
					}).
					Value(&name),
				huh.NewInput().
					Title("Schema").
					Prompt("? ").
					Value(&schema),
			),
		)

		if err := form.Run(); err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		err = MkMicroservice(packageName, name, schema)
		if err != nil {
			fmt.Printf("Command failed %v\n", err)
			return
		}
	},
}

func MkMicroservice(projectName, name, schema string) error {
	ProgressAdd(6)

	ProgressNext()
	return nil
}
