package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type RequestPost struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	Slug       string    `json:"slug"`
	Url        string    `json:"url"`
	CoverImage string    `json:"cover_image"`
	UpdatedAt  time.Time `json:"updated_at"`
	Published  bool      `json:"published"`
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

func Main(data map[string]interface{}) map[string]Post {
	fmt.Println("Starting editPost")
	body := data["__ow_body"].(string)
	var r_post RequestPost
	err := json.Unmarshal([]byte(body), &r_post)
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	if r_post.ID == "" {
		fmt.Println("Creating new post")
		query := "insert into post (title,body,slug,url,cover_image,published) values ($1,$2,$3,$4,$5,$6) returning id"
		var id int
		err := conn.QueryRow(ctx, query, r_post.Title, r_post.Body, r_post.Slug, r_post.Url, r_post.CoverImage, r_post.Published).Scan(&id)
		if err != nil {
			fmt.Println(err)
		}
		r_post.ID = fmt.Sprintf("%d", id)
	} else {
		insertQuery := "UPDATE post set title = $1, body = $2, slug = $3, url = $4, cover_image = $5, updated_at = $6, published = $7 WHERE id = $8"
		_, err = conn.Exec(ctx, insertQuery, r_post.Title, r_post.Body, r_post.Slug, r_post.Url, r_post.CoverImage, time.Now(), r_post.Published, r_post.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
	}

	var post Post
	response := map[string]Post{}
	selectQuery := "select id,title,body,slug,url,cover_image,updated_at,published from post where id = $1"
	err = pgxscan.Get(ctx, conn, &post, selectQuery, r_post.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	response["body"] = post
	return response
}
