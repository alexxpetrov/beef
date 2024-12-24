# Beef - A Chat Application Proof of Concept

Beef is a proof-of-concept chat application showcasing a robust, scalable, and modern tech stack. The project demonstrates best practices in microservice architecture and distributed systems while utilizing cutting-edge tools to optimize performance, reliability, and scalability.

This repository contains the backend services and infrastructure for Beef, designed with:

- **Golang** for the backend.
- **PostgreSQL** for chat room management and metadata.
- **ScyllaDB** (a Cassandra alternative) for chat messages.
- **Kafka** as a message queue.
- **Redis** for caching.
- **Docker Compose** for streamlined local and production deployments.

## Features
- **Golang Backend**: Efficient, lightweight, and built for concurrency.  
- **ScyllaDB Replication**: Production replication factor of 3 ensures fault tolerance.  
- **PostgreSQL Metadata**: Reliable and structured metadata storage.  
- **Message Queue with Kafka**: Decouples message handling and ensures reliability.  
- **Redis Caching**: Accelerates chat application performance.  
- **Authentication via [Identia](https://github.com/alexxpetrov/identia-be)**: Seamless authentication with support for biometry and traditional credentials.
- **Logging with [Erdtree](https://github.com/alexxpetrov/erdtree)**: Centralized logging for easier debugging and monitoring.

## Why This Tech Stack?

### **PostgreSQL**
- PostgreSQL is an **ACID-compliant** relational database, ensuring data consistency and reliability.  
- It excels at handling structured data like chat rooms and user metadata.  
- **Indexing, strong query support**, and extensive features like JSONB make PostgreSQL a perfect choice for chat room metadata management.  
- Its ability to scale horizontally and vertically ensures smooth operation as data grows.

### **ScyllaDB**
- ScyllaDB is a **NoSQL database** compatible with Cassandra but offers significantly better performance due to its **low-latency architecture** and C++ implementation.  
- Ideal for high-throughput workloads like chat messages where speed and scalability are critical.  
- **Replication factor of 3 in production** ensures fault tolerance, while local development uses a single replica for simplicity.  
- Offers automatic sharding, making it perfect for distributed environments.

### **Kafka**
- Kafka is a **distributed event-streaming platform** that ensures reliable and asynchronous message delivery.  
- Messages are written to Kafka topics, making them available to multiple consumers (e.g., notification services, archiving).  
- Built-in **scalability and fault tolerance** make Kafka a natural fit for high-throughput systems like chat applications.

### **Redis**
- Redis serves as a **caching layer** to minimize database queries and enhance performance.  
- Perfect for caching recent messages and frequently accessed chat-related data.  
- Built-in support for **Pub/Sub** could be leveraged for real-time updates.

### **Microservice Architecture**
- The backend is composed of multiple microservices, ensuring better modularity and fault isolation.  
- **Docker Compose** is used to orchestrate services locally and in production environments, with configurations tailored for both.

## How to Launch Locally

1. Ensure **Docker Desktop** is running.  
2. Run the following commands in your terminal:

```bash
# Start all containers
make docker-local

# Wait for containers to initialize, then create PostgreSQL tables
make migrate-up

# Set up ScyllaDB keyspace and tables
make scylla-migrate-up

# Create Kafka topics
make create-kafka-topic-local
```

3. For additional commands, refer to the `Makefile` in the root directory.

## Dependencies
### **Identia**  
Handles all authentication-related workflows (registration, login, WebAuthn).  
**GitHub**: [Identia](https://github.com/alexxpetrov/identia-be)

### **Erdtree**  
Manages centralized logging and provides Redis caching.  
**GitHub**: [Erdtree](https://github.com/alexxpetrov/erdtree)

## Live Demo
Try out the application: [Beef Demo](https://www.alexspetrov.com/chat)

---
