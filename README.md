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

## Database Schema

### Tables and Relationships

#### Users
- `id` (UUID, PK): Unique identifier
- `username` (VARCHAR): Unique username for login
- `email` (VARCHAR): Unique email address
- `password_hash` (VARCHAR): Bcrypt hashed password
- `full_name` (VARCHAR): User's full name
- `role` (VARCHAR): User role (default: 'user')
- `preferences` (JSONB): User preferences stored as JSON
- `is_active` (BOOLEAN): Account status
- `last_login_at` (TIMESTAMP): Last login timestamp
- `created_at`, `updated_at` (TIMESTAMP): Record timestamps
- `is_deleted` (BOOLEAN): Soft delete flag

#### Products
- `id` (UUID, PK): Unique identifier
- `sku` (VARCHAR): Unique product code
- `name` (VARCHAR): Product name
- `description` (TEXT): Product description
- `price` (DECIMAL): Regular price
- `sale_price` (DECIMAL): Discounted price (if applicable)
- `cost_price` (DECIMAL): Product cost
- `stock_quantity` (INT): Available stock
- `status` (VARCHAR): Product status (active, out_of_stock, etc.)
- `attributes` (JSONB): Product attributes stored as JSON
- `metadata` (JSONB): Additional metadata
- `created_at`, `updated_at` (TIMESTAMP): Record timestamps
- `is_deleted` (BOOLEAN): Soft delete flag

#### Categories
- `id` (UUID, PK): Unique identifier
- `name` (VARCHAR): Category name
- `slug` (VARCHAR): URL-friendly unique identifier
- `description` (TEXT): Category description
- `parent_id` (UUID, FK): Self-referencing foreign key for hierarchical categories
- `is_active` (BOOLEAN): Category status
- `display_order` (INT): Ordering for display
- `created_at`, `updated_at` (TIMESTAMP): Record timestamps
- `is_deleted` (BOOLEAN): Soft delete flag

#### Product Categories (Junction)
- `product_id` (UUID, FK): Reference to products
- `category_id` (UUID, FK): Reference to categories
- Combined primary key (product_id, category_id)

#### Reviews
- `id` (UUID, PK): Unique identifier
- `product_id` (UUID, FK): Reference to product
- `user_id` (UUID, FK): Reference to user
- `rating` (SMALLINT): Rating (1-5)
- `title` (VARCHAR): Review title
- `comment` (TEXT): Review content
- `is_verified_purchase` (BOOLEAN): Verified purchase flag
- `helpful_votes` (INT): Number of helpful votes
- `created_at`, `updated_at` (TIMESTAMP): Record timestamps
- `is_deleted` (BOOLEAN): Soft delete flag

#### Wishlist
- `user_id` (UUID, FK): Reference to user
- `product_id` (UUID, FK): Reference to product
- `added_at` (TIMESTAMP): When item was added
- Combined primary key (user_id, product_id)

#### Inventory Movements
- `id` (UUID, PK): Unique identifier
- `product_id` (UUID, FK): Reference to product
- `quantity` (INT): Quantity changed
- `movement_type` (VARCHAR): Type of movement (purchase, sale, adjustment, return)
- `reference_id` (UUID): Reference to related entity
- `notes` (TEXT): Additional notes
- `created_by` (UUID, FK): User who created the record
- `created_at` (TIMESTAMP): Record timestamp

### Entity Relationship Diagram (ERD)

```
Users 1──────────┐
   │             │
   │             ▼
   │         Reviews
   │             ▲
   │             │
Products ◄──────┘
   ▲
   │
   ├────► Inventory Movements
   │
   │     Categories
   │     ▲   │
   │     │   │
   └─────┴───┘
      via
 Product Categories
```

### Key Relationships

- Many-to-Many: Products <-> Categories (via product_categories junction table)
- One-to-Many: Users -> Reviews
- One-to-Many: Products -> Reviews
- Many-to-Many: Users <-> Products (via wishlist)
- One-to-Many: Products -> Inventory Movements
- Hierarchical: Categories -> Categories (self-referencing via parent_id)

## Security Considerations

- All passwords are hashed using bcrypt
- Authentication is handled via JWT tokens
- Role-based access control for protected endpoints
- Input validation for all API endpoints
- Prepared statements to prevent SQL injection
- CORS protection for API endpoints

## Future Improvements

The following technologies are planned for implementation to enhance system performance, search capabilities, and observability:

### Redis for Caching

Redis will be integrated to implement a robust caching strategy:

- **API Response Caching**: Cache frequently requested data to reduce database load
- **Session Management**: Store user sessions for faster authentication
- **Rate Limiting**: Implement rate limiting for API endpoints
- **Background Job Queue**: Manage asynchronous tasks efficiently
- **Leaderboards/Counters**: Fast access to dynamic statistics

Implementation plans include:
- Cache invalidation strategies to ensure data consistency
- Tiered caching approach for different types of data
- Horizontal scaling with Redis Cluster for high availability

### ElasticSearch for Full-Text Search

ElasticSearch will be implemented to enhance search capabilities:

- **Advanced Product Search**: Fuzzy matching, synonyms, and contextual search
- **Faceted Navigation**: Dynamic filtering based on product attributes
- **Search Analytics**: Track and optimize search patterns
- **Multi-language Support**: Search across multiple languages
- **Real-time Indexing**: Immediate indexing of new products

Implementation plans include:
- Custom analyzers for e-commerce specific terminology
- Relevance tuning to improve search result quality
- Integration with PostgreSQL via change data capture (CDC)

### Prometheus & Grafana for Metrics

Prometheus and Grafana will be implemented to provide comprehensive monitoring and metrics:

- **System Metrics**: CPU, memory, disk usage, and network performance
- **Application Metrics**: Request rates, response times, and error rates
- **Business Metrics**: Order volume, revenue, and conversion rates
- **Custom Dashboards**: Tailored visualizations for different stakeholders
- **Alerting**: Proactive notifications for system anomalies

Implementation plans include:
- Custom instrumentation for critical business processes
- SLO/SLI tracking for service reliability
- Long-term metrics storage for trend analysis

These technologies will progressively enhance the platform's performance, user experience, and operational visibility as the system scales.

