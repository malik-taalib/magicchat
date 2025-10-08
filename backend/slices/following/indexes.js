// MongoDB indexes for the following slice
// Run this in MongoDB shell or using mongosh:
// mongosh magicchat < indexes.js

// Switch to the magicchat database
use magicchat;

// Create indexes for the follows collection
db.follows.createIndex(
  { "follower_id": 1, "following_id": 1 },
  {
    unique: true,
    name: "unique_follow_relationship"
  }
);

db.follows.createIndex(
  { "following_id": 1, "_id": -1 },
  {
    name: "following_id_pagination"
  }
);

db.follows.createIndex(
  { "follower_id": 1, "_id": -1 },
  {
    name: "follower_id_pagination"
  }
);

// Create indexes for the users collection (if not already present)
db.users.createIndex(
  { "username": 1 },
  {
    unique: true,
    name: "unique_username"
  }
);

db.users.createIndex(
  { "email": 1 },
  {
    unique: true,
    name: "unique_email"
  }
);

print("✅ All indexes created successfully!");
print("\nIndexes on 'follows' collection:");
printjson(db.follows.getIndexes());

print("\nIndexes on 'users' collection:");
printjson(db.users.getIndexes());
