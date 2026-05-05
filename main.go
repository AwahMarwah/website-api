package main

import (
	"log"
	"website-api/cache"
	"website-api/database"
	"website-api/router"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SqlDb.Close()

	// REDIS
	redisClient := cache.NewRedis()

	if err = router.Run(db, redisClient); err != nil {
		log.Fatal(err)
	}

}
