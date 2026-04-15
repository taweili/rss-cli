package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"rss-cli/pkg/database"
	"rss-cli/pkg/rss"
	"rss-cli/pkg/ui"
)

var articleFetchCmd = &cobra.Command{
	Use:   "fetch [id]",
	Short: "Fetch full article content from source URL and convert to markdown",
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

		// Fetch article from database
		article, err := db.GetArticleByID(id)
		if err != nil {
			return printer.Error(fmt.Sprintf("Failed to retrieve article: %v", err))
		}

		if article.Link == "" {
			return printer.Error("Article has no URL to fetch")
		}

		// Fetch and convert article from source URL
		markdown, err := rss.FetchAndConvertArticle(article.Link)
		if err != nil {
			return printer.Error(fmt.Sprintf("Failed to fetch article content: %v", err))
		}

		// Cache the markdown content to database
		if err := db.UpdateArticleContent(id, markdown); err != nil {
			return printer.Error(fmt.Sprintf("Failed to cache article content: %v", err))
		}

		// Output the markdown content
		return printer.Output(map[string]interface{}{
			"status":  "success",
			"id":      article.ID,
			"title":   article.Title,
			"link":    article.Link,
			"content": markdown,
		})
	},
}
