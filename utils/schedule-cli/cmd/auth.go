package cmd

import (
	"fmt"
	"schedule-cli/internal/config"
	"schedule-cli/internal/models"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
    Use:   "auth",
    Short: "Authentication commands",
}

var loginCmd = &cobra.Command{
    Use:   "login [username] [password]",
    Short: "Login and obtain tokens",
    Args:  cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {

		err := apiClient.Login(models.LoginRequest{
			Username: args[0],
			Password: args[1],
		})
        if err != nil {
            fmt.Println("Login failed:", err)
            return
        }
        fmt.Println("Login successful")
    },
}

var logoutCmd = &cobra.Command{
    Use:   "logout",
    Short: "Logout current session",
    Run: func(cmd *cobra.Command, args []string) {
        err := apiClient.Get("/logout", map[string]string{}, nil)
        if err != nil {
            fmt.Println("Logout error:", err)
        } else {
            fmt.Println("Logged out")
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
            fmt.Println("Refresh failed:", err)
            return
        }
        fmt.Println("Token refreshed")
    },
}

func init() {
    authCmd.AddCommand(loginCmd, logoutCmd, refreshCmd)
}
