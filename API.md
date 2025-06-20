# Calories App API Documentation

This document provides a detailed description of the Calories App API, including authentication, protected endpoints, and data models.

**Base URL**: The base URL for all endpoints is `/api/v1`.

## Authentication

Authentication is handled via JWT (JSON Web Tokens). Currently, the API primarily uses **Token-Based authentication** through the User Module endpoints.

**Current Authentication Method:**
- **Token-Based**: The client receives token IDs (UUIDs) and sends them back in the `Authorization` header for subsequent requests. This is suitable for both mobile and web clients.

**Note**: Cookie-based authentication through the Auth Module (`/api/v1/auth/*`) is temporarily disabled. Use the User Module endpoints (`/api/v1/login`, `/api/v1/logout`) instead.

### Register a New User
- **URL**: `/api/v1/register`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepassword",
    "age": 30,
    "gender": "Male",
    "weight": 70.5,
    "height": 175.0,
    "goal": "Weight Loss",
    "activity_level": "Moderate"
  }
  ```
- **Success Response (Code `201 Created`)**:
    ```json
    {
      "status": "success",
      "message": "User registered successfully",
      "data": {
        "user": {
          "id": 1,
          "name": "John Doe",
          "email": "john@example.com",
          "role": "user"
        }
      }
    }
    ```
- **Error Response**: 
  - **Code**: 400 Bad Request
    ```json
    {
      "status": "error",
      "message": "Email already in use"
    }
    ```
  - **Code**: 400 Bad Request
    ```json
    {
      "status": "error",
      "message": "Invalid request format",
      "error": "Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
    }
    ```

### User Login
- **URL**: `/api/v1/login`
- **Method**: `POST`
- **Description**: Authenticates a user and returns JWT token IDs for token-based authentication.
- **Request Body**:
  ```json
  {
    "email": "john@example.com",
    "password": "securepassword"
  }
  ```
- **Success Response (Code `200 OK`)**:
  - **Content**: 
    ```json
    {
      "status": "success",
      "message": "Login successful",
      "data": {
        "user": {
          "id": 1,
          "name": "John Doe",
          "email": "john@example.com",
          "role": "user"
        },
        "tokens": {
          "access_token_id": "0097fcb4-abeb-4140-8683-181f7d796755",
          "refresh_token_id": "15aff74e-bafe-410b-a7f7-fa5a96e530a6",
          "expires_in": 3600
        }
      }
    }
    ```

---

- **Common Error Responses**: 
  - **Code `401 Unauthorized`**: `{"status": "error", "message": "Invalid credentials"}`
  - **Code `500 Internal Server Error`**: `{"status": "error", "message": "Authentication failed", "error": "..."}`


### User Logout
- **URL**: `/api/v1/logout`
- **Method**: `POST`
- **Description**: Invalidates the current user's access token and logs them out
- **Authentication**: Required (Bearer token with access_token_id)
- **Headers**: 
  - `Authorization: Bearer {access_token_id}`
- **Success Response (Code `200 OK`)**:
  - **Content**: 
    ```json
    {
      "status": "success",
      "message": "Logged out successfully"
    }
    ```
- **Error Responses**: 
  - **Code `400 Bad Request`**: `{"status": "error", "message": "No active session"}`
  - **Code `401 Unauthorized`**: `{"error": "Authentication required"}`
  - **Code `500 Internal Server Error`**: `{"status": "error", "message": "Failed to invalidate token"}`

---

## Auth Module Endpoints (Temporarily Disabled)

> **Note**: The following Auth Module endpoints (`/api/v1/auth/*`) are temporarily disabled. Use the User Module endpoints above instead.

### Token Refresh (Temporarily Disabled)
- **URL**: `/api/v1/auth/refresh`
- **Method**: `POST`
- **Description**: Generates a new access token using the refresh token. Requires cookie-based authentication.
- **Status**: ⚠️ **Temporarily Disabled**
- **Required Cookies**:
  - `refresh-id`: Contains the refresh token ID.
- **Success Response (Code `200 OK`)**:
  - **Cookies Updated**: 
    - `jwt-id`: A new access token ID is set.
  - **Content**: `{"status": "success", "message": "Token refreshed successfully"}`
- **Error Response (`401 Unauthorized`)**: `{"status": "error", "message": "Invalid refresh token"}`


### Auth Logout (Temporarily Disabled)
- **URL**: `/api/v1/auth/logout`
- **Method**: `POST`
- **Description**: Invalidates both access and refresh tokens by deleting them from the server's store (Redis).
- **Status**: ⚠️ **Temporarily Disabled** - Use `/api/v1/logout` instead
- **Authentication**: This endpoint works with both cookie-based and token-based authentication. The server will automatically detect the method used.
- **Success Response (Code `200 OK`)**:
  - **Cookies Cleared**: If using cookie-based auth, `jwt-id` and `refresh-id` cookies are cleared.
  - **Content**: `{"status": "success", "message": "Logged out successfully"}`
- **Error Response (`400 Bad Request`)**: `{"status": "error", "message": "No active session"}`


### Token Validation for Protected Endpoints
All protected endpoints validate authentication credentials on every request.

**Current Method: Token-Based (Header)**
- The client must include an `Authorization` header.
- The token can be either the **full JWT** or the **access\_token\_id (UUID)**. The backend middleware is designed to handle both formats seamlessly.
- **Format**: `Authorization: Bearer {jwt_or_uuid}`

**Temporarily Disabled: Cookie-Based**
- ⚠️ Cookie-based authentication is temporarily disabled.
- Previously, browsers could send the `jwt-id` cookie with each request, but this functionality is currently unavailable.

- **Common Error Response (`401 Unauthorized`)**:
  ```json

  {
    "error": "Authentication required"
  }
  ```

## Protected Endpoints

### Get User Profile
- **URL**: `/api/v1/profile`
- **Method**: `GET`
- **Authentication**: Required (`Authorization: Bearer` header with access_token_id)
- **Success Response (Code `200 OK`)**:
  - **Content**: 
    ```json
    {
      "status": "success",
      "data": {
        "user": {
          "id": 1,
          "name": "John Doe",
          "email": "john@example.com",
          "age": 30,
          "gender": "Male",
          "weight": 70.5,
          "height": 175.0,
          "goal": "Weight Loss",
          "activity_level": "Moderate",
          "role": "user",
          "created_at": "2023-09-22T10:30:45Z"
        }
      }
    }
    ```
- **Error Response (`401 Unauthorized`)**: `{"error": "Not authenticated"}`

### Get Auth User Profile (Example)
- **URL**: `/api/auth/profile`
- **Method**: `GET`
- **Headers**: 
  - `Authorization: Bearer {access_token}`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "user_id": 1,
      "email": "john@example.com",
      "role": "user"
    }
    ```
- **Error Response**: 
  - **Code**: 401 Unauthorized
    ```json
    {
      "error": "Not authenticated"
    }
    ```

## User Endpoints

### Create a new user
- **URL**: `/api/users`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com",
    "password_hash": "securepassword",
    "age": 30,
    "gender": "Male",
    "weight": 70.5,
    "height": 175.0,
    "goal": "Weight Loss",
    "activity_level": "Moderate",
    "role": "user"
  }
  ```
- **Success Response**: 
  - **Code**: 201 Created
  - **Content**: User object with ID

### Get all users
- **URL**: `/api/v1/users`
- **Method**: `GET`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of user objects

### Get user by ID
- **URL**: `/api/v1/users/{id}`
- **Method**: `GET`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: User object
- **Error Response**: 
  - **Code**: 404 Not Found

### Get user by email
- **URL**: `/api/v1/users/email/{email}`
- **Method**: `GET`
- **URL Parameters**: `email=[string]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: User object
- **Error Response**: 
  - **Code**: 404 Not Found

### Get users by role
- **URL**: `/api/v1/users/role/{role}`
- **Method**: `GET`
- **URL Parameters**: `role=[string]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of user objects

### Update a user
- **URL**: `/api/v1/users/{id}`
- **Method**: `PUT`
- **URL Parameters**: `id=[uint]`
- **Request Body**: Same as create user
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Updated user object
- **Error Response**: 
  - **Code**: 404 Not Found

### Delete a user
- **URL**: `/api/v1/users/{id}`
- **Method**: `DELETE`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 204 No Content
- **Error Response**: 
  - **Code**: 404 Not Found or 500 Internal Server Error

## Food Endpoints

### Create a new food
- **URL**: `/api/v1/foods`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "Apple",
    "serving_size_gram": 100,
    "source": "USDA",
    "image_url": "https://example.com/apple.jpg"
  }
  ```
- **Success Response**: 
  - **Code**: 201 Created
  - **Content**: 
    ```json
    {
      "id": 1,
      "name": "Apple",
      "serving_size_gram": 100,
      "source": "USDA",
      "image_url": "https://example.com/apple.jpg"
    }
    ```
- **Error Response**: 
  - **Code**: 400 Bad Request
    ```json
    {
      "error": "Invalid request format"
    }
    ```
  - **Code**: 400 Bad Request
    ```json
    {
      "error": "Image URL must be a valid URL"
    }
    ```
  - **Code**: 500 Internal Server Error
    ```json
    {
      "error": "Failed to create food"
    }
    ```

### Get all foods
- **URL**: `/api/v1/foods`
- **Method**: `GET`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of food objects
- **Error Response**: 
  - **Code**: 500 Internal Server Error
    ```json
    {
      "error": "Failed to retrieve foods"
    }
    ```

### Get food by ID
- **URL**: `/api/v1/foods/{id}`
- **Method**: `GET`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Food object
- **Error Response**: 
  - **Code**: 404 Not Found
    ```json
    {
      "error": "Food not found"
    }
    ```

### Update a food
- **URL**: `/api/v1/foods/{id}`
- **Method**: `PUT`
- **URL Parameters**: `id=[uint]`
- **Request Body**: Same as create food
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Updated food object
- **Error Response**: 
  - **Code**: 404 Not Found or 500 Internal Server Error

### Delete a food
- **URL**: `/api/v1/foods/{id}`
- **Method**: `DELETE`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "message": "Food deleted successfully"
    }
    ```
- **Error Response**: 
  - **Code**: 404 Not Found or 500 Internal Server Error

## Nutrient Endpoints

### Create a new nutrient
- **URL**: `/api/v1/nutrients`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "Vitamin C",
    "category": "Vitamins"
  }
  ```
- **Success Response**: 
  - **Code**: 201 Created
  - **Content**: 
    ```json
    {
      "id": 1,
      "name": "Vitamin C",
      "category": "Vitamins"
    }
    ```
- **Error Response**: 
  - **Code**: 400 Bad Request or 500 Internal Server Error

### Get all nutrients
- **URL**: `/api/v1/nutrients`
- **Method**: `GET`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of nutrient objects
- **Error Response**: 
  - **Code**: 500 Internal Server Error

### Get nutrient by ID
- **URL**: `/api/v1/nutrients/{id}`
- **Method**: `GET`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Nutrient object
- **Error Response**: 
  - **Code**: 404 Not Found

### Get nutrients by category
- **URL**: `/api/v1/nutrients/category/{category}`
- **Method**: `GET`
- **URL Parameters**: `category=[string]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of nutrient objects
- **Error Response**: 
  - **Code**: 500 Internal Server Error

### Update a nutrient
- **URL**: `/api/v1/nutrients/{id}`
- **Method**: `PUT`
- **URL Parameters**: `id=[uint]`
- **Request Body**: Same as create nutrient
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Updated nutrient object
- **Error Response**: 
  - **Code**: 404 Not Found or 500 Internal Server Error

### Delete a nutrient
- **URL**: `/api/v1/nutrients/{id}`
- **Method**: `DELETE`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "message": "Nutrient deleted successfully"
    }
    ```
- **Error Response**: 
  - **Code**: 404 Not Found or 500 Internal Server Error

## Food Nutrients Endpoints

### Create a new food nutrient
- **URL**: `/api/v1/food-nutrients`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "food_id": 1,
    "nutrient_id": 1,
    "amount_per_100g": 4.5
  }
  ```
- **Success Response**: 
  - **Code**: 201 Created
  - **Content**: Food nutrient object
- **Error Response**: 
  - **Code**: 400 Bad Request or 500 Internal Server Error

### Get all food nutrients
- **URL**: `/api/v1/food-nutrients`
- **Method**: `GET`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of food nutrient objects
- **Error Response**: 
  - **Code**: 500 Internal Server Error

### Get food nutrient by ID
- **URL**: `/api/v1/food-nutrients/{id}`
- **Method**: `GET`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Food nutrient object
- **Error Response**: 
  - **Code**: 404 Not Found

### Get food nutrients by food ID
- **URL**: `/api/v1/food-nutrients/food/{foodId}`
- **Method**: `GET`
- **URL Parameters**: `foodId=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of food nutrient objects
- **Error Response**: 
  - **Code**: 500 Internal Server Error

### Get food nutrients by nutrient ID
- **URL**: `/api/v1/food-nutrients/nutrient/{nutrientId}`
- **Method**: `GET`
- **URL Parameters**: `nutrientId=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of food nutrient objects
- **Error Response**: 
  - **Code**: 500 Internal Server Error

### Update a food nutrient
- **URL**: `/api/v1/food-nutrients/{id}`
- **Method**: `PUT`
- **URL Parameters**: `id=[uint]`
- **Request Body**: Same as create food nutrient
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Updated food nutrient object
- **Error Response**: 
  - **Code**: 404 Not Found or 500 Internal Server Error

### Delete a food nutrient
- **URL**: `/api/v1/food-nutrients/{id}`
- **Method**: `DELETE`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "message": "Food nutrient deleted successfully"
    }
    ```
- **Error Response**: 
  - **Code**: 404 Not Found or 500 Internal Server Error

## Meal Log Endpoints

### Create a new meal log
- **URL**: `/api/v1/meal-logs`
- **Method**: `POST`
- **Authentication**: Requires JWT token
- **Request Body**:
  ```json
  {
    "meal_type": "Breakfast",
    "items": [
      {
        "food_id": 1,
        "quantity": 2,
        "quantity_grams": 100.5
      },
      {
        "food_id": 3,
        "quantity": 1,
        "quantity_grams": 50.0
      }
    ]
  }
  ```
- **Success Response**: 
  - **Code**: 201 Created
  - **Content**: 
    ```json
    {
      "meal_log": {
        "id": 1,
        "user_id": 5,
        "created_at": "2023-09-22T10:30:45Z",
        "meal_type": "Breakfast"
      },
      "items": [
        {
          "id": 1,
          "meal_log_id": 1,
          "food_id": 1,
          "quantity": 2,
          "quantity_grams": 100.5
        },
        {
          "id": 2,
          "meal_log_id": 1,
          "food_id": 3,
          "quantity": 1,
          "quantity_grams": 50.0
        }
      ]
    }
    ```
- **Error Response**: 
  - **Code**: 400 Bad Request or 500 Internal Server Error
  - **Code**: 401 Unauthorized
    ```json
    {
      "error": "Unauthorized"
    }
    ```

### Get meal log by ID
- **URL**: `/api/v1/meal-logs/{id}`
- **Method**: `GET`
- **Authentication**: Requires JWT token
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Meal log object
- **Error Response**: 
  - **Code**: 404 Not Found
  - **Code**: 401 Unauthorized
    ```json
    {
      "error": "Unauthorized"
    }
    ```

### Get meal logs by user ID
- **URL**: `/api/v1/meal-logs/user`
- **Method**: `GET`
- **Authentication**: Requires JWT token (user ID is extracted from token)
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of meal log objects
- **Error Response**: 
  - **Code**: 500 Internal Server Error
  - **Code**: 401 Unauthorized
    ```json
    {
      "error": "Unauthorized"
    }
    ```

### Get meal logs by user ID and date
- **URL**: `/api/v1/meal-logs/user/date/{date}`
- **Method**: `GET`
- **Authentication**: Requires JWT token (user ID is extracted from token)
- **URL Parameters**: 
  - `date=[YYYY-MM-DD]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of meal log objects
- **Error Response**: 
  - **Code**: 500 Internal Server Error
  - **Code**: 401 Unauthorized
    ```json
    {
      "error": "Unauthorized"
    }
    ```

### Get meal logs by user ID and date range
- **URL**: `/api/v1/meal-logs/user/date-range`
- **Method**: `GET`
- **Authentication**: Requires JWT token (user ID is extracted from token)
- **Query Parameters**:
  - `startDate=[YYYY-MM-DD]`
  - `endDate=[YYYY-MM-DD]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of meal log objects
- **Error Response**: 
  - **Code**: 400 Bad Request or 500 Internal Server Error
  - **Code**: 401 Unauthorized
    ```json
    {
      "error": "Unauthorized"
    }
    ```

### Update a meal log
- **URL**: `/api/v1/meal-logs/{id}`
- **Method**: `PUT`
- **Authentication**: Requires JWT token (only the owner of the meal log can update it)
- **URL Parameters**: `id=[uint]`
- **Request Body**: Same as create meal log
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Updated meal log object
- **Error Response**: 
  - **Code**: 404 Not Found or 500 Internal Server Error
  - **Code**: 401 Unauthorized
    ```json
    {
      "error": "Unauthorized"
    }
    ```
  - **Code**: 403 Forbidden
    ```json
    {
      "error": "You are not allowed to update this meal log"
    }
    ```

### Delete a meal log
- **URL**: `/api/v1/meal-logs/{id}`
- **Method**: `DELETE`
- **Authentication**: Requires JWT token (only the owner of the meal log can delete it)
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "message": "Meal log deleted successfully"
    }
    ```
- **Error Response**: 
  - **Code**: 404 Not Found or 500 Internal Server Error
  - **Code**: 401 Unauthorized
    ```json
    {
      "error": "Unauthorized"
    }
    ```
  - **Code**: 403 Forbidden
    ```json
    {
      "error": "You are not allowed to delete this meal log"
    }
    ```

## Meal Log Items Endpoints

### Create a new meal log item
- **URL**: `/api/v1/meal-log-items`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "meal_log_id": 1,
    "food_id": 1,
    "quantity": 1,
    "quantity_": 1,
    "quantity_grams": 100
  }
  ```
- **Success Response**: 
  - **Code**: 201 Created
  - **Content**: Meal log item object
- **Error Response**: 
  - **Code**: 400 Bad Request or 500 Internal Server Error

### Get meal log item by ID
- **URL**: `/api/v1/meal-log-items/{id}`
- **Method**: `GET`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Meal log item object
- **Error Response**: 
  - **Code**: 404 Not Found

### Get meal log items by meal log ID
- **URL**: `/api/v1/meal-log-items/meal-log/{mealLogId}`
- **Method**: `GET`
- **URL Parameters**: `mealLogId=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of meal log item objects
- **Error Response**: 
  - **Code**: 500 Internal Server Error

### Get meal log items by food ID
- **URL**: `/api/v1/meal-log-items/food/{foodId}`
- **Method**: `GET`
- **URL Parameters**: `foodId=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of meal log item objects
- **Error Response**: 
  - **Code**: 500 Internal Server Error

### Update a meal log item
- **URL**: `/api/v1/meal-log-items/{id}`
- **Method**: `PUT`
- **URL Parameters**: `id=[uint]`
- **Request Body**: Same as create meal log item
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Updated meal log item object
- **Error Response**: 
  - **Code**: 404 Not Found or 500 Internal Server Error

### Delete a meal log item
- **URL**: `/api/v1/meal-log-items/{id}`
- **Method**: `DELETE`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "message": "Meal log item deleted successfully"
    }
    ```
- **Error Response**: 
  - **Code**: 404 Not Found or 500 Internal Server Error

### Delete meal log items by meal log ID
- **URL**: `/api/v1/meal-log-items/meal-log/{mealLogId}`
- **Method**: `DELETE`
- **URL Parameters**: `mealLogId=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "message": "All items for meal log deleted successfully"
    }
    ```
- **Error Response**: 
  - **Code**: 500 Internal Server Error

### Add Items to Meal Log
- **URL**: `/api/v1/meal-logs/:id/items`
- **Method**: `POST`
- **Authentication**: Requires JWT token (Cookie-based authentication with jwt-id cookie)
- **URL Parameters**: `id=[uint]` (Meal Log ID)
- **Request Body**:
  ```json
  {
    "items": [
      {
        "food_id": 1,
        "quantity": 2,
        "quantity_grams": 100.5
      },
      {
        "food_id": 3,
        "quantity": 1,
        "quantity_grams": 50.0
      }
    ]
  }
  ```
- **Success Response**: 
  - **Code**: 201 Created
  - **Content**: Array of created meal log item objects
    ```json
    [
      {
        "id": 1,
        "meal_log_id": 5,
        "food_id": 1,
        "quantity": 2,
        "quantity_grams": 100.5
      },
      {
        "id": 2,
        "meal_log_id": 5,
        "food_id": 3,
        "quantity": 1,
        "quantity_grams": 50.0
      }
    ]
    ```
- **Error Response**: 
  - **Code**: 400 Bad Request
    ```json
    {
      "error": "No items provided"
    }
    ```
  - **Code**: 401 Unauthorized
    ```json
    {
      "error": "Unauthorized"
    }
    ```
  - **Code**: 403 Forbidden
    ```json
    {
      "error": "You are not allowed to modify this meal log"
    }
    ```
  - **Code**: 500 Internal Server Error
    ```json
    {
      "error": "Failed to add items to meal log: <error details>"
    }
    ```

## User Biometrics Endpoints

### Create a new user biometric
- **URL**: `/api/v1/user-biometrics`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "user_id": 1,
    "created_at": "2023-09-22T10:30:45Z",
    "type": "Weight",
    "value": 70.5,
    "unit": "kg"
  }
  ```
- **Success Response**: 
  - **Code**: 201 Created
  - **Content**: User biometric object
- **Error Response**: 
  - **Code**: 400 Bad Request or 500 Internal Server Error

### Get user biometric by ID
- **URL**: `/api/v1/user-biometrics/{id}`
- **Method**: `GET`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: User biometric object
- **Error Response**: 
  - **Code**: 404 Not Found

### Get all user biometrics by user ID
- **URL**: `/api/v1/user-biometrics/user/{userId}`
- **Method**: `GET`
- **URL Parameters**: `userId=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of user biometric objects
- **Error Response**: 
  - **Code**: 500 Internal Server Error

### Get user biometrics by user ID and type
- **URL**: `/api/v1/user-biometrics/user/{userId}/type/{type}`
- **Method**: `GET`
- **URL Parameters**: 
  - `userId=[uint]`
  - `type=[string]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of user biometric objects
- **Error Response**: 
  - **Code**: 500 Internal Server Error

### Get user biometrics by user ID, type and date range
- **URL**: `/api/v1/user-biometrics/user/{userId}/type/{type}/date-range`
- **Method**: `GET`
- **URL Parameters**: 
  - `userId=[uint]`
  - `type=[string]`
- **Query Parameters**:
  - `startDate=[YYYY-MM-DD]`
  - `endDate=[YYYY-MM-DD]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of user biometric objects
- **Error Response**: 
  - **Code**: 400 Bad Request or 500 Internal Server Error

### Get latest user biometric by user ID and type
- **URL**: `/api/v1/user-biometrics/user/{userId}/type/{type}/latest`
- **Method**: `GET`
- **URL Parameters**: 
  - `userId=[uint]`
  - `type=[string]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: User biometric object
- **Error Response**: 
  - **Code**: 404 Not Found

### Update a user biometric
- **URL**: `/api/v1/user-biometrics/{id}`
- **Method**: `PUT`
- **URL Parameters**: `id=[uint]`
- **Request Body**: Same as create user biometric
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Updated user biometric object
- **Error Response**: 
  - **Code**: 404 Not Found or 500 Internal Server Error

### Delete a user biometric
- **URL**: `/api/v1/user-biometrics/{id}`
- **Method**: `DELETE`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "message": "User biometric deleted successfully"
    }
    ```
- **Error Response**: 
  - **Code**: 404 Not Found or 500 Internal Server Error

### Get biometric progress
- **URL**: `/api/v1/user-biometrics/user/{userId}/progress/{type}`
- **Method**: `GET`
- **URL Parameters**: 
  - `userId=[uint]`
  - `type=[string]`
- **Query Parameters**:
  - `startDate=[YYYY-MM-DD]` (optional, defaults to 30 days ago)
  - `endDate=[YYYY-MM-DD]` (optional, defaults to today)
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "type": "weight",
      "unit": "kg",
      "current_value": 75.5,
      "previous_value": 78.0,
      "overall_change": -2.5,
      "percent_change": -3.21,
      "trend": "down",
      "data_points": [
        {
          "date": "2023-10-01T00:00:00Z",
          "value": 78.0,
          "change": 0,
          "trend": "stable"
        },
        {
          "date": "2023-10-15T00:00:00Z",
          "value": 76.5,
          "change": -1.5,
          "trend": "down"
        },
        {
          "date": "2023-10-30T00:00:00Z",
          "value": 75.5,
          "change": -1.0,
          "trend": "down"
        }
      ],
      "start_date": "2023-10-01T00:00:00Z",
      "end_date": "2023-10-30T23:59:59Z"
    }
    ```
- **Error Response**: 
  - **Code**: 400 Bad Request or 500 Internal Server Error

### Get chart data
- **URL**: `/api/v1/user-biometrics/user/{userId}/chart/{type}`
- **Method**: `GET`
- **URL Parameters**: 
  - `userId=[uint]`
  - `type=[string]`
- **Query Parameters**:
  - `startDate=[YYYY-MM-DD]` (optional, defaults to 30 days ago)
  - `endDate=[YYYY-MM-DD]` (optional, defaults to today)
  - `maxPoints=[int]` (optional, defaults to 50)
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "type": "weight",
      "unit": "kg",
      "labels": ["2023-10-01", "2023-10-15", "2023-10-30"],
      "values": [78.0, 76.5, 75.5],
      "start_date": "2023-10-01T00:00:00Z",
      "end_date": "2023-10-30T23:59:59Z"
    }
    ```
- **Error Response**: 
  - **Code**: 400 Bad Request or 500 Internal Server Error

### Get advanced metrics
- **URL**: `/api/v1/user-biometrics/user/{userId}/advanced-metrics`
- **Method**: `GET`
- **URL Parameters**: `userId=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "bmi": 23.5,
      "body_fat_percentage": 15.5,
      "muscle_mass": 65.2,
      "waist_to_hip_ratio": 0.85,
      "body_water_percentage": 60.5,
      "bmi_category": "Normal weight",
      "body_fat_category": "Fitness",
      "health_risk": "Low"
    }
    ```
- **Error Response**: 
  - **Code**: 400 Bad Request or 500 Internal Server Error

### Get biometric summary
- **URL**: `/api/v1/user-biometrics/user/{userId}/summary`
- **Method**: `GET`
- **URL Parameters**: `userId=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "user_id": 1,
      "latest_biometrics": {
        "weight": {
          "id": 15,
          "user_id": 1,
          "created_at": "2023-10-30T00:00:00Z",
          "type": "weight",
          "value": 75.5,
          "unit": "kg"
        },
        "body_fat_percentage": {
          "id": 16,
          "user_id": 1,
          "created_at": "2023-10-28T00:00:00Z",
          "type": "body_fat_percentage",
          "value": 15.5,
          "unit": "%"
        }
      },
      "progress_data": {
        "weight": {
          "type": "weight",
          "unit": "kg",
          "current_value": 75.5,
          "previous_value": 78.0,
          "overall_change": -2.5,
          "percent_change": -3.21,
          "trend": "down",
          "data_points": [],
          "start_date": "2023-10-01T00:00:00Z",
          "end_date": "2023-10-30T23:59:59Z"
        }
      },
      "goals": {
        "goals": [],
        "overall_progress": 0,
        "achieved_goals": 0,
        "total_goals": 0
      },
      "last_updated": "2023-10-30T12:00:00Z"
    }
    ```
- **Error Response**: 
  - **Code**: 400 Bad Request or 500 Internal Server Error

### Get available biometric types for user
- **URL**: `/api/v1/user-biometrics/user/{userId}/types`
- **Method**: `GET`
- **URL Parameters**: `userId=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "types": ["weight", "height", "body_fat_percentage", "muscle_mass"]
    }
    ```
- **Error Response**: 
  - **Code**: 400 Bad Request or 500 Internal Server Error

### Get all supported biometric types
- **URL**: `/api/v1/user-biometrics/types`
- **Method**: `GET`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "Weight": "weight",
      "Height": "height",
      "BodyFatPercentage": "body_fat_percentage",
      "MuscleMass": "muscle_mass",
      "BMI": "bmi",
      "WaistCircumference": "waist_circumference",
      "HipCircumference": "hip_circumference",
      "ChestCircumference": "chest_circumference",
      "ArmCircumference": "arm_circumference",
      "ThighCircumference": "thigh_circumference",
      "BloodPressureSystolic": "blood_pressure_systolic",
      "BloodPressureDiastolic": "blood_pressure_diastolic",
      "RestingHeartRate": "resting_heart_rate",
      "BodyWaterPercentage": "body_water_percentage",
      "BoneDensity": "bone_density"
    }
    ```

## Dashboard Endpoints

### Get User Dashboard
- **URL**: `/api/v1/dashboard`
- **Method**: `GET`
- **Authentication**: Requires JWT token (Cookie-based authentication with jwt-id cookie)
- **Query Parameters**:
  - `date=[YYYY-MM-DD]` (optional, defaults to today)
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "date": "2023-10-25",
      "total_calories": 1875.5,
      "number_of_meals": 3,
      "meal_logs": [
        {
          "id": 1,
          "meal_type": "Breakfast",
          "created_at": "2023-10-25T08:30:00Z",
          "total_calories": 450.5,
          "food_items": [
            {
              "id": 1,
              "food_id": 12,
              "food_name": "Oatmeal",
              "quantity": 1,
              "quantity_grams": 100.0,
              "calories": 350.5
            },
            {
              "id": 2,
              "food_id": 15,
              "food_name": "Apple",
              "quantity": 1,
              "quantity_grams": 100.0,
              "calories": 100.0
            }
          ]
        },
        {
          "id": 2,
          "meal_type": "Lunch",
          "created_at": "2023-10-25T12:30:00Z",
          "total_calories": 750.0,
          "food_items": [
            {
              "id": 3,
              "food_id": 18,
              "food_name": "Chicken Breast",
              "quantity": 1,
              "quantity_grams": 200.0,
              "calories": 330.0
            },
            {
              "id": 4,
              "food_id": 20,
              "food_name": "Brown Rice",
              "quantity": 1,
              "quantity_grams": 150.0,
              "calories": 220.0
            },
            {
              "id": 5,
              "food_id": 25,
              "food_name": "Broccoli",
              "quantity": 1,
              "quantity_grams": 100.0,
              "calories": 200.0
            }
          ]
        },
        {
          "id": 3,
          "meal_type": "Dinner",
          "created_at": "2023-10-25T19:30:00Z",
          "total_calories": 675.0,
          "food_items": [
            {
              "id": 6,
              "food_id": 30,
              "food_name": "Salmon",
              "quantity": 1,
              "quantity_grams": 150.0,
              "calories": 375.0
            },
            {
              "id": 7,
              "food_id": 35,
              "food_name": "Sweet Potato",
              "quantity": 1,
              "quantity_grams": 150.0,
              "calories": 150.0
            },
            {
              "id": 8,
              "food_id": 40,
              "food_name": "Asparagus",
              "quantity": 1,
              "quantity_grams": 100.0,
              "calories": 150.0
            }
          ]
        }
      ],
      "total_macronutrients": {
        "protein": 120.5,
        "carbohydrate": 185.3,
        "fat": 45.7
      }
    }
    ```
- **Error Response**: 
  - **Code**: 401 Unauthorized
    ```json
    {
      "error": "Unauthorized"
    }
    ```
  - **Code**: 400 Bad Request
    ```json
    {
      "error": "Invalid date format. Use YYYY-MM-DD"
    }
    ```
  - **Code**: 500 Internal Server Error
    ```json
    {
      "error": "Failed to get dashboard data: [error message]"
    }
    ``` 