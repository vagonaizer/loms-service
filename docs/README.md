# LOMS Service (Logistics and Order Management System)

## Overview
LOMS is a microservice responsible for managing orders and stock inventory in an e-commerce system. It provides a robust API for order management, stock tracking, and inventory control through gRPC communication.

## Architecture

### Clean Architecture
The service follows Clean Architecture principles, separating the codebase into distinct layers:

1. **Domain Layer** (`internal/domain/`)
   - Core business entities and interfaces
   - Defines the business rules and models
   - Independent of external frameworks and tools

2. **Use Case Layer** (`internal/usecase/`)
   - Implements business logic
   - Orchestrates the flow of data between entities
   - Depends only on the domain layer

3. **Infrastructure Layer** (`internal/infrastructure/`)
   - Implements interfaces defined in the domain layer
   - Contains external dependencies and implementations
   - Includes repositories, API handlers, and external service clients

### Key Components

#### Repositories
- **OrderRepository**: Manages order data in memory
- **StockRepository**: Handles stock inventory management
  - Tracks total and reserved quantities
  - Provides atomic operations for stock reservation and release

#### API Layer
- gRPC server implementation
- Protocol buffer definitions for service contracts
- Handles client requests and responses

#### Service Integration
- Client implementation for cart-service communication
- gRPC-based inter-service communication

## Integration with Cart Service

### Communication Flow
1. **Stock Validation**
   - Cart service checks product availability via `StocksInfo`
   - LOMS returns available quantity (total - reserved)

2. **Order Creation**
   - Cart service initiates checkout via `OrderCreate`
   - LOMS reserves items and creates order
   - Returns order ID to cart service

3. **Order Management**
   - Cart service can query order status via `OrderInfo`
   - Supports order payment and cancellation

### API Endpoints

#### OrderCreate
- Creates new order from cart items
- Reserves required stock
- Returns order ID
- Status transitions: new → awaiting payment/failed

#### OrderInfo
- Retrieves order details
- Returns status, user info, and items

#### OrderPay
- Marks order as paid
- Reduces stock quantities
- Status: awaiting payment → paid

#### OrderCancel
- Cancels pending order
- Releases reserved stock
- Status: awaiting payment → cancelled

#### StocksInfo
- Returns available stock quantity
- Considers reserved items
- Used for stock validation


### Storage
- In-memory storage for orders and stock
- Initial stock data loaded from JSON file
- Thread-safe operations with mutex locks

## Development

### Prerequisites
- Go 1.21 or higher
- Protocol Buffers compiler
- gRPC tools

### Building
```bash
make build
```

### Running
```bash
make run
```

### Testing
```bash
make test
```

## Configuration
Service configuration is managed through `config/config.yaml`:
- gRPC server port
- Service addresses
- Other runtime parameters

## Best Practices
1. **Atomic Operations**
   - Stock reservation and release are atomic
   - Prevents race conditions

2. **Error Handling**
   - Clear error types
   - Proper error propagation

3. **Logging**
   - Comprehensive operation logging
   - Debug information
