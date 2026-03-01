package database

type Feed struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	LastUpdated *string `json:"last_updated,omitempty"` // Store as string in ISO format
	ErrorCount  int     `json:"error_count"`
}

func (d *DB) AddFeed(title, url string) error {
	query := `INSERT INTO feeds (title, url) VALUES (?, ?)`
	_, err := d.Exec(query, title, url)
	return err
}

func (d *DB) GetAllFeeds() ([]Feed, error) {
	rows, err := d.Query(`SELECT id, title, url, last_updated, error_count FROM feeds`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []Feed
	for rows.Next() {
		var f Feed
		err := rows.Scan(&f.ID, &f.Title, &f.URL, &f.LastUpdated, &f.ErrorCount)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, f)
	}

	return feeds, nil
}

func (d *DB) GetFeedByID(id int) (*Feed, error) {
	var f Feed
	row := d.QueryRow(`SELECT id, title, url, last_updated, error_count FROM feeds WHERE id = ?`, id)
	err := row.Scan(&f.ID, &f.Title, &f.URL, &f.LastUpdated, &f.ErrorCount)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (d *DB) DeleteFeed(id int) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete associated articles first (due to foreign key constraint)
	_, err = tx.Exec(`DELETE FROM articles WHERE feed_id = ?`, id)
	if err != nil {
		return err
	}

	// Delete the feed
	_, err = tx.Exec(`DELETE FROM feeds WHERE id = ?`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Update feed's last updated timestamp
func (d *DB) UpdateFeedTimestamp(id int, timestamp string) error {
	query := `UPDATE feeds SET last_updated = ? WHERE id = ?`
	_, err := d.Exec(query, timestamp, id)
	return err
}

// Error handling for feeds
func (d *DB) IncrementErrorCount(id int) error {
	query := `UPDATE feeds SET error_count = error_count + 1 WHERE id = ?`
	_, err := d.Exec(query, id)
	return err
}
