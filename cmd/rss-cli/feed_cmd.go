package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"rss-cli/pkg/database"
	"rss-cli/pkg/rss"
	"rss-cli/pkg/ui"
)

var feedCmd = &cobra.Command{
	Use:   "feed",
	Short: "Manage RSS feeds",
}

var feedAddCmd = &cobra.Command{
	Use:   "add [url]",
	Short: "Add a new RSS feed",
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

		feedURL := args[0]

		feedData, err := rss.FetchAndParseFeed(feedURL)
		if err != nil {
			return printer.Error(fmt.Sprintf("Failed to fetch feed: %v", err))
		}

		title := strings.TrimSpace(feedData.Title)
		if title == "" {
			// Fallback to hostname if feed has no title
			title = feedURL
			if parsedURL, err := url.Parse(feedURL); err == nil {
				title = parsedURL.Host
			}
		}

		if err := db.AddFeed(title, feedURL); err != nil {
			return printer.Error(fmt.Sprintf("Failed to add feed: %v", err))
		}

		return printer.Output(map[string]interface{}{
			"status": "success",
			"feed": map[string]string{
				"title": title,
				"url":   feedURL,
			},
		})
	},
}

var feedListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all RSS feeds",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonMode, _ := cmd.Flags().GetBool("json")
		printer := ui.NewPrinter(jsonMode)

		dbPath, _ := cmd.Flags().GetString("db-path")
		db, err := database.NewDB(dbPath)
		if err != nil {
			return printer.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		}
		defer db.Close()

		feeds, err := db.GetAllFeeds()
		if err != nil {
			return printer.Error(fmt.Sprintf("Failed to retrieve feeds: %v", err))
		}

		return printer.Output(map[string]interface{}{
			"feeds": feeds,
			"count": len(feeds),
		})
	},
}

var feedRemoveCmd = &cobra.Command{
	Use:   "remove [id]",
	Short: "Remove an RSS feed by ID",
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
			return printer.Error("Invalid feed ID")
		}

		if err := db.DeleteFeed(id); err != nil {
			return printer.Error(fmt.Sprintf("Failed to remove feed: %v", err))
		}

		return printer.Output(map[string]string{
			"status": "success",
			"msg":    fmt.Sprintf("Removed feed ID %d", id),
		})
	},
}

var feedUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update an RSS feed by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateSingleFeed(cmd, args[0])
	},
}

var feedUpdateAllCmd = &cobra.Command{
	Use:   "update-all",
	Short: "Update all RSS feeds",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonMode, _ := cmd.Flags().GetBool("json")
		printer := ui.NewPrinter(jsonMode)

		dbPath, _ := cmd.Flags().GetString("db-path")
		db, err := database.NewDB(dbPath)
		if err != nil {
			return printer.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		}
		defer db.Close()

		feeds, err := db.GetAllFeeds()
		if err != nil {
			return printer.Error(fmt.Sprintf("Failed to retrieve feeds: %v", err))
		}

		results := make([]map[string]interface{}, 0)
		for _, feed := range feeds {
			result := updateFeed(db, &feed, jsonMode)
			results = append(results, result)
		}

		return printer.Output(map[string]interface{}{
			"results": results,
			"updated": len(results),
		})
	},
}

func init() {
	feedCmd.AddCommand(feedAddCmd)
	feedCmd.AddCommand(feedListCmd)
	feedCmd.AddCommand(feedRemoveCmd)
	feedCmd.AddCommand(feedUpdateCmd)
	feedCmd.AddCommand(feedUpdateAllCmd)
}

func updateSingleFeed(cmd *cobra.Command, idStr string) error {
	jsonMode, _ := cmd.Flags().GetBool("json")
	printer := ui.NewPrinter(jsonMode)

	dbPath, _ := cmd.Flags().GetString("db-path")
	db, err := database.NewDB(dbPath)
	if err != nil {
		return printer.Error(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer db.Close()

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return printer.Error("Invalid feed ID")
	}

	feed, err := db.GetFeedByID(id)
	if err != nil {
		return printer.Error("Feed not found")
	}

	result := updateFeed(db, feed, jsonMode)

	return printer.Output(result)
}

func updateFeed(db *database.DB, feed *database.Feed, jsonMode bool) map[string]interface{} {
	result := map[string]interface{}{
		"feed_id": feed.ID,
		"title":   feed.Title,
	}

	// Actually retrieve and parse the feed
	parsedFeed, err := rss.FetchAndParseFeed(feed.URL)
	if err != nil {
		db.IncrementErrorCount(feed.ID)
		result["status"] = "error"
		result["error"] = err.Error()
		return result
	}

	count := 0
	for _, item := range parsedFeed.Items {
		// Use Published date since UpdatedAt doesn't exist in gofeed Item struct
		pubTime := item.Published

		// If we couldn't get published date from feed, use current time
		timeToUse := time.Now()
		if pubTime != "" {
			if parsedTime, err := parseTime(pubTime); err == nil {
				timeToUse = parsedTime
			}
		}

		// Use GUID or Link as the unique identifier
		guid := item.GUID
		if guid == "" {
			guid = item.Link
		}

		if err := db.AddArticle(
			feed.ID,
			guid,
			strings.TrimSpace(item.Title),
			item.Description, // Using description as content
			item.Link,
			timeToUse, // Using proper time.Time
			false,     // New articles are initially unread
		); err != nil {
			// Still continue adding other articles
			continue
		}
		count++
	}

	// Update the feed's last updated timestamp
	if count > 0 {
		db.UpdateFeedTimestamp(feed.ID, time.Now().Format(time.RFC3339))
	}

	result["status"] = "success"
	result["added_articles"] = count

	return result
}

func parseTime(pubTimeString string) (time.Time, error) {
	// RFC3339 is the most common format
	if t, err := time.Parse(time.RFC3339, pubTimeString); err == nil {
		return t, nil
	}

	// Try some other common formats
	formats := []string{
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		"Mon Jan 2 15:04:05 -0700 2006",
		"Jan 2, 2006 3:04 PM MST",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02 Jan 2006 15:04:05 GMT",
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"02 Jan 2006 15:04:05 MST",
		"January 2, 2006",
		"Jan 2, 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, pubTimeString); err == nil {
			return t, nil
		}
	}

	// If nothing matches, fallback to current time
	return time.Now(), nil
}
