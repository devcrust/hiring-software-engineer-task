# Software Engineering Task: Ad Bidding Service

## Overview

You are tasked with extending an ad bidding service responsible for managing and serving advertisements.

The service is built using Go with the Fiber framework and provides basic functionality for creating ad line items.

## Task Requirements

Your challenge is to:

1. **Implement ad selection logic** to find winning ads based on various criteria
2. **Implement a tracking endpoint** to record user interactions with ads
3. **Add relevancy scoring** to improve ad matching quality
4. **Implement appropriate validation** for all endpoints
5. **Document your approach** and any assumptions made

## Prerequisites

* [Go](https://golang.org/doc/install) 1.24+
* [Docker](https://docs.docker.com/engine/install/)
* [Compose](https://docs.docker.com/compose/install/)
* [golangci-lint](https://golangci-lint.run/)

## Setup & Environment

This repository provides a basic service structure to get started:

```bash
# Build and start the service
docker-compose up -d

# Check service status
curl http://localhost:8080/health

# Test creating a line item
curl -X POST http://localhost:8080/api/v1/lineitems \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Summer Sale Banner",
    "advertiser_id": "adv123",
    "bid": 2.5,
    "budget": 1000.0,
    "placement": "homepage_top",
    "categories": ["electronics", "sale"],
    "keywords": ["summer", "discount"]
  }'

# Get winning ads for a placement (you'll need to implement this)
curl -X GET "http://localhost:8080/api/v1/ads?placement=homepage_top&category=electronics&keyword=discount"
```

## Configuration

The service uses environment variables for configuration, using [Kelsey Hightower's envconfig](https://github.com/kelseyhightower/envconfig) library.

Available environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| APP_NAME | Application name | "Ad Bidding Service" |
| APP_ENVIRONMENT | Running environment | "development" |
| APP_LOG_LEVEL | Log level (debug, info, warn, error) | "info" |
| APP_VERSION | Application version | "1.0.0" |
| SERVER_PORT | HTTP server port | 8080 |
| SERVER_TIMEOUT | Server timeout for requests | "30s" |

## API Structure

The service exposes the following endpoints:

- **POST /api/v1/lineitems**: Create new ad line items with bidding parameters
- **GET /api/v1/ads**: Get winning ads for a specific placement with optional filters (you'll need to implement this)
- **POST /api/v1/tracking**: Record ad interactions (you'll need to implement this)

The complete API specification is available in the OpenAPI document at `api/openapi.yaml`.

## Data Model

The core data model includes:

- **LineItem**: An advertisement with associated bid information
  - `id`: Unique identifier
  - `name`: Display name of the line item
  - `advertiser_id`: ID of the advertiser
  - `bid`: Maximum bid amount (CPM)
  - `budget`: Daily budget for the line item
  - `placement`: Target placement identifier
  - `categories`: List of associated categories
  - `keywords`: List of associated keywords

## Deliverables

Please provide the following:

1. **Ad Selection Logic**: Implement the logic to select winning ads based on placement, categories, and keywords
2. **Tracking Endpoint**: Implement an endpoint to record impressions, clicks, and conversions
3. **Relevancy System**: Develop a scoring mechanism to determine ad relevance
4. **Input Validation**: Add appropriate validation for all API endpoints
5. **Documentation**: Update the README and API docs with your changes

## Evaluation Criteria

Your solution will be evaluated based on:

- **Code quality**: Clean, well-structured, and maintainable code
- **API design**: RESTful design, appropriate error handling, and documentation
- **Implementation quality**: Performance, reliability, and adherence to Go best practices
- **Documentation**: Clear explanation of your approach, design decisions, and trade-offs
- **Testing**: Comprehensive test coverage and consideration of edge cases
- **Innovation**: Creative solutions to the technical challenges presented

## Technical Requirements

- Your solution should be containerized and runnable with docker-compose
- All code should follow Go best practices and conventions
- The API should handle appropriate error cases with meaningful status codes and messages
- Your implementation should consider performance and scaling aspects
- Update the OpenAPI specification to match your implementation

## Storage Solutions

The current implementation uses in-memory storage for simplicity, but this is not suitable for production. You are free to use any storage solution you prefer.
Choose solutions that best fit the requirements and consider factors like scalability, reliability, and performance.

## Scaling Considerations

As part of your solution, please include a section in your documentation addressing the following questions:

1. How would you scale this service to handle millions of ad requests per minute?
2. What bottlenecks do you anticipate and how would you address them?
3. How would you design the system to ensure high availability and fault tolerance?
4. What data storage and access patterns would you recommend for different components (line items, tracking events, etc.)?
5. How would you implement caching to improve performance?

## Getting Started

1. Clone this repository
2. Explore the existing code to understand the current implementation
3. Run the service locally using docker-compose
4. Implement the required features
5. Update tests and documentation
6. Submit your solution

For local development:
- Build and run: `go run ./cmd/server`
- Run tests: `go test ./...`
- Build binary: `go build -o adserver ./cmd/server`
- Linting: `golangci-lint run`

## Project Structure

```
.
├── api/                    # API documentation and OpenAPI spec
├── cmd/                    # Application entrypoints
│   └── server/             # Main server application
├── internal/               # Private application code
│   ├── config/             # Configuration handling
│   ├── handler/            # HTTP handlers
│   ├── model/              # Data models
│   └── service/            # Business logic
├── docker-compose.yml      # Docker Compose configuration
├── Dockerfile              # Docker build configuration
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── README.md               # Project documentation
```

Good luck!