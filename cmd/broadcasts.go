package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "broadcasts",
		Short: "Manage broadcasts",
	}

	ls := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List broadcasts",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/broadcasts", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	get := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a broadcast",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/broadcasts/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	trigger := &cobra.Command{
		Use:   "trigger <id>",
		Short: "Trigger a broadcast",
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
			if body == nil {
				body = json.RawMessage(`{}`)
			}
			data, err := c.Post(fmt.Sprintf("/v1/campaigns/%s/triggers", args[0]), body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(trigger)

	triggers := &cobra.Command{
		Use:   "triggers <id>",
		Short: "List broadcast triggers",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/broadcasts/%s/triggers", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	triggerStatus := &cobra.Command{
		Use:   "trigger-status <id> <trigger-id>",
		Short: "Get broadcast trigger status",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/broadcasts/%s/triggers/%s", args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	triggerErrors := &cobra.Command{
		Use:   "trigger-errors <id> <trigger-id>",
		Short: "Get broadcast trigger errors",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/broadcasts/%s/triggers/%s/errors", args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	bcastActions := &cobra.Command{
		Use:   "actions <id>",
		Short: "List broadcast actions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/broadcasts/%s/actions", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	bcastAction := &cobra.Command{
		Use:   "action <id> <action-id>",
		Short: "Get or update a broadcast action",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/v1/broadcasts/%s/actions/%s", args[0], args[1])
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
	addBodyFlag(bcastAction)

	bcastMetrics := &cobra.Command{
		Use:   "metrics <id>",
		Short: "Get broadcast metrics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/broadcasts/%s/metrics", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	bcastLinkMetrics := &cobra.Command{
		Use:   "link-metrics <id>",
		Short: "Get broadcast link metrics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/broadcasts/%s/metrics/links", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	bcastActionMetrics := &cobra.Command{
		Use:   "action-metrics <id> <action-id>",
		Short: "Get broadcast action metrics",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/broadcasts/%s/actions/%s/metrics", args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	bcastActionLinkMetrics := &cobra.Command{
		Use:   "action-link-metrics <id> <action-id>",
		Short: "Get broadcast action link metrics",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/broadcasts/%s/actions/%s/metrics/links", args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	bcastMessages := &cobra.Command{
		Use:   "messages <id>",
		Short: "Get broadcast messages",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/broadcasts/%s/messages", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	bcastTranslation := &cobra.Command{
		Use:   "translation <id> <action-id> <lang>",
		Short: "Get or update a broadcast translation",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/v1/broadcasts/%s/actions/%s/language/%s", args[0], args[1], args[2])
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
	addBodyFlag(bcastTranslation)

	parent.AddCommand(ls, get, trigger, triggers, triggerStatus, triggerErrors, bcastActions, bcastAction, bcastMetrics, bcastLinkMetrics, bcastActionMetrics, bcastActionLinkMetrics, bcastMessages, bcastTranslation)
	rootCmd.AddCommand(parent)
}
