// API Client for Magic Chat
const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';
const WS_URL = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/api';

export class APIError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'APIError';
  }
}

interface APIResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  metadata?: {
    cursor?: string;
    has_more?: boolean;
    total?: number;
  };
}

class APIClient {
  private token: string | null = null;

  constructor() {
    // Load token from localStorage if available
    if (typeof window !== 'undefined') {
      this.token = localStorage.getItem('auth_token');
    }
  }

  setToken(token: string) {
    this.token = token;
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', token);
    }
  }

  clearToken() {
    this.token = null;
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token');
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const headers: Record<string, string> = {
      ...(options.headers as Record<string, string>),
    };

    // Add auth token if available
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    // Add Content-Type for JSON requests
    if (options.body && typeof options.body === 'string') {
      headers['Content-Type'] = 'application/json';
    }

    const response = await fetch(`${API_URL}${endpoint}`, {
      ...options,
      headers,
    });

    // Check if response is ok before parsing
    if (!response.ok) {
      let errorMessage = `Request failed with status ${response.status}`;
      try {
        const errorData = await response.json();
        errorMessage = errorData.error || errorMessage;
      } catch {
        // If we can't parse the error as JSON, use the status text
        errorMessage = response.statusText || errorMessage;
      }
      throw new APIError(response.status, errorMessage);
    }

    // Parse JSON response
    let data: APIResponse<T>;
    try {
      data = await response.json();
    } catch {
      throw new APIError(response.status, 'Invalid JSON response from server');
    }

    if (!data.success) {
      throw new APIError(response.status, data.error || 'Request failed');
    }

    return data.data as T;
  }

  // Auth endpoints
  async register(username: string, email: string, password: string, displayName: string) {
    return this.request<{ token: string; user: User }>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ username, email, password, display_name: displayName }),
    });
  }

  async login(email: string, password: string) {
    return this.request<{ token: string; user: User }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
  }

  async logout() {
    await this.request('/auth/logout', { method: 'POST' });
    this.clearToken();
  }

  async getCurrentUser() {
    return this.request<User>('/auth/me');
  }

  // Video feed endpoints
  async getForYouFeed(cursor?: string, limit = 10) {
    const params = new URLSearchParams();
    if (cursor) params.append('cursor', cursor);
    params.append('limit', limit.toString());

    return this.request<FeedResponse>(`/feed/for-you?${params}`);
  }

  async getFollowingFeed(cursor?: string, limit = 10) {
    const params = new URLSearchParams();
    if (cursor) params.append('cursor', cursor);
    params.append('limit', limit.toString());

    return this.request<FeedResponse>(`/feed/following?${params}`);
  }

  async getVideo(id: string) {
    return this.request<Video>(`/feed/${id}`);
  }

  // Video upload endpoints
  async uploadVideo(file: File, title: string, description: string, hashtags: string[]) {
    const formData = new FormData();
    formData.append('video', file);
    formData.append('title', title);
    formData.append('description', description);
    hashtags.forEach(tag => formData.append('hashtags', tag));

    const response = await fetch(`${API_URL}/videos/upload`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.token}`,
      },
      body: formData,
    });

    const data: APIResponse<{ video_id: string; status: string }> = await response.json();
    if (!data.success) {
      throw new APIError(response.status, data.error || 'Upload failed');
    }
    return data.data!;
  }

  // Engagement endpoints
  async likeVideo(videoId: string) {
    return this.request(`/videos/${videoId}/like`, { method: 'POST' });
  }

  async unlikeVideo(videoId: string) {
    return this.request(`/videos/${videoId}/like`, { method: 'DELETE' });
  }

  async getComments(videoId: string, limit = 20, offset = 0) {
    return this.request<Comment[]>(`/videos/${videoId}/comments?limit=${limit}&offset=${offset}`);
  }

  async createComment(videoId: string, text: string, parentId?: string) {
    return this.request<Comment>(`/videos/${videoId}/comments`, {
      method: 'POST',
      body: JSON.stringify({ text, parent_id: parentId }),
    });
  }

  async shareVideo(videoId: string) {
    return this.request(`/videos/${videoId}/share`, { method: 'POST' });
  }

  // Following endpoints
  async followUser(userId: string) {
    return this.request(`/users/${userId}/follow`, { method: 'POST' });
  }

  async unfollowUser(userId: string) {
    return this.request(`/users/${userId}/follow`, { method: 'DELETE' });
  }

  async getFollowers(userId: string, cursor?: string, limit = 20) {
    const params = new URLSearchParams();
    if (cursor) params.append('cursor', cursor);
    params.append('limit', limit.toString());

    return this.request<FollowListResponse>(`/users/${userId}/followers?${params}`);
  }

  async getFollowing(userId: string, cursor?: string, limit = 20) {
    const params = new URLSearchParams();
    if (cursor) params.append('cursor', cursor);
    params.append('limit', limit.toString());

    return this.request<FollowListResponse>(`/users/${userId}/following?${params}`);
  }

  // Search endpoints
  async search(query: string, type: 'users' | 'videos' | 'hashtags', cursor?: string, limit = 20) {
    const params = new URLSearchParams({ q: query, type, limit: limit.toString() });
    if (cursor) params.append('cursor', cursor);

    return this.request<SearchResponse>(`/search?${params}`);
  }

  async getTrendingHashtags(limit = 20) {
    return this.request<Hashtag[]>(`/trending/hashtags?limit=${limit}`);
  }

  async getVideosByHashtag(tag: string, cursor?: string, limit = 20) {
    const params = new URLSearchParams({ limit: limit.toString() });
    if (cursor) params.append('cursor', cursor);

    return this.request<FeedResponse>(`/hashtags/${tag}/videos?${params}`);
  }

  // Notifications endpoints
  async getNotifications(cursor?: string, limit = 20) {
    const params = new URLSearchParams({ limit: limit.toString() });
    if (cursor) params.append('cursor', cursor);

    return this.request<NotificationListResponse>(`/notifications?${params}`);
  }

  async markNotificationAsRead(notificationId: string) {
    return this.request(`/notifications/${notificationId}/read`, { method: 'PUT' });
  }

  async markAllNotificationsAsRead() {
    return this.request('/notifications/read-all', { method: 'PUT' });
  }

  // WebSocket connection for notifications
  connectNotifications(onMessage: (notification: Notification) => void): WebSocket | null {
    if (typeof window === 'undefined' || !this.token) return null;

    const ws = new WebSocket(`${WS_URL}/notifications/stream`);

    ws.onopen = () => {
      // Send auth token
      ws.send(JSON.stringify({ type: 'auth', token: this.token }));
    };

    ws.onmessage = (event) => {
      try {
        const notification = JSON.parse(event.data);
        onMessage(notification);
      } catch (error) {
        console.error('Failed to parse notification:', error);
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    ws.onclose = () => {
      console.log('WebSocket connection closed');
    };

    return ws;
  }
}

// Type definitions
export interface User {
  id: string;
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
  id: string;
  user_id: string;
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
  processing_status: string;
  created_at: string;
  updated_at: string;
  user?: {
    username: string;
    display_name: string;
    avatar_url: string;
  };
}

export interface FeedResponse {
  videos: Video[];
  next_cursor?: string;
  has_more: boolean;
}

export interface Comment {
  id: string;
  user_id: string;
  video_id: string;
  text: string;
  parent_id?: string;
  replies: Comment[];
  created_at: string;
  updated_at: string;
}

export interface FollowListResponse {
  users: User[];
  next_cursor?: string;
  has_more: boolean;
  total: number;
}

export interface SearchResponse {
  users?: User[];
  videos?: Video[];
  hashtags?: Hashtag[];
  next_cursor?: string;
  has_more: boolean;
}

export interface Hashtag {
  tag: string;
  video_count: number;
  trending_score: number;
}

export interface Notification {
  id: string;
  user_id: string;
  type: 'like' | 'comment' | 'follow' | 'mention';
  actor_id: string;
  video_id?: string;
  comment_id?: string;
  text: string;
  read: boolean;
  created_at: string;
  actor: {
    id: string;
    username: string;
    display_name: string;
    avatar_url: string;
  };
}

export interface NotificationListResponse {
  notifications: Notification[];
  unread_count: number;
  next_cursor?: string;
  has_more: boolean;
}

// Export singleton instance
const apiClient = new APIClient();
export default apiClient;
