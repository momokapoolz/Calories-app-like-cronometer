# Calories App with GORM and PostgreSQL

This project uses GORM with PostgreSQL to manage database operations for a calorie tracking application.

## Project Structure

- `models/user.go`: Contains the User model definition
- `database/db.go`: Handles database connection and migrations
- `repository/user_repository.go`: Provides methods to perform CRUD operations on the User model
- `main.go`: Application entry point
- `main_test.go`: Tests for the database operations

## Setup PostgreSQL

1. Install PostgreSQL if you haven't already.

2. Create a new database:
   ```
   createdb calories_app
   ```

3. Configure the database connection in `database/db.go` if needed:
   ```go
   dsn := "host=localhost user=postgres password=postgres dbname=calories_app port=5432 sslmode=disable"
   ```

## Getting Started

1. Install dependencies:
   ```
   go mod tidy
   ```

2. Run the application:
   ```
   go run main.go
   ```

3. Run tests:
   ```
   go test -v
   ```

## User Model

The User model has the following fields:

- `ID`: Primary key, auto-incremented, unique integer
- `Name`: User's name (varchar)
- `Email`: User's email (varchar)
- `Password`: User's password (text)
- `Age`: User's age (integer)
- `Gender`: User's gender (varchar)
- `Weight`: User's weight (decimal)
- `Height`: User's height (decimal)
- `Goal`: User's fitness goal (varchar)
- `ActivityLevel`: User's activity level (varchar)
- `CreatedAt`: Timestamp when the user was created

## Using the Repository

The UserRepository provides the following methods:

- `Create`: Create a new user
- `FindByID`: Find a user by ID
- `FindByEmail`: Find a user by email
- `Update`: Update an existing user
- `Delete`: Delete a user
- `FindAll`: Get all users

Example usage:

```go
// Initialize the repository
userRepo := repository.NewUserRepository()

// Create a new user
user := models.User{
    Name: "John Doe",
    Email: "john@example.com",
    Password: "password123",
    Age: 30,
    Gender: "Male",
    Weight: 70.5,
    Height: 175.0,
    Goal: "Weight Loss",
    ActivityLevel: "Moderate",
}

// Save the user
err := userRepo.Create(&user)
if err != nil {
    log.Fatal("Failed to create user:", err)
}

// Find a user by ID
foundUser, err := userRepo.FindByID(user.ID)
if err != nil {
    log.Fatal("Failed to find user:", err)
}

// Update a user
user.Name = "Jane Doe"
err = userRepo.Update(&user)
if err != nil {
    log.Fatal("Failed to update user:", err)
}

// Delete a user
err = userRepo.Delete(user.ID)
if err != nil {
    log.Fatal("Failed to delete user:", err)
}
``` 