package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "campaigns",
		Short: "Manage campaigns",
	}

	ls := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List campaigns",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/campaigns", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	get := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a campaign",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/campaigns/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	actions := &cobra.Command{
		Use:   "actions <id>",
		Short: "List campaign actions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/campaigns/%s/actions", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	action := &cobra.Command{
		Use:   "action <id> <action-id>",
		Short: "Get or update a campaign action",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/v1/campaigns/%s/actions/%s", args[0], args[1])
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
	addBodyFlag(action)

	metrics := &cobra.Command{
		Use:   "metrics <id>",
		Short: "Get campaign metrics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/campaigns/%s/metrics", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	linkMetrics := &cobra.Command{
		Use:   "link-metrics <id>",
		Short: "Get campaign link metrics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/campaigns/%s/metrics/links", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	actionMetrics := &cobra.Command{
		Use:   "action-metrics <id> <action-id>",
		Short: "Get campaign action metrics",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/campaigns/%s/actions/%s/metrics", args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	actionLinkMetrics := &cobra.Command{
		Use:   "action-link-metrics <id> <action-id>",
		Short: "Get campaign action link metrics",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/campaigns/%s/actions/%s/metrics/links", args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	journeyMetrics := &cobra.Command{
		Use:   "journey-metrics <id>",
		Short: "Get campaign journey metrics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/campaigns/%s/journey_metrics", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	messages := &cobra.Command{
		Use:   "messages <id>",
		Short: "Get campaign messages",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/campaigns/%s/messages", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	translation := &cobra.Command{
		Use:   "translation <id> <action-id> <lang>",
		Short: "Get or update a campaign translation",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/v1/campaigns/%s/actions/%s/language/%s", args[0], args[1], args[2])
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
	addBodyFlag(translation)

	parent.AddCommand(ls, get, actions, action, metrics, linkMetrics, actionMetrics, actionLinkMetrics, journeyMetrics, messages, translation)
	rootCmd.AddCommand(parent)
}
