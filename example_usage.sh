#!/bin/bash

# Example usage of the RSS CLI tool
echo "Testing RSS CLI functionality..."

# List current feeds (should be empty)
echo "--- Current feeds ---"
go run ./cmd/rss-cli feed list

# Try to add a test feed (using Hacker News)
echo -e "\n--- Adding Hacker News feed ---"
go run ./cmd/rss-cli feed add "https://news.ycombinator.com/rss"

# List feeds after addition
echo -e "\n--- Feeds after adding ---"
go run ./cmd/rss-cli feed list

# Update the feed to get articles
echo -e "\n--- Updating feeds to fetch articles ---"
go run ./cmd/rss-cli feed update-all

# List articles
echo -e "\n--- Articles ---"
go run ./cmd/rss-cli article list

# Show just a few latest articles 
echo -e "\n--- Latest 3 articles ---"
go run ./cmd/rss-cli article list -l 3