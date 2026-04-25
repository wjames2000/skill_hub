import { api } from './client';
import type { SyncTask, Skill, PaginatedResponse } from '@/src/types';

function mapSyncTask(raw: Record<string, unknown>): SyncTask {
  const nameMap: Record<string, string> = { full: '全量同步', incremental: '增量同步' };
  const total = (raw.total_repos as number) || 0;
  const found = (raw.found_repos as number) || 0;
  return {
    id: raw.id as number,
    name: nameMap[raw.type as string] || (raw.type as string),
    type: raw.type as 'full' | 'incremental',
    status: (raw.status as SyncTask['status']) || 'pending',
    progress: total > 0 ? Math.round(((found) / total) * 100) : 0,
    startedAt: raw.started_at as string,
    completedAt: raw.finished_at as string,
    errorMessage: raw.error_message as string,
  };
}

export const adminApi = {
  getSyncTasks: async () => {
    const res = await api.get<{ tasks: Record<string, unknown>[]; total: number; page: number; size: number }>('/admin/sync-tasks');
    return (res.tasks || []).map(mapSyncTask);
  },

  triggerSync: async (type: 'full' | 'incremental') => {
    const res = await api.post<Record<string, unknown>>('/admin/sync-tasks', { type });
    return mapSyncTask(res);
  },

  getPendingReviews: () =>
    api.get<Skill[]>('/admin/pending-review'),

  approveSkill: (id: number) =>
    api.put<void>(`/admin/skills/${id}/approve`),

  rejectSkill: (id: number, reason?: string) =>
    api.put<void>(`/admin/skills/${id}/reject`, { reason }),

  getSystemLogs: (lines?: number) =>
    api.get<{ timestamp: string; level: string; message: string }[]>('/admin/logs', { lines }),

  getUsers: async (params?: { page?: number; pageSize?: number }) => {
    const res = await api.get<{
      users: Record<string, unknown>[];
      total: number;
      page: number;
      size: number;
    }>('/admin/users', params as Record<string, string | number | boolean | undefined>);
    const users = (res.users || []).map(u => ({
      id: u.id as number,
      username: u.username as string,
      email: u.email as string,
      role: u.role as string,
      createdAt: u.created_at as string,
    }));
    return { data: users, total: res.total, page: res.page, pageSize: res.size };
  },
};
