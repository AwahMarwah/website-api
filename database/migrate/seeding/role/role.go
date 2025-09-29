package main

import (
	"log"
	"website-api/common"
	"website-api/database"
	roleRepo "website-api/repository/role"
	"website-api/service/role"

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
	roleService := role.NewService(roleRepo.NewRepo(db.GormDb))
	if err = roleService.Seed(); err != nil {
		log.Fatal(err)
	}
	log.Print(common.SuccessfullyCreated)
}
