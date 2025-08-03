import axios from 'axios';

// Determine the backend URL based on how the frontend is accessed
const getBackendURL = () => {
  if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
    return 'http://localhost:8080';
  } else {
    // Use the same IP as the frontend but port 8080
    return `http://${window.location.hostname}:8080`;
  }
};

const API_BASE_URL = getBackendURL();

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor to include auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  login: (username, password) => 
    api.post('/users/login', { username, password }),
};

// Items API
export const itemsAPI = {
  getAll: () => api.get('/items'),
  create: (name, status = 'active') => 
    api.post('/items', { name, status }),
};

// Cart API
export const cartAPI = {
  addItem: (item_id) => api.post('/carts', { item_id }),
  getUserCart: () => api.get('/carts/user'),
  getAll: () => api.get('/carts'),
};

// Orders API
export const ordersAPI = {
  create: () => api.post('/orders'),
  getUserOrders: () => api.get('/orders/user'),
  getAll: () => api.get('/orders'),
};

export default api;
