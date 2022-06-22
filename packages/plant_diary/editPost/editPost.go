package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type RequestPost struct {
	ID         int       `json:"id"`
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

func Main(r_post RequestPost) map[string]interface{} {
	fmt.Println("Starting editPost")
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

	if r_post.ID == 0 {
		fmt.Println("Creating new post")
		query := "insert into post (title,body,slug,url,cover_image,published) values ($1,$2,$3,$4,$5,$6) returning id"
		err := conn.QueryRow(ctx, query, r_post.Title, r_post.Body, r_post.Slug, r_post.Url, r_post.CoverImage, r_post.Published).Scan(&r_post.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Insert failed: %v\n", err)
			response["body"] = fmt.Sprintf("Insert failed: %v\n", err)
			response["statusCode"] = 500
			return response
			// os.Exit(1)
		}
	} else {
		insertQuery := "UPDATE post set title = $1, body = $2, slug = $3, url = $4, cover_image = $5, updated_at = $6, published = $7 WHERE id = $8"
		_, err = conn.Exec(ctx, insertQuery, r_post.Title, r_post.Body, r_post.Slug, r_post.Url, r_post.CoverImage, time.Now(), r_post.Published, r_post.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			response["body"] = fmt.Sprintf("QueryRow failed: %v\n", err)
			response["statusCode"] = 500
			return response
		}
	}

	var post Post
	selectQuery := "select id,title,body,slug,url,cover_image,updated_at,published from post where id = $1"
	err = pgxscan.Get(ctx, conn, &post, selectQuery, r_post.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		response["body"] = fmt.Sprintf("QueryRow failed: %v\n", err)
		response["statusCode"] = 500
		return response
	}

	response["body"] = post
	return response
}
