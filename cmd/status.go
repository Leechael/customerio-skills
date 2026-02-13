package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Check API token and connectivity",
		RunE: func(cmd *cobra.Command, args []string) error {
			token := os.Getenv("CUSTOMERIO_API_TOKEN")
			if token == "" {
				fmt.Fprintln(os.Stderr, "CUSTOMERIO_API_TOKEN is not set.")
				fmt.Fprintln(os.Stderr, "")
				fmt.Fprintln(os.Stderr, "Set it with:")
				fmt.Fprintln(os.Stderr, "  export CUSTOMERIO_API_TOKEN=\"your-app-api-key\"")
				fmt.Fprintln(os.Stderr, "")
				fmt.Fprintln(os.Stderr, "Or use 1Password CLI:")
				fmt.Fprintln(os.Stderr, "  op run --env-file=.env -- cio status")
				return fmt.Errorf("missing CUSTOMERIO_API_TOKEN")
			}

			c, err := newClient()
			if err != nil {
				return err
			}

			_, err = c.Get("/v1/info/ip_addresses", nil)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Authentication failed.")
				fmt.Fprintln(os.Stderr, "")
				fmt.Fprintln(os.Stderr, "Your CUSTOMERIO_API_TOKEN may be invalid or expired.")
				fmt.Fprintln(os.Stderr, "Get a new key from: https://fly.customer.io/settings/api_credentials")
				return fmt.Errorf("authentication failed")
			}

			masked := "***"
			if len(token) >= 8 {
				masked = token[:4] + "..." + token[len(token)-4:]
			}

			if jsonOutput {
				return printObject(map[string]any{
					"authenticated": true,
					"region":        region,
					"token":         masked,
				})
			}

			fmt.Printf("Authenticated (%s)\n", masked)
			fmt.Printf("Region: %s\n", region)
			return nil
		},
	})
}
