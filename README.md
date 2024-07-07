# gRPC User Service

The "# gRPC User Service" repository you're referring to is a Golang-based service designed for managing user details using gRPC (Google's Remote Procedure Call). Here's a breakdown of its key components and functionalities based on the provided documentation:

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Build and Run](#build-and-run)
- [Accessing gRPC Endpoints](#accessing-grpc-endpoints)
- [Docker](#docker)

## Overview

The service primarily focuses on Read(aspect of CRUD functionality) operations for managing user details. Additionally, it includes search functionality based on specific criteria such as first name, city, phone number, and marital status.

## Prerequisites

Before running the application, ensure you have the following installed:

- Go (Golang) 1.19 or higher
- Docker (optional, for containerization)

## Build and Run

To build and run the application locally:

1. Clone the repository:

   ```bash
   git clone https://github.com/VinikaAnthwal/grpc-user-service.git
   cd grpc-user-service
   ```

2. Build the Docker image:

   ```bash
   docker build -t grpc-user-service .
   ```

   This command builds a Docker image named `grpc-user-service` using the Dockerfile provided in the repository.

3. Run the Docker container:

   ```bash
   docker run -p 8080:8080 grpc-user-service
   ```

   This command starts a Docker container from the `grpc-user-service` image, exposing port 8080 to access the gRPC server.

4. After the service is running in Docker, you can run the client application. For example, assuming your client code is located in client/main.go and it connects to port 8080:

   ```bash
   go run client/main.go
   ```

Accessing gRPC Endpoints
Once the service is running, gRPC endpoints are accessible on port 8080. You can configure your gRPC client applications to communicate with these endpoints to perform CRUD operations and utilize search functionality.

Docker
The service supports Docker for containerization. This allows easy deployment and scalability of the application using Docker images.
