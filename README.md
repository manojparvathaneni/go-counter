# Go Counter Project 🧮

A minimal, idiomatic Go web service that keeps track of visits via a REST API. It stores data in a JSON file safely using atomic writes, Unix file locking, and periodic auto-saving.

This project is built entirely with the **Go standard library**, demonstrating best practices in file I/O, concurrency, configuration, and graceful shutdowns — all in one tiny microservice.

---

## ✨ Features

- ✅ `/api/counter` endpoint to increment and return visit count
- ✅ `sync/atomic` for safe counter access
- ✅ JSON file persistence
- ✅ **Atomic file writes** using `os.CreateTemp` + `os.Rename`
- ✅ **File locking** with `syscall.Flock` (Unix-safe)
- ✅ Graceful shutdown on SIGINT/SIGTERM
- ✅ Background auto-save loop (`COUNTER_AUTOSAVE_INTERVAL`)
- ✅ Command-line client to test the server
- ✅ Clean Go project layout with modules and tests

---

## 📁 Project Structure

```
go-counter/
│
├── cmd/
│   └── server/         # Main HTTP service
│       └── main.go
│
├── internal/
│   └── counter/        # Shared logic for JSON and file handling
│       └── counter.go
│
├── client/             # CLI tool to hit the API
│   └── main.go
│
├── tests/              # Unit tests for counter save/load
│   └── counter_test.go
│
├── counter.json        # Created at runtime
├── go.mod              # Go module file
└── README.md           # You're reading it!
```

---

## 🚀 Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/yourusername/go-counter.git
cd go-counter
```

### 2. Build and run the server

```bash
go run ./cmd/server
```

You’ll see:

```
Counter server running on http://localhost:8090
Auto-save interval: 10s
```

Then visit [http://localhost:8090/api/counter](http://localhost:8090/api/counter) in your browser or use the CLI tool.

---

## 🛠 Configuration

You can control the auto-save interval with an environment variable:

```bash
export COUNTER_AUTOSAVE_INTERVAL=5  # save every 5 seconds
go run ./cmd/server
```

---

## 🖥️ Using the CLI Tool

### Run from source:

```bash
go run ./client -times 3 -delay 500ms
```

### Flags:

- `-host`: URL of your API server
- `-times`: Number of requests to send
- `-delay`: Delay between requests (e.g. `1s`, `200ms`)

---

## 🧪 Running Tests

```bash
go test ./tests
```

This runs unit tests for counter saving and loading logic using isolated test files.

---

## 📘 Concepts Demonstrated

| Concept | Explanation |
|--------|-------------|
| **Atomic Counter** | Uses `sync/atomic` to safely increment visit count |
| **JSON Serialization** | Reads and writes `counter.json` using `encoding/json` |
| **Atomic File Write** | Writes to a temp file + renames to avoid corruption |
| **Unix File Locking** | Uses `syscall.Flock` for safe multi-process access (Unix-only) |
| **Graceful Shutdown** | Saves the counter on termination (SIGINT/SIGTERM) |
| **Auto-Save** | Background goroutine saves counter every N seconds |
| **HTTP Server** | Built with `net/http`, idiomatic handler setup |
| **CLI Testing** | Simulate load using Go-based CLI |
| **Testing** | `testing` + `httptest` used for isolated verification |

---

## 📖 Evolution of the Project

| Milestone | Description |
|----------|-------------|
| ✅ Initial counter | JSON read/write with basic file handling |
| ✅ API server | Added `/api/counter` |
| ✅ Atomic counter | Used `sync/atomic` for concurrency |
| ✅ Atomic save | Used temp file + rename |
| ✅ File locking | Added `syscall.Flock` for exclusive access |
| ✅ Graceful shutdown | Captures `SIGINT` and `SIGTERM`, final save |
| ✅ Auto-save loop | Background ticker with configurable interval |
| ✅ CLI tool | `client/main.go` to test the endpoint |
| ✅ Tests | Unit tests in `tests/` for `Save` and `Load` |

---

## 🙋‍♀️ Why This Project?

This repo was built to teach and demonstrate:

- How to use the Go standard library effectively
- Safe concurrency with files and memory
- Idiomatic Go layout and module structure
- Real-world microservice concerns (shutdown, state, I/O)

It's a clean starting point or reference for anyone learning Go.

---

## 🤝 Contributing

Suggestions and PRs welcome! To propose changes:

1. Fork the repo
2. Create a feature branch
3. Submit a pull request

---

## 📜 License

MIT © 2025 Manoj Parvathaneni
