package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "segments",
		Short: "Manage segments",
	}

	ls := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List segments",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/segments", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	get := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a segment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/segments/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "Create a segment",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := requireBody(cmd)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/segments", body)
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
		Short:   "Delete a segment",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Delete(fmt.Sprintf("/v1/segments/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	count := &cobra.Command{
		Use:   "count <id>",
		Short: "Get segment customer count",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/segments/%s/customer_count", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	members := &cobra.Command{
		Use:   "members <id>",
		Short: "Get segment membership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/segments/%s/membership", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	deps := &cobra.Command{
		Use:   "deps <id>",
		Short: "Get segment dependencies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/segments/%s/dependencies", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	parent.AddCommand(ls, get, create, rm, count, members, deps)
	rootCmd.AddCommand(parent)
}
