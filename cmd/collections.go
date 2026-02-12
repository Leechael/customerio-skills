package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "collections",
		Short: "Manage collections",
	}

	ls := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List collections",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/collections", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	get := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/collections/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "Create a collection",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := requireBody(cmd)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/collections", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(create)

	update := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := requireBody(cmd)
			if err != nil {
				return err
			}
			data, err := c.Put(fmt.Sprintf("/v1/collections/%s", args[0]), body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(update)

	rm := &cobra.Command{
		Use:     "rm <id>",
		Aliases: []string{"delete"},
		Short:   "Delete a collection",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Delete(fmt.Sprintf("/v1/collections/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	content := &cobra.Command{
		Use:   "content <id>",
		Short: "Get or update collection content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/v1/collections/%s/content", args[0])
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
	addBodyFlag(content)

	parent.AddCommand(ls, get, create, update, rm, content)
	rootCmd.AddCommand(parent)
}
