package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"rss-cli/pkg/database"
)

type OutputFormat bool

const (
	FormatText OutputFormat = false
	FormatJSON OutputFormat = true
)

type Printer struct {
	format OutputFormat
	writer io.Writer
}

func NewPrinter(jsonMode bool) *Printer {
	return &Printer{
		format: OutputFormat(jsonMode),
		writer: os.Stdout,
	}
}

func (p *Printer) Output(data interface{}) error {
	if p.format == FormatJSON {
		return p.outputJSON(data)
	}
	return p.outputText(data)
}

func (p *Printer) outputJSON(data interface{}) error {
	encoder := json.NewEncoder(p.writer)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(data)
}

func (p *Printer) outputText(data interface{}) error {
	switch v := data.(type) {
	case map[string]interface{}:
		return p.outputMap(v)
	case map[string]string:
		return p.outputStringMap(v)
	case *database.Feed:
		return p.outputFeed(v)
	case []database.Feed:
		return p.outputFeeds(v)
	case *database.Article:
		return p.outputArticle(v)
	case []database.Article:
		return p.outputArticles(v)
	default:
		return p.outputJSON(data)
	}
}

func (p *Printer) outputMap(m map[string]interface{}) error {
	if status, ok := m["status"]; ok {
		switch status {
		case "success":
			if msg, ok := m["msg"]; ok {
				fmt.Fprintf(p.writer, "Success: %s\n", msg)
				return nil
			}
			if msg, ok := m["message"]; ok {
				fmt.Fprintf(p.writer, "Success: %s\n", msg)
				return nil
			}
			if imported, ok := m["imported"]; ok {
				added, _ := m["added"]
				fmt.Fprintf(p.writer, "Imported: %d feeds, %d added\n", imported, added)
				return nil
			}
			if feed, ok := m["feed"]; ok {
				fmt.Fprintf(p.writer, "Feed added successfully\n")
				return p.outputText(feed)
			}
			fmt.Fprintln(p.writer, "Success")
			return nil
		}
	}

	if feeds, ok := m["feeds"].([]database.Feed); ok {
		return p.outputFeeds(feeds)
	}
	if articles, ok := m["articles"].([]database.Article); ok {
		return p.outputArticles(articles)
	}
	if results, ok := m["results"].([]map[string]interface{}); ok {
		return p.outputResults(results)
	}

	return p.outputJSON(m)
}

func (p *Printer) outputStringMap(m map[string]string) error {
	if status, ok := m["status"]; ok && status == "success" {
		if msg, ok := m["msg"]; ok {
			fmt.Fprintf(p.writer, "Success: %s\n", msg)
			return nil
		}
		if msg, ok := m["message"]; ok {
			fmt.Fprintf(p.writer, "Success: %s\n", msg)
			return nil
		}
		fmt.Fprintln(p.writer, "Success")
		return nil
	}
	return p.outputJSON(m)
}

func (p *Printer) outputFeed(feed *database.Feed) error {
	fmt.Fprintf(p.writer, "ID: %d\n", feed.ID)
	fmt.Fprintf(p.writer, "Title: %s\n", feed.Title)
	fmt.Fprintf(p.writer, "URL: %s\n", feed.URL)
	if feed.LastUpdated != nil {
		fmt.Fprintf(p.writer, "Last Updated: %s\n", *feed.LastUpdated)
	}
	fmt.Fprintf(p.writer, "Error Count: %d\n", feed.ErrorCount)
	return nil
}

func (p *Printer) outputFeeds(feeds []database.Feed) error {
	if len(feeds) == 0 {
		fmt.Fprintln(p.writer, "No feeds found.")
		return nil
	}

	w := tabwriter.NewWriter(p.writer, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTitle\tURL\tLast Updated\tErrors")

	for _, feed := range feeds {
		lastUpdated := "-"
		if feed.LastUpdated != nil {
			lastUpdated = *feed.LastUpdated
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\n",
			feed.ID,
			truncate(feed.Title, 30),
			truncate(feed.URL, 40),
			lastUpdated,
			feed.ErrorCount,
		)
	}
	return w.Flush()
}

func (p *Printer) outputArticle(article *database.Article) error {
	fmt.Fprintf(p.writer, "ID: %d\n", article.ID)
	fmt.Fprintf(p.writer, "Feed ID: %d\n", article.FeedID)
	fmt.Fprintf(p.writer, "Title: %s\n", article.Title)
	if article.Link != "" {
		fmt.Fprintf(p.writer, "Link: %s\n", article.Link)
	}
	fmt.Fprintf(p.writer, "Published: %s\n", article.PublishedAt)
	fmt.Fprintf(p.writer, "Read: %t\n", article.Read)
	if article.Content != "" {
		fmt.Fprintf(p.writer, "\n%s\n", article.Content)
	}
	return nil
}

func (p *Printer) outputArticles(articles []database.Article) error {
	if len(articles) == 0 {
		fmt.Fprintln(p.writer, "No articles found.")
		return nil
	}

	w := tabwriter.NewWriter(p.writer, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tFeed\tTitle\tPublished\tRead")

	for _, article := range articles {
		read := "no"
		if article.Read {
			read = "yes"
		}
		fmt.Fprintf(w, "%d\t%d\t%s\t%s\t%s\n",
			article.ID,
			article.FeedID,
			truncate(article.Title, 40),
			article.PublishedAt,
			read,
		)
	}
	return w.Flush()
}

func (p *Printer) outputResults(results []map[string]interface{}) error {
	for _, r := range results {
		feedID, _ := r["feed_id"].(int)
		title, _ := r["title"].(string)
		status, _ := r["status"].(string)

		if status == "error" {
			errMsg, _ := r["error"].(string)
			fmt.Fprintf(p.writer, "Feed %d (%s): ERROR - %s\n", feedID, title, errMsg)
		} else {
			added, _ := r["added_articles"].(int)
			fmt.Fprintf(p.writer, "Feed %d (%s): updated, %d articles added\n", feedID, title, added)
		}
	}
	return nil
}

func (p *Printer) Error(message string) error {
	if p.format == FormatJSON {
		return p.outputJSON(map[string]string{"error": message})
	}
	fmt.Fprintf(p.writer, "Error: %s\n", message)
	return nil
}

func truncate(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
