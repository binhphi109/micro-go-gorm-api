package cmd

import (
	"github.com/spf13/cobra"

	"github.com/sample/sample-server/api"
	"github.com/sample/sample-server/core"
)

var serverCmd = &cobra.Command{
	Use:          "server",
	Short:        "Run the Kasse server",
	RunE:         serverCmdF,
	SilenceUsage: true,
}

func init() {
	RootCmd.AddCommand(serverCmd)
	RootCmd.RunE = serverCmdF
}

func serverCmdF(command *cobra.Command, args []string) error {
	return runServer()
}

func runServer() error {

	server := core.NewServer()

	api.Init(server)

	server.Start()

	return nil
}
