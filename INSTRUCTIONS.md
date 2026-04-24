# Panduan Sistem: Fix Fluktuasi & Manajemen Telegraf

File ini berisi dokumentasi cara setup dan penggunaan fitur baru di sistem Bandwith Monitor.

## 1. Fitur Perbaikan Fluktuasi Data (Anti-Spike)
Sistem sekarang memiliki filter otomatis pada `ProcessTraffic` untuk membuang lonjakan data yang tidak realistis (misal: lonjakan 10Gbps pada link 1Gbps).

- **Threshold**: 5x Kapasitas terdaftar atau maksimal **2.500 Mbps** (Global Ceiling).
- **Hasil**: Data anomali akan dibuang dan dicatat di log backend sebagai `Warning`.

## 2. Fitur Manajemen Telegraf via Dashboard
Anda sekarang bisa menambah IP Agent SNMP langsung dari menu "Telegraf Settings" di GUI browser. 

### Persiapan di VPS (PENTING):
Agar fitur "Auto-Write" ke file `/etc/telegraf/telegraf.conf` berjalan, jalankan perintah berikut di VPS:

```bash
# 1. Beri izin tulis ke user 'noc' untuk file telegraf
sudo chown noc:noc /etc/telegraf/telegraf.conf

# 2. Izinkan reload tanpa password
sudo visudo
# Tambahkan baris ini di paling bawah:
noc ALL=(ALL) NOPASSWD: /usr/bin/systemctl reload telegraf
```

### Cara Penggunaan Fitur Baru:
1.  **Sync from File**: Jika Anda baru pertama kali memakai GUI, klik tombol **"Sync from File"**. Aplikasi akan menarik semua IP yang sudah ada di file config ke database.
2.  **Tambah IP Baru**: Klik **"Add SNMP Agent"**. Masukkan IP mikrotik/perangkat Anda.
3.  **Hapus IP**: (Fitur Delete akan otomatis mengupdate file config).
4.  **Auto-Update**: Setiap kali Anda menambah IP, aplikasi Go akan menulis ulang file `/etc/telegraf/telegraf.conf` menggunakan template produksi lengkap Anda dan mereload servis Telegraf.

## 3. Langkah Instalasi (Rebuild)
Setiap ada update kode di branch `settingtelegraf` atau `documentation`, lakukan ini di VPS:

```bash
# Update Kode
git pull origin settingtelegraf

# Build Frontend
cd frontend && npm install && npm run build

# Restart Backend
cd ..
go run cmd/app/main.go
```

---
*Dokumentasi ini dibuat untuk memastikan operasional NMS berjalan lancar.*
