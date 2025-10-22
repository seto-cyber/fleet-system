package geofence

import (
	"encoding/json"
	"log"
	"math"
	"time"

	mq "github.com/streadway/amqp"
	"sistem-manajemen-armada/backend/mqtt"
)

var (
	rabbitURL     = "amqp://guest:guest@rabbitmq:5672/"
	geofencePoint = struct{ Lat, Lon float64 }{-6.2088, 106.8456}
	radiusMeters  = 50.0

	conn *mq.Connection
	ch   *mq.Channel
)

// Distance menghitung jarak antar koordinat (Haversine formula)
func distance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// InitRabbitMQ menyiapkan koneksi dan exchange
func InitRabbitMQ() {
	var err error
	for i := 1; i <= 10; i++ {
		conn, err = mq.Dial(rabbitURL)
		if err == nil {
			ch, err = conn.Channel()
			if err == nil {
				err = ch.ExchangeDeclare(
					"fleet.events",
					"fanout",
					true, false, false, false, nil,
				)
				if err == nil {
					log.Println("✅ Terhubung ke RabbitMQ dan exchange siap.")
					return
				}
			}
		}
		log.Printf("⏳ Gagal konek RabbitMQ (percobaan %d/10): %v", i, err)
		time.Sleep(3 * time.Second)
	}
	log.Fatalf("❌ Gagal konek ke RabbitMQ setelah 10 percobaan: %v", err)
}

// Check apakah kendaraan masuk geofence
func Check(loc mqtt.LocationMessage) {
	if ch == nil {
		log.Println("⚠️ RabbitMQ belum siap, lewati geofence check.")
		return
	}

	if distance(loc.Latitude, loc.Longitude, geofencePoint.Lat, geofencePoint.Lon) <= radiusMeters {
		body, _ := json.Marshal(map[string]interface{}{
			"vehicle_id": loc.VehicleID,
			"event":      "geofence_entry",
			"location": map[string]float64{
				"latitude":  loc.Latitude,
				"longitude": loc.Longitude,
			},
			"timestamp": loc.Timestamp,
		})

		err := ch.Publish("fleet.events", "", false, false, mq.Publishing{
			ContentType: "application/json",
			Body: body,
		})
		if err != nil {
			log.Println("❌ Gagal publish event:", err)
		} else {
			log.Printf("📡 Geofence event sent for %s\n", loc.VehicleID)
		}
	}
}

// Worker untuk membaca queue geofence_alerts
func StartWorker() {
	if ch == nil {
		log.Println("⚠️ RabbitMQ belum siap, worker tidak bisa dijalankan.")
		return
	}

	q, err := ch.QueueDeclare(
		"geofence_alerts",
		true, false, false, false, nil,
	)
	if err != nil {
		log.Println("❌ Gagal deklarasi queue:", err)
		return
	}

	// Bind queue ke exchange
	err = ch.QueueBind(
		q.Name,        // queue
		"",            // routing key kosong untuk fanout
		"fleet.events", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Println("❌ Gagal bind queue:", err)
		return
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Println("❌ Gagal consume:", err)
		return
	}

	go func() {
		for d := range msgs {
			log.Println("📨 Received geofence alert:", string(d.Body))
		}
	}()
}
