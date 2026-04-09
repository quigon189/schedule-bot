package cmd

import (
	"encoding/json"
	"fmt"
	"schedule-cli/internal/models"
	"schedule-cli/internal/prompts"
	"time"

	"github.com/spf13/cobra"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Administrator methods",
}

var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Get sessions list from server",
	Run: func(cmd *cobra.Command, args []string) {
		type session struct {
			ID        string      `json:"uuid"`
			UserAgent string      `json:"user_agent"`
			ClientIP  string      `json:"client_ip"`
			CreatedAt time.Time   `json:"created_at"`
			User      models.User `json:"user"`
		}
		sessions := []session{}
		query := make(map[string]string)
		if err := apiClient.Get("/admin/sessions", query, &sessions); err != nil {
			prompts.ShowError(fmt.Sprintf("Failed to get sessions: %v", err))
			return
		}
		msg, _ := json.MarshalIndent(sessions, "", "  ")
		prompts.ShowSuccess(fmt.Sprintf("Sessions:\n%s\n", msg))

	},
}

func init() {
	adminCmd.AddCommand(sessionsCmd)
}
