import { Link } from "react-router-dom";
import type { Skill } from "../../types";
import { useDescription } from "../../stores/LanguageContext";

interface Props {
  skill: Skill;
  showMatchScore?: boolean;
}

const iconMap: Record<string, string> = {
  'data_object': 'data_object',
  'terminal': 'terminal',
  'analytics': 'analytics',
  'language': 'language',
  'table_chart': 'table_chart',
  'query_stats': 'query_stats',
  'insights': 'insights',
  'description': 'description',
  'code': 'code',
  'database': 'database',
  'finance': 'finance',
  'palette': 'palette',
  'robot_2': 'robot_2',
  'translate': 'translate',
  'default': 'extension',
};

export function SkillCard({ skill, showMatchScore }: Props) {
  const icon = iconMap[skill.icon] || iconMap.default;
  const desc = useDescription(skill);

  return (
    <Link
      to={`/skill/${skill.id}`}
      className="card p-5 flex flex-col h-full relative group overflow-hidden"
    >
      <div className="flex items-start justify-between mb-3 relative z-10">
        <div className="flex items-center gap-3">
          <div className={`w-10 h-10 rounded flex items-center justify-center ${skill.iconBg || 'bg-brand-50'} ${skill.iconColor || 'text-brand-600'}`}>
            <span className="material-symbols-outlined">{icon}</span>
          </div>
          <div>
            <h3 className="font-bold text-base text-slate-900 group-hover:text-brand-600 transition-colors line-clamp-1">
              {skill.title}
            </h3>
            <p className="text-xs text-slate-500">by {skill.author}</p>
          </div>
        </div>
        <div className="flex items-center gap-1 bg-slate-50 px-2 py-1 rounded border border-slate-200">
          <span className="material-symbols-outlined text-[14px] text-amber-500 fill">star</span>
          <span className="font-mono text-[12px] font-medium text-slate-700">{skill.rating.toFixed(1)}</span>
        </div>
      </div>

      <p className="text-sm text-slate-600 line-clamp-2 mb-4 flex-1">
        {desc}
      </p>

      <div className="flex items-center justify-between mt-auto pt-4 border-t border-slate-100">
        <div className="flex gap-2 flex-wrap">
          {skill.tags.slice(0, 2).map(tag => (
            <span key={tag} className="px-2 py-1 bg-slate-50 text-slate-600 rounded font-mono text-[11px] border border-slate-200">
              {tag}
            </span>
          ))}
          {skill.tags.length > 2 && (
            <span className="px-2 py-1 bg-slate-50 text-slate-400 rounded font-mono text-[11px] border border-slate-200">
              +{skill.tags.length - 2}
            </span>
          )}
        </div>
        <div className="flex items-center gap-3 text-xs text-slate-400">
          <span className="flex items-center gap-1">
            <span className="material-symbols-outlined text-[14px]">download</span>
            {(skill.downloads / 1000).toFixed(1)}k
          </span>
          {showMatchScore && skill.matchScore != null && (
            <span className="font-bold text-green-500">{skill.matchScore}%</span>
          )}
        </div>
      </div>

      {skill.safe && (
        <div className="absolute top-3 right-3">
          <span className="material-symbols-outlined text-[16px] text-green-500 fill">verified</span>
        </div>
      )}
    </Link>
  );
}
