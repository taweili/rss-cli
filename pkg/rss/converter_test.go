package rss

import (
	"testing"

	"github.com/kreuzberg-dev/html-to-markdown/packages/go/v2/htmltomarkdown"
)

func TestHTMLToMarkdownConversion(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantMD   string
		wantErr  bool
	}{
		{
			name:    "simple heading and paragraph",
			html:    "<h1>Hello World</h1><p>This is a paragraph.</p>",
			wantMD:  "# Hello World\n\nThis is a paragraph.\n",
			wantErr: false,
		},
		{
			name:    "bold and italic text",
			html:    "<p>This is <strong>bold</strong> and <em>italic</em>.</p>",
			wantMD:  "This is **bold** and *italic*.\n",
			wantErr: false,
		},
		{
			name:    "links",
			html:    "<p>Visit <a href=\"https://example.com\">Example</a></p>",
			wantMD:  "Visit [Example](https://example.com)\n",
			wantErr: false,
		},
		{
			name:    "lists",
			html:    "<ul><li>Item 1</li><li>Item 2</li></ul>",
			wantMD:  "- Item 1\n- Item 2\n",
			wantErr: false,
		},
		{
			name:    "empty html",
			html:    "",
			wantMD:  "",
			wantErr: false,
		},
		{
			name:    "nested structure",
			html:    "<article><h1>Title</h1><p>Intro</p><h2>Section</h2><p>Content</p></article>",
			wantMD:  "# Title\n\nIntro\n\n## Section\n\nContent\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md, err := htmltomarkdown.Convert(tt.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if md != tt.wantMD {
				t.Errorf("Convert() = %q, want %q", md, tt.wantMD)
			}
		})
	}
}

func TestFetchArticleContent_InvalidURL(t *testing.T) {
	_, err := FetchArticleContent("not-a-valid-url")
	if err == nil {
		t.Error("FetchArticleContent() expected error for invalid URL, got nil")
	}
}

func TestFetchArticleContent_404(t *testing.T) {
	_, err := FetchArticleContent("https://httpbin.org/status/404")
	if err == nil {
		t.Error("FetchArticleContent() expected error for 404, got nil")
	}

	if httpErr, ok := err.(*HTTPError); ok {
		if httpErr.StatusCode != 404 {
			t.Errorf("Expected status code 404, got %d", httpErr.StatusCode)
		}
	} else {
		t.Error("Expected HTTPError type")
	}
}
