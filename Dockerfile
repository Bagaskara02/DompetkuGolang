# Stage 1: Build binary aplikasi Go
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Cache dependensi Go
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode proyek
COPY . .

# Kompilasi aplikasi menjadi binary mandiri (statically linked)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server main.go

# Stage 2: Runtime image yang sangat ringan
FROM alpine:3.20

# Install sertifikat CA untuk koneksi database aman (SSL/TLS)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Salin executable dari builder
COPY --from=builder /app/server .

# Salin folder web karena Gin menyajikan file frontend langsung dari ./web
COPY --from=builder /app/web ./web

# Expose port (Hostless akan mengatur port melalui environment variable PORT)
EXPOSE 8080

# Jalankan server
CMD ["./server"]
