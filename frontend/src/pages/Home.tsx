import { Link, useNavigate } from "react-router-dom";
import { useState, useEffect, type FormEvent } from "react";
import { SkillCard } from "../components/ui/SkillCard";
import { CategoryTree } from "../components/ui/CategoryTree";
import { ErrorBanner } from "../components/ui/ErrorBanner";
import { skillsApi } from "../lib/api/skills";
import { statsApi } from "../lib/api/stats";
import { useLanguage } from "../stores/LanguageContext";
import type { Skill, Stats, Category } from "../types";

function pickDesc(lang: string, skill: Skill): string {
  if (lang === 'zh' && skill.zhDescription) return skill.zhDescription;
  if (lang === 'en' && skill.enDescription) return skill.enDescription;
  return skill.description;
}

export function Home() {
  const navigate = useNavigate();
  const { language, t } = useLanguage();
  const [searchMode, setSearchMode] = useState<'keyword' | 'semantic'>('keyword');
  const [searchQuery, setSearchQuery] = useState("");
  const [stats, setStats] = useState<Stats | null>(null);
  const [categories, setCategories] = useState<Category[]>([]);
  const [trendingSkills, setTrendingSkills] = useState<Skill[]>([]);
  const [latestSkills, setLatestSkills] = useState<Skill[]>([]);
  const [topError, setTopError] = useState<string | null>(null);

  useEffect(() => {
    statsApi.getOverview().then(setStats).catch(() => setTopError(t('加载数据失败', 'Failed to load data')));
    skillsApi.getCategories().then(setCategories).catch(() => {});
    skillsApi.getTrending().then(setTrendingSkills).catch(() => {});
    skillsApi.getLatest().then(setLatestSkills).catch(() => {});
  }, []);

  const handleSearch = (e: FormEvent) => {
    e.preventDefault();
    if (!searchQuery.trim()) return;
    const params = new URLSearchParams({ q: searchQuery.trim() });
    if (searchMode === 'semantic') {
      params.set('mode', 'semantic');
    }
    navigate(`/search?${params.toString()}`);
  };

  const semanticCards = [
    { label: t('可用技能总量', 'Total Skills'), value: stats?.totalSkills?.toLocaleString() ?? '-', change: '12%', icon: 'extension', color: 'text-brand-600', bg: 'bg-brand-50' },
    { label: t('月活跃开发者', 'Monthly Active Devs'), value: stats ? (stats.monthlyActiveDevs / 1000).toFixed(1) + 'k+' : '-', change: '8%', icon: 'group', color: 'text-brand-600', bg: 'bg-brand-50' },
    { label: t('累计 API 调用', 'Total API Calls'), value: stats ? (stats.totalApiCalls / 1e9).toFixed(1) + 'B' : '-', changeText: t('历史总计', 'All Time'), icon: 'api', color: 'text-brand-600', bg: 'bg-brand-50' },
    { label: t('VS Code 插件安装', 'VS Code Plugin Installs'), value: stats ? (stats.pluginInstalls / 1000).toFixed(0) + 'k+' : '-', change: '24%', icon: 'download', color: 'text-brand-600', bg: 'bg-brand-50' },
  ];

  return (
    <div className="flex w-full">
      <aside className="hidden md:flex flex-col bg-slate-50 w-56 lg:w-64 border-r border-slate-200 overflow-y-auto">
        <div className="px-4 lg:px-6 mb-4 pt-8">
          <h3 className="text-slate-900 font-bold text-sm">{t('探索分类', 'Explore Categories')}</h3>
          <p className="text-slate-500 text-xs mt-1">{t('按领域浏览', 'Browse by domain')}</p>
        </div>
        <nav className="flex flex-col w-full text-sm">
          <CategoryTree categories={categories} linkMode language={language} />
        </nav>
      </aside>

      <main className="flex-1 px-4 md:px-8 py-8 flex flex-col gap-10 max-w-[1280px] mx-auto w-full">
        <ErrorBanner message={topError} onDismiss={() => setTopError(null)} />
        <section className="flex flex-col items-center justify-center text-center py-12 lg:py-16 px-4 md:px-8 bg-white rounded-xl border border-slate-200 shadow-card relative overflow-hidden">
          <div className="absolute inset-0 bg-gradient-to-br from-brand-50/50 to-white pointer-events-none" />
          <div className="relative z-10 max-w-3xl flex flex-col items-center gap-5">
            <div className="inline-flex items-center gap-2 px-3 py-1 bg-brand-50 text-brand-600 rounded-full text-sm font-medium border border-brand-200/50">
              <span className="material-symbols-outlined text-sm fill">new_releases</span>
              {t('v2.4 发布: 支持本地模型接入', 'v2.4 Release: Local Model Support')}
            </div>
            <h1 className="text-2xl md:text-4xl font-bold text-slate-900 tracking-tight leading-tight">
              {t('让 AI 能力的获取像安装 IDE 插件一样简单', 'Get AI capabilities as easily as installing IDE plugins')}
            </h1>
            <p className="text-sm md:text-base text-slate-500 max-w-2xl">
              {t('探索、测试并一键集成数以千计的专业 AI 技能。专为开发者打造的高效智能体生态系统。', 'Explore, test, and integrate thousands of professional AI skills with one click. An efficient agent ecosystem built for developers.')}
            </p>

            <div className="w-full max-w-xl mt-2">
              <div className="flex items-center gap-1 mb-3">
                <button
                  onClick={() => setSearchMode('keyword')}
                  aria-pressed={searchMode === 'keyword'}
                  className={`px-3 py-1.5 text-xs font-medium rounded-full transition-all ${
                    searchMode === 'keyword'
                      ? 'bg-brand-600 text-white'
                      : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                  }`}
                >
                  {t('关键词搜索', 'Keyword Search')}
                </button>
                <button
                  onClick={() => setSearchMode('semantic')}
                  aria-pressed={searchMode === 'semantic'}
                  className={`px-3 py-1.5 text-xs font-medium rounded-full transition-all ${
                    searchMode === 'semantic'
                      ? 'bg-brand-600 text-white'
                      : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                  }`}
                >
                  {t('智能语义搜索', 'Semantic Search')}
                </button>
              </div>

              <form onSubmit={handleSearch} className="relative flex items-center shadow-sm hover:shadow-md transition-shadow duration-300 rounded-lg bg-white border border-slate-200 focus-within:border-brand-600 focus-within:ring-2 focus-within:ring-brand-50">
                <span className="material-symbols-outlined absolute left-4 text-slate-400">
                  {searchMode === 'semantic' ? 'psychology' : 'search'}
                </span>
                <input
                  type="text"
                  value={searchQuery}
                  onChange={e => setSearchQuery(e.target.value)}
                  className="w-full py-3.5 md:py-4 pl-12 pr-32 bg-transparent border-none text-slate-900 text-sm focus:ring-0 outline-none placeholder:text-slate-500"
                  placeholder={searchMode === 'semantic' ? t("用自然语言描述您需要的技能...", "Describe the skill you need in natural language...") : t("搜索技能，如 'Python 代码重构'", "Search skills, e.g. 'Python code refactoring'")}
                />
                <button
                  type="submit"
                  disabled={!searchQuery.trim()}
                  className="absolute right-2 px-4 py-1.5 md:py-2 bg-brand-600 text-white text-sm font-medium rounded-md hover:bg-brand-700 transition-colors disabled:opacity-50"
                >
                  {searchMode === 'semantic' ? t('智能匹配', 'Smart Match') : t('搜索', 'Search')}
                </button>
              </form>
            </div>
          </div>
        </section>

        <section className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {semanticCards.map((stat, idx) => (
            <div key={idx} className="card p-5 md:p-6">
              <div className="flex items-center justify-between mb-3 md:mb-4">
                <span className="text-xs md:text-sm text-slate-500 font-medium">{stat.label}</span>
                <div className={`p-2 rounded-md ${stat.bg} ${stat.color}`}>
                  <span className="material-symbols-outlined text-[20px]">{stat.icon}</span>
                </div>
              </div>
              <div className="flex items-baseline gap-2">
                <span className="text-xl md:text-2xl font-bold text-slate-900">{stat.value}</span>
                {'change' in stat ? (
                  <span className="text-xs text-green-500 flex items-center">
                    <span className="material-symbols-outlined text-[14px]">arrow_upward</span> {stat.change}
                  </span>
                ) : (
                  <span className="text-xs text-slate-500">{stat.changeText}</span>
                )}
              </div>
            </div>
          ))}
        </section>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <section className="lg:col-span-2 flex flex-col gap-5">
            <div className="flex items-center justify-between">
              <h2 className="text-lg md:text-xl font-bold text-slate-900 flex items-center gap-2">
                <span className="material-symbols-outlined text-amber-500 fill">local_fire_department</span>
                {t('热门技能榜', 'Trending Skills')}
              </h2>
              <Link to="/search?sort=downloads" className="text-sm text-brand-600 hover:text-brand-700 font-medium flex items-center">
                {t('查看全部', 'View All')}
                <span className="material-symbols-outlined text-[16px]">chevron_right</span>
              </Link>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              {trendingSkills.length > 0 ? trendingSkills.map(skill => (
                <SkillCard key={skill.id} skill={skill} />
                )) : (
                  <div className="col-span-full text-center text-sm text-slate-500 py-8">{t('暂无热门技能', 'No trending skills available')}</div>
                )}
              </div>
          </section>

          <section className="flex flex-col gap-5">
              <h2 className="text-lg md:text-xl font-bold text-slate-900 flex items-center gap-2">
                <span className="material-symbols-outlined text-green-500">update</span>
                {t('最新收录', 'Latest Additions')}
              </h2>
            <div className="card p-5">
              <div className="relative border-l-2 border-slate-100 ml-3 flex flex-col gap-5 py-2">
                {latestSkills.length > 0 ? latestSkills.map((item, idx) => (
                  <div key={idx} className="relative pl-6">
                    <div className="absolute w-3 h-3 rounded-full -left-[7px] top-1.5 ring-4 ring-white bg-slate-300" />
                    <div className="flex flex-col gap-0.5">
                      <div className="flex items-baseline justify-between gap-2">
                        <Link to={`/skill/${item.id}`} className="text-sm font-medium text-slate-900 hover:text-brand-600 transition-colors">
                          {item.title}
                        </Link>
                        <span className="text-[11px] text-slate-500 shrink-0">{item.createdAt}</span>
                      </div>
                      <p className="text-xs text-slate-500 line-clamp-1">{pickDesc(language, item)}</p>
                    </div>
                  </div>
                  )) : (
                    <div className="text-center text-sm text-slate-500 py-8">{t('暂无最新技能', 'No new skills available')}</div>
                  )}
                </div>
                <Link to="/search?sort=newest" className="block text-center w-full mt-4 py-2 text-brand-600 text-sm font-medium bg-brand-50 hover:bg-brand-100 rounded transition-colors">
                  {t('查看动态流', 'View Activity Feed')}
                </Link>
            </div>
          </section>
        </div>

        <section className="flex flex-col gap-5 pt-8 border-t border-slate-200">
          <h2 className="text-xl md:text-2xl font-bold text-slate-900">{t('按领域探索', 'Explore by Domain')}</h2>
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4">
            {categories.map(cat => {
            const catName = language === 'en' && cat.enName ? cat.enName : (language === 'zh' && cat.zhName ? cat.zhName : cat.name);
            return (
              <Link
                key={cat.id}
                to={`/search?category=${cat.slug}`}
                className="flex flex-col items-center justify-center p-5 md:p-6 bg-white border border-slate-200 rounded-lg hover:border-brand-600 hover:shadow-md group transition-all"
              >
                <span className="material-symbols-outlined text-[28px] md:text-[32px] text-slate-400 group-hover:text-brand-600 mb-2 md:mb-3 transition-colors">{cat.icon}</span>
                <span className="text-xs md:text-sm font-medium text-slate-900 group-hover:text-brand-600 transition-colors">{catName}</span>
              </Link>
            );
          })}
          </div>
        </section>
      </main>
    </div>
  );
}
