package cmd

import (
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
		prompts.ShowSuccess("Sessions:", sessions)

	},
}

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Users of schedule service",
}

var usersListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List users",
	Run: func(cmd *cobra.Command, args []string) {
		users := []map[string]any{}
		fullName := prompts.AskString("Full Name", false, func(val any) error { return nil })
		query := map[string]string{}
		if fullName != "" {
			query["full_name"] = fullName
		}
		if err := apiClient.Get("/admin/users", query, &users); err != nil {
			prompts.ShowError(fmt.Sprintf("Failed to get users: %v", err))
			return
		}

		prompts.ShowSuccess("Users:", users)
	},
}

var changeUserPasswordCmd = &cobra.Command{
	Use:   "passsword",
	Short: "Change user password",
	Run: func(cmd *cobra.Command, args []string) {
		username := prompts.AskUsername("")
		password := prompts.AskPassword(true)

		users := []map[string]any{}

		query := map[string]string{
			"username": username,
		}
		if err := apiClient.Get("/admin/users", query, &users); err != nil {
			prompts.ShowError(fmt.Sprintf("Failed to get users: %v", err))
			return
		}

		if len(users) == 0 {
			prompts.ShowError("Invalid username")
			return
		}

		userID, ok := users[0]["id"].(float64)
		if !ok {
			prompts.ShowError("Invalid user id in response")
			return
		}

		body := map[string]any{
			"user_id":           int(userID),
			"new_password": password,
		}

		if err := apiClient.Patch("/admin/users/password", body, nil); err != nil {
			prompts.ShowError(fmt.Sprintf("Failed to patch password: %v", err))
			return
		}

		prompts.ShowSuccess("Password changed", nil)
	},
}

func init() {
	usersCmd.AddCommand(usersListCmd, changeUserPasswordCmd)
	adminCmd.AddCommand(sessionsCmd, usersCmd)
}
