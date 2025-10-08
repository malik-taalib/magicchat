# Following Slice Integration Guide

## Overview
This slice implements user following functionality with a complete vertical slice architecture including models, repository, service, handler, and routes.

## Architecture

```
following/
├── models.go       # Data structures (Follow, UserProfile, responses)
├── repository.go   # Database operations
├── service.go      # Business logic
├── handler.go      # HTTP handlers
└── routes.go       # Route definitions
```

## Features

1. **Follow/Unfollow Users**
   - Create and remove follow relationships
   - Automatic follower/following count updates
   - Prevents self-following

2. **Paginated Lists**
   - Get followers of a user
   - Get users that a user follows
   - Cursor-based pagination
   - Configurable page size (default 20, max 100)

3. **Follow Status Check**
   - Check if one user follows another
   - Useful for UI state management

## API Endpoints

All endpoints require authentication via Bearer token.

### Follow a User
```
POST /users/{id}/follow
Authorization: Bearer <token>

Response:
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

### Unfollow a User
```
DELETE /users/{id}/follow
Authorization: Bearer <token>

Response:
{
  "success": true,
  "data": {
    "success": true,
    "is_following": false,
    "follower_count": 149,
    "following_count": 74,
    "message": "Successfully unfollowed user"
  }
}
```

### Get Followers
```
GET /users/{id}/followers?cursor=<cursor>&limit=20
Authorization: Bearer <token>

Response:
{
  "success": true,
  "data": {
    "users": [
      {
        "id": "507f1f77bcf86cd799439011",
        "username": "john_doe",
        "display_name": "John Doe",
        "avatar_url": "https://...",
        "follower_count": 100,
        "following_count": 50
      }
    ],
    "next_cursor": "507f1f77bcf86cd799439012",
    "has_more": true,
    "total": 150
  }
}
```

### Get Following
```
GET /users/{id}/following?cursor=<cursor>&limit=20
Authorization: Bearer <token>

Response:
{
  "success": true,
  "data": {
    "users": [
      {
        "id": "507f1f77bcf86cd799439011",
        "username": "jane_smith",
        "display_name": "Jane Smith",
        "avatar_url": "https://...",
        "follower_count": 200,
        "following_count": 100
      }
    ],
    "next_cursor": "507f1f77bcf86cd799439012",
    "has_more": true,
    "total": 75
  }
}
```

### Check Follow Status
```
GET /users/{id}/following/check
Authorization: Bearer <token>

Response:
{
  "success": true,
  "data": {
    "is_following": true,
    "user_id": "507f1f77bcf86cd799439011"
  }
}
```

## Integration Example

### In your main application or gateway:

```go
package main

import (
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "go.mongodb.org/mongo-driver/mongo"

    "magicchat/backend/slices/auth"
    "magicchat/backend/slices/following"
)

func main() {
    // Connect to MongoDB
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(context.Background())

    db := client.Database("magicchat")

    // Create main router
    r := chi.NewRouter()

    // Mount auth routes
    r.Mount("/api/auth", auth.Routes(db))

    // Mount following routes
    r.Mount("/api/users", following.Routes(db))

    log.Println("Server running on :8080")
    http.ListenAndServe(":8080", r)
}
```

## Database Collections

### follows collection
```json
{
  "_id": ObjectId("..."),
  "follower_id": ObjectId("..."),
  "following_id": ObjectId("..."),
  "created_at": ISODate("2024-01-01T00:00:00Z")
}
```

**Indexes needed:**
```javascript
db.follows.createIndex({ "follower_id": 1, "following_id": 1 }, { unique: true })
db.follows.createIndex({ "following_id": 1, "_id": -1 })
db.follows.createIndex({ "follower_id": 1, "_id": -1 })
```

### users collection
The slice expects these fields in the users collection:
- `_id` (ObjectId)
- `username` (string)
- `display_name` (string)
- `avatar_url` (string)
- `follower_count` (int)
- `following_count` (int)

## Dependencies

- `go.mongodb.org/mongo-driver/mongo` - MongoDB driver
- `github.com/go-chi/chi/v5` - HTTP router
- `magicchat/backend/slices/auth` - Authentication middleware

## Error Handling

The slice provides comprehensive error handling for:
- Invalid user IDs
- Self-follow attempts
- Duplicate follow relationships
- Non-existent users
- Unauthorized requests
- Invalid pagination parameters

## Testing Recommendations

1. **Unit Tests**
   - Test repository methods with mock MongoDB
   - Test service business logic
   - Test handler request/response handling

2. **Integration Tests**
   - Test full follow/unfollow flow
   - Test pagination with various cursor values
   - Test concurrent follow operations

3. **Edge Cases**
   - Following non-existent users
   - Double-following
   - Unfollowing when not following
   - Large follower/following lists
   - Invalid ObjectIDs

## Performance Considerations

1. **Pagination**: Uses cursor-based pagination for efficient large dataset handling
2. **Indexes**: Ensure proper indexes on `follower_id` and `following_id` fields
3. **Batch Operations**: User profile fetching is batched to minimize database queries
4. **Count Updates**: Follower/following counts are updated atomically

## Future Enhancements

- Follow request system for private accounts
- Follow suggestions based on mutual connections
- Activity feed for followed users
- Bulk follow operations
- Follow analytics and insights
