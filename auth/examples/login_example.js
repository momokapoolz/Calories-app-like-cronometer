// Example of how to use the login API from a client application

/**
 * Example login function that authenticates a user and returns the JWT token
 * @param {string} email - User's email
 * @param {string} password - User's password
 * @returns {Promise<Object>} Token and user data
 */
async function loginUser(email, password) {
  try {
    const response = await fetch('http://localhost:8080/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: email,
        password: password,
      }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Login failed');
    }

    const data = await response.json();
    
    // Store tokens in localStorage or secure cookie
    localStorage.setItem('access_token', data.data.tokens.access_token);
    localStorage.setItem('refresh_token', data.data.tokens.refresh_token);
    
    // Return the user data and tokens
    return {
      user: data.data.user,
      tokens: data.data.tokens,
    };
  } catch (error) {
    console.error('Login error:', error);
    throw error;
  }
}

/**
 * Example of making an authenticated API request using the JWT token
 * @param {string} url - API endpoint
 * @returns {Promise<Object>} Response data
 */
async function makeAuthenticatedRequest(url) {
  try {
    // Get the token from storage
    const token = localStorage.getItem('access_token');
    
    if (!token) {
      throw new Error('No authentication token found');
    }
    
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    });
    
    // If token is expired or invalid
    if (response.status === 401) {
      // Attempt to refresh the token
      const refreshed = await refreshToken();
      if (refreshed) {
        // Retry the request with the new token
        return makeAuthenticatedRequest(url);
      } else {
        throw new Error('Session expired, please login again');
      }
    }
    
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Request failed');
    }
    
    return response.json();
  } catch (error) {
    console.error('Request error:', error);
    throw error;
  }
}

/**
 * Example of refreshing an expired token
 * @returns {Promise<boolean>} Success status
 */
async function refreshToken() {
  try {
    const refreshToken = localStorage.getItem('refresh_token');
    
    if (!refreshToken) {
      return false;
    }
    
    const response = await fetch('http://localhost:8080/auth/refresh', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        refresh_token: refreshToken,
      }),
    });
    
    if (!response.ok) {
      // Clear tokens if refresh failed
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      return false;
    }
    
    const data = await response.json();
    
    // Update tokens in storage
    localStorage.setItem('access_token', data.access_token);
    localStorage.setItem('refresh_token', data.refresh_token);
    
    return true;
  } catch (error) {
    console.error('Token refresh error:', error);
    return false;
  }
}

/**
 * Example of logging out a user
 */
function logoutUser() {
  // Clear tokens from storage
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
}

// Example usage:
// Login example
document.getElementById('login-form')?.addEventListener('submit', async (e) => {
  e.preventDefault();
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;
  
  try {
    const result = await loginUser(email, password);
    console.log('Login successful:', result);
    
    // Update UI or redirect
    window.location.href = '/dashboard';
  } catch (error) {
    console.error('Login failed:', error);
    // Show error message to user
  }
});

// Protected API call example
async function getUserProfile() {
  try {
    // Updated path to avoid conflicts
    const profileData = await makeAuthenticatedRequest('http://localhost:8080/api/auth/profile');
    console.log('User profile:', profileData);
    
    // Update UI with profile data
    document.getElementById('user-name').textContent = profileData.data.user.name;
  } catch (error) {
    console.error('Failed to get profile:', error);
    
    // If authentication error, redirect to login
    if (error.message === 'Session expired, please login again') {
      window.location.href = '/login';
    }
  }
} 