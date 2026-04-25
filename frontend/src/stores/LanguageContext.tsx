import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';

export type Language = 'zh' | 'en';

interface LanguageContextType {
  language: Language;
  setLanguage: (lang: Language) => void;
  t: (zh: string, en: string) => string;
}

const LanguageContext = createContext<LanguageContextType | undefined>(undefined);

export function LanguageProvider({ children }: { children: ReactNode }) {
  const [language, setLanguage] = useState<Language>('zh');

  useEffect(() => {
    document.documentElement.lang = language === 'zh' ? 'zh-CN' : 'en';
  }, [language]);

  const t = (zh: string, en: string) => language === 'zh' ? zh : en;

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  );
}

export function useLanguage() {
  const ctx = useContext(LanguageContext);
  if (!ctx) throw new Error('useLanguage must be used within LanguageProvider');
  return ctx;
}

export function useDescription(skill: { zhDescription: string; enDescription: string; description: string }): string {
  const { language } = useLanguage();
  if (language === 'zh' && skill.zhDescription) return skill.zhDescription;
  if (language === 'en' && skill.enDescription) return skill.enDescription;
  return skill.description;
}
