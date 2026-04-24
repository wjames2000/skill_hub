import { api, setToken, clearToken } from './client';
import type { LoginRequest, RegisterRequest, AuthResponse, User, ApiKey } from '@/src/types';

export const authApi = {
  login: async (data: LoginRequest) => {
    const res = await api.post<AuthResponse>('/auth/login', data);
    setToken(res.token);
    return res;
  },

  register: async (data: RegisterRequest) => {
    const res = await api.post<AuthResponse>('/auth/register', data);
    setToken(res.token);
    return res;
  },

  logout: () => {
    clearToken();
  },

  getProfile: () =>
    api.get<User>('/auth/profile'),

  updateProfile: (data: Partial<User>) =>
    api.put<User>('/auth/profile', data),

  githubLogin: () => {
    window.location.href = `${import.meta.env.VITE_API_BASE || '/api/v1'}/auth/github`;
  },

  getApiKeys: () =>
    api.get<ApiKey[]>('/auth/api-keys'),

  createApiKey: (name: string) =>
    api.post<ApiKey>('/auth/api-keys', { name }),

  revokeApiKey: (id: number) =>
    api.delete<void>(`/auth/api-keys/${id}`),
};
