import { apiClient } from './client';
import { Video } from '@/types/models';

interface FeedResponse {
  videos: Video[];
  cursor?: string;
  has_more: boolean;
}

export const videosApi = {
  async getForYouFeed(cursor?: string): Promise<FeedResponse> {
    const params = cursor ? { cursor } : {};
    const response = await apiClient.get('/api/feed/for-you', { params });
    return response.data.data;
  },

  async getFollowingFeed(cursor?: string): Promise<FeedResponse> {
    const params = cursor ? { cursor } : {};
    const response = await apiClient.get('/api/feed/following', { params });
    return response.data.data;
  },

  async getVideoById(id: string): Promise<Video> {
    const response = await apiClient.get(`/api/videos/${id}`);
    return response.data.data;
  },

  async likeVideo(videoId: string): Promise<void> {
    await apiClient.post(`/api/engage/${videoId}/like`);
  },

  async unlikeVideo(videoId: string): Promise<void> {
    await apiClient.delete(`/api/engage/${videoId}/like`);
  },

  async getComments(videoId: string, cursor?: string): Promise<any> {
    const params = cursor ? { cursor } : {};
    const response = await apiClient.get(`/api/engage/${videoId}/comments`, { params });
    return response.data.data;
  },

  async addComment(videoId: string, text: string): Promise<any> {
    const response = await apiClient.post(`/api/engage/${videoId}/comments`, { text });
    return response.data.data;
  },

  async shareVideo(videoId: string): Promise<void> {
    await apiClient.post(`/api/engage/${videoId}/share`);
  },
};