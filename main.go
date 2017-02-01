package main

import (
	"./server"
)

func main() {
	s := server.New("0.0.0.0", "8080")
	s.Execute()
}
