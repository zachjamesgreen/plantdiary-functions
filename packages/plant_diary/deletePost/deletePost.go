package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

type Payload struct {
	ID string
}

func Main(json Payload) map[string]interface{} {
	fmt.Println("Starting deletePost")
	ctx := context.Background()
	response := make(map[string]interface{})
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		response["body"] = fmt.Sprintf("Unable to connect to database: %v\n", err)
		response["statusCode"] = 500
		return response
	}
	defer conn.Close(ctx)

	command, err := conn.Exec(ctx, "DELETE from post where id = $1", json.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Delete failed for <%v>: %v\n", json.ID, err)
		response["body"] = fmt.Sprintf("Delete failed for <%v>: %v\n", json.ID, err)
		response["statusCode"] = 500
		return response
	}
	if command.RowsAffected() != 1 {
		fmt.Fprintf(os.Stderr, "No row found to delete for <%v>: %v\n", json.ID, err)
		response["body"] = fmt.Sprintf("No row found to delete for <%v>: %v\n", json.ID, err)
		response["statusCode"] = 404
		return response
	}

	response["statusCode"] = 204
	return response
}
