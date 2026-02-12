package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "transactional",
		Short: "Manage transactional messages",
	}

	ls := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List transactional messages",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/transactional", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	get := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a transactional message",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/transactional/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	txMetrics := &cobra.Command{
		Use:   "metrics <id>",
		Short: "Get transactional message metrics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/transactional/%s/metrics", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	txLinkMetrics := &cobra.Command{
		Use:   "link-metrics <id>",
		Short: "Get transactional message link metrics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/transactional/%s/metrics/links", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	txContent := &cobra.Command{
		Use:   "content <id>",
		Short: "Get or update transactional message content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/v1/transactional/%s/content", args[0])
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
	addBodyFlag(txContent)

	txTranslation := &cobra.Command{
		Use:   "translation <id> <lang>",
		Short: "Get or update transactional message translation",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/v1/transactional/%s/language/%s", args[0], args[1])
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
	addBodyFlag(txTranslation)

	deliveries := &cobra.Command{
		Use:   "deliveries <id>",
		Short: "Get transactional message deliveries",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/transactional/%s/deliveries", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	parent.AddCommand(ls, get, txMetrics, txLinkMetrics, txContent, txTranslation, deliveries)
	rootCmd.AddCommand(parent)
}
