package main

import (
	"log"

	"github.com/alextonkonogov/gb-go-backend-1/lesson-2/task-2/internal/server"
)

func main() {
	s := server.NewServer()
	err := s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
