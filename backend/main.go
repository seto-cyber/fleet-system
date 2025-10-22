package main

import (
	"log"
	"time"

	"sistem-manajemen-armada/backend/api"
	"sistem-manajemen-armada/backend/db"
	"sistem-manajemen-armada/backend/geofence"
	"sistem-manajemen-armada/backend/mqtt"
)

func main() {
	// Inisialisasi koneksi database
	db.InitDB()
	defer db.CloseDB()

	// Inisialisasi koneksi RabbitMQ
	geofence.InitRabbitMQ()

	// Jalankan worker untuk menerima event geofence (RabbitMQ consumer)
	go geofence.StartWorker()

	// Jalankan subscriber MQTT
	go mqtt.StartSubscriber(func(loc mqtt.LocationMessage) {
		db.SaveLocation(loc)  // Simpan ke PostgreSQL
		loc1 := mqtt.LocationMessage{
    VehicleID: "B1234XYZ",
    Latitude:  -6.2088,
    Longitude: 106.8456,
    Timestamp: time.Now().Unix(),
}
		geofence.Check(loc1) // Cek apakah kendaraan masuk geofence
	})

	// Jalankan REST API (Fiber)
	app := api.SetupRouter()
	log.Fatal(app.Listen(":8080"))
}
