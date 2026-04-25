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
    api.get<User>('/user/profile'),

  updateProfile: (data: Partial<User>) =>
    api.put<User>('/user/profile', data),

  getApiKeys: () =>
    api.get<ApiKey[]>('/user/api-keys'),

  createApiKey: (name: string) =>
    api.post<ApiKey>('/user/api-keys', { name }),

  revokeApiKey: (id: number) =>
    api.delete<void>(`/user/api-keys/${id}`),
};
