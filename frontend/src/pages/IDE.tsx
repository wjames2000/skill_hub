import { Link } from "react-router-dom";
import { useState, useEffect } from "react";
import { skillsApi } from "../lib/api/skills";
import { useLanguage } from "../stores/LanguageContext";
import type { Skill } from "../types";

function pickDesc(lang: string, s: { zhDescription: string; enDescription: string; description: string }): string {
  if (lang === 'zh' && s.zhDescription) return s.zhDescription;
  if (lang === 'en' && s.enDescription) return s.enDescription;
  return s.description;
}

export function IDE() {
  const { language, t } = useLanguage();
  const [skills, setSkills] = useState<Skill[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    skillsApi.getTrending()
      .then(setSkills)
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);
  return (
    <div className="bg-white text-slate-900 font-sans h-screen flex overflow-hidden w-full">
      {/* Extension Sidebar */}
      <nav className="bg-slate-50 h-screen w-64 border-r border-slate-200 flex flex-col overflow-y-auto shrink-0 z-10">
        <div className="px-4 py-4 border-b border-slate-200">
          <h1 className="text-lg font-black text-slate-900 uppercase tracking-wider block"><Link to="/">SkillHub Pro</Link></h1>
          <p className="text-xs font-medium text-slate-500 mt-1">{t('IDE 扩展', 'IDE Extension')}</p>
        </div>
        <div className="flex-1 py-2 flex flex-col gap-1">
          <button className="flex items-center gap-3 px-4 py-2 bg-blue-50 text-blue-700 border-l-4 border-blue-600 font-medium text-xs">
            <span className="material-symbols-outlined text-[18px]">explore</span>
            {t('发现', 'Discover')}
          </button>
          <button className="flex items-center gap-3 px-4 py-2 text-slate-500 hover:text-slate-900 hover:bg-slate-100 border-l-4 border-transparent transition-all text-xs">
            <span className="material-symbols-outlined text-[18px]">extension</span>
            {t('我的技能', 'My Skills')}
          </button>
        </div>
      </nav>

      {/* Extension Panel View */}
      <main className="flex-1 flex flex-col bg-white max-w-[400px] border-r border-slate-200 shadow-sm h-full shrink-0 z-0">
        <div className="p-4 border-b border-slate-200 flex flex-col gap-2 shrink-0">
          <div className="relative">
            <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[16px]">search</span>
            <input 
              type="text" 
               className="w-full pl-10 pr-3 py-2 bg-slate-50 text-slate-900 border border-slate-200 rounded focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 text-sm placeholder:text-slate-500 transition-colors" 
              placeholder={t('搜索技能...', 'Search skills...')} 
            />
          </div>
          <div className="bg-blue-50 border border-blue-100 rounded p-3 flex items-start gap-2 mt-2">
            <span className="text-[14px]">💡</span>
            <p className="text-xs text-slate-600 leading-relaxed">
              {t('推荐：根据当前', 'Tip: Based on your current')} <code className="font-mono text-blue-600 bg-white px-1.5 py-0.5 rounded border border-blue-100 text-[11px]">.py</code> {t('文件，建议安装 Python 代码优化技能', 'file, consider installing a Python code optimization skill')}
            </p>
          </div>
        </div>
        <div className="flex-1 overflow-y-auto p-4 flex flex-col gap-3">
          {loading ? (
            <div className="flex flex-col gap-3">
              {[1, 2, 3, 4].map(i => (
                <div key={i} className="border border-slate-200 rounded p-3 flex items-start gap-3">
                  <div className="w-8 h-8 rounded skeleton shrink-0" />
                  <div className="flex-1 flex flex-col gap-2">
                    <div className="h-4 w-32 skeleton" />
                    <div className="h-3 w-full skeleton" />
                  </div>
                </div>
              ))}
            </div>
          ) : skills.length === 0 ? (
            <div className="text-center text-sm text-slate-500 py-8">{t('暂无推荐技能', 'No recommended skills')}</div>
          ) : (
            skills.map(skill => (
              <div key={skill.id} className="group border border-slate-200 rounded p-3 flex items-start gap-3 hover:bg-slate-50 hover:border-blue-200 transition-all cursor-pointer shadow-sm hover:shadow">
                <div className={`w-8 h-8 rounded flex items-center justify-center shrink-0 ${skill.iconBg || 'bg-slate-50'}`}>
                  <span className={`material-symbols-outlined text-[18px] ${skill.iconColor || 'text-slate-600'}`}>{skill.icon}</span>
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex justify-between items-center mb-1">
                    <h3 className="text-sm font-bold text-slate-900 truncate">{skill.title}</h3>
                    <span className="text-[11px] text-slate-500 flex items-center gap-0.5">
                      <span className="material-symbols-outlined text-[12px]">download</span> {(skill.downloads / 1000).toFixed(0)}k
                    </span>
                  </div>
                  <p className="text-xs text-slate-500 line-clamp-2 leading-relaxed">
                    {pickDesc(language, skill)}
                  </p>
                </div>
                <button className="shrink-0 w-7 h-7 flex items-center justify-center rounded bg-blue-600 text-white hover:bg-blue-700 transition-colors shadow-sm" title={t('安装', 'Install')}>
                  <span className="material-symbols-outlined text-[16px]">add</span>
                </button>
              </div>
            ))
          )}
        </div>
      </main>

      {/* Empty Editor Background */}
      <div className="flex-1 bg-[#1e1e1e] flex flex-col relative z-0">
         <div className="w-full flex">
            <div className="px-4 py-2 border-r border-[#333] border-t-2 border-t-blue-500 bg-[#1e1e1e] text-slate-300 text-xs font-mono">
               main.py
            </div>
            <div className="flex-1 border-b border-[#333]"></div>
         </div>
         <div className="flex-1 flex p-4 relative font-mono text-xs">
            <div className="flex flex-col gap-1 text-[#666] select-none pr-4 text-right">
               <span>1</span><span>2</span><span>3</span><span>4</span><span>5</span>
            </div>
            <div className="flex flex-col gap-1 text-slate-300">
               <pre>
<span className="text-purple-400">def</span> <span className="text-blue-300">calculate_metrics</span>(data):
    <span className="text-green-400">"""
    Calculate performance metrics.
    """</span>
    <span className="text-pink-400">pass</span>
               </pre>
               <div className="mt-8 flex items-center justify-center opacity-30 select-none">
                 <p className="text-slate-400 text-[14px]">IDE Editor Canvas</p>
               </div>
            </div>
         </div>
      </div>
    </div>
  );
}
