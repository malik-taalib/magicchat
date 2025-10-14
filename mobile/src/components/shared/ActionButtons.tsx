import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Image } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Video } from '@/types/models';

interface ActionButtonsProps {
  video: Video;
  onLike: () => void;
  onComment: () => void;
  onShare: () => void;
  onAvatarPress?: () => void;
  onFollowPress?: () => void;
  isLiked?: boolean;
  isFollowing?: boolean;
}

export default function ActionButtons({
  video,
  onLike,
  onComment,
  onShare,
  onAvatarPress,
  onFollowPress,
  isLiked = false,
  isFollowing = false,
}: ActionButtonsProps) {
  const formatCount = (count: number): string => {
    if (count >= 1000000) {
      return `${(count / 1000000).toFixed(1)}M`;
    }
    if (count >= 1000) {
      return `${(count / 1000).toFixed(1)}K`;
    }
    return count.toString();
  };

  return (
    <View style={styles.container}>
      {/* Avatar */}
      <TouchableOpacity
        style={styles.avatarContainer}
        onPress={onAvatarPress}
        activeOpacity={0.7}
      >
        <Image
          source={{ uri: video.avatar_url || 'https://i.pravatar.cc/100' }}
          style={styles.avatar}
        />
        {!isFollowing && onFollowPress && (
          <TouchableOpacity
            style={styles.followButton}
            onPress={(e) => {
              e.stopPropagation();
              onFollowPress();
            }}
            activeOpacity={0.8}
          >
            <Ionicons name="add" size={16} color="#fff" />
          </TouchableOpacity>
        )}
      </TouchableOpacity>

      {/* Like button */}
      <TouchableOpacity style={styles.actionButton} onPress={onLike}>
        <Ionicons
          name={isLiked ? 'heart' : 'heart-outline'}
          size={32}
          color={isLiked ? '#ff0050' : '#fff'}
        />
        <Text style={styles.actionText}>{formatCount(video.like_count)}</Text>
      </TouchableOpacity>

      {/* Comment button */}
      <TouchableOpacity style={styles.actionButton} onPress={onComment}>
        <Ionicons name="chatbubble-outline" size={28} color="#fff" />
        <Text style={styles.actionText}>{formatCount(video.comment_count)}</Text>
      </TouchableOpacity>

      {/* Share button */}
      <TouchableOpacity style={styles.actionButton} onPress={onShare}>
        <Ionicons name="arrow-redo-outline" size={28} color="#fff" />
        <Text style={styles.actionText}>{formatCount(video.share_count)}</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    position: 'absolute',
    right: 12,
    bottom: 100,
    alignItems: 'center',
    gap: 24,
  },
  avatarContainer: {
    marginBottom: 8,
  },
  avatar: {
    width: 48,
    height: 48,
    borderRadius: 24,
    borderWidth: 2,
    borderColor: '#fff',
  },
  followButton: {
    position: 'absolute',
    bottom: -8,
    left: 12,
    width: 24,
    height: 24,
    borderRadius: 12,
    backgroundColor: '#ff0050',
    justifyContent: 'center',
    alignItems: 'center',
  },
  actionButton: {
    alignItems: 'center',
    gap: 4,
  },
  actionText: {
    color: '#fff',
    fontSize: 12,
    fontWeight: '600',
  },
});