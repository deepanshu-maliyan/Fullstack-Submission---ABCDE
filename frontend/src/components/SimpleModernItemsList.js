import React, { useState, useEffect } from 'react';
import { useAuth } from '../AuthContext';
import { itemsAPI, cartAPI, ordersAPI } from '../api';
import { toast } from 'react-toastify';
import './SimpleModernItemsList.css';

const SimpleModernItemsList = () => {
  const [items, setItems] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [cartCount, setCartCount] = useState(0);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [cartItems, setCartItems] = useState([]);
  const [showCart, setShowCart] = useState(false);
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
      console.error('Fetch items error:', error);
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
      console.log('Cart fetch error (this is normal):', error);
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

  const checkout = async () => {
    try {
      await ordersAPI.createOrder();
      toast.success('🎉 Order placed successfully!');
      fetchCartItems(); // Refresh cart
      setShowCart(false);
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
      <div className="simple-modern-loading">
        <div className="loading-spinner"></div>
        <p>Loading amazing products...</p>
      </div>
    );
  }

  return (
    <div className="simple-modern-container">
      {/* Header */}
      <header className="simple-modern-header">
        <div className="header-content">
          <div className="header-left">
            <h1 className="app-title">🛒 ModernStore</h1>
            <p className="welcome-text">Welcome back, {user?.username}!</p>
          </div>
          
          <div className="header-actions">
            <button
              className="header-btn cart-btn"
              onClick={() => setShowCart(!showCart)}
            >
              🛒 Cart {cartCount > 0 && <span className="cart-badge">{cartCount}</span>}
            </button>
            
            <button
              className="header-btn"
              onClick={showOrderHistory}
            >
              📦 Orders
            </button>
            
            <button
              className="header-btn checkout-btn"
              onClick={checkout}
            >
              💳 Checkout
            </button>
            
            <button
              className="header-btn logout-btn"
              onClick={logout}
            >
              🚪 Logout
            </button>
          </div>
        </div>
      </header>

      {/* Search and Filters */}
      <div className="simple-modern-filters">
        <div className="search-section">
          <input
            type="text"
            placeholder="🔍 Search products..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="search-input"
          />
        </div>

        <div className="category-filters">
          {categories.map(category => (
            <button
              key={category}
              className={`category-btn ${selectedCategory === category ? 'active' : ''}`}
              onClick={() => setSelectedCategory(category)}
            >
              {category === 'all' ? 'All Products' : category.charAt(0).toUpperCase() + category.slice(1)}
            </button>
          ))}
        </div>
      </div>

      {/* Products Grid */}
      <div className="products-container">
        {filteredItems.length === 0 ? (
          <div className="no-products">
            <h3>No products found</h3>
            <p>Try adjusting your search or filters</p>
          </div>
        ) : (
          filteredItems.map((item) => (
            <div key={item.id} className="simple-product-card">
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
              </div>
              
              <div className="product-content">
                <h3 className="product-name">{item.name}</h3>
                <div className="product-rating">
                  ⭐⭐⭐⭐⭐ <span className="rating-text">4.8 (124 reviews)</span>
                </div>
                <p className="product-status">
                  <span className={`status-badge ${item.status}`}>
                    {item.status}
                  </span>
                </p>
                <button 
                  className="add-to-cart-btn"
                  onClick={() => addToCart(item.id)}
                >
                  ➕ Add to Cart
                </button>
              </div>
            </div>
          ))
        )}
      </div>

      {/* Cart Sidebar */}
      {showCart && (
        <div className="cart-overlay" onClick={() => setShowCart(false)}>
          <div className="cart-sidebar" onClick={(e) => e.stopPropagation()}>
            <div className="cart-header">
              <h3>Your Cart</h3>
              <button 
                className="close-cart"
                onClick={() => setShowCart(false)}
              >
                ✖️
              </button>
            </div>
            
            <div className="cart-content">
              {cartItems.length === 0 ? (
                <div className="empty-cart">
                  <p>🛒 Your cart is empty</p>
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
                            e.target.onerror = null;
                          }}
                        />
                        <div className="cart-item-details">
                          <h4>{cartItem.item.name}</h4>
                          <p>Status: {cartItem.item.status}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                  
                  <button 
                    className="checkout-cart-btn"
                    onClick={checkout}
                  >
                    💳 Checkout ({cartItems.length} items)
                  </button>
                </>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default SimpleModernItemsList;
