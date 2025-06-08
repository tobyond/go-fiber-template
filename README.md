# Fresh - Go Fiber App Template

A modern, production-ready Go web application template built with Fiber, GORM, PostgreSQL, and Tailwind CSS.

## ✨ Features

- **🚀 Go Fiber** - Fast, Express-inspired web framework
- **🗄️ GORM + PostgreSQL** - Type-safe database operations
- **🎨 Tailwind CSS** - Modern, utility-first CSS framework
- **⚡ esbuild** - Lightning-fast JavaScript bundling
- **🧪 Comprehensive Testing** - Unit, integration, and HTTP tests
- **🔒 Authentication** - User registration, login, sessions
- **📝 Template System** - Clean HTML templates with layouts
- **⚙️ Configuration** - Rails-style database.yml with environment variables
- **🛠️ Development Tools** - Hot reload, asset watching, testing

## 🏗️ Project Structure

```
fresh/
├── app/
│   ├── controllers/     # HTTP request handlers
│   ├── models/          # Database models and repositories
│   ├── services/        # Business logic layer
│   └── middleware/      # Custom middleware
├── config/              # Configuration files
│   ├── database.yml     # Database configuration
│   └── database.go      # Database initialization
├── routes/              # Route definitions
├── tests/               # Test files
├── web/
│   ├── assets/          # Source assets (CSS, JS)
│   ├── static/          # Built assets (generated)
│   └── templates/       # HTML templates
├── build.js             # Frontend build script
├── package.json         # Node.js dependencies
├── go.mod              # Go dependencies
└── Makefile            # Development commands
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+**
- **Node.js 18+**
- **PostgreSQL**
- **Make** (optional, but recommended)

### 1. Clone and Setup

```bash
# Clone the template
git clone <your-repo-url> myapp
cd myapp

# Run full setup (creates database, installs dependencies, builds assets)
make setup
```

### 2. Configure Database

Edit `config/database.yml` or set environment variables:

```bash
# Option 1: Edit config/database.yml directly
# Option 2: Set environment variables
export DB_NAME=myapp_development
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export DB_HOST=localhost
```

### 3. Start Development

```bash
# Start with hot reload and asset watching
make dev

# Or start normally
make run
```

Visit `http://localhost:3000` 🎉

## 🔧 Customizing for Your App

### Change App Name

1. **Update go.mod:**
   ```go
   module myapp  // Change from 'fresh'
   ```

2. **Update imports in all Go files:**
   ```bash
   # Option 1: Using find with proper escaping
   find . -name "*.go" -type f -exec sed -i 's|fresh/|myapp/|g' {} \;
   
   # Option 2: Using grep + xargs (more reliable)
   grep -r -l "fresh/" --include="*.go" . | xargs sed -i 's|fresh/|myapp/|g'
   
   # Option 3: If you have ripgrep installed (fastest)
   rg -l "fresh/" --type go | xargs sed -i 's|fresh/|myapp/|g'
   ```

3. **Update package.json:**
   ```json
   {
     "name": "myapp",
     "description": "My awesome Go Fiber app"
   }
   ```

4. **Update templates:**
   - Edit `web/templates/layout.html` to change app name in navbar
   - Update page titles and branding

### Database Configuration

**Development (config/database.yml):**
```yaml
development:
  adapter: postgres
  host: ${DB_HOST:localhost}
  port: ${DB_PORT:5432}
  database: ${DB_NAME:myapp_development}
  username: ${DB_USER:postgres}
  password: ${DB_PASSWORD:}
  sslmode: ${DB_SSLMODE:disable}
```

**Environment Variables (.env):**
```bash
ENV=development
DB_NAME=myapp_development
DB_USER=postgres
DB_PASSWORD=secret
PORT=3000
```

**Production:**
Set required environment variables:
```bash
ENV=production
DB_HOST=your-db-host
DB_NAME=myapp_production
DB_USER=your-user
DB_PASSWORD=your-secure-password
DB_SSLMODE=require
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Run with coverage report
make test-coverage

# Run specific test
make test-run TEST=TestUserRepository_Create
```

**Test Structure:**
- `tests/user_test.go` - Model layer tests
- `tests/auth_service_test.go` - Business logic tests  
- `tests/routes_test.go` - HTTP endpoint tests
- `tests/helpers.go` - Test utilities

## 🎨 Frontend Development

**Development:**
```bash
# Watch and rebuild assets automatically
npm run dev

# Build once
npm run build
```

**CSS (Tailwind):**
- Source: `web/assets/css/input.css`
- Output: `web/static/css/styles.css`
- Uses Tailwind CLI for processing

**JavaScript:**
- Source: `web/assets/js/app.js`
- Output: `web/static/js/app.js`
- Uses esbuild for bundling

**Custom Styles:**
Add custom components in `web/assets/css/input.css`:
```css
@layer components {
  .btn-custom {
    @apply bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded;
  }
}
```

## 📝 Adding New Features

### 1. Add a Model

```go
// app/models/post.go
type Post struct {
    ID     uint   `gorm:"primarykey"`
    Title  string `gorm:"not null"`
    Content string
    UserID uint
    User   User
}
```

### 2. Add a Service

```go
// app/services/post_service.go
type PostService struct {
    postRepo *models.PostRepository
}

func (s *PostService) CreatePost(title, content string, userID uint) (*models.Post, error) {
    // Business logic here
}
```

### 3. Add a Controller

```go
// app/controllers/post_controller.go
type PostController struct {
    postService *services.PostService
    templateService services.TemplateRenderer
}

func (pc *PostController) Create(c *fiber.Ctx) error {
    // Handle HTTP request
}
```

### 4. Add Routes

```go
// routes/routes.go
protected.Get("/posts", postController.Index)
protected.Post("/posts", postController.Create)
```

### 5. Add Templates

```html
<!-- web/templates/posts.html -->
{{template "layout" .}}
{{define "content"}}
<h1>Posts</h1>
<!-- Your content here -->
{{end}}
```

### 6. Write Tests

```go
// tests/post_test.go
func TestPostService_Create(t *testing.T) {
    // Test your new feature
}
```

## 🛠️ Available Commands

```bash
# Development
make dev          # Start with hot reload + asset watching
make run          # Start normally
make test         # Run tests
make assets-dev   # Watch and rebuild assets

# Database (uses DB_NAME, DB_USER, DB_HOST env vars or defaults)
make db-create    # Create database
make db-drop      # Drop database  
make db-reset     # Drop and recreate database

# Override database settings:
DB_NAME=myapp_dev DB_USER=myuser make db-create

# Production
make build        # Build optimized binary + assets
make assets-prod  # Build production assets

# Utilities
make clean        # Clean build artifacts
make setup        # Full project setup
```

## 🚀 Deployment

### Build for Production

```bash
# Build everything
make build

# This creates:
# - ./fresh (binary)
# - web/static/ (optimized assets)
```

### Environment Variables

Required for production:
```bash
ENV=production
PORT=8080
DB_HOST=your-db-host
DB_NAME=your-app-prod
DB_USER=your-user  
DB_PASSWORD=your-password
DB_SSLMODE=require
```

### Docker Example

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o fresh *.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/fresh .
COPY --from=builder /app/web ./web
COPY --from=builder /app/config ./config
CMD ["./fresh"]
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for your changes
4. Run `make test` to ensure all tests pass
5. Submit a pull request

## 📄 License

MIT License - feel free to use this template for your projects!

---

**Happy coding! 🎉**

For issues or questions, please open an issue on GitHub.
