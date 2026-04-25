import { useState } from 'react';
import { Link } from 'react-router-dom';
import type { Category } from '@/src/types';
import type { Language } from '@/src/stores/LanguageContext';

function pickCatName(lang: Language, cat: Category): string {
  if (lang === 'zh' && cat.zhName) return cat.zhName;
  if (lang === 'en' && cat.enName) return cat.enName;
  return cat.name;
}

interface Props {
  categories: Category[];
  selected?: string;
  onSelect?: (cat: Category) => void;
  linkMode?: boolean;
  depth?: number;
  language?: Language;
}

export function CategoryTree({ categories, selected, onSelect, linkMode = false, depth = 0, language = 'zh' }: Props) {
  if (!categories || categories.length === 0) return null;

  return (
    <div className="flex flex-col">
      {categories.map(cat => (
        <CategoryTreeNode
          key={cat.id}
          category={cat}
          selected={selected}
          onSelect={onSelect}
          linkMode={linkMode}
          depth={depth}
          language={language}
        />
      ))}
    </div>
  );
}

function CategoryTreeNode({ category, selected, onSelect, linkMode, depth, language }: { category: Category; selected?: string; onSelect?: (cat: Category) => void; linkMode: boolean; depth: number; language: Language }) {
  const [expanded, setExpanded] = useState(depth < 1 && category.children.length > 0);
  const hasChildren = category.children && category.children.length > 0;
  const isSelected = selected === category.slug;

  const toggle = () => setExpanded(!expanded);

  const content = (
    <div
      className={`flex items-center gap-2 px-4 lg:px-6 py-2.5 text-sm transition-all cursor-pointer ${
        isSelected
          ? 'bg-brand-50 text-brand-700 border-r-2 border-brand-600 font-medium'
          : 'text-slate-500 hover:bg-slate-100 hover:text-slate-900'
      }`}
      style={{ paddingLeft: depth > 0 ? `${16 + depth * 16}px` : undefined }}
      onClick={linkMode ? undefined : () => { onSelect?.(category); }}
    >
      {hasChildren && (
        <button onClick={(e) => { e.stopPropagation(); toggle(); }} className="text-slate-400 hover:text-slate-600 p-0.5" aria-expanded={expanded} aria-label={expanded ? 'Collapse category' : 'Expand category'}>
          <span className="material-symbols-outlined text-[16px] transition-transform" style={{ transform: expanded ? 'rotate(90deg)' : undefined }}>
            chevron_right
          </span>
        </button>
      )}
      {!hasChildren && <span className="w-5" />}
      <span className="material-symbols-outlined text-lg">{category.icon}</span>
      <span className="flex-1 truncate">{pickCatName(language, category)}</span>
      <span className="text-xs text-slate-400">({category.count})</span>
    </div>
  );

  return (
    <div>
      {linkMode ? (
        <Link to={`/search?category=${category.slug}`}>
          {content}
        </Link>
      ) : (
        <div onClick={hasChildren ? toggle : undefined}>{content}</div>
      )}
      {hasChildren && expanded && (
        <CategoryTree
          categories={category.children}
          selected={selected}
          onSelect={onSelect}
          linkMode={linkMode}
          depth={depth + 1}
          language={language}
        />
      )}
    </div>
  );
}
