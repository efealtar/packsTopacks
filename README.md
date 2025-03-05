## Getting Started

### Prerequisites

- Go 1.x
- Docker, Docker compose

### Installation

1. Clone the repos:

   ```bash
   git clone https://github.com/efealtar/packsTopacks.git
   ```

   ```bash
   cd packTopack
   ```

2. Be sure you have docker and docker compose v2, To build:

   ```bash
   docker compose up --build -d
   ```

3. Application is running at:
   ```bash
   http://localhost:3000
   ```

### To Clear Docker Packages

    ```bash
    docker compose down
    ```

### To Test Application with default values

Note: Application will automatically test before every successful build

    ```bash
    go test -v ./packserver/...
    ```

# Pack Size Calculator API

This project is a solution for calculating optimal pack sizes for product orders. The application provides a REST API built with Go that determines the most efficient combination of packs to fulfill customer orders.

## Problem Description

Given a set of standard pack sizes (250, 500, 1000, 2000, 5000 items), the application calculates the optimal combination of packs for any order quantity following these rules:

1. Only whole packs can be sent. Packs cannot be broken open.
2. Send out the least amount of items to fulfill the order (primary priority).
3. Send out as few packs as possible to fulfill each order (secondary priority).

### Example Scenarios

| Items Ordered | Correct Packs               | Why This is Optimal             |
| ------------- | --------------------------- | ------------------------------- |
| 1             | 1 x 250                     | Smallest possible pack size     |
| 250           | 1 x 250                     | Exact match                     |
| 251           | 1 x 500                     | Minimum items above order       |
| 501           | 1 x 500, 1 x 250            | Minimum combination above order |
| 12001         | 2 x 5000, 1 x 2000, 1 x 250 | Optimal combination             |

## Features

- RESTful API endpoints for pack calculation
- Flexible configuration for pack sizes
- Web UI for easy interaction
- Comprehensive unit tests
- Docker support

## Technology Stack

- Backend: Go
- Frontend: React

Important NOTE: To avoid unnecessary build in every compose command in /client folder static build already implemented and local dockerFile is managed to use static build.
