import React, { useState, useEffect, useCallback } from 'react';
import {
  View,
  Text,
  Image,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
  Alert,
  FlatList,
  Dimensions,
  ActivityIndicator,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useFocusEffect } from '@react-navigation/native';
import { useAuthStore } from '@/store';
import { authApi, videosApi, usersApi } from '@/services/api';
import { Video } from '@/types/models';

const { width } = Dimensions.get('window');
const ITEM_WIDTH = (width - 6) / 3;

type TabType = 'videos' | 'likes' | 'bookmarks';

export default function ProfileScreen({ navigation }: any) {
  const { user, setUser, logout } = useAuthStore();
  const [activeTab, setActiveTab] = useState<TabType>('videos');
  const [videos, setVideos] = useState<Video[]>([]);
  const [loading, setLoading] = useState(false);

  // Refresh user data when screen is focused
  useFocusEffect(
    useCallback(() => {
      refreshUserData();
    }, [])
  );

  useEffect(() => {
    if (user) {
      loadVideos();
    }
  }, [user, activeTab]);

  const refreshUserData = async () => {
    try {
      const freshUser = await authApi.getCurrentUser();
      setUser(freshUser);
    } catch (error) {
      console.error('Error refreshing user data:', error);
    }
  };

  const loadVideos = async () => {
    if (!user) return;

    try {
      setLoading(true);

      if (activeTab === 'videos') {
        // Fetch user's own videos by filtering feed by user_id
        const response = await videosApi.getForYouFeed();
        const userVideos = (response.videos || []).filter(
          (video) => video.user_id === user.id
        );
        setVideos(userVideos);
      } else if (activeTab === 'likes') {
        // Fetch user's liked videos
        const response = await usersApi.getUserLikedVideos(user.id);
        setVideos(response.videos || []);
      } else if (activeTab === 'bookmarks') {
        // For now, show empty until we implement backend endpoint
        // TODO: Implement getUserBookmarkedVideos API endpoint
        setVideos([]);
      }
    } catch (error) {
      console.error('Error loading videos:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    Alert.alert('Logout', 'Are you sure you want to logout?', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Logout',
        style: 'destructive',
        onPress: async () => {
          try {
            await authApi.logout();
            logout();
          } catch (error) {
            console.error('Logout error:', error);
            logout(); // Logout anyway
          }
        },
      },
    ]);
  };

  const handleVideoPress = (video: Video) => {
    // Navigate to feed screen starting at this video
    // We'll pass the video data and navigate to the feed
    navigation.navigate('Feed', {
      initialVideoId: video.id,
      videos: videos
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

  if (!user) {
    return (
      <View style={styles.container}>
        <Text style={styles.emptyText}>Not logged in</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <ScrollView
        contentContainerStyle={styles.content}
        stickyHeaderIndices={[2]}
        showsVerticalScrollIndicator={false}
      >
        {/* Header */}
        <View style={styles.header}>
          <TouchableOpacity style={styles.settingsButton}>
            <Ionicons name="settings-outline" size={24} color="#fff" />
          </TouchableOpacity>
        </View>

        {/* Profile Info */}
        <View style={styles.profileSection}>
          <Image
            source={{ uri: user.avatar_url || 'https://i.pravatar.cc/200' }}
            style={styles.avatar}
          />
          <Text style={styles.displayName}>{user.display_name}</Text>
          <Text style={styles.username}>@{user.username}</Text>

          {user.bio && <Text style={styles.bio}>{user.bio}</Text>}

          {/* Stats */}
          <View style={styles.statsContainer}>
            <View style={styles.stat}>
              <Text style={styles.statNumber}>
                {formatNumber(user.following_count || 0)}
              </Text>
              <Text style={styles.statLabel}>Following</Text>
            </View>
            <View style={styles.stat}>
              <Text style={styles.statNumber}>
                {formatNumber(user.follower_count || 0)}
              </Text>
              <Text style={styles.statLabel}>Followers</Text>
            </View>
            <View style={styles.stat}>
              <Text style={styles.statNumber}>
                {formatNumber(user.total_likes || 0)}
              </Text>
              <Text style={styles.statLabel}>Likes</Text>
            </View>
          </View>

          {/* Action Buttons */}
          <View style={styles.actionButtons}>
            <TouchableOpacity
              style={styles.editButton}
              onPress={() => navigation.navigate('EditProfile')}
            >
              <Text style={styles.editButtonText}>Edit Profile</Text>
            </TouchableOpacity>
            <TouchableOpacity style={styles.logoutButton} onPress={handleLogout}>
              <Ionicons name="log-out-outline" size={20} color="#fff" />
            </TouchableOpacity>
          </View>
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
            <TouchableOpacity
              style={[styles.tab, activeTab === 'bookmarks' && styles.activeTab]}
              onPress={() => setActiveTab('bookmarks')}
            >
              <Ionicons
                name="bookmark-outline"
                size={20}
                color={activeTab === 'bookmarks' ? '#fff' : '#666'}
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
              name={
                activeTab === 'videos'
                  ? 'videocam-outline'
                  : activeTab === 'likes'
                  ? 'heart-outline'
                  : 'bookmark-outline'
              }
              size={64}
              color="#333"
            />
            <Text style={styles.emptyVideosText}>
              {activeTab === 'videos'
                ? 'No videos yet'
                : activeTab === 'likes'
                ? 'No liked videos'
                : 'No bookmarks'}
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
  header: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    padding: 16,
    paddingTop: 60,
  },
  settingsButton: {
    padding: 8,
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
  editButton: {
    flex: 1,
    backgroundColor: '#1a1a1a',
    borderWidth: 1,
    borderColor: '#333',
    borderRadius: 8,
    paddingVertical: 12,
    alignItems: 'center',
  },
  editButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  logoutButton: {
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
  loadingContainer: {
    padding: 40,
    alignItems: 'center',
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