package main

import (
	"fmt"
	"testBackend/internal/post"
	"testBackend/internal/serverapi"
)

func main() {
	fmt.Println("API Server started. Hello!")
	var a post.Sender
	a = &serverapi.Queue{}
	serverapi.SetQueue(a)
	serverapi.Run()
}
