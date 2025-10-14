import { apiClient } from './client';
import { User, Video } from '@/types/models';

interface UserVideosResponse {
  videos: Video[];
  cursor?: string;
  has_more: boolean;
}

export const usersApi = {
  async followUser(userId: string): Promise<void> {
    await apiClient.post(`/api/users/${userId}/follow`);
  },

  async unfollowUser(userId: string): Promise<void> {
    await apiClient.delete(`/api/users/${userId}/follow`);
  },

  async getFollowers(userId: string, cursor?: string): Promise<any> {
    const params = cursor ? { cursor } : {};
    const response = await apiClient.get(`/api/users/${userId}/followers`, { params });
    return response.data.data;
  },

  async getFollowing(userId: string, cursor?: string): Promise<any> {
    const params = cursor ? { cursor } : {};
    const response = await apiClient.get(`/api/users/${userId}/following`, { params });
    return response.data.data;
  },

  async getUserProfile(username: string): Promise<User> {
    const response = await apiClient.get(`/api/users/${username}`);
    return response.data.data;
  },

  async getUserVideos(userId: string, cursor?: string): Promise<UserVideosResponse> {
    const params = cursor ? { cursor } : {};
    // Use search endpoint to filter by user
    const response = await apiClient.get(`/api/search`, {
      params: { ...params, type: 'videos', q: userId, limit: 20 }
    });
    return response.data.data;
  },

  async getUserLikedVideos(userId: string, limit = 20, offset = 0): Promise<{ videos: Video[] }> {
    const response = await apiClient.get(`/api/users/${userId}/likes`, {
      params: { limit, offset }
    });
    return response.data.data;
  },
};