---
# Go Simple API
`go-simple-api` adalah RESTful API sederhana yang dibangun menggunakan Golang, Gin Framework, dan GORM, dengan integrasi dokumentasi API menggunakan Swagger.

## 🧱 Struktur Data

### 🧑 User
Disimpan dalam tabel `tblUser`:
- `idUser` (string, primary key)
- `nameUser` (string)
- `emailUser` (string)
- `passwordUser` (string, hidden dari JSON response)
- `createUser` (timestamp, auto-create)
- `updateUser` (timestamp, auto-update)
- Relasi: memiliki banyak Wallet

### 💳 Wallet
Disimpan dalam tabel `tblWallet`:
- `idWallet` (string, primary key)
- `idUser` (string, foreign key)
- `addressWallet` (string)
- `createWallet` (timestamp, auto-create)
- `updateWallet` (timestamp, auto-update)

## 📦 Dependencies
- Gin Gonic (Web framework)
- GORM (ORM untuk Golang)
- Swaggo (Swagger documentation)
- JWT (untuk autentikasi)

## 📁 API Endpoint Structure

### 🔐 Auth
| Method | Endpoint             | Description         |
|--------|----------------------|---------------------|
| POST   | `/api/v1/auth/register` | Register user       |
| POST   | `/api/v1/auth/signin`  | Login & get token   |

### 👤 User (requires Auth)
| Method | Endpoint      | Description       |
|--------|---------------|-------------------|
| GET    | `/api/v1/user` | Get user profile  |
| PUT    | `/api/v1/user` | Update user data *(belum diimplementasi)* |

### 💳 Wallet (requires Auth)
| Method | Endpoint       | Description        |
|--------|----------------|--------------------|
| POST   | `/api/v1/wallet` | Add wallet *(belum diimplementasi)* |
| DELETE | `/api/v1/wallet` | Delete wallet *(belum diimplementasi)* |

### 📚 Swagger
| Method | Endpoint     | Description              |
|--------|--------------|--------------------------|
| GET    | `/swagger/*` | Swagger API Documentation |

## 🚀 Cara Menjalankan

1. Clone repositori ini:
   bash
   git clone https://github.com/username/go-simple-api.git
   cd go-simple-api
   
2. Install dependensi:
   bash
   go mod tidy
   
3. Jalankan aplikasi:
   bash
   go run main.go
   
4. Buka dokumentasi Swagger di:
   http://localhost:8080/swagger/index.html
   

## 🔐 Middleware
* `AuthMiddleware`: digunakan untuk mengamankan endpoint `/user` dan `/wallet`.

## 🗂️ Struktur Folder
src/
├── config/
├── controller/
│   └── api/
├── middleware/
├── model/
│   └── model_database/
├── routes/
└── main.go


## 📌 Catatan
* Pastikan environment variable seperti koneksi database dikonfigurasi di `config`.
* Endpoint `PUT /user`, `POST /wallet`, dan `DELETE /wallet` belum diimplementasikan secara penuh.
* Password dienkripsi (rekomendasi: gunakan bcrypt).

---

## 🧑‍💻 Kontributor
* Nur Wahid Azhar

## 📝 Lisensi
MIT License
---