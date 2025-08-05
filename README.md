# Ticket

Welcome to **Ticket** – an Event-driven Microservice written in Go where each service is implemented based on Hexagonal architecture principles. This practice project is designed for simplicity and as an exercise with event driven architecture(not just a simple pub/sub), making it approachable for anyone wanting to learn,explore and contribute.

---

## 🚄 What is Ticket?

Ticket is a train ticket booking system built with an **event-driven architecture**. It allows users to register, book and cancel tickets, manage trains and thier seats, and receive notifications on specfic operations and so on.

---

## 🛠 Technologies Used

- **Go** (Golang): The entire backend is written in Go.
- **RESTful APIs**: Each domain (user, train, ticket) exposes clear rest API for operations utilize Fiber(it would be changen when API gateway completed. services will communicate via events and API gateway will change http requests to commands).
- **Nats/JetStream**: Utilized as event bus exposing robust features and also written in Go
- **Postgresql**: As database(Where each service has own database)
- **Protocol Buffers**: Used as event payload serialization format because of it's elegant performance.

---

## 📝 TODO List

1. **Add authentication & authorization** (secure endpoints)
2. **Improve error handling**
3. **Complete API Gateway**
4. **Implement Travel,Commnets and City services** 
5. **Comprehensive testing** (unit & integration)
6. **Notifications for booking tickets** 
7. **Front-end client** (React/Vue/Angular)
8. **Documentation** – expand with API specs and code walkthroughs

---

## 📦 Getting Started
### Prerequisites
- Go
- PostgreSQL
- Nats(JetStream enabled)

1. **Clone the repository**  
   `git clone https://github.com/AssassinRobot/Ticket.git`

2. **Run**  
   - See .env for configuration format and config based on provided format
   - `go run main.go` in each module
---

## 🙌 Contributing

Beginner or expert, all contributions are welcome! Open an issue or submit a pull request.
