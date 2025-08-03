# Admin User Credentials 🔑

## Default Admin Account

**Username**: `admin`  
**Password**: `Admin@123`

## How It Works

✅ **Auto-Created**: The admin user is automatically created when the backend starts  
✅ **Pre-configured**: No need to register - the admin account is ready to use  
✅ **Secure**: Password is properly hashed using bcrypt  
✅ **Persistent**: Admin user persists in the in-memory database during runtime  

## Login Instructions

1. **Start the backend**: Run `start-backend.bat`
2. **Start the frontend**: Run `start-frontend.bat`  
3. **Login with admin credentials**:
   - Username: `admin`
   - Password: `Admin@123`

## Admin Features

Once logged in as admin, you have access to:
- ✅ All items in the store
- ✅ Shopping cart functionality
- ✅ Order processing
- ✅ Full API access with JWT token

## Backend Logs

When the backend starts, you'll see:
```
Created admin user (username: admin, password: Admin@123)
```

If the admin user already exists:
```
Admin user already exists
```

## Security Note

🔒 The password is securely hashed using bcrypt before storage  
🔒 JWT tokens are generated for authentication  
🔒 All protected endpoints require valid Bearer token  

## Testing the Admin Account

You can test the admin login by:
1. Using the frontend login form
2. Making a direct API call:
   ```bash
   curl -X POST http://localhost:8080/users/login \
   -H "Content-Type: application/json" \
   -d '{"username":"admin","password":"Admin@123"}'
   ```

**The admin account is ready to use immediately when you start the backend!** ✅
