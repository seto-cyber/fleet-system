package mqtt

import (
	"encoding/json"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var broker = "tcp://mosquitto:1883"

// Structure data lokasi kendaraan
type LocationMessage struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

// Callback type
type LocationCallback func(LocationMessage)

// fungsi start subscriber
func StartSubscriber(callback LocationCallback) {
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("fleet-subscriber")
	client := mqtt.NewClient(opts)
	for {
	token := client.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Println("Gagal konek ke MQTT broker:", token.Error())
		log.Println("Coba lagi dalam 3 detik...")
		time.Sleep(3 * time.Second)
		continue
	}

	log.Println("✅ Berhasil konek ke MQTT broker!")
	break
	}

	topic := "/fleet/vehicle/+/location"
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		var loc LocationMessage
		if err := json.Unmarshal(msg.Payload(), &loc); err != nil {
			log.Println("Invalid message:", err)
			return
		}
		callback(loc)
	})

	log.Println("MQTT Subscriber started")
	select {}
}
