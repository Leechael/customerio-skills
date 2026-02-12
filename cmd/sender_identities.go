package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "sender-identities",
		Short: "Manage sender identities",
	}

	ls := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List sender identities",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/sender_identities", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	get := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a sender identity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/sender_identities/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	usedBy := &cobra.Command{
		Use:   "used-by <id>",
		Short: "Get sender identity usage",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/sender_identities/%s/used_by", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	parent.AddCommand(ls, get, usedBy)
	rootCmd.AddCommand(parent)
}
