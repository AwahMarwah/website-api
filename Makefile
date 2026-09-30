.PHONY: run run_worker build vet swag \
	run_db_migrate_up run_db_migrate_down \
	run_db_seed_role run_db_seed_content_page run_db_seed_menu \
	run_db_seed_merchant run_db_seed_catalog run_db_seed

run:
	go run main.go

# Worker background (Asynq): auto-cancel order kedaluwarsa + rekonsiliasi pembayaran.
# Harus dijalankan terpisah dari `run`, bukan dipanggil olehnya.
run_worker:
	go run cmd/worker/main.go

build:
	go build -o bin/api main.go
	go build -o bin/worker cmd/worker/main.go

vet:
	go vet ./...

# docs/ adalah hasil generate. Jalankan setelah mengubah anotasi swagger.
swag:
	swag init -g main.go -o docs

run_db_migrate_up:
	go run database/migrate/up/up.go

run_db_migrate_down:
	go run database/migrate/down/down.go

run_db_seed_role:
	go run database/migrate/seeding/role/role.go

run_db_seed_content_page:
	go run database/migrate/seeding/content-page/content_page.go

run_db_seed_menu:
	go run database/migrate/seeding/menu/menu.go

run_db_seed_merchant:
	go run database/migrate/seeding/merchant/merchant.go

run_db_seed_catalog:
	go run database/migrate/seeding/catalog/catalog.go

run_db_seed: run_db_seed_role run_db_seed_menu run_db_seed_merchant run_db_seed_catalog run_db_seed_content_page
