package main

import "github.com/alextonkonogov/gb-go-backend-1/lesson-2/task-1/internal/client"

func main() {
	c := client.NewClient()
	err := c.Start()
	if err != nil {

	}
}
