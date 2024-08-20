package cmd

import (
	"github.com/spf13/cobra"
	"surge/internal/api"
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
	surgeAPI := api.NewSurgeAPI()
	surgeAPI.ListenAndServe(cmd.Context(), "0.0.0.0:3000")
	return nil
}
