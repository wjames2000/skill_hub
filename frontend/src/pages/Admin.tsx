import { Link } from "react-router-dom";

export function Admin() {
  return (
    <div className="font-sans text-sm text-slate-900 bg-slate-50 min-h-screen">
      {/* TopAppBar */}
      <header className="bg-white/80 backdrop-blur-md fixed top-0 right-0 left-0 flex items-center justify-between px-6 h-16 md:ml-64 border-b border-slate-200 z-50">
        <div className="flex items-center gap-4 flex-1">
          <Link to="/" className="text-lg font-bold tracking-tight text-slate-900 md:hidden">SkillHub Pro</Link>
          <div className="relative w-64 ml-8 hidden md:block">
            <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[18px]">search</span>
            <input 
              type="text" 
              className="w-full pl-10 pr-4 py-1.5 bg-slate-100 border-transparent rounded-full text-sm focus:border-blue-500 focus:bg-white focus:ring-0 transition-colors" 
              placeholder="Search..." 
            />
          </div>
        </div>
        <div className="flex items-center gap-4">
          <button className="p-2 text-slate-500 hover:bg-slate-50 rounded-full flex items-center justify-center transition-colors">
            <span className="material-symbols-outlined">notifications</span>
          </button>
          <button className="p-2 text-slate-500 hover:bg-slate-50 rounded-full flex items-center justify-center transition-colors">
            <span className="material-symbols-outlined">terminal</span>
          </button>
          <button className="p-2 text-slate-500 hover:bg-slate-50 rounded-full flex items-center justify-center transition-colors">
            <span className="material-symbols-outlined">settings</span>
          </button>
          <div className="h-8 w-8 rounded-full bg-slate-200 overflow-hidden ml-2 flex-shrink-0 border border-slate-300">
            <img 
              src="https://lh3.googleusercontent.com/aida-public/AB6AXuBaxje5gZ0GlhEXHKwvdnKx05-54jQFSZ0BXjbiJFFp8EFH4MJaxX2S0dGUvSzJwUMMT3XFfyGLQia88QYjqBa-j_00EL8f4gBRYDeSkez5EZoOBxl5GZ1uiPc8RAvgiJo40mM0Xz04DhCKqQQ3YOvwdT7FeSMA2UifmYN-IcBuh-jJNClgJx334NWv9YbAl4GdBDOhil6e0TQdWYTzj9jZv6QwP4TSoCFd1yRVvoHkz3syvh_Oitog6BIjzOpoBFdgQ4z4jcXC5F8" 
              alt="Admin" 
              className="w-full h-full object-cover" 
            />
          </div>
        </div>
      </header>

      {/* SideNavBar */}
      <nav className="bg-white font-medium h-screen w-64 border-r border-slate-200 fixed left-0 top-0 bottom-0 md:flex flex-col pt-16 z-40 hidden">
        <div className="p-6 pb-2">
          <Link to="/" className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-blue-100 flex items-center justify-center">
              <span className="material-symbols-outlined text-blue-600">admin_panel_settings</span>
            </div>
            <div>
              <h2 className="text-xl font-black tracking-tighter text-blue-600 leading-tight">SkillHub Pro</h2>
              <span className="text-xs text-slate-500">Admin Console</span>
            </div>
          </Link>
        </div>
        <div className="flex-1 overflow-y-auto py-4 px-3 flex flex-col gap-1">
          {[
            { label: '仪表盘', icon: 'dashboard', active: true },
            { label: '爬虫任务', icon: 'pest_control' },
            { label: '技能审核', icon: 'verified' },
            { label: '配置管理', icon: 'settings_input_component' },
            { label: '数据统计', icon: 'monitoring' },
            { label: '用户管理', icon: 'group' },
          ].map((item, idx) => (
            <a 
              key={idx} 
              href="#" 
              className={`flex items-center gap-3 px-3 py-2 rounded-lg cursor-pointer transition-all ${item.active ? 'bg-blue-50 text-blue-700 border-r-2 border-blue-600' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'}`}
            >
              <span className="material-symbols-outlined text-[20px]">{item.icon}</span>
              <span>{item.label}</span>
            </a>
          ))}
        </div>
        <div className="p-3 border-t border-slate-200 flex flex-col gap-1">
          <a href="#" className="flex items-center gap-3 px-3 py-2 rounded-lg text-slate-600 hover:bg-slate-100 hover:text-slate-900 transition-all cursor-pointer">
            <span className="material-symbols-outlined text-[20px]">description</span>
            <span>系统日志</span>
          </a>
          <a href="#" className="flex items-center gap-3 px-3 py-2 rounded-lg text-slate-600 hover:bg-slate-100 hover:text-slate-900 transition-all cursor-pointer">
            <span className="material-symbols-outlined text-[20px]">menu_book</span>
            <span>技术文档</span>
          </a>
        </div>
      </nav>

      {/* Main Content */}
      <main className="md:ml-64 pt-16 min-h-screen p-6 max-w-[1440px] mx-auto">
        <div className="mb-8 mt-6">
          <h2 className="text-2xl font-bold text-slate-900 mb-2">系统概览</h2>
          <p className="text-slate-500 text-sm">系统核心指标与最新动态监控。</p>
        </div>

        {/* Metrics Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
          {[
            { label: '技能库总量', value: '12,480', change: '12%', icon: 'library_books', color: 'text-blue-600', bg: 'bg-blue-50' },
            { label: '今日新增技能', value: '86', change: '5%', icon: 'add_circle', color: 'text-blue-600', bg: 'bg-blue-50' },
            { label: 'API 24h 调用量', value: '1.2M', tags: ['正常'], icon: 'api', color: 'text-amber-500', bg: 'bg-amber-50' },
            { label: '当前爬虫状态', value: '3', suffix: '个任务进行中', icon: 'bug_report', color: 'text-green-500', bg: 'bg-green-50' },
          ].map((m, idx) => (
            <div key={idx} className="bg-white rounded-lg border border-slate-200 p-6 shadow-sm hover:shadow-md transition-shadow">
              <div className="flex items-center justify-between mb-4">
                <span className="text-slate-500 text-sm font-medium">{m.label}</span>
                <span className={`material-symbols-outlined p-1.5 rounded-md ${m.color} ${m.bg}`}>{m.icon}</span>
              </div>
              <div className="flex items-baseline gap-2">
                <span className="text-3xl font-bold text-slate-900">{m.value}</span>
                {m.change && (
                  <span className="text-green-500 text-sm font-medium flex items-center">
                    <span className="material-symbols-outlined text-[16px]">arrow_upward</span>{m.change}
                  </span>
                )}
                {m.tags && m.tags.map(tag => (
                  <span key={tag} className="text-green-600 text-xs font-medium ml-2 bg-green-50 px-2 py-0.5 rounded border border-green-200">{tag}</span>
                ))}
                {m.suffix && <span className="text-slate-500 text-sm ml-1">{m.suffix}</span>}
              </div>
              {m.suffix && (
                <div className="mt-3 w-full bg-slate-100 rounded-full h-1.5 border border-slate-200">
                  <div className="bg-green-500 h-1.5 rounded-full" style={{ width: '100%' }}></div>
                </div>
              )}
            </div>
          ))}
        </div>

        {/* Chart Area */}
        <div className="bg-white rounded-lg border border-slate-200 p-6 shadow-sm mb-6">
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-xl font-bold text-slate-900">近 7 日技能增长与 API 调用趋势</h3>
            <button className="text-blue-600 hover:text-blue-800 text-sm flex items-center gap-1 font-medium">
              查看详情 <span className="material-symbols-outlined text-[18px]">chevron_right</span>
            </button>
          </div>
          <div className="h-64 w-full flex items-end gap-2 pt-8 relative border-b border-l border-slate-200 pl-4 pb-4">
            <div className="absolute left-[-30px] bottom-4 top-8 flex flex-col justify-between text-xs text-slate-400 h-full">
              <span>10k</span>
              <span>5k</span>
              <span>0</span>
            </div>
            {[
              { label: '周一', h1: '40%', h2: '60%' },
              { label: '周二', h1: '45%', h2: '65%' },
              { label: '周三', h1: '35%', h2: '50%' },
              { label: '周四', h1: '50%', h2: '70%' },
              { label: '周五', h1: '60%', h2: '75%' },
              { label: '周六', h1: '75%', h2: '85%' },
              { label: '周日', h1: '85%', h2: '95%' },
            ].map((col, idx) => (
              <div key={idx} className="flex-1 flex items-end justify-around h-full relative group">
                <div className="w-4 bg-blue-600/30 rounded-t hover:bg-blue-600 transition-colors cursor-pointer" style={{ height: col.h1 }}></div>
                <div className="w-4 bg-blue-400/50 rounded-t ml-1 hover:bg-blue-400 transition-colors cursor-pointer" style={{ height: col.h2 }}></div>
                <span className="absolute -bottom-6 text-xs text-slate-400">{col.label}</span>
              </div>
            ))}
            <div className="absolute top-0 right-0 flex gap-4 text-xs font-medium text-slate-500">
              <div className="flex items-center gap-1"><div className="w-3 h-3 bg-blue-600 rounded-sm"></div>新增技能</div>
              <div className="flex items-center gap-1"><div className="w-3 h-3 bg-blue-400 rounded-sm"></div>API 调用量</div>
            </div>
          </div>
        </div>

        {/* Lists Area */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="bg-white rounded-lg border border-slate-200 shadow-sm flex flex-col h-[400px]">
            <div className="p-4 border-b border-slate-200 flex justify-between items-center bg-slate-50/50 rounded-t-lg">
              <h3 className="text-lg font-bold text-slate-900 flex items-center gap-2">
                <span className="material-symbols-outlined text-amber-500 text-[20px]">pending_actions</span>
                待审核技能
              </h3>
              <a href="#" className="text-blue-600 hover:text-blue-800 text-sm font-medium">全部审核</a>
            </div>
            <div className="flex-1 overflow-y-auto p-0">
              <ul className="divide-y divide-slate-100">
                {[
                  { title: 'Data Processing Pipeline', author: 'alice.dev', time: '10分钟前' },
                  { title: 'AWS S3 Bucket Analyzer', author: 'bob_smith', time: '1小时前' },
                  { title: 'Natural Language Query to SQL', author: 'data_wizard', time: '2小时前' },
                  { title: 'React Component Generator', author: 'ui_ninja', time: '4小时前' },
                  { title: 'Automated PR Reviewer', author: 'ci_cd_bot', time: '5小时前' },
                ].map((req, idx) => (
                  <li key={idx} className="p-4 hover:bg-slate-50 transition-colors flex justify-between items-center group cursor-pointer">
                    <div>
                      <h4 className="font-medium text-slate-900 mb-1 group-hover:text-blue-600 transition-colors">{req.title}</h4>
                      <div className="flex items-center gap-3 text-xs text-slate-500">
                        <span className="flex items-center gap-1"><span className="material-symbols-outlined text-[14px]">person</span> {req.author}</span>
                        <span className="flex items-center gap-1"><span className="material-symbols-outlined text-[14px]">schedule</span> {req.time}</span>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <button className="px-3 py-1.5 text-xs font-medium text-green-600 border border-green-200 bg-white hover:bg-green-50 rounded transition-colors shadow-sm">通过</button>
                      <button className="px-3 py-1.5 text-xs font- medium text-red-600 border border-red-200 bg-white hover:bg-red-50 rounded transition-colors shadow-sm">驳回</button>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          </div>

          <div className="bg-white rounded-lg border border-slate-200 shadow-sm flex flex-col h-[400px]">
            <div className="p-4 border-b border-slate-200 flex justify-between items-center bg-slate-50/50 rounded-t-lg">
              <h3 className="text-lg font-bold text-slate-900 flex items-center gap-2">
                <span className="material-symbols-outlined text-blue-500 text-[20px]">terminal</span>
                实时爬虫日志
              </h3>
              <div className="flex items-center gap-2 border border-slate-200 bg-white px-2 py-1 rounded-full shadow-sm text-xs">
                <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
                <span className="text-slate-600 font-medium whitespace-nowrap">运行中 (3)</span>
              </div>
            </div>
            <div className="flex-1 overflow-y-auto p-4 bg-[#1e1e1e] rounded-b-lg font-mono text-sm text-slate-300 shadow-inner">
              <div className="flex flex-col gap-2">
                <div className="flex gap-3 items-start"><span className="text-slate-500 whitespace-nowrap">[14:32:01]</span><span className="text-green-400">[INFO]</span><span className="text-slate-200">Synced repo anthropics/skills successfully.</span></div>
                <div className="flex gap-3 items-start"><span className="text-slate-500 whitespace-nowrap">[14:31:45]</span><span className="text-blue-400">[SYS]</span><span className="text-slate-200">Vector index updated (batch #492).</span></div>
                <div className="flex gap-3 items-start"><span className="text-slate-500 whitespace-nowrap">[14:30:12]</span><span className="text-amber-400">[WARN]</span><span className="text-slate-200">Rate limit approaching for github_api_worker_2. Backing off 5s.</span></div>
                <div className="flex gap-3 items-start"><span className="text-slate-500 whitespace-nowrap">[14:28:55]</span><span className="text-green-400">[INFO]</span><span className="text-slate-200">Extracted 14 new skill definitions from lang-chain/community.</span></div>
                <div className="flex gap-3 items-start"><span className="text-slate-500 whitespace-nowrap">[14:25:30]</span><span className="text-blue-400">[SYS]</span><span className="text-slate-200">Initiating scheduled crawl for tag: 'data-analysis'.</span></div>
                <div className="flex gap-3 items-start opacity-70"><span className="text-slate-500 whitespace-nowrap">[14:20:00]</span><span className="text-green-400">[INFO]</span><span className="text-slate-200">Database optimization routine completed.</span></div>
                <div className="flex gap-3 items-start opacity-50"><span className="text-slate-500 whitespace-nowrap">[14:15:10]</span><span className="text-blue-400">[SYS]</span><span className="text-slate-200">Worker health check passed (12/12 online).</span></div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
