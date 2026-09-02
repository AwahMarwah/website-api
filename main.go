package main

import (
	"log"
	"website-api/cache"
	"website-api/database"
	"website-api/router"

	_ "github.com/joho/godotenv/autoload"
)

// @title		Website Api
// @version		1.0
// @description	Api Documentation for Website
// @host		localhost:8085
// @BasePath	/
//
// @tag.name 1. Health Check
// @tag.name 2. Content Page
// @tag.name 3. Master
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
