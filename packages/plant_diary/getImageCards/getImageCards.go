package main

import (
	"context"
	"fmt"
	"os"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type Post struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	PublishedAt bool   `json:"published_at" db:"published_at"`
	Url         string `json:"url"`
	CoverImage  string `json:"cover_image" `
}

func Main(args map[string]interface{}) map[string]interface{} {
	fmt.Println("Starting getImageCards")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	query := "SELECT id,title,slug,published_at,url,cover_image FROM post WHERE published = true ORDER BY id DESC"
	rows, err := conn.Query(ctx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	var posts []Post
	if err := pgxscan.ScanAll(&posts, rows); err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	msg := make(map[string]interface{})
	msg["body"] = posts
	return msg
}
