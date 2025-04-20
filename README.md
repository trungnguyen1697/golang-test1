# E-commerce Go API

This repository contains a RESTful API for an e-commerce platform built with Go, PostgreSQL, and Docker.

## Features

- 🔐 **User Authentication & Authorization**: JWT-based auth with role-based access control
- 👤 **User Management**: Registration, profile management, and admin control
- 📦 **Product Management**: CRUD operations with category management
- 🔖 **Category Management**: Hierarchical categories with product assignments
- ⭐ **Review System**: Product reviews with ratings
- ❤️ **Wishlist**: User wishlists for saving favorite products
- 📊 **Dashboard**: Admin dashboard with statistics
- 🔍 **Search & Filtering**: Advanced product filtering and search

## Technology Stack

- **Backend**: Go (Golang) 1.24+
- **Database**: PostgreSQL 17.4
- **API**: RESTful API with JSON responses
- **Authentication**: JWT (JSON Web Tokens)
- **Containerization**: Docker & Docker Compose
- **Documentation**: OpenAPI/Swagger

## Project Structure

```
ecommerce-go/
├── cmd/
│   └── api/                  # Application entry point
│       └── main.go           # Main application file
├── internal/
│   ├── config/               # Configuration management
│   ├── database/             # Database connections and migrations
│   ├── handlers/             # HTTP handlers for API endpoints
│   ├── middleware/           # HTTP middleware components
│   ├── models/               # Data models and DTOs
│   ├── repository/           # Database access layer
│   ├── routes/               # API route definitions
│   ├── services/             # Business logic layer
│   └── utils/                # Utility functions
├── pkg/
│   └── logger/               # Logging package
├── .env.example              # Example environment variables
├── .gitignore                # Git ignore file
├── docker-compose.yml        # Docker Compose configuration
├── Dockerfile                # Docker build instructions
├── go.mod                    # Go module file
└── README.md                 # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- PostgreSQL (if running locally)
- golang-migrate
- swag
- make

### Running with Docker

1. Clone the repository:
   ```bash
   git clone https://github.com/trungnguyen1697/golang-test1.git
   cd golang-test1
   ```

2. Create an `.env` file from the example:
   ```bash
   cp .env.example .env
   ```

3. Build and start the containers:
   ```bash
   make docker.run
   ```

4. The API will be available at `http://localhost:5000` or `http://localhost:5000/swagger/index.html` for the Swagger UI.

### Running Locally

1. Clone the repository:
   ```bash
   git clone https://github.com/trungnguyen1697/golang-test1.git
   cd golang-test1
   ```

2. Create an `.env` file from the example:
   ```bash
   cp .env.example .env
   ```

3. Set up the PostgreSQL database:
      ```bash
    make docker.setup
   make docker.postgres
    ```
   
4. Migrate db & seed some demo data:
    ```bash
    make migrate.up
    ```
   
5. Run the application:
   ```bash
   make run
   ```
   
Common HTTP status codes:
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

## Security Considerations

- All passwords are hashed using bcrypt
- Authentication is handled via JWT tokens
- Role-based access control for protected endpoints
- Input validation for all API endpoints
- Prepared statements to prevent SQL injection
- CORS protection for API endpoints

## License

This project is licensed under the MIT License - see the LICENSE file for details.
