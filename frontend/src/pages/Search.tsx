import { useSearchParams, Link } from "react-router-dom";
import { useState, useEffect } from "react";
import { skillsApi } from "../lib/api/skills";
import { routerApi } from "../lib/api/router";
import { SkillCard } from "../components/ui/SkillCard";
import { Pagination } from "../components/ui/Pagination";
import type { Skill, SearchFilters, RouterMatchResult } from "../types";

const categories = [
  { value: '', label: '全部' },
  { value: '数据处理', label: '数据处理' },
  { value: '可视化', label: '可视化' },
  { value: '自动化脚本', label: '自动化脚本' },
  { value: '代码生成', label: '代码生成' },
  { value: '前端', label: '前端' },
  { value: 'DevOps', label: 'DevOps' },
  { value: '数据库', label: '数据库' },
];

const mockSkills: Skill[] = [
  { id: 1, title: 'Pandas Excel Master', description: '深度集成 Pandas 库的智能助手，能够通过自然语言指令快速分析庞大的 Excel 文件，自动生成数据清洗脚本和可视化图表代码。', author: 'Microsoft Official', icon: 'table_chart', iconColor: 'text-indigo-600', iconBg: 'bg-indigo-50', tags: ['数据处理', 'Excel'], category: '数据处理', version: 'v2.4.1', rating: 4.9, downloads: 124000, installCount: 0, source: 'official', safe: true, createdAt: '', updatedAt: '', matchScore: 98 },
  { id: 2, title: 'Excel VBA Auto-Coder', description: '专为传统 Excel 用户设计，将自然语言需求直接转化为高质量、带注释的 VBA 宏代码。', author: 'dev_community_x', icon: 'analytics', iconColor: 'text-teal-600', iconBg: 'bg-teal-50', tags: ['VBA', 'Excel'], category: '自动化脚本', version: 'v1.0.8', rating: 4.6, downloads: 56000, installCount: 0, source: 'github', safe: true, createdAt: '', updatedAt: '', matchScore: 85 },
  { id: 3, title: 'JSON Schema 生成器', description: '自动将 JSON 数据结构推断转换为符合标准规范的 JSON Schema。', author: 'DataCraft_AI', icon: 'data_object', iconColor: 'text-blue-600', iconBg: 'bg-blue-50', tags: ['JSON', '数据处理'], category: '数据处理', version: 'v3.1.0', rating: 4.9, downloads: 12400, installCount: 0, source: 'official', safe: true, createdAt: '', updatedAt: '', matchScore: 92 },
  { id: 4, title: 'Python 代码重构', description: '遵循 PEP 8 规范，自动识别代码异味并重构 Python 脚本。', author: 'CodeGuru', icon: 'terminal', iconColor: 'text-green-600', iconBg: 'bg-green-50', tags: ['Python', '重构'], category: '代码生成', version: 'v2.0.0', rating: 4.8, downloads: 9800, installCount: 0, source: 'official', safe: true, createdAt: '', updatedAt: '', matchScore: 78 },
  { id: 5, title: 'SQL 慢查询分析器', description: '解析慢查询日志，提供索引优化建议和查询重写方案。', author: 'DB_Doctor', icon: 'analytics', iconColor: 'text-purple-600', iconBg: 'bg-purple-50', tags: ['SQL', '数据库'], category: '数据库', version: 'v1.5.0', rating: 4.7, downloads: 7600, installCount: 0, source: 'github', safe: false, createdAt: '', updatedAt: '', matchScore: 65 },
  { id: 6, title: 'i18n 多语言自动提取', description: '扫描 React/Vue 组件库，自动提取硬编码文本并生成 i18n 文件。', author: 'FrontEnd_Ninja', icon: 'language', iconColor: 'text-orange-600', iconBg: 'bg-orange-50', tags: ['前端', 'i18n'], category: '前端', version: 'v1.2.0', rating: 4.9, downloads: 5400, installCount: 0, source: 'official', safe: true, createdAt: '', updatedAt: '', matchScore: 71 },
];

export function Search() {
  const [searchParams, setSearchParams] = useSearchParams();
  const query = searchParams.get('q') || '';
  const mode = searchParams.get('mode') || 'keyword';

  const [skills, setSkills] = useState<Skill[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [category, setCategory] = useState('');
  const [safeOnly, setSafeOnly] = useState(true);
  const [sort, setSort] = useState<'relevance' | 'rating' | 'downloads'>('relevance');
  const [semanticResults, setSemanticResults] = useState<RouterMatchResult[]>([]);

  const pageSize = 10;

  useEffect(() => {
    if (!query) return;
    setLoading(true);

    if (mode === 'semantic') {
      routerApi.match({ query, topK: 6 })
        .then(results => {
          setSemanticResults(results);
          setSkills(results.map(r => r.skill));
          setTotal(results.length);
        })
        .catch(() => {
          setSemanticResults([]);
          setSkills(mockSkills);
          setTotal(mockSkills.length);
        })
        .finally(() => setLoading(false));
    } else {
      const filters: SearchFilters = { query, page, pageSize };
      if (category) filters.category = category;
      if (safeOnly) filters.safe = true;
      if (sort) filters.sort = sort;

      skillsApi.search(filters)
        .then(res => {
          setSkills(res.data);
          setTotal(res.total);
        })
        .catch(() => {
          const filtered = mockSkills.filter(s =>
            (!category || s.category === category) &&
            (!safeOnly || s.safe)
          );
          setSkills(filtered);
          setTotal(filtered.length);
        })
        .finally(() => setLoading(false));
    }
  }, [query, page, category, safeOnly, sort, mode]);

  const handleFilterChange = (key: string, value: string | boolean) => {
    setPage(1);
    if (key === 'category') setCategory(value as string);
    else if (key === 'safeOnly') setSafeOnly(value as boolean);
    else if (key === 'sort') setSort(value as 'relevance' | 'rating' | 'downloads');
  };

  return (
    <div className="w-full max-w-[1440px] mx-auto px-4 md:px-6 py-6 md:py-8 flex flex-col md:flex-row gap-6 md:gap-8 items-start">
      <aside className="w-full md:w-[240px] lg:w-[260px] flex-shrink-0 flex flex-col gap-6 md:sticky md:top-20">
        <div className="card p-4 flex flex-col gap-5">
          <div className="flex items-center justify-between border-b border-slate-100 pb-2">
            <h2 className="font-semibold text-slate-900">筛选</h2>
            <button onClick={() => { setCategory(''); setSafeOnly(true); setSort('relevance'); }} className="text-brand-600 text-sm hover:text-brand-700 transition-colors">
              重置
            </button>
          </div>

          <div className="flex flex-col gap-3">
            <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider">分类</h3>
            <div className="flex flex-col gap-1.5 text-sm">
              {categories.map(cat => (
                <label key={cat.value} className="flex items-center gap-2 cursor-pointer group">
                  <input
                    type="radio"
                    name="category"
                    checked={category === cat.value}
                    onChange={() => handleFilterChange('category', cat.value)}
                    className="w-4 h-4 border-slate-300 text-brand-600 focus:ring-brand-500"
                  />
                  <span className="text-slate-700 group-hover:text-brand-600">{cat.label}</span>
                </label>
              ))}
            </div>
          </div>

          <div className="flex flex-col gap-3">
            <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider">安全</h3>
            <label className="flex items-center gap-2 cursor-pointer group">
              <input
                type="checkbox"
                checked={safeOnly}
                onChange={e => handleFilterChange('safeOnly', e.target.checked)}
                className="w-4 h-4 rounded border-slate-300 text-green-500 focus:ring-green-500"
              />
              <span className="material-symbols-outlined text-[16px] text-green-500">verified_user</span>
              <span className="text-slate-700 group-hover:text-brand-600">仅显示安全</span>
            </label>
          </div>

          <div className="flex flex-col gap-3">
            <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider">排序</h3>
            <div className="flex flex-col gap-1 bg-slate-50 rounded p-1 text-sm border border-slate-200">
              {[
                { value: 'relevance', label: '匹配度最高' },
                { value: 'rating', label: '评分最高' },
                { value: 'downloads', label: '下载量最多' },
              ].map(opt => (
                <button
                  key={opt.value}
                  onClick={() => handleFilterChange('sort', opt.value)}
                  className={`py-1.5 px-3 rounded text-left transition-all ${
                    sort === opt.value
                      ? 'bg-white shadow-sm border border-brand-200 text-brand-600 font-medium'
                      : 'text-slate-600 hover:bg-slate-100'
                  }`}
                >
                  {opt.label}
                </button>
              ))}
            </div>
          </div>
        </div>
      </aside>

      <section className="flex-1 min-w-0 flex flex-col gap-5">
        {query && (
          <div className="flex flex-col gap-3">
            <h1 className="text-xl md:text-2xl font-bold text-slate-900 flex items-center gap-2 flex-wrap">
              {mode === 'semantic' ? '智能匹配结果' : (
                <>为您找到 <span className="text-brand-600">{total}</span> 个关于
                <span className="font-mono bg-brand-50 text-brand-700 px-2 py-0.5 rounded border border-brand-100">"{query}"</span>
                的技能</>
              )}
            </h1>
            <div className="flex items-center gap-2 flex-wrap">
              <span className="text-xs text-slate-400">已选筛选器:</span>
              {category && (
                <span className="inline-flex items-center gap-1 bg-white text-slate-700 text-xs px-2.5 py-1 rounded-full border border-slate-200">
                  分类: {category}
                  <button onClick={() => handleFilterChange('category', '')}>
                    <span className="material-symbols-outlined text-[14px] cursor-pointer hover:text-red-500">close</span>
                  </button>
                </span>
              )}
            </div>
          </div>
        )}

        {mode === 'semantic' && semanticResults.length > 0 && (
          <div className="bg-brand-50/50 border border-brand-200 rounded-lg p-4 flex items-start gap-4">
            <div className="w-10 h-10 rounded-full bg-brand-600 text-white flex items-center justify-center shrink-0">
              <span className="material-symbols-outlined fill">lightbulb</span>
            </div>
            <div className="pt-1">
              <h3 className="text-base font-semibold text-brand-900 mb-1">智能语义分析</h3>
              <p className="text-sm text-brand-800/80 leading-relaxed">
                基于您的自然语言描述，系统通过向量检索和语义匹配为您推荐以下最相关的 AI 技能。
              </p>
            </div>
          </div>
        )}

        {loading ? (
          <div className="flex flex-col gap-4">
            {[1, 2, 3].map(i => (
              <div key={i} className="card p-5">
                <div className="flex gap-4">
                  <div className="w-16 h-16 rounded-lg skeleton" />
                  <div className="flex-1 flex flex-col gap-3">
                    <div className="h-5 w-48 skeleton" />
                    <div className="h-3 w-32 skeleton" />
                    <div className="h-4 w-full skeleton" />
                    <div className="h-4 w-3/4 skeleton" />
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : skills.length === 0 ? (
          <div className="card p-12 text-center">
            <span className="material-symbols-outlined text-[48px] text-slate-300">search_off</span>
            <h3 className="text-lg font-semibold text-slate-900 mt-4">未找到相关技能</h3>
            <p className="text-sm text-slate-500 mt-1">请尝试调整搜索关键词或筛选条件</p>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 gap-4">
              {skills.map(skill => (
                <Link
                  key={skill.id}
                  to={`/skill/${skill.id}`}
                  className="card p-4 md:p-5 flex flex-col md:flex-row gap-4 md:gap-6 relative group"
                >
                  <div className="w-12 h-12 md:w-16 md:h-16 rounded-lg flex items-center justify-center shrink-0 group-hover:scale-105 transition-transform"
                    style={{ backgroundColor: skill.iconBg || '#eff6ff' }}
                  >
                    <span className="material-symbols-outlined text-[24px] md:text-[32px] fill"
                      style={{ color: skill.iconColor || '#2563eb' }}
                    >
                      {skill.icon}
                    </span>
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex flex-col md:flex-row md:justify-between md:items-start gap-2">
                      <div>
                        <h3 className="text-lg md:text-xl font-bold text-slate-900 group-hover:text-brand-600 transition-colors">
                          {skill.title}
                        </h3>
                        <div className="flex items-center gap-2 text-xs text-slate-500 mt-1 flex-wrap">
                          <span className="flex items-center gap-1">
                            <span className="material-symbols-outlined text-[14px]">person</span>
                            {skill.author}
                          </span>
                          <span className="w-1 h-1 rounded-full bg-slate-300" />
                          <span>{skill.version}</span>
                          {skill.tags.slice(0, 2).map(t => (
                            <span key={t} className="font-mono bg-slate-100 px-1.5 py-0.5 rounded border border-slate-200">{t}</span>
                          ))}
                        </div>
                      </div>
                      <div className={`hidden md:flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-medium ${
                        skill.safe ? 'bg-green-50 text-green-600 border border-green-200/50' : 'bg-amber-50 text-amber-600 border border-amber-200/50'
                      }`}>
                        <span className="material-symbols-outlined text-[14px]">{skill.safe ? 'shield' : 'warning'}</span>
                        {skill.safe ? '已扫描安全' : '存在潜在风险'}
                      </div>
                    </div>
                    <p className="text-sm text-slate-600 line-clamp-2 mt-2 leading-relaxed">{skill.description}</p>
                  </div>
                  <div className="flex md:flex-col items-center justify-between gap-3 md:border-l border-slate-100 md:pl-4 md:min-w-[140px]">
                    {skill.matchScore != null && (
                      <div className="w-full hidden md:block">
                        <div className="flex justify-between items-center mb-1">
                          <span className="text-xs text-slate-400">匹配度</span>
                          <span className={`text-sm font-bold ${skill.matchScore >= 90 ? 'text-green-500' : 'text-brand-600'}`}>
                            {skill.matchScore}%
                          </span>
                        </div>
                        <div className="w-full h-1.5 bg-slate-100 rounded-full mb-2">
                          <div className={`h-full rounded-full ${skill.matchScore >= 90 ? 'bg-green-500' : 'bg-brand-600'}`}
                            style={{ width: `${skill.matchScore}%` }} />
                        </div>
                      </div>
                    )}
                    <div className="flex items-center gap-3 text-xs text-slate-500">
                      <span className="flex items-center gap-1">
                        <span className="material-symbols-outlined text-[14px]">download</span>
                        {(skill.downloads / 1000).toFixed(0)}k
                      </span>
                      <span className="flex items-center gap-1 text-amber-500">
                        <span className="material-symbols-outlined text-[14px] fill">star</span>
                        {skill.rating.toFixed(1)}
                      </span>
                    </div>
                    <button className="btn-primary text-xs !py-1.5 !px-3 w-full hidden md:flex">
                      <span className="material-symbols-outlined text-[14px]">add_box</span>
                      安装
                    </button>
                  </div>
                </Link>
              ))}
            </div>

            <Pagination
              page={page}
              total={total || skills.length}
              pageSize={pageSize}
              onChange={setPage}
            />
          </>
        )}
      </section>
    </div>
  );
}
