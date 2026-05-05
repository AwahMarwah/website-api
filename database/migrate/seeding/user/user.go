package user

import (
	"log"
	"website-api/common"
	"website-api/database"
	roleRepo "website-api/repository/role"
	userRepo "website-api/repository/user"
	"website-api/service/user"
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
	userService := user.NewService(userRepo.NewRepo(db.GormDb), roleRepo.NewRepo(db.GormDb))
	if err = userService.Seed(); err != nil {
		log.Fatal(err)
	}
	log.Print(common.SuccessfullyCreated)
}
