# Go Microservices Architecture: gRPC & GraphQL

This project demonstrates a robust microservices architecture using **gRPC** for inter-service communication and **GraphQL** as the API gateway. It includes services for account management, product catalog, and order processing.

---

## Table of Contents

- [Overview](#overview)
- [Project Structure](#project-structure)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [gRPC Protobuf Generation](#grpc-protobuf-generation)
- [References](#references)
- [Contributing](#contributing)

---

## Overview

This repository showcases a production-style Go microservices system, designed for stability and maintainability. The project uses older, proven versions of Go and dependencies to minimize issues from frequent upstream changes—an approach recommended for solo developers or small teams.

---

## Project Structure

- **Account Service**  
  - Handles user account management.  
  - **Database:** PostgreSQL

- **Catalog Service**  
  - Manages the product catalog.  
  - **Database:** Elasticsearch

- **Order Service**  
  - Processes customer orders.  
  - **Database:** PostgreSQL

- **GraphQL API Gateway**  
  - Exposes a unified API for clients and routes requests to the appropriate microservices.

---

## Tech Stack

- **Go** (Golang)
- **gRPC** (Inter-service communication)
- **GraphQL** (API Gateway)
- **PostgreSQL** (Account & Order services)
- **Elasticsearch** (Catalog service)
- **Docker Compose** (Service orchestration)

---

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd <project-directory>
```

### 2. Start Services with Docker Compose

```bash
docker-compose up -d --build
```

### 3. Access the GraphQL Playground

Visit: [http://localhost:8000/playground](http://localhost:8000/playground)

---

## gRPC Protobuf Generation

To generate Go code from your `.proto` files, follow these steps:

1. **Install Protocol Buffers Compiler**

   ```bash
   wget https://github.com/protocolbuffers/protobuf/releases/download/v23.0/protoc-23.0-linux-x86_64.zip
   unzip protoc-23.0-linux-x86_64.zip -d protoc
   sudo mv protoc/bin/protoc /usr/local/bin/
   ```

2. **Install Go Plugins for Protobuf and gRPC**

   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

3. **Update Your PATH**

   ```bash
   export PATH="$PATH:$(go env GOPATH)/bin"
   source ~/.bashrc
   ```

4. **Prepare Your Project**

   - Create a `pb` folder in your project root.
   - In your `account.proto`, add:
     ```
     option go_package = "./pb";
     ```

5. **Generate gRPC Code**

   ```bash
   protoc --go_out=./ --go-grpc_out=./ account.proto
   ```

---

## References

- **Planning Drawboard:**  
  [Whimsical Board](https://whimsical.com/graphql-grpc-go-microservice-LdA8wTyHe3pUaUnEdH99cj)

- **Reference Code:**  
  [AkhilSharma90/go-grpc-graphql-microservices](https://github.com/AkhilSharma90/go-grpc-graphql-microservices)

---

## Contributing

Contributions are welcome! Please open issues or submit pull requests for improvements and bug fixes.

