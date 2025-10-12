package command

import (
	"log"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/cli/command/start"
	"github.com/spf13/cobra"
)

// Init performs all the available commands initialization.
func Init() {
	root := &cobra.Command{
		Use:   "fate-seekers-server-cli",
		Short: "CLI tool to manage and run the FateSeekers game server",
		Long:  `FateSeekers Server CLI provides a command launch the FateSeekers multiplayer game server`,
	}

	root.CompletionOptions.DisableDefaultCmd = true

	start.Init(root)

	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}
