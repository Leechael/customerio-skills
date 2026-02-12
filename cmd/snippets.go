package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "snippets",
		Short: "Manage snippets",
	}

	ls := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List snippets",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/snippets", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	upsert := &cobra.Command{
		Use:   "upsert",
		Short: "Create or update snippets",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := requireBody(cmd)
			if err != nil {
				return err
			}
			data, err := c.Put("/v1/snippets", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(upsert)

	rm := &cobra.Command{
		Use:     "rm <name>",
		Aliases: []string{"delete"},
		Short:   "Delete a snippet",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Delete(fmt.Sprintf("/v1/snippets/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	parent.AddCommand(ls, upsert, rm)
	rootCmd.AddCommand(parent)
}
