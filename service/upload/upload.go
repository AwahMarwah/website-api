package upload

import (
	"fmt"
	"net/http"
	productModel "website-api/model/product"
	uploadModel "website-api/model/upload"
	"website-api/repository/product"
	userRepo "website-api/repository/user"
	"website-api/third-party/provider/minio"

	"github.com/google/uuid"
)

const (
	ScopeAvatar  = "avatar"
	ScopeProduct = "product"

	MaxFileSize = 2 * 1024 * 1024 // 2MB
)

type (
	IService interface {
		UploadImage(req *UploadImageRequest) (uploadModel.UploadImageResponse, int, error)
		UploadProductImage(req *UploadProductImageRequest) (uploadModel.UploadImageResponse, int, error)
	}

	service struct {
		minioProvider minio.Provider
		userRepo      userRepo.IRepo
		productRepo   product.IRepo
	}
)

type UploadImageRequest struct {
	Scope       string
	ResourceID  string
	ContentType string
	Size        int64
	Ext         string
	Body        []byte
}

type UploadProductImageRequest struct {
	ProductID   string
	ContentType string
	Size        int64
	Ext         string
	Body        []byte
}

func NewService(minioProvider minio.Provider, userRepo userRepo.IRepo, productRepo product.IRepo) IService {
	return &service{minioProvider: minioProvider, userRepo: userRepo, productRepo: productRepo}
}

// UploadImage digunakan untuk upload avatar & gambar produk (generic scope).
func (s *service) UploadImage(req *UploadImageRequest) (uploadModel.UploadImageResponse, int, error) {
	if req.Size > MaxFileSize {
		return uploadModel.UploadImageResponse{}, http.StatusBadRequest, fmt.Errorf("file terlalu besar (maks 2MB)")
	}
	if req.ContentType != "image/jpeg" && req.ContentType != "image/png" {
		return uploadModel.UploadImageResponse{}, http.StatusBadRequest, fmt.Errorf("tipe file harus jpg/png")
	}
	if len(req.Body) == 0 {
		return uploadModel.UploadImageResponse{}, http.StatusBadRequest, fmt.Errorf("file kosong")
	}
	if req.ResourceID == "" {
		return uploadModel.UploadImageResponse{}, http.StatusBadRequest, fmt.Errorf("resource_id wajib diisi")
	}

	var objectKey, scope string
	switch req.Scope {
	case ScopeAvatar:
		scope = "users"
		objectKey = fmt.Sprintf("users/%s/%s.%s", req.ResourceID, uuid.NewString(), req.Ext)
	case ScopeProduct:
		scope = "products"
		objectKey = fmt.Sprintf("products/%s/%s.%s", req.ResourceID, uuid.NewString(), req.Ext)
		if _, err := s.productRepo.FindByID(req.ResourceID); err != nil {
			return uploadModel.UploadImageResponse{}, http.StatusBadRequest, fmt.Errorf("produk tidak ditemukan")
		}
	default:
		return uploadModel.UploadImageResponse{}, http.StatusBadRequest, fmt.Errorf("scope tidak valid")
	}

	url, err := s.minioProvider.UploadFile(objectKey, req.Body, req.ContentType)
	if err != nil {
		return uploadModel.UploadImageResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal upload ke storage: %w", err)
	}

	switch scope {
	case "users":
		if err := s.userRepo.Update(&req.ResourceID, &map[string]any{"picture": url}); err != nil {
			return uploadModel.UploadImageResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal simpan avatar: %w", err)
		}
	case "products":
		count, _ := s.productRepo.ProductImageCount(req.ResourceID)
		img := &productModel.ProductImage{
			ProductID:  req.ResourceID,
			ImageURL:   url,
			IsPrimary:  count == 0,
			SortOrder:  int(count),
		}
		if err := s.productRepo.CreateProductImage(img); err != nil {
			return uploadModel.UploadImageResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal simpan gambar produk: %w", err)
		}
	}

	return uploadModel.UploadImageResponse{
		URL:       url,
		ObjectKey: objectKey,
		Scope:     req.Scope,
		Resource:  req.ResourceID,
	}, http.StatusOK, nil
}

// UploadProductImage untuk endpoint khusus produk: selalu secondary (primary ditentukan via endpoint set-primary).
func (s *service) UploadProductImage(req *UploadProductImageRequest) (uploadModel.UploadImageResponse, int, error) {
	if req.Size > MaxFileSize {
		return uploadModel.UploadImageResponse{}, http.StatusBadRequest, fmt.Errorf("file terlalu besar (maks 2MB)")
	}
	if req.ContentType != "image/jpeg" && req.ContentType != "image/png" {
		return uploadModel.UploadImageResponse{}, http.StatusBadRequest, fmt.Errorf("tipe file harus jpg/png")
	}
	if len(req.Body) == 0 {
		return uploadModel.UploadImageResponse{}, http.StatusBadRequest, fmt.Errorf("file kosong")
	}

	if _, err := s.productRepo.FindByID(req.ProductID); err != nil {
		return uploadModel.UploadImageResponse{}, http.StatusBadRequest, fmt.Errorf("produk tidak ditemukan")
	}

	objectKey := fmt.Sprintf("products/%s/%s.%s", req.ProductID, uuid.NewString(), req.Ext)

	url, err := s.minioProvider.UploadFile(objectKey, req.Body, req.ContentType)
	if err != nil {
		return uploadModel.UploadImageResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal upload ke storage: %w", err)
	}

	maxSort, err := s.productRepo.MaxSortOrder(req.ProductID)
	if err != nil {
		return uploadModel.UploadImageResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal menghitung sort_order: %w", err)
	}

	// upload pertama otomatis primary
	count, _ := s.productRepo.ProductImageCount(req.ProductID)

	img := &productModel.ProductImage{
		ProductID:  req.ProductID,
		ImageURL:   url,
		IsPrimary:  count == 0,
		SortOrder:  maxSort + 1,
	}
	if err := s.productRepo.CreateProductImage(img); err != nil {
		return uploadModel.UploadImageResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal simpan gambar produk: %w", err)
	}

	return uploadModel.UploadImageResponse{
		URL:       url,
		ObjectKey: objectKey,
		Scope:     "product",
		Resource:  req.ProductID,
	}, http.StatusOK, nil
}
