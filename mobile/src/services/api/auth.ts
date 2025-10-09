import { apiClient } from './client';
import AsyncStorage from '@react-native-async-storage/async-storage';

interface LoginCredentials {
  email: string;
  password: string;
}

interface RegisterData {
  username: string;
  email: string;
  password: string;
  display_name: string;
}

interface AuthResponse {
  token: string;
  user: any;
}

export const authApi = {
  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    const response = await apiClient.post('/api/auth/login', credentials);
    const { token } = response.data.data;
    await AsyncStorage.setItem('auth_token', token);
    return response.data.data;
  },

  async register(data: RegisterData): Promise<AuthResponse> {
    const response = await apiClient.post('/api/auth/register', data);
    const { token } = response.data.data;
    await AsyncStorage.setItem('auth_token', token);
    return response.data.data;
  },

  async logout(): Promise<void> {
    await apiClient.post('/api/auth/logout');
    await AsyncStorage.removeItem('auth_token');
  },

  async getCurrentUser(): Promise<any> {
    const response = await apiClient.get('/api/auth/me');
    return response.data.data;
  },
};