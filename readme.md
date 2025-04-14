# Product Service

## 1. Persiapan

Pastikan kamu sudah memiliki PostgreSQL yang terinstal di sistem atau menggunakan Docker untuk menjalankan PostgreSQL.

Jika menggunakan Docker, kamu bisa menjalankan PostgreSQL menggunakan perintah berikut:

```bash
docker run --name postgresql -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -v /var/lib/postgresql/data -d postgres
```

Untuk melakukan koneksi ke PostgreSQL, pastikan kamu menggunakan \`host=localhost\`, \`port=5432\`, \`user=postgres\`, dan \`password=postgres\`.

## 2. Database Name / Ekstensi UUID

Untuk mendukung penggunaan tipe data UUID di PostgreSQL, pastikan ekstensi \`uuid-ossp\` sudah diaktifkan. Jalankan perintah berikut pada database PostgreSQL untuk mengaktifkan ekstensi \`uuid-ossp\`:

```sql
CREATE DATABASE geb1_product_service;
```

```sql
CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";
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

## After created proto file

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

## How to generate from proto file

export PATH="$PATH:$(go env GOPATH)/bin"

protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative proto/product/product.proto
protoc --proto_path=proto --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/product/product.proto

protoc proto/product/product.proto \
  --go_out=. \
  --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative
