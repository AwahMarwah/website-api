package main

import (
	"log"
	"time"
	"website-api/database"

	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

type merchantSeed struct {
	id            string
	name          string
	slug          string
	destinationID int64
	cityID        string
	address       string
	productID     string
}

var merchants = []merchantSeed{
	{name: "Toko Nike Sporting", slug: "toko-nike-sporting", destinationID: 17547, cityID: "3174", address: "Jl. Senayan Utara, Jakarta Selatan", productID: "prod-001"},
	{name: "Adidas Store Bandung", slug: "adidas-store-bandung", destinationID: 4916, cityID: "3273", address: "Jl. Cihampelas, Bandung", productID: "prod-002"},
	{name: "Uniqlo Central Java", slug: "uniqlo-central-java", destinationID: 65005, cityID: "3374", address: "Jl. Pandanaran, Semarang", productID: "prod-003"},
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

	merchantIDByName := map[string]string{}
	for _, m := range merchants {
		id := uuid.NewString()
		merchantIDByName[m.name] = id
		if err := db.GormDb.Exec(
			"INSERT INTO merchants (id, name, slug, destination_id, city_id, address, is_active, created_at) VALUES (?, ?, ?, ?, ?, ?, TRUE, ?) ON CONFLICT (slug) DO NOTHING",
			id, m.name, m.slug, m.destinationID, m.cityID, m.address, now,
		).Error; err != nil {
			log.Fatalf("gagal seed merchant %s: %v", m.name, err)
		}
	}

	// Ambil id merchant yang benar-benar ada (untuk idempotent), lalu petakan produk
	for _, m := range merchants {
		var mid string
		if err := db.GormDb.Raw("SELECT id FROM merchants WHERE slug = ? LIMIT 1", m.slug).Scan(&mid).Error; err != nil {
			log.Fatalf("gagal ambil id merchant %s: %v", m.name, err)
		}

		if err := db.GormDb.Exec(
			"UPDATE products SET merchant_id = ? WHERE id = ?",
			mid, m.productID,
		).Error; err != nil {
			log.Printf("gagal link produk %s ke merchant %s: %v", m.productID, m.name, err)
		}
	}

	log.Print("merchant seeding selesai")
}