// Data models matching the backend API

export interface User {
  id: string; // Backend returns "id" not "_id"
  username: string;
  email: string;
  display_name: string;
  bio: string;
  avatar_url: string;
  follower_count: number;
  following_count: number;
  video_count: number;
  total_likes: number;
  created_at: string;
  updated_at: string;
}

export interface Video {
  id: string; // Backend uses "id" not "_id" for feed videos
  user_id: string;
  username: string;
  display_name: string;
  avatar_url: string;
  title: string;
  description: string;
  video_url: string;
  thumbnail_url: string;
  duration: number;
  hashtags: string[];
  view_count: number;
  like_count: number;
  comment_count: number;
  share_count: number;
  processing_status: 'pending' | 'processing' | 'completed' | 'failed';
  created_at: string;
  updated_at: string;
}

export interface Comment {
  _id: string;
  user_id: string;
  video_id: string;
  text: string;
  likes: number;
  created_at: string;
}

export interface Notification {
  _id: string;
  user_id: string;
  type: 'like' | 'comment' | 'follow' | 'mention';
  from_user_id: string;
  video_id?: string;
  text: string;
  read: boolean;
  created_at: string;
}