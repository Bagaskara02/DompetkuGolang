# Project Specification — Vibe Coding

> **Tugas 1:** Isi spesifikasi aplikasi sebelum mulai melakukan vibe coding.
>
> Tentukan kebutuhan aplikasi dan pilih teknologi yang akan digunakan. Jangan mulai membuat kode sebelum spesifikasi selesai.

## 1. Informasi Aplikasi

**Nama aplikasi:**  
`DompetKu — Pencatat Pengeluaran Harian`

**Deskripsi aplikasi:**  
`Aplikasi web sederhana untuk membantu pengguna mencatat dan memantau pengeluaran keuangan harian mereka. Pengguna dapat menambahkan riwayat transaksi, melihat total pengeluaran, serta mengelola catatan keuangan pribadi dengan mudah dan cepat.`

**Target pengguna:**  
`Mahasiswa, pekerja, atau masyarakat umum yang ingin mencatat pengeluaran harian secara praktis tanpa fitur akuntansi yang rumit.`

**Masalah yang ingin diselesaikan:**  
`Banyak orang sering lupa ke mana uang mereka habis karena malas menggunakan aplikasi keuangan yang terlalu kompleks, sehingga membutuhkan tempat pencatatan pengeluaran yang cepat, simpel, dan langsung ke intinya.`

## 2. Kebutuhan Aplikasi

Fitur utama aplikasi:

1. `Autentikasi pengguna (Registrasi akun baru, Login, dan Logout)`
2. `CRUD Catatan Pengeluaran (Menambah, melihat daftar, mengedit, dan menghapus data pengeluaran berupa nama pengeluaran, nominal, kategori, dan tanggal)`
3. `Ringkasan Total Pengeluaran (Menampilkan jumlah total nominal pengeluaran secara otomatis di halaman utama)`
4. `Pencarian & Filter Kategori (Mencari catatan pengeluaran berdasarkan nama atau memfilter berdasarkan kategori seperti Makanan, Transportasi, dan Lainnya)`

**Halaman aplikasi:**

- Login
- Register
- Main Page — My Expenses (Halaman Utama Catatan Pengeluaran)

Pada Main Page, pengguna dapat melihat ringkasan total pengeluaran, daftar catatan pengeluaran, menambah catatan baru, mengedit, menghapus, dan logout.

## 3. Pilihan Teknologi

Pilih satu opsi untuk setiap komponen yang diperlukan. Sesuaikan dengan kebutuhan aplikasi. isiskan [x] hanya satu saja

### A. Backend Framework

- [ ] 1. Go Standard Library (`net/http`)
- [x] 2. Gin
- [ ] 3. Lainnya sebutkan

**Pilihan saya:** `2`

### B. ORM (Object-Relational Mapping)

- [ ] 1. Tidak menggunakan ORM (SQL manual)
- [x] 2. GORM
- [ ] 3. Lainnya sebutkan

**Pilihan saya:** `2`

**Ketentuan:** `Menggunakan GORM`

### C. Database

- [x] 1. PostgreSQL
- [ ] 2. MySQL
- [ ] 3. SQLite

**Pilihan saya:** `1`

### D. Frontend

- [x] 1. HTML, CSS, JavaScript murni
- [ ] 2. React
- [ ] 3. Vue

**Pilihan saya:** `1`

### E. Sistem Autentikasi

- [ ] 1. Session
- [x] 2. JWT

**Pilihan saya:** `1`

**Ketentuan:** `Gunakan JWT berbasis server untuk mengelola status login pengguna.`

### F. Penyimpanan Sesi atau Token

- [x] 1. Cookie
- [ ] 2. localStorage
- [ ] 3. Tidak diperlukan

**Pilihan saya:** 1

**Ketentuan:** Simpan ID session dalam Cookie dengan pengaturan keamanan yang sesuai.

### G. API

- [x] 1. REST API
- [ ] 2. Tidak membutuhkan API terpisah

**Pilihan saya:** 1

**Ketentuan:** Backend menyediakan REST API untuk autentikasi dan pengelolaan catatan. Frontend mengakses API menggunakan `fetch()`.

### H. Version Control

- [x] 1. Git
- [ ] 2. Tidak menggunakan Git

**Pilihan saya:** 1

## 4. Detail Autentikasi

**Registrasi akun:** Ya

**Login:** Ya

**Logout:** Ya

**Role pengguna:** User saja. Tidak ada role Admin.

**Aturan akses:**

1. Pengguna harus login untuk mengakses catatan.
2. Pengguna yang belum login hanya dapat mengakses halaman Login dan Register.
3. Setiap catatan harus memiliki `user_id` sebagai pemiliknya.
4. Backend harus memeriksa kepemilikan catatan pada setiap operasi melihat detail, mengedit, dan menghapus.
5. Pengguna tidak boleh melihat atau mengubah catatan pengguna lain, termasuk dengan mengubah ID pada URL atau request API.
6. Password harus disimpan dalam bentuk hash menggunakan algoritma yang sesuai, bukan sebagai teks biasa.
7. Setelah logout, session harus dinonaktifkan dan cookie session dibersihkan.

## 5. Konfigurasi Environment dan Database

**Ketentuan wajib:**

1. Semua konfigurasi aplikasi yang sensitif harus menggunakan environment variable.
2. Gunakan file `.env` untuk konfigurasi saat pengembangan lokal.
3. Sediakan file `.env.example` sebagai contoh konfigurasi tanpa kredensial asli.
4. File `.env` tidak boleh dimasukkan ke Git.
5. URL atau connection string PostgreSQL harus dibaca dari environment variable bernama `DATABASE_URL`.
6. Jangan menulis URL database, username, password database, session secret, atau kredensial lainnya secara langsung di kode.
7. Aplikasi harus memberikan pesan error yang jelas jika konfigurasi wajib belum tersedia, tanpa membocorkan password atau rahasia lainnya.

Contoh `.env.example`:

```dotenv
APP_PORT=8080
DATABASE_URL=postgres://username:password@localhost:5432/personal_notes?sslmode=disable
SESSION_SECRET=replace-with-a-random-secret
```

**Catatan:** Nilai di atas hanya contoh untuk pengembangan lokal. Gunakan kredensial dan session secret milik sendiri pada file `.env`. Untuk database produksi, gunakan konfigurasi SSL sesuai persyaratan penyedia database.

## 6. Batasan Pengembangan (BOLEH DIUBAH SESUAI KEINGINAN PESERTA PELATIHAN)

- [X] Struktur folder sederhana.
- [X] Kode mudah dipahami pemula.
- [X] Tidak menambahkan teknologi di luar pilihan tanpa persetujuan.
- [X] Memisahkan frontend dan backend.
- [X] Menggunakan environment variable untuk konfigurasi.
- [X] Membuat pengujian dasar.
- [X] Menjelaskan kode dan perubahan yang dibuat AI.

**Batasan tambahan:**
`Desain antarmuka (UI) dibuat bersih, minimalis, dan responsif hanya menggunakan CSS murni tanpa framework CSS berat agar mudah dipelajari.`

## 6. Fitur yang Tidak Dibutuhkan

Tuliskan fitur yang tidak perlu dibuat agar ruang lingkup proyek tetap jelas.

- `Grafik atau chart visualisasi statistik keuangan yang kompleks.`
- `Fitur ekspor laporan ke PDF, Excel, atau CSV.`
- `Manajemen multi-mata uang, koneksi rekening bank, atau pengingat tagihan otomatis.`
- `Fitur lupa password via email (reset password) atau verifikasi OTP.`

## 7. Pemeriksaan Sebelum Vibe Coding (BOLEH DIUBAH SESUAI KEINGINAN PESERTA PELATIHAN)

- [X] Nama dan tujuan aplikasi sudah jelas.
- [X] Fitur utama sudah ditentukan.
- [X] Pilihan teknologi sudah diisi.
- [X] Kebutuhan autentikasi sudah ditentukan.
- [X] Batasan proyek sudah ditentukan.
- [X] Penggunaan `.env` dan `DATABASE_URL` sudah ditentukan.

**Status spesifikasi:** Siap untuk diperiksa dan disetujui sebelum implementasi.

---

## 8. Instruksi untuk Agy

Baca dan pahami file `docs/prompts/00-project-spec.md` secara lengkap.

Gunakan spesifikasi ini sebagai acuan utama pengembangan.

1. Periksa apakah spesifikasi sudah lengkap dan konsisten.
2. Tanyakan jika ada informasi penting yang belum jelas.
3. Usulkan struktur folder sederhana sesuai teknologi yang dipilih.
4. Jelaskan fungsi folder dan file utama.
5. Jangan langsung membuat seluruh aplikasi sekaligus.
6. Jangan membuat, mengubah, atau menghapus file sebelum saya menyetujui usulan struktur folder.
7. Implementasikan aplikasi secara bertahap setelah mendapat persetujuan.
8. Jelaskan perubahan kode dan cara menguji setiap tahap.
9. Jangan mengganti teknologi atau menambahkan dependensi tanpa menjelaskan alasannya dan meminta persetujuan saya.
10. Pastikan `DATABASE_URL` dibaca dari environment variable dan tidak pernah di-hardcode.
11. Pastikan `.env` masuk ke `.gitignore`, sementara `.env.example` boleh disimpan di Git.
12. Jangan menampilkan nilai rahasia dari `.env` pada log atau respons API.
