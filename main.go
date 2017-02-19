package main

import (
	"./server"
	"log"
)

// Test
func main() {
	// TODO: Change to env. vars.
	s, err := server.New("0.0.0.0", "8080", "unix:///var/run/docker.sock")
	if err != nil {
		log.Fatal(err.Error())
	}
	s.Cron.RescheduleJobs()
	s.Execute()
}
