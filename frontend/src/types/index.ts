export interface Skill {
  id: number;
  title: string;
  description: string;
  zhDescription: string;
  enDescription: string;
  author: string;
  icon: string;
  iconColor: string;
  iconBg: string;
  tags: string[];
  category: string;
  version: string;
  rating: number;
  downloads: number;
  installCount: number;
  matchScore?: number;
  source: 'official' | 'github';
  safe: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface SkillDetail extends Skill {
  readme: string;
  files: SkillFile[];
  similarSkills: Skill[];
  reviews: Review[];
}

export interface SkillFile {
  name: string;
  path: string;
  type: 'file' | 'dir';
  size?: number;
}

export interface Review {
  id: number;
  userId: number;
  userName: string;
  userAvatar: string;
  rating: number;
  comment: string;
  createdAt: string;
}

export interface Category {
  id: number;
  name: string;
  slug: string;
  icon: string;
  count: number;
  children: Category[];
}

export type FlatCategory = Omit<Category, 'children'>;

export interface Stats {
  totalSkills: number;
  monthlyActiveDevs: number;
  totalApiCalls: number;
  pluginInstalls: number;
  todayNew: number;
  api24hCalls: string;
  crawlerRunning: number;
}

export interface User {
  id: number;
  username: string;
  email: string;
  avatar: string;
  bio: string;
  role: 'user' | 'admin';
  createdAt: string;
}

export interface ApiKey {
  id: number;
  name: string;
  key: string;
  createdAt: string;
  lastUsed: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
}

export interface SearchFilters {
  query: string;
  category?: string;
  tags?: string[];
  source?: string;
  safe?: boolean;
  sort?: 'relevance' | 'rating' | 'downloads';
  page?: number;
  pageSize?: number;
}

export interface RouterMatchRequest {
  query: string;
  topK?: number;
}

export interface RouterMatchResult {
  skill: Skill;
  score: number;
  reason: string;
}

export interface RouterExecuteRequest {
  skillId: number;
  params: Record<string, unknown>;
}

export interface RouterExecuteResult {
  output: string;
  duration: number;
  tokenUsage: number;
}

export interface SyncTask {
  id: number;
  name: string;
  status: 'running' | 'completed' | 'failed' | 'pending';
  type: 'full' | 'incremental';
  progress: number;
  startedAt: string;
  completedAt?: string;
  log: string[];
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}
