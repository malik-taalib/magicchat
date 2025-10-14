#!/usr/bin/env node

const { MongoClient } = require('mongodb');
const fs = require('fs');
const path = require('path');

// MongoDB connection URL
const MONGO_URL = process.env.MONGO_URL || 'mongodb://localhost:27017';
const DB_NAME = process.env.DB_NAME || 'magicchat';

async function seedDatabase() {
  const client = new MongoClient(MONGO_URL);

  try {
    console.log('🔌 Connecting to MongoDB...');
    await client.connect();
    console.log('✅ Connected to MongoDB');

    const db = client.db(DB_NAME);

    // Read seed data
    const seedDataPath = path.join(__dirname, 'seed_data.json');
    const seedData = JSON.parse(fs.readFileSync(seedDataPath, 'utf8'));

    // Clear existing data (optional - comment out if you want to keep existing data)
    console.log('\n🗑️  Clearing existing data...');
    await db.collection('users').deleteMany({});
    await db.collection('videos').deleteMany({});
    await db.collection('comments').deleteMany({});
    await db.collection('likes').deleteMany({});
    await db.collection('follows').deleteMany({});
    await db.collection('hashtags').deleteMany({});
    console.log('✅ Existing data cleared');

    // Insert users
    console.log('\n👥 Inserting users...');
    if (seedData.users && seedData.users.length > 0) {
      await db.collection('users').insertMany(seedData.users);
      console.log(`✅ Inserted ${seedData.users.length} users`);
    }

    // Insert videos
    console.log('\n🎥 Inserting videos...');
    if (seedData.videos && seedData.videos.length > 0) {
      await db.collection('videos').insertMany(seedData.videos);
      console.log(`✅ Inserted ${seedData.videos.length} videos`);
    }

    // Insert comments
    console.log('\n💬 Inserting comments...');
    if (seedData.comments && seedData.comments.length > 0) {
      await db.collection('comments').insertMany(seedData.comments);
      console.log(`✅ Inserted ${seedData.comments.length} comments`);
    }

    // Insert likes
    console.log('\n❤️  Inserting likes...');
    if (seedData.likes && seedData.likes.length > 0) {
      await db.collection('likes').insertMany(seedData.likes);
      console.log(`✅ Inserted ${seedData.likes.length} likes`);
    }

    // Insert follows
    console.log('\n👤 Inserting follows...');
    if (seedData.follows && seedData.follows.length > 0) {
      await db.collection('follows').insertMany(seedData.follows);
      console.log(`✅ Inserted ${seedData.follows.length} follows`);
    }

    // Insert hashtags
    console.log('\n🏷️  Inserting hashtags...');
    if (seedData.hashtags && seedData.hashtags.length > 0) {
      await db.collection('hashtags').insertMany(seedData.hashtags);
      console.log(`✅ Inserted ${seedData.hashtags.length} hashtags`);
    }

    // Create indexes for better performance
    console.log('\n📊 Creating indexes...');
    await db.collection('users').createIndex({ username: 1 }, { unique: true });
    await db.collection('users').createIndex({ email: 1 }, { unique: true });
    await db.collection('videos').createIndex({ user_id: 1 });
    await db.collection('videos').createIndex({ created_at: -1 });
    await db.collection('videos').createIndex({ hashtags: 1 });
    await db.collection('comments').createIndex({ video_id: 1 });
    await db.collection('likes').createIndex({ user_id: 1, video_id: 1 }, { unique: true });
    await db.collection('follows').createIndex({ follower_id: 1 });
    await db.collection('follows').createIndex({ following_id: 1 });
    await db.collection('hashtags').createIndex({ tag: 1 }, { unique: true });
    console.log('✅ Indexes created');

    console.log('\n🎉 Database seeding completed successfully!');
    console.log('\n📊 Summary:');
    console.log(`   Users: ${seedData.users.length}`);
    console.log(`   Videos: ${seedData.videos.length}`);
    console.log(`   Comments: ${seedData.comments.length}`);
    console.log(`   Likes: ${seedData.likes.length}`);
    console.log(`   Follows: ${seedData.follows.length}`);
    console.log(`   Hashtags: ${seedData.hashtags.length}`);
    console.log('\n💡 Test credentials:');
    console.log('   Email: sarah@magicchat.com');
    console.log('   Password: password123');

  } catch (error) {
    console.error('❌ Error seeding database:', error);
    process.exit(1);
  } finally {
    await client.close();
    console.log('\n🔌 Disconnected from MongoDB');
  }
}

// Run the seed function
seedDatabase();