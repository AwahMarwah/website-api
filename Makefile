run :
	go run main.go

run_db_migrate_up:
	go run database/migrate/up/up.go

run_db_migrate_down:
	go run database/migrate/down/down.go

run_db_seed_role:
	go run database/migrate/seeding/role/role.go

run_db_seed_content_page:
	go run database/migrate/seeding/content-page/content_page.go


# Jalankan test dengan environment yang sudah disetup
go test -v ./service/user/

# Jalankan dengan timeout yang lebih pendek
go test -v -timeout 30s ./service/user/

# Jalankan dengan parallel testing
go test -v -parallel 4 ./service/user/