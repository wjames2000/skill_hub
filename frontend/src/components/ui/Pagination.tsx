import { cn } from "../../lib/utils";
import { useLanguage } from "../../stores/LanguageContext";

interface Props {
  page: number;
  total: number;
  pageSize: number;
  onChange: (page: number) => void;
}

export function Pagination({ page, total, pageSize, onChange }: Props) {
  const { t } = useLanguage();
  const totalPages = Math.ceil(total / pageSize);
  if (totalPages <= 1) return null;

  const pages: (number | 'ellipsis')[] = [];
  for (let i = 1; i <= totalPages; i++) {
    if (i === 1 || i === totalPages || (i >= page - 1 && i <= page + 1)) {
      pages.push(i);
    } else if (pages[pages.length - 1] !== 'ellipsis') {
      pages.push('ellipsis');
    }
  }

  return (
    <div className="flex items-center justify-center gap-2 mt-8" role="navigation" aria-label="Pagination">
      <button
        onClick={() => onChange(page - 1)}
        disabled={page <= 1}
        aria-label={t('上一页', 'Previous page')}
        className="w-8 h-8 rounded border border-slate-200 flex items-center justify-center text-slate-400 disabled:opacity-40 hover:bg-slate-50 transition-colors disabled:cursor-not-allowed"
      >
        <span className="material-symbols-outlined text-[18px]">chevron_left</span>
      </button>
      {pages.map((p, idx) =>
        p === 'ellipsis' ? (
          <span key={`e-${idx}`} className="text-slate-400 px-1">...</span>
        ) : (
          <button
            key={p}
            onClick={() => onChange(p)}
            aria-current={p === page ? 'page' : undefined}
            className={cn(
              "w-8 h-8 rounded text-sm font-medium transition-colors",
              p === page
                ? 'bg-brand-600 text-white'
                : 'border border-slate-200 text-slate-600 hover:bg-slate-50'
            )}
          >
            {p}
          </button>
        ),
      )}
      <button
        onClick={() => onChange(page + 1)}
        disabled={page >= totalPages}
        aria-label={t('下一页', 'Next page')}
        className="w-8 h-8 rounded border border-slate-200 flex items-center justify-center text-slate-400 disabled:opacity-40 hover:bg-slate-50 transition-colors disabled:cursor-not-allowed"
      >
        <span className="material-symbols-outlined text-[18px]">chevron_right</span>
      </button>
    </div>
  );
}
