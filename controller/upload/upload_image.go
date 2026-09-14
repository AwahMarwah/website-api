package upload

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"
	"website-api/service/upload"

	"github.com/gin-gonic/gin"
)

// @Summary Upload Image
// @Description Upload gambar (avatar user atau produk) ke object storage. Tipe jpg/png, maks 2MB.
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File gambar"
// @Param scope formData string true "avatar | product"
// @Param resource_id formData string true "ID user (avatar) atau ID produk (product)"
// @Success 200 {object} upload.UploadImageResponse "Upload berhasil"
// @Router /upload/image [post]
func (c *controller) UploadImage(ctx *gin.Context) {
	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	scope := ctx.PostForm("scope")
	resourceID := ctx.PostForm("resource_id")

	if scope == "avatar" {
		// hanya pemilik akun (atau admin)
		if resourceID != userInfo.UserID && userInfo.Role != "super_admin" && userInfo.Role != "admin" {
			response.Error(ctx, http.StatusForbidden, "forbidden")
			return
		}
	}
	if scope == "product" && userInfo.Role != "super_admin" && userInfo.Role != "admin" {
		response.Error(ctx, http.StatusForbidden, "forbidden")
		return
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "file wajib diisi")
		return
	}

	// validasi ukuran
	if fileHeader.Size > upload.MaxFileSize {
		response.Error(ctx, http.StatusBadRequest, "file terlalu besar (maks 2MB)")
		return
	}

	// validasi tipe via header
	ctype := fileHeader.Header.Get("Content-Type")
	if ctype != "image/jpeg" && ctype != "image/png" {
		response.Error(ctx, http.StatusBadRequest, "tipe file harus jpg/png")
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "gagal membuka file")
		return
	}
	defer src.Close()

	body := make([]byte, fileHeader.Size)
	if _, err := src.Read(body); err != nil {
		response.Error(ctx, http.StatusBadRequest, "gagal membaca file")
		return
	}

	ext := "jpg"
	if ctype == "image/png" {
		ext = "png"
	}

	resData, statusCode, err := c.uploadService.UploadImage(&upload.UploadImageRequest{
		Scope:       scope,
		ResourceID:  resourceID,
		ContentType: ctype,
		Size:        fileHeader.Size,
		Ext:         ext,
		Body:        body,
	})
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}

// @Summary Upload Product Image (Dedicated)
// @Description Upload gambar produk (selalu secondary). Pertama kali otomatis primary. Admin only.
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Product ID"
// @Param file formData file true "File gambar (jpg/png, maks 2MB)"
// @Success 200 {object} upload.UploadImageResponse "Upload berhasil"
// @Router /product/{id}/image [post]
func (c *controller) UploadProductImage(ctx *gin.Context) {
	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	if userInfo.Role != "super_admin" && userInfo.Role != "admin" {
		response.Error(ctx, http.StatusForbidden, "forbidden")
		return
	}

	productID := ctx.Param("id")

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "file wajib diisi")
		return
	}

	if fileHeader.Size > upload.MaxFileSize {
		response.Error(ctx, http.StatusBadRequest, "file terlalu besar (maks 2MB)")
		return
	}

	ctype := fileHeader.Header.Get("Content-Type")
	if ctype != "image/jpeg" && ctype != "image/png" {
		response.Error(ctx, http.StatusBadRequest, "tipe file harus jpg/png")
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "gagal membuka file")
		return
	}
	defer src.Close()

	body := make([]byte, fileHeader.Size)
	if _, err := src.Read(body); err != nil {
		response.Error(ctx, http.StatusBadRequest, "gagal membaca file")
		return
	}

	ext := "jpg"
	if ctype == "image/png" {
		ext = "png"
	}

	resData, statusCode, err := c.uploadService.UploadProductImage(&upload.UploadProductImageRequest{
		ProductID:   productID,
		ContentType: ctype,
		Size:        fileHeader.Size,
		Ext:         ext,
		Body:        body,
	})
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}
