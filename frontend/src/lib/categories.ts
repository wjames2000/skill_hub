import type { Language } from '../stores/LanguageContext';

const categoryNames: Record<string, { zh: string; en: string }> = {
  'content-creation': { zh: '内容创作类', en: 'Content Creation' },
  'copywriting': { zh: '文案撰写', en: 'Copywriting' },
  'visual-design': { zh: '视觉设计', en: 'Visual Design' },
  'multimedia-production': { zh: '多媒体制作', en: 'Multimedia Production' },
  'information-processing': { zh: '信息处理类', en: 'Information Processing' },
  'data-collection': { zh: '数据采集', en: 'Data Collection' },
  'data-analysis': { zh: '数据分析', en: 'Data Analysis' },
  'data-cleaning': { zh: '数据清洗', en: 'Data Cleaning' },
  'document-processing': { zh: '文档处理', en: 'Document Processing' },
  'translation': { zh: '翻译', en: 'Translation' },
  'software-engineering': { zh: '软件工程类', en: 'Software Engineering' },
  'code-quality': { zh: '代码质量', en: 'Code Quality' },
  'pr-review': { zh: 'PR 检查', en: 'PR Review' },
  'security-scan': { zh: '安全扫描', en: 'Security Scan' },
  'lint-check': { zh: '规范校验', en: 'Lint Check' },
  'templates': { zh: '模板脚手架', en: 'Templates & Scaffolding' },
  'project-init': { zh: '项目初始化', en: 'Project Init' },
  'doc-generation': { zh: '文档生成', en: 'Docs Generation' },
  'code-generation': { zh: '代码生成', en: 'Code Generation' },
  'deployment': { zh: '部署运维', en: 'Deployment & Operations' },
  'release-check': { zh: '上线检查', en: 'Release Check' },
  'rollback': { zh: '回滚管理', en: 'Rollback Management' },
  'env-verification': { zh: '环境验证', en: 'Env Verification' },
  'team-collaboration': { zh: '团队协作类', en: 'Team Collaboration' },
  'process-automation': { zh: '流程自动化', en: 'Process Automation' },
  'requirement-review': { zh: '需求评审', en: 'Requirement Review' },
  'release-review': { zh: '发版检查', en: 'Release Review' },
  'weekly-report': { zh: '周报流程', en: 'Weekly Report' },
  'knowledge-management': { zh: '知识管理', en: 'Knowledge Management' },
  'knowledge-base': { zh: '内部知识库', en: 'Knowledge Base' },
  'sdk-guide': { zh: 'SDK 指南', en: 'SDK Guide' },
  'complaint-sop': { zh: '客诉 SOP', en: 'Complaint SOP' },
  'project-management': { zh: '项目管理', en: 'Project Management' },
  'goal-tracking': { zh: '目标追踪', en: 'Goal Tracking' },
  'log-generation': { zh: '日志生成', en: 'Log Generation' },
  'infrastructure': { zh: '基础设施类', en: 'Infrastructure' },
  'resource-inspection': { zh: '资源巡检', en: 'Resource Inspection' },
  'cluster-management': { zh: '集群管理', en: 'Cluster Management' },
  'troubleshooting': { zh: '故障排查', en: 'Troubleshooting' },
  'env-repair': { zh: '环境修复', en: 'Env Repair' },
  'ai-agent': { zh: 'AI 智能体类', en: 'AI Agents' },
  'atomic-skill': { zh: '原子型技能', en: 'Atomic Skills' },
  'workflow-skill': { zh: '工作流型技能', en: 'Workflow Skills' },
  'dedicated-skill': { zh: '专属型技能', en: 'Dedicated Skills' },
  'reference-materials': { zh: '参考型资料', en: 'Reference Materials' },
  'rules': { zh: '规则说明', en: 'Rules' },
  'sdk-usage': { zh: 'SDK 使用方法', en: 'SDK Usage' },
};

export function getCategoryName(slug: string, language: Language): string {
  // Handle path format: "大类 > 子类" — split, look up each segment, rejoin
  if (slug.includes(' > ')) {
    return slug.split(' > ').map(s => {
      const trimmed = s.trim();
      const cat = categoryNames[trimmed];
      return cat ? (language === 'zh' ? cat.zh : cat.en) : trimmed;
    }).join(' > ');
  }
  const cat = categoryNames[slug];
  if (!cat) return slug;
  return language === 'zh' ? cat.zh : cat.en;
}
