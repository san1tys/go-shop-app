# GoShop — Fullstack E-commerce App (Go + React + Docker)

GoShop is a learning-oriented full-stack e-commerce project built with **Go**, **Postgres**, **React** и **Docker**.  
The entire application is containerized, easy to run with a single command, and ready for team collaboration.

---

## Tech Stack

### Backend (Go)
- Go 1.24.0+
- Gin Web Framework
- PostgreSQL (using lib/pq)
- JWT Authentication
- Docker containerization
- SQL migrations

### Frontend
- React + Vite
- Pnpm

### DevOps
- Docker
- Docker Compose
- Makefile

---

## How to Run the Project (One Command)

### 1. Install Docker
https://www.docker.com/products/docker-desktop/

### 2. Start the full project:

```bash
make dev

---

## Seed migration with demo data

- The file backend/migrations/002_seed.sql creates demo data:
    - an administrator account: admin@example.com with password admin123 and role admin;
    - several test products for the storefront.

