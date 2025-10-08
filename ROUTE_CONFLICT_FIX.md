# Route Conflict Fix - Magic Chat

## Problem

You had **duplicate routes** mounted on `/videos`:

```go
// ❌ CONFLICT - Both trying to use /videos
r.Mount("/videos", videoupload.Routes(...))      // Line 90
r.Route("/videos", func(r chi.Router) {          // Line 96
    r.Mount("/", engagement.Routes(db))
})
```

This caused the panic:
```
panic: chi: attempting to Mount() a handler on an existing path, '/videos'
```

---

## Solution

Changed engagement routes from `/videos` to `/engage` to avoid conflict:

```go
// ✅ FIXED - Separate paths
r.Mount("/videos", videoupload.Routes(...))      // Video upload
r.Mount("/engage", engagement.Routes(db))        // Engagement (likes, comments)
```

---

## Updated API Endpoints

### Video Upload Routes (`/api/videos`)
```
POST   /api/videos/upload          - Upload video
GET    /api/videos/:id/status      - Get processing status
POST   /api/videos/process         - Processing webhook
```

### Video Feed Routes (`/api/feed`)
```
GET    /api/feed/for-you           - For You feed
GET    /api/feed/following         - Following feed
GET    /api/feed/:id               - Get single video
```

### Engagement Routes (`/api/engage`) ⚠️ **CHANGED**
```
POST   /api/engage/:id/like        - Like video (was /api/videos/:id/like)
DELETE /api/engage/:id/like        - Unlike video (was /api/videos/:id/like)
POST   /api/engage/:id/comments    - Create comment (was /api/videos/:id/comments)
GET    /api/engage/:id/comments    - Get comments (was /api/videos/:id/comments)
GET    /api/engage/comments/:id/replies  - Get replies (was /api/videos/comments/:id/replies)
POST   /api/engage/:id/share       - Share video (was /api/videos/:id/share)
```

---

## Breaking Change

⚠️ **Frontend API client needs to be updated!**

Change all engagement endpoints from `/videos/` to `/engage/`:

### Before (Old):
```typescript
// ❌ OLD
apiClient.likeVideo(videoId)
// Calls: POST /api/videos/:id/like
```

### After (New):
```typescript
// ✅ NEW
apiClient.likeVideo(videoId)
// Should call: POST /api/engage/:id/like
```

---

## Alternative Solutions (Not Used)

### Option 1: Nest engagement under videos
```go
r.Route("/videos", func(r chi.Router) {
    r.Post("/upload", uploadHandler)
    r.Get("/{id}/status", statusHandler)
    r.Post("/{id}/like", likeHandler)      // Nested
    r.Post("/{id}/comments", commentHandler)  // Nested
})
```
**Pros**: Clean URL structure
**Cons**: Requires merging all handlers into one route file

### Option 2: Use different prefixes
```go
r.Mount("/videos", videoupload.Routes())      // /videos/upload
r.Mount("/video-engagement", engagement.Routes())  // /video-engagement/:id/like
```
**Pros**: Clear separation
**Cons**: Longer URLs

### Option 3: Separate by resource (CHOSEN ✅)
```go
r.Mount("/videos", videoupload.Routes())      // Video management
r.Mount("/engage", engagement.Routes())       // User engagement
```
**Pros**: Logical separation of concerns
**Cons**: Engagement not under /videos

---

## To Update Frontend

Update `frontend/lib/api/client.ts`:

```typescript
// Change engagement endpoints
async likeVideo(videoId: string) {
  return this.request(`/engage/${videoId}/like`, { method: 'POST' });
}

async unlikeVideo(videoId: string) {
  return this.request(`/engage/${videoId}/like`, { method: 'DELETE' });
}

async getComments(videoId: string, limit = 20, offset = 0) {
  return this.request<Comment[]>(`/engage/${videoId}/comments?limit=${limit}&offset=${offset}`);
}

async createComment(videoId: string, text: string, parentId?: string) {
  return this.request<Comment>(`/engage/${videoId}/comments`, {
    method: 'POST',
    body: JSON.stringify({ text, parent_id: parentId }),
  });
}

async shareVideo(videoId: string) {
  return this.request(`/engage/${videoId}/share`, { method: 'POST' });
}
```

---

## Status

✅ **Backend Fixed** - No more route conflicts
✅ **Build Successful** - Server compiles without errors
⚠️ **Frontend Needs Update** - API client endpoints changed

---

## Quick Test

Start the server and check registered routes:
```bash
make run

# Server will print all registered routes including:
# POST   /api/videos/upload
# POST   /api/engage/{id}/like
# POST   /api/engage/{id}/comments
```
