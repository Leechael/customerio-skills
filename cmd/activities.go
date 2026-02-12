package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "activities",
		Short: "Manage activities",
	}

	ls := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List activities",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/activities", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	parent.AddCommand(ls)
	rootCmd.AddCommand(parent)
}
