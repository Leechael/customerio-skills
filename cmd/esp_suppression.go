package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "esp-suppression",
		Short: "Manage ESP suppressions",
	}

	search := &cobra.Command{
		Use:   "search",
		Short: "Search ESP suppressions",
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
			data, err := c.Post("/v1/esp_suppression/search", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(search)

	get := &cobra.Command{
		Use:   "get <email>",
		Short: "Get ESP suppression for an email",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/esp_suppression/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	suppress := &cobra.Command{
		Use:   "suppress <email>",
		Short: "Suppress an email",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := readBody(cmd)
			if err != nil {
				return err
			}
			data, err := c.Put(fmt.Sprintf("/v1/esp_suppression/%s", args[0]), body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(suppress)

	unsuppress := &cobra.Command{
		Use:   "unsuppress <email>",
		Short: "Unsuppress an email",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Delete(fmt.Sprintf("/v1/esp_suppression/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	parent.AddCommand(search, get, suppress, unsuppress)
	rootCmd.AddCommand(parent)
}
