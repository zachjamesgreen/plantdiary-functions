package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type Payload struct {
	Slug string
}

type Post struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	Slug       string    `json:"slug"`
	Url        string    `json:"url"`
	CoverImage string    `json:"cover_image"`
	UpdatedAt  time.Time `json:"updated_at"`
	Published  bool      `json:"published"`
}

func Main(json Payload) map[string]interface{} {
	fmt.Println("Starting getPost")
	var posts []Post
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

	if err = pgxscan.Select(ctx, conn, &posts, `SELECT id,title,body,slug,url,cover_image,updated_at,published from post`); err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed for <%v>: %v\n", json.Slug, err)
		response["body"] = fmt.Sprintf("QueryRow failed for <%v>: %v\n", json.Slug, err)
		response["statusCode"] = 500
		return response
	}

	response["body"] = posts
	return response
}
