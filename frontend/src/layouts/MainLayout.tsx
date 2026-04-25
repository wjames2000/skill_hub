import { Link, Outlet, useNavigate } from "react-router-dom";
import { useState, useEffect, useRef, type FormEvent } from "react";
import { useAuth } from "../stores/AuthContext";
import { useLanguage, type Language } from "../stores/LanguageContext";

export function MainLayout() {
  const { isAuthenticated, user, logout } = useAuth();
  const { language, setLanguage, t } = useLanguage();
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState("");
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setDropdownOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === "Escape") setDropdownOpen(false);
    };
    document.addEventListener("keydown", handleEsc);
    return () => document.removeEventListener("keydown", handleEsc);
  }, []);

  const handleSearch = (e: FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      navigate(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
    }
  };

  return (
    <div className="bg-slate-50 min-h-screen flex flex-col text-slate-900 font-sans">
      <a
        href="#main-content"
        className="sr-only focus:not-sr-only focus:absolute focus:top-0 focus:left-0 focus:z-[100] focus:px-4 focus:py-2 focus:bg-brand-600 focus:text-white focus:font-medium focus:text-sm"
      >
        {t('跳到主内容', 'Skip to main content')}
      </a>
      <header className="bg-white/90 backdrop-blur-md sticky top-0 z-50 border-b border-slate-200">
        <div className="flex items-center justify-between px-4 md:px-6 h-16 w-full max-w-[1440px] mx-auto">
          <div className="flex items-center gap-4 md:gap-6">
            <Link to="/" className="text-xl font-bold tracking-tighter text-slate-900 flex items-center gap-2 shrink-0">
              <span className="material-symbols-outlined text-brand-600 fill">hexagon</span>
              <span className="hidden sm:inline">SkillHub Pro</span>
            </Link>

            <form onSubmit={handleSearch} className="hidden md:flex relative items-center w-48 lg:w-72 xl:w-96">
              <span className="material-symbols-outlined absolute left-3 text-slate-400 text-sm">search</span>
              <input
                    type="text"
                    value={searchQuery}
                    onChange={e => setSearchQuery(e.target.value)}
                    className="w-full pl-10 pr-4 py-2 bg-slate-100 border-none rounded-md text-sm focus:ring-2 focus:ring-brand-500 transition-all text-slate-900 placeholder:text-slate-500"
                    placeholder={t('搜索技能、作者...', 'Search skills, authors...')}
                  />
            </form>
          </div>

          <div className="flex items-center gap-2 md:gap-4">
            <div className="hidden md:flex gap-4 mr-2">
              <Link to="/search" className="text-slate-600 font-medium hover:text-brand-600 transition-colors text-sm">{t('浏览', 'Browse')}</Link>
              <Link to="/ide" className="text-slate-600 font-medium hover:text-brand-600 transition-colors text-sm">{t('IDE 插件', 'IDE Plugin')}</Link>
              {user?.role === 'admin' && (
                <Link to="/admin" className="text-slate-600 font-medium hover:text-brand-600 transition-colors text-sm">{t('管理', 'Admin')}</Link>
              )}
            </div>

            <button
              onClick={() => setLanguage(language === 'zh' ? 'en' : 'zh')}
              className="px-2 py-1 text-xs font-medium rounded border border-slate-300 text-slate-600 hover:bg-slate-100 transition-colors"
              title="切换语言"
            >
              {language === 'zh' ? 'EN' : '中'}
            </button>

            <button
              onClick={() => setSearchQuery('')}
              className="md:hidden p-2 text-slate-600 hover:bg-slate-100 rounded-full transition-colors"
              aria-label="搜索"
            >
              <span className="material-symbols-outlined">search</span>
            </button>

            {isAuthenticated ? (
              <div className="flex items-center gap-3">
                <button className="text-slate-600 hover:bg-slate-100 transition-colors p-2 rounded-full hidden sm:block" aria-label={t('通知', 'Notifications')}>
                  <span className="material-symbols-outlined">notifications</span>
                </button>
                <div className="relative" ref={dropdownRef}>
                  <button
                    onClick={() => setDropdownOpen(prev => !prev)}
                    className="h-8 w-8 rounded-full overflow-hidden border border-slate-200 cursor-pointer"
                    aria-haspopup="true"
                    aria-expanded={dropdownOpen}
                    aria-label={t('用户菜单', 'User menu')}
                  >
                    <img
                      src={user?.avatar || 'https://ui-avatars.com/api/?name=User&background=6366f1&color=fff'}
                      alt=""
                      className="w-full h-full object-cover"
                    />
                  </button>
                  <div
                    className={`absolute right-0 top-full mt-2 w-48 bg-white border border-slate-200 rounded-xl shadow-elevated transition-all duration-200 py-1 z-50 ${
                      dropdownOpen ? 'opacity-100 visible' : 'opacity-0 invisible pointer-events-none'
                    }`}
                    role="menu"
                  >
                    <Link to="/profile" onClick={() => setDropdownOpen(false)} className="flex items-center gap-3 px-4 py-2.5 text-sm text-slate-700 hover:bg-slate-50 transition-colors" role="menuitem">
                      <span className="material-symbols-outlined text-[18px]">person</span>
                      {t('个人中心', 'Profile')}
                    </Link>
                    {user?.role === 'admin' && (
                      <Link to="/admin" onClick={() => setDropdownOpen(false)} className="flex items-center gap-3 px-4 py-2.5 text-sm text-slate-700 hover:bg-slate-50 transition-colors" role="menuitem">
                        <span className="material-symbols-outlined text-[18px]">admin_panel_settings</span>
                        {t('管理后台', 'Admin Panel')}
                      </Link>
                    )}
                    <hr className="my-1 border-slate-100" />
                    <button onClick={() => { logout(); setDropdownOpen(false); }} className="flex items-center gap-3 px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 transition-colors w-full text-left" role="menuitem">
                      <span className="material-symbols-outlined text-[18px]">logout</span>
                      {t('退出登录', 'Logout')}
                    </button>
                  </div>
                </div>
              </div>
            ) : (
              <div className="flex items-center gap-2">
                <Link to="/login" className="btn-ghost text-sm hidden sm:flex">{t('登录', 'Login')}</Link>
                <Link to="/register" className="btn-primary text-sm !py-1.5 !px-4">{t('注册', 'Register')}</Link>
              </div>
            )}
          </div>
        </div>
      </header>

      <div id="main-content" className="flex flex-1 w-full max-w-[1440px] mx-auto">
        <Outlet />
      </div>

      <footer className="bg-white text-slate-500 text-xs w-full border-t border-slate-200 mt-auto">
        <div className="flex flex-col md:flex-row justify-between items-center px-6 md:px-8 py-8 md:py-12 w-full max-w-[1280px] mx-auto gap-4 md:gap-0">
          <div className="font-bold text-slate-400">
            © 2024 SkillHub Pro. All rights reserved.
          </div>
          <nav className="flex flex-wrap justify-center gap-4 md:gap-6">
            <a href="#" className="hover:text-brand-600 transition-colors">{t('文档', 'Docs')}</a>
            <a href="#" className="hover:text-brand-600 transition-colors">{t('更新日志', 'Changelog')}</a>
            <a href="#" className="hover:text-brand-600 transition-colors">{t('API 状态', 'API Status')}</a>
            <a href="#" className="hover:text-brand-600 transition-colors">{t('隐私', 'Privacy')}</a>
            <a href="#" className="hover:text-brand-600 transition-colors">{t('条款', 'Terms')}</a>
          </nav>
        </div>
      </footer>
    </div>
  );
}
