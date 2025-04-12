package commands

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:               "cli-quickstart",
		Short:             "A quickstart example",
		DisableAutoGenTag: true,
		Run: func(c *cobra.Command, _ []string) {

		},
	}
	return command
}
