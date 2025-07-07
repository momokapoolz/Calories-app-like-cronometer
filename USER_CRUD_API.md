# User CRUD API Documentation

This document provides comprehensive documentation for the User CRUD (Create, Read, Update, Delete) API with JWT authentication. Users can only manage their own information.

**Base URL**: `/api/v1`

## Overview

The User CRUD API allows authenticated users to:
- ✅ **Create** their account (Register)
- ✅ **Read** their profile information
- ✅ **Update** their profile information
- ✅ **Delete** their account
- ✅ **Change** their password

**Security Features:**
- JWT token-based authentication required for all operations (except registration)
- Users can only access and modify their own data
- Input validation and sanitization
- Secure password hashing with bcrypt

---

## Authentication

All endpoints (except registration) require JWT authentication via the `Authorization` header:

```
Authorization: Bearer {access_token_id}
```

You receive the `access_token_id` after successful login or registration.

---

## Endpoints

### 1. Create Account (Register)

**URL**: `POST /api/v1/register`
**Authentication**: Not required
**Description**: Creates a new user account and returns authentication tokens.

**Request Body**:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123",
  "age": 30,
  "gender": "Male",
  "weight": 70.5,
  "height": 175.0,
  "goal": "Weight Loss",
  "activity_level": "Moderate"
}
```

**Success Response (201 Created)**:
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
      "access_token_id": "0097fcb4-abeb-4140-8683-181f7d796755",
      "refresh_token_id": "15aff74e-bafe-410b-a7f7-fa5a96e530a6",
      "expires_in": 3600
    }
  }
}
```

**Error Responses**:
- `400 Bad Request`: Email already in use or invalid data
- `500 Internal Server Error`: Server error during registration

---

### 2. Read Profile

**URL**: `GET /api/v1/profile`
**Authentication**: Required
**Description**: Retrieves the authenticated user's complete profile information.

**Headers**:
```
Authorization: Bearer {access_token_id}
```

**Success Response (200 OK)**:
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

**Error Responses**:
- `401 Unauthorized`: Invalid or missing authentication token
- `404 Not Found`: User not found

---

### 3. Update Profile

**URL**: `PUT /api/v1/profile`
**Authentication**: Required
**Description**: Updates the authenticated user's profile information. Only provided fields will be updated (partial updates supported).

**Headers**:
```
Authorization: Bearer {access_token_id}
```

**Request Body** (all fields optional):
```json
{
  "name": "John Smith",
  "email": "johnsmith@example.com",
  "age": 31,
  "gender": "Male",
  "weight": 68.0,
  "height": 175.0,
  "goal": "Muscle Gain",
  "activity_level": "Active"
}
```

**Success Response (200 OK)**:
```json
{
  "status": "success",
  "message": "Profile updated successfully",
  "data": {
    "user": {
      "id": 1,
      "name": "John Smith",
      "email": "johnsmith@example.com",
      "age": 31,
      "gender": "Male",
      "weight": 68.0,
      "height": 175.0,
      "goal": "Muscle Gain",
      "activity_level": "Active",
      "role": "user",
      "created_at": "2023-09-22T10:30:45Z"
    }
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid data format or email already in use
- `401 Unauthorized`: Invalid or missing authentication token
- `404 Not Found`: User not found
- `500 Internal Server Error`: Failed to update profile

---

### 4. Change Password

**URL**: `PUT /api/v1/password`
**Authentication**: Required
**Description**: Changes the authenticated user's password.

**Headers**:
```
Authorization: Bearer {access_token_id}
```

**Request Body**:
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**Success Response (200 OK)**:
```json
{
  "status": "success",
  "message": "Password changed successfully"
}
```

**Error Responses**:
- `400 Bad Request`: Current password is incorrect or invalid data format
- `401 Unauthorized`: Invalid or missing authentication token
- `404 Not Found`: User not found
- `500 Internal Server Error`: Failed to update password

---

### 5. Delete Account

**URL**: `DELETE /api/v1/account`
**Authentication**: Required
**Description**: Permanently deletes the authenticated user's account.

**Headers**:
```
Authorization: Bearer {access_token_id}
```

**Success Response (200 OK)**:
```json
{
  "status": "success",
  "message": "Account deleted successfully"
}
```

**Error Responses**:
- `401 Unauthorized`: Invalid or missing authentication token
- `404 Not Found`: User not found
- `500 Internal Server Error`: Failed to delete account

---

## Usage Examples

### Complete User Journey

#### 1. Register a new account
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "password": "securepass123",
    "age": 28,
    "gender": "Female",
    "weight": 65.0,
    "height": 168.0,
    "goal": "Weight Loss",
    "activity_level": "Moderate"
  }'
```

#### 2. Get profile information
```bash
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer 0097fcb4-abeb-4140-8683-181f7d796755"
```

#### 3. Update profile (partial update)
```bash
curl -X PUT http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer 0097fcb4-abeb-4140-8683-181f7d796755" \
  -H "Content-Type: application/json" \
  -d '{
    "weight": 63.5,
    "goal": "Maintenance"
  }'
```

#### 4. Change password
```bash
curl -X PUT http://localhost:8080/api/v1/password \
  -H "Authorization: Bearer 0097fcb4-abeb-4140-8683-181f7d796755" \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "securepass123",
    "new_password": "newsecurepass456"
  }'
```

#### 5. Delete account
```bash
curl -X DELETE http://localhost:8080/api/v1/account \
  -H "Authorization: Bearer 0097fcb4-abeb-4140-8683-181f7d796755"
```

---

## Security Features

### Authentication & Authorization
- **JWT Token-based Authentication**: All protected endpoints require valid JWT tokens
- **User Isolation**: Users can only access and modify their own data
- **Context Validation**: User ID is extracted from JWT token and validated on every request

### Input Validation
- **Request Binding**: JSON request bodies are validated using Go struct tags
- **Email Format**: Email addresses must be in valid format
- **Password Requirements**: Minimum 6 characters for passwords
- **Data Sanitization**: All input data is validated before processing

### Password Security
- **bcrypt Hashing**: Passwords are hashed using bcrypt with default cost
- **Current Password Verification**: Password changes require current password verification
- **Secure Storage**: Only password hashes are stored, never plain text passwords

### Data Integrity
- **Email Uniqueness**: Email addresses must be unique across all users
- **Partial Updates**: Update operations only modify provided fields
- **Atomic Operations**: Database operations are atomic to prevent data corruption

---

## Error Handling

All endpoints follow a consistent error response format:

```json
{
  "status": "error",
  "message": "Error description",
  "error": "Detailed error information (optional)"
}
```

### Common HTTP Status Codes
- `200 OK`: Successful operation
- `201 Created`: Account created successfully
- `400 Bad Request`: Invalid request data or validation error
- `401 Unauthorized`: Authentication required or invalid token
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

---

## Integration with Frontend

### JavaScript/TypeScript Example

```javascript
class UserAPI {
  constructor(baseURL, token) {
    this.baseURL = baseURL;
    this.token = token;
  }

  async getProfile() {
    const response = await fetch(`${this.baseURL}/profile`, {
      headers: {
        'Authorization': `Bearer ${this.token}`,
        'Content-Type': 'application/json'
      }
    });
    return response.json();
  }

  async updateProfile(profileData) {
    const response = await fetch(`${this.baseURL}/profile`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${this.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(profileData)
    });
    return response.json();
  }

  async changePassword(currentPassword, newPassword) {
    const response = await fetch(`${this.baseURL}/password`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${this.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        current_password: currentPassword,
        new_password: newPassword
      })
    });
    return response.json();
  }

  async deleteAccount() {
    const response = await fetch(`${this.baseURL}/account`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${this.token}`
      }
    });
    return response.json();
  }
}

// Usage
const userAPI = new UserAPI('http://localhost:8080/api/v1', 'your-token-here');
const profile = await userAPI.getProfile();
```

---

## Database Schema

The User model includes the following fields:

```sql
CREATE TABLE "User" (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  age BIGINT NOT NULL,
  gender VARCHAR(255) NOT NULL,
  weight DOUBLE PRECISION NOT NULL,
  height DOUBLE PRECISION NOT NULL,
  goal VARCHAR(255) NOT NULL,
  activity_level VARCHAR(255) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  role VARCHAR(255) NOT NULL DEFAULT 'user'
);
```

---

## Best Practices

### For Frontend Developers
1. **Token Management**: Store JWT tokens securely and handle token expiration
2. **Error Handling**: Implement proper error handling for all API responses
3. **Input Validation**: Validate input on the frontend before sending to API
4. **Loading States**: Show loading indicators during API requests
5. **Security**: Never store passwords in frontend code or logs

### For Backend Integration
1. **Rate Limiting**: Implement rate limiting for API endpoints
2. **Logging**: Log important events (registration, profile updates, account deletion)
3. **Monitoring**: Monitor API usage and error rates
4. **Backup**: Ensure user data is properly backed up before deletion
5. **GDPR Compliance**: Consider data protection regulations when implementing deletion

---

## Testing

### Manual Testing Checklist
- [ ] User registration with valid data
- [ ] User registration with duplicate email (should fail)
- [ ] Profile retrieval with valid token
- [ ] Profile update with partial data
- [ ] Profile update with duplicate email (should fail)
- [ ] Password change with correct current password
- [ ] Password change with incorrect current password (should fail)
- [ ] Account deletion with valid token
- [ ] All endpoints with invalid/expired tokens (should fail)

### Automated Testing
Consider implementing automated tests for:
- Unit tests for individual controller methods
- Integration tests for complete API flows
- Authentication and authorization tests
- Input validation tests
- Error handling tests

---

This completes the User CRUD API documentation. The API provides a secure, complete solution for user self-management with proper authentication and authorization controls. 