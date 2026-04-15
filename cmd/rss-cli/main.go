package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rss-cli",
	Short: "A command line RSS reader with SQLite backend",
	Long: `A command line RSS reader with SQLite backend.

Manages RSS feeds and articles in a local SQLite database.
Supports feed management, article reading, and OPML import/export.`,
}

// Define flags at root command level
func init() {
	// Define the default DB path
	homeDir, _ := os.UserHomeDir()
	defaultDbPath := filepath.Join(homeDir, ".rss-cli.db")

	rootCmd.PersistentFlags().StringP("db-path", "d", defaultDbPath, "Database file path")
	rootCmd.PersistentFlags().BoolP("json", "j", true, "Output in JSON format (default)")
	rootCmd.PersistentFlags().BoolP("text", "t", false, "Output in plain text format")

	// Expand tilde in path and handle output format flags
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		dbPath, _ := rootCmd.PersistentFlags().GetString("db-path")
		if dbPath != "" && dbPath[0:1] == "~" {
			home, _ := os.UserHomeDir()
			dbPath = filepath.Join(home, dbPath[1:])
			rootCmd.PersistentFlags().Set("db-path", dbPath)
		}

		// --text overrides --json
		textMode, _ := rootCmd.PersistentFlags().GetBool("text")
		if textMode {
			rootCmd.PersistentFlags().Set("json", "false")
		}
	}

	// Initialize commands (they are accessible globally as variables exported from other files)
	rootCmd.AddCommand(feedCmd)
	rootCmd.AddCommand(articleCmd)
	rootCmd.AddCommand(importCmd)

	// Export is a subcommand of import
	importCmd.AddCommand(exportCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
