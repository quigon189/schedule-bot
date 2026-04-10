package cmd

import (
	"fmt"
	"schedule-cli/internal/prompts"

	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Get current user info",
	Run: func(cmd *cobra.Command, args []string) {
		resp := make(map[string]any)
		query := make(map[string]string)
		if err := apiClient.Get("/user", query, &resp); err != nil {
			prompts.ShowError(fmt.Sprintf("Failed to get user info: %v", err))
			return
		}
		prompts.ShowSuccess("User info:", resp)
	},
}
