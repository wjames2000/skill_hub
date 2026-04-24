import { api } from './client';
import type { Skill, SkillDetail, SearchFilters, PaginatedResponse } from '@/src/types';

export const skillsApi = {
  list: (params?: { page?: number; pageSize?: number; category?: string; sort?: string }) =>
    api.get<PaginatedResponse<Skill>>('/skills', params as Record<string, string | number | boolean | undefined>),

  getById: (id: number) =>
    api.get<SkillDetail>(`/skills/${id}`),

  search: (filters: SearchFilters) =>
    api.post<PaginatedResponse<Skill>>('/skills/search', filters),

  getTrending: () =>
    api.get<Skill[]>('/skills/trending'),

  getCategories: () =>
    api.get<{ id: number; name: string; icon: string; count: number }[]>('/skills/categories'),

  getLatest: () =>
    api.get<Skill[]>('/skills/latest'),
};
