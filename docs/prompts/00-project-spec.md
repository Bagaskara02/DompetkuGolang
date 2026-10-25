# Project Specification — Vibe Coding

> **Tugas 1:** Isi spesifikasi aplikasi sebelum mulai melakukan vibe coding.
>
> Tentukan kebutuhan aplikasi dan pilih teknologi yang akan digunakan. Jangan mulai membuat kode sebelum spesifikasi selesai.

## 1. Informasi Aplikasi

**Nama aplikasi:**  
`[Isi nama aplikasi]`

**Deskripsi aplikasi:**  
`[Jelaskan tujuan aplikasi dalam 1–3 kalimat]`

**Target pengguna:**  
`[Siapa yang akan menggunakan aplikasi?]`

**Masalah yang ingin diselesaikan:**  
`[Jelaskan masalah yang ingin diselesaikan oleh aplikasi]`

## 2. Kebutuhan Aplikasi

Tuliskan fitur yang ingin dibuat.

1. `[Fitur utama 1]`
2. `[Fitur utama 2]`
3. `[Fitur utama 3]`
4. `[Fitur tambahan, jika ada]`

## 3. Pilihan Teknologi

Pilih satu opsi untuk setiap komponen yang diperlukan. Sesuaikan dengan kebutuhan aplikasi.

### A. Backend Framework

- [ ] 1. Go Standard Library (`net/http`)
- [ ] 2. Gin
- [ ] 3. Lainnya sebutkan

**Pilihan saya:** `[Isi nomor pilihan]`

### B. ORM (Object-Relational Mapping)

- [ ] 1. Tidak menggunakan ORM (SQL manual)
- [ ] 2. GORM
- [ ] 3. Lainnya sebutkan

**Pilihan saya:** `[Isi nomor pilihan]`

### C. Database

- [ ] 1. PostgreSQL
- [ ] 2. MySQL
- [ ] 3. SQLite

**Pilihan saya:** `[Isi nomor pilihan]`

### D. Frontend

- [ ] 1. HTML, CSS, JavaScript murni
- [ ] 2. React
- [ ] 3. Vue

**Pilihan saya:** `[Isi nomor pilihan]`

### E. Sistem Autentikasi

- [ ] 1. Tidak membutuhkan autentikasi
- [ ] 2. Session
- [ ] 3. JWT

**Pilihan saya:** `[Isi nomor pilihan]`

### F. Penyimpanan Sesi atau Token

- [ ] 1. Cookie
- [ ] 2. localStorage
- [ ] 3. Tidak diperlukan

**Pilihan saya:** `[Isi nomor pilihan]`

### G. API

- [ ] 1. REST API
- [ ] 2. Tidak membutuhkan API terpisah

**Pilihan saya:** `[Isi nomor pilihan]`

### H. Version Control

- [ ] 1. Git
- [ ] 2. Tidak menggunakan Git

**Pilihan saya:** `[Isi nomor pilihan]`

## 4. Detail Autentikasi

Isi bagian ini jika aplikasi membutuhkan autentikasi.

**Registrasi akun:** `[Ya / Tidak]`

**Login:** `[Ya / Tidak]`

**Logout:** `[Ya / Tidak]`

**Role pengguna:**
`[Contoh: Admin dan User / Tidak ada]`

**Aturan akses:**
`[Contoh: User hanya dapat melihat dan mengubah data miliknya sendiri]`

## 5. Batasan Pengembangan

Pilih aturan yang ingin diterapkan.

- [ ] Struktur folder sederhana.
- [ ] Kode mudah dipahami pemula.
- [ ] Tidak menambahkan teknologi di luar pilihan tanpa persetujuan.
- [ ] Memisahkan tampilan dan logika backend.
- [ ] Menggunakan environment variable untuk konfigurasi rahasia.
- [ ] Membuat pengujian dasar.
- [ ] Menjelaskan kode dan perubahan yang dibuat AI.

**Batasan tambahan:**
`[Isi jika ada]`

## 6. Fitur yang Tidak Dibutuhkan

Tuliskan fitur yang tidak perlu dibuat agar ruang lingkup proyek tetap jelas.

- `[Fitur yang tidak dibutuhkan 1]`
- `[Fitur yang tidak dibutuhkan 2]`

## 7. Pemeriksaan Sebelum Vibe Coding

- [ ] Nama dan tujuan aplikasi sudah jelas.
- [ ] Fitur utama sudah ditentukan.
- [ ] Pilihan teknologi sudah diisi.
- [ ] Kebutuhan autentikasi sudah ditentukan jika diperlukan.
- [ ] Batasan proyek sudah ditentukan.

**Status spesifikasi:** `[Belum siap / Siap untuk vibe coding]`

---

## 8. Instruksi untuk Agy

Setelah dokumen ini diisi, berikan instruksi berikut kepada Agy:

Baca dan pahami file `docs/prompts/00-project-spec.md`.

Gunakan pilihan teknologi dan kebutuhan yang sudah ditentukan sebagai acuan utama.

1. Periksa apakah spesifikasi sudah lengkap dan konsisten.
2. Tanyakan jika ada informasi penting yang belum jelas.
3. Usulkan struktur folder yang sederhana sesuai pilihan teknologi.
4. Jangan langsung membuat seluruh aplikasi sekaligus.
5. Tunggu persetujuan saya sebelum mulai implementasi.
6. Implementasikan aplikasi secara bertahap dan jelaskan cara menguji setiap bagian.

Jangan mengganti teknologi yang telah dipilih tanpa menjelaskan alasannya dan meminta persetujuan saya.
