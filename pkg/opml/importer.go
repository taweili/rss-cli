package opml

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

type Head struct {
	Title       string `xml:"title"`
	DateCreated string `xml:"dateCreated"`
}

type Body struct {
	Outlines []Outline `xml:"outline"`
}

type Outline struct {
	Text        string    `xml:"text,attr"`
	Title       string    `xml:"title,attr"`
	Type        string    `xml:"type,attr"`
	XMLURL      string    `xml:"xmlUrl,attr"`
	HTMLURL     string    `xml:"htmlUrl,attr"`
	Description string    `xml:"description,attr"`
	Outlines    []Outline `xml:"outline"`
}

// Import Feeds from OPML file
func Import(opmlPath string) ([]string, error) {
	file, err := os.Open(opmlPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var opml OPML
	err = xml.Unmarshal(bytes, &opml)
	if err != nil {
		return nil, err
	}

	var urls []string
	extractFeedURLs(opml.Body.Outlines, &urls)
	return urls, nil
}

func extractFeedURLs(outlines []Outline, urls *[]string) {
	for _, outline := range outlines {
		// If this outline has an XMLURL, it's a feed
		if outline.XMLURL != "" {
			*urls = append(*urls, outline.XMLURL)
		}
		// If this outline has nested outlines, recursively process them
		if len(outline.Outlines) > 0 {
			extractFeedURLs(outline.Outlines, urls)
		}
	}
}
