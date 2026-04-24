import { useParams, Link } from "react-router-dom";
import { useState, useEffect } from "react";
import { skillsApi } from "../lib/api/skills";
import { ReviewSection } from "../components/ui/ReviewSection";
import { StarRating } from "../components/ui/StarRating";
import type { SkillDetail } from "../types";

const mockDetail: SkillDetail = {
  id: 1,
  title: 'Excel 趋势分析器',
  description: '自动分析 Excel 数据集，识别季节性趋势、异常值，并生成执行摘要报告。完美集成至您的数据科学工作流。',
  author: 'Anthropic',
  icon: 'finance',
  iconColor: 'text-teal-600',
  iconBg: 'bg-teal-50',
  tags: ['数据分析', 'Excel', '可视化'],
  category: '数据分析',
  version: 'v1.2.0',
  rating: 4.8,
  downloads: 12500,
  installCount: 12500,
  source: 'official',
  safe: true,
  createdAt: '2024-01-15',
  updatedAt: '2024-03-20',
  readme: `# Excel Trend Analyzer

A powerful skill to automatically analyze trends in large Excel datasets using advanced LLM reasoning.

## Installation

\`skillhub install anthropic/excel-trend\`

## Usage Example

\`\`\`typescript
import { SkillHub } from '@skillhub/core';

// Initialize the analyzer
const analyzer = new SkillHub('anthropic/excel-trend');

// Run analysis on sales data
const results = await analyzer.execute({
  file: './Q3_Sales.xlsx',
  targetColumn: 'Revenue',
  detectAnomalies: true
});

console.log(results.insights);
\`\`\`

## Features

- Automatic trend detection
- Anomaly identification
- Executive summary generation
- Multi-sheet support
`,
  files: [
    { name: 'src', path: '/src', type: 'dir' },
    { name: 'src/index.ts', path: '/src/index.ts', type: 'file', size: 2048 },
    { name: 'src/analyzer.ts', path: '/src/analyzer.ts', type: 'file', size: 4096 },
    { name: 'README.md', path: '/README.md', type: 'file', size: 1024 },
    { name: 'package.json', path: '/package.json', type: 'file', size: 512 },
  ],
  similarSkills: [
    { id: 2, title: 'CSV 数据清洗', author: 'Google', icon: 'table_chart', iconColor: 'text-green-600', iconBg: 'bg-green-50', tags: ['数据处理'], category: '数据处理', version: 'v2.1', rating: 4.8, downloads: 8200, installCount: 0, source: 'official', safe: true, createdAt: '', updatedAt: '', description: 'Automatically detects and fixes formatting issues in CSV files.' },
    { id: 3, title: 'SQL 查询生成器', author: 'OpenAI', icon: 'query_stats', iconColor: 'text-purple-600', iconBg: 'bg-purple-50', tags: ['SQL'], category: '数据库', version: 'v1.0', rating: 4.9, downloads: 21000, installCount: 0, source: 'official', safe: true, createdAt: '', updatedAt: '', description: 'Convert natural language questions into optimized SQL.' },
    { id: 4, title: '图表自动生成', author: 'Anthropic', icon: 'insights', iconColor: 'text-orange-600', iconBg: 'bg-orange-50', tags: ['可视化'], category: '可视化', version: 'v0.9', rating: 4.5, downloads: 3400, installCount: 0, source: 'official', safe: true, createdAt: '', updatedAt: '', description: 'Analyzes data arrays and suggests best charts.' },
  ],
  reviews: [
    { id: 1, userId: 101, userName: '数据达人', userAvatar: '', rating: 5, comment: '非常好用的工具，帮我节省了大量数据分析时间。趋势识别非常准确！', createdAt: '3天前' },
    { id: 2, userId: 102, userName: 'Pythonista', userAvatar: '', rating: 4, comment: '功能强大，但处理超大 Excel 文件时略慢，希望能优化性能。', createdAt: '1周前' },
    { id: 3, userId: 103, userName: '分析师小王', userAvatar: '', rating: 5, comment: '生成的报告可以直接用于客户汇报，非常专业。', createdAt: '2周前' },
  ],
};

export function Detail() {
  const { id } = useParams<{ id: string }>();
  const [skill, setSkill] = useState<SkillDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'readme' | 'files' | 'details'>('readme');
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    setLoading(true);
    if (id) {
      skillsApi.getById(Number(id))
        .then(setSkill)
        .catch(() => setSkill(mockDetail))
        .finally(() => setLoading(false));
    }
  }, [id]);

  const handleCopyInstall = () => {
    navigator.clipboard.writeText(`skillhub install ${skill?.author?.toLowerCase() || 'example'}/${(skill?.title || 'skill').toLowerCase().replace(/\s+/g, '-')}`);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  if (loading) {
    return (
      <div className="w-full max-w-[1280px] mx-auto px-6 py-8">
        <div className="flex flex-col gap-6">
          <div className="h-4 w-48 skeleton" />
          <div className="card p-6">
            <div className="flex gap-6">
              <div className="w-24 h-24 rounded-lg skeleton" />
              <div className="flex-1 flex flex-col gap-3">
                <div className="h-8 w-64 skeleton" />
                <div className="h-4 w-48 skeleton" />
                <div className="h-4 w-full skeleton" />
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!skill) {
    return (
      <div className="w-full max-w-[1280px] mx-auto px-6 py-16 text-center">
        <span className="material-symbols-outlined text-[64px] text-slate-300">error</span>
        <h2 className="text-xl font-bold text-slate-900 mt-4">技能未找到</h2>
        <Link to="/" className="text-brand-600 font-medium mt-2 inline-block">返回首页</Link>
      </div>
    );
  }

  return (
    <div className="w-full max-w-[1280px] mx-auto px-4 md:px-6 py-6 pb-24 lg:pb-6 flex flex-col gap-6">
      <nav className="flex text-xs text-slate-500 items-center gap-2 flex-wrap">
        <Link to="/" className="hover:text-brand-600 transition-colors">首页</Link>
        <span className="material-symbols-outlined text-[16px]">chevron_right</span>
        <a href="#" className="hover:text-brand-600 transition-colors">{skill.category}</a>
        <span className="material-symbols-outlined text-[16px]">chevron_right</span>
        <span className="text-slate-900 font-medium truncate max-w-[200px]">{skill.title}</span>
      </nav>

      <div className="card p-4 md:p-6 flex flex-col md:flex-row md:items-start justify-between gap-4 md:gap-6">
        <div className="flex gap-4 md:gap-6">
          <div className="w-16 h-16 md:w-24 md:h-24 rounded-lg flex items-center justify-center shrink-0"
            style={{ backgroundColor: skill.iconBg || '#f8fafc' }}
          >
            <span className="material-symbols-outlined text-[32px] md:text-[48px] fill"
              style={{ color: skill.iconColor || '#2563eb' }}
            >
              {skill.icon}
            </span>
          </div>
          <div className="min-w-0">
            <div className="flex items-center gap-2 md:gap-3 mb-2 flex-wrap">
              <h1 className="text-xl md:text-3xl font-bold text-slate-900">{skill.title}</h1>
              {skill.safe && (
                <span className="badge-green text-[11px]">
                  <span className="material-symbols-outlined text-[14px]">verified_user</span>
                  已通过安全扫描
                </span>
              )}
            </div>

            <div className="flex flex-wrap items-center gap-2 md:gap-4 text-xs md:text-sm text-slate-500">
              <div className="flex items-center gap-1">
                <span className="material-symbols-outlined text-[16px] md:text-[18px]">person</span>
                <span className="font-medium text-slate-900">{skill.author}</span>
              </div>
              <span className="w-1 h-1 rounded-full bg-slate-300" />
              <div className="flex items-center gap-1">
                <span className="material-symbols-outlined text-[16px] md:text-[18px]">tag</span>
                <span>{skill.version}</span>
              </div>
              <span className="w-1 h-1 rounded-full bg-slate-300" />
              <div className="flex items-center gap-1">
                <span className="material-symbols-outlined text-[16px] md:text-[18px]">download</span>
                <span>{(skill.downloads / 1000).toFixed(1)}k 安装</span>
              </div>
              <span className="w-1 h-1 rounded-full bg-slate-300" />
              <StarRating rating={skill.rating} size="sm" />
            </div>

            <p className="mt-3 text-sm md:text-base text-slate-600 max-w-2xl leading-relaxed">
              {skill.description}
            </p>
          </div>
        </div>

        <div className="flex flex-col gap-2 shrink-0 w-full md:w-auto">
          <button className="btn-primary w-full md:w-auto">
            <span className="material-symbols-outlined text-[20px]">code_blocks</span>
            在 VS Code 中安装
          </button>
          <div className="flex gap-2">
            <button onClick={handleCopyInstall} className="btn-secondary flex-1 md:flex-none text-sm">
              <span className="material-symbols-outlined text-[18px]">{copied ? 'check' : 'content_copy'}</span>
              {copied ? '已复制' : '复制命令'}
            </button>
            <button className="btn-secondary p-2.5 flex items-center justify-center">
              <span className="material-symbols-outlined text-[20px]">star</span>
            </button>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div className="lg:col-span-8 xl:col-span-9 flex flex-col gap-4">
          <div className="flex border-b border-slate-200 overflow-x-auto">
            {([
              { key: 'readme', label: 'SKILL.md', icon: 'description' },
              { key: 'files', label: '文件目录', icon: 'folder' },
              { key: 'details', label: '详情', icon: 'info' },
            ] as const).map(tab => (
              <button
                key={tab.key}
                onClick={() => setActiveTab(tab.key)}
                className={`px-3 md:px-4 py-3 border-b-2 font-medium text-sm flex items-center gap-1.5 transition-all whitespace-nowrap ${
                  activeTab === tab.key
                    ? 'border-brand-600 text-brand-600'
                    : 'border-transparent text-slate-500 hover:text-slate-900 hover:bg-slate-50'
                }`}
              >
                <span className="material-symbols-outlined text-[16px] md:text-[18px]">{tab.icon}</span>
                {tab.label}
              </button>
            ))}
          </div>

          {activeTab === 'readme' && (
            <div className="bg-[#1e1e1e] rounded-lg border border-slate-800 p-4 md:p-6 shadow-sm text-slate-300 overflow-x-auto">
              <button onClick={handleCopyInstall} className="absolute top-4 right-4 text-slate-400 hover:text-white bg-slate-800 p-1.5 rounded opacity-0 group-hover:opacity-100 transition-opacity hidden md:flex items-center">
                <span className="material-symbols-outlined text-[16px]">content_copy</span>
              </button>
              {skill.readme.split('\n').map((line, idx) => {
                if (line.startsWith('# ')) return <h1 key={idx} className="text-2xl font-bold text-white mb-4 border-b border-slate-700 pb-2">{line.slice(2)}</h1>;
                if (line.startsWith('## ')) return <h2 key={idx} className="text-xl font-bold text-white mt-6 mb-3">{line.slice(3)}</h2>;
                if (line.startsWith('### ')) return <h3 key={idx} className="text-lg font-bold text-white mt-4 mb-2">{line.slice(4)}</h3>;
                if (line.startsWith('```')) return null;
                if (line.startsWith('- ')) return <li key={idx} className="text-slate-400 ml-4 list-disc">{line.slice(2)}</li>;
                if (line.trim()) return <p key={idx} className="text-slate-400 mb-2">{line}</p>;
                return <div key={idx} className="h-2" />;
              })}
            </div>
          )}

          {activeTab === 'files' && (
            <div className="card p-4">
              <div className="flex flex-col gap-1">
                {skill.files.map(file => (
                  <div key={file.path} className="flex items-center gap-3 px-3 py-2 rounded hover:bg-slate-50 transition-colors cursor-pointer text-sm">
                    <span className="material-symbols-outlined text-[18px] text-slate-400">
                      {file.type === 'dir' ? 'folder' : 'description'}
                    </span>
                    <span className="text-slate-700">{file.name}</span>
                    {file.size && (
                      <span className="text-xs text-slate-400 ml-auto">
                        {file.size >= 1024 ? `${(file.size / 1024).toFixed(1)} KB` : `${file.size} B`}
                      </span>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}

          {activeTab === 'details' && (
            <div className="card p-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="text-xs text-slate-500 uppercase tracking-wider font-semibold">作者</label>
                  <p className="text-slate-900 font-medium mt-1">{skill.author}</p>
                </div>
                <div>
                  <label className="text-xs text-slate-500 uppercase tracking-wider font-semibold">版本</label>
                  <p className="text-slate-900 font-medium mt-1">{skill.version}</p>
                </div>
                <div>
                  <label className="text-xs text-slate-500 uppercase tracking-wider font-semibold">分类</label>
                  <p className="text-slate-900 font-medium mt-1">{skill.category}</p>
                </div>
                <div>
                  <label className="text-xs text-slate-500 uppercase tracking-wider font-semibold">来源</label>
                  <p className="text-slate-900 font-medium mt-1 flex items-center gap-1">
                    {skill.source === 'official' ? (
                      <span className="badge-blue">官方认证</span>
                    ) : (
                      <span className="badge bg-slate-50 text-slate-600 border border-slate-200">社区开源</span>
                    )}
                  </p>
                </div>
                <div>
                  <label className="text-xs text-slate-500 uppercase tracking-wider font-semibold">创建时间</label>
                  <p className="text-slate-900 font-medium mt-1">{skill.createdAt}</p>
                </div>
                <div>
                  <label className="text-xs text-slate-500 uppercase tracking-wider font-semibold">更新时间</label>
                  <p className="text-slate-900 font-medium mt-1">{skill.updatedAt}</p>
                </div>
                <div className="md:col-span-2">
                  <label className="text-xs text-slate-500 uppercase tracking-wider font-semibold">标签</label>
                  <div className="flex gap-2 mt-2 flex-wrap">
                    {skill.tags.map(tag => (
                      <span key={tag} className="badge-blue">{tag}</span>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          )}

          <ReviewSection reviews={skill.reviews} skillId={skill.id} />
        </div>

        <aside className="lg:col-span-4 xl:col-span-3 flex flex-col gap-4">
          <h3 className="text-lg font-bold text-slate-900 flex items-center gap-2">
            <span className="material-symbols-outlined text-brand-600">psychology</span>
            相似技能推荐
          </h3>
          <div className="flex flex-col gap-3">
            {skill.similarSkills.map(s => (
              <Link
                key={s.id}
                to={`/skill/${s.id}`}
                className="card p-3 md:p-4 flex flex-col gap-2"
              >
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-2 min-w-0">
                    <div className={`w-8 h-8 rounded flex items-center justify-center shrink-0 ${s.iconBg || 'bg-slate-50'}`}
                      style={{ color: s.iconColor }}
                    >
                      <span className="material-symbols-outlined text-[18px]">{s.icon}</span>
                    </div>
                    <div className="min-w-0">
                      <h4 className="font-medium text-slate-900 text-sm truncate group-hover:text-brand-600 transition-colors">{s.title}</h4>
                      <p className="text-xs text-slate-500">{s.author}</p>
                    </div>
                  </div>
                  <span className="text-xs text-slate-400 bg-slate-100 px-1.5 py-0.5 rounded font-mono shrink-0">{s.version}</span>
                </div>
                <p className="text-xs text-slate-600 line-clamp-2">{s.description}</p>
                <div className="flex items-center gap-3 text-xs text-slate-500 mt-1">
                  <span className="flex items-center gap-1">
                    <span className="material-symbols-outlined text-[14px]">star</span> {s.rating.toFixed(1)}
                  </span>
                  <span className="flex items-center gap-1">
                    <span className="material-symbols-outlined text-[14px]">download</span> {(s.downloads / 1000).toFixed(1)}k
                  </span>
                </div>
              </Link>
            ))}
          </div>
        </aside>
      </div>
    </div>
  );
}
