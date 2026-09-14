package upload

import (
	productRepo "website-api/repository/product"
	userRepo "website-api/repository/user"
	"website-api/service/upload"
	"website-api/third-party/provider/minio"

	"gorm.io/gorm"
)

type controller struct {
	uploadService upload.IService
}

func NewController(db *gorm.DB, minioProvider minio.Provider) *controller {
	return &controller{
		uploadService: upload.NewService(minioProvider, userRepo.NewRepo(db), productRepo.NewRepo(db)),
	}
}