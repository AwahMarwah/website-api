package main

import (
	"fmt"
	"log"
	"time"
	"website-api/database"

	_ "github.com/joho/godotenv/autoload"
)

type brandSeed struct {
	id   string
	name string
	slug string
}

type catSeed struct {
	id       string
	name     string
	slug     string
	parentID string
}

type prodSeed struct {
	id          string
	name        string
	slug        string
	description string
	basePrice   float64
	sku         string
	brandID     string
	merchantIdx int // 0..2 -> toko dummy
}

type variantSeed struct {
	prodID string
	sku    string
	name   string
	price  float64
	stock  int
	weightKg float64
}

var brands = []brandSeed{
	{id: "brand-001", name: "Nike", slug: "nike"},
	{id: "brand-002", name: "Adidas", slug: "adidas"},
	{id: "brand-003", name: "Uniqlo", slug: "uniqlo"},
	{id: "brand-004", name: "H&M", slug: "hm"},
	{id: "brand-005", name: "Zara", slug: "zara"},
	{id: "brand-006", name: "Levi's", slug: "levis"},
	{id: "brand-007", name: "Converse", slug: "converse"},
	{id: "brand-008", name: "Puma", slug: "puma"},
	{id: "brand-009", name: "New Balance", slug: "new-balance"},
	{id: "brand-010", name: "Vans", slug: "vans"},
	{id: "brand-011", name: "The North Face", slug: "the-north-face"},
}

var categories = []catSeed{
	{id: "cat-006", name: "Electronics", slug: "electronics"},
	{id: "cat-007", name: "Accessories", slug: "accessories"},
	{id: "cat-008", name: "Outerwear", slug: "outerwear"},
	{id: "cat-009", name: "Pants", slug: "pants"},
	{id: "cat-010", name: "Bags", slug: "bags", parentID: "cat-007"},
	{id: "cat-011", name: "Smartphone", slug: "smartphone", parentID: "cat-006"},
	{id: "cat-012", name: "Headphone", slug: "headphone", parentID: "cat-006"},
	{id: "cat-013", name: "Jacket", slug: "jacket", parentID: "cat-008"},
	{id: "cat-014", name: "Jeans", slug: "jeans", parentID: "cat-009"},
	{id: "cat-015", name: "Backpack", slug: "backpack", parentID: "cat-010"},
}

var products = []prodSeed{
	// Men Fashion
	{id: "prod-004", name: "H&M Slim Fit Chino", slug: "hm-slim-fit-chino", sku: "HM-CHINO-001", brandID: "brand-004", merchantIdx: 1, basePrice: 349000,
		description: "Celana chino slim fit dari bahan stretch cotton premium. Nyaman untuk aktivitas harian maupun semi-formal, tersedia beberapa warna."},
	{id: "prod-005", name: "Zara Linen Shirt", slug: "zara-linen-shirt", sku: "ZR-LS-001", brandID: "brand-005", merchantIdx: 2, basePrice: 429000,
		description: "Kemeja linen ringan dengan tekstur natural. Cocok untuk iklim tropis, model relaxed fit dengan kerah mandarin."},
	{id: "prod-006", name: "Levi's 501 Original Jeans", slug: "levis-501-original", sku: "LV-501-001", brandID: "brand-006", merchantIdx: 0, basePrice: 899000,
		description: "Celana jeans legendaris Levi's 501 dengan straight fit dan denim premium 100% katun."},
	{id: "prod-007", name: "The North Face Apex Jacket", slug: "tnf-apex-jacket", sku: "TNF-APX-001", brandID: "brand-011", merchantIdx: 1, basePrice: 1250000,
		description: "Jaket windproof dengan teknologi gore-tex. Tahan angin & air ringan, cocok untuk outdoor."},

	// Women Fashion
	{id: "prod-008", name: "H&M Floral Dress", slug: "hm-floral-dress", sku: "HM-DRS-001", brandID: "brand-004", merchantIdx: 2, basePrice: 299000,
		description: "Gaun floral dengan bahan ringan dan potongan A-line. Elegan untuk acara santai maupun formal."},
	{id: "prod-009", name: "Zara Wool Coat", slug: "zara-wool-coat", sku: "ZR-COAT-001", brandID: "brand-005", merchantIdx: 0, basePrice: 1500000,
		description: "Mantel wol premium dengan double-breasted design. Hangat dan berkelas untuk musim dingin."},
	{id: "prod-010", name: "Levi's Women Denim Jacket", slug: "levis-denim-jacket-women", sku: "LV-DJ-W-001", brandID: "brand-006", merchantIdx: 1, basePrice: 749000,
		description: "Jaket denim wanita dengan classic button closure, bahan stretch denim yang nyaman."},

	// Shoes
	{id: "prod-011", name: "Converse Chuck Taylor All Star", slug: "converse-chuck-taylor", sku: "CV-70-001", brandID: "brand-007", merchantIdx: 2, basePrice: 899000,
		description: "Sepatu kanvas ikonik Chuck Taylor All Star. Vegan leather, tahan lama, cocok untuk gaya kasual."},
	{id: "prod-012", name: "Puma Suede Classic", slug: "puma-suede-classic", sku: "PU-SD-001", brandID: "brand-008", merchantIdx: 0, basePrice: 799000,
		description: "Sepatu suede klasik Puma dengan midsole cushioning. Gaya retro yang timeless."},
	{id: "prod-013", name: "New Balance 574", slug: "new-balance-574", sku: "NB-574-001", brandID: "brand-009", merchantIdx: 1, basePrice: 1250000,
		description: "Sepatu lifestyle NB 574 dengan ENCAP midsole technology. Nyaman untuk sehari-hari."},
	{id: "prod-014", name: "Vans Old Skool", slug: "vans-old-skool", sku: "VN-OS-001", brandID: "brand-010", merchantIdx: 2, basePrice: 949000,
		description: "Sepatu skate ikonik Vans Old Skool dengan suede dan canvas, fitur padded collar."},

	// Electronics
	{id: "prod-015", name: "New Balance FuelCell Sneakers", slug: "nb-fuelcell", sku: "NB-FC-001", brandID: "brand-009", merchantIdx: 0, basePrice: 1750000,
		description: "Sepatu lari FuelCell dengan energi return tinggi, ringan dan responsif."},
	{id: "prod-016", name: "The North Face Borealis Backpack", slug: "tnf-borealis", sku: "TNF-BR-001", brandID: "brand-011", merchantIdx: 1, basePrice: 850000,
		description: "Ransel Borealis dengan kompartemen laptop 15\", tahan air, ergonomis."},
	{id: "prod-017", name: "Converse Run Star Hike", slug: "converse-run-star-hike", sku: "CV-RSH-001", brandID: "brand-007", merchantIdx: 2, basePrice: 1100000,
		description: "Siluet futuristik Converse dengan platform outsole bergerigi, perpaduan sneaker & boot."},
}

var variants = []variantSeed{
	// prod-004 Chino
	{prodID: "prod-004", sku: "HM-CHINO-001-BLK-30", name: "Black - 30", price: 349000, stock: 25, weightKg: 0.4},
	{prodID: "prod-004", sku: "HM-CHINO-001-NAVY-32", name: "Navy - 32", price: 349000, stock: 20, weightKg: 0.4},
	{prodID: "prod-004", sku: "HM-CHINO-001-KHAKI-32", name: "Khaki - 32", price: 349000, stock: 15, weightKg: 0.4},
	// prod-005 Zara Shirt
	{prodID: "prod-005", sku: "ZR-LS-001-M", name: "White - M", price: 429000, stock: 18, weightKg: 0.2},
	{prodID: "prod-005", sku: "ZR-LS-001-L", name: "Beige - L", price: 429000, stock: 12, weightKg: 0.2},
	// prod-006 Levi's 501
	{prodID: "prod-006", sku: "LV-501-001-30", name: "Wash Blue - 30", price: 899000, stock: 20, weightKg: 0.8},
	{prodID: "prod-006", sku: "LV-501-001-32", name: "Wash Blue - 32", price: 899000, stock: 30, weightKg: 0.8},
	{prodID: "prod-006", sku: "LV-501-001-34", name: "Dark Blue - 34", price: 899000, stock: 16, weightKg: 0.8},
	// prod-007 TNF Jacket
	{prodID: "prod-007", sku: "TNF-APX-001-M", name: "Black - M", price: 1250000, stock: 10, weightKg: 0.9},
	{prodID: "prod-007", sku: "TNF-APX-001-L", name: "Navy - L", price: 1250000, stock: 8, weightKg: 0.9},
	// prod-008 H&M Dress
	{prodID: "prod-008", sku: "HM-DRS-001-M", name: "Rose - M", price: 299000, stock: 22, weightKg: 0.3},
	{prodID: "prod-008", sku: "HM-DRS-001-S", name: "Coral - S", price: 299000, stock: 14, weightKg: 0.3},
	// prod-009 Zara Coat
	{prodID: "prod-009", sku: "ZR-COAT-001-S", name: "Camel - S", price: 1500000, stock: 6, weightKg: 1.1},
	{prodID: "prod-009", sku: "ZR-COAT-001-M", name: "Camel - M", price: 1500000, stock: 5, weightKg: 1.1},
	// prod-010 Levi's Denim Jacket
	{prodID: "prod-010", sku: "LV-DJ-W-001-M", name: "Wash Denim - M", price: 749000, stock: 12, weightKg: 0.6},
	{prodID: "prod-010", sku: "LV-DJ-W-001-L", name: "Wash Denim - L", price: 749000, stock: 9, weightKg: 0.6},
	// prod-011 Converse Chuck Taylor
	{prodID: "prod-011", sku: "CV-70-001-40", name: "Black - 40", price: 899000, stock: 30, weightKg: 0.5},
	{prodID: "prod-011", sku: "CV-70-001-41", name: "White - 41", price: 899000, stock: 26, weightKg: 0.5},
	{prodID: "prod-011", sku: "CV-70-001-42", name: "Black - 42", price: 899000, stock: 18, weightKg: 0.5},
	// prod-012 Puma Suede
	{prodID: "prod-012", sku: "PU-SD-001-40", name: "Classic Red - 40", price: 799000, stock: 22, weightKg: 0.5},
	{prodID: "prod-012", sku: "PU-SD-001-41", name: "Black - 41", price: 799000, stock: 20, weightKg: 0.5},
	// prod-013 NB 574
	{prodID: "prod-013", sku: "NB-574-001-41", name: "Grey - 41", price: 1250000, stock: 15, weightKg: 0.6},
	{prodID: "prod-013", sku: "NB-574-001-42", name: "Grey - 42", price: 1250000, stock: 12, weightKg: 0.6},
	{prodID: "prod-013", sku: "NB-574-001-43", name: "Navy - 43", price: 1250000, stock: 10, weightKg: 0.6},
	// prod-014 Vans Old Skool
	{prodID: "prod-014", sku: "VN-OS-001-40", name: "Black/White - 40", price: 949000, stock: 24, weightKg: 0.5},
	{prodID: "prod-014", sku: "VN-OS-001-42", name: "Black/White - 42", price: 949000, stock: 20, weightKg: 0.5},
	// prod-015 NB FuelCell
	{prodID: "prod-015", sku: "NB-FC-001-41", name: "Electric - 41", price: 1750000, stock: 8, weightKg: 0.3},
	{prodID: "prod-015", sku: "NB-FC-001-42", name: "White - 42", price: 1750000, stock: 7, weightKg: 0.3},
	// prod-016 TNF Borealis
	{prodID: "prod-016", sku: "TNF-BR-001", name: "Black", price: 850000, stock: 20, weightKg: 0.9},
	{prodID: "prod-016", sku: "TNF-BR-001-GRY", name: "Grey", price: 850000, stock: 12, weightKg: 0.9},
	// prod-017 Converse Run Star
	{prodID: "prod-017", sku: "CV-RSH-001-40", name: "Black - 40", price: 1100000, stock: 14, weightKg: 0.6},
	{prodID: "prod-017", sku: "CV-RSH-001-41", name: "White - 41", price: 1100000, stock: 11, weightKg: 0.6},
}

// merchantByIndex memetakan index seed merchant
func merchantByIndex() []string {
	return []string{
		"fe4c58b4-d58f-4ac4-a00f-a73e05a233d0", // Toko Nike Sporting
		"46e9ede4-b1ff-4664-b03f-60303462f670", // Adidas Store Bandung
		"ceeacc20-a635-4efd-9971-eef331f91eab", // Uniqlo Central Java
	}
}

// productCategoryMap: slug produk -> kategori id
var productCategoryMap = map[string][]string{
	"prod-004": {"cat-001", "cat-009"},
	"prod-005": {"cat-001", "cat-004"},
	"prod-006": {"cat-001", "cat-014"},
	"prod-007": {"cat-001", "cat-013"},
	"prod-008": {"cat-002"},
	"prod-009": {"cat-002", "cat-013"},
	"prod-010": {"cat-002", "cat-008"},
	"prod-011": {"cat-003", "cat-005"},
	"prod-012": {"cat-003", "cat-005"},
	"prod-013": {"cat-003", "cat-005"},
	"prod-014": {"cat-003", "cat-005"},
	"prod-015": {"cat-003", "cat-005"},
	"prod-016": {"cat-015", "cat-007"},
	"prod-017": {"cat-003"},
}

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

	now := time.Now()
	merchants := merchantByIndex()

	// 1. Brands (upsert by slug; set logo_url)
	for _, b := range brands {
		logoURL := fmt.Sprintf("https://picsum.photos/seed/brand-%s/300/300", b.slug)
		if err := db.GormDb.Exec(
			"INSERT INTO brands (id, name, slug, logo_url, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name, logo_url = EXCLUDED.logo_url, updated_at = EXCLUDED.updated_at",
			b.id, b.name, b.slug, logoURL, now, now,
		).Error; err != nil {
			log.Printf("brand %s: %v", b.name, err)
		}
	}

	// 2. Categories (insert if not exists)
	for _, c := range categories {
		parent := c.parentID
		if parent == "" {
			parent = "NULL"
		}
		if err := db.GormDb.Exec(
			"INSERT INTO categories (id, name, slug, parent_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT (slug) DO NOTHING",
			c.id, c.name, c.slug, nullable(c.parentID), now, now,
		).Error; err != nil {
			log.Printf("category %s: %v", c.name, err)
		}
	}

	// 3. Products
	for _, p := range products {
		if err := db.GormDb.Exec(
			"INSERT INTO products (id, brand_id, merchant_id, sku, name, slug, description, base_price, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'active', ?, ?) ON CONFLICT (slug) DO NOTHING",
			p.id, p.brandID, merchants[p.merchantIdx], p.sku, p.name, p.slug, p.description, p.basePrice, now, now,
		).Error; err != nil {
			log.Printf("product %s: %v", p.name, err)
		}
	}

	// 4. Variants
	for _, v := range variants {
		if err := db.GormDb.Exec(
			"INSERT INTO product_variants (id, product_id, sku, variant_name, price, stock, weight, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, TRUE, ?, ?) ON CONFLICT (sku) DO NOTHING",
			fmt.Sprintf("var-%s", v.sku), v.prodID, v.sku, v.name, v.price, v.stock, v.weightKg, now, now,
		).Error; err != nil {
			log.Printf("variant %s: %v", v.sku, err)
		}
	}

	// 5. Product images (picsum) — reset per product agar bersih & realistis
	for _, p := range products {
		// hapus gambar lama
		db.GormDb.Exec("DELETE FROM product_images WHERE product_id = ?", p.id)

		imageCount := 3
		for i := 1; i <= imageCount; i++ {
			imgURL := fmt.Sprintf("https://picsum.photos/seed/%s-%d/600/600", p.slug, i)
			isPrimary := i == 1
			if err := db.GormDb.Exec(
				"INSERT INTO product_images (id, product_id, image_url, is_primary, sort_order, created_at) VALUES (?, ?, ?, ?, ?, ?)",
				fmt.Sprintf("img-%s-%d", p.id, i), p.id, imgURL, isPrimary, i, now,
			).Error; err != nil {
				log.Printf("image %s: %v", p.slug, err)
			}
		}
	}

	// 6. Product categories (relink)
	for prodID, catIDs := range productCategoryMap {
		db.GormDb.Exec("DELETE FROM product_categories WHERE product_id = ?", prodID)
		for _, cid := range catIDs {
			if err := db.GormDb.Exec(
				"INSERT INTO product_categories (product_id, category_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				prodID, cid,
			).Error; err != nil {
				log.Printf("pcat %s-%s: %v", prodID, cid, err)
			}
		}
	}

	log.Printf("catalog seeding selesai: %d brand, %d kategori, %d produk, %d varian",
		len(brands), len(categories), len(products), len(variants))
}

func nullable(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}