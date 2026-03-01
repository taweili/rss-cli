package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

// Initialize DB with required tables
func NewDB(dbPath string) (*DB, error) {
	dir := filepath.Dir(dbPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func createTables(db *sql.DB) error {
	// Create feeds table
	feedTable := `
	CREATE TABLE IF NOT EXISTS feeds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		url TEXT UNIQUE NOT NULL,
		last_updated TIMESTAMP,
		error_count INTEGER DEFAULT 0
	);
	`

	_, err := db.Exec(feedTable)
	if err != nil {
		return err
	}

	// Create articles table
	articleTable := `
	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		feed_id INTEGER NOT NULL,
		guid TEXT NOT NULL,
		title TEXT NOT NULL,
		content TEXT,
		link TEXT,
		published_at TIMESTAMP,
		read BOOLEAN DEFAULT 0,
		FOREIGN KEY (feed_id) REFERENCES feeds (id),
		UNIQUE (feed_id, guid)
	);
	`

	_, err = db.Exec(articleTable)
	if err != nil {
		return err
	}

	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_articles_feed_id ON articles (feed_id);",
		"CREATE INDEX IF NOT EXISTS idx_articles_read_status ON articles (read);",
		"CREATE INDEX IF NOT EXISTS idx_articles_published_at ON articles (published_at DESC);",
		"CREATE INDEX IF NOT EXISTS idx_feeds_url ON feeds (url);",
	}

	for _, index := range indexes {
		if _, err := db.Exec(index); err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) Close() error {
	return d.DB.Close()
}
