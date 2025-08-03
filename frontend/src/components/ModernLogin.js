import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { useAuth } from '../AuthContext';
import { Lock, User, Eye, EyeOff, ShoppingBag, Sparkles } from 'lucide-react';
import { toast } from 'react-toastify';
import './ModernLogin.css';

const ModernLogin = () => {
  const [formData, setFormData] = useState({
    username: '',
    password: ''
  });
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    
    try {
      await login(formData.username, formData.password);
      toast.success('Welcome to ModernStore! 🎉');
    } catch (error) {
      toast.error(error.response?.data?.error || 'Login failed');
    } finally {
      setIsLoading(false);
    }
  };

  const quickLogin = () => {
    setFormData({
      username: 'admin',
      password: 'Admin@123'
    });
    toast.info('Quick login credentials filled! 🚀');
  };

  return (
    <div className="modern-login-container">
      {/* Animated Background */}
      <div className="background-animation">
        <div className="floating-shape shape-1"></div>
        <div className="floating-shape shape-2"></div>
        <div className="floating-shape shape-3"></div>
        <div className="floating-shape shape-4"></div>
        <div className="floating-shape shape-5"></div>
      </div>

      {/* Login Card */}
      <motion.div
        className="login-card"
        initial={{ opacity: 0, y: 50, scale: 0.9 }}
        animate={{ opacity: 1, y: 0, scale: 1 }}
        transition={{ 
          duration: 0.8, 
          ease: "easeOut",
          type: "spring",
          stiffness: 100
        }}
      >
        {/* Header */}
        <motion.div
          className="login-header"
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.2 }}
        >
          <motion.div
            className="logo"
            whileHover={{ scale: 1.1, rotate: 5 }}
            transition={{ type: "spring", stiffness: 300 }}
          >
            <ShoppingBag size={48} />
            <Sparkles className="sparkle" size={24} />
          </motion.div>
          <h1>ModernStore</h1>
          <p>Your Premium Shopping Experience</p>
        </motion.div>

        {/* Admin Info Banner */}
        <motion.div
          className="admin-info"
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6, delay: 0.4 }}
        >
          <div className="admin-badge">
            <Lock size={16} />
            <span>Admin Portal</span>
          </div>
          <p>This is an admin-only system. Use the credentials below to access the store.</p>
          
          <motion.button
            className="quick-login-btn"
            onClick={quickLogin}
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <Sparkles size={16} />
            Quick Login
          </motion.button>
        </motion.div>

        {/* Login Form */}
        <motion.form
          className="login-form"
          onSubmit={handleSubmit}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.6 }}
        >
          <div className="form-group">
            <div className="input-container">
              <User className="input-icon" size={20} />
              <input
                type="text"
                name="username"
                placeholder="Username"
                value={formData.username}
                onChange={handleChange}
                className="form-input"
                required
              />
            </div>
          </div>

          <div className="form-group">
            <div className="input-container">
              <Lock className="input-icon" size={20} />
              <input
                type={showPassword ? "text" : "password"}
                name="password"
                placeholder="Password"
                value={formData.password}
                onChange={handleChange}
                className="form-input"
                required
              />
              <button
                type="button"
                className="password-toggle"
                onClick={() => setShowPassword(!showPassword)}
              >
                {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
              </button>
            </div>
          </div>

          <motion.button
            type="submit"
            className="login-btn"
            disabled={isLoading}
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
          >
            {isLoading ? (
              <motion.div
                className="loading-spinner"
                animate={{ rotate: 360 }}
                transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
              />
            ) : (
              <>
                <Lock size={18} />
                Access Store
              </>
            )}
          </motion.button>
        </motion.form>

        {/* Credentials Display */}
        <motion.div
          className="credentials-display"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.8 }}
        >
          <h4>Demo Credentials</h4>
          <div className="credential-item">
            <span className="label">Username:</span>
            <code>admin</code>
          </div>
          <div className="credential-item">
            <span className="label">Password:</span>
            <code>Admin@123</code>
          </div>
        </motion.div>

        {/* Footer */}
        <motion.div
          className="login-footer"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.6, delay: 1 }}
        >
          <p>Powered by Modern Technology Stack</p>
          <div className="tech-stack">
            <span>React</span>
            <span>Go</span>
            <span>JWT</span>
          </div>
        </motion.div>
      </motion.div>
    </div>
  );
};

export default ModernLogin;
