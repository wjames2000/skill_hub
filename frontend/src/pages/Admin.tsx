import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { adminApi } from "../lib/api/admin";
import { statsApi } from "../lib/api/stats";
import { ErrorBanner } from "../components/ui/ErrorBanner";
import { cn } from "../lib/utils";
import type { Stats, SyncTask } from "../types";

type AdminTab = 'dashboard' | 'sync' | 'review' | 'users' | 'logs';

export function Admin() {
  const [activeTab, setActiveTab] = useState<AdminTab>('dashboard');
  const [stats, setStats] = useState<Stats | null>(null);
  const [tasks, setTasks] = useState<SyncTask[]>([]);
  const [pendingSkills, setPendingSkills] = useState<{ id: number; title: string; author: string; time: string }[]>([]);
  const [logs, setLogs] = useState<{ timestamp: string; level: string; message: string }[]>([]);
  const [users, setUsers] = useState<{ id: number; username: string; email: string; role: string; createdAt: string }[]>([]);
  const [triggering, setTriggering] = useState<'full' | 'incremental' | null>(null);
  const [topError, setTopError] = useState<string | null>(null);
  const [sidebarOpen, setSidebarOpen] = useState(false);

  const handleApprove = async (id: number) => {
    try {
      await adminApi.approveSkill(id);
      setPendingSkills(prev => prev.filter(s => s.id !== id));
    } catch { setTopError('审核操作失败'); }
  };

  const handleReject = async (id: number) => {
    try {
      await adminApi.rejectSkill(id);
      setPendingSkills(prev => prev.filter(s => s.id !== id));
    } catch { setTopError('审核操作失败'); }
  };

  const tabs: { key: AdminTab; label: string; icon: string }[] = [
    { key: 'dashboard', label: '仪表盘', icon: 'dashboard' },
    { key: 'sync', label: '同步任务', icon: 'pest_control' },
    { key: 'review', label: '技能审核', icon: 'verified' },
    { key: 'users', label: '用户管理', icon: 'group' },
    { key: 'logs', label: '系统日志', icon: 'terminal' },
  ];

  useEffect(() => {
    statsApi.getOverview().then(setStats).catch(() => {});
    adminApi.getSyncTasks().then(setTasks).catch(() => {});
    adminApi.getPendingReviews().then(skills => {
      setPendingSkills(skills.map(s => ({
        id: s.id,
        title: s.title,
        author: s.author,
        time: '刚刚',
      })));
    }).catch(() => {});
    adminApi.getSystemLogs(50).then(setLogs).catch(() => {});
    adminApi.getUsers().then(res => setUsers(res.data)).catch(() => {});
  }, []);

  const handleTriggerSync = async (type: 'full' | 'incremental') => {
    setTriggering(type);
    try {
      await adminApi.triggerSync(type);
      const updated = await adminApi.getSyncTasks();
      setTasks(updated);
    } catch { setTopError('触发同步任务失败，请检查后端服务是否正常运行'); }
    setTriggering(null);
  };

  const renderContent = () => {
    switch (activeTab) {
      case 'dashboard':
        return (
          <div className="flex flex-col gap-6">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6">
              {[
                { label: '技能库总量', value: stats?.totalSkills?.toLocaleString() ?? '-', change: '12%', icon: 'library_books', color: 'text-brand-600', bg: 'bg-brand-50' },
                { label: '今日新增', value: stats?.todayNew != null ? String(stats.todayNew) : '-', change: '5%', icon: 'add_circle', color: 'text-brand-600', bg: 'bg-brand-50' },
                { label: 'API 24h 调用量', value: stats?.api24hCalls ?? '-', tags: ['正常'], icon: 'api', color: 'text-amber-500', bg: 'bg-amber-50' },
                { label: '当前爬虫状态', value: stats?.crawlerRunning != null ? String(stats.crawlerRunning) : '-', suffix: '个任务进行中', icon: 'bug_report', color: 'text-green-500', bg: 'bg-green-50' },
              ].map((m, idx) => (
                <div key={idx} className="card p-5 md:p-6">
                  <div className="flex items-center justify-between mb-3 md:mb-4">
                    <span className="text-sm text-slate-500 font-medium">{m.label}</span>
                    <span className={`material-symbols-outlined p-1.5 rounded-md ${m.color} ${m.bg}`}>{m.icon}</span>
                  </div>
                  <div className="flex items-baseline gap-2">
                    <span className="text-2xl md:text-3xl font-bold text-slate-900">{m.value}</span>
                    {'change' in m && m.change && (
                      <span className="text-green-500 text-sm font-medium flex items-center">
                        <span className="material-symbols-outlined text-[16px]">arrow_upward</span>{m.change}
                      </span>
                    )}
                    {'tags' in m && m.tags?.map(tag => (
                      <span key={tag} className="badge-green text-[11px]">{tag}</span>
                    ))}
                    {'suffix' in m && m.suffix && <span className="text-slate-500 text-sm">{m.suffix}</span>}
                  </div>
                </div>
              ))}
            </div>

            <div className="card p-6">
              <div className="flex items-center justify-between mb-6">
                <h3 className="text-lg md:text-xl font-bold text-slate-900">近 7 日趋势</h3>
              </div>
              <div className="h-48 md:h-64 w-full flex items-end gap-2 pt-8 relative border-b border-l border-slate-200 pl-4 pb-4 overflow-x-auto">
                <div className="absolute left-[-30px] bottom-4 top-8 flex flex-col justify-between text-xs text-slate-400 h-full">
                  <span>10k</span>
                  <span>5k</span>
                  <span>0</span>
                </div>
                {['周一','周二','周三','周四','周五','周六','周日'].map((day, idx) => (
                  <div key={day} className="flex-1 min-w-[40px] flex items-end justify-around h-full relative group">
                    <div className="w-3 md:w-4 bg-brand-600/30 rounded-t hover:bg-brand-600 transition-colors cursor-pointer"
                      style={{ height: `${[40, 45, 35, 50, 60, 75, 85][idx]}%` }}
                    />
                    <div className="w-3 md:w-4 bg-brand-400/50 rounded-t ml-0.5 md:ml-1 hover:bg-brand-400 transition-colors cursor-pointer"
                      style={{ height: `${[60, 65, 50, 70, 75, 85, 95][idx]}%` }}
                    />
                    <span className="absolute -bottom-6 text-[10px] md:text-xs text-slate-400 whitespace-nowrap">{day}</span>
                  </div>
                ))}
                <div className="absolute top-0 right-0 flex gap-3 md:gap-4 text-xs font-medium text-slate-500">
                  <div className="flex items-center gap-1"><div className="w-3 h-3 bg-brand-600 rounded-sm" />新增技能</div>
                  <div className="flex items-center gap-1"><div className="w-3 h-3 bg-brand-400 rounded-sm" />API 调用量</div>
                </div>
              </div>
            </div>
          </div>
        );

      case 'sync':
        return (
          <div className="flex flex-col gap-6">
            <div className="flex items-center justify-between">
              <h3 className="text-lg md:text-xl font-bold text-slate-900">同步任务管理</h3>
              <div className="flex gap-2">
                <button
                  onClick={() => handleTriggerSync('incremental')}
                  disabled={triggering !== null}
                  className="btn-secondary text-sm"
                >
                  {triggering === 'incremental' ? '触发中...' : '增量同步'}
                </button>
                <button
                  onClick={() => handleTriggerSync('full')}
                  disabled={triggering !== null}
                  className="btn-primary text-sm"
                >
                  {triggering === 'full' ? '触发中...' : '全量同步'}
                </button>
              </div>
            </div>
            <div className="card">
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-slate-200 bg-slate-50/50">
                      <th className="text-left px-4 py-3 font-medium text-slate-600">任务名称</th>
                      <th className="text-left px-4 py-3 font-medium text-slate-600">类型</th>
                      <th className="text-left px-4 py-3 font-medium text-slate-600">状态</th>
                      <th className="text-left px-4 py-3 font-medium text-slate-600">进度</th>
                      <th className="text-left px-4 py-3 font-medium text-slate-600">启动时间</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100">
                    {tasks.map(task => (
                      <tr key={task.id} className="hover:bg-slate-50 transition-colors">
                        <td className="px-4 py-3 font-medium text-slate-900">{task.name}</td>
                        <td className="px-4 py-3">
                          <span className="badge bg-slate-100 text-slate-600 text-[11px]">
                            {task.type === 'full' ? '全量' : '增量'}
                          </span>
                        </td>
                        <td className="px-4 py-3">
                          <span className={`inline-flex items-center gap-1 text-xs font-medium px-2 py-0.5 rounded-full ${
                            task.status === 'running' ? 'bg-blue-50 text-blue-700' :
                            task.status === 'completed' ? 'bg-green-50 text-green-700' :
                            task.status === 'failed' ? 'bg-red-50 text-red-700' :
                            'bg-slate-50 text-slate-600'
                          }`}>
                            {task.status === 'running' && <span className="w-1.5 h-1.5 rounded-full bg-blue-500 animate-pulse" />}
                            {task.status === 'running' ? '运行中' : task.status === 'completed' ? '已完成' : task.status === 'failed' ? '失败' : '等待中'}
                          </span>
                        </td>
                        <td className="px-4 py-3">
                          <div className="flex items-center gap-2">
                            <div className="w-24 bg-slate-200 rounded-full h-1.5">
                              <div className={`h-1.5 rounded-full ${
                                task.status === 'completed' ? 'bg-green-500' :
                                task.status === 'failed' ? 'bg-red-500' : 'bg-brand-600'
                              }`} style={{ width: `${task.progress}%` }} />
                            </div>
                            <span className="text-xs text-slate-500">{task.progress}%</span>
                          </div>
                        </td>
                        <td className="px-4 py-3 text-xs text-slate-500">{task.startedAt || '-'}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        );

      case 'review':
        return (
          <div className="card">
            <div className="p-4 border-b border-slate-200">
              <h3 className="text-lg md:text-xl font-bold text-slate-900">待审核技能</h3>
            </div>
            <div className="divide-y divide-slate-100">
              {pendingSkills.map(item => (
                <div key={item.id} className="p-4 hover:bg-slate-50 transition-colors flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 group">
                  <div>
                    <h4 className="font-medium text-slate-900 group-hover:text-brand-600 transition-colors">{item.title}</h4>
                    <div className="flex items-center gap-3 text-xs text-slate-500 mt-0.5">
                      <span className="flex items-center gap-1"><span className="material-symbols-outlined text-[14px]">person</span>{item.author}</span>
                      <span className="flex items-center gap-1"><span className="material-symbols-outlined text-[14px]">schedule</span>{item.time}</span>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <button onClick={() => handleApprove(item.id)} className="px-3 py-1.5 text-xs font-medium text-green-600 border border-green-200 bg-white hover:bg-green-50 rounded transition-colors">通过</button>
                    <button onClick={() => handleReject(item.id)} className="px-3 py-1.5 text-xs font-medium text-red-600 border border-red-200 bg-white hover:bg-red-50 rounded transition-colors">驳回</button>
                  </div>
                </div>
              ))}
              {pendingSkills.length === 0 && !stats && (
                <div className="p-8 text-center text-sm text-slate-500">暂无待审核技能</div>
              )}
            </div>
          </div>
        );

      case 'users':
        return (
          <div className="card">
            <div className="p-4 border-b border-slate-200">
              <h3 className="text-lg md:text-xl font-bold text-slate-900">用户管理</h3>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200 bg-slate-50/50">
                    <th className="text-left px-4 py-3 font-medium text-slate-600">ID</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">用户名</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600 hidden md:table-cell">邮箱</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">角色</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600 hidden md:table-cell">注册时间</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {users.map(user => (
                    <tr key={user.id} className="hover:bg-slate-50 transition-colors">
                      <td className="px-4 py-3 text-slate-500">#{user.id}</td>
                      <td className="px-4 py-3 font-medium text-slate-900">{user.username}</td>
                      <td className="px-4 py-3 text-slate-600 hidden md:table-cell">{user.email}</td>
                      <td className="px-4 py-3">
                        <span className={`badge text-[11px] ${user.role === 'admin' ? 'badge-blue' : 'bg-slate-50 text-slate-600 border border-slate-200'}`}>
                          {user.role === 'admin' ? '管理员' : '用户'}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-xs text-slate-500 hidden md:table-cell">{user.createdAt}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        );

      case 'logs':
        return (
          <div className="card">
            <div className="p-4 border-b border-slate-200 flex justify-between items-center">
              <h3 className="text-lg md:text-xl font-bold text-slate-900">系统日志</h3>
              <div className="flex items-center gap-2 border border-slate-200 bg-white px-2 py-1 rounded-full shadow-sm text-xs">
                <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
                <span className="text-slate-600 font-medium">实时</span>
              </div>
            </div>
            <div className="p-4 bg-[#1e1e1e] rounded-b-lg font-mono text-xs md:text-sm text-slate-300 max-h-[600px] overflow-y-auto">
              {logs.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-12 text-slate-500">
                  <span className="material-symbols-outlined text-[40px] mb-2">terminal</span>
                  <p className="text-sm">暂无系统日志</p>
                  <p className="text-xs mt-1 text-slate-600">日志数据将在 sync-worker 运行后显示</p>
                </div>
              ) : (
                logs.map((log, idx) => (
                <div key={idx} className="flex gap-2 md:gap-3 items-start py-0.5">
                  <span className="text-slate-500 whitespace-nowrap">[{log.timestamp}]</span>
                  <span className={`${
                    log.level === 'INFO' ? 'text-green-400' :
                    log.level === 'WARN' ? 'text-amber-400' :
                    log.level === 'SYS' ? 'text-blue-400' : 'text-slate-400'
                  }`}>[{log.level}]</span>
                  <span className="text-slate-200 break-all">{log.message}</span>
                </div>
              ))
              )}
            </div>
          </div>
        );
    }
  };

  return (
    <div className="font-sans text-sm text-slate-900 bg-slate-50 min-h-screen">
      <header className="bg-white/80 backdrop-blur-md fixed top-0 right-0 left-0 flex items-center justify-between px-4 md:px-6 h-16 md:ml-56 lg:ml-64 border-b border-slate-200 z-50">
        <div className="flex items-center gap-4">
          <button
            onClick={() => setSidebarOpen(prev => !prev)}
            className="md:hidden p-2 text-slate-500 hover:bg-slate-50 rounded-lg transition-colors"
            aria-label="Toggle navigation"
          >
            <span className="material-symbols-outlined">{sidebarOpen ? 'close' : 'menu'}</span>
          </button>
          <Link to="/" className="text-lg font-bold tracking-tight text-slate-900 md:hidden">SkillHub Pro</Link>
        </div>
        <div className="flex items-center gap-2 md:gap-4">
          <button className="p-2 text-slate-500 hover:bg-slate-50 rounded-full transition-colors">
            <span className="material-symbols-outlined">notifications</span>
          </button>
          <button className="p-2 text-slate-500 hover:bg-slate-50 rounded-full transition-colors">
            <span className="material-symbols-outlined">settings</span>
          </button>
          <Link to="/" className="p-2 text-slate-500 hover:bg-slate-50 rounded-full transition-colors" title="返回首页">
            <span className="material-symbols-outlined">home</span>
          </Link>
        </div>
      </header>

      {/* Mobile overlay backdrop */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black/30 z-30 md:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      <nav className={`bg-white font-medium h-screen w-56 lg:w-64 border-r border-slate-200 fixed left-0 top-0 bottom-0 flex-col pt-16 z-40 transition-transform duration-200 ${
        sidebarOpen ? 'translate-x-0' : '-translate-x-full'
      } md:translate-x-0 md:flex`}>
        <div className="p-4 lg:p-6 pb-2">
          <Link to="/" className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-brand-100 flex items-center justify-center">
              <span className="material-symbols-outlined text-brand-600">admin_panel_settings</span>
            </div>
            <div>
              <h2 className="text-xl font-black tracking-tighter text-brand-600 leading-tight">SkillHub Pro</h2>
              <span className="text-xs text-slate-500">Admin Console</span>
            </div>
          </Link>
        </div>
        <div className="flex-1 overflow-y-auto py-4 px-2 lg:px-3 flex flex-col gap-1">
          {tabs.map(tab => (
            <button
              key={tab.key}
              onClick={() => { setActiveTab(tab.key); setSidebarOpen(false); }}
              aria-current={activeTab === tab.key ? 'page' : undefined}
              className={cn(
                "flex items-center gap-3 px-3 py-2 rounded-lg transition-all",
                activeTab === tab.key
                  ? 'bg-brand-50 text-brand-700 border-r-2 border-brand-600'
                  : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'
              )}
            >
              <span className="material-symbols-outlined text-[20px]">{tab.icon}</span>
              <span className="text-sm">{tab.label}</span>
            </button>
          ))}
        </div>
        <div className="p-2 lg:p-3 border-t border-slate-200">
          <Link to="/" className="flex items-center gap-3 px-2 lg:px-3 py-2 rounded-lg text-slate-600 hover:bg-slate-100 hover:text-slate-900 transition-all text-sm">
            <span className="material-symbols-outlined text-[20px]">arrow_back</span>
            <span>返回前台</span>
          </Link>
        </div>
      </nav>

      <main className="md:ml-56 lg:ml-64 pt-16 min-h-screen px-4 md:px-6 pb-4 md:pb-6 max-w-[1440px] mx-auto">
        <ErrorBanner message={topError} onDismiss={() => setTopError(null)} />
        <div className="mb-6 mt-4 md:mt-6">
          <h1 className="text-xl md:text-2xl font-bold text-slate-900 mb-1">
            {tabs.find(t => t.key === activeTab)?.label || '管理后台'}
          </h1>
          <p className="text-slate-500 text-xs md:text-sm">系统管理与监控面板</p>
        </div>
        {renderContent()}
      </main>
    </div>
  );
}
