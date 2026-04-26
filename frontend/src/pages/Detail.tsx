import { useParams, Link } from "react-router-dom";
import { useState, useEffect } from "react";
import { skillsApi } from "../lib/api/skills";
import { ReviewSection } from "../components/ui/ReviewSection";
import { StarRating } from "../components/ui/StarRating";
import { ErrorBanner } from "../components/ui/ErrorBanner";
import { useLanguage } from "../stores/LanguageContext";
import { MarkdownRenderer } from "../components/ui/MarkdownRenderer";
import { getCategoryName } from "../lib/categories";
import type { SkillDetail } from "../types";

function pickDesc(lang: string, s: { zhDescription: string; enDescription: string; description: string }): string {
  if (lang === 'zh' && s.zhDescription) return s.zhDescription;
  if (lang === 'en' && s.enDescription) return s.enDescription;
  return s.description;
}

export function Detail() {
  const { id } = useParams<{ id: string }>();
  const { language, t } = useLanguage();
  const [skill, setSkill] = useState<SkillDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'readme' | 'files' | 'details'>('readme');
  const [copied, setCopied] = useState(false);
  const [topError, setTopError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    if (id) {
      skillsApi.getById(Number(id))
        .then(setSkill)
        .catch(() => setTopError('加载技能详情失败'))
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
        <h2 className="text-xl font-bold text-slate-900 mt-4">{t('技能未找到', 'Skill Not Found')}</h2>
        <Link to="/" className="text-brand-600 font-medium mt-2 inline-block">{t('返回首页', 'Back to Home')}</Link>
      </div>
    );
  }

  return (
    <div className="w-full max-w-[1280px] mx-auto px-4 md:px-6 py-6 pb-24 lg:pb-6 flex flex-col gap-6">
      <ErrorBanner message={topError} onDismiss={() => setTopError(null)} />
      <nav className="flex text-xs text-slate-500 items-center gap-2 flex-wrap">
        <Link to="/" className="hover:text-brand-600 transition-colors">{t('首页', 'Home')}</Link>
        <span className="material-symbols-outlined text-[16px]">chevron_right</span>
        <a href="#" className="hover:text-brand-600 transition-colors">{getCategoryName(skill.category, language)}</a>
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
                  {t('已通过安全扫描', 'Safety Verified')}
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
                <span>{(skill.downloads / 1000).toFixed(1)}k {t('安装', 'installs')}</span>
              </div>
              <span className="w-1 h-1 rounded-full bg-slate-300" />
              <StarRating rating={skill.rating} size="sm" />
            </div>

            <p className="mt-3 text-sm md:text-base text-slate-600 max-w-2xl leading-relaxed">
              {pickDesc(language, skill)}
            </p>
          </div>
        </div>

        <div className="flex flex-col gap-2 shrink-0 w-full md:w-auto">
          <button className="btn-primary w-full md:w-auto">
            <span className="material-symbols-outlined text-[20px]">code_blocks</span>
            {t('在 VS Code 中安装', 'Install in VS Code')}
          </button>
          <div className="flex gap-2">
            <button onClick={handleCopyInstall} className="btn-secondary flex-1 md:flex-none text-sm" aria-label={t('复制命令', 'Copy command')}>
              <span className="material-symbols-outlined text-[18px]">{copied ? 'check' : 'content_copy'}</span>
              {copied ? t('已复制', 'Copied') : t('复制命令', 'Copy Command')}
            </button>
            <button className="btn-secondary p-2.5 flex items-center justify-center" aria-label={t('收藏', 'Favorite')}>
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
              { key: 'files', label: t('文件目录', 'Files'), icon: 'folder' },
              { key: 'details', label: t('详情', 'Details'), icon: 'info' },
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
            <div className="bg-[#1e1e1e] rounded-lg border border-slate-800 p-4 md:p-6 shadow-sm overflow-x-auto">
              <MarkdownRenderer content={skill.readme} />
            </div>
          )}

          {activeTab === 'files' && (
            <div className="card p-8 text-center">
              <span className="material-symbols-outlined text-[48px] text-slate-300">folder_open</span>
              <p className="text-sm text-slate-500 mt-3">{t('文件目录将在技能同步后显示', 'File listing will be available after skill sync')}</p>
            </div>
          )}

          {activeTab === 'details' && (
            <div className="card p-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="text-xs text-slate-500 uppercase tracking-wider font-semibold">{t('作者', 'Author')}</label>
                  <p className="text-slate-900 font-medium mt-1">{skill.author}</p>
                </div>
                <div>
                  <label className="text-xs text-slate-500 uppercase tracking-wider font-semibold">{t('版本', 'Version')}</label>
                  <p className="text-slate-900 font-medium mt-1">{skill.version}</p>
                </div>
                <div>
                  <label className="text-xs text-slate-500 uppercase tracking-wider font-semibold">{t('分类', 'Category')}</label>
                  <p className="text-slate-900 font-medium mt-1">{getCategoryName(skill.category, language)}</p>
                </div>
                <div>
                  <label className="text-xs text-slate-500 uppercase tracking-wider font-semibold">{t('来源', 'Source')}</label>
                  <p className="text-slate-900 font-medium mt-1 flex items-center gap-1">
                    {skill.source === 'official' ? (
                      <span className="badge-blue">{t('官方认证', 'Official')}</span>
                    ) : (
                      <span className="badge bg-slate-50 text-slate-600 border border-slate-200">{t('社区开源', 'Community')}</span>
                    )}
                  </p>
                </div>
                <div>
                  <label className="text-xs text-slate-500 uppercase tracking-wider font-semibold">{t('创建时间', 'Created')}</label>
                  <p className="text-slate-900 font-medium mt-1">{skill.createdAt}</p>
                </div>
                <div>
                  <label className="text-xs text-slate-500 uppercase tracking-wider font-semibold">{t('更新时间', 'Updated')}</label>
                  <p className="text-slate-900 font-medium mt-1">{skill.updatedAt}</p>
                </div>
                <div className="md:col-span-2">
                  <label className="text-xs text-slate-500 uppercase tracking-wider font-semibold">{t('标签', 'Tags')}</label>
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
            {t('相似技能推荐', 'Similar Skills')}
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
                <p className="text-xs text-slate-600 line-clamp-2">{pickDesc(language, s)}</p>
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
