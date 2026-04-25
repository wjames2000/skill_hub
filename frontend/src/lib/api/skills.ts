import { api } from './client';
import { mapSkill, mapSkillList } from './mappers';
import type { Review, Skill, SkillDetail, SearchFilters } from '@/src/types';

export const skillsApi = {
  list: async (params?: { page?: number; pageSize?: number; category?: string; sort?: string }) => {
    const res = await api.get<{ skills: Record<string, unknown>[]; total: number; page: number; size: number }>('/skills', params as Record<string, string | number | boolean | undefined>);
    return { data: mapSkillList(res.skills), total: res.total, page: res.page, pageSize: res.size };
  },

  getById: async (id: number) => {
    const raw = await api.get<Record<string, unknown>>(`/skills/${id}`);
    const skill = mapSkill(raw);
    return {
      ...skill,
      readme: raw.readme as string || '',
      files: (raw.files || []) as SkillDetail['files'],
      similarSkills: mapSkillList((raw.similar_skills || []) as Record<string, unknown>[]),
      reviews: (raw.reviews || []) as Review[],
    } as SkillDetail;
  },

  search: async (filters: SearchFilters) => {
    const res = await api.post<{ skills: Record<string, unknown>[]; total: number; page: number; size: number }>('/skills/search', filters);
    return { data: mapSkillList(res.skills), total: res.total, page: res.page, pageSize: res.size };
  },

  getTrending: async () => {
    const res = await api.get<{ skills: Record<string, unknown>[] }>('/skills/trending');
    return mapSkillList(res.skills);
  },

  getCategories: () =>
    api.get<{ id: number; name: string; icon: string; count: number }[]>('/skills/categories'),

  getLatest: async () => {
    const res = await api.get<{ skills: Record<string, unknown>[] }>('/skills/latest');
    return mapSkillList(res.skills);
  },

  addReview: (skillId: number, data: { rating: number; comment: string }) =>
    api.post<Review>(`/skills/${skillId}/reviews`, data),
};
