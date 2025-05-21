# Product Service

## Persiapan

Pastikan kamu sudah memiliki PostgreSQL yang terinstal di sistem atau menggunakan Docker untuk menjalankan PostgreSQL.

Jika menggunakan Docker, kamu bisa menjalankan PostgreSQL menggunakan perintah berikut:

```bash
docker run --name postgresql -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -v /var/lib/postgresql/data -d postgres
```

Untuk melakukan koneksi ke PostgreSQL, pastikan kamu menggunakan \`host=localhost\`, \`port=5432\`, \`user=postgres\`, dan \`password=postgres\`.

## 2. Database Name / Ekstensi UUID

Untuk mendukung penggunaan tipe data UUID di PostgreSQL, pastikan ekstensi \`uuid-ossp\` sudah diaktifkan. Jalankan perintah berikut pada database PostgreSQL untuk mengaktifkan ekstensi \`uuid-ossp\`:

```sql
CREATE DATABASE goescape_marketplace;
```

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

Ekstensi ini digunakan untuk menghasilkan UUID secara otomatis.

## 3. Persiapan Tabel products

```sql
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    description TEXT,
    qty INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 4. Persiapan untuk gRPC

### Instalasi protoc

Untuk menginstal `protoc` (Protocol Buffers Compiler), ikuti langkah-langkah berikut:

#### macOS (menggunakan Homebrew)

```bash
brew install protobuf
```

#### Linux

1. **Download Protoc**: Kunjungi [halaman rilis Protoc](https://github.com/protocolbuffers/protobuf/releases) dan unduh versi terbaru untuk sistem operasi kamu. Pilih file yang sesuai dengan sistem operasi kamu (misalnya, `protoc-[versi]-linux-x86_64.zip` untuk Linux).
2. **Ekstrak File**:

   ```bash
   unzip protoc-*-linux-x86_64.zip -d protoc.
   ```

3. **Pindahkan File**:

   ```bash
   sudo mv protoc/bin/protoc /usr/local/bin/
   sudo mv protoc/include/* /usr/local/include/
   ```

#### Windows

1. **Download Protoc**: Kunjungi [halaman rilis Protoc](https://github.com/protocolbuffers/protobuf/releases) dan unduh versi terbaru untuk Windows. Pilih file yang sesuai dengan sistem operasi kamu (misalnya, `protoc-[versi]-win64.zip` untuk Windows 64-bit).
2. **Ekstrak File**: Ekstrak file zip yang telah diunduh.
3. **Tambahkan ke PATH**: Tambahkan lokasi binary `protoc.exe` ke variabel lingkungan PATH.

### Verifikasi instalasi

Setelah menginstal `protoc`, verifikasi instalasi dengan menjalankan perintah berikut di terminal:

```bash
protoc --version
```

Jika instalasi berhasil, kamu akan melihat versi `protoc` yang terinstal.

### Instal Plugin Go untuk Protocol Buffers

#### Instal plugin untuk menghasilkan kode Go dari proto

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

#### Instal plugin untuk menghasilkan kode gRPC

```bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

#### Pastikan binary yang diinstal tersedia di PATH

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## 5. Generate kode Go dari file proto

Setelah membuat file proto (misalnya product.proto), gunakan perintah berikut untuk menghasilkan kode Go:

### Metode 1: Perintah terpisah

```bash
protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative proto/product/product.proto
protoc --proto_path=proto --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/product/product.proto
```

### Metode 2: Perintah gabungan (rekomendasi)

```bash
protoc proto/product/product.proto \
  --go_out=. \
  --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative
```

Perintah ini akan menghasilkan dua file:

- product.pb.go: berisi definisi pesan dan kode untuk serialisasi/deserialisasi.
- product_grpc.pb.go: berisi definisi layanan gRPC dan kode untuk memanggil layanan

### Metode 3: Generate semua file proto di dalam folder dan subfolder

```bash
protoc proto/**/*.proto \
  --go_out=. \
  --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative
```

Perintah ini akan menghasilkan semua file proto yang ada di dalam folder dan subfolder ke dalam folder yang sama dengan file proto tersebut.
