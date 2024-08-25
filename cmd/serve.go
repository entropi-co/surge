package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"surge/internal/api"
	"surge/internal/conf"
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
	config, err := conf.LoadFromEnvironments()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration from environments\n")
		return nil
	} else {
		logrus.Println("Loaded configuration from environments")
	}

	if config.Logging.EnableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugln("EnableDebug debugging")
	}

	surgeAPI := api.NewSurgeAPI(config)
	defer surgeAPI.CloseDatabaseConnection()

	surgeAPI.ListenAndServe(cmd.Context(), "0.0.0.0:3000")
	return nil
}
