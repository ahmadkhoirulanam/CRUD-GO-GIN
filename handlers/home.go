package handlers

import (
	"database/sql" // Untuk koneksi dan operasi database SQL
	"log"          // Untuk mencetak log error atau informasi
	"net/http"     // Untuk konstanta dan fungsi HTTP
	"path/filepath"

	"github.com/gin-gonic/gin"         // Framework web Gin
	_ "github.com/go-sql-driver/mysql" // Driver MySQL untuk Go, underscore berarti hanya untuk efek samping (init)
)

var db *sql.DB // Variabel global untuk koneksi database

// Fungsi init akan dipanggil secara otomatis saat package di-load
func init() {
	// DSN (Data Source Name) untuk koneksi ke database MySQL
	dsn := "root:@tcp(127.0.0.1:3306)/upgris"
	var err error
	// Membuka koneksi ke database
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err) // Jika gagal koneksi, hentikan program
	}
}

// Handler untuk menampilkan halaman login
func LoginPage(c *gin.Context) {
	// Render template login.html
	c.HTML(http.StatusOK, "login.html", nil)
}

// Handler untuk menangani form login
func LoginHandler(c *gin.Context) {
	// Ambil nilai dari form login
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validasi user ke database
	var user string
	err := db.QueryRow("SELECT username FROM users WHERE username = ? AND password = ?", username, password).Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			// Jika user tidak ditemukan
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
			return
		}
		// Error lain (misalnya query gagal)
		log.Println("Error checking credentials:", err)
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Internal server error"})
		return
	}

	// Jika login berhasil, arahkan ke halaman home
	c.Redirect(http.StatusFound, "/home")
}

// Handler untuk menampilkan halaman home beserta data produk
func HomePage(c *gin.Context) {
	// Query semua produk dari tabel (tambahkan kolom gambar)
	rows, err := db.Query("SELECT id, nama_produk, harga, gambar FROM produk")
	if err != nil {
		log.Println("Error fetching produk:", err)
		c.HTML(http.StatusInternalServerError, "home.html", gin.H{"error": "Gagal mengambil data produk"})
		return
	}
	defer rows.Close()

	// Struct lokal untuk menyimpan data produk
	type Produk struct {
		ID         int
		NamaProduk string
		Harga      float64
		Gambar     string
	}

	var produks []Produk

	for rows.Next() {
		var p Produk
		// Scan ke semua kolom termasuk gambar
		if err := rows.Scan(&p.ID, &p.NamaProduk, &p.Harga, &p.Gambar); err != nil {
			log.Println("Error scanning produk:", err)
			continue
		}
		produks = append(produks, p)
	}

	var jumlahProduk int
	err = db.QueryRow("SELECT COUNT(*) FROM produk").Scan(&jumlahProduk)
	if err != nil {
		log.Println("Error menghitung jumlah produk:", err)
		jumlahProduk = 0 // fallback
	}

	// Kirim data ke template
	c.HTML(http.StatusOK, "home.html", gin.H{
		"produks": produks,
		"jumlah1": jumlahProduk,
	})
}

// Handler untuk menambahkan produk baru
func TambahProduk(c *gin.Context) {
	// Ambil data dari form
	nama := c.PostForm("nama_produk")
	harga := c.PostForm("harga")

	// Ambil file gambar
	file, err := c.FormFile("gambar")
	if err != nil {
		log.Println("Gagal mengambil file gambar:", err)
		c.String(http.StatusBadRequest, "Gagal upload gambar")
		return
	}

	// Simpan file ke folder "uploads"
	filename := filepath.Base(file.Filename)
	path := filepath.Join("uploads", filename) // Pastikan folder "uploads" sudah ada
	if err := c.SaveUploadedFile(file, path); err != nil {
		log.Println("Gagal menyimpan file:", err)
		c.String(http.StatusInternalServerError, "Gagal menyimpan gambar")
		return
	}

	// Simpan data ke database (misalnya dengan nama file gambar)
	_, err = db.Exec("INSERT INTO produk (nama_produk, harga, gambar) VALUES (?, ?, ?)", nama, harga, filename)
	if err != nil {
		log.Println("Error inserting produk:", err)
	}

	// Redirect ke halaman home
	c.Redirect(http.StatusFound, "/home")
}

// Handler untuk mengedit produk
func EditProduk(c *gin.Context) {
	// Ambil data dari form edit
	id := c.PostForm("id")
	nama := c.PostForm("nama_produk")
	harga := c.PostForm("harga")

	// Update data produk berdasarkan ID
	_, err := db.Exec("UPDATE produk SET nama_produk = ?, harga = ? WHERE id = ?", nama, harga, id)
	if err != nil {
		log.Println("Error updating produk:", err)
	}

	// Redirect kembali ke halaman home
	c.Redirect(http.StatusFound, "/home")
}

// Handler untuk menghapus produk berdasarkan ID
func HapusProduk(c *gin.Context) {
	id := c.Param("id") // Ambil ID dari parameter URL

	// Hapus data dari database
	_, err := db.Exec("DELETE FROM produk WHERE id = ?", id)
	if err != nil {
		log.Println("Error deleting produk:", err)
	}

	// Redirect ke halaman home
	c.Redirect(http.StatusFound, "/home")
}
