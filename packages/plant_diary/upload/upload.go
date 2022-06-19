package main

import (
	"fmt"
)

func Main(args map[string]interface{}) map[string]interface{} {
	fmt.Println("Starting Function")
	fmt.Println(args)

	msg := make(map[string]interface{})
	msg["body"] = "got here"
	return msg
}
