package cmd

import (
	"fmt"
	"schedule-cli/internal/prompts"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure base parameters",
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Set base url for server",
	Run: func(cmd *cobra.Command, args []string) {
		baseURL := prompts.AskString("Base URL:", true, func(val any) error { return nil })
		timeout := prompts.AskInt("Timeout (sec):", 30, 100000)
		if err := cfgManager.SetBaseURL(baseURL); err != nil {
			msg := fmt.Sprintf("Failed to set base url: %v\n", err)
			prompts.ShowError(msg)
			return
		}
		if err := cfgManager.SetTimeout(timeout); err != nil {
			msg := fmt.Sprintf("Failed to set timeout: %v\n", err)
			prompts.ShowError(msg)
			return
		}
		msg := fmt.Sprintf("Set server config:\n")
		msg += fmt.Sprintf("    Base URL: %s\n", baseURL)
		msg += fmt.Sprintf("    Timeout: %d\n", timeout)
		prompts.ShowSuccess(msg, nil)
	},
}

func init() {
	configCmd.AddCommand(serverCmd)
}
