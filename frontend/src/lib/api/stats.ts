import { api } from './client';
import type { Stats } from '@/src/types';

export const statsApi = {
  getOverview: () =>
    api.get<Stats>('/stats'),

  getTrend: (days?: number) =>
    api.get<{ date: string; newSkills: number; apiCalls: number }[]>('/stats/trend', { days }),
};
