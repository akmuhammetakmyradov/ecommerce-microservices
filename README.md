# 🛒 E-Commerce Microservices Backend

## 📌 Overview

This project is a **microservices-based backend** for an e-commerce platform, implemented in **Go**.
It contains two main services and supporting infrastructure:

- **Cart Service** – Manages shopping cart operations.
- **Stocks Service** – Manages product inventory and pricing.

The system uses a **clean architecture** approach and supports:

- **gRPC** + **gRPC-Gateway** (REST)
- **PostgreSQL** for data storage
- **Kafka** for event streaming
- **Prometheus + Grafana** for monitoring
- **Jaeger** for distributed tracing

---

## 🏗 Architecture

![Architecture Diagram](docs/img/General%20Project%20Architecture.png)

**Key Features:**

- Layered architecture (internal services, repositories, delivery)
- Communication via **gRPC** (with optional HTTP REST gateway)
- Observability with **logging, tracing, and metrics**
- Dockerized deployment for dev & prod
- Makefile automation for build, test, and lint

---

## 🚀 Getting Started

### **Clone & Setup**

```bash
git clone https://github.com/akmuhammetakmyradov/ecommerce-microservices.git
cd ecommerce-microservices
```

### **Build All Services**

```bash
make build
```

### **Run All Services**

```bash
make run
```

Starts `cart` and `stocks` services with dependencies.

### **Run Linter**

```bash
make lint
```

---

## 🌐 Service Endpoints

| Service | Port | Protocols   |
| ------- | ---- | ----------- |
| Cart    | 8080 | gRPC + REST |
| Stocks  | 8081 | gRPC + REST |

**Example REST call via gRPC-Gateway:**

```bash
curl "http://localhost:8080/v1/cart?id=123"
```

**Example gRPC call:**

```bash
grpcurl -plaintext localhost:9090 list
```

---

## 🧪 Testing

Run **unit tests**:

```bash
make test
```

Run **integration tests**:

```bash
INTEGRATION_TEST=1 make integration-test
```

---

## 📦 Deployment with Docker

**Build & run with Docker Compose:**

```bash
docker-compose up --build
```

Each service has its own Dockerfile and `docker-compose.yml` for development and production.

---

## 📊 Observability

- **Logging** – Structured logs with Zap
- **Metrics** – Prometheus counters & histograms
- **Tracing** – Distributed tracing with Jaeger

---

## 🛠 Tech Stack

- **Language:** Go
- **Frameworks:** gRPC, gRPC-Gateway
- **Database:** PostgreSQL
- **Messaging:** Kafka
- **Monitoring:** Prometheus, Grafana
- **Tracing:** OpenTelemetry, Jaeger
- **Deployment:** Docker, Docker Compose

---

## 📄 License

This project was developed as part of a **Bootcamp Final Project** and is open for learning and portfolio purposes.

---