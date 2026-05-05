package main

import (
	"log"
	"website-api/common"
	"website-api/database"
	contentPageRepo "website-api/repository/content-page"
	"website-api/service/content_page"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = db.SqlDb.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	contentPageService := content_page.NewService(contentPageRepo.NewRepo(db.GormDb))
	if err = contentPageService.Seed(); err != nil {
		log.Fatal(err)
	}
	log.Print(common.SuccessfullyCreated)
}
