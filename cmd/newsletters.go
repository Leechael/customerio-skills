package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "newsletters",
		Short: "Manage newsletters",
	}

	ls := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List newsletters",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/newsletters", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	get := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a newsletter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/newsletters/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	rm := &cobra.Command{
		Use:     "rm <id>",
		Aliases: []string{"delete"},
		Short:   "Delete a newsletter",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Delete(fmt.Sprintf("/v1/newsletters/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	contents := &cobra.Command{
		Use:   "contents <id>",
		Short: "List newsletter contents",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/newsletters/%s/contents", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	content := &cobra.Command{
		Use:   "content <id> <content-id>",
		Short: "Get or update newsletter content",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/v1/newsletters/%s/contents/%s", args[0], args[1])
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

	nlMetrics := &cobra.Command{
		Use:   "metrics <id>",
		Short: "Get newsletter metrics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/newsletters/%s/metrics", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	nlLinkMetrics := &cobra.Command{
		Use:   "link-metrics <id>",
		Short: "Get newsletter link metrics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/newsletters/%s/metrics/links", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	contentMetrics := &cobra.Command{
		Use:   "content-metrics <id> <content-id>",
		Short: "Get newsletter content metrics",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/newsletters/%s/contents/%s/metrics", args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	contentLinkMetrics := &cobra.Command{
		Use:   "content-link-metrics <id> <content-id>",
		Short: "Get newsletter content link metrics",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/newsletters/%s/contents/%s/metrics/links", args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	nlMessages := &cobra.Command{
		Use:   "messages <id>",
		Short: "Get newsletter messages",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/newsletters/%s/messages", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	nlTranslation := &cobra.Command{
		Use:   "translation <id> <lang>",
		Short: "Get or update a newsletter translation",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/v1/newsletters/%s/language/%s", args[0], args[1])
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
	addBodyFlag(nlTranslation)

	testGroups := &cobra.Command{
		Use:   "test-groups <id>",
		Short: "Get newsletter test groups",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/newsletters/%s/test_groups", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	testGroupTranslation := &cobra.Command{
		Use:   "test-group-translation <id> <group-id> <lang>",
		Short: "Get or update a newsletter test group translation",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/v1/newsletters/%s/test_groups/%s/language/%s", args[0], args[1], args[2])
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
	addBodyFlag(testGroupTranslation)

	parent.AddCommand(ls, get, rm, contents, content, nlMetrics, nlLinkMetrics, contentMetrics, contentLinkMetrics, nlMessages, nlTranslation, testGroups, testGroupTranslation)
	rootCmd.AddCommand(parent)
}
