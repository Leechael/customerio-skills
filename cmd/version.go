package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version of cio",
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonOutput {
				return printObject(map[string]string{
					"cli":     "cio",
					"version": version,
					"commit":  commit,
					"built":   date,
				})
			}
			fmt.Printf("cio %s (commit: %s, built: %s)\n", version, commit, date)
			return nil
		},
	})
}
