import { api } from './client';
import type { SyncTask, Skill, PaginatedResponse } from '@/src/types';

export const adminApi = {
  getSyncTasks: () =>
    api.get<SyncTask[]>('/admin/sync-tasks'),

  triggerSync: (type: 'full' | 'incremental') =>
    api.post<SyncTask>('/admin/sync-tasks', { type }),

  getPendingReviews: () =>
    api.get<Skill[]>('/admin/pending-review'),

  approveSkill: (id: number) =>
    api.put<void>(`/admin/skills/${id}/approve`),

  rejectSkill: (id: number, reason?: string) =>
    api.put<void>(`/admin/skills/${id}/reject`, { reason }),

  getSystemLogs: (lines?: number) =>
    api.get<{ timestamp: string; level: string; message: string }[]>('/admin/logs', { lines }),

  getUsers: (params?: { page?: number; pageSize?: number }) =>
    api.get<PaginatedResponse<{ id: number; username: string; email: string; role: string; createdAt: string }>>('/admin/users', params as Record<string, string | number | boolean | undefined>),
};
