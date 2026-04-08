package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure base parameters",
}

var serverCmd = &cobra.Command{
	Use:   "server [base_url]",
	Short: "Set base url for server",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := cfgManager.SetBaseURL(args[0]); err != nil {
			fmt.Printf("Failed to set base url: %v\n", err)
			return
		}
		fmt.Printf("Set base url to %s\n", args[0])
	},
}

func init() {
	configCmd.AddCommand(serverCmd)
}
