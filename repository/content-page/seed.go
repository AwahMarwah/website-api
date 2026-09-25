package content_page

import (
	"fmt"
	"time"

	content_page "website-api/model/content-page"

	"github.com/google/uuid"
)

func (r *repo) SeedCmsPage() (err error) {
	return r.db.Create([]content_page.CmsPage{
		{
			Slug: "about-us", Title: "About Us", Status: true,
			Content: `<p><strong>Website API</strong> adalah platform e-commerce sederhana yang menyediakan solusi backend yang aman dan cepat.</p>`,
		},
		{
			Slug: "terms-and-conditions", Title: "Terms & Conditions", Status: true,
			Content: `<p>Syarat & ketentuan penggunaan platform ini.</p><ul><li>Data harus benar</li><li>Dilarang menyalahgunakan sistem</li></ul>`,
		},
		{
			Slug: "privacy-policy", Title: "Privacy Policy", Status: true,
			Content: `<p>Kami menghargai privasi pengguna dan berkomitmen melindungi data pribadi Anda.</p>`,
		},
	}).Error
}

func (r *repo) SeedCmsFaq() error {
	faqs := []struct {
		Question string
		Answer   string
		OrderNo  int
	}{
		{"Bagaimana cara memesan produk?", "Pilih produk yang diinginkan, tambahkan ke keranjang, lalu lakukan checkout. Pilih alamat pengiriman dan metode pembayaran, kemudian bayar.", 1},
		{"Metode pembayaran apa yang didukung?", "Kami mendukung pembayaran melalui Midtrans (transfer bank, e-wallet, kartu kredit/debit).", 2},
		{"Bagaimana cara melacak pesanan?", "Masuk ke Dashboard → Pesanan Saya. Status diperbarui dari PENDING menjadi PAID, PROCESSING, SHIPPED, hingga COMPLETED.", 3},
		{"Bagaimana cara membatalkan pesanan?", "Pembatalan hanya dapat dilakukan pada status PENDING sebelum pembayaran.", 4},
		{"Bagaimana cara mengatur alamat pengiriman?", "Masuk ke Dashboard → Alamat Saya → Tambah Alamat. Isi provinsi, kota, kecamatan, kelurahan, alamat lengkap, dan kode pos.", 5},
		{"Apakah data pribadi saya aman?", "Kami menggunakan token autentikasi JWT dan tidak membagikan data pribadi kepada pihak ketiga.", 6},
	}

	for _, f := range faqs {
		id := uuid.NewString()
		faq := content_page.CmsFaq{
			Id:        id,
			Question:  f.Question,
			Answer:    f.Answer,
			OrderNo:   f.OrderNo,
			IsActive:  true,
			CreatedAt: time.Now(),
		}
		if err := r.db.Where("question = ?", f.Question).FirstOrCreate(&faq).Error; err != nil {
			return fmt.Errorf("seed FAQ failed: %w", err)
		}
	}
	return nil
}