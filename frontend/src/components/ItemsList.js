import React, { useState, useEffect } from 'react';
import { useAuth } from '../AuthContext';
import { itemsAPI, cartAPI, ordersAPI } from '../api';
import { toast } from 'react-toastify';
import './ItemsList.css';

const ItemsList = () => {
  const [items, setItems] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const { user, logout } = useAuth();

  useEffect(() => {
    fetchItems();
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

  const addToCart = async (itemId) => {
    try {
      await cartAPI.addItem(itemId);
      toast.success('Item added to cart!');
    } catch (error) {
      if (error.response?.status === 409) {
        toast.warning('Item already in cart');
      } else {
        toast.error('Failed to add item to cart');
      }
    }
  };

  const showCart = async () => {
    try {
      const response = await cartAPI.getUserCart();
      const cart = response.data;
      
      if (cart.cart_items && cart.cart_items.length > 0) {
        const cartDetails = cart.cart_items.map(item => 
          `Cart ID: ${item.cart_id}, Item: ${item.item.name} (ID: ${item.item_id})`
        ).join('\n');
        window.alert(`Cart Items:\n${cartDetails}`);
      } else {
        window.alert('Your cart is empty');
      }
    } catch (error) {
      window.alert('Failed to fetch cart details');
    }
  };

  const showOrderHistory = async () => {
    try {
      const response = await ordersAPI.getUserOrders();
      const orders = response.data;
      
      if (orders && orders.length > 0) {
        const orderIds = orders.map(order => `Order ID: ${order.id}`).join('\n');
        window.alert(`Order History:\n${orderIds}`);
      } else {
        window.alert('No orders found');
      }
    } catch (error) {
      window.alert('Failed to fetch order history');
    }
  };

  const checkout = async () => {
    try {
      await ordersAPI.create();
      toast.success('Order successful!');
    } catch (error) {
      if (error.response?.status === 400) {
        toast.error('Cart is empty');
      } else {
        toast.error('Failed to create order');
      }
    }
  };

  if (isLoading) {
    return <div className="loading">Loading items...</div>;
  }

  return (
    <div className="items-container">
      <header className="items-header">
        <h1>Welcome, {user?.username}!</h1>
        <div className="header-buttons">
          <button onClick={checkout} className="checkout-btn">
            Checkout
          </button>
          <button onClick={showCart} className="cart-btn">
            Cart
          </button>
          <button onClick={showOrderHistory} className="history-btn">
            Order History
          </button>
          <button onClick={logout} className="logout-btn">
            Logout
          </button>
        </div>
      </header>

      <div className="items-grid">
        {items.length === 0 ? (
          <p>No items available</p>
        ) : (
          items.map((item) => (
            <div key={item.id} className="item-card">
              <div className="item-image">
                <img 
                  src={`http://localhost:8080${item.image}`} 
                  alt={item.name}
                  onError={(e) => {
                    e.target.src = 'https://via.placeholder.com/300x300/e9ecef/6c757d?text=' + encodeURIComponent(item.name);
                    e.target.onerror = null;
                  }}
                />
              </div>
              <div className="item-details">
                <h3>{item.name}</h3>
                <p>Status: {item.status}</p>
                <button 
                  onClick={() => addToCart(item.id)}
                  className="add-to-cart-btn"
                >
                  Add to Cart
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default ItemsList;
