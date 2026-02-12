package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "index",
		Short: "Data index lookups",
	}

	attributes := &cobra.Command{
		Use:   "attributes",
		Short: "List indexed attributes",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/index/attributes", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	events := &cobra.Command{
		Use:   "events",
		Short: "List indexed events",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/index/events", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	parent.AddCommand(attributes, events)
	rootCmd.AddCommand(parent)
}
