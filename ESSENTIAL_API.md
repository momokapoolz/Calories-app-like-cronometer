# Calories App - Essential API Documentation

This document outlines the core APIs needed for the Calories App nutrition tracking system.

## Authentication Endpoints

### Register a new user
- **URL**: `/register`
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
- **Success Response**: 
  - **Code**: 201 Created
  - **Content**: User object with JWT tokens

### User Login
- **URL**: `/login`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "email": "john@example.com",
    "password": "securepassword"
  }
  ```
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: User object with JWT tokens

### Token Refresh
- **URL**: `/auth/refresh`
- **Method**: `POST`
- **Description**: Generates new access token using refresh token
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: New access token information

### Logout
- **URL**: `/auth/logout`
- **Method**: `POST`
- **Description**: Invalidates both access and refresh tokens
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Logout confirmation

## User Endpoints

### Get User Profile
- **URL**: `/api/profile`
- **Method**: `GET`
- **Authentication**: Required
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: User profile data

### Update User Profile
- **URL**: `/api/users/{id}`
- **Method**: `PUT`
- **Authentication**: Required
- **Request Body**: User profile data
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Updated user profile

## Food Endpoints

### Create a new food
- **URL**: `/api/v1/foods`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "Apple",
    "serving_size_gram": 100,
    "source": "USDA"
  }
  ```
- **Success Response**: 
  - **Code**: 201 Created
  - **Content**: Created food object

### Get all foods
- **URL**: `/api/v1/foods`
- **Method**: `GET`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of food objects

### Get food by ID
- **URL**: `/api/v1/foods/{id}`
- **Method**: `GET`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Food object

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
  - **Content**: Created nutrient object

### Get all nutrients
- **URL**: `/api/v1/nutrients`
- **Method**: `GET`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of nutrient objects

### Get nutrients by category
- **URL**: `/api/v1/nutrients/category/{category}`
- **Method**: `GET`
- **URL Parameters**: `category=[string]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of nutrient objects

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

### Get food nutrients by food ID
- **URL**: `/api/v1/food-nutrients/food/{foodId}`
- **Method**: `GET`
- **URL Parameters**: `foodId=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of food nutrient objects

## Meal Log Endpoints

### Create a new meal log
- **URL**: `/api/v1/meal-logs`
- **Method**: `POST`
- **Authentication**: Required
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
  - **Content**: Created meal log with items

### Get meal logs by user ID and date
- **URL**: `/api/v1/meal-logs/user/date/{date}`
- **Method**: `GET`
- **Authentication**: Required
- **URL Parameters**: 
  - `date=[YYYY-MM-DD]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of meal log objects

### Get meal logs by user ID and date range
- **URL**: `/api/v1/meal-logs/user/date-range`
- **Method**: `GET`
- **Authentication**: Required
- **Query Parameters**:
  - `startDate=[YYYY-MM-DD]`
  - `endDate=[YYYY-MM-DD]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of meal log objects

### Update a meal log
- **URL**: `/api/v1/meal-logs/{id}`
- **Method**: `PUT`
- **Authentication**: Required
- **URL Parameters**: `id=[uint]`
- **Request Body**: Meal log data
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Updated meal log object

### Delete a meal log
- **URL**: `/api/v1/meal-logs/{id}`
- **Method**: `DELETE`
- **Authentication**: Required
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Deletion confirmation

### Add Items to Meal Log
- **URL**: `/api/v1/meal-logs/:id/items`
- **Method**: `POST`
- **Authentication**: Required
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

## User Biometrics Endpoints

### Create a new user biometric
- **URL**: `/api/v1/user-biometrics`
- **Method**: `POST`
- **Authentication**: Required
- **Request Body**:
  ```json
  {
    "type": "Weight",
    "value": 70.5,
    "unit": "kg"
  }
  ```
- **Success Response**: 
  - **Code**: 201 Created
  - **Content**: User biometric object

### Get user biometrics by user ID, type and date range
- **URL**: `/api/v1/user-biometrics/user/{userId}/type/{type}/date-range`
- **Method**: `GET`
- **Authentication**: Required
- **URL Parameters**: 
  - `userId=[uint]`
  - `type=[string]`
- **Query Parameters**:
  - `startDate=[YYYY-MM-DD]`
  - `endDate=[YYYY-MM-DD]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of user biometric objects

### Get latest user biometric by user ID and type
- **URL**: `/api/v1/user-biometrics/user/{userId}/type/{type}/latest`
- **Method**: `GET`
- **Authentication**: Required
- **URL Parameters**: 
  - `userId=[uint]`
  - `type=[string]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: User biometric object

## Dashboard Endpoints

### Get User Dashboard
- **URL**: `/api/v1/dashboard`
- **Method**: `GET`
- **Authentication**: Required
- **Query Parameters**:
  - `date=[YYYY-MM-DD]` (optional, defaults to today)
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Dashboard data including daily summary, meals, and nutrition totals 