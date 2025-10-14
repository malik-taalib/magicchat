import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  Image,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
  FlatList,
  Dimensions,
  ActivityIndicator,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useAuthStore } from '@/store';
import { usersApi, videosApi } from '@/services/api';
import { Video, User } from '@/types/models';

const { width } = Dimensions.get('window');
const ITEM_WIDTH = (width - 6) / 3;

type TabType = 'videos' | 'likes';

interface UserProfileScreenProps {
  navigation: any;
  route: {
    params: {
      userId: string;
      username: string;
    };
  };
}

export default function UserProfileScreen({ navigation, route }: UserProfileScreenProps) {
  const { user: currentUser } = useAuthStore();
  const { userId, username } = route.params;

  const [profileUser, setProfileUser] = useState<User | null>(null);
  const [activeTab, setActiveTab] = useState<TabType>('videos');
  const [videos, setVideos] = useState<Video[]>([]);
  const [loading, setLoading] = useState(false);
  const [isFollowing, setIsFollowing] = useState(false);

  useEffect(() => {
    loadProfileData();
  }, [userId]);

  useEffect(() => {
    if (profileUser) {
      loadVideos();
    }
  }, [profileUser, activeTab]);

  const loadProfileData = async () => {
    try {
      setLoading(true);
      // Fetch full user profile from API
      const user = await usersApi.getUserProfile(userId);
      setProfileUser(user);
    } catch (error) {
      console.error('Error loading profile:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadVideos = async () => {
    if (!profileUser) return;

    try {
      setLoading(true);

      if (activeTab === 'videos') {
        // Fetch user's videos by filtering feed
        const response = await videosApi.getForYouFeed();
        const userVideos = (response.videos || []).filter(
          (video) => video.user_id === profileUser.id
        );
        setVideos(userVideos);
      } else if (activeTab === 'likes') {
        // Fetch user's liked videos
        const response = await usersApi.getUserLikedVideos(profileUser.id);
        setVideos(response.videos || []);
      }
    } catch (error) {
      console.error('Error loading videos:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleFollowToggle = async () => {
    if (!profileUser) return;

    try {
      if (isFollowing) {
        await usersApi.unfollowUser(profileUser.id);
        setIsFollowing(false);
      } else {
        await usersApi.followUser(profileUser.id);
        setIsFollowing(true);
      }
    } catch (error) {
      console.error('Error toggling follow:', error);
    }
  };

  const handleVideoPress = (video: Video) => {
    navigation.navigate('Feed', {
      initialVideoId: video.id,
      videos: videos,
    });
  };

  const formatNumber = (num: number): string => {
    if (num >= 1000000) {
      return `${(num / 1000000).toFixed(1)}M`;
    }
    if (num >= 1000) {
      return `${(num / 1000).toFixed(1)}K`;
    }
    return num.toString();
  };

  const renderVideoItem = ({ item }: { item: Video }) => (
    <TouchableOpacity
      style={styles.videoItem}
      onPress={() => handleVideoPress(item)}
      activeOpacity={0.8}
    >
      <Image
        source={{ uri: item.thumbnail_url }}
        style={styles.videoThumbnail}
        resizeMode="cover"
      />
      <View style={styles.videoOverlay}>
        <View style={styles.videoStats}>
          <Ionicons name="play" size={14} color="#fff" />
          <Text style={styles.videoStatsText}>
            {formatNumber(item.view_count)}
          </Text>
        </View>
      </View>
    </TouchableOpacity>
  );

  if (loading && !profileUser) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#ff0050" />
      </View>
    );
  }

  if (!profileUser) {
    return (
      <View style={styles.container}>
        <Text style={styles.emptyText}>User not found</Text>
      </View>
    );
  }

  // Check if viewing own profile
  const isOwnProfile = currentUser?.id === profileUser.id;

  return (
    <View style={styles.container}>
      <ScrollView
        contentContainerStyle={styles.content}
        stickyHeaderIndices={[2]}
        showsVerticalScrollIndicator={false}
      >
        {/* Header */}
        <View style={styles.header}>
          <TouchableOpacity
            style={styles.backButton}
            onPress={() => navigation.goBack()}
          >
            <Ionicons name="arrow-back" size={24} color="#fff" />
          </TouchableOpacity>
          <Text style={styles.headerTitle}>@{profileUser.username}</Text>
          <View style={styles.headerButton} />
        </View>

        {/* Profile Info */}
        <View style={styles.profileSection}>
          <Image
            source={{ uri: profileUser.avatar_url || 'https://i.pravatar.cc/200' }}
            style={styles.avatar}
          />
          <Text style={styles.displayName}>{profileUser.display_name}</Text>
          <Text style={styles.username}>@{profileUser.username}</Text>

          {profileUser.bio && <Text style={styles.bio}>{profileUser.bio}</Text>}

          {/* Stats */}
          <View style={styles.statsContainer}>
            <View style={styles.stat}>
              <Text style={styles.statNumber}>
                {formatNumber(profileUser.following_count || 0)}
              </Text>
              <Text style={styles.statLabel}>Following</Text>
            </View>
            <View style={styles.stat}>
              <Text style={styles.statNumber}>
                {formatNumber(profileUser.follower_count || 0)}
              </Text>
              <Text style={styles.statLabel}>Followers</Text>
            </View>
            <View style={styles.stat}>
              <Text style={styles.statNumber}>
                {formatNumber(profileUser.total_likes || 0)}
              </Text>
              <Text style={styles.statLabel}>Likes</Text>
            </View>
          </View>

          {/* Action Buttons */}
          {!isOwnProfile && (
            <View style={styles.actionButtons}>
              <TouchableOpacity
                style={[styles.followButton, isFollowing && styles.followingButton]}
                onPress={handleFollowToggle}
              >
                <Text style={styles.followButtonText}>
                  {isFollowing ? 'Following' : 'Follow'}
                </Text>
              </TouchableOpacity>
              <TouchableOpacity style={styles.messageButton}>
                <Ionicons name="chatbubble-outline" size={20} color="#fff" />
              </TouchableOpacity>
            </View>
          )}
        </View>

        {/* Tabs */}
        <View style={styles.tabsContainer}>
          <View style={styles.tabs}>
            <TouchableOpacity
              style={[styles.tab, activeTab === 'videos' && styles.activeTab]}
              onPress={() => setActiveTab('videos')}
            >
              <Ionicons
                name="grid-outline"
                size={20}
                color={activeTab === 'videos' ? '#fff' : '#666'}
              />
            </TouchableOpacity>
            <TouchableOpacity
              style={[styles.tab, activeTab === 'likes' && styles.activeTab]}
              onPress={() => setActiveTab('likes')}
            >
              <Ionicons
                name="heart-outline"
                size={20}
                color={activeTab === 'likes' ? '#fff' : '#666'}
              />
            </TouchableOpacity>
          </View>
        </View>

        {/* Videos Grid */}
        {loading ? (
          <View style={styles.loadingContainer}>
            <ActivityIndicator size="large" color="#ff0050" />
          </View>
        ) : videos.length === 0 ? (
          <View style={styles.emptyVideosContainer}>
            <Ionicons
              name={activeTab === 'videos' ? 'videocam-outline' : 'heart-outline'}
              size={64}
              color="#333"
            />
            <Text style={styles.emptyVideosText}>
              {activeTab === 'videos' ? 'No videos yet' : 'No liked videos'}
            </Text>
          </View>
        ) : (
          <FlatList
            data={videos}
            renderItem={renderVideoItem}
            keyExtractor={(item, index) => `${item.id}-${index}`}
            numColumns={3}
            scrollEnabled={false}
            columnWrapperStyle={styles.row}
            contentContainerStyle={styles.videosGrid}
          />
        )}
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000',
  },
  content: {
    paddingBottom: 40,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#000',
    padding: 40,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 16,
    paddingTop: 60,
  },
  backButton: {
    padding: 8,
  },
  headerTitle: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  headerButton: {
    width: 40,
  },
  profileSection: {
    alignItems: 'center',
    paddingHorizontal: 24,
  },
  avatar: {
    width: 100,
    height: 100,
    borderRadius: 50,
    marginBottom: 16,
  },
  displayName: {
    color: '#fff',
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  username: {
    color: '#999',
    fontSize: 16,
    marginBottom: 16,
  },
  bio: {
    color: '#fff',
    fontSize: 14,
    textAlign: 'center',
    marginBottom: 24,
    lineHeight: 20,
  },
  statsContainer: {
    flexDirection: 'row',
    marginBottom: 24,
    gap: 40,
  },
  stat: {
    alignItems: 'center',
  },
  statNumber: {
    color: '#fff',
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  statLabel: {
    color: '#999',
    fontSize: 14,
  },
  actionButtons: {
    flexDirection: 'row',
    gap: 12,
    width: '100%',
  },
  followButton: {
    flex: 1,
    backgroundColor: '#ff0050',
    borderRadius: 8,
    paddingVertical: 12,
    alignItems: 'center',
  },
  followingButton: {
    backgroundColor: '#1a1a1a',
    borderWidth: 1,
    borderColor: '#333',
  },
  followButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  messageButton: {
    backgroundColor: '#1a1a1a',
    borderWidth: 1,
    borderColor: '#333',
    borderRadius: 8,
    paddingVertical: 12,
    paddingHorizontal: 16,
    justifyContent: 'center',
    alignItems: 'center',
  },
  tabsContainer: {
    backgroundColor: '#000',
  },
  tabs: {
    flexDirection: 'row',
    borderBottomWidth: 1,
    borderBottomColor: '#333',
    marginTop: 32,
  },
  tab: {
    flex: 1,
    alignItems: 'center',
    paddingVertical: 12,
    borderBottomWidth: 2,
    borderBottomColor: 'transparent',
  },
  activeTab: {
    borderBottomColor: '#fff',
  },
  videosGrid: {
    paddingTop: 2,
  },
  row: {
    gap: 2,
  },
  videoItem: {
    width: ITEM_WIDTH,
    height: ITEM_WIDTH * 1.5,
    marginBottom: 2,
    position: 'relative',
  },
  videoThumbnail: {
    width: '100%',
    height: '100%',
    backgroundColor: '#1a1a1a',
  },
  videoOverlay: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    padding: 8,
  },
  videoStats: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  videoStatsText: {
    color: '#fff',
    fontSize: 12,
    fontWeight: '600',
    textShadowColor: 'rgba(0, 0, 0, 0.75)',
    textShadowOffset: { width: 0, height: 1 },
    textShadowRadius: 3,
  },
  emptyVideosContainer: {
    padding: 60,
    alignItems: 'center',
    justifyContent: 'center',
  },
  emptyVideosText: {
    color: '#666',
    fontSize: 16,
    marginTop: 16,
  },
  emptyText: {
    color: '#666',
    fontSize: 16,
    textAlign: 'center',
    marginTop: 100,
  },
});