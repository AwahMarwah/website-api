package main

import (
	"github.com/joho/godotenv"
	"log"
	"website-api/worker"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}
}

func main() {
	worker.StartWorker()
}
