package commands

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:               "academy",
		Short:             "A academy example",
		DisableAutoGenTag: true,
		Run: func(c *cobra.Command, _ []string) {

		},
	}
	return command
}
