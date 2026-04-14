package cmd

import "github.com/spf13/cobra"

var getCmd = &cobra.Command{
	Use: "get [path]",
	Short: "send get request",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}
