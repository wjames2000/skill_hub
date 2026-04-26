import { useLanguage } from "../../stores/LanguageContext";

export function LoadingScreen() {
  const { t } = useLanguage();
  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-50">
      <div className="flex flex-col items-center gap-4">
        <div className="w-10 h-10 border-4 border-brand-200 border-t-brand-600 rounded-full animate-spin" />
        <p className="text-slate-500 text-sm font-medium">{t('加载中...', 'Loading...')}</p>
      </div>
    </div>
  );
}
