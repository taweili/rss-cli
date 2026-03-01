package database

import (
	"time"
)

type Article struct {
	ID          int    `json:"id"`
	FeedID      int    `json:"feed_id"`
	GUID        string `json:"guid"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Link        string `json:"link"`
	PublishedAt string `json:"published_at"`
	Read        bool   `json:"read"`
}

type ArticleFilter struct {
	FeedID *int
	Read   *bool
	Limit  *int
	SortBy string // "published_at" or "id"
	Order  string // "ASC" or "DESC"
}

func (d *DB) AddArticle(feedID int, guid, title, content, link string, pubDate time.Time, read bool) error {
	query := `
		INSERT INTO articles (feed_id, guid, title, content, link, published_at, read)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (feed_id, guid) 
		DO UPDATE SET title=excluded.title, content=excluded.content, 
		              link=excluded.link, published_at=excluded.published_at
	`
	_, err := d.Exec(query, feedID, guid, title, content, link, pubDate.Format(time.RFC3339), read)
	return err
}

func (d *DB) GetArticles(filter *ArticleFilter) ([]Article, error) {
	query := `SELECT id, feed_id, guid, title, content, link, published_at, read
	          FROM articles
	          WHERE 1=1 `

	args := []interface{}{}

	if filter.FeedID != nil {
		query += ` AND feed_id = ?`
		args = append(args, *filter.FeedID)
	}

	if filter.Read != nil {
		query += ` AND read = ?`
		args = append(args, *filter.Read)
	}

	// Default sort by published_at DESC
	sortBy := "published_at"
	order := "DESC"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	if filter.Order != "" {
		order = filter.Order
	}

	query += ` ORDER BY ` + sortBy + ` ` + order

	if filter.Limit != nil {
		query += ` LIMIT ?`
		args = append(args, *filter.Limit)
	}

	rows, err := d.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var art Article
		var publishedAtTemp string
		err := rows.Scan(
			&art.ID, &art.FeedID, &art.GUID, &art.Title,
			&art.Content, &art.Link, &publishedAtTemp, &art.Read,
		)
		if err != nil {
			return nil, err
		}

		// Parse the timestamp back to time.Time then convert to string
		// (could be simplified to store timestamp strings throughout but staying consistent with SQL)
		t, _ := time.Parse(time.RFC3339, publishedAtTemp)
		art.PublishedAt = t.Format(time.RFC3339)

		articles = append(articles, art)
	}

	return articles, nil
}

func (d *DB) SetArticleReadStatus(id int, read bool) error {
	query := `UPDATE articles SET read = ? WHERE id = ?`
	_, err := d.Exec(query, read, id)
	return err
}
