package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rss",
	Short: "A command line RSS reader with SQLite backend",
	Long: `A command line RSS reader with SQLite backend

A comprehensive RSS feed reader that stores feeds and articles in a SQLite database.
Supports feed management, article reading, and OPML import/export.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show comprehensive help when running just "rss" or "rss -h"
		showComprehensiveHelp(cmd)
	},
}

func showComprehensiveHelp(cmd *cobra.Command) {
	output := `
================================================================================
RSS CLI - A Command Line RSS Reader
================================================================================

A comprehensive RSS feed reader with SQLite backend for managing feeds and articles.

================================================================================
GLOBAL FLAGS
================================================================================
  -d, --db-path string   Database file path (default: ~/.rss-cli.db)
  -j, --json             Output in JSON format (default: true)
  -t, --text             Output in plain text format
  -h, --help             Show help for any command

================================================================================
COMMANDS
================================================================================

┌─────────────────────────────────────────────────────────────────────────────┐
│ FEED COMMANDS (rss feed)                                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│ rss feed add [url]                                                          │
│   Add a new RSS feed by URL                                                 │
│                                                                             │
│ rss feed list                                                               │
│   List all RSS feeds in the database                                        │
│                                                                             │
│ rss feed remove [id]                                                        │
│   Remove an RSS feed by its ID                                              │
│                                                                             │
│ rss feed update [id]                                                        │
│   Update a specific RSS feed and import new articles                        │
│                                                                             │
│ rss feed update-all                                                         │
│   Update all RSS feeds and import new articles                              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ ARTICLE COMMANDS (rss article)                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ rss article list [flags]                                                    │
│   List RSS articles with optional filtering                                 │
│   Flags:                                                                    │
│     --unread        Show only unread articles                               │
│     --read          Show only read articles                                 │
│     -f, --feed id   Filter by feed ID                                       │
│     -l, --limit n   Limit number of results                                 │
│                                                                             │
│ rss article mark [id] [read|unread]                                         │
│   Mark an article as read or unread                                         │
│   Arguments:                                                                │
│     id            Article ID                                                │
│     read|unread   State to mark the article                                 │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ IMPORT/EXPORT COMMANDS (rss import)                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│ rss import [opml-file]                                                      │
│   Import RSS feeds from an OPML file                                        │
│   Arguments:                                                                │
│     opml-file    Path to the OPML file to import                            │
│                                                                             │
│ rss import export [opml-file]                                               │
│   Export RSS feeds to an OPML file                                          │
│   Arguments:                                                                │
│     opml-file    Path to the OPML file to export to                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ UTILITY COMMANDS                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│ rss completion [bash|zsh|fish|powershell]                                   │
│   Generate autocompletion script for the specified shell                    │
│                                                                             │
│ rss help [command]                                                          │
│   Show help about any command                                               │
└─────────────────────────────────────────────────────────────────────────────┘

================================================================================
EXAMPLES
================================================================================

  # Add a new RSS feed
  rss feed add https://example.com/feed.xml

  # List all feeds
  rss feed list

  # Update all feeds
  rss feed update-all

  # List unread articles
  rss article list --unread

  # Mark article as read
  rss article mark 123 read

  # Import feeds from OPML
  rss import feeds.opml

  # Export feeds to OPML
  rss import export feeds.opml

================================================================================
Use "rss [command] --help" for more information about a specific command.
================================================================================
`
	fmt.Println(strings.TrimSpace(output))
}

// Define flags at root command level
func init() {
	// Override the help function to show comprehensive help only for root command
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd == rootCmd {
			showComprehensiveHelp(cmd)
		} else {
			// Use default help for subcommands
			cmd.Root().SetHelpFunc(nil)
			cmd.Help()
		}
	})

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
