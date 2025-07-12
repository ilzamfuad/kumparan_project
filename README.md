# kumparan_project
This is project for test from kumparan

## Requirements

- **Golang**: Make sure you have Go installed. [Download Go](https://golang.org/dl/)
- **Mockgen**: Used for generating mocks. Install it using:
  ```bash
  go install github.com/golang/mock/mockgen@v1.6.0
  ```
- **Docker**: Required for containerization. [Install Docker](https://www.docker.com/get-started)
- **Docker Compose**: Used for managing multi-container Docker applications. [Install Docker Compose](https://docs.docker.com/compose/install/)
- **Gin**: A web framework for Go. Install it using:
  ```bash
  go get -u github.com/gin-gonic/gin
  ```
- **Postgres**: The database used in this project.

## Installation

1. **Install Golang**:
   - Follow the instructions on the [official Go website](https://golang.org/doc/install).

2. **Setup Configuration**:
   - Copy the sample environment file:
     ```bash
     cp sample.env .env
     ```
   - Update the `.env` file with your configuration.

3. **Run the Application**:
   - Install requirement with Docker Compose:
     ```bash
     docker-compose up --build
     ```

   - Run your application:
     ```bash
     go run main.go
     ```
     you can access the API with localhost:8080

## Database Migration

- Use the `Makefile` to handle database migrations. Run the following command:
  ```bash
  make migrate
  ```

## Generating Mocks

- To generate mocks for a service, use the following command:
  ```bash
  mockgen -source=service/article_service.go -destination=tests/mock/service/article_service.go
  ```
- Replace `service/article_service.go` with the path to the file you want to mock.

## Unit Test

- To run the unit test and check the coverage, use the following command:
  ```bash
  make unit-test
  ```