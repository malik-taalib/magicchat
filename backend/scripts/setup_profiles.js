// Script to create full user profiles and associate videos
const { MongoClient, ObjectId } = require('mongodb');

const uri = 'mongodb://admin:password@localhost:27017/magicchat?authSource=admin';

async function setupProfiles() {
  const client = new MongoClient(uri);

  try {
    await client.connect();
    const db = client.db('magicchat');

    // Define complete user profiles
    const profiles = [
      {
        _id: new ObjectId('68edddb55eb09dae606aa122'),
        username: 'sarah_adventures',
        display_name: 'Sarah Smith',
        bio: '🌍 Travel enthusiast | Exploring the world one adventure at a time ✈️',
        avatar_url: 'https://i.pravatar.cc/300?img=1',
        follower_count: 12500,
        following_count: 450,
        video_count: 1,
        total_likes: 45000,
      },
      {
        _id: new ObjectId('68eddf355eb09dae606aa123'),
        username: 'chef_mike',
        display_name: 'Chef Mike',
        bio: '👨‍🍳 Professional chef | Quick recipes for busy people | DM for collabs',
        avatar_url: 'https://i.pravatar.cc/300?img=12',
        follower_count: 28000,
        following_count: 320,
        video_count: 1,
        total_likes: 120000,
      },
      {
        _id: new ObjectId('68eddf355eb09dae606aa124'),
        username: 'fitness_emma',
        display_name: 'Emma Fitness',
        bio: '💪 Certified personal trainer | Transform your body in 5 mins a day 🔥',
        avatar_url: 'https://i.pravatar.cc/300?img=5',
        follower_count: 45000,
        following_count: 180,
        video_count: 1,
        total_likes: 230000,
      },
      {
        _id: new ObjectId('68eddf355eb09dae606aa125'),
        username: 'comedy_jay',
        display_name: 'Jay Comedy',
        bio: '😂 Making you laugh daily | Stand-up comedian | Tag me in your fails',
        avatar_url: 'https://i.pravatar.cc/300?img=13',
        follower_count: 67000,
        following_count: 890,
        video_count: 1,
        total_likes: 580000,
      },
    ];

    // Update users with full profile data
    console.log('Updating user profiles...');
    for (const profile of profiles) {
      await db.collection('users').updateOne(
        { _id: profile._id },
        {
          $set: {
            username: profile.username,
            display_name: profile.display_name,
            bio: profile.bio,
            avatar_url: profile.avatar_url,
            follower_count: profile.follower_count,
            following_count: profile.following_count,
            video_count: profile.video_count,
            total_likes: profile.total_likes,
            updated_at: new Date(),
          }
        },
        { upsert: false }
      );
      console.log(`Updated profile for ${profile.username}`);
    }

    // Verify video associations are correct
    console.log('\nChecking video associations...');
    const videos = await db.collection('videos').find({}).toArray();
    for (const video of videos) {
      console.log(`Video "${video.title}" -> user_id: ${video.user_id}`);

      // Find the user
      const user = await db.collection('users').findOne({ _id: video.user_id });
      if (user) {
        console.log(`  ✓ Associated with ${user.username} (${user.display_name})`);
      } else {
        console.log(`  ✗ User not found!`);
      }
    }

    console.log('\n✅ Profile setup complete!');

  } catch (error) {
    console.error('Error:', error);
  } finally {
    await client.close();
  }
}

setupProfiles();
