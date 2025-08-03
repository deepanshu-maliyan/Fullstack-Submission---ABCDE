# Product Images Setup

## Quick Setup for Testing

Since you don't have product images yet, here are some options:

### Option 1: Download Free Stock Images
Visit these sites for free product images:
- **Unsplash**: https://unsplash.com/s/photos/laptop
- **Pexels**: https://www.pexels.com/search/electronics/
- **Pixabay**: https://pixabay.com/images/search/computer/

### Option 2: Use Placeholder Services (Temporary)
For quick testing, you can temporarily modify the image URLs in the database to use:
- `https://via.placeholder.com/400x400/007bff/ffffff?text=Laptop`
- `https://via.placeholder.com/400x400/28a745/ffffff?text=Phone`

### Option 3: Create Simple Placeholders
I can create simple colored placeholder images for testing.

## Expected Image Files

Place these images in `/assets/products/`:
1. `laptop.jpg` - Laptop image
2. `smartphone.jpg` - Smartphone image  
3. `headphones.jpg` - Headphones image
4. `keyboard.jpg` - Keyboard image
5. `mouse.jpg` - Mouse image
6. `monitor.jpg` - Monitor image
7. `tablet.jpg` - Tablet image
8. `webcam.jpg` - Webcam image

## Image Requirements
- **Format**: JPG, PNG, or WebP
- **Size**: 400x400px recommended (square)
- **File size**: Under 500KB each
- **Naming**: Lowercase, matching product names

## Backend Configuration
✅ Backend is already configured to serve static files
✅ Database includes image URLs
✅ Frontend displays product images

## Testing Without Images
If images are missing, the frontend will show a broken image icon, but the application will still work normally.

## Next Steps
1. Add your product images to `/assets/products/`
2. Restart the backend server
3. Refresh the frontend to see images
