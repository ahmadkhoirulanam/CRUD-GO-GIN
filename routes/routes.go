package routes

import (
	"framework/handlers"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(r *gin.Engine) {
	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Route untuk menampilkan halaman login
	r.GET("/", handlers.LoginPage)

	// Route untuk menangani login
	r.POST("/login", handlers.LoginHandler)

	// Route untuk menampilkan halaman home setelah login berhasil
	r.GET("/home", handlers.HomePage)
	r.POST("/produk/tambah", handlers.TambahProduk)
	r.POST("/produk/edit", handlers.EditProduk)
	r.GET("/produk/hapus/:id", handlers.HapusProduk)
}
