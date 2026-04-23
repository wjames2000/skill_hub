import { Link, Outlet } from "react-router-dom";

export function MainLayout() {
  return (
    <div className="bg-slate-50 min-h-screen flex flex-col text-slate-900 font-sans">
      {/* Top Navbar */}
      <header className="bg-white/90 backdrop-blur-md sticky top-0 z-50 border-b border-slate-200">
        <div className="flex items-center justify-between px-6 h-16 w-full max-w-[1440px] mx-auto">
          <div className="flex items-center gap-6">
            <Link to="/" className="text-xl font-bold tracking-tighter text-slate-900 flex items-center gap-2">
              <span className="material-symbols-outlined text-blue-600 fill">hexagon</span>
              SkillHub Pro
            </Link>
            
            <div className="hidden md:flex relative items-center w-64 lg:w-96">
              <span className="material-symbols-outlined absolute left-3 text-slate-400 text-sm">search</span>
              <input
                type="text"
                className="w-full pl-10 pr-4 py-2 bg-slate-100 border-none rounded-md text-sm focus:ring-2 focus:ring-blue-500 transition-all text-slate-900"
                placeholder="Search skills, authors..."
              />
              <div className="absolute right-3 hidden lg:flex items-center gap-1">
                <kbd className="px-1.5 py-0.5 bg-white border border-slate-200 rounded text-xs font-mono text-slate-500">⌘</kbd>
                <kbd className="px-1.5 py-0.5 bg-white border border-slate-200 rounded text-xs font-mono text-slate-500">K</kbd>
              </div>
            </div>
          </div>
          
          <div className="flex items-center gap-4">
            <div className="hidden md:flex gap-6 mr-4">
              <Link to="/ide" className="text-slate-600 font-medium hover:text-blue-600 transition-colors">IDE插件</Link>
              <Link to="/admin" className="text-slate-600 font-medium hover:text-blue-600 transition-colors">管理后台</Link>
            </div>
            
            <button className="text-slate-600 hover:bg-slate-50 hover:text-blue-600 transition-colors p-2 rounded-full">
              <span className="material-symbols-outlined">notifications</span>
            </button>
            <div className="h-8 w-8 rounded-full overflow-hidden border border-slate-200 cursor-pointer">
              <img
                src="https://lh3.googleusercontent.com/aida-public/AB6AXuCCPfE0v0JhSF7vFOAhPofV1uSZ_nz0eMhokS5XvTMMfRtt_u-OuA6BIOunTXy8F_U0Nc1CcKM6PeF3dxC2cyiBpY5LTjVy-tDXB6OwtaPhRQpllnmAvpwpHNqctJmdioEXSUHfPd4lT9X3gZSwCAIA4wfxGBrZyYxUARq0-0OR6wkZqcVkX_wHrUlHIQIBkdF7BqRs8QuWmEwfddMVnFvB6edzxul1NNhRPWSkO9gTQi4SMYlaSm4n9jTBK7cpde4UCpBsSaLvEgg"
                alt="User Avatar"
                className="w-full h-full object-cover"
              />
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <div className="flex flex-1 w-full max-w-[1440px] mx-auto">
        <Outlet />
      </div>

      {/* Footer */}
      <footer className="bg-white text-slate-500 text-xs w-full border-t border-slate-200 mt-auto">
        <div className="flex flex-col md:flex-row justify-between items-center px-8 py-12 w-full max-w-[1280px] mx-auto">
          <div className="font-bold text-slate-400 mb-4 md:mb-0">
            © 2024 SkillHub Pro. All rights reserved.
          </div>
          <nav className="flex flex-wrap justify-center gap-6">
            <a href="#" className="hover:text-blue-600 transition-colors">Documentation</a>
            <a href="#" className="hover:text-blue-600 transition-colors">Changelog</a>
            <a href="#" className="hover:text-blue-600 transition-colors">API Status</a>
            <a href="#" className="hover:text-blue-600 transition-colors">Privacy</a>
            <a href="#" className="hover:text-blue-600 transition-colors">Terms</a>
          </nav>
        </div>
      </footer>
    </div>
  );
}
