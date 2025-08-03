import React, { useState } from 'react';
import { useAuth } from '../AuthContext';
import { authAPI } from '../api';
import './Login.css';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isLoginLoading, setIsLoginLoading] = useState(false);
  const { login } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoginLoading(true);

    try {
      console.log('Attempting login with:', { username, password: 'hidden' });
      const response = await authAPI.login(username, password);
      console.log('Login response:', response);
      const { token, user } = response.data;
      login(user, token);
    } catch (error) {
      console.error('Login error:', error);
      console.error('Error response:', error.response);
      if (error.response) {
        window.alert(`Login failed: ${error.response.data.error || error.response.statusText}`);
      } else if (error.request) {
        window.alert('Login failed: Unable to connect to server. Please ensure the backend is running on http://localhost:8080');
      } else {
        window.alert('Login failed: ' + error.message);
      }
    } finally {
      setIsLoginLoading(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-form">
        <h2>E-commerce Login</h2>
        <p className="admin-info">Use admin credentials: <strong>admin / Admin@123</strong></p>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Username:</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              disabled={isLoginLoading}
            />
          </div>
          <div className="form-group">
            <label>Password:</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              disabled={isLoginLoading}
            />
          </div>
          <div className="button-group">
            <button type="submit" disabled={isLoginLoading}>
              {isLoginLoading ? 'Logging in...' : 'Login'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default Login;
