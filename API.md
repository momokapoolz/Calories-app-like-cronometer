# Calories App API Documentation

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
  - **Content**: 
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
        },
        "tokens": {
          "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
          "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
          "expires_in": 3600
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
          "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
          "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
          "expires_in": 3600
        }
      }
    }
    ```
- **Error Response**: 
  - **Code**: 401 Unauthorized
    ```json
    {
      "status": "error",
      "message": "Invalid credentials"
    }
    ```

### Refresh Token
- **URL**: `/auth/refresh`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 3600
    }
    ```
- **Error Response**: 
  - **Code**: 401 Unauthorized
    ```json
    {
      "status": "error",
      "message": "Invalid refresh token"
    }
    ```

## Protected Endpoints

### Get User Profile
- **URL**: `/api/profile`
- **Method**: `GET`
- **Headers**: 
  - `Authorization: Bearer {access_token}`
- **Success Response**: 
  - **Code**: 200 OK
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
- **Error Response**: 
  - **Code**: 401 Unauthorized
    ```json
    {
      "error": "Authorization header required"
    }
    ```

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
- **URL**: `/api/users`
- **Method**: `GET`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of user objects

### Get user by ID
- **URL**: `/api/users/{id}`
- **Method**: `GET`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: User object
- **Error Response**: 
  - **Code**: 404 Not Found

### Get user by email
- **URL**: `/api/users/email/{email}`
- **Method**: `GET`
- **URL Parameters**: `email=[string]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: User object
- **Error Response**: 
  - **Code**: 404 Not Found

### Get users by role
- **URL**: `/api/users/role/{role}`
- **Method**: `GET`
- **URL Parameters**: `role=[string]`
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Array of user objects

### Update a user
- **URL**: `/api/users/{id}`
- **Method**: `PUT`
- **URL Parameters**: `id=[uint]`
- **Request Body**: Same as create user
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: Updated user object
- **Error Response**: 
  - **Code**: 404 Not Found

### Delete a user
- **URL**: `/api/users/{id}`
- **Method**: `DELETE`
- **URL Parameters**: `id=[uint]`
- **Success Response**: 
  - **Code**: 204 No Content
- **Error Response**: 
  - **Code**: 404 Not Found or 500 Internal Server Error 