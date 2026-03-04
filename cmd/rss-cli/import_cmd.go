package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"rss-cli/pkg/database"
	"rss-cli/pkg/opml"
	"rss-cli/pkg/ui"
)

var importCmd = &cobra.Command{
	Use:   "import [opml-file]",
	Short: "Import RSS feeds from OPML file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonMode, _ := cmd.Flags().GetBool("json")
		printer := ui.NewPrinter(jsonMode)

		dbPath, _ := cmd.Flags().GetString("db-path")
		db, err := database.NewDB(dbPath)
		if err != nil {
			return printer.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		}
		defer db.Close()

		opmlPath := args[0]

		feedURLs, err := opml.Import(opmlPath)
		if err != nil {
			return printer.Error(fmt.Sprintf("Failed to import OPML: %v", err))
		}

		added := 0
		for _, url := range feedURLs {
			// Try to add the feed to the database
			// For simplicity, we'll use the URL as title for now since we're not
			// fetching details at import time
			if err := db.AddFeed(url, url); err != nil {
				// Continue adding other feeds even if one fails
				continue
			}
			added++
		}

		return printer.Output(map[string]interface{}{
			"status":   "success",
			"imported": len(feedURLs),
			"added":    added,
		})
	},
}

var exportCmd = &cobra.Command{
	Use:   "export [opml-file]",
	Short: "Export RSS feeds to OPML file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonMode, _ := cmd.Flags().GetBool("json")
		printer := ui.NewPrinter(jsonMode)

		dbPath, _ := cmd.Flags().GetString("db-path")
		db, err := database.NewDB(dbPath)
		if err != nil {
			return printer.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		}
		defer db.Close()

		// This is a placeholder implementation;
		// a real export implementation would need the file writing function

		// For now, the export command is just to show it's available
		// Actual implementation would require importing encoding/xml and adding file I/O
		return printer.Error("Export command not implemented in this version")
	},
}

func init() {
	// db-path flag is inherited from root command's PersistentFlags
}
