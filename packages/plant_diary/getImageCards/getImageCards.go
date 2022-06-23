package main

import (
	"context"
	"fmt"
	"os"

	// "time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type Post struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	Slug        string      `json:"slug"`
	PublishedAt pgtype.Time `json:"published_at" db:"published_at"`
	Url         string      `json:"url"`
	CoverImage  string      `json:"cover_image" `
}

func Main(args map[string]interface{}) map[string]interface{} {
	fmt.Println("Starting getImageCards")
	response := make(map[string]interface{})
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		response["body"] = fmt.Sprintf("Unable to connect to database: %v\n", err)
		response["statusCode"] = 500
		return response
	}
	defer conn.Close(ctx)

	query := "SELECT id,title,slug,published_at,url,cover_image FROM post WHERE published = true ORDER BY id DESC"
	rows, err := conn.Query(ctx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		response["body"] = fmt.Sprintf("QueryRow failed: %v\n", err)
		response["statusCode"] = 500
		return response
	}
	var posts []Post
	if err := pgxscan.ScanAll(&posts, rows); err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		response["body"] = fmt.Sprintf("QueryRow failed: %v\n", err)
		response["statusCode"] = 500
		return response
	}

	response["body"] = posts
	return response
}
