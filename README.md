# Magic Chat 🎬

A TikTok-style short-form video social platform built with **Vertical Slice Architecture**, combining Go backend services, Next.js frontend, and MongoDB for data persistence.

## 🏗️ Architecture

Magic Chat follows the **Vertical Slice Architecture** pattern where each feature is organized as a complete vertical slice containing all layers (API, business logic, data access) rather than organizing by technical layer.

### Benefits
- ✅ Independent feature deployment
- ✅ Reduced coupling between features
- ✅ Easier testing and maintenance
- ✅ Team autonomy per feature

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.21+
- **Router**: Chi v5
- **Database**: MongoDB 6.0+
- **Caching**: Redis 7
- **Storage**: S3-compatible (AWS S3, MinIO)
- **WebSocket**: Gorilla WebSocket

### Frontend
- **Framework**: Next.js 14+ (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Video Player**: HTML5 native

## 📦 Features

### ✅ Implemented Vertical Slices

1. **Authentication Slice** (`/slices/auth`)
   - User registration & login
   - JWT-based authentication
   - Password hashing with bcrypt
   - Auth middleware for protected routes

2. **Video Upload Slice** (`/slices/video-upload`)
   - Multipart file upload
   - Video storage (S3/MinIO)
   - Processing status tracking
   - File validation (size, format, duration)

3. **Video Feed Slice** (`/slices/video-feed`)
   - For You feed (algorithmic)
   - Following feed
   - Engagement-based ranking
   - Cursor-based pagination

4. **Engagement Slice** (`/slices/engagement`)
   - Like/unlike videos
   - Comments with nested replies
   - Share tracking
   - Real-time count updates

5. **Following Slice** (`/slices/following`)
   - Follow/unfollow users
   - Followers & following lists
   - Follow status checks
   - Follower count tracking

6. **Search & Discovery Slice** (`/slices/search`)
   - Search users, videos, hashtags
   - Trending hashtags
   - Videos by hashtag
   - Relevance-based ranking

7. **Notifications Slice** (`/slices/notifications`)
   - Real-time WebSocket notifications
   - Notification types: like, comment, follow, mention
   - Read/unread tracking
   - Notification history with pagination

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ (for frontend)

### Local Development with Docker

```bash
# Clone the repository
git clone <repository-url>
cd magicchat

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

Services will be available at:
- **Backend API**: http://localhost:8080
- **Frontend**: http://localhost:3000
- **MinIO Console**: http://localhost:9001 (admin/minioadmin)
- **MongoDB**: localhost:27017

### Local Development without Docker

#### Backend

```bash
cd backend

# Copy environment file
cp .env.example .env

# Install dependencies
make install

# Run the server
make run

# Or with hot reload (requires air)
make dev
```

#### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api
```

### Authentication Endpoints

```http
POST   /api/auth/register    # Register new user
POST   /api/auth/login       # Login user
POST   /api/auth/logout      # Logout (client-side)
GET    /api/auth/me          # Get current user (protected)
```

### Video Endpoints

```http
POST   /api/videos/upload              # Upload video (protected)
GET    /api/videos/:id/status          # Get processing status
GET    /api/feed/for-you               # For You feed (protected)
GET    /api/feed/following             # Following feed (protected)
GET    /api/feed/:id                   # Get single video
```

### Engagement Endpoints

```http
POST   /api/videos/:id/like            # Like video (protected)
DELETE /api/videos/:id/like            # Unlike video (protected)
POST   /api/videos/:id/comments        # Create comment (protected)
GET    /api/videos/:id/comments        # Get comments
GET    /api/videos/comments/:id/replies # Get comment replies
POST   /api/videos/:id/share           # Share video (protected)
```

### Social Endpoints

```http
POST   /api/users/:id/follow           # Follow user (protected)
DELETE /api/users/:id/follow           # Unfollow user (protected)
GET    /api/users/:id/followers        # Get followers
GET    /api/users/:id/following        # Get following
```

### Search Endpoints

```http
GET    /api/search?q=query&type=users|videos|hashtags    # Search
GET    /api/trending/hashtags                             # Trending hashtags
GET    /api/hashtags/:tag/videos                          # Videos by hashtag
```

### Notification Endpoints

```http
GET    /api/notifications              # Get notifications (protected)
PUT    /api/notifications/:id/read     # Mark as read (protected)
PUT    /api/notifications/read-all     # Mark all as read (protected)
WS     /api/notifications/stream       # WebSocket connection (protected)
```

## 🗄️ Database Schema

### Users Collection
```javascript
{
  _id: ObjectId,
  username: String (unique),
  email: String (unique),
  password_hash: String,
  display_name: String,
  bio: String,
  avatar_url: String,
  follower_count: Number,
  following_count: Number,
  video_count: Number,
  total_likes: Number,
  created_at: Date,
  updated_at: Date
}
```

### Videos Collection
```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  title: String,
  description: String,
  video_url: String,
  thumbnail_url: String,
  duration: Number,
  hashtags: [String],
  view_count: Number,
  like_count: Number,
  comment_count: Number,
  share_count: Number,
  processing_status: Enum,
  created_at: Date,
  updated_at: Date
}
```

See [backend/migrations/init-indexes.js](backend/migrations/init-indexes.js) for complete schema and indexes.

## 🏗️ Project Structure

```
magicchat/
├── backend/
│   ├── cmd/
│   │   └── server/          # Main application entry
│   ├── slices/              # Vertical slices
│   │   ├── auth/
│   │   ├── video-upload/
│   │   ├── video-feed/
│   │   ├── engagement/
│   │   ├── following/
│   │   ├── search/
│   │   └── notifications/
│   ├── pkg/                 # Shared packages
│   │   ├── config/
│   │   ├── database/
│   │   ├── cache/
│   │   └── storage/
│   ├── migrations/          # Database migrations
│   ├── go.mod
│   └── Dockerfile
├── frontend/
│   ├── app/                 # Next.js app directory
│   ├── components/
│   ├── lib/
│   └── package.json
└── docker-compose.yml
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run linter
make lint
```

## 🔧 Configuration

### Environment Variables

Backend configuration is done via environment variables. See [backend/.env.example](backend/.env.example) for all available options.

Key variables:
- `MONGODB_URI` - MongoDB connection string
- `REDIS_URL` - Redis connection string
- `JWT_SECRET` - Secret key for JWT tokens
- `STORAGE_PROVIDER` - `minio` or `s3`
- `MINIO_ENDPOINT` - MinIO server endpoint

## 🚢 Deployment

### Production Considerations

1. **Database**: Use MongoDB Atlas or managed MongoDB cluster
2. **Storage**: AWS S3 with CloudFront CDN
3. **Cache**: Redis cluster or managed Redis service
4. **Backend**: Kubernetes deployment with horizontal scaling
5. **Frontend**: Vercel or self-hosted with Docker
6. **Security**:
   - Change default passwords
   - Use strong JWT secrets
   - Enable rate limiting
   - Set up proper CORS

## 📝 Development Commands

```bash
# Backend
make help          # Show available commands
make build         # Build the application
make run           # Run the server
make dev           # Run with hot reload
make test          # Run tests
make docker-up     # Start Docker services
make docker-down   # Stop Docker services

# Frontend
npm run dev        # Development server
npm run build      # Production build
npm run start      # Start production server
npm run lint       # Run ESLint
```

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. Create a new vertical slice for new features
2. Follow the established patterns
3. Write tests for your code
4. Update documentation
5. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details

## 👥 Authors

- Your Name - Initial work

## 🙏 Acknowledgments

- Built following Vertical Slice Architecture principles
- Inspired by TikTok's short-form video platform
- Uses best practices from Go and React communities
