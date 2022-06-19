package main

import (
	"fmt"
)

func Main(post map[string]interface{}) map[string]interface{} {
	fmt.Println("Starting savePost")

	msg := make(map[string]interface{})
	msg["body"] = "savePost"
	return msg
}
