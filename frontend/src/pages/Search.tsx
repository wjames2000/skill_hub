import { useSearchParams, Link, useNavigate } from "react-router-dom";
import { useState, useEffect, type FormEvent } from "react";
import { skillsApi } from "../lib/api/skills";
import { routerApi } from "../lib/api/router";
import { SkillCard } from "../components/ui/SkillCard";
import { Pagination } from "../components/ui/Pagination";
import { CategoryTree } from "../components/ui/CategoryTree";
import { useLanguage } from "../stores/LanguageContext";
import type { Skill, SearchFilters, RouterMatchResult, Category } from "../types";

function pickDesc(lang: string, s: { zhDescription: string; enDescription: string; description: string }): string {
  if (lang === 'zh' && s.zhDescription) return s.zhDescription;
  if (lang === 'en' && s.enDescription) return s.enDescription;
  return s.description;
}

export function Search() {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const { language, t } = useLanguage();
  const query = searchParams.get('q') || '';
  const mode = searchParams.get('mode') || 'keyword';

  const [searchInput, setSearchInput] = useState(query);
  const [skills, setSkills] = useState<Skill[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [category, setCategory] = useState('');
  const [tags, setTags] = useState('');
  const [safeOnly, setSafeOnly] = useState(true);
  const [sort, setSort] = useState<'relevance' | 'rating' | 'downloads'>('relevance');
  const [semanticResults, setSemanticResults] = useState<RouterMatchResult[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [pageSize, setPageSize] = useState(10);


  const handleSearch = (e: FormEvent) => {
    e.preventDefault();
    if (!searchInput.trim()) return;
    const params: Record<string, string> = { q: searchInput.trim() };
    if (mode === 'semantic') params.mode = 'semantic';
    navigate(`/search?${new URLSearchParams(params).toString()}`);
  };

  useEffect(() => {
    setSearchInput(query);
  }, [query]);

  useEffect(() => {
    skillsApi.getCategories().then(setCategories).catch(() => {});
  }, []);

  useEffect(() => {
    setLoading(true);
    setLoadError(null);

    if (!query) {
      const listParams: { sort: string; page: number; pageSize: number; safe: boolean; category?: string; tags?: string[] } = { sort: 'installs', page, pageSize, safe: safeOnly };
      if (category) listParams.category = category;
      if (tags.trim()) {
        listParams.tags = tags.split(',').map(t => t.trim()).filter(Boolean);
      }
      skillsApi.list(listParams)
        .then(res => {
          setSkills(res.data);
          setTotal(res.total);
        })
        .catch(() => setLoadError(t('加载失败，请重试', 'Failed to load, please retry')))
        .finally(() => setLoading(false));
      return;
    }

    if (mode === 'semantic') {
      routerApi.match({ query, topK: 6 })
        .then(results => {
          setSemanticResults(results);
          setSkills(results.map(r => r.skill));
          setTotal(results.length);
        })
        .catch(() => {
          setSemanticResults([]);
          setLoadError(t('搜索失败，请重试', 'Search failed, please retry'));
        })
        .finally(() => setLoading(false));
    } else {
      const filters: SearchFilters = { query, page, pageSize };
      if (category) filters.category = category;
      filters.safe = safeOnly;
      if (sort) filters.sort = sort;
      if (tags.trim()) {
        filters.tags = tags.split(',').map(t => t.trim()).filter(Boolean);
      }

      skillsApi.search(filters)
        .then(res => {
          setSkills(res.data);
          setTotal(res.total);
        })
        .catch(() => setLoadError(t('搜索失败，请重试', 'Search failed, please retry')))
        .finally(() => setLoading(false));
    }
  }, [query, page, category, tags, safeOnly, sort, mode, pageSize]);

  const handleFilterChange = (key: string, value: string | boolean) => {
    setPage(1);
    if (key === 'category') setCategory(value as string);
    else if (key === 'safeOnly') setSafeOnly(value as boolean);
    else if (key === 'sort') setSort(value as 'relevance' | 'rating' | 'downloads');
  };

  return (
    <div className="w-full max-w-[1440px] mx-auto px-4 md:px-6 py-6 md:py-8 flex flex-col gap-6">

      <div className="flex flex-col md:flex-row gap-6 md:gap-8 items-start">
      <aside className="w-full md:w-[240px] lg:w-[260px] flex-shrink-0 flex flex-col gap-6 md:sticky md:top-20">
        <div className="card p-4 flex flex-col gap-5">
          <div className="flex items-center justify-between border-b border-slate-100 pb-2">
            <h2 className="font-semibold text-slate-900">{t('筛选', 'Filters')}</h2>
            <button onClick={() => { setCategory(''); setTags(''); setSafeOnly(true); setSort('relevance'); }} className="text-brand-600 text-sm hover:text-brand-700 transition-colors">
              {t('重置', 'Reset')}
            </button>
          </div>

          <div className="flex flex-col gap-3">
            <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider">{t('分类', 'Category')}</h3>
            <div className="flex flex-col text-sm">
              <label className="flex items-center gap-2 cursor-pointer group px-3 py-2 hover:bg-slate-50 rounded">
                <input
                  type="radio"
                  name="category"
                  checked={category === ''}
                  onChange={() => handleFilterChange('category', '')}
                  className="w-4 h-4 border-slate-300 text-brand-600 focus:ring-brand-500"
                />
                <span className="text-slate-700 group-hover:text-brand-600">{t('全部', 'All')}</span>
              </label>
              <CategoryTree
                categories={categories}
                selected={category}
                onSelect={(cat) => handleFilterChange('category', cat.slug)}
                language={language}
              />
            </div>
          </div>

          <div className="flex flex-col gap-3">
            <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider">{t('标签', 'Tags')}</h3>
            <input
              type="text"
              value={tags}
              onChange={e => { setTags(e.target.value); setPage(1); }}
              placeholder={t("用逗号分隔多个标签", "Use commas to separate tags")}
               className="w-full px-3 py-2 text-sm border border-slate-200 rounded-md focus:border-brand-600 focus:ring-2 focus:ring-brand-50 outline-none placeholder:text-slate-500"
            />
          </div>

          <div className="flex flex-col gap-3">
            <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider">{t('安全', 'Safety')}</h3>
            <label className="flex items-center gap-2 cursor-pointer group">
              <input
                type="checkbox"
                checked={safeOnly}
                onChange={e => handleFilterChange('safeOnly', e.target.checked)}
                className="w-4 h-4 rounded border-slate-300 text-green-500 focus:ring-green-500"
              />
              <span className="material-symbols-outlined text-[16px] text-green-500">verified_user</span>
              <span className="text-slate-700 group-hover:text-brand-600">{t('仅显示安全', 'Safe Only')}</span>
            </label>
          </div>

          <div className="flex flex-col gap-3">
            <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider">{t('排序', 'Sort')}</h3>
            <div className="flex flex-col gap-1 bg-slate-50 rounded p-1 text-sm border border-slate-200">
              {[
                { value: 'relevance', label: t('匹配度最高', 'Best Match') },
                { value: 'rating', label: t('评分最高', 'Top Rated') },
                { value: 'downloads', label: t('下载量最多', 'Most Downloaded') },
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
        <form onSubmit={handleSearch} className="relative flex items-center shadow-sm rounded-lg bg-white border border-slate-200 focus-within:border-brand-600 focus-within:ring-2 focus-within:ring-brand-50 w-full shrink-0">
          <span className="material-symbols-outlined absolute left-4 text-slate-400">search</span>
          <input
            type="text"
            value={searchInput}
            onChange={e => setSearchInput(e.target.value)}
            className="w-full py-3 pl-12 pr-32 bg-transparent border-none text-slate-900 text-sm focus:ring-0 outline-none rounded-lg placeholder:text-slate-500"
            placeholder={mode === 'semantic' ? t("用自然语言描述您需要的技能...", "Describe the skill you need in natural language...") : t("搜索技能，如 'Python 代码重构'", "Search skills, e.g. 'Python code refactoring'")}
          />
          <button
            type="submit"
            disabled={!searchInput.trim()}
            className="absolute right-2 px-4 py-1.5 bg-brand-600 text-white text-sm font-medium rounded-md hover:bg-brand-700 transition-colors disabled:opacity-50"
          >
            {mode === 'semantic' ? t('智能匹配', 'Smart Match') : t('搜索', 'Search')}
          </button>
        </form>
        {query && (
          <div className="flex flex-col gap-3">
            <h1 className="text-xl md:text-2xl font-bold text-slate-900 flex items-center gap-2 flex-wrap">
              {mode === 'semantic' ? t('智能匹配结果', 'Semantic Search Results') : (
                <>{t('为您找到', 'Found')} <span className="text-brand-600">{total}</span> {t('个关于', ' results for')}
                <span className="font-mono bg-brand-50 text-brand-700 px-2 py-0.5 rounded border border-brand-100">"{query}"</span>
                {t('的技能', '')}</>
              )}
            </h1>
            <div className="flex items-center gap-2 flex-wrap">
              <span className="text-xs text-slate-400">{t('已选筛选器:', 'Active filters:')}</span>
              {category && (
                <span className="inline-flex items-center gap-1 bg-white text-slate-700 text-xs px-2.5 py-1 rounded-full border border-slate-200">
                  {t('分类', 'Category')}: {category}
                  <button onClick={() => handleFilterChange('category', '')} aria-label={t('移除筛选', 'Remove filter')}>
                    <span className="material-symbols-outlined text-[14px] cursor-pointer hover:text-red-500">close</span>
                  </button>
                </span>
              )}
              {tags.trim() && (
                <span className="inline-flex items-center gap-1 bg-white text-slate-700 text-xs px-2.5 py-1 rounded-full border border-slate-200">
                  {t('标签', 'Tags')}: {tags}
                  <button onClick={() => setTags('')} aria-label={t('移除筛选', 'Remove filter')}>
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
              <h3 className="text-base font-semibold text-brand-900 mb-1">{t('智能语义分析', 'Smart Semantic Analysis')}</h3>
              <p className="text-sm text-brand-800/80 leading-relaxed">
                {t('基于您的自然语言描述，系统通过向量检索和语义匹配为您推荐以下最相关的 AI 技能。', 'Based on your natural language description, the system recommends the most relevant AI skills through vector search and semantic matching.')}
              </p>
            </div>
          </div>
        )}

        <div role="status" aria-live="polite">
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
        ) : loadError ? (
          <div className="card p-12 text-center">
            <span className="material-symbols-outlined text-[48px] text-slate-300">cloud_off</span>
            <h3 className="text-lg font-semibold text-slate-900 mt-4">{loadError}</h3>
            <button
              onClick={() => setLoadError(null)}
              className="btn-secondary text-sm mt-4"
            >
              {t('关闭', 'Dismiss')}
            </button>
          </div>
        ) : skills.length === 0 ? (
          <div className="card p-12 text-center">
            <span className="material-symbols-outlined text-[48px] text-slate-300">{query ? 'search_off' : 'folder_off'}</span>
            <h3 className="text-lg font-semibold text-slate-900 mt-4">
              {query
                ? t('未找到相关技能', 'No matching skills found')
                : t('该分类下暂无可用技能', 'No skills available in this category')
              }
            </h3>
            <p className="text-sm text-slate-500 mt-1">
              {query
                ? t('请尝试调整搜索关键词或筛选条件', 'Try adjusting your search keywords or filters')
                : t('请尝试切换其他分类', 'Try switching to another category')
              }
            </p>
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
                        {skill.safe ? t('已扫描安全', 'Safety Verified') : t('存在潜在风险', 'Potential Risk')}
                      </div>
                    </div>
                    <p className="text-sm text-slate-600 line-clamp-2 mt-2 leading-relaxed">{pickDesc(language, skill)}</p>
                  </div>
                  <div className="flex md:flex-col items-center justify-between gap-3 md:border-l border-slate-100 md:pl-4 md:min-w-[140px]">
                    {skill.matchScore != null && (
                      <div className="w-full hidden md:block">
                        <div className="flex justify-between items-center mb-1">
                          <span className="text-xs text-slate-400">{t('匹配度', 'Match')}</span>
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
                      {t('安装', 'Install')}
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
              onPageSizeChange={(size) => { setPageSize(size); setPage(1); }}
            />
          </>
        )}
      </div>
      </section>
      </div>
    </div>
  );
}
