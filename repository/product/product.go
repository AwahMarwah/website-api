package product

import (
	productModel "website-api/model/product"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		FindByID(id string) (resData productModel.Product, err error)
		FindDetailByID(id string) (resData productModel.ProductDetailResponse, err error)
		GetProduct(reqQuery *productModel.GetListProductReqQuerry) (resData []productModel.ListProductResponse, count int64, err error)
		Create(product productModel.Product) error
		UpdateProduct(id string, values map[string]any) error
		SetInactive(id string) error
		RebindCategories(productID string, catIDs []string) error
		CreateProductImage(img *productModel.ProductImage) error
		ProductImageCount(productID string) (int64, error)
		FindImagesByProductID(productID string) ([]productModel.ProductImage, error)
		MaxSortOrder(productID string) (int, error)
		SetPrimary(productID, imageID string) error
		DeleteImage(imageID string) error
		SetImagePrimaryFalse(imageID string) error
		UpdateSortOrders(productID string, orders []productModel.ImageSortOrder) error
		Update(product productModel.Product) error
		WithTx(tx *gorm.DB) IRepo
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}
