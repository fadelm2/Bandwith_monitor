# 🌐 WAN Monitoring System

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![React Version](https://img.shields.io/badge/React-19-61DAFB?style=for-the-badge&logo=react)](https://react.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)

A high-performance, real-time Wide Area Network (WAN) monitoring system designed to provide visibility into bandwidth usage, traffic patterns, and network capacity. Built with Go for high-concurrency data ingestion and React for a modern, responsive dashboard.

---

## ✨ Key Features

- **🚀 Real-time Traffic Ingestion**: Seamless integration with NATS for high-throughput network metrics processing.
- **📊 Interactive Dashboard**: Modern React 19 frontend with real-time data visualization.
- **🛡️ Secure Authentication**: JWT-based authentication system with RBAC support.
- **📉 Capacity Analysis**: Automated monitoring of WAN capacity thresholds (e.g., highlighting 99% usage).
- **🧪 Built-in Simulator**: Easy testing with a dedicated traffic simulation tool.
- **📖 API Documentation**: Integrated OpenAPI (Swagger) documentation.

## 🛠️ Tech Stack

### Backend
- **Core**: Go (Golang)
- **Framework**: [Fiber](https://gofiber.io/) (High-performance web framework)
- **Database**: MySQL with GORM ORM
- **Messaging**: [NATS](https://nats.io/) (Real-time data streaming)
- **Config**: Viper (Environment agnostic configuration)
- **Logging**: Logrus

### Frontend
- **Core**: React 19
- **Build Tool**: Vite
- **Styling**: Tailwind CSS / Vanilla CSS
- **Network**: Axios

---

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or higher
- Node.js & npm
- MySQL Server
- NATS Server

### Installation & Setup

1. **Clone the Repository**
   ```bash
   git clone https://github.com/fadelm2/Bandwith_monitor.git
   cd wan-system
   ```

2. **Backend Configuration**
   Edit `config.json` in the root directory:
   ```json
   {
     "database": {
       "username": "root",
       "password": "yourpassword",
       "host": "127.0.0.1",
       "port": 3306,
       "name": "network"
     },
     "nats": {
       "url": "nats://127.0.0.1:4222"
     },
     "app": {
       "port": 8082
     }
   }
   ```

3. **Frontend Configuration**
   Update `frontend/.env`:
   ```env
   VITE_API_URL=http://localhost:8082/internal
   VITE_AUTH_API_URL=http://localhost:8082/api
   ```

### Execution

#### Running the Backend
```bash
go run cmd/app/main.go
```
*API will be available at `http://localhost:8082`*

#### Running the Frontend
```bash
cd frontend
npm run dev
```
*Dashboard will be available at `http://localhost:5173` (or next available port)*

#### Running the Simulator
To start generating live network data:
```bash
go run cmd/simulate/main.go
```

---

## 🔐 Default Credentials (Seeded)

| Username | Password |
| :--- | :--- |
| `fadel` | `12345679` |
| `nanang` | `12345679` |

---

## 📅 Architecture Overview

The system follows **Clean Architecture** principles:
- **`cmd/`**: Entry points for the application and simulation tools.
- **`internal/config/`**: Configuration management.
- **`internal/delivery/`**: Delivery layers (HTTP Controllers & NATS Consumers).
- **`internal/usecase/`**: Business logic.
- **`internal/repository/`**: Data access layer.
- **`internal/model/`**: Input/Output models & Entities.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---
Developed by **Fadel Muhammad**.
