import { api } from './client';
import { mapStats } from './mappers';
import type { Stats } from '@/src/types';

export const statsApi = {
  getOverview: async () =>
    mapStats(await api.get<Record<string, unknown>>('/stats')),

  getTrend: async (days?: number) => {
    const res = await api.get<{ days: number; daily: { date: string; count: number }[] }>('/stats/trend', { days });
    return res.daily || [];
  },
};
