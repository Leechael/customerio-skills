package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	parent := &cobra.Command{
		Use:   "info",
		Short: "General information",
	}

	ipAddresses := &cobra.Command{
		Use:   "ip-addresses",
		Short: "Get IP addresses",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/info/ip_addresses", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}

	parent.AddCommand(ipAddresses)
	rootCmd.AddCommand(parent)
}
