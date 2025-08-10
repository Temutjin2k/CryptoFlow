# CryptoFlow — Real-Time Cryptocurrency Market Tracker

CryptoFlow is a high-performance backend system with a modern lightweight dashboard that processes and visualizes real-time cryptocurrency market data using Go, Redis, PostgreSQL, and a vanilla JS frontend.

![Architecture](https://img.shields.io/badge/hexagonal-architecture-blue?style=for-the-badge)
![Go Version](https://img.shields.io/badge/go-blue?style=for-the-badge)
![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)
![PostgreSql](https://img.shields.io/badge/postgresql-4169e1?style=for-the-badge&logo=postgresql&logoColor=white)

---

## Overview

CryptoFlow collects cryptocurrency price data from multiple exchanges, aggregates it, caches it, and makes it accessible via a REST API. A beautiful built-in frontend lets users view real-time prices.

### Features
- **Hexagonal Architecture**
- Real-Time & Historical **Price Stats**
- Supports **Live/Test Mode** switching via API
- Uses **Redis** for real-time cache and **PostgreSQL** for aggregates
- Built-in **Vanilla HTML/CSS/JS Frontend**
- Dockerized and easy to run

---

## ⚙️ Installation

### Prerequisites
- Docker

### Getting started

```bash
git clone https://github.com/Temutjin2k/CryptoFlow.git
cd CryptoFlow

# Load Exchange Images
docker load -i exchanges/exchange1_amd64.tar
docker load -i exchanges/exchange2_amd64.tar
docker load -i exchanges/exchange3_amd64.tar

# Copy .env.example file to .env file
cp .env.example .env

docker compose up --build
```

## Postman Collection

You can find a ready-to-use Postman collection inside the `api/` directory:

```
api/CryptoFlow.postman_collection.json
```

Import this collection into Postman to quickly test all available endpoints including switching modes, fetching price stats, and running health checks.

---

## Frontend Features

* Real-time price stream from Redis
* Filter by exchange/symbol/metric/period
* Theme toggle: 🌞 Light & 🌚 Dark
* Simple vanilla HTML/CSS/JS

---

## Dashboard Screenshots

### Light Mode

<img width="1547" height="981" alt="image" src="https://github.com/user-attachments/assets/c0b2857b-08e6-4654-a6e8-0ed218127655" />

### Dark Mode
<img width="1549" height="989" alt="image" src="https://github.com/user-attachments/assets/9142ffd0-0fb1-426a-90d1-c91dbb43a0c6" />

---

## Concurrency Implementation and Patterns:
This project heavily relies on concurrency to handle large volumes of real-time data. Key concurrency patterns that were used include:

- **Fan-in:** Aggregating multiple market data streams into a single channel for centralized processing.

- **Fan-out:** Distributing incoming data updates to multiple workers to process them in parallel.

- **[Worker Pool](https://gobyexample.com/worker-pools):** Managing a set of workers that process live updates efficiently, ensuring balanced workload distribution.

- **Generator:** Implementing a generator to produce synthetic market data for `Test Mode`.



                      +---------------+       +---------------+       +---------------+
                      |  Source 1     |       |  Source 2     |       |  Source 3     |
                      |  (Generator)  |       |  (Generator)  |       |  (Generator)  |
                      +-------+-------+       +-------+-------+       +-------+-------+
                              |                       |                       |
                              v                       v                       v
                      +-------+-------+       +-------+-------+       +-------+-------+
                      |   Fan-Out 1   |       |   Fan-Out 2   |       |   Fan-Out 3   |
                      |  (Distributor)|       |  (Distributor)|       |  (Distributor)|
                      +---+---+---+---+       +---+---+---+---+       +---+---+---+---+
                          |   |   |               |   |   |               |   |   |
          +---------------+-+-+-+-+-+-------------+-+-+-+-+-+---------------------------+
          |               | | | | | |                 | | | | | |                       |
          v               v v v v v v                 v v v v v v                       v
      +---+---+       +---+---+---+---+-----+       +---+---+---+---+---+---+       +---+---+
      |Worker1|       |Worker2| ... |WorkerN|       |WorkerN+1| ... |WorkerM|      |WorkerM+1|
      +---+---+       +---+---+---+---+-----+       +---+---+---+---+---+---+       +---+---+
          |               |                   ...       |                   ...         |
          +---------------+-----------------------------+-------------------------------+
                              | (all output channels)
                              v
                      +-------+-------+
                      |    Fan-In     |
                      |  (Aggregator) |
                      +-------+-------+
                              | (resultCh)
                              v
                      +-------+-------+
                      |   Receiver    |
                      | (Collector)   |
                      +---------------+

--- 

## Authors

* **Meruyert** — Database, API, frontend
* **Temutjin** — System architecture, service design
--- 
