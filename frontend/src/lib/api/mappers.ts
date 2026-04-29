import type { Skill, Stats } from '@/src/types';

const iconMap: Record<string, string> = {
  'python': 'code',
  'javascript': 'code',
  'typescript': 'code',
  'go': 'code',
  'rust': 'code',
  'shell': 'terminal',
  'sql': 'database',
};

export function mapStats(raw: Record<string, unknown>): Stats {
  return {
    totalSkills: (raw.total_skills as number) || 0,
    monthlyActiveDevs: (raw.active_skills as number) || 0,
    totalApiCalls: 0,
    pluginInstalls: (raw.total_installs as number) || 0,
    todayNew: (raw.today_new as number) || 0,
    api24hCalls: '0',
    crawlerRunning: 0,
  };
}

export function mapSkill(raw: Record<string, unknown>): Skill {
  return {
    id: raw.id as number,
    title: (raw.display_name || raw.name) as string,
    description: (raw.description || '') as string,
    zhDescription: (raw.zh_description || '') as string,
    enDescription: (raw.en_description || '') as string,
    author: (raw.author || '') as string,
    icon: iconMap[raw.language as string] || 'extension',
    iconColor: '',
    iconBg: '',
    tags: (raw.topics || []) as string[],
    category: (raw.category || '') as string,
    version: (raw.version || '') as string,
    rating: raw.score as number || 0,
    downloads: raw.installs as number || 0,
    installCount: raw.installs as number || 0,
    matchScore: raw.match_score as number | undefined,
    source: raw.is_official ? 'official' : 'github',
    safe: raw.scan_passed as boolean || false,
    createdAt: raw.created_at as string || '',
    updatedAt: raw.updated_at as string || '',
  };
}

export function mapSkillList(raw: Record<string, unknown>[]): Skill[] {
  return raw.map(mapSkill);
}
