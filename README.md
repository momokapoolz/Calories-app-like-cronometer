# Calories App

A calorie and nutrition tracking application built with Go, Gin, GORM, PostgreSQL, and Redis.

## Project Structure

This project follows a modular approach with a layered architecture for each module.

### Modules

The application is divided into the following modules:

1. **User** - User authentication and management
2. **Food** - Food items management
3. **Nutrient** - Nutrient information management
4. **Food Nutrients** - Relationship between foods and nutrients
5. **Meal Log** - User meal logging
6. **Meal Log Items** - Individual food items in meal logs
7. **User Biometrics** - User health metrics tracking

### Architecture

Each module follows a layered architecture:

- **Models** - Data models with GORM tags
- **Repository** - Database operations
- **Services** - Business logic
- **Controllers** - API endpoints
- **Routes** - Route definitions

## API Endpoints

### Food Module
- `POST /api/v1/foods` - Create a new food
- `GET /api/v1/foods` - Get all foods
- `GET /api/v1/foods/:id` - Get a specific food
- `PUT /api/v1/foods/:id` - Update a food
- `DELETE /api/v1/foods/:id` - Delete a food

### Nutrient Module
- `POST /api/v1/nutrients` - Create a new nutrient
- `GET /api/v1/nutrients` - Get all nutrients
- `GET /api/v1/nutrients/:id` - Get a specific nutrient
- `GET /api/v1/nutrients/category/:category` - Get nutrients by category
- `PUT /api/v1/nutrients/:id` - Update a nutrient
- `DELETE /api/v1/nutrients/:id` - Delete a nutrient

### Food Nutrients Module
- `POST /api/v1/food-nutrients` - Create a new food nutrient
- `GET /api/v1/food-nutrients` - Get all food nutrients
- `GET /api/v1/food-nutrients/:id` - Get a specific food nutrient
- `GET /api/v1/food-nutrients/food/:foodId` - Get food nutrients by food ID
- `GET /api/v1/food-nutrients/nutrient/:nutrientId` - Get food nutrients by nutrient ID
- `PUT /api/v1/food-nutrients/:id` - Update a food nutrient
- `DELETE /api/v1/food-nutrients/:id` - Delete a food nutrient

### Meal Log Module
- `POST /api/v1/meal-logs` - Create a new meal log
- `GET /api/v1/meal-logs/:id` - Get a specific meal log
- `GET /api/v1/meal-logs/user/:userId` - Get meal logs by user ID
- `GET /api/v1/meal-logs/user/:userId/date/:date` - Get meal logs by user ID and date
- `GET /api/v1/meal-logs/user/:userId/date-range` - Get meal logs by user ID and date range
- `PUT /api/v1/meal-logs/:id` - Update a meal log
- `DELETE /api/v1/meal-logs/:id` - Delete a meal log

### Meal Log Items Module
- `POST /api/v1/meal-log-items` - Create a new meal log item
- `GET /api/v1/meal-log-items/:id` - Get a specific meal log item
- `GET /api/v1/meal-log-items/meal-log/:mealLogId` - Get meal log items by meal log ID
- `GET /api/v1/meal-log-items/food/:foodId` - Get meal log items by food ID
- `PUT /api/v1/meal-log-items/:id` - Update a meal log item
- `DELETE /api/v1/meal-log-items/:id` - Delete a meal log item
- `DELETE /api/v1/meal-log-items/meal-log/:mealLogId` - Delete all items for a meal log

### User Biometrics Module
- `POST /api/v1/user-biometrics` - Create a new user biometric
- `GET /api/v1/user-biometrics/:id` - Get a specific user biometric
- `GET /api/v1/user-biometrics/user/:userId` - Get user biometrics by user ID
- `GET /api/v1/user-biometrics/user/:userId/type/:type` - Get user biometrics by user ID and type
- `GET /api/v1/user-biometrics/user/:userId/type/:type/date-range` - Get user biometrics by user ID, type and date range
- `GET /api/v1/user-biometrics/user/:userId/type/:type/latest` - Get latest user biometric by user ID and type
- `PUT /api/v1/user-biometrics/:id` - Update a user biometric
- `DELETE /api/v1/user-biometrics/:id` - Delete a user biometric

## Getting Started

### Prerequisites
- Go 1.18 or later
- PostgreSQL
- Redis

### Environment Variables
Create a `.env` file in the root directory with the following variables:

```
# Database Settings
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=caloriesapp
DB_PORT=5432

# Redis Settings
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Settings
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRATION=24h

# Server Settings
PORT=8080
ENV=development
```

### Running the Application

1. Install dependencies:
   ```
   go mod download
   ```

2. Run the application:
   ```
   go run main.go
   ```

The server will start on `http://localhost:8080`. 