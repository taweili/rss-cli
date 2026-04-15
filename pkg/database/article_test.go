package database

import (
	"os"
	"testing"
	"time"
)

// testingTBTime returns a fixed time for testing
func testingTBTime(t *testing.T) time.Time {
	t.Helper()
	return time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
}

func TestUpdateArticleContent(t *testing.T) {
	// Create a temporary database
	tmpfile, err := os.CreateTemp("", "test_*.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	db, err := NewDB(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Add a test article first
	feedID := 1
	guid := "test-guid"
	title := "Test Article"
	initialContent := "Initial content"
	link := "https://example.com/article"
	pubDate := testingTBTime(t)
	read := false

	err = db.AddArticle(feedID, guid, title, initialContent, link, pubDate, read)
	if err != nil {
		t.Fatal(err)
	}

	// Get the article to verify initial state
	article, err := db.GetArticleByID(1)
	if err != nil {
		t.Fatal(err)
	}

	if article.Content != initialContent {
		t.Errorf("Expected initial content %q, got %q", initialContent, article.Content)
	}

	// Update the content
	newContent := "Updated markdown content"
	err = db.UpdateArticleContent(1, newContent)
	if err != nil {
		t.Fatal(err)
	}

	// Verify the content was updated
	updatedArticle, err := db.GetArticleByID(1)
	if err != nil {
		t.Fatal(err)
	}

	if updatedArticle.Content != newContent {
		t.Errorf("Expected updated content %q, got %q", newContent, updatedArticle.Content)
	}
}

func TestUpdateArticleContent_NonExistentArticle(t *testing.T) {
	// Create a temporary database
	tmpfile, err := os.CreateTemp("", "test_*.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	db, err := NewDB(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Try to update a non-existent article
	err = db.UpdateArticleContent(999, "some content")
	if err != nil {
		t.Errorf("Expected no error for non-existent article, got %v", err)
	}
}
