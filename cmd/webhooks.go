package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "webhooks",
		Short: "Manage reporting webhooks",
	}

	ls := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List reporting webhooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/reporting_webhooks", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	get := &cobra.Command{
		Use:   "get <id>",
		Short: "Get or update a reporting webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/v1/reporting_webhooks/%s", args[0])
			body, err := readBody(cmd)
			if err != nil {
				return err
			}
			if body != nil {
				data, err := c.Put(path, body)
				if err != nil {
					return err
				}
				return printJSON(data)
			}
			data, err := c.Get(path, nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(get)

	create := &cobra.Command{
		Use:   "create",
		Short: "Create a reporting webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := requireBody(cmd)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/reporting_webhooks", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(create)

	rm := &cobra.Command{
		Use:     "rm <id>",
		Aliases: []string{"delete"},
		Short:   "Delete a reporting webhook",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Delete(fmt.Sprintf("/v1/reporting_webhooks/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	parent.AddCommand(ls, get, create, rm)
	rootCmd.AddCommand(parent)
}
