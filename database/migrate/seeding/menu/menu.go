package main

import (
	"log"
	"time"
	"website-api/database"

	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

type menuSeed struct {
	id          string
	name        string
	displayName string
	icon        string
	path        string
	sortOrder   int
}

var menus = []menuSeed{
	{name: "dashboard", displayName: "Dashboard", icon: "dashboard", path: "/admin", sortOrder: 1},
	{name: "produk", displayName: "Produk", icon: "package", path: "/admin/products", sortOrder: 2},
	{name: "pesanan", displayName: "Pesanan", icon: "cart", path: "/admin/orders", sortOrder: 3},
	{name: "pengguna", displayName: "Pengguna", icon: "users", path: "/admin/users", sortOrder: 4},
	{name: "role", displayName: "Role", icon: "shield", path: "/admin/roles", sortOrder: 5},
	{name: "menu", displayName: "Menu", icon: "menu", path: "/admin/menus", sortOrder: 6},
}

type permissionSeed struct {
	name        string
	displayName string
	menu        string
}

var permissions = []permissionSeed{
	{name: "product.view", displayName: "Lihat Produk", menu: "produk"},
	{name: "order.view", displayName: "Lihat Pesanan", menu: "pesanan"},
	{name: "user.view", displayName: "Lihat Pengguna", menu: "pengguna"},
	{name: "user.update", displayName: "Ubah Pengguna", menu: "pengguna"},
	{name: "role.manage", displayName: "Kelola Role", menu: "role"},
	{name: "menu.manage", displayName: "Kelola Menu", menu: "menu"},
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

	// Peta menu id berdasarkan name (gunakan id yang benar-benar ada di DB)
	menuIDByName := map[string]string{}
	for _, m := range menus {
		if err := db.GormDb.Exec(
			"INSERT INTO menus (id, parent_id, name, display_name, icon, path, sort_order, is_active, created_at) VALUES (?, NULL, ?, ?, ?, ?, ?, TRUE, ?) ON CONFLICT (name) DO NOTHING",
			uuid.NewString(), m.name, m.displayName, m.icon, m.path, m.sortOrder, now,
		).Error; err != nil {
			log.Fatalf("gagal seed menu %s: %v", m.name, err)
		}
		var id string
		if err := db.GormDb.Raw("SELECT id FROM menus WHERE name = ? LIMIT 1", m.name).Scan(&id).Error; err != nil {
			log.Fatalf("gagal ambil id menu %s: %v", m.name, err)
		}
		menuIDByName[m.name] = id
	}

	// Seed permissions, link ke menu
	for _, p := range permissions {
		menuID := menuIDByName[p.menu]
		if err := db.GormDb.Exec(
			"INSERT INTO permissions (id, name, display_name, description, menu_id, created_at) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT (name) DO NOTHING",
			uuid.NewString(), p.name, p.displayName, "", menuID, now,
		).Error; err != nil {
			log.Fatalf("gagal seed permission %s: %v", p.name, err)
		}
	}

	// Assign semua menu ke role super_admin & admin
	roleNames := []string{"super_admin", "admin"}
	for _, roleName := range roleNames {
		var roleID string
		if err := db.GormDb.Raw("SELECT id FROM roles WHERE name = ? LIMIT 1", roleName).Scan(&roleID).Error; err != nil {
			log.Fatalf("gagal cari role %s: %v", roleName, err)
		}
		if roleID == "" {
			log.Printf("role %s tidak ditemukan, dilewati", roleName)
			continue
		}
		for _, m := range menus {
		if err := db.GormDb.Exec(
			"INSERT INTO role_menus (id, role_id, menu_id, created_at) VALUES (?, ?, ?, ?) ON CONFLICT (role_id, menu_id) DO NOTHING",
			uuid.NewString(), roleID, menuIDByName[m.name], now,
		).Error; err != nil {
			log.Printf("gagal assign menu %s ke role %s: %v", m.name, roleName, err)
		}
		}
	}

	log.Print("menu & permission seeding selesai")
}
