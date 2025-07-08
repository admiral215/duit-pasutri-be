# 💰 Duit Pasutri

**Duit Pasutri** adalah aplikasi pencatatan keuangan untuk pasangan suami istri agar lebih transparan dan teratur dalam mengelola pemasukan dan pengeluaran rumah tangga.

> Project ini dibuat sebagai side-project pribadi, sekaligus eksplorasi teknologi Golang, GraphQL, dan PostgreSQL.

---

## ✨ Fitur

- 👤 Autentikasi User
- 👥 Role (admin, user)
- 💳 Pembuatan dompet pribadi / bersama
- 📌 Mencatat transaksi (pemasukan / pengeluaran)
- 📅 Pencatatan dengan tanggal dan deskripsi
- 🔍 Query fleksibel via GraphQL
- 📊 Siap dikembangkan ke sisi frontend (web/mobile)

---

## 🛠️ Teknologi yang Digunakan

| Komponen     | Teknologi        |
|--------------|------------------|
| Bahasa       | Go (Golang)      |
| API Layer    | gqlgen (GraphQL) |
| ORM          | GORM             |
| Database     | PostgreSQL       |


---

## 🗂️ Struktur Proyek

```
duit-pasutri/
├── graph/              # Skema dan resolver GraphQL
├── internal/           # Modul internal untuk config, logic, dan database
├── gqlgen.yml          # Konfigurasi gqlgen
├── go.mod              # Dependensi Go
├── server.go           # Entry point aplikasi
└── README.md           # Penjelasan projek
```

---

## ⚙️ Instalasi dan Menjalankan

### 1. Clone Repositori
```bash
git clone https://github.com/admiral215/duit-pasutri-be.git
cd duit-pasutri
```

### 2. Buat File `.env`
Contoh :
```env
DATABASE_URL="host=localhost user=postgres password=yourpassword dbname=duit_pasutri port=5432 sslmode=disable TimeZone=Asia/Jakarta"
JWT_SECRET="jwt-secret-key"
PORT="8000"
```

### 3. Jalankan Aplikasi
```bash
go run server.go
```

Aplikasi akan berjalan di `http://localhost:8000/query`  
Kunjungi playground GraphQL di `http://localhost:8000`

---

## 🗺️ Rencana Ke Depan

- [ ] Pembuatan frontend (pertimbangan: Svelte / Flutter)
- [ ] Fitur OCR untuk scan struk belanja

---

## 🤝 Kontribusi

Pull request dan saran terbuka lebar. Silakan fork dan kirimkan PR jika kamu ingin berkontribusi 🙌

---


## 👨‍💻 Dibuat oleh

[Rully Admiral](https://www.linkedin.com/in/rully-admiral-1aa315287/)  
Software Engineer