import { Link } from "react-router-dom";

export function Home() {
  return (
    <div className="flex w-full">
      {/* SideNavBar */}
      <aside className="hidden md:flex flex-col bg-slate-50 w-64 border-r border-slate-200">
        <div className="px-6 mb-4 pt-8">
          <h3 className="text-slate-900 font-bold text-sm tracking-normal capitalize">Explore Skills</h3>
          <p className="text-slate-500 text-xs tracking-normal mt-1 capitalize">Browse by domain</p>
        </div>
        <nav className="flex flex-col gap-1 w-full text-sm">
          <a href="#" className="flex items-center gap-3 px-6 py-3 text-slate-500 hover:bg-slate-100 transition-all">
            <span className="material-symbols-outlined text-lg">description</span>
            Documentation
          </a>
          <a href="#" className="flex items-center gap-3 px-6 py-3 text-slate-500 hover:bg-slate-100 transition-all">
            <span className="material-symbols-outlined text-lg">code</span>
            Code
          </a>
          <a href="#" className="flex items-center gap-3 px-6 py-3 bg-blue-50 text-blue-700 border-r-2 border-blue-600">
            <span className="material-symbols-outlined text-lg">database</span>
            Data Analysis
          </a>
          <a href="#" className="flex items-center gap-3 px-6 py-3 text-slate-500 hover:bg-slate-100 transition-all">
            <span className="material-symbols-outlined text-lg">palette</span>
            Design
          </a>
          <a href="#" className="flex items-center gap-3 px-6 py-3 text-slate-500 hover:bg-slate-100 transition-all">
            <span className="material-symbols-outlined text-lg">robot_2</span>
            Automation
          </a>
          <a href="#" className="flex items-center gap-3 px-6 py-3 text-slate-500 hover:bg-slate-100 transition-all">
            <span className="material-symbols-outlined text-lg">translate</span>
            Translation
          </a>
        </nav>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 px-4 md:px-8 py-8 flex flex-col gap-12 max-w-[1280px] mx-auto w-full">
        {/* Hero Section */}
        <section className="flex flex-col items-center justify-center text-center py-16 px-4 bg-white rounded-xl border border-neutral-200 shadow-sm relative overflow-hidden">
          <div className="absolute inset-0 bg-gradient-to-br from-blue-50/50 to-white pointer-events-none"></div>
          <div className="relative z-10 max-w-3xl flex flex-col items-center gap-6">
            <div className="inline-flex items-center gap-2 px-3 py-1 bg-blue-50 text-blue-600 rounded-full text-sm font-medium border border-blue-200/50 mb-2">
              <span className="material-symbols-outlined text-sm fill">new_releases</span>
              v2.4 发布: 支持本地模型接入
            </div>
            <h1 className="text-3xl md:text-4xl font-bold text-slate-900 tracking-tight leading-tight">让 AI 能力的获取像安装 IDE 插件一样简单</h1>
            <p className="text-base text-slate-500 max-w-2xl">探索、测试并一键集成数以千计的专业 AI 技能。专为开发者打造的高效智能体生态系统。</p>
            
            <div className="w-full max-w-xl mt-4 relative flex items-center shadow-sm hover:shadow-md transition-shadow duration-300 rounded-lg bg-white border border-neutral-200 focus-within:border-blue-600 focus-within:ring-2 focus-within:ring-blue-50/50">
              <span className="material-symbols-outlined absolute left-4 text-neutral-400">search</span>
              <input 
                type="text" 
                className="w-full py-4 pl-12 pr-32 bg-transparent border-none text-slate-900 text-sm focus:ring-0 outline-none" 
                placeholder="寻找特定技能，例如 'Python 代码重构' 或 'SQL 优化'..." 
              />
              <Link to="/search" className="absolute right-2 px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded hover:bg-blue-700 transition-colors">
                搜索技能
              </Link>
            </div>
          </div>
        </section>

        {/* Stats Grid */}
        <section className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {[
            { label: '可用技能总量', value: '12,480+', change: '12%', icon: 'extension', color: 'text-blue-600', bg: 'bg-blue-50' },
            { label: '月活跃开发者', value: '85.2k+', change: '8%', icon: 'group', color: 'text-blue-600', bg: 'bg-blue-50' },
            { label: '累计 API 调用', value: '1.2B', changeText: '历史总计', icon: 'api', color: 'text-blue-600', bg: 'bg-blue-50' },
            { label: 'VS Code 插件安装', value: '150k+', change: '24%', icon: 'download', color: 'text-blue-600', bg: 'bg-blue-50' },
          ].map((stat, idx) => (
            <div key={idx} className="bg-white p-6 rounded-lg border border-neutral-200 shadow-sm hover:shadow-md hover:border-blue-100 transition-all">
              <div className="flex items-center justify-between mb-4">
                <span className="text-sm text-slate-500 font-medium">{stat.label}</span>
                <div className={`p-2 rounded-md ${stat.bg} ${stat.color}`}>
                  <span className="material-symbols-outlined text-[20px]">{stat.icon}</span>
                </div>
              </div>
              <div className="flex items-baseline gap-2">
                <span className="text-2xl font-bold text-slate-900">{stat.value}</span>
                {stat.change ? (
                  <span className="text-xs text-green-500 flex items-center">
                    <span className="material-symbols-outlined text-[14px]">arrow_upward</span> {stat.change}
                  </span>
                ) : (
                  <span className="text-xs text-slate-500">{stat.changeText}</span>
                )}
              </div>
            </div>
          ))}
        </section>

        {/* Main Content Grid: Trending Left, Timeline Right */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <section className="lg:col-span-2 flex flex-col gap-6">
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-bold text-slate-900 flex items-center gap-2">
                <span className="material-symbols-outlined text-amber-500 fill">local_fire_department</span>
                热门技能榜
              </h2>
              <a href="#" className="text-sm text-blue-600 hover:text-blue-700 font-medium flex items-center">
                查看全部 <span className="material-symbols-outlined text-[16px]">chevron_right</span>
              </a>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {[
                { title: 'JSON Schema 生成器', author: 'DataCraft_AI', icon: 'data_object', iconColor: 'text-blue-600', iconBg: 'bg-blue-50', tags: ['数据处理', 'JSON'], desc: '自动将复杂的 JSON 数据结构推断并转换为符合标准规范的 JSON Schema，支持嵌套对象处理。' },
                { title: 'Python 代码重构优化', author: 'CodeGuru', icon: 'terminal', iconColor: 'text-green-600', iconBg: 'bg-green-50', tags: ['代码生成', 'Python'], desc: '遵循 PEP 8 规范，自动识别代码异味并重构 Python 脚本，提升可读性与执行效率。' },
                { title: 'SQL 慢查询分析器', author: 'DB_Doctor', icon: 'analytics', iconColor: 'text-purple-600', iconBg: 'bg-purple-50', tags: ['数据库', 'SQL'], desc: '解析慢查询日志，自动提供索引优化建议和查询重写方案，支持 MySQL 和 PostgreSQL。' },
                { title: 'i18n 多语言自动提取', author: 'FrontEnd_Ninja', icon: 'language', iconColor: 'text-orange-600', iconBg: 'bg-orange-50', tags: ['前端', '工具'], desc: '扫描 React/Vue 组件库，自动提取硬编码文本并生成 i18n 键值对文件，大幅缩短国际化工作流。' },
              ].map((skill, idx) => (
                <Link to="/skill/1" key={idx} className="bg-white border border-neutral-200 rounded-lg p-5 hover:shadow-md hover:border-blue-200 transition-all group flex flex-col h-full relative overflow-hidden">
                  <div className="flex items-start justify-between mb-3 relative z-10">
                    <div className="flex items-center gap-3">
                      <div className={`w-10 h-10 rounded flex items-center justify-center border border-white mix-blend-multiply ${skill.iconBg} ${skill.iconColor}`}>
                        <span className="material-symbols-outlined">{skill.icon}</span>
                      </div>
                      <div>
                        <h3 className="font-bold text-base text-slate-900 group-hover:text-blue-600 transition-colors">{skill.title}</h3>
                        <p className="text-xs text-slate-500">by {skill.author}</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-1 bg-slate-50 px-2 py-1 rounded border border-neutral-200">
                      <span className="material-symbols-outlined text-[14px] text-amber-500 fill">star</span>
                      <span className="font-mono text-[12px] font-medium text-slate-700">4.9</span>
                    </div>
                  </div>
                  <p className="text-sm text-slate-600 line-clamp-2 mb-4 flex-1">
                    {skill.desc}
                  </p>
                  <div className="flex items-center justify-between mt-auto pt-4 border-t border-neutral-100">
                    <div className="flex gap-2">
                      {skill.tags.map(tag => (
                        <span key={tag} className="px-2 py-1 bg-slate-50 text-slate-600 rounded font-mono text-[11px] border border-neutral-200">{tag}</span>
                      ))}
                    </div>
                    <span className="text-blue-600 text-sm font-medium opacity-0 group-hover:opacity-100 transition-opacity flex items-center gap-1">
                      安装 <span className="material-symbols-outlined text-[16px]">download</span>
                    </span>
                  </div>
                </Link>
              ))}
            </div>
          </section>

          {/* Timeline */}
          <section className="lg:col-span-1 flex flex-col gap-6">
            <h2 className="text-xl font-bold text-slate-900 flex items-center gap-2">
              <span className="material-symbols-outlined text-green-500">update</span>
              最新收录
            </h2>
            <div className="bg-white border border-neutral-200 rounded-lg p-5">
              <div className="relative border-l-2 border-slate-100 ml-3 flex flex-col gap-6 py-2">
                {[
                  { title: 'Docker Compose 生成器', time: '10分钟前', desc: 'DevOps / 自动化部署', active: true },
                  { title: 'React 组件测试桩代码', time: '1小时前', desc: '前端 / 单元测试' },
                  { title: 'Git 提交信息美化', time: '3小时前', desc: '工作流 / 工具' },
                  { title: 'Nginx 伪静态转换', time: '5小时前', desc: '运维 / 配置' },
                ].map((item, idx) => (
                  <div key={idx} className="relative pl-6">
                    <div className={`absolute w-3 h-3 rounded-full -left-[7px] top-1.5 ring-4 ring-white ${item.active ? 'bg-blue-600' : 'bg-slate-300'}`}></div>
                    <div className="flex flex-col gap-1">
                      <div className="flex items-baseline justify-between">
                        <a href="#" className="text-sm font-medium text-slate-900 hover:text-blue-600 transition-colors">{item.title}</a>
                        <span className="text-[11px] text-slate-500">{item.time}</span>
                      </div>
                      <p className="text-xs text-slate-500 line-clamp-1">{item.desc}</p>
                    </div>
                  </div>
                ))}
              </div>
              <a href="#" className="block text-center w-full mt-4 py-2 text-blue-600 text-sm font-medium bg-blue-50 hover:bg-blue-100 rounded transition-colors">
                查看动态流
              </a>
            </div>
          </section>
        </div>

        {/* Discover by Concept Row */}
        <section className="flex flex-col gap-6 pt-8 border-t border-neutral-200 mt-4">
          <h2 className="text-2xl font-bold text-slate-900">按领域探索</h2>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
            {[
              { label: '文档生成', icon: 'description' },
              { label: '代码编写', icon: 'code' },
              { label: '数据分析', icon: 'database' },
              { label: '设计辅助', icon: 'palette' },
              { label: '自动化流', icon: 'robot_2' },
              { label: '本地化翻译', icon: 'translate' },
            ].map(cat => (
              <a key={cat.label} href="#" className="flex flex-col items-center justify-center p-6 bg-white border border-neutral-200 rounded-lg hover:border-blue-600 hover:shadow-md group transition-all">
                <span className="material-symbols-outlined text-[32px] text-slate-400 group-hover:text-blue-600 mb-3 transition-colors">{cat.icon}</span>
                <span className="text-sm font-medium text-slate-900 group-hover:text-blue-600 transition-colors">{cat.label}</span>
              </a>
            ))}
          </div>
        </section>

      </main>
    </div>
  );
}
