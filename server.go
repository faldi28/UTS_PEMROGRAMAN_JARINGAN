// server.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients     = make(map[*websocket.Conn]bool) // Menyimpan koneksi WebSocket aktif
	broadcast   = make(chan Donation)            // Channel untuk mengirim data donation ke semua klien
	upgrader    = websocket.Upgrader{}           // WebSocket upgrader untuk mengupgrade koneksi HTTP menjadi WebSocket
	clientMutex = sync.Mutex{}                   // Mutex untuk mengamankan akses ke map 'clients' dari banyak goroutine
	balances    = make(map[string]float64)       // Menyimpan saldo masing-masing pengguna
)

type Donation struct {
	From    string  `json:"from"`    // Nama pengirim
	Amount  float64 `json:"amount"`  // Jumlah donasi
	Message string  `json:"message"` // Pesan dari pengirim
}

func main() {
	// Menjalankan server pada berbagai protokol secara bersamaan (TCP, UDP, WebSocket)
	go handleTCP()                               // Menangani koneksi TCP
	go handleUDP()                               // Menangani koneksi UDP
	go handleWebSocket()                         // Menangani koneksi WebSocket
	http.HandleFunc("/ws", wsHandler)            // Menangani request WebSocket di endpoint "/ws"
	log.Println("Server started on port :8080")  // Menampilkan pesan server sudah mulai berjalan
	log.Fatal(http.ListenAndServe(":8080", nil)) // Menjalankan server HTTP di port 8080
}

// Fungsi untuk menangani koneksi WebSocket
func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // Memungkinkan koneksi dari origin manapun
	conn, err := upgrader.Upgrade(w, r, nil)                          // Upgrade koneksi HTTP menjadi WebSocket
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	clientMutex.Lock()
	clients[conn] = true // Menambahkan koneksi WebSocket yang baru ke dalam map clients
	clientMutex.Unlock()

	// Mendengarkan dan menerima donation dari klien WebSocket
	for {
		var donation Donation
		err := conn.ReadJSON(&donation) // Membaca data JSON dari WebSocket
		if err != nil {
			log.Printf("Client disconnected: %v", err)
			clientMutex.Lock()
			delete(clients, conn) // Menghapus koneksi yang terputus dari map clients
			clientMutex.Unlock()
			conn.Close() // Menutup koneksi WebSocket
			break
		}
		broadcast <- donation // Mengirim donation yang diterima ke channel broadcast
	}
}

// Fungsi untuk menangani koneksi TCP
func handleTCP() {
	listener, _ := net.Listen("tcp", ":8081") // Menunggu koneksi TCP pada port 8081
	defer listener.Close()                    // Menutup listener saat selesai
	for {
		conn, _ := listener.Accept() // Menerima koneksi dari klien
		go func(c net.Conn) {
			defer c.Close() // Menutup koneksi TCP setelah selesai

			var donation Donation
			reader := bufio.NewReader(c)         // Membaca data dari koneksi TCP
			line, err := reader.ReadString('\n') // Membaca input baris per baris
			if err != nil {
				log.Println("Error reading from TCP client:", err)
				return
			}

			// Memisahkan input berdasarkan spasi
			parts := strings.Fields(line)
			if len(parts) < 3 {
				log.Println("Invalid input format:", line)
				return
			}

			// Memproses pengirim dan jumlah donasi
			donation.From = parts[0]
			donation.Amount, err = strconv.ParseFloat(parts[1], 64) // Mengonversi jumlah donasi dari string ke float
			if err != nil {
				log.Println("Invalid amount format:", parts[1])
				return
			}

			// Menggabungkan sisa bagian sebagai pesan
			donation.Message = strings.Join(parts[2:], " ")

			// Jika pesan adalah "TOP_UP", top-up saldo pengguna
			if donation.Message == "TOP_UP" {
				balances[donation.From] += donation.Amount
			} else {
				balances[donation.From] -= donation.Amount
				broadcast <- donation // Mengirim donasi ke channel broadcast untuk diteruskan ke WebSocket
			}
		}(conn)
	}
}

// Fungsi untuk menangani koneksi UDP
func handleUDP() {
	addr, _ := net.ResolveUDPAddr("udp", ":8082") // Menentukan alamat UDP
	conn, _ := net.ListenUDP("udp", addr)         // Menunggu koneksi UDP pada port 8082
	defer conn.Close()                            // Menutup koneksi UDP setelah selesai
	for {
		buffer := make([]byte, 1024)                 // Menyediakan buffer untuk menerima data UDP
		n, remoteAddr, _ := conn.ReadFromUDP(buffer) // Membaca data dari UDP
		message := string(buffer[:n])                // Mengonversi data menjadi string

		// Memisahkan perintah dengan spasi
		parts := strings.Fields(message)
		if len(parts) < 2 {
			continue
		}
		username := parts[0]
		command := parts[1]

		// Perintah untuk memeriksa saldo
		if command == "CHECK_BALANCE" {
			balance := balances[username]
			response := fmt.Sprintf("Saldo anda saat ini: %.2f", balance)
			conn.WriteToUDP([]byte(response), remoteAddr) // Mengirimkan saldo ke klien UDP
		}
	}
}

// Fungsi untuk menangani broadcast donasi ke semua klien WebSocket
func handleWebSocket() {
	for {
		donation := <-broadcast // Menunggu donasi yang diterima
		clientMutex.Lock()
		// Mengirimkan donasi ke semua klien WebSocket
		for client := range clients {
			err := client.WriteJSON(donation) // Mengirim data donasi ke WebSocket client
			if err != nil {
				log.Printf("Error broadcasting to client: %v", err)
				client.Close()          // Menutup koneksi WebSocket jika terjadi kesalahan
				delete(clients, client) // Menghapus klien yang terputus
			}
		}
		clientMutex.Unlock()
	}
}
