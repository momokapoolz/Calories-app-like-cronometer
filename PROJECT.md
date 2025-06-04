# Calories App - Nutrition Tracking System

## Project Overview

The Calories App is a comprehensive nutrition and calorie tracking system designed to help users monitor their food intake, track nutritional information, and achieve their health and fitness goals. Built with modern technologies including Go, Gin, GORM, PostgreSQL, and Redis, this application provides a robust backend API for managing user data, food information, and nutritional tracking.

## Technologies Used

### Backend
- **Go (Golang)**: Main programming language
- **Gin**: HTTP web framework
- **GORM**: Object-Relational Mapping library for database operations
- **PostgreSQL**: Primary database for storing application data
- **Redis**: In-memory data store for session management and caching
- **JWT**: JSON Web Tokens for secure authentication

### Development Tools
- **Environment Variables**: Configuration management via .env files
- **Testing**: Built-in Go testing framework

## System Architecture

The application follows a modular, layered architecture that promotes separation of concerns and maintainability:

### Modular Structure
Each functional area is organized as a separate module:

1. **User Module**: Handles user registration, authentication, and profile management
2. **Food Module**: Manages food item data including names, descriptions, and base nutritional information
3. **Nutrient Module**: Maintains a database of nutrients and their properties
4. **Food Nutrients Module**: Maps the relationship between food items and their nutrient composition
5. **Meal Log Module**: Tracks user meal entries with timestamps and categories
6. **Meal Log Items Module**: Records individual food items within meal logs
7. **User Biometrics Module**: Stores and tracks user health metrics over time

### Layered Architecture
Each module follows a consistent layered approach:

- **Models Layer**: Data structures with GORM tags for database mapping
- **Repository Layer**: Database operations and queries
- **Service Layer**: Business logic and data processing
- **Controller Layer**: API endpoint handlers
- **Routes Layer**: API route definitions and middleware configuration

## Database Schema

The application uses PostgreSQL with the following core tables:

1. **users**: User account information and authentication details
2. **food**: Food items with basic information
3. **nutrient**: Nutritional components like proteins, carbohydrates, vitamins, etc.
4. **food_nutrient**: Junction table linking foods to their nutrient composition
5. **meal_log**: User meal entries with metadata
6. **meal_log_item**: Individual food items within meal logs
7. **user_biometric**: User health metrics like weight, BMI, body fat percentage, etc.

## Authentication System

The application implements a robust authentication system with:

1. **JWT-based Authentication**: Secure token-based authentication
2. **Cookie-based Sessions**: Alternative session management using secure cookies
3. **Token Refresh Mechanism**: Automatic renewal of authentication tokens
4. **Redis-backed Session Store**: Fast, in-memory session data storage

## API Endpoints

### Authentication Endpoints
- User registration with profile creation
- Login with JWT or cookie-based authentication
- Token refresh mechanism
- Secure logout

### User Management
- User profile CRUD operations
- Password management

### Food and Nutrition
- Food item management
- Nutrient information retrieval
- Food-nutrient relationship management

### Meal Logging
- Create, read, update, and delete meal logs
- Add/remove food items from meal logs
- Query meal logs by date ranges

### User Biometrics
- Track health metrics over time
- Support for multiple biometric types
- Historical data analysis

## Features

1. **Comprehensive Food Database**: Extensive database of food items with detailed nutritional information
2. **Personalized Meal Tracking**: Users can log their meals with specific food items and portions
3. **Nutritional Analysis**: Automatic calculation of nutritional intake based on logged meals
4. **Progress Monitoring**: Track health metrics and nutritional patterns over time
5. **User Goal Setting**: Set and monitor progress towards health and fitness goals
6. **Secure Authentication**: Robust user authentication and authorization system

## Deployment

The application is designed to be deployed in various environments:

- **Development**: Local setup with environment variables
- **Production**: Scalable deployment with proper security configurations
- **Database Options**: Support for both local PostgreSQL and Supabase cloud database

## Getting Started

### Prerequisites
- Go 1.18 or later
- PostgreSQL database
- Redis server

### Environment Setup
Create a `.env` file with the following variables:

```
# Database Settings
POSTGRES_DB_CONNECTION_STRING=your_supabase_connection_string
# Or for local development:
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

### Installation and Running

1. Clone the repository
2. Install dependencies: `go mod download`
3. Set up the database and Redis
4. Configure environment variables
5. Run the application: `go run main.go`

## Future Enhancements

1. **Mobile Application**: Develop companion mobile apps for iOS and Android
2. **Machine Learning Integration**: Implement ML for personalized nutritional recommendations
3. **Social Features**: Add social sharing and community features
4. **Barcode Scanning**: Implement barcode scanning for easy food logging
5. **Integration with Fitness Trackers**: Connect with popular fitness tracking devices
6. **Meal Planning**: Add meal planning and recipe suggestion features
7. **Internationalization**: Support for multiple languages and regional food databases

## Project Status

This application is being developed as a Final Year Project (FYP) and is currently in active development.

## License

[Specify your license information here]

## Contributors

[List of contributors and acknowledgments] 