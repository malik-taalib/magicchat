#!/usr/bin/env node

const { MongoClient, ObjectId } = require('mongodb');
const axios = require('axios');

const MONGO_URL = process.env.MONGO_URL || 'mongodb://localhost:27017';
const DB_NAME = process.env.DB_NAME || 'magicchat';
const API_URL = process.env.API_URL || 'http://localhost:8080';

// Sample users data
const users = [
  {
    username: 'sarah_adventures',
    email: 'sarah@magicchat.com',
    password: 'password123',
    display_name: 'Sarah Adventures',
    bio: '🌍 Travel vlogger | 📸 Photography | ✨ Living my best life',
    avatar_url: 'https://i.pravatar.cc/300?img=1',
  },
  {
    username: 'chef_mike',
    email: 'mike@magicchat.com',
    password: 'password123',
    display_name: 'Chef Mike',
    bio: '👨‍🍳 Professional Chef | 🍜 Quick recipes | DM for collabs',
    avatar_url: 'https://i.pravatar.cc/300?img=12',
  },
  {
    username: 'fitness_emma',
    email: 'emma@magicchat.com',
    password: 'password123',
    display_name: 'Emma Fitness',
    bio: '💪 Certified trainer | 🏋️‍♀️ Home workouts | Transform your body!',
    avatar_url: 'https://i.pravatar.cc/300?img=5',
  },
  {
    username: 'comedy_jay',
    email: 'jay@magicchat.com',
    password: 'password123',
    display_name: 'Jay Comedy',
    bio: '😂 Making you laugh daily | 🎭 Comedian | Follow for daily dose of humor',
    avatar_url: 'https://i.pravatar.cc/300?img=15',
  },
  {
    username: 'tech_alex',
    email: 'alex@magicchat.com',
    password: 'password123',
    display_name: 'Tech Alex',
    bio: '💻 Software Engineer | 📱 Tech reviews | 🚀 Latest gadgets',
    avatar_url: 'https://i.pravatar.cc/300?img=8',
  },
  {
    username: 'music_lisa',
    email: 'lisa@magicchat.com',
    password: 'password123',
    display_name: 'Lisa Music',
    bio: '🎵 Singer | 🎸 Guitar covers | 🎤 Original songs',
    avatar_url: 'https://i.pravatar.cc/300?img=9',
  },
  {
    username: 'art_marcus',
    email: 'marcus@magicchat.com',
    password: 'password123',
    display_name: 'Marcus Art',
    bio: '🎨 Digital artist | ✏️ Drawing tutorials | 🖼️ Commissions open',
    avatar_url: 'https://i.pravatar.cc/300?img=13',
  },
  {
    username: 'dance_nina',
    email: 'nina@magicchat.com',
    password: 'password123',
    display_name: 'Nina Dance',
    bio: '💃 Professional dancer | 🩰 Choreography | Join my dance challenges!',
    avatar_url: 'https://i.pravatar.cc/300?img=20',
  },
];

async function createUsers() {
  console.log('\n👥 Creating users via API...');
  const createdUsers = [];

  for (const user of users) {
    try {
      const response = await axios.post(`${API_URL}/api/auth/register`, user);
      createdUsers.push({
        ...user,
        id: response.data.data.user.id,
        objectId: new ObjectId(response.data.data.user.id),
      });
      console.log(`✅ Created user: ${user.username}`);
    } catch (error) {
      console.error(`❌ Failed to create ${user.username}:`, error.response?.data || error.message);
    }
  }

  return createdUsers;
}

async function seedDatabase() {
  const client = new MongoClient(MONGO_URL);

  try {
    console.log('🔌 Connecting to MongoDB...');
    await client.connect();
    console.log('✅ Connected to MongoDB');

    const db = client.db(DB_NAME);

    // Create users via API (this ensures proper password hashing)
    const createdUsers = await createUsers();

    if (createdUsers.length === 0) {
      console.error('❌ No users created. Exiting.');
      return;
    }

    // Update users with additional profile data
    console.log('\n📝 Updating user profiles...');
    for (const user of createdUsers) {
      await db.collection('users').updateOne(
        { _id: user.objectId },
        {
          $set: {
            bio: user.bio,
            avatar_url: user.avatar_url,
            follower_count: Math.floor(Math.random() * 50000) + 1000,
            following_count: Math.floor(Math.random() * 500) + 50,
            video_count: Math.floor(Math.random() * 200) + 10,
            total_likes: Math.floor(Math.random() * 1000000) + 10000,
          },
        }
      );
    }
    console.log('✅ Updated user profiles');

    // Insert videos with proper user references
    console.log('\n🎥 Inserting videos...');
    const videos = [
      {
        user_id: createdUsers[0].objectId, // sarah
        title: 'Sunset in Santorini',
        description: 'The most beautiful sunset I\'ve ever seen! 🌅 #travel #greece #santorini',
        video_url: 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4',
        thumbnail_url: 'https://picsum.photos/400/600?random=1',
        duration: 30,
        hashtags: ['travel', 'greece', 'santorini', 'sunset'],
        view_count: 45230,
        like_count: 3421,
        comment_count: 234,
        share_count: 156,
        processing_status: 'completed',
        created_at: new Date('2024-10-10T14:30:00Z'),
        updated_at: new Date(),
      },
      {
        user_id: createdUsers[1].objectId, // chef_mike
        title: '30-Second Pasta Recipe',
        description: 'Quick and delicious pasta hack! 🍝 #cooking #recipe #pasta #quickmeals',
        video_url: 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4',
        thumbnail_url: 'https://picsum.photos/400/600?random=2',
        duration: 45,
        hashtags: ['cooking', 'recipe', 'pasta', 'quickmeals'],
        view_count: 128450,
        like_count: 8932,
        comment_count: 567,
        share_count: 892,
        processing_status: 'completed',
        created_at: new Date('2024-10-11T10:15:00Z'),
        updated_at: new Date(),
      },
      {
        user_id: createdUsers[2].objectId, // fitness_emma
        title: '5-Min Ab Workout',
        description: 'Burn those abs! No equipment needed 💪 #fitness #workout #abs #homeworkout',
        video_url: 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4',
        thumbnail_url: 'https://picsum.photos/400/600?random=3',
        duration: 60,
        hashtags: ['fitness', 'workout', 'abs', 'homeworkout'],
        view_count: 234560,
        like_count: 15678,
        comment_count: 1234,
        share_count: 2341,
        processing_status: 'completed',
        created_at: new Date('2024-10-12T07:00:00Z'),
        updated_at: new Date(),
      },
      {
        user_id: createdUsers[3].objectId, // comedy_jay
        title: 'When you forget your password',
        description: 'We\'ve all been there 😂 #comedy #funny #relatable #humor',
        video_url: 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4',
        thumbnail_url: 'https://picsum.photos/400/600?random=4',
        duration: 25,
        hashtags: ['comedy', 'funny', 'relatable', 'humor'],
        view_count: 567890,
        like_count: 34567,
        comment_count: 2345,
        share_count: 4567,
        processing_status: 'completed',
        created_at: new Date('2024-10-12T16:45:00Z'),
        updated_at: new Date(),
      },
      {
        user_id: createdUsers[4].objectId, // tech_alex
        title: 'iPhone 15 Pro Review',
        description: 'Is it worth the upgrade? 📱 #tech #iphone #review #apple',
        video_url: 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerFun.mp4',
        thumbnail_url: 'https://picsum.photos/400/600?random=5',
        duration: 90,
        hashtags: ['tech', 'iphone', 'review', 'apple'],
        view_count: 89234,
        like_count: 5678,
        comment_count: 789,
        share_count: 456,
        processing_status: 'completed',
        created_at: new Date('2024-10-11T12:30:00Z'),
        updated_at: new Date(),
      },
      {
        user_id: createdUsers[5].objectId, // music_lisa
        title: 'Acoustic Cover - Perfect',
        description: 'Ed Sheeran cover 🎸 Hope you like it! #music #cover #guitar #singing',
        video_url: 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerJoyrides.mp4',
        thumbnail_url: 'https://picsum.photos/400/600?random=6',
        duration: 120,
        hashtags: ['music', 'cover', 'guitar', 'singing'],
        view_count: 456789,
        like_count: 28934,
        comment_count: 1567,
        share_count: 3421,
        processing_status: 'completed',
        created_at: new Date('2024-10-10T18:00:00Z'),
        updated_at: new Date(),
      },
      {
        user_id: createdUsers[6].objectId, // art_marcus
        title: 'Speed Drawing a Portrait',
        description: 'Watch me draw in 60 seconds ✏️ #art #drawing #timelapse #artist',
        video_url: 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerMeltdowns.mp4',
        thumbnail_url: 'https://picsum.photos/400/600?random=7',
        duration: 60,
        hashtags: ['art', 'drawing', 'timelapse', 'artist'],
        view_count: 67890,
        like_count: 4567,
        comment_count: 345,
        share_count: 234,
        processing_status: 'completed',
        created_at: new Date('2024-10-12T09:20:00Z'),
        updated_at: new Date(),
      },
      {
        user_id: createdUsers[7].objectId, // dance_nina
        title: 'Hip Hop Dance Tutorial',
        description: 'Learn this move! 💃 #dance #hiphop #tutorial #dancechallenge',
        video_url: 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/Sintel.mp4',
        thumbnail_url: 'https://picsum.photos/400/600?random=8',
        duration: 75,
        hashtags: ['dance', 'hiphop', 'tutorial', 'dancechallenge'],
        view_count: 345678,
        like_count: 23456,
        comment_count: 1234,
        share_count: 2345,
        processing_status: 'completed',
        created_at: new Date('2024-10-11T15:40:00Z'),
        updated_at: new Date(),
      },
    ];

    const insertedVideos = await db.collection('videos').insertMany(videos);
    console.log(`✅ Inserted ${Object.keys(insertedVideos.insertedIds).length} videos`);

    // Insert hashtags
    console.log('\n🏷️  Inserting hashtags...');
    const hashtags = [
      { tag: 'travel', video_count: 234, trending_score: 8.5, created_at: new Date(), updated_at: new Date() },
      { tag: 'cooking', video_count: 567, trending_score: 9.2, created_at: new Date(), updated_at: new Date() },
      { tag: 'fitness', video_count: 892, trending_score: 9.8, created_at: new Date(), updated_at: new Date() },
      { tag: 'comedy', video_count: 1234, trending_score: 9.9, created_at: new Date(), updated_at: new Date() },
      { tag: 'music', video_count: 678, trending_score: 8.7, created_at: new Date(), updated_at: new Date() },
    ];

    await db.collection('hashtags').insertMany(hashtags);
    console.log(`✅ Inserted ${hashtags.length} hashtags`);

    // Create indexes
    console.log('\n📊 Creating indexes...');
    await db.collection('users').createIndex({ username: 1 }, { unique: true });
    await db.collection('users').createIndex({ email: 1 }, { unique: true });
    await db.collection('videos').createIndex({ user_id: 1 });
    await db.collection('videos').createIndex({ created_at: -1 });
    await db.collection('videos').createIndex({ hashtags: 1 });
    await db.collection('hashtags').createIndex({ tag: 1 }, { unique: true });
    console.log('✅ Indexes created');

    console.log('\n🎉 Database seeding completed successfully!');
    console.log('\n📊 Summary:');
    console.log(`   Users: ${createdUsers.length}`);
    console.log(`   Videos: ${videos.length}`);
    console.log(`   Hashtags: ${hashtags.length}`);
    console.log('\n💡 Test credentials:');
    console.log('   Email: sarah@magicchat.com');
    console.log('   Password: password123');
    console.log('\n   (All users have password: password123)');

  } catch (error) {
    console.error('❌ Error seeding database:', error);
    process.exit(1);
  } finally {
    await client.close();
    console.log('\n🔌 Disconnected from MongoDB');
  }
}

seedDatabase();