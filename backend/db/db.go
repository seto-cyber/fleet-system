package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"sistem-manajemen-armada/backend/mqtt"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Struktur tabel penyimpanan lokasi kendaraan
type VehicleLocation struct {
	ID        uint   `gorm:"primaryKey"`
	VehicleID string `gorm:"index"`
	Latitude  float64
	Longitude float64
	Timestamp int64
}

// InitDB melakukan koneksi ke database dan auto migrate tabel
func InitDB() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Gagal membaca .env, menggunakan default environment variables")
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")


	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	var err error
	
	// untuk pengecekan 10 kali koneksi ke postgres
	for i := 1; i <= 10; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("✅ Terhubung ke PostgreSQL")

			// Lakukan migrasi otomatis
			log.Println("🚀 Melakukan migrasi tabel...")
			if err := DB.AutoMigrate(&VehicleLocation{}); err != nil {
				log.Fatalf("❌ Gagal migrasi tabel: %v", err)
			}

			log.Println("✅ Migrasi tabel selesai.")
			return
		}

		log.Printf("⏳ Gagal konek ke database (percobaan %d/10): %v", i, err)
		time.Sleep(3 * time.Second)
	}

	log.Fatalf("❌ Gagal konek ke database setelah 10 percobaan: %v", err)
}

// Tutup koneksi database saat aplikasi berhenti
func CloseDB() {
	if DB == nil {
		log.Println("⚠️ DB belum terhubung, tidak bisa ditutup.")
		return
	}
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("⚠️ Gagal mendapatkan koneksi SQL:", err)
		return
	}
	sqlDB.Close()
}

// Simpan data lokasi kendaraan ke tabel
func SaveLocation(loc mqtt.LocationMessage) {
	if DB == nil {
		log.Println("⚠️ DB belum terhubung, data tidak disimpan.")
		return
	}

	record := VehicleLocation{
		VehicleID: loc.VehicleID,
		Latitude:  loc.Latitude,
		Longitude: loc.Longitude,
		Timestamp: loc.Timestamp,
	}

	if err := DB.Create(&record).Error; err != nil {
		log.Printf("❌ Gagal menyimpan lokasi: %v", err)
	} else {
		log.Printf("✅ Lokasi kendaraan %s tersimpan di database.", loc.VehicleID)
	}
}
