package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AaronRoethe/go-journal-client/message"
	"github.com/AaronRoethe/go-journal-client/pocket"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-journal",
	Short: "A brief description of mycli",
	Long:  "A longer description of mycli",
	Run: func(cmd *cobra.Command, args []string) {
		message.Journal()
	},
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	pocket.Auth_refresh()
	rootCmd.AddCommand(pocket.LoginCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
