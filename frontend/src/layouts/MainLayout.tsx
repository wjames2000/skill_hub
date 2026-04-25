import { Link, Outlet, useNavigate } from "react-router-dom";
import { useState, type FormEvent } from "react";
import { useAuth } from "../stores/AuthContext";
import { useLanguage, type Language } from "../stores/LanguageContext";

export function MainLayout() {
  const { isAuthenticated, user, logout } = useAuth();
  const { language, setLanguage } = useLanguage();
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState("");

  const handleSearch = (e: FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      navigate(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
    }
  };

  return (
    <div className="bg-slate-50 min-h-screen flex flex-col text-slate-900 font-sans">
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
                className="w-full pl-10 pr-4 py-2 bg-slate-100 border-none rounded-md text-sm focus:ring-2 focus:ring-brand-500 transition-all text-slate-900"
                placeholder="搜索技能、作者..."
              />
            </form>
          </div>

          <div className="flex items-center gap-2 md:gap-4">
            <div className="hidden md:flex gap-4 mr-2">
              <Link to="/search" className="text-slate-600 font-medium hover:text-brand-600 transition-colors text-sm">浏览</Link>
              <Link to="/ide" className="text-slate-600 font-medium hover:text-brand-600 transition-colors text-sm">IDE 插件</Link>
              {user?.role === 'admin' && (
                <Link to="/admin" className="text-slate-600 font-medium hover:text-brand-600 transition-colors text-sm">管理</Link>
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
                <button className="text-slate-600 hover:bg-slate-100 transition-colors p-2 rounded-full hidden sm:block">
                  <span className="material-symbols-outlined">notifications</span>
                </button>
                <div className="relative group">
                  <button className="h-8 w-8 rounded-full overflow-hidden border border-slate-200 cursor-pointer">
                    <img
                      src={user?.avatar || 'https://ui-avatars.com/api/?name=User&background=6366f1&color=fff'}
                      alt=""
                      className="w-full h-full object-cover"
                    />
                  </button>
                  <div className="absolute right-0 top-full mt-2 w-48 bg-white border border-slate-200 rounded-xl shadow-elevated opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 py-1 z-50">
                    <Link to="/profile" className="flex items-center gap-3 px-4 py-2.5 text-sm text-slate-700 hover:bg-slate-50 transition-colors">
                      <span className="material-symbols-outlined text-[18px]">person</span>
                      个人中心
                    </Link>
                    {user?.role === 'admin' && (
                      <Link to="/admin" className="flex items-center gap-3 px-4 py-2.5 text-sm text-slate-700 hover:bg-slate-50 transition-colors">
                        <span className="material-symbols-outlined text-[18px]">admin_panel_settings</span>
                        管理后台
                      </Link>
                    )}
                    <hr className="my-1 border-slate-100" />
                    <button onClick={logout} className="flex items-center gap-3 px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 transition-colors w-full text-left">
                      <span className="material-symbols-outlined text-[18px]">logout</span>
                      退出登录
                    </button>
                  </div>
                </div>
              </div>
            ) : (
              <div className="flex items-center gap-2">
                <Link to="/login" className="btn-ghost text-sm hidden sm:flex">登录</Link>
                <Link to="/register" className="btn-primary text-sm !py-1.5 !px-4">注册</Link>
              </div>
            )}
          </div>
        </div>
      </header>

      <div className="flex flex-1 w-full max-w-[1440px] mx-auto">
        <Outlet />
      </div>

      <footer className="bg-white text-slate-500 text-xs w-full border-t border-slate-200 mt-auto">
        <div className="flex flex-col md:flex-row justify-between items-center px-6 md:px-8 py-8 md:py-12 w-full max-w-[1280px] mx-auto gap-4 md:gap-0">
          <div className="font-bold text-slate-400">
            © 2024 SkillHub Pro. All rights reserved.
          </div>
          <nav className="flex flex-wrap justify-center gap-4 md:gap-6">
            <a href="#" className="hover:text-brand-600 transition-colors">文档</a>
            <a href="#" className="hover:text-brand-600 transition-colors">更新日志</a>
            <a href="#" className="hover:text-brand-600 transition-colors">API 状态</a>
            <a href="#" className="hover:text-brand-600 transition-colors">隐私</a>
            <a href="#" className="hover:text-brand-600 transition-colors">条款</a>
          </nav>
        </div>
      </footer>
    </div>
  );
}
