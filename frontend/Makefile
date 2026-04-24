# Makefile untuk mempermudah menjalankan Frontend via Docker

# Konfigurasi Nama
IMAGE_NAME = frontend-app
CONTAINER_NAME = frontend-container
PORT = 8085

.PHONY: all build run stop clean restart logs

# Target default jika hanya mengetik 'make'
all: build run

# 1. Build image Docker
build:
	@echo "=> Membangun (Building) Docker image..."
	docker build -t $(IMAGE_NAME) .

# 2. Menjalankan container (Hanya jika belum berjalan)
run:
	@echo "=> Menjalankan container di port $(PORT)..."
	docker run -d -p $(PORT):8085 --name $(CONTAINER_NAME) $(IMAGE_NAME)
	@echo "============================================="
	@echo "  Frontend dapat diakses di: http://localhost:$(PORT)"
	@echo "============================================="

# 3. Menghentikan & menghapus container yang sedang berjalan
stop:
	@echo "=> Menghentikan container yang sedang berjalan..."
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

# 4. Restart (Stop lalu jalankan ulang container)
restart: stop run

# 5. Melihat log dari Nginx / Frontend
logs:
	docker logs -f $(CONTAINER_NAME)

# 6. Menghapus (Clean) container beserta image-nya
clean: stop
	@echo "=> Menghapus Docker image..."
	docker rmi $(IMAGE_NAME) || true
