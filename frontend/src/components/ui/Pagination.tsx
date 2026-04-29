import { cn } from "../../lib/utils";
import { useLanguage } from "../../stores/LanguageContext";

interface Props {
  page: number;
  total: number;
  pageSize: number;
  onChange: (page: number) => void;
  onPageSizeChange?: (size: number) => void;
}

export function Pagination({ page, total, pageSize, onChange, onPageSizeChange }: Props) {
  const { t } = useLanguage();
  const totalPages = Math.ceil(total / pageSize);
  if (totalPages <= 1 && !onPageSizeChange) return null;

  const from = (page - 1) * pageSize + 1;
  const to = Math.min(page * pageSize, total);
  const pageSizes = [5, 10, 20, 50];

  const pages: (number | 'ellipsis')[] = [];
  if (totalPages > 1) {
    for (let i = 1; i <= totalPages; i++) {
      if (i === 1 || i === totalPages || (i >= page - 1 && i <= page + 1)) {
        pages.push(i);
      } else if (pages[pages.length - 1] !== 'ellipsis') {
        pages.push('ellipsis');
      }
    }
  }

  return (
    <div className="flex flex-col sm:flex-row items-center justify-between gap-4 mt-8">
      <div className="text-sm text-slate-500">
        {t('共', '')}<span className="font-medium text-slate-700">{total}</span>{t('条结果，第', ' results, page ')}
        <span className="font-medium text-slate-700">{page}/{totalPages > 0 ? totalPages : 1}</span>{t('页', '')}
        <span className="mx-2 text-slate-300">|</span>
        {t('每页', 'Per page: ')}
        <select
          value={pageSize}
          onChange={e => onPageSizeChange?.(Number(e.target.value))}
          className="ml-1 border border-slate-200 rounded px-1.5 py-0.5 text-sm text-slate-700 bg-white focus:outline-none focus:border-brand-500"
        >
          {pageSizes.map(s => (
            <option key={s} value={s}>{s}</option>
          ))}
        </select>
      </div>

      {totalPages > 1 && (
        <div className="flex items-center gap-1.5" role="navigation" aria-label="Pagination">
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
      )}
    </div>
  );
}
