# Following Slice

A complete vertical slice implementation for user following functionality in the MagicChat application.

## Overview

This slice follows the vertical slice architecture pattern, encapsulating all layers (models, repository, service, handler, routes) required for the following feature.

## Structure

```
following/
├── models.go           # Data models and DTOs
├── repository.go       # Database operations layer
├── service.go          # Business logic layer
├── handler.go          # HTTP handlers
├── routes.go           # Route definitions
├── interfaces.go       # Repository interface for testing
├── service_test.go     # Example unit tests
├── indexes.js          # MongoDB index creation script
├── INTEGRATION.md      # Detailed integration guide
└── README.md           # This file
```

## Features

- **Follow/Unfollow**: Users can follow and unfollow other users
- **Follow Lists**: Paginated lists of followers and following
- **Follow Counts**: Automatic updating of follower/following counts
- **Follow Status**: Check if a user is following another user
- **Authentication**: All endpoints protected with JWT middleware

## Quick Start

### 1. Create MongoDB Indexes

```bash
mongosh magicchat < indexes.js
```

### 2. Mount Routes in Your Application

```go
import "magicchat/slices/following"

// In your main.go or router setup
db := client.Database("magicchat")
r.Mount("/api/users", following.Routes(db))
```

### 3. API Endpoints

All endpoints require `Authorization: Bearer <token>` header.

- `POST /api/users/{id}/follow` - Follow a user
- `DELETE /api/users/{id}/follow` - Unfollow a user
- `GET /api/users/{id}/followers` - Get user's followers
- `GET /api/users/{id}/following` - Get users that a user follows
- `GET /api/users/{id}/following/check` - Check follow status

## API Examples

### Follow a User

```bash
curl -X POST http://localhost:8080/api/users/507f1f77bcf86cd799439011/follow \
  -H "Authorization: Bearer <token>"
```

Response:
```json
{
  "success": true,
  "data": {
    "success": true,
    "is_following": true,
    "follower_count": 150,
    "following_count": 75,
    "message": "Successfully followed user"
  }
}
```

### Get Followers (with pagination)

```bash
curl http://localhost:8080/api/users/507f1f77bcf86cd799439011/followers?limit=20 \
  -H "Authorization: Bearer <token>"
```

Response:
```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": "507f1f77bcf86cd799439012",
        "username": "john_doe",
        "display_name": "John Doe",
        "avatar_url": "https://...",
        "follower_count": 100,
        "following_count": 50
      }
    ],
    "next_cursor": "507f1f77bcf86cd799439013",
    "has_more": true,
    "total": 150
  }
}
```

## Testing

Run the included tests:

```bash
cd backend/slices/following
go test -v
```

The service layer uses the `RepositoryInterface` which allows for easy mocking in tests.

## Database Schema

### follows Collection

```json
{
  "_id": ObjectId("..."),
  "follower_id": ObjectId("..."),    // User who is following
  "following_id": ObjectId("..."),   // User being followed
  "created_at": ISODate("...")
}
```

**Required Indexes:**
- `{follower_id: 1, following_id: 1}` (unique)
- `{following_id: 1, _id: -1}` (for followers pagination)
- `{follower_id: 1, _id: -1}` (for following pagination)

### users Collection (partial)

The slice expects these fields:
- `follower_count` (int)
- `following_count` (int)

These counts are automatically maintained by the repository layer.

## Architecture Layers

### 1. Models (`models.go`)
- Data structures
- Request/Response DTOs
- No business logic

### 2. Repository (`repository.go`)
- Database operations
- MongoDB queries
- Data persistence
- Count management

### 3. Service (`service.go`)
- Business logic
- Validation
- Error handling
- Uses `RepositoryInterface` for testability

### 4. Handler (`handler.go`)
- HTTP request handling
- Parameter extraction
- Response formatting
- Authentication context

### 5. Routes (`routes.go`)
- Route definitions
- Middleware application
- Dependency injection

## Dependencies

- `go.mongodb.org/mongo-driver/mongo` - MongoDB driver
- `github.com/go-chi/chi/v5` - HTTP router
- `magicchat/slices/auth` - Authentication middleware

## Error Handling

The slice provides clear error messages for:
- Invalid user IDs
- Self-follow attempts
- Already following/not following
- User not found
- Unauthorized access
- Invalid pagination parameters

## Performance Considerations

1. **Cursor-based Pagination**: Efficient for large datasets
2. **Indexed Queries**: All queries use indexed fields
3. **Batch User Fetching**: Profiles fetched in a single query
4. **Atomic Count Updates**: Follower/following counts updated atomically

## Future Enhancements

- [ ] Follow request system for private accounts
- [ ] Mutual follow detection
- [ ] Follow suggestions
- [ ] Activity notifications
- [ ] Bulk operations
- [ ] Follow analytics

## License

Part of the MagicChat application.

## See Also

- [INTEGRATION.md](./INTEGRATION.md) - Detailed integration guide
- [service_test.go](./service_test.go) - Test examples
