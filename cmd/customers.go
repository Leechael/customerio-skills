package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "customers",
		Short: "Manage customers",
	}

	get := &cobra.Command{
		Use:   "get <id>",
		Short: "Get customer attributes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/customers/%s/attributes", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	search := &cobra.Command{
		Use:   "search",
		Short: "Search customers by email or filter",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			email, _ := cmd.Flags().GetString("email")
			if email != "" {
				q := url.Values{"email": {email}}
				data, err := c.Get("/v1/customers", q)
				if err != nil {
					return err
				}
				return printJSON(data)
			}
			body, err := requireBody(cmd)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/customers", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	search.Flags().String("email", "", "Search by email address")
	addBodyFlag(search)

	ls := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List customers by attributes",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := requireBody(cmd)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/customers/attributes", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(ls)

	activities := &cobra.Command{
		Use:   "activities <id>",
		Short: "Get customer activities",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/customers/%s/activities", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	messages := &cobra.Command{
		Use:   "messages <id>",
		Short: "Get customer messages",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/customers/%s/messages", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	segments := &cobra.Command{
		Use:   "segments <id>",
		Short: "Get customer segments",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/customers/%s/segments", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	relationships := &cobra.Command{
		Use:   "relationships <id>",
		Short: "Get customer relationships",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/customers/%s/relationships", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	subPrefs := &cobra.Command{
		Use:   "sub-prefs <id>",
		Short: "Get customer subscription preferences",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get(fmt.Sprintf("/v1/customers/%s/subscription_preferences", args[0]), nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	parent.AddCommand(get, search, ls, activities, messages, segments, relationships, subPrefs)
	rootCmd.AddCommand(parent)
}
