package cmd

import (
	"fmt"
	"schedule-cli/internal/config"
	"schedule-cli/internal/models"
	"schedule-cli/internal/prompts"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
    Use:   "auth",
    Short: "Authentication commands",
}

var loginCmd = &cobra.Command{
    Use:   "login [username] [password]",
    Short: "Login and obtain tokens",
    Run: func(cmd *cobra.Command, args []string) {

		username := prompts.AskUsername("")
		password := prompts.AskPassword(false)

		err := apiClient.Login(models.LoginRequest{
			Username: username,
			Password: password,
		})
        if err != nil {
			prompts.ShowError(fmt.Sprintf("Login failed: %v", err))
            return
        }
		prompts.ShowSuccess("Loggin successful")
    },
}

var logoutCmd = &cobra.Command{
    Use:   "logout",
    Short: "Logout current session",
    Run: func(cmd *cobra.Command, args []string) {
        err := apiClient.Get("/logout", map[string]string{}, nil)
        if err != nil {
			prompts.ShowError("Logout error:" + err.Error())
        } else {
            prompts.ShowSuccess("Logged out")
        }
       	cfgManager.UpdateToken(config.UpdateToken{
			AccessToken: "",
			RefreshToken: "",
			SessionID: "",
			ExpiresAt: 0,
		}) 
    },
}

var refreshCmd = &cobra.Command{
    Use:   "refresh",
    Short: "Refresh access token",
    Run: func(cmd *cobra.Command, args []string) {
        err := apiClient.Refresh()
        if err != nil {
			prompts.ShowError("Refresh failed:" + err.Error())
            return
        }
		prompts.ShowSuccess("Token refreshed")
    },
}

func init() {
    authCmd.AddCommand(loginCmd, logoutCmd, refreshCmd)
}
