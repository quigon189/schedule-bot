package cmd

import (
    "log"
    "schedule-cli/internal/client"
    "schedule-cli/internal/config"

    "github.com/spf13/cobra"
)

var (
    cfgManager *config.Manager
    apiClient  *client.APIClient
    formatFlag string
)

var rootCmd = &cobra.Command{
    Use:   "schedule-cli",
    Short: "CLI client for schedule management service",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        var err error
        cfgManager, err = config.NewManager()
        if err != nil {
            log.Fatalf("config error: %v", err)
        }
        apiClient = client.NewAPIClient(cfgManager)
    },
}

func Execute() {
    rootCmd.PersistentFlags().StringVarP(&formatFlag, "output", "o", "table", "output format (table, json)")
    rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(configCmd)
    rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(adminCmd)
	rootCmd.AddCommand(chatCmd)
    // ... добавить остальные команды
    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}
