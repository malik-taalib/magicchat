import React, { useState, useEffect, useRef, useCallback } from 'react';
import {
  View,
  Text,
  FlatList,
  StyleSheet,
  Dimensions,
  ActivityIndicator,
  TouchableOpacity,
  Alert,
} from 'react-native';
import { useFocusEffect } from '@react-navigation/native';
import VideoPlayer from '@/components/VideoPlayer';
import ActionButtons from '@/components/shared/ActionButtons';
import { videosApi, usersApi } from '@/services/api';
import { Video } from '@/types/models';
import { useAuthStore } from '@/store';

const { height } = Dimensions.get('window');

export default function FeedScreen({ route, navigation }: any) {
  const { user } = useAuthStore();
  const [videos, setVideos] = useState<Video[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [likedVideos, setLikedVideos] = useState<Set<string>>(new Set());
  const [followedUsers, setFollowedUsers] = useState<Set<string>>(new Set());
  const flatListRef = useRef<FlatList>(null);

  const loadVideos = async () => {
    try {
      setLoading(true);

      // Check if we're navigating from profile with specific videos
      if (route?.params?.videos && route?.params?.initialVideoId) {
        setVideos(route.params.videos);

        // Find the index of the initial video and scroll to it
        const initialIndex = route.params.videos.findIndex(
          (v: Video) => v.id === route.params.initialVideoId
        );

        if (initialIndex >= 0) {
          setCurrentIndex(initialIndex);
          // Scroll to the video after render
          setTimeout(() => {
            flatListRef.current?.scrollToIndex({
              index: initialIndex,
              animated: false,
            });
          }, 100);
        }

        // Clear navigation params
        navigation.setParams({ videos: undefined, initialVideoId: undefined });
      } else {
        // Normal feed loading
        const response = await videosApi.getForYouFeed();
        setVideos(response.videos || []);
      }
    } catch (error) {
      console.error('Error loading videos:', error);
      Alert.alert('Error', 'Could not load videos');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadVideos();
    loadUserLikes();
  }, [route?.params?.videos]);

  const loadUserLikes = async () => {
    if (!user) return;

    try {
      const response = await usersApi.getUserLikedVideos(user.id);
      const likedVideoIds = new Set(response.videos.map(v => v.id));
      setLikedVideos(likedVideoIds);
    } catch (error) {
      console.error('Error loading user likes:', error);
    }
  };

  const onRefresh = async () => {
    setRefreshing(true);
    await loadVideos();
    await loadUserLikes();
    setRefreshing(false);
  };

  const onViewableItemsChanged = useRef(({ viewableItems }: any) => {
    if (viewableItems.length > 0) {
      setCurrentIndex(viewableItems[0].index || 0);
    }
  }).current;

  const viewabilityConfig = useRef({
    itemVisiblePercentThreshold: 80,
  }).current;

  const handleLike = async (videoId: string) => {
    const isLiked = likedVideos.has(videoId);

    // Optimistic UI update - do this FIRST for instant feedback
    if (isLiked) {
      setLikedVideos((prev) => {
        const newSet = new Set(prev);
        newSet.delete(videoId);
        return newSet;
      });
      setVideos((prev) =>
        prev.map((v) =>
          v.id === videoId
            ? { ...v, like_count: Math.max(0, v.like_count - 1) }
            : v
        )
      );
    } else {
      setLikedVideos((prev) => new Set(prev).add(videoId));
      setVideos((prev) =>
        prev.map((v) =>
          v.id === videoId
            ? { ...v, like_count: v.like_count + 1 }
            : v
        )
      );
    }

    // Then make API call in background
    try {
      if (isLiked) {
        await videosApi.unlikeVideo(videoId);
      } else {
        await videosApi.likeVideo(videoId);
      }
    } catch (error: any) {
      // Revert on error
      if (isLiked) {
        setLikedVideos((prev) => new Set(prev).add(videoId));
        setVideos((prev) =>
          prev.map((v) =>
            v.id === videoId
              ? { ...v, like_count: v.like_count + 1 }
              : v
          )
        );
      } else {
        setLikedVideos((prev) => {
          const newSet = new Set(prev);
          newSet.delete(videoId);
          return newSet;
        });
        setVideos((prev) =>
          prev.map((v) =>
            v.id === videoId
              ? { ...v, like_count: Math.max(0, v.like_count - 1) }
              : v
          )
        );
      }

      // Only show error if it's not a 409 (already liked/unliked)
      if (error?.response?.status !== 409) {
        console.error('Error toggling like:', error);
      }
    }
  };

  const handleComment = (videoId: string) => {
    Alert.alert('Comments', 'Comment feature coming soon!');
  };

  const handleShare = async (videoId: string) => {
    try {
      await videosApi.shareVideo(videoId);
      Alert.alert('Shared!', 'Video shared successfully');
    } catch (error) {
      console.error('Error sharing video:', error);
    }
  };

  const handleFollow = async (userId: string) => {
    try {
      // Optimistic UI update
      setFollowedUsers((prev) => new Set(prev).add(userId));

      // Make API call
      await usersApi.followUser(userId);
    } catch (error) {
      // Revert on error
      setFollowedUsers((prev) => {
        const newSet = new Set(prev);
        newSet.delete(userId);
        return newSet;
      });
      console.error('Error following user:', error);
      Alert.alert('Error', 'Could not follow user');
    }
  };

  const renderVideo = ({ item, index }: { item: Video; index: number }) => {
    const isCurrentVideo = index === currentIndex;

    return (
      <View style={styles.videoContainer}>
        <VideoPlayer uri={item.video_url} shouldPlay={isCurrentVideo} />

        {/* Video info */}
        <View style={styles.infoContainer}>
          <Text style={styles.username}>@{item.username}</Text>
          <Text style={styles.description} numberOfLines={2}>
            {item.description}
          </Text>
          {item.hashtags && item.hashtags.length > 0 && (
            <Text style={styles.hashtags}>
              {item.hashtags.map((tag) => `#${tag}`).join(' ')}
            </Text>
          )}
        </View>

        {/* Action buttons */}
        <ActionButtons
          video={item}
          onLike={() => handleLike(item.id)}
          onComment={() => handleComment(item.id)}
          onShare={() => handleShare(item.id)}
          onAvatarPress={() =>
            navigation.navigate('UserProfile', {
              userId: item.user_id,
              username: item.username,
            })
          }
          onFollowPress={() => handleFollow(item.user_id)}
          isLiked={likedVideos.has(item.id)}
          isFollowing={followedUsers.has(item.user_id)}
        />
      </View>
    );
  };

  if (loading && videos.length === 0) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#ff0050" />
        <Text style={styles.loadingText}>Loading videos...</Text>
      </View>
    );
  }

  if (videos.length === 0) {
    return (
      <View style={styles.emptyContainer}>
        <Text style={styles.emptyText}>No videos available</Text>
        <TouchableOpacity style={styles.retryButton} onPress={loadVideos}>
          <Text style={styles.retryButtonText}>Retry</Text>
        </TouchableOpacity>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <FlatList
        ref={flatListRef}
        data={videos}
        renderItem={renderVideo}
        keyExtractor={(item, index) => `${item.id}-${index}`}
        pagingEnabled
        showsVerticalScrollIndicator={false}
        snapToInterval={height}
        snapToAlignment="start"
        decelerationRate="fast"
        onViewableItemsChanged={onViewableItemsChanged}
        viewabilityConfig={viewabilityConfig}
        refreshing={refreshing}
        onRefresh={onRefresh}
        getItemLayout={(_, index) => ({
          length: height,
          offset: height * index,
          index,
        })}
        removeClippedSubviews={false}
        windowSize={3}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000',
  },
  videoContainer: {
    height,
    position: 'relative',
  },
  infoContainer: {
    position: 'absolute',
    left: 12,
    bottom: 100,
    right: 80,
  },
  username: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 8,
  },
  description: {
    color: '#fff',
    fontSize: 14,
    marginBottom: 8,
  },
  hashtags: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  loadingContainer: {
    flex: 1,
    backgroundColor: '#000',
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingText: {
    color: '#fff',
    fontSize: 16,
    marginTop: 16,
  },
  emptyContainer: {
    flex: 1,
    backgroundColor: '#000',
    justifyContent: 'center',
    alignItems: 'center',
  },
  emptyText: {
    color: '#fff',
    fontSize: 16,
    marginBottom: 24,
  },
  retryButton: {
    backgroundColor: '#ff0050',
    paddingHorizontal: 32,
    paddingVertical: 12,
    borderRadius: 8,
  },
  retryButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
});