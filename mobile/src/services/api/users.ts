import { apiClient } from './client';
import { User } from '@/types/models';

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
};