package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "send",
		Short: "Send messages",
	}

	email := &cobra.Command{
		Use:   "email",
		Short: "Send an email",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := requireBody(cmd)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/send/email", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(email)

	push := &cobra.Command{
		Use:   "push",
		Short: "Send a push notification",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := requireBody(cmd)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/send/push", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(push)

	sms := &cobra.Command{
		Use:   "sms",
		Short: "Send an SMS",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			body, err := requireBody(cmd)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/send/sms", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	addBodyFlag(sms)

	parent.AddCommand(email, push, sms)
	rootCmd.AddCommand(parent)
}
