# User and Authentication API Documentation

## Authentication Endpoints

### Register a new user
- **URL**: `/api/v1/register`
- **Method**: `POST`
- **Description**: Creates a new user account and returns authentication tokens
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
          "access_token_id": 12345678901234,
          "refresh_token_id": 98765432109876,
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
  - **Code**: 400 Bad Request
    ```json
    {
      "status": "error",
      "message": "Invalid request format",
      "error": "Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
    }
    ```
  - **Code**: 500 Internal Server Error
    ```json
    {
      "status": "error",
      "message": "Failed to create user account",
      "error": "Database connection error"
    }
    ```

### User Login
- **URL**: `/api/v1/login`
- **Method**: `POST`
- **Description**: Authenticates a user and returns JWT tokens
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
          "access_token_id": 12345678901234,
          "refresh_token_id": 98765432109876,
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
  - **Code**: 500 Internal Server Error
    ```json
    {
      "status": "error",
      "message": "Failed to generate authentication token"
    }
    ```

### Auth Login (Cookie-Based)
- **URL**: `/api/v1/auth/login`
- **Method**: `POST`
- **Description**: Authenticates a user and sets secure cookies for session management
- **Request Body**:
  ```json
  {
    "email": "john@example.com",
    "password": "securepassword"
  }
  ```
- **Success Response**: 
  - **Code**: 200 OK
  - **Cookies Set**: 
    - `jwt-id` - Contains the access token ID (HttpOnly)
    - `refresh-id` - Contains the refresh token ID (HttpOnly)
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
        "expires_in": 3600
      }
    }
    ```
- **Error Response**: 
  - **Code**: 401 Unauthorized
    ```json
    {
      "error": "Invalid credentials"
    }
    ```
  - **Code**: 500 Internal Server Error
    ```json
    {
      "error": "Failed to retrieve token"
    }
    ```

### Token Refresh
- **URL**: `/api/v1/auth/refresh`
- **Method**: `POST`
- **Description**: Generates new access token using refresh token stored in cookies
- **Required Cookies**:
  - `refresh-id` - Contains the refresh token ID
- **Success Response**: 
  - **Code**: 200 OK
  - **Cookies Updated**: 
    - `jwt-id` - Contains the new access token ID (HttpOnly)
  - **Content**: 
    ```json
    {
      "message": "Token refreshed successfully",
      "expires_in": 3600
    }
    ```
- **Error Response**: 
  - **Code**: 400 Bad Request
    ```json
    {
      "error": "Refresh token ID cookie required"
    }
    ```
  - **Code**: 401 Unauthorized
    ```json
    {
      "error": "Invalid refresh token"
    }
    ```

### Logout
- **URL**: `/api/v1/auth/logout`
- **Method**: `POST`
- **Description**: Invalidates the current token and clears authentication cookies
- **Authentication**: Requires valid session/token
- **Success Response**: 
  - **Code**: 200 OK
  - **Cookies Cleared**: 
    - `jwt-id` 
    - `refresh-id`
  - **Content**: 
    ```json
    {
      "message": "Logged out successfully"
    }
    ```
- **Error Response**: 
  - **Code**: 400 Bad Request
    ```json
    {
      "error": "No active session"
    }
    ```
  - **Code**: 500 Internal Server Error
    ```json
    {
      "error": "Failed to invalidate token"
    }
    ```

## User Endpoints

### Get User Profile
- **URL**: `/api/v1/profile`
- **Method**: `GET`
- **Description**: Retrieves the authenticated user's profile information
- **Authentication**: Requires valid JWT token
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
      "status": "error",
      "message": "Not authenticated"
    }
    ```
  - **Code**: 404 Not Found
    ```json
    {
      "status": "error",
      "message": "User not found"
    }
    ```

### Update User Password
- **URL**: `/api/v1/user/password/update`
- **Method**: `POST`
- **Description**: Allows users to update their password
- **Authentication**: Requires valid JWT token
- **Request Body**:
  ```json
  {
    "current_password": "oldpassword",
    "new_password": "newpassword"
  }
  ```
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "status": "success",
      "message": "Password updated successfully"
    }
    ```
- **Error Response**: 
  - **Code**: 401 Unauthorized
    ```json
    {
      "status": "error",
      "message": "Current password is incorrect"
    }
    ```
  - **Code**: 400 Bad Request
    ```json
    {
      "status": "error",
      "message": "Invalid request format",
      "error": "Key: 'UpdatePasswordRequest.NewPassword' Error:Field validation for 'NewPassword' failed on the 'min' tag"
    }
    ```
  - **Code**: 500 Internal Server Error
    ```json
    {
      "status": "error",
      "message": "Failed to update password",
      "error": "Database error"
    }
    ```

## Admin Endpoints

### Admin Update User Password
- **URL**: `/api/v1/admin/user/password/update`
- **Method**: `POST`
- **Description**: Allows administrators to update any user's password
- **Authentication**: Requires valid JWT token with admin role
- **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "new_password": "newpassword"
  }
  ```
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "status": "success",
      "message": "Password updated successfully"
    }
    ```
- **Error Response**: 
  - **Code**: 401 Unauthorized
    ```json
    {
      "status": "error",
      "message": "Not authenticated"
    }
    ```
  - **Code**: 403 Forbidden
    ```json
    {
      "status": "error",
      "message": "Insufficient permissions"
    }
    ```
  - **Code**: 400 Bad Request
    ```json
    {
      "status": "error",
      "message": "Invalid request format",
      "error": "Key: 'AdminUpdatePasswordRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
    }
    ```
  - **Code**: 500 Internal Server Error
    ```json
    {
      "status": "error",
      "message": "Failed to update password",
      "error": "Database error"
    }
    ``` 