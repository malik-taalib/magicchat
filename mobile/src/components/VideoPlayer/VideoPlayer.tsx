import React, { useState, useRef, useEffect } from 'react';
import { View, StyleSheet, Dimensions, TouchableOpacity } from 'react-native';
import { VideoView, useVideoPlayer } from 'expo-video';
import { Ionicons } from '@expo/vector-icons';

const { width, height } = Dimensions.get('window');

interface VideoPlayerProps {
  uri: string;
  shouldPlay?: boolean;
}

export default function VideoPlayer({
  uri,
  shouldPlay = false,
}: VideoPlayerProps) {
  const [isMuted, setIsMuted] = useState(false);
  const [isPlaying, setIsPlaying] = useState(shouldPlay);

  const player = useVideoPlayer(uri, (player) => {
    player.loop = true;
    player.muted = isMuted;
  });

  useEffect(() => {
    if (shouldPlay && !player.playing) {
      player.play();
      setIsPlaying(true);
    } else if (!shouldPlay && player.playing) {
      player.pause();
      setIsPlaying(false);
    }
  }, [shouldPlay, player]);

  useEffect(() => {
    if (player) {
      player.muted = isMuted;
    }
  }, [isMuted, player]);

  const togglePlayPause = () => {
    if (player.playing) {
      player.pause();
      setIsPlaying(false);
    } else {
      player.play();
      setIsPlaying(true);
    }
  };

  const toggleMute = () => {
    setIsMuted(!isMuted);
  };

  return (
    <View style={styles.container}>
      <VideoView
        style={styles.video}
        player={player}
        contentFit="cover"
        nativeControls={false}
      />

      {/* Play/Pause overlay */}
      <TouchableOpacity
        style={styles.overlay}
        onPress={togglePlayPause}
        activeOpacity={1}
      >
        {!isPlaying && (
          <View style={styles.playButton}>
            <Ionicons name="play" size={60} color="#fff" />
          </View>
        )}
      </TouchableOpacity>

      {/* Mute button */}
      <TouchableOpacity style={styles.muteButton} onPress={toggleMute}>
        <Ionicons
          name={isMuted ? 'volume-mute' : 'volume-high'}
          size={24}
          color="#fff"
        />
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    width,
    height,
    backgroundColor: '#000',
  },
  video: {
    width: '100%',
    height: '100%',
  },
  overlay: {
    ...StyleSheet.absoluteFillObject,
    justifyContent: 'center',
    alignItems: 'center',
  },
  playButton: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  muteButton: {
    position: 'absolute',
    top: 60,
    right: 16,
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'center',
    alignItems: 'center',
  },
});
