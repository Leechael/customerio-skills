package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "objects",
		Short: "Manage objects",
	}

	types := &cobra.Command{
		Use:   "types",
		Short: "List object types",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/object_types", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	search := &cobra.Command{
		Use:   "search",
		Short: "Search objects",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := readBody(cmd)
			if err != nil {
				return err
			}
			if body == nil {
				body = json.RawMessage(`{}`)
			}
			data, err := c.Post("/v1/objects", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(search)

	get := &cobra.Command{
		Use:   "get <type-id> <object-id>",
		Short: "Get object attributes",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/objects/%s/%s/attributes", args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	relationships := &cobra.Command{
		Use:   "relationships <type-id> <object-id>",
		Short: "Get object relationships",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/objects/%s/%s/relationships", args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	parent.AddCommand(types, search, get, relationships)
	rootCmd.AddCommand(parent)
}
