package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var _rootCommand = cobra.Command{
	Use:   "surge serve",
	Short: "Surge: authentication and identity server written in go",
	RunE:  handleRootCommand,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var nameString = color.New(color.Underline).Sprintf(func() string {
			if cmd.Root() == cmd {
				return cmd.Root().Name()
			}
			return cmd.Root().Name() + " " + cmd.Name()
		}())

		var versionString = func() string {
			if cmd.Version == "" {
				return ""
			}
			return color.New(color.Bold).Sprintf(" v" + cmd.Version)
		}()

		fmt.Printf(color.New(color.Bold).Sprintf("Running ")+"%s%s\n", nameString, versionString)
	},
}

func BuildRootCommand() *cobra.Command {
	_rootCommand.AddCommand(buildServeCommand(), buildMigrateCommand())

	return &_rootCommand
}

func handleRootCommand(cmd *cobra.Command, args []string) error {
	cmd.SetArgs([]string{"migrate"})
	if err := cmd.Execute(); err != nil {
		return err
	}
	cmd.SetArgs([]string{"serve"})
	return cmd.Execute()
}
