package cmd

import (
	"github.com/spf13/cobra"
)

var serveCommand = cobra.Command{
	Use:   "serve",
	Short: "Start Surge and listen to requests",
	RunE:  handleServeCommand,
}

func buildServeCommand() *cobra.Command {
	return &serveCommand
}

func handleServeCommand(cmd *cobra.Command, args []string) error {
	return nil
}
