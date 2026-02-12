package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "imports",
		Short: "Manage imports",
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "Create an import",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := requireBody(cmd)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/imports", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(create)

	get := &cobra.Command{
		Use:   "get <id>",
		Short: "Get an import",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/imports/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	parent.AddCommand(create, get)
	rootCmd.AddCommand(parent)
}
