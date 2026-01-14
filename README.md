# Real-time Study Room

A real-time web app where users can join a room with a code and:

- 💬 Chat in real time (WebSockets)
- 👥 See who’s currently online (“presence”)
- 🔗 Share links (optional)
- 📎 Send small files (optional)
- ✋ Use quick actions like “raise hand” / reactions (optional)

---

## Tech Stack

- **Backend:** Go (WebSockets + REST)
- **Frontend:** TypeScript (Vite)
- **Networking:** WebSocket protocol

---

## Getting Started

### Prerequisites
Make sure you have installed:
- [Go](https://go.dev/dl/)
- [Node.js (LTS)](https://nodejs.org/)
- Git

---

## Run Locally

### 1) Start the Go backend
In one terminal:

```bash
cd server
go mod download
go run .
```
---
### 2) Start the Typescript client
In another terminal:

```bash
cd web
npm install
npm run dev
```

.