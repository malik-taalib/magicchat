// MongoDB initialization script for MagicChat
// This script creates all necessary indexes for optimal performance

db = db.getSiblingDB('magicchat');

print('Creating indexes for MagicChat database...');

// ===================================
// USERS COLLECTION
// ===================================
print('Creating users indexes...');
db.users.createIndex({ email: 1 }, { unique: true });
db.users.createIndex({ username: 1 }, { unique: true });
db.users.createIndex({ username: 'text', display_name: 'text' });
db.users.createIndex({ follower_count: -1 });
db.users.createIndex({ created_at: -1 });

// ===================================
// VIDEOS COLLECTION
// ===================================
print('Creating videos indexes...');
db.videos.createIndex({ user_id: 1, created_at: -1 });
db.videos.createIndex({ processing_status: 1 });
db.videos.createIndex({ created_at: -1 });
db.videos.createIndex({ hashtags: 1 });
db.videos.createIndex({ title: 'text', description: 'text' });
// Compound index for feed queries
db.videos.createIndex({ processing_status: 1, created_at: -1 });
// Engagement metrics for For You algorithm
db.videos.createIndex({
  processing_status: 1,
  view_count: -1,
  like_count: -1,
  created_at: -1
});

// ===================================
// FOLLOWS COLLECTION
// ===================================
print('Creating follows indexes...');
// Unique constraint: one follow relationship per pair
db.follows.createIndex({ follower_id: 1, following_id: 1 }, { unique: true });
// For getting followers list
db.follows.createIndex({ following_id: 1, created_at: -1 });
// For getting following list
db.follows.createIndex({ follower_id: 1, created_at: -1 });

// ===================================
// LIKES COLLECTION
// ===================================
print('Creating likes indexes...');
db.likes.createIndex({ user_id: 1, video_id: 1 }, { unique: true });
db.likes.createIndex({ video_id: 1, created_at: -1 });
db.likes.createIndex({ user_id: 1, created_at: -1 });

// ===================================
// COMMENTS COLLECTION
// ===================================
print('Creating comments indexes...');
db.comments.createIndex({ video_id: 1, created_at: -1 });
db.comments.createIndex({ user_id: 1, created_at: -1 });
db.comments.createIndex({ parent_id: 1, created_at: 1 });
// For nested comment queries
db.comments.createIndex({ video_id: 1, parent_id: 1, created_at: -1 });

// ===================================
// SHARES COLLECTION
// ===================================
print('Creating shares indexes...');
db.shares.createIndex({ user_id: 1, video_id: 1, created_at: -1 });
db.shares.createIndex({ video_id: 1, created_at: -1 });

// ===================================
// HASHTAGS COLLECTION
// ===================================
print('Creating hashtags indexes...');
db.hashtags.createIndex({ tag: 1 }, { unique: true });
db.hashtags.createIndex({ video_count: -1 });
db.hashtags.createIndex({ trending_score: -1, video_count: -1 });
db.hashtags.createIndex({ last_used: -1 });
db.hashtags.createIndex({ tag: 'text' });

// ===================================
// NOTIFICATIONS COLLECTION
// ===================================
print('Creating notifications indexes...');
db.notifications.createIndex({ user_id: 1, created_at: -1 });
db.notifications.createIndex({ user_id: 1, read: 1 });
// For duplicate detection
db.notifications.createIndex({
  user_id: 1,
  actor_id: 1,
  type: 1,
  video_id: 1,
  created_at: -1
});

// ===================================
// USER_INTERACTIONS COLLECTION (for feed algorithm)
// ===================================
print('Creating user_interactions indexes...');
db.user_interactions.createIndex({ user_id: 1, video_id: 1 }, { unique: true });
db.user_interactions.createIndex({ user_id: 1, updated_at: -1 });
db.user_interactions.createIndex({ video_id: 1, watch_time: -1 });

print('✓ All indexes created successfully!');

// Create default admin user (optional)
print('Creating default test user...');
db.users.insertOne({
  username: 'testuser',
  email: 'test@magicchat.com',
  password_hash: '$2a$10$YourHashedPasswordHere',
  display_name: 'Test User',
  bio: 'Welcome to MagicChat!',
  avatar_url: '',
  follower_count: 0,
  following_count: 0,
  video_count: 0,
  total_likes: 0,
  created_at: new Date(),
  updated_at: new Date()
});

print('✓ Database initialization complete!');
