package create

import (
	"fmt"

	"github.com/cgalvisleon/et/utility"
	"github.com/spf13/cobra"
)

var CmdMicro = &cobra.Command{
	Use:   "micro [name schema]",
	Short: "Create project base type microservice.",
	Long:  "Template project to microservice include folder cmd, deployments, pkg, rest, test and web, with files .go required for making a microservice.",
	Run: func(cmd *cobra.Command, args []string) {
		packageName, err := utility.GetGoMod("module")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		name, err := PrompStr("Name", true)
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		schema, err := PrompStr("Schema", false)
		if err != nil {
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
