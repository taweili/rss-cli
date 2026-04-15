package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
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

var articleOpenCmd = &cobra.Command{
	Use:   "open [id]",
	Short: "Open an article URL in the default browser",
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

		if article.Link == "" {
			return printer.Error("Article has no URL")
		}

		// Validate URL format
		if _, err := url.ParseRequestURI(article.Link); err != nil {
			return printer.Error(fmt.Sprintf("Invalid article URL: %v", err))
		}

		// Check if custom browser is specified
		browser, _ := cmd.Flags().GetString("browser")
		if browser != "" {
			// Use custom browser
			if err := openBrowserWithCustom(article.Link, browser); err != nil {
				return printer.Error(fmt.Sprintf("Failed to open article in %s: %v", browser, err))
			}
		} else {
			// Use default browser detection
			if err := openBrowser(article.Link); err != nil {
				return printer.Error(fmt.Sprintf("Failed to open article: %v", err))
			}
		}

		return printer.Output(map[string]interface{}{
			"status": "success",
			"msg":    fmt.Sprintf("Opened article %d in browser", id),
			"url":    article.Link,
		})
	},
}

// detectBrowser returns the platform-specific browser command
func detectBrowser() (string, error) {
	switch {
	case isLinux():
		return "xdg-open", nil
	case isMacOS():
		return "open", nil
	case isWindows():
		return "cmd", nil
	default:
		return "", fmt.Errorf("unsupported platform")
	}
}

// isLinux checks if the current OS is Linux
func isLinux() bool {
	return os.Getenv("XDG_CURRENT_DESKTOP") != "" || os.Getenv("DISPLAY") != ""
}

// isMacOS checks if the current OS is macOS
func isMacOS() bool {
	_, err := exec.LookPath("open")
	return err == nil
}

// isWindows checks if the current OS is Windows
func isWindows() bool {
	return os.Getenv("OS") == "Windows_NT"
}

// openBrowser opens a URL in the default browser using platform-specific commands
func openBrowser(url string) error {
	// First try the webbrowser package (cross-platform)
	if err := webbrowser.Open(url); err == nil {
		return nil
	}

	// Fallback to platform-specific commands
	browserCmd, err := detectBrowser()
	if err != nil {
		return fmt.Errorf("no browser found: %w", err)
	}

	var cmd *exec.Cmd
	switch browserCmd {
	case "xdg-open":
		cmd = exec.Command("xdg-open", url)
	case "open":
		cmd = exec.Command("open", url)
	case "cmd":
		cmd = exec.Command("cmd", "/c", "start", url)
	}

	if cmd != nil {
		// Detach the process so it doesn't block
		cmd.Stdout = nil
		cmd.Stderr = nil
		return cmd.Start()
	}

	return fmt.Errorf("unable to launch browser")
}

// openBrowserWithCustom opens a URL using a custom browser command
func openBrowserWithCustom(url string, browser string) error {
	cmd := exec.Command(browser, url)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Start()
}

func init() {
	articleCmd.AddCommand(articleListCmd)
	articleCmd.AddCommand(articleMarkCmd)
	articleCmd.AddCommand(articleViewCmd)
	articleCmd.AddCommand(articleOpenCmd)
	articleCmd.AddCommand(articleFetchCmd)

	// Flags for article list
	articleListCmd.Flags().Bool("unread", false, "Show only unread articles")
	articleListCmd.Flags().Bool("read", false, "Show only read articles")
	articleListCmd.Flags().StringP("feed", "f", "", "Filter by feed ID")
	articleListCmd.Flags().StringP("limit", "l", "", "Limit number of results")

	// Flags for article view
	articleViewCmd.Flags().Bool("open", false, "Open article URL in default browser")

	// Flags for article open
	articleOpenCmd.Flags().String("browser", "", "Custom browser command to use (e.g., 'firefox', 'chrome')")
}
