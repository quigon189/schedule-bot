package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Interactive chat with LLM assistant",
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("Starting interactive chat. Type 'exit' to quit.\n")
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print(color.CyanString(" >"))
			input, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			input = strings.TrimSpace(input)
			if input == "exit" {
				break
			}
			if input == "" {
				continue
			}
			var resp map[string]any
			err = apiClient.Post("/chat", map[string]any{"message": input}, &resp)
			if err != nil {
				color.Red("Error: %v", err)
				continue
			}

			reply, ok := resp["reply"].(string)
			if ok {
				fmt.Println(reply)
			} else {
				color.Red("Reply: %v", resp)
			}
		}
	},
}
