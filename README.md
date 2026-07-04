# Restaurant Management API

This project is a backend REST API for managing a restaurant’s daily operations. It is built with Go, Gin Gonic, MongoDB, JWT authentication, and Stripe integration for payment-related functionality.

The system is designed to handle core restaurant workflows such as user authentication, food and menu management, table handling, order creation, order items, and invoice generation.

## Features

- User registration and authentication
- JWT-based authorization for protected routes
- Food item management
- Menu management
- Table management
- Order and order item management
- Invoice generation
- Environment-based configuration
- Stripe test-mode payment integration
- Payment intent creation and payment lookup

## Tech Stack

- Go
- Gin Gonic
- MongoDB
- JWT
- Godotenv
- Stripe Go SDK

## Project Structure

```text
server/
├── main.go
├── controllers/
│   ├── foodController.go
│   ├── inviceController.go
│   ├── menuController.go
│   ├── orderController.go
│   ├── orderitemController.go
│   ├── paymentController.go
│   ├── tableController.go
│   └── userController.go
├── database/
│   └── databaseConnection.go
├── helpers/
│   └── tokenHelper.go
├── middleware/
│   └── authMiddleware.go
├── models/
│   ├── foodModel.go
│   ├── invoiceModel.go
│   ├── menuModel.go
│   ├── noteModel.go
│   ├── ordeModel.go
│   ├── orderItemModel.go
│   ├── paymentModel.go
│   ├── tableModel.go
│   └── userModel.go
├── routes/
│   ├── foodRouter.go
│   ├── invoiceRouter.go
│   ├── menuRouter.go
│   ├── orderitemRouter.go
│   ├── orderRouter.go
│   ├── paymentRouter.go
│   ├── tableRouter.go
│   └── userRouter.go
└── go.mod
```

## Folder Explanation

- controllers/: Contains request handling logic for each module.
- database/: Handles database connection setup.
- helpers/: Utility functions such as JWT token generation and validation.
- middleware/: Authentication and request protection logic.
- models/: Defines data structures for users, foods, menus, orders, invoices, and tables.
- routes/: Maps API endpoints to their respective controller functions.

## Prerequisites

Before running the project, ensure you have:

- Go installed on your machine
- MongoDB running or a MongoDB Atlas connection string
- A Stripe account if you want to use payment-related features

## Installation

1. Open the project folder.
2. Navigate to the server directory:

```bash
cd server
```

3. Install dependencies:

```bash
go mod tidy
```

## Environment Variables

Create a .env file inside the server directory with values similar to the following:

```env
PORT=8080
MONGODB=your_mongodb_connection_string
DB_Name=restaurant_db
SECRET_KEY=your_secret_key
STRIPE_SECRET_KEY=your_stripe_test_secret_key
STRIPE_PUBLISHABLE_KEY=your_stripe_test_publishable_key
STRIPE_WEBHOOK_SECRET=your_stripe_webhook_secret
```

## Running the Project

Start the server with:

```bash
go run main.go
```

The server will start on the port defined in the .env file, or on port 8000 if no port is set.

## API Modules

- Users: login and authentication-related actions
- Foods: create, update, and manage food items
- Menus: define restaurant menus
- Tables: manage table availability and details
- Orders: create and manage customer orders
- Order Items: handle items inside each order
- Invoices: generate billing and invoice records
- Payments: create Stripe payment intents and retrieve payment records

## Stripe Payment Integration

The project now includes a Stripe-based payment flow for development mode.

### Available Payment Endpoints

- GET /payments/config
  - Returns the Stripe publishable key and mode information
- POST /payments/create-payment-intent
  - Creates a Stripe PaymentIntent for an order
- POST /payments/webhook
  - Handles Stripe webhook events such as payment success or failure
- GET /payments/:payment_id
  - Retrieves a payment record by payment ID
- GET /payments/order/:order_id
  - Retrieves a payment record by order ID

### Example Payment Request

```json
{
  "order_id": "64f1abc123",
  "user_id": "64f2def456",
  "amount": 5000,
  "currency": "inr",
  "payment_method_types": ["card", "upi"]
}
```

### Development Mode Notes

- Use Stripe test keys from your Stripe dashboard.
- The payment flow is intended for development and testing.
- Card and UPI-style payment methods can be passed through the payment intent request.

## Authentication

Protected routes require a JWT token in the request header:

```http
token: <your_jwt_token>
```

## Notes

- Make sure your MongoDB connection string is valid.
- The project uses middleware to protect routes that require authentication.
- Stripe keys should be kept private and never committed to source control.

## Future Improvements

Possible enhancements for this project include:

- Admin dashboard integration
- Better order status tracking
- Payment flow improvements
- Swagger API documentation
- Unit and integration tests
