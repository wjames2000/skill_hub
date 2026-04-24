import { api } from './client';
import type { RouterMatchRequest, RouterMatchResult, RouterExecuteRequest, RouterExecuteResult } from '@/src/types';

export const routerApi = {
  match: (data: RouterMatchRequest) =>
    api.post<RouterMatchResult[]>('/router/match', data),

  execute: (data: RouterExecuteRequest) =>
    api.post<RouterExecuteResult>('/router/execute', data),
};
