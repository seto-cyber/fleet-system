package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	broker := "tcp://mosquitto:1883" // pakai nama service di docker-compose
	clientID := "fleet-publisher"

	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID)
	client := mqtt.NewClient(opts)

	// Retry loop kalau broker belum siap
	for {
		token := client.Connect()
		token.Wait()
		if token.Error() != nil {
			log.Println("Gagal konek ke MQTT broker:", token.Error())
			log.Println("Coba lagi dalam 3 detik...")
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	log.Println("✅ Publisher berhasil konek ke MQTT broker!")

	for {
		msg := map[string]interface{}{
			"vehicle_id": "B1234XYZ",
			"latitude":   -6.2088,
			"longitude":  106.8456,
			"timestamp":  time.Now().Unix(),
		}
		data, _ := json.Marshal(msg)
		token := client.Publish("/fleet/vehicle/B1234XYZ/location", 0, false, data)
		token.Wait()
		fmt.Println("📡 Published:", string(data))
		time.Sleep(2 * time.Second)
	}
}
