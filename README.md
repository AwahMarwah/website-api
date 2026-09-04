# Website API — Backend Ecommerce

Backend RESTful API untuk platform e-commerce sederhana, dibangun dengan clean architecture.

## Tech Stack

| Komponen | Teknologi |
|---|---|
| Bahasa | Go 1.24.3 |
| HTTP Framework | Gin 1.10.1 |
| ORM | GORM + PostgreSQL 1.6.0 |
| Cache | Redis (go-redis v9) |
| Background Job | Asynq (Redis) |
| Autentikasi | JWT (HS256) + bcrypt |
| Pembayaran | Midtrans (Snap + Payment Link + Webhook) |
| Pengiriman | RajaOngkir / Komerce API |
| API Docs | Swaggo (Swagger) |

## Struktur Proyek

```
website-api/
├── main.go                          # Entry point
├── Makefile                         # Perintah umum
├── .env                             # Konfigurasi lingkungan
│
├── cache/                           # Redis cache layer
├── common/                          # Konstanta global
├── controller/                      # HTTP handlers (Gin)
│   ├── auth/                        # Login, register, forgot password
│   ├── brand/                       # Brand
│   ├── cart/                        # Keranjang belanja
│   ├── category/                    # Kategori produk
│   ├── content-page/                # CMS pages & FAQ
│   ├── health-check/                # Health check
│   ├── master/                      # Data master (provinsi, kota, kecamatan, kelurahan)
│   ├── menu/                        # Menu RBAC
│   ├── order/                       # Order, checkout, payment
│   ├── permission/                  # Permission list
│   ├── product/                     # Produk + detail
│   ├── role/                        # Manajemen role
│   ├── user/                        # User management
│   └── user_address/                # Alamat pengguna
│
├── database/
│   ├── database.go                  # Koneksi PostgreSQL
│   ├── redis.go                     # Koneksi Redis
│   ├── encrypt/                     # JWT & hashing
│   ├── transaction/                 # DB transaction wrapper
│   └── migrate/                     # SQL migration & seeding
│       ├── sql/                     # File migrasi SQL
│       ├── up/up.go                 # Jalankan migrasi
│       ├── down/down.go             # Rollback migrasi
│       └── seeding/                 # Seed data
│
├── library/
│   ├── cache/                       # Cache helper
│   ├── helper/email/                # Kirim email
│   ├── helper/filter/               # GORM scope filters
│   ├── pagination/                  # Helper pagination
│   └── response/                    # Standardized JSON response
│
├── middleware/                       # Auth, RBAC, CORS, Ngrok
├── model/                           # Entity, request, response structs
├── repository/                      # Database access layer (GORM)
├── router/                          # Route registration
├── service/                         # Business logic
├── task/                            # Task type & payload (Asynq)
├── queue/                           # Queue client
├── worker/                          # Worker handler + scheduler
├── third-party/provider/            # Integrasi eksternal
│   ├── midtrans/                    # Midtrans payment
│   └── rajaongkir/                  # RajaOngkir shipping
├── utils/                           # Utility functions
├── templates/email/                 # Template HTML email
└── docs/                            # Swagger (auto-generated)
```

## Prasyarat

- **Go** 1.24+
- **PostgreSQL** (tested: v17)
- **Redis** (untuk cache & background jobs)
- **SMTP Server** (opsional, untuk email verifikasi & reset password — di dev bisa pakai MailHog: `localhost:1025`)

## Memulai (Getting Started)

### 1. Clone & Install

```bash
git clone https://github.com/username/website-api.git
cd website-api
go mod tidy
```

### 2. Setup Environment

```bash
cp .env.example .env
```

Edit `.env` sesuai konfigurasi lokal Anda.

### 3. Database Migration

```bash
make run_db_migrate_up
```

### 4. Seed Data

```bash
# Seed role (super_admin, admin, staff, finance, merchant, consumer, guest)
make run_db_seed_role

# Seed content page (CMS)
make run_db_seed_content_page

# Seed menu & permission untuk RBAC
make run_db_seed_menu
```

### 5. Jalankan Server

```bash
make run
# atau
go run main.go
```

Server berjalan di `http://localhost:8085` (atau sesuai `PORT` di `.env`).

Swagger UI: `http://localhost:8085/swagger/index.html`

## Database

- **25 tabel** PostgreSQL dengan foreign keys
- **10 migrasi** SQL (termasuk menu/RBAC terbaru)
- Struktur relasi:

```
users ──────┐
             ├── orders ──── order_items
             ├── cart_items
             └── user_addresses ── city_shipping_mappings
                     │
                     ├── provinces → cities → districts → subdistricts

products ──┐
            ├── product_variants
            ├── product_images
            ├── product_categories → categories
            └── reviews

roles ──────┐
            ├── role_menus ── menus
            └── permissions (menu_id)

brands ──── products
```

## API Endpoints

### Public (Tanpa Autentikasi)

| Method | Endpoint | Deskripsi |
|---|---|---|
| GET | `/health` | Health check |
| POST | `/user/sign-up` | Daftar akun baru |
| POST | `/user/sign-in` | Masuk (login) |
| GET | `/user/verify-email?token=` | Verifikasi email via link |
| POST | `/user/verify-email` | Verifikasi email via JSON |
| POST | `/auth/forgot-password` | Lupa password |
| POST | `/auth/reset-password` | Reset password |
| POST | `/auth/resend-verification` | Ulangi verifikasi email |
| GET | `/product` | Daftar produk (filter, search, pagination) |
| GET | `/product/:id` | Detail produk + variants |
| GET | `/brand` | Daftar brand |
| GET | `/brand/:slug` | Detail brand |
| GET | `/category` | Daftar kategori |
| GET | `/master/provincies` | Daftar provinsi |
| GET | `/master/cities` | Daftar kota (filter provinsi via `search`) |
| GET | `/master/district` | Daftar kecamatan (filter kota via `search`) |
| GET | `/master/subdistricts` | Daftar kelurahan/desa (filter kecamatan via `search`) |

### Private (Autentikasi Required)

| Method | Endpoint | Deskripsi |
|---|---|---|
| DELETE | `/user/sign-out` | Keluar (logout) |
| GET | `/user` | Daftar pengguna (admin) |
| GET | `/user/:id` | Detail pengguna |
| PUT | `/user/:id` | Update profil |
| GET | `/user-address` | Daftar alamat |
| POST | `/user-address` | Tambah alamat |
| GET | `/user-address/:id` | Detail alamat |
| PUT | `/user-address/:id` | Update alamat |
| DELETE | `/user-address/:id` | Hapus alamat |
| PATCH | `/user-address/:id/primary` | Ubah alamat utama |
| GET | `/cart` | Keranjang belanja |
| POST | `/cart` | Tambah item ke keranjang |
| DELETE | `/cart` | Hapus item dari keranjang |
| POST | `/order` | Checkout (buat order + Midtrans Snap) |
| GET | `/order` | Riwayat order saya |
| GET | `/order/:id` | Detail order |
| POST | `/order/:id/payment-link` | Buat payment link Midtrans |
| GET | `/menu/my` | Menu navigasi untuk role saya |

### Admin Only (super_admin/admin)

| Method | Endpoint | Deskripsi |
|---|---|---|
| POST | `/role` | Buat role |
| GET | `/role` | Daftar role |
| GET | `/role/:id` | Detail role |
| PUT | `/role/:id` | Update role |
| DELETE | `/role/:id` | Hapus role |
| GET | `/admin/orders` | Semua order (admin) |
| GET | `/menu` | Daftar menu |
| GET | `/menu/tree` | Tree menu |
| POST | `/menu` | Buat menu |
| PUT | `/menu/:id` | Update menu |
| DELETE | `/menu/:id` | Hapus menu |
| GET | `/role/:id/menus` | Menu yang ditetapkan ke role |
| PUT | `/role/:id/menus` | Tetapkan menu ke role |
| GET | `/permission` | Daftar permission |

### Webhook

| Method | Endpoint | Deskripsi |
|---|---|---|
| POST | `/order/notification` | Webhook notifikasi Midtrans (tanpa JWT) |

## Arsitektur

### Layering (Clean Architecture)

```
controller/ → service/ → repository/ → model/
     │            │            │
     │            │            └── GORM / PostgreSQL
     │            ├── database/transaction (opsional)
     │            ├── cache/ (opsional)
     │            └── third-party/ (opsional)
     ├── library/response (standardized JSON response)
     └── middleware/ (Auth, RBAC)
```

### Middleware

| Middleware | Fungsi |
|---|---|
| `AuthMiddleware(db)` | Verifikasi JWT, set context `user_id`, `role_name` |
| `RequireRole(roles...)` | Batasi akses berdasarkan role |
| `Permission(db, name)` | Batasi akses berdasarkan permission (RBAC) |

### Background Jobs (Asynq)

| Task | Fungsi | Jadwal |
|---|---|---|
| `email:reset_password` | Kirim email reset password | On-demand |
| `email:payment_success` | Kirim email invoice pembayaran | On-demand |
| `order:cancel_expired` | Auto-cancel order PENDING yang expired | Setiap 5 menit |

### Payment (Midtrans)

- **Snap**: Checkout → Snap token → Redirect ke halaman pembayaran Midtrans
- **Payment Link**: Buat link pembayaran yang bisa dibagikan
- **Webhook**: `POST /order/notification` — verifikasi signature SHA512 + update status order

### RBAC (Role-Based Access Control)

```
roles → role_menus → menus
                    ↓
              permissions (menu_id)
```

- Menu ditetapkan ke role melalui tabel `role_menus`
- Setiap menu memiliki permission yang terhubung
- `super_admin` dan `admin` bypass semua permission check
- Menu bisa hierarkis (parent-child)

### Seed Data

| Seed | Isi |
|---|---|
| Role | super_admin, admin, staff, finance, merchant, consumer, guest |
| Menu | Dashboard, Produk, Pesanan, Pengguna, Role, Menu |
| Permission | product.view, order.view, user.view, user.update, role.manage, menu.manage |
| Role-Menu | Semua menu ditetapkan ke super_admin & admin |

## Perintah Makefile

| Perintah | Deskripsi |
|---|---|
| `make run` | Jalankan server |
| `make run_db_migrate_up` | Jalankan migrasi ke atas |
| `make run_db_migrate_down` | Rollback migrasi |
| `make run_db_seed_role` | Seed role |
| `make run_db_seed_content_page` | Seed content page |
| `make run_db_seed_menu` | Seed menu & permission |

## Known Issues

- JWT secret masih development (`simple_ecommerce_change_me_in_production`) — ganti untuk production
- RajaOngkir `CalculateCost` belum terimplementasi (stub)
- Tidak ada user profile update endpoint (hanya admin bisa update phone/role)
- Brand slug route belum dipakai oleh frontend
- CORS `AllowAllOrigins: true` — perlu diketatkan untuk production
