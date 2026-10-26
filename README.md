# DompetKu — Pencatat Pengeluaran Harian

**DompetKu** adalah aplikasi web sederhana dan modern untuk membantu pengguna mencatat dan memantau pengeluaran keuangan pribadi sehari-hari. Dibuat dengan backend Go (Gin + GORM) dan antarmuka web murni (HTML, CSS, JavaScript) tanpa dependensi frontend yang berat.

---

## 🚀 Fitur Utama

1. **Autentikasi Pengguna**:
   - Registrasi akun baru (`username`, `email`, `password`).
   - Login dengan `username` dan `password`.
   - Penyimpanan sesi aman menggunakan **JWT di dalam HTTP-only Cookie** (mencegah risiko XSS).
   - Logout (membersihkan cookie token).
2. **Catatan Pengeluaran (CRUD)**:
   - Tambah pengeluaran: Nama pengeluaran, nominal (Rp), kategori, dan tanggal.
   - Daftar riwayat transaksi terurut dari yang terbaru.
   - Edit pengeluaran yang sudah dicatat.
   - Hapus pengeluaran dengan konfirmasi.
   - **Isolasi Data Ketat**: Setiap pengguna hanya dapat melihat dan mengelola datanya sendiri (`user_id`).
3. **Ringkasan Otomatis**:
   - Total nominal akumulasi pengeluaran.
   - Jumlah transaksi tercatat.
4. **Pencarian & Filter Kategori**:
   - Pencarian berdasarkan nama/keterangan pengeluaran.
   - Filter dropdown kategori: *Makanan, Transportasi, Belanja, Tagihan, Hiburan, Lainnya*.
   - Input kategori dinamis jika memilih opsi *Lainnya*.

---

## 🛠️ Pilihan Teknologi

* **Backend**: Go (`golang`) dengan framework [Gin](https://github.com/gin-gonic/gin)
* **ORM & Database**: [GORM](https://gorm.io/) dengan driver **PostgreSQL**
* **Autentikasi**: JSON Web Token ([golang-jwt/jwt](https://github.com/golang-jwt/jwt)) & [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
* **Frontend**: HTML5, CSS3 murni (Modern & Responsif), Vanilla JavaScript (`fetch` API)
* **Konfigurasi**: [godotenv](https://github.com/joho/godotenv) (`.env`)

---

## 📁 Struktur Direktori

```text
go_vibe/
├── main.go                       # Entry point utama aplikasi
├── internal/
│   ├── config/
│   │   ├── config.go             # Loader konfigurasi dari .env
│   │   └── config_test.go        # Unit test konfigurasi
│   ├── database/
│   │   └── database.go           # Koneksi PostgreSQL & AutoMigrate GORM
│   ├── handlers/
│   │   ├── auth_handler.go       # Handler HTTP autentikasi (Register, Login, Logout, Me)
│   │   └── expense_handler.go    # Handler HTTP CRUD catatan pengeluaran & ringkasan
│   ├── middleware/
│   │   ├── auth_middleware.go    # Middleware JWT HttpOnly Cookie & Bcrypt hashing
│   │   └── auth_middleware_test.go # Unit test JWT & Bcrypt
│   └── models/
│       ├── user.go               # Model User & input binding
│       └── expense.go            # Model Expense & input binding
├── web/                          # Folder antarmuka web (disajikan oleh Gin)
│   ├── css/
│   │   └── style.css             # Styling CSS responsif
│   ├── js/
│   │   ├── auth.js               # Logika form login & register
│   │   └── app.js                # Logika dashboard, modal, CRUD, dan filter
│   ├── index.html                # Halaman utama (Dashboard My Expenses)
│   ├── login.html                # Halaman login
│   └── register.html             # Halaman registrasi
├── docs/
│   └── prompts/
│       └── 00-project-spec.md    # Spesifikasi proyek
├── .env.example                  # Template environment variables
├── .gitignore                    # Mengabaikan file sensitif (.env) dan binary
├── go.mod                        # Modul Go
├── go.sum                        # Checksum dependensi Go
└── README.md                     # Dokumentasi proyek
```

---

## ⚙️ Panduan Menjalankan Aplikasi

### 1. Prasyarat
* Telah menginstal [Go (v1.20+)](https://go.dev/dl/)
* Server database **PostgreSQL** yang sedang berjalan dan database baru yang telah dibuat (misal: nama database `dompetku`).

### 2. Konfigurasi Environment (`.env`)
Salin file `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```
Sesuaikan konfigurasi pada file `.env`:
```dotenv
APP_PORT=8080
DATABASE_URL=postgres://postgres:password_kamu@localhost:5432/dompetku?sslmode=disable
JWT_SECRET=buat-string-acak-yang-panjang-dan-rahasia
```

> **Catatan Keamanan**: File `.env` sudah masuk ke `.gitignore` sehingga tidak akan terunggah ke Git repository.

### 3. Menjalankan Unit Test
Untuk memastikan logika hashing, JWT, dan konfigurasi berjalan baik:
```bash
go test ./...
```

### 4. Menjalankan Server
Jalankan aplikasi dengan perintah:
```bash
go run main.go
```
Jika berhasil, buka browser di:
```text
http://localhost:8080
```

---

## 📡 Daftar REST API Endpoints

### Autentikasi (Publik)
* `POST /api/register` — Mendaftarkan akun baru (`username`, `email`, `password`)
* `POST /api/login` — Masuk ke akun (`username`, `password`) dan menerima HTTP-only Cookie
* `POST /api/logout` — Keluar dan membersihkan cookie sesi

### Catatan Pengeluaran (Wajib Login)
* `GET /api/me` — Melihat profil pengguna yang sedang login
* `GET /api/expenses` — Mengambil daftar pengeluaran (`?search=...&category=...`)
* `GET /api/expenses/summary` — Mengambil total nominal & jumlah transaksi
* `GET /api/expenses/:id` — Mengambil detail 1 catatan pengeluaran
* `POST /api/expenses` — Menambahkan catatan pengeluaran baru
* `PUT /api/expenses/:id` — Mengubah catatan pengeluaran yang sudah ada
* `DELETE /api/expenses/:id` — Menghapus catatan pengeluaran
