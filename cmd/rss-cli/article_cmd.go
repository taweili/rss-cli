package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/toqueteos/webbrowser"
	"rss-cli/pkg/database"
	"rss-cli/pkg/ui"
)

var articleCmd = &cobra.Command{
	Use:   "article",
	Short: "Manage RSS articles",
}

var articleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List RSS articles",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonMode, _ := cmd.Flags().GetBool("json")
		printer := ui.NewPrinter(jsonMode)

		dbPath, _ := cmd.Flags().GetString("db-path")
		db, err := database.NewDB(dbPath)
		if err != nil {
			return printer.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		}
		defer db.Close()

		unread, _ := cmd.Flags().GetBool("unread")
		read, _ := cmd.Flags().GetBool("read")
		feedIDStr, _ := cmd.Flags().GetString("feed")
		limitStr, _ := cmd.Flags().GetString("limit")

		filter := &database.ArticleFilter{}

		if feedIDStr != "" {
			feedID, err := strconv.Atoi(feedIDStr)
			if err != nil {
				return printer.Error("Invalid feed ID")
			}
			filter.FeedID = &feedID
		}

		if unread {
			status := false
			filter.Read = &status
		} else if read {
			status := true
			filter.Read = &status
		}

		if limitStr != "" {
			limit, err := strconv.Atoi(limitStr)
			if err != nil {
				return printer.Error("Invalid limit")
			}
			filter.Limit = &limit
		}

		articles, err := db.GetArticles(filter)
		if err != nil {
			return printer.Error(fmt.Sprintf("Failed to retrieve articles: %v", err))
		}

		return printer.Output(map[string]interface{}{
			"articles": articles,
			"count":    len(articles),
		})
	},
}

var articleMarkCmd = &cobra.Command{
	Use:   "mark [id] [read|unread]",
	Short: "Mark an article as read or unread",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonMode, _ := cmd.Flags().GetBool("json")
		printer := ui.NewPrinter(jsonMode)

		dbPath, _ := cmd.Flags().GetString("db-path")
		db, err := database.NewDB(dbPath)
		if err != nil {
			return printer.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		}
		defer db.Close()

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return printer.Error("Invalid article ID")
		}

		readState := strings.ToLower(args[1])
		var read bool

		switch readState {
		case "read":
			read = true
		case "unread":
			read = false
		default:
			return printer.Error("Invalid state, use 'read' or 'unread'")
		}

		if err := db.SetArticleReadStatus(id, read); err != nil {
			return printer.Error(fmt.Sprintf("Failed to update article: %v", err))
		}

		return printer.Output(map[string]string{
			"status": "success",
			"msg": fmt.Sprintf("Article %d marked as %s",
				id, map[bool]string{true: "read", false: "unread"}[read]),
		})
	},
}

var articleViewCmd = &cobra.Command{
	Use:   "view [id]",
	Short: "View an article by ID",
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

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return printer.Error("Invalid article ID")
		}

		article, err := db.GetArticleByID(id)
		if err != nil {
			return printer.Error(fmt.Sprintf("Failed to retrieve article: %v", err))
		}

		open, _ := cmd.Flags().GetBool("open")
		if open {
			if err := webbrowser.Open(article.Link); err != nil {
				return printer.Error(fmt.Sprintf("Failed to open article in browser: %v", err))
			}
		}

		return printer.Output(map[string]interface{}{
			"status": "success",
			"article": article,
		})
	},
}

func init() {
	articleCmd.AddCommand(articleListCmd)
	articleCmd.AddCommand(articleMarkCmd)
	articleCmd.AddCommand(articleViewCmd)

	// Flags for article list
	articleListCmd.Flags().Bool("unread", false, "Show only unread articles")
	articleListCmd.Flags().Bool("read", false, "Show only read articles")
	articleListCmd.Flags().StringP("feed", "f", "", "Filter by feed ID")
	articleListCmd.Flags().StringP("limit", "l", "", "Limit number of results")

	// Flags for article view
	articleViewCmd.Flags().Bool("open", false, "Open article URL in default browser")
}
