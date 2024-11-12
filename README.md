UTS PEMEROGRAMAN JARINGAN

APP-SAWERIA SEDERHANA DENGAN GOLANG

1. Yohanes J Palis / 213400009
2. Link Youtube: https://youtu.be/qNSQLww5pY4
3. Langkah-Langkah Menjalankan Aplikasi
   1. Jalankan Perintah go run server.go di terminal untuk menjalankan server
      ![image](https://github.com/user-attachments/assets/b805594d-8e10-4a34-b211-0cb0aa2be709)
   2. Setelah server di jalankan lalu jalankan perintah go run client1.go ini untuk client yang menerima donasi
      ![image](https://github.com/user-attachments/assets/d1eee584-2070-455d-8df2-468955503001)
      
      Kemudian untuk yang index.HTML silakan jalankan dengan klik kanan pada mouse kemudian cari tulisan open with live server atau ALT+L ATAU ALT+O
   4. Setelah client1.go Di jalankan lalu jalankan perintah go run client.go, ini untuk client yang memberikan donasi
      ![image](https://github.com/user-attachments/assets/8cafd1b9-ea5f-48ec-a6bb-ad18ad8a7641)
      setelah itu silakan pilih menu masuk, lalu masukan username.
      ![image](https://github.com/user-attachments/assets/22af7ea8-b0fd-44bb-9dcf-83320a80b88b)
      
      setelah username di isi maka akan langsung masuk ke menu donasi.
      
      ![image](https://github.com/user-attachments/assets/98acb9c8-7c8d-41fd-9675-9b6499220eb3)
      di dalam menu donasi ada 3 menu utama yaitu, beri donasi, isi saldo,dan cek saldo.
      lakukan pengisian saldo dulu sebelum memberikan donasi, karena ketika pertama kali masuk maka saldo akan di defaut dari 0.

   5. setelah semua perintah di atas di lakukan maka coba untuk mengirim donasi dari client.go ke client1.go maka lihat apa yang akan terjadi di terminal client1.go.
      ![image](https://github.com/user-attachments/assets/99b487f8-da54-46fc-b981-90c45e2617fd)
4. Protokol-protokol yang digunakan
   1. tcp (Transmission Control Protocol)
      Fungsi:
      - Menerima donasi melalui koneksi TCP.
      - Memproses pesan seperti "TOP_UP" untuk memperbarui saldo pengguna.
      - Mengirimkan donasi ke WebSocket untuk distribusi ke klien.
   2. UDP (User Datagram Protocol)
      Fungsi:
      - Digunakan untuk Mengecek saldo.
      - Mengirimkan saldo pengguna dengan perintah CHECK_BALANCE.
   3. WEBSOCCET
      Fungsi:
      - Menerima donasi dari klien WebSocket.
      - Menyebarkan donasi ke semua klien yang terhubung secara langsung.
      - Menangani koneksi yang terputus dan menghapus klien yang tidak aktif.
      


      
      
      
   
   
