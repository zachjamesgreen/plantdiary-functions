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
	ID string
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

func Main(json Payload) map[string]Post {
	fmt.Println("Starting getPostById")
	var post Post
	response := map[string]Post{"body": post}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	query := "select id,title,body,slug,url,cover_image,updated_at,published from post where id = $1"
	if err = pgxscan.Get(ctx, conn, &post, query, json.ID); err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	response["body"] = post
	return response
}
