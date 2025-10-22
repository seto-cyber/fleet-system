## Sistem manajemen armada  

berbasis Golang, PostgreSQL, MQTT (Eclipse Mosquitto)  dan RabbitMQ dengan Docker.

---

## Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

---

## Jalankan Aplikasi

1. Clone repository:

```bash
git clone https://github.com/seto-cyber/fleet-system.git
cd fleet-system


---
2. Buat file .env:
   Di folder project, buat file bernama .env
   Isi contoh untuk project fleet-system:
   # Database
    POSTGRES_HOST=
    POSTGRES_USER=
    POSTGRES_PASSWORD=
    POSTGRES_DB=
    POSTGRES_PORT=

    sesuaikan setting postgres
---
3. Jalankan Docker
    pastikan docker sudah berjalan di komputer
    setalah docker berjalan, jalankan :
    ```bash
    docker compose build --no-cache
    docker compose up

4. Testing API
    setelah docker berjalan maka bisa testing API menggunakan postman

Contoh endpoint:

# Mendapatkan lokasi terakhir kendaraan
curl http://localhost:8080/vehicles/B1234XYZ/location

# Mendapatkan riwayat kendaraan
curl http://localhost:8080/vehicles/B1234XYZ/history?start=1715000000&end=1715009999

Sesuaikan B1234XYZ dengan Vehicle ID yang ada.

dengan README ini, semoga menjadi petunjuk dalam penggunakan aplikasi ini.
