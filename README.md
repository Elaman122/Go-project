# Currency Exchange Project

## Overview
This project is a RESTful API for managing currency information.

## Endpoints
- `POST /api/v1/menus`:  Create a new currency.
- `GET /api/v1/menus/:id`: Get currency information by identifier.
- `PUT /api/v1/menus/:id`: Update currency information by identifier.
- `DELETE /api/v1/menus/:id`: Delete currency by identifier.

## Database Structure
The database structure is represented as follows:

### Таблица `currency`
- `id`: unique currency identifier (integer).
- `code`: currency code (string).
- `rate`: exchange rate (decimal number).
- `timestamp`: timestamp (timestamp).

## Usage
To use this API, follow these steps:

1. Install all necessary dependencies using go mod tidy.
2. Run the application using go run .
3. Use the corresponding HTTP requests to interact with the API.

## Technologies Used
- Go: programming language.
- Gin: web framework for building web applications.
- PostgreSQL: relational database management system.

## Author
1. name : `Daniar Elaman`
2. ID : `22B031171` 


