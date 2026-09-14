package minio

type Provider interface {
	// UploadFile mengunggah file ke bucket dan mengembalikan URL publik permanen.
	UploadFile(objectKey string, body []byte, contentType string) (url string, err error)
	// PublicURL membuat URL publik dari object key (tanpa memeriksa eksistensi).
	PublicURL(objectKey string) string
}