# Magic Chat Database Seed Data

This directory contains sample data to populate your Magic Chat database for testing and development.

## Sample Data Includes

- **8 Users** - Various content creators (travel, cooking, fitness, comedy, tech, music, art, dance)
- **12 Videos** - Sample videos with real URLs from Google's test video bucket
- **8 Comments** - Sample comments on videos
- **5 Likes** - Sample video likes
- **5 Follows** - Sample user follows
- **5 Hashtags** - Trending hashtags with stats

## Test Credentials

All users have the same password for testing:

| Email | Username | Password |
|-------|----------|----------|
| sarah@magicchat.com | sarah_adventures | password123 |
| mike@magicchat.com | chef_mike | password123 |
| emma@magicchat.com | fitness_emma | password123 |
| jay@magicchat.com | comedy_jay | password123 |
| alex@magicchat.com | tech_alex | password123 |
| lisa@magicchat.com | music_lisa | password123 |
| marcus@magicchat.com | art_marcus | password123 |
| nina@magicchat.com | dance_nina | password123 |

## How to Use

### Prerequisites

Make sure MongoDB is running:

```bash
# If using Docker Compose from project root
docker-compose up -d mongodb

# Or if running MongoDB locally
mongosh # to verify it's running
```

### Option 1: Run with npm (from backend directory)

```bash
cd backend
npm run seed
```

### Option 2: Run with Node.js directly

```bash
cd backend/migrations
node seed.js
```

### Option 3: Custom MongoDB URL

```bash
MONGO_URL="mongodb://localhost:27017" DB_NAME="magicchat" node seed.js
```

## What It Does

1. **Connects** to MongoDB
2. **Clears** existing data in all collections (optional - can be disabled in script)
3. **Inserts** all sample data:
   - Users
   - Videos
   - Comments
   - Likes
   - Follows
   - Hashtags
4. **Creates indexes** for better query performance
5. **Displays summary** of inserted data

## Video URLs

The seed data uses free sample videos from Google's test video bucket:
- These are real MP4 files that will play in your app
- Perfect for testing video playback
- No copyright issues

Sample videos include:
- Big Buck Bunny
- Elephant's Dream
- Sintel
- Tears of Steel

## Customizing the Data

Edit `seed_data.json` to:
- Add more users
- Add more videos
- Change video URLs to your own videos
- Modify user profiles
- Add more interactions

## Troubleshooting

### "Connection refused"
Make sure MongoDB is running on the specified URL.

### "Duplicate key error"
The script clears data by default. If you disabled clearing, make sure usernames and emails are unique.

### "Module not found"
Install MongoDB driver:
```bash
cd backend
npm install mongodb
```

## Database Collections

After seeding, your database will have:

| Collection | Count | Purpose |
|-----------|-------|---------|
| users | 8 | User accounts and profiles |
| videos | 12 | Video metadata and stats |
| comments | 8 | Video comments |
| likes | 5 | User-video like relationships |
| follows | 5 | User follow relationships |
| hashtags | 5 | Trending hashtags |

## Next Steps

After seeding:

1. Start your Go backend server
2. Test authentication with any of the test users
3. Browse the For You feed
4. Test video playback
5. Try liking and commenting on videos
6. Test the search functionality

## Production Warning

⚠️ **DO NOT run this seed script in production!**

This will delete all existing data. Only use for development and testing.