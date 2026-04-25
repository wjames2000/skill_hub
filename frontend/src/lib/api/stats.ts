import { api } from './client';
import { mapStats } from './mappers';
import type { Stats } from '@/src/types';

export const statsApi = {
  getOverview: async () =>
    mapStats(await api.get<Record<string, unknown>>('/stats')),

  getTrend: (days?: number) =>
    api.get<{ date: string; newSkills: number; apiCalls: number }[]>('/stats/trend', { days }),
};
