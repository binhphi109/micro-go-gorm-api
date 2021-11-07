package cmd

import (
	"github.com/spf13/cobra"
)

type Command = cobra.Command

func Run(args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.Execute()
}

var RootCmd = &cobra.Command{
	Use:   "sample",
	Short: "sample",
}

func init() {
	RootCmd.PersistentFlags().StringP("config", "c", "", "Configuration file to use.")
}
