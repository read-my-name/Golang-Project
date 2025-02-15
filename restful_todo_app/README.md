# Go Todo List API

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A robust RESTful Todo List API built with Go (Golang) featuring concurrent processing, CSV persistence, and modern API design patterns.

## 📌 Key Features

- **RESTful API Endpoints** with proper HTTP methods and status codes
- **Concurrent Processing** using Goroutines and Channels
- **CSV Storage** with thread-safe operations (mutex protected)
- **Structured Project Layout** following Go best practices
- **Advanced Todo Management**:
  - Status tracking (Not Started, In Progress, On Hold, Completed)
  - Priority levels (Low, Medium, High)
  - Labels/tags system
  - Subtasks support
  - Due dates and time-based filtering
- Request validation and error handling
- Comprehensive filtering system (status, priority, labels, time periods)

## 🚀 Installation

1. **Clone Repository**
```bash
git clone https://github.com/read-my-name/restful-todo-app.git
cd restful-todo-app
```

2. **Install Dependencies**

3. **Run Application**
```bash
go run cmd/server/main.go
```

4. **Create Data Directory**
```bash
mkdir -p data
```

## 📚 API Documentation

### Endpoints

| Method | Endpoint                | Description                     |
|--------|-------------------------|---------------------------------|
| GET    | /todos                  | List all todos                  |
| POST   | /todos                  | Create new todo                 |
| GET    | /todos/period/{period}  | Filter by period (today) |
| GET    | /todos/filter           | Filter by status/priority/labels|
| PUT    | /todos/{id}             | Update todo                     |
| DELETE | /todos/{id}             | Delete todo                     |

### Example Requests

**Create Todo**
```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "API Security Review",
    "description": "Complete security audit",
    "status": "in_progress",
    "priority": "high",
    "due_date": "2023-12-31T00:00:00Z",
    "labels": ["security", "critical"],
    "subtasks": [{
      "title": "Implement authentication",
      "completed": false
    }]
  }'
```

**Filter Todos**
```bash
# Get high priority security todos
curl "http://localhost:8080/todos?priority=high&label=security"

# # Get this week's todos
# curl http://localhost:8080/todos/period/week
# ```

## 🏗 Project Structure

```
restful-todo-app/
├── api/
│   └── handlers/      # HTTP request handlers
├── internal/
│   ├── service/       # Business logic and worker pools
│   └── storage/       # CSV persistence implementation
├── pkg/
│   └── models/        # Data structures and enums
└── cmd/
    └── server/        # Main application entry point
```

## 🛠 Technologies Used

- **Go** (Golang
- **Gin** Web Framework
- **UUID** for unique identifiers
- **CSV** for persistent storage
- **Standard Library**:
  - net/http
  - sync (mutex)
  - time
  - encoding/csv

## 🔒 Error Handling & Logging

- Structured error responses
- Concurrent-safe operations using sync.RWMutex
- Worker pool pattern for async processing
- Detailed request/response logging
- Graceful shutdown handling

<!-- ## 🔮 Future Enhancements

- [ ] SQL/Redis database integration
- [ ] User authentication (JWT/OAuth)
- [ ] Web frontend (React/Vue)
- [ ] Comprehensive test suite
- [ ] Docker containerization
- [ ] Rate limiting and API throttling -->

## 🤝 Contributing

Contributions welcome! Please follow the [Contributor Covenant](https://www.contributor-covenant.org/) code of conduct.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.