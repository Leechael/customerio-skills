package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "exports",
		Short: "Manage exports",
	}

	ls := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List exports",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/exports", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	createCustomers := &cobra.Command{
		Use:   "create-customers",
		Short: "Create a customer export",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := requireBody(cmd)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/exports/customers", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(createCustomers)

	createDeliveries := &cobra.Command{
		Use:   "create-deliveries",
		Short: "Create a deliveries export",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := requireBody(cmd)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/exports/deliveries", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(createDeliveries)

	get := &cobra.Command{
		Use:   "get <id>",
		Short: "Get an export",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/exports/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	download := &cobra.Command{
		Use:   "download <id>",
		Short: "Download an export",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/exports/%s/download", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	parent.AddCommand(ls, createCustomers, createDeliveries, get, download)
	rootCmd.AddCommand(parent)
}
