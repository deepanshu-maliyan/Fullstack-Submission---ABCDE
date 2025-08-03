# Assets Folder

This folder contains all static assets for the e-commerce application.

## Structure

```
assets/
├── products/          # Product images
│   ├── laptop.jpg     # Example product images
│   ├── smartphone.jpg
│   ├── headphones.jpg
│   ├── keyboard.jpg
│   ├── mouse.jpg
│   ├── monitor.jpg
│   ├── tablet.jpg
│   └── webcam.jpg
└── README.md          # This file
```

## Product Images

### Naming Convention
- Use lowercase names matching the product names
- Use common image formats: `.jpg`, `.png`, `.webp`
- Keep file sizes reasonable (< 1MB for web performance)

### Expected Images
Based on the seeded products in the database:
1. `laptop.jpg` - Laptop product image
2. `smartphone.jpg` - Smartphone product image  
3. `headphones.jpg` - Headphones product image
4. `keyboard.jpg` - Keyboard product image
5. `mouse.jpg` - Mouse product image
6. `monitor.jpg` - Monitor product image
7. `tablet.jpg` - Tablet product image
8. `webcam.jpg` - Webcam product image

### Usage
- Images are served statically by the backend server
- Frontend references images via URL: `/assets/products/{filename}`
- Images are automatically displayed in the product listings

## Adding New Images

1. **Add image file** to `/assets/products/` folder
2. **Use descriptive filename** matching product name
3. **Update product data** if adding new products
4. **Restart backend** to serve new static files

## Image Guidelines

### Recommended Specifications:
- **Size**: 400x400px or 600x600px (square aspect ratio)
- **Format**: JPG for photos, PNG for graphics with transparency
- **Quality**: 80-90% compression for web optimization
- **File Size**: Under 500KB per image

### Example Filenames:
- `laptop.jpg`
- `gaming-laptop.jpg` 
- `wireless-mouse.jpg`
- `mechanical-keyboard.jpg`

## Backend Configuration

The backend serves static files from this folder via:
```go
router.Static("/assets", "./assets")
```

This allows frontend to access images at:
```
http://localhost:8080/assets/products/laptop.jpg
```
