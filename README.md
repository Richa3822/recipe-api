## ⚙️ Getting Started

### Prerequisites
- Docker
- Docker Compose

### Run with Docker

```bash
# Clone the repo
git clone https://github.com/Richa3822/recipe-api.git
cd recipe-api

# Start everything
docker-compose up --build
```

API is now running at `http://localhost:8080`

### Run locally (without Docker)

```bash
# Make sure PostgreSQL is running, then:
go mod download
cd cmd
go run main.go
```

## 🔑 Environment Variables

Create a `.env` file in the root:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=recipedb
JWT_SECRET=your-super-secret-key
```

## 📡 API Endpoints

### Auth

| Method | Endpoint | Description | Auth |
|---|---|---|---|
| POST | `/api/auth/register` | Register new user | ❌ |
| POST | `/api/auth/login` | Login, returns JWT token | ❌ |

### Recipes

| Method | Endpoint | Description | Auth |
|---|---|---|---|
| POST | `/api/recipes` | Create a recipe | ✅ |
| GET | `/api/recipes` | Get all your recipes | ✅ |
| GET | `/api/recipes?search=pasta` | Search by title or ingredient | ✅ |
| GET | `/api/recipes?page=2&limit=5` | Paginated results | ✅ |
| GET | `/api/recipes/:id` | Get recipe by ID | ✅ |
| PUT | `/api/recipes/:id` | Update your recipe | ✅ |
| DELETE | `/api/recipes/:id` | Delete your recipe | ✅ |

### Auth Header
All protected routes require: