import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useAuth } from '../AuthContext';
import { itemsAPI, cartAPI, ordersAPI } from '../api';
import { toast } from 'react-toastify';
import { 
  ShoppingCart, 
  User, 
  LogOut, 
  Package, 
  CreditCard, 
  Plus,
  Star,
  Heart,
  Filter,
  Search,
  Grid,
  List
} from 'lucide-react';
import './ModernItemsList.css';

const ModernItemsList = () => {
  const [items, setItems] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [cartCount, setCartCount] = useState(0);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [viewMode, setViewMode] = useState('grid');
  const [favorites, setFavorites] = useState(new Set());
  const [showCart, setShowCart] = useState(false);
  const [cartItems, setCartItems] = useState([]);
  const { user, logout } = useAuth();

  useEffect(() => {
    fetchItems();
    fetchCartItems();
  }, []);

  const fetchItems = async () => {
    setIsLoading(true);
    try {
      const response = await itemsAPI.getAll();
      setItems(response.data);
    } catch (error) {
      toast.error('Failed to fetch items');
    } finally {
      setIsLoading(false);
    }
  };

  const fetchCartItems = async () => {
    try {
      const response = await cartAPI.getUserCart();
      const cart = response.data;
      setCartItems(cart.cart_items || []);
      setCartCount(cart.cart_items?.length || 0);
    } catch (error) {
      console.log('Cart fetch error:', error);
    }
  };

  const addToCart = async (itemId) => {
    try {
      await cartAPI.addItem(itemId);
      toast.success('✨ Item added to cart!');
      fetchCartItems(); // Refresh cart count
    } catch (error) {
      if (error.response?.status === 409) {
        toast.warning('Item already in cart');
      } else {
        toast.error('Failed to add item to cart');
      }
    }
  };

  const toggleFavorite = (itemId) => {
    const newFavorites = new Set(favorites);
    if (favorites.has(itemId)) {
      newFavorites.delete(itemId);
      toast.info('Removed from favorites');
    } else {
      newFavorites.add(itemId);
      toast.success('Added to favorites ❤️');
    }
    setFavorites(newFavorites);
  };

  const checkout = async () => {
    try {
      await ordersAPI.createOrder();
      toast.success('🎉 Order placed successfully!');
      fetchCartItems(); // Refresh cart
    } catch (error) {
      if (error.response?.status === 400) {
        toast.error('Cart is empty');
      } else {
        toast.error('Failed to create order');
      }
    }
  };

  const showOrderHistory = async () => {
    try {
      const response = await ordersAPI.getUserOrders();
      const orders = response.data;
      
      if (orders && orders.length > 0) {
        const orderDetails = orders.map(order => 
          `Order ID: ${order.id}, Status: ${order.status}, Created: ${new Date(order.created_at).toLocaleDateString()}`
        ).join('\n');
        window.alert(`Order History:\n${orderDetails}`);
      } else {
        toast.info('No order history found');
      }
    } catch (error) {
      toast.error('Failed to fetch order history');
    }
  };

  const filteredItems = items.filter(item => {
    const matchesSearch = item.name.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesCategory = selectedCategory === 'all' || 
      item.name.toLowerCase().includes(selectedCategory.toLowerCase());
    return matchesSearch && matchesCategory;
  });

  const categories = ['all', 'laptop', 'phone', 'headphones', 'keyboard', 'mouse', 'monitor', 'tablet', 'webcam'];

  if (isLoading) {
    return (
      <div className="modern-loading">
        <motion.div
          className="loading-spinner"
          animate={{ rotate: 360 }}
          transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
        />
        <p>Loading amazing products...</p>
      </div>
    );
  }

  return (
    <div className="modern-container">
      {/* Modern Header */}
      <motion.header 
        className="modern-header"
        initial={{ y: -50, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.5 }}
      >
        <div className="header-content">
          <div className="header-left">
            <motion.h1 
              className="app-title"
              whileHover={{ scale: 1.05 }}
            >
              🛒 ModernStore
            </motion.h1>
            <p className="welcome-text">Welcome back, {user?.username}!</p>
          </div>
          
          <div className="header-actions">
            <motion.button
              className="header-btn cart-btn"
              onClick={() => setShowCart(!showCart)}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              <ShoppingCart size={20} />
              <span>Cart</span>
              {cartCount > 0 && <span className="cart-badge">{cartCount}</span>}
            </motion.button>
            
            <motion.button
              className="header-btn"
              onClick={showOrderHistory}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              <Package size={20} />
              <span>Orders</span>
            </motion.button>
            
            <motion.button
              className="header-btn checkout-btn"
              onClick={checkout}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              <CreditCard size={20} />
              <span>Checkout</span>
            </motion.button>
            
            <motion.button
              className="header-btn logout-btn"
              onClick={logout}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              <LogOut size={20} />
              <span>Logout</span>
            </motion.button>
          </div>
        </div>
      </motion.header>

      {/* Modern Search and Filters */}
      <motion.div 
        className="modern-filters"
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.5, delay: 0.2 }}
      >
        <div className="search-section">
          <div className="search-container">
            <Search size={20} className="search-icon" />
            <input
              type="text"
              placeholder="Search products..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="search-input"
            />
          </div>
          
          <div className="view-toggle">
            <button
              className={`view-btn ${viewMode === 'grid' ? 'active' : ''}`}
              onClick={() => setViewMode('grid')}
            >
              <Grid size={18} />
            </button>
            <button
              className={`view-btn ${viewMode === 'list' ? 'active' : ''}`}
              onClick={() => setViewMode('list')}
            >
              <List size={18} />
            </button>
          </div>
        </div>

        <div className="category-filters">
          {categories.map(category => (
            <motion.button
              key={category}
              className={`category-btn ${selectedCategory === category ? 'active' : ''}`}
              onClick={() => setSelectedCategory(category)}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              {category === 'all' ? 'All Products' : category.charAt(0).toUpperCase() + category.slice(1)}
            </motion.button>
          ))}
        </div>
      </motion.div>

      {/* Modern Products Grid */}
      <motion.div 
        className={`products-container ${viewMode}`}
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.5, delay: 0.4 }}
      >
        <AnimatePresence>
          {filteredItems.length === 0 ? (
            <motion.div 
              className="no-products"
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.8 }}
            >
              <Package size={64} />
              <h3>No products found</h3>
              <p>Try adjusting your search or filters</p>
            </motion.div>
          ) : (
            filteredItems.map((item, index) => (
              <motion.div
                key={item.id}
                className="modern-product-card"
                initial={{ opacity: 0, y: 50 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -50 }}
                transition={{ duration: 0.5, delay: index * 0.1 }}
                whileHover={{ y: -5, scale: 1.02 }}
                layout
              >
                <div className="product-image-container">
                  <img 
                    src={`http://localhost:8080${item.image}`} 
                    alt={item.name}
                    className="product-image"
                    onError={(e) => {
                      e.target.src = 'https://via.placeholder.com/300x300/e9ecef/6c757d?text=' + encodeURIComponent(item.name);
                      e.target.onerror = null;
                    }}
                  />
                  <motion.button
                    className={`favorite-btn ${favorites.has(item.id) ? 'active' : ''}`}
                    onClick={() => toggleFavorite(item.id)}
                    whileHover={{ scale: 1.1 }}
                    whileTap={{ scale: 0.9 }}
                  >
                    <Heart size={18} fill={favorites.has(item.id) ? '#ff4757' : 'none'} />
                  </motion.button>
                </div>
                
                <div className="product-content">
                  <h3 className="product-name">{item.name}</h3>
                  <div className="product-rating">
                    {[...Array(5)].map((_, i) => (
                      <Star key={i} size={14} fill="#ffd700" />
                    ))}
                    <span className="rating-text">4.8 (124 reviews)</span>
                  </div>
                  <p className="product-status">
                    <span className={`status-badge ${item.status}`}>
                      {item.status}
                    </span>
                  </p>
                  <motion.button 
                    className="add-to-cart-btn"
                    onClick={() => addToCart(item.id)}
                    whileHover={{ scale: 1.05 }}
                    whileTap={{ scale: 0.95 }}
                  >
                    <Plus size={18} />
                    Add to Cart
                  </motion.button>
                </div>
              </motion.div>
            ))
          )}
        </AnimatePresence>
      </motion.div>

      {/* Modern Cart Sidebar */}
      <AnimatePresence>
        {showCart && (
          <motion.div
            className="cart-overlay"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={() => setShowCart(false)}
          >
            <motion.div
              className="cart-sidebar"
              initial={{ x: '100%' }}
              animate={{ x: 0 }}
              exit={{ x: '100%' }}
              transition={{ type: 'spring', damping: 20 }}
              onClick={(e) => e.stopPropagation()}
            >
              <div className="cart-header">
                <h3>Your Cart</h3>
                <button 
                  className="close-cart"
                  onClick={() => setShowCart(false)}
                >
                  ×
                </button>
              </div>
              
              <div className="cart-content">
                {cartItems.length === 0 ? (
                  <div className="empty-cart">
                    <ShoppingCart size={48} />
                    <p>Your cart is empty</p>
                  </div>
                ) : (
                  <>
                    <div className="cart-items">
                      {cartItems.map((cartItem) => (
                        <div key={cartItem.id} className="cart-item">
                          <img 
                            src={`http://localhost:8080${cartItem.item.image}`} 
                            alt={cartItem.item.name}
                            className="cart-item-image"
                            onError={(e) => {
                              e.target.src = 'https://via.placeholder.com/60x60/e9ecef/6c757d?text=IMG';
                            }}
                          />
                          <div className="cart-item-details">
                            <h4>{cartItem.item.name}</h4>
                            <p>Status: {cartItem.item.status}</p>
                          </div>
                        </div>
                      ))}
                    </div>
                    
                    <motion.button 
                      className="checkout-cart-btn"
                      onClick={checkout}
                      whileHover={{ scale: 1.05 }}
                      whileTap={{ scale: 0.95 }}
                    >
                      <CreditCard size={18} />
                      Checkout ({cartItems.length} items)
                    </motion.button>
                  </>
                )}
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default ModernItemsList;
