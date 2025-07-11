# 📈 MarketFlow

> Real-Time Cryptocurrency Market Data Processing System

MarketFlow is a real-time backend system designed to ingest, process, cache, and serve cryptocurrency market data. Built with Go, Redis, and PostgreSQL, it simulates a real trading platform backend — efficiently handling concurrent data streams and exposing a REST API for querying prices and aggregated stats.

---

## 🌟 Features

- 🧠 **Hexagonal Architecture**  
  Cleanly separates domain, services, adapters, and ports.

- 🚀 **Real-Time Processing**  
  Supports both **Live Mode** (Docker exchanges) and **Test Mode** (synthetic generator).

- 🧵 **Concurrency Mastery**  
  Uses Go’s concurrency primitives: Worker Pool, Fan-in, Fan-out, and Generator patterns.

- ⚡ **Redis Caching**  
  Stores recent price updates for ultra-fast access.

- 🗃️ **PostgreSQL Storage**  
  Aggregates data every minute and persists it with min/avg/max stats.

- 📡 **REST API**  
  Query latest, highest, lowest, and average prices with filters by exchange and time period.

- ❤️ **Health Monitoring**  
  Exposes health-check and mode-switching endpoints.

---

## 🛠️ Setup

### 🧱 Requirements

- Go 1.22+
- Docker & Docker Compose
- Redis
- PostgreSQL

---

### 📦 Installation

Clone the repo:

```bash
git clone git@git.platform.alem.school:tkoszhan/marketflow.git
cd marketflow
````

---

### 🐳 Running Services

Start PostgreSQL, Redis, and the app:

```bash
docker-compose up --build
```

> ℹ️ Make sure to load and run the provided exchange Docker images beforehand. See **Live Mode Setup**.

---

### 🧪 Test Mode Setup

Run MarketFlow in **Test Mode** (synthetic data):

```bash
curl -X POST http://localhost:8080/mode/test
```

---

### 📡 Live Mode Setup

1. **Load Exchange Docker images**:

```bash
docker load -i exchange1_amd64.tar
docker load -i exchange2_amd64.tar
docker load -i exchange3_amd64.tar
```

2. **Run the exchanges**:

```bash
docker run -p 40101:40101 -d exchange1-arch
docker run -p 40102:40102 -d exchange2-arch
docker run -p 40103:40103 -d exchange3-arch
```

3. **Switch to Live Mode**:

```bash
curl -X POST http://localhost:8080/mode/live
```

---

## 🧪 Usage

```bash
./marketflow --help

Usage:
  marketflow [--port <N>]
  marketflow --help

Options:
  --port N     Port number
```

---

## 📖 API Endpoints

### 🔹 Latest Prices

* `GET /prices/latest/{symbol}`
* `GET /prices/latest/{exchange}/{symbol}`

### 🔸 Highest Prices

* `GET /prices/highest/{symbol}?period=1m`
* `GET /prices/highest/{exchange}/{symbol}?period=1m`

### 🔻 Lowest Prices

* `GET /prices/lowest/{symbol}?period=1m`
* `GET /prices/lowest/{exchange}/{symbol}?period=1m`

### 📊 Average Prices

* `GET /prices/average/{symbol}`
* `GET /prices/average/{exchange}/{symbol}`

### ⚙️ System

* `POST /mode/test` — switch to **Test Mode**
* `POST /mode/live` — switch to **Live Mode**
* `GET /health` — system health info

---

## 🔐 Configuration

MarketFlow reads configs from a YAML/JSON/TOML file. You must provide:

* PostgreSQL settings: host, port, user, pass, dbname
* Redis settings: host, port, password
* Exchange addresses: for test and live sources

---

## 📝 Logging

MarketFlow uses `log/slog` for structured logs.

Levels:

* `Info` — Startup, shutdown, routine ops
* `Warn` — No data, retries
* `Error` — Redis/Postgres failures, parsing errors

Example:

```json
{
  "time": "2025-07-10T17:29:10Z",
  "level": "WARN",
  "msg": "No prices found wait for 1 minute",
  "exchange": "exchange1",
  "symbol": "BTCUSDT"
}
```

---

## 🧼 Graceful Shutdown

Press `Ctrl+C` or send `SIGINT`, `SIGTERM` to cleanly shut down all services and close Redis/PostgreSQL connections.

---

## 🧑‍💻 Authors

Made with love and panic-free code by:

* **Meruyert**
* **Temutjin**

> "Concurrency is not parallelism — but with MarketFlow, you’ll master both."

---

## 💡 Resources

* [Go Concurrency Patterns](https://go.dev/blog/pipelines)
* [Worker Pool in Go](https://gobyexample.com/worker-pools)
* [Redis Docs](https://redis.io/docs/latest/)
* [Effective Go](https://go.dev/doc/effective_go#concurrency)

---

## 📁 Postman Collection

Postman collection for all API routes is available in the `/api` directory. Import it to test endpoints quickly.

---