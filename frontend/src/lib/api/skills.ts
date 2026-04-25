import { api } from './client';
import { mapSkill, mapSkillList } from './mappers';
import type { Review, Skill, SkillDetail, SearchFilters, Category } from '@/src/types';

export const skillsApi = {
  list: async (params?: { page?: number; pageSize?: number; category?: string; sort?: string; tags?: string[]; safe?: boolean }) => {
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
    const body: Record<string, unknown> = { query: filters.query, page: filters.page, pageSize: filters.pageSize };
    if (filters.category) body.category = filters.category;
    if (filters.tags && filters.tags.length > 0) body.tags = filters.tags;
    body.safe = filters.safe ?? false;
    if (filters.sort) {
      const sortMap: Record<string, string> = { relevance: '', rating: 'score', downloads: 'installs' };
      body.sort = sortMap[filters.sort] || '';
    }
    const res = await api.post<{ skills: Record<string, unknown>[]; total: number; page: number; size: number }>('/skills/search', body);
    return { data: mapSkillList(res.skills), total: res.total, page: res.page, pageSize: res.size };
  },

  getTrending: async () => {
    const res = await api.get<{ skills: Record<string, unknown>[] }>('/skills/trending');
    return mapSkillList(res.skills);
  },

  getCategories: async () => {
    const raw = await api.get<{ id: number; name: string; zh_name: string; en_name: string; slug: string; icon: string; count: number; children: Record<string, unknown>[] }[]>('/skills/categories');
    const mapNode = (c: typeof raw[number]): Category => ({
      id: c.id,
      name: c.name,
      zhName: (c.zh_name || c.name) as string,
      enName: (c.en_name || c.name) as string,
      slug: c.slug,
      icon: c.icon || 'category',
      count: c.count,
      children: (c.children || []).map(mapNode),
    });
    return raw.map(mapNode);
  },

  getLatest: async () => {
    const res = await api.get<{ skills: Record<string, unknown>[] }>('/skills/latest');
    return mapSkillList(res.skills);
  },

  addReview: (skillId: number, data: { rating: number; comment: string }) =>
    api.post<Review>(`/skills/${skillId}/reviews`, data),
};
