package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/leechael/cio/internal/client"
	"github.com/leechael/cio/internal/output"
	"github.com/spf13/cobra"
)

var (
	region string
	jqExpr string
)

var rootCmd = &cobra.Command{
	Use:   "cio",
	Short: "CLI for Customer.io App API",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&region, "region", "us", "API region: us or eu")
	rootCmd.PersistentFlags().StringVar(&jqExpr, "jq", "", "jq expression to filter JSON output")
}

var newClient = func() (*client.Client, error) {
	return client.New(region)
}

func printJSON(data json.RawMessage) error {
	return output.Print(data, jqExpr)
}

func readBody(cmd *cobra.Command) (json.RawMessage, error) {
	body, _ := cmd.Flags().GetString("body")
	if body != "" {
		return json.RawMessage(body), nil
	}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		if len(data) > 0 {
			return json.RawMessage(data), nil
		}
	}
	return nil, nil
}

func requireBody(cmd *cobra.Command) (json.RawMessage, error) {
	body, err := readBody(cmd)
	if err != nil {
		return nil, err
	}
	if body == nil {
		return nil, fmt.Errorf("request body is required (use -body flag or pipe via stdin)")
	}
	return body, nil
}

func addBodyFlag(cmd *cobra.Command) {
	cmd.Flags().String("body", "", "JSON request body")
}
