package main

import (
	"./server"
)

func main() {
	s := server.NewServer()
	s.ListenAndServe()
}
