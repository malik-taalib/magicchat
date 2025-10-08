# Quick Start Guide for Magic Chat

This guide will help you get Magic Chat up and running in under 5 minutes.

## Prerequisites

- Docker and Docker Compose installed
- Go 1.21+ (optional, for local development)
- Node.js 18+ (optional, for frontend development)

## Option 1: Docker Compose (Recommended)

### Start Everything with One Command

```bash
# Start all services (MongoDB, Redis, MinIO, Backend, Frontend)
docker-compose up -d

# View logs
docker-compose logs -f backend
```

That's it! Services will be available at:
- **Backend API**: http://localhost:8080
- **Frontend**: http://localhost:3000
- **MinIO Console**: http://localhost:9001
- **API Health Check**: http://localhost:8080/health

### Stop Services

```bash
docker-compose down

# Remove all data volumes (caution: deletes all data)
docker-compose down -v
```

## Option 2: Local Development

### 1. Start Infrastructure Services

```bash
# Start only MongoDB, Redis, and MinIO
docker-compose up -d mongodb redis minio
```

### 2. Run Backend Locally

```bash
cd backend

# Copy environment file
cp .env.example .env

# Install dependencies
go mod download

# Run server
go run cmd/server/main.go

# Or use Make
make run
```

Backend will start on http://localhost:8080

### 3. Run Frontend Locally

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

Frontend will start on http://localhost:3000

## First Steps

### 1. Create a Test User

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "display_name": "Test User"
  }'
```

Response:
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "...",
      "username": "testuser",
      ...
    }
  }
}
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. Get Current User (Protected Route)

```bash
# Replace YOUR_TOKEN with the token from registration/login
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4. Upload a Video

```bash
curl -X POST http://localhost:8080/api/videos/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "video=@/path/to/your/video.mp4" \
  -F "title=My First Video" \
  -F "description=This is a test video"
```

### 5. Get For You Feed

```bash
curl -X GET http://localhost:8080/api/feed/for-you \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Testing the Frontend

1. Open http://localhost:3000 in your browser
2. You should see the Magic Chat interface
3. Sign up or log in with your test account
4. Explore the features!

## Troubleshooting

### Port Already in Use

If you get port conflicts:

```bash
# Change ports in docker-compose.yml
# For example, change "8080:8080" to "8081:8080"
```

### Cannot Connect to Database

```bash
# Check if services are running
docker-compose ps

# Restart services
docker-compose restart mongodb redis minio
```

### Backend Won't Start

```bash
# Check logs
docker-compose logs backend

# Or if running locally
cd backend
go mod tidy
go run cmd/server/main.go
```

### Frontend Build Errors

```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run dev
```

## Environment Variables

### Backend (.env)

Key variables you might want to change:

```bash
# Server
PORT=8080
ENV=development

# MongoDB
MONGODB_URI=mongodb://admin:password@localhost:27017
MONGODB_DATABASE=magicchat

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=change-this-to-a-secure-secret
JWT_EXPIRY=24h

# Storage (MinIO local or AWS S3)
STORAGE_PROVIDER=minio
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
S3_BUCKET=magicchat-videos

# Video Limits
MAX_VIDEO_SIZE_MB=100
MAX_VIDEO_DURATION_SECONDS=180
ALLOWED_VIDEO_FORMATS=mp4,mov,avi,webm

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

### Frontend (.env.local)

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_WS_URL=ws://localhost:8080/api
```

## Next Steps

- Read the [main README](README.md) for detailed documentation
- Check out the [API documentation](#) (coming soon)
- Explore the vertical slices in `backend/slices/`
- Build your first frontend component

## Useful Commands

```bash
# Backend
make help          # Show all available commands
make run           # Run backend server
make test          # Run tests
make docker-up     # Start Docker services
make docker-logs   # View logs

# Frontend
npm run dev        # Development mode
npm run build      # Production build
npm run lint       # Check code quality

# Docker
docker-compose up -d              # Start all services
docker-compose down               # Stop all services
docker-compose logs -f [service]  # View logs
docker-compose restart [service]  # Restart a service
docker-compose exec backend sh    # Access backend container
```

## Accessing Services

### MongoDB

```bash
# Using Docker
docker-compose exec mongodb mongosh -u admin -p password

# Or using local mongosh
mongosh mongodb://admin:password@localhost:27017
```

### Redis

```bash
# Using Docker
docker-compose exec redis redis-cli

# Or using local redis-cli
redis-cli
```

### MinIO

Open http://localhost:9001 in your browser
- Username: `minioadmin`
- Password: `minioadmin`

## Production Deployment

For production deployment, see the deployment section in the [main README](README.md#-deployment).

Key steps:
1. Use managed services (MongoDB Atlas, Redis Cloud, AWS S3)
2. Change all default passwords and secrets
3. Enable HTTPS/TLS
4. Set up proper CORS and rate limiting
5. Use environment-specific configurations
6. Set up monitoring and logging

Happy coding! 🚀
