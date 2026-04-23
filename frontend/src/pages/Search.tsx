import { Link } from "react-router-dom";
import { useState } from "react";

export function Search() {
  return (
    <div className="w-full max-w-[1440px] mx-auto px-6 py-8 flex gap-8 items-start">
      {/* Left Sidebar Filters */}
      <aside className="w-[260px] flex-shrink-0 flex flex-col gap-6 sticky top-[80px]">
        <div className="bg-white rounded-lg border border-slate-200 p-4 flex flex-col gap-6 shadow-sm">
          <div className="flex items-center justify-between border-b border-slate-100 pb-2">
            <h2 className="font-semibold text-lg text-slate-900">筛选</h2>
            <button className="text-blue-600 text-sm hover:text-blue-700 transition-colors">重置</button>
          </div>

          <div className="flex flex-col gap-3">
            <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider">分类</h3>
            <div className="flex flex-col gap-2 text-sm">
              <label className="flex items-center gap-2 cursor-pointer group">
                <input type="checkbox" defaultChecked className="w-4 h-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500" />
                <span className="text-slate-900 group-hover:text-blue-600">数据处理</span>
                <span className="text-xs text-slate-400 ml-auto">128</span>
              </label>
              <label className="flex items-center gap-2 cursor-pointer group">
                <input type="checkbox" className="w-4 h-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500" />
                <span className="text-slate-900 group-hover:text-blue-600">可视化</span>
                <span className="text-xs text-slate-400 ml-auto">45</span>
              </label>
              <label className="flex items-center gap-2 cursor-pointer group">
                <input type="checkbox" className="w-4 h-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500" />
                <span className="text-slate-900 group-hover:text-blue-600">自动化脚本</span>
                <span className="text-xs text-slate-400 ml-auto">89</span>
              </label>
            </div>
          </div>

          <div className="flex flex-col gap-3">
            <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider">来源</h3>
            <div className="flex flex-col gap-2 text-sm">
              <label className="flex items-center gap-2 cursor-pointer group">
                <input type="checkbox" defaultChecked className="w-4 h-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500" />
                <span className="text-slate-900 group-hover:text-blue-600">官方认证 (Official)</span>
              </label>
              <label className="flex items-center gap-2 cursor-pointer group">
                <input type="checkbox" className="w-4 h-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500" />
                <span className="text-slate-900 group-hover:text-blue-600">社区开源 (GitHub)</span>
              </label>
            </div>
          </div>

          <div className="flex flex-col gap-3">
            <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider">安全状态</h3>
            <div className="flex flex-col gap-2 text-sm">
              <label className="flex items-center gap-2 cursor-pointer group">
                <input type="checkbox" defaultChecked className="w-4 h-4 rounded border-green-300 text-green-500 focus:ring-green-500" />
                <span className="material-symbols-outlined text-[16px] text-green-500">verified_user</span>
                <span className="text-slate-900 group-hover:text-blue-600">已扫描安全</span>
              </label>
              <label className="flex items-center gap-2 cursor-pointer group">
                <input type="checkbox" className="w-4 h-4 rounded border-amber-300 text-amber-500 focus:ring-amber-500" />
                <span className="material-symbols-outlined text-[16px] text-amber-500">warning</span>
                <span className="text-slate-900 group-hover:text-blue-600">存在潜在风险</span>
              </label>
            </div>
          </div>

          <div className="flex flex-col gap-3">
            <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider">排序规则</h3>
            <div className="flex flex-col gap-1 bg-slate-50 rounded p-1 text-sm border border-slate-200">
              <label className="flex items-center justify-center py-1.5 bg-white rounded shadow-sm cursor-pointer border border-blue-200">
                <span className="text-blue-600 font-medium">匹配度最高</span>
              </label>
              <label className="flex items-center justify-center py-1.5 hover:bg-slate-100 rounded cursor-pointer text-slate-600">
                <span>评分最高</span>
              </label>
              <label className="flex items-center justify-center py-1.5 hover:bg-slate-100 rounded cursor-pointer text-slate-600">
                <span>下载量最多</span>
              </label>
            </div>
          </div>
        </div>
      </aside>

      {/* Main Content */}
      <section className="flex-1 flex flex-col gap-6">
        <div className="flex flex-col gap-3">
          <h1 className="text-2xl font-bold text-slate-900 flex items-center gap-2">
            为您找到 <span className="text-blue-600">156</span> 个关于 
            <span className="font-mono bg-blue-50 text-blue-700 px-2 py-0.5 rounded border border-blue-100 text-lg mx-1">"分析 Excel"</span> 
            的技能
          </h1>
          <div className="flex items-center gap-2">
            <span className="text-xs text-slate-400">已选筛选器:</span>
            <div className="flex items-center gap-1.5 bg-white text-slate-700 text-xs px-2.5 py-1 flex items-center rounded-full border border-slate-200">
              分类: 数据处理
              <span className="material-symbols-outlined text-[14px] cursor-pointer hover:text-red-500">close</span>
            </div>
            <div className="flex items-center gap-1.5 bg-white text-slate-700 text-xs px-2.5 py-1 flex items-center rounded-full border border-slate-200">
              来源: 官方认证
              <span className="material-symbols-outlined text-[14px] cursor-pointer hover:text-red-500">close</span>
            </div>
          </div>
        </div>

        {/* Suggestion Box */}
        <div className="bg-blue-50/50 border border-blue-100 rounded-lg p-4 flex items-start gap-4 shadow-sm">
          <div className="w-10 h-10 rounded-full bg-blue-600 text-white flex items-center justify-center shrink-0 shadow-sm">
            <span className="material-symbols-outlined fill">lightbulb</span>
          </div>
          <div className="pt-1">
            <h3 className="text-base font-semibold text-blue-900 mb-1">💡 推荐理由</h3>
            <p className="text-sm text-blue-800/80 leading-relaxed">
              基于您的数据分析需求，以下技能在处理结构化数据、生成透视表以及大批量 Excel 文件读取方面表现最佳。它们均已通过代码扫描，可安全集成到您的 VS Code 环境中。
            </p>
          </div>
        </div>

        {/* Result Cards */}
        <div className="flex flex-col gap-4">
          {/* Card 1 */}
          <Link to="/skill/1" className="bg-white border border-slate-200 rounded-xl p-5 hover:shadow-lg hover:border-blue-300 transition-all flex flex-col md:flex-row gap-6 relative group overflow-hidden block">
            <div className="w-16 h-16 rounded-lg bg-indigo-50 border border-indigo-100 flex items-center justify-center shrink-0 group-hover:scale-105 transition-transform">
              <span className="material-symbols-outlined text-indigo-600 text-[32px] fill">table_chart</span>
            </div>
            <div className="flex-1 flex flex-col gap-2">
              <div className="flex justify-between items-center">
                <h3 className="text-xl font-bold text-slate-900 group-hover:text-blue-600 transition-colors">Pandas Excel Master</h3>
                <div className="hidden md:flex items-center gap-1.5 bg-green-50 text-green-600 px-2.5 py-1 rounded-md text-xs font-medium border border-green-200/50">
                  <span className="material-symbols-outlined text-[14px]">shield</span>
                  已扫描安全
                </div>
              </div>
              <div className="flex items-center gap-3 text-xs text-slate-500">
                <span className="flex items-center gap-1 text-blue-600 bg-blue-50 px-1.5 py-0.5 rounded">
                  <span className="material-symbols-outlined text-[14px]">verified</span> Microsoft Official
                </span>
                <span className="w-1 h-1 rounded-full bg-slate-300"></span>
                <span>v2.4.1</span>
                <span className="w-1 h-1 rounded-full bg-slate-300"></span>
                <span className="font-mono bg-slate-100 px-1.5 py-0.5 rounded border border-slate-200">python</span>
              </div>
              <p className="text-sm text-slate-600 line-clamp-2 mt-1 leading-relaxed">
                深度集成 Pandas 库的智能助手，能够通过自然语言指令快速分析庞大的 Excel 文件，自动生成数据清洗脚本和可视化图表代码。支持多表合并与透视分析。
              </p>
            </div>
            <div className="flex flex-row md:flex-col items-center md:items-end justify-between gap-4 md:border-l border-slate-100 md:pl-6 min-w-[180px]">
              <div className="w-full flex justify-between items-center mb-1">
                <span className="text-xs text-slate-400">匹配度</span>
                <span className="text-sm font-bold text-green-500">98%</span>
              </div>
              <div className="w-full h-1.5 bg-slate-100 rounded-full mb-2">
                <div className="h-full bg-green-500 rounded-full" style={{ width: '98%' }}></div>
              </div>
              <div className="flex items-center gap-4 text-xs text-slate-500 mb-2">
                <span className="flex items-center gap-1"><span className="material-symbols-outlined text-[14px]">download</span> 124k</span>
                <span className="flex items-center gap-1 text-amber-500"><span className="material-symbols-outlined text-[14px] fill">star</span> 4.9</span>
              </div>
              <button className="w-full bg-blue-600 text-white text-sm font-medium py-2 rounded shadow-sm hover:bg-blue-700 transition-colors flex justify-center items-center gap-2">
                <span className="material-symbols-outlined text-[16px]">add_box</span>
                在 VS Code 中安装
              </button>
            </div>
          </Link>
          
          {/* Card 2 */}
          <Link to="/skill/2" className="bg-white border border-slate-200 rounded-xl p-5 hover:shadow-lg hover:border-blue-300 transition-all flex flex-col md:flex-row gap-6 relative group overflow-hidden block">
            <div className="w-16 h-16 rounded-lg bg-teal-50 border border-teal-100 flex items-center justify-center shrink-0 group-hover:scale-105 transition-transform">
              <span className="material-symbols-outlined text-teal-600 text-[32px] fill">analytics</span>
            </div>
            <div className="flex-1 flex flex-col gap-2">
              <div className="flex justify-between items-center">
                <h3 className="text-xl font-bold text-slate-900 group-hover:text-blue-600 transition-colors">Excel VBA Auto-Coder</h3>
                <div className="hidden md:flex items-center gap-1.5 bg-green-50 text-green-600 px-2.5 py-1 rounded-md text-xs font-medium border border-green-200/50">
                  <span className="material-symbols-outlined text-[14px]">shield</span>
                  已扫描安全
                </div>
              </div>
              <div className="flex items-center gap-3 text-xs text-slate-500">
                <span className="flex items-center gap-1 bg-slate-50 px-1.5 py-0.5 rounded border border-slate-200">
                  <span className="material-symbols-outlined text-[14px]">person</span> dev_community_x
                </span>
                <span className="w-1 h-1 rounded-full bg-slate-300"></span>
                <span>v1.0.8</span>
                <span className="w-1 h-1 rounded-full bg-slate-300"></span>
                <span className="font-mono bg-slate-100 px-1.5 py-0.5 rounded border border-slate-200">vba</span>
              </div>
              <p className="text-sm text-slate-600 line-clamp-2 mt-1 leading-relaxed">
                专为传统 Excel 用户设计，将自然语言需求直接转化为高质量、带注释的 VBA 宏代码。擅长处理复杂的单元格格式化和跨工作簿数据迁移。
              </p>
            </div>
            <div className="flex flex-row md:flex-col items-center md:items-end justify-between gap-4 md:border-l border-slate-100 md:pl-6 min-w-[180px]">
              <div className="w-full flex justify-between items-center mb-1">
                <span className="text-xs text-slate-400">匹配度</span>
                <span className="text-sm font-bold text-blue-600">85%</span>
              </div>
              <div className="w-full h-1.5 bg-slate-100 rounded-full mb-2">
                <div className="h-full bg-blue-600 rounded-full" style={{ width: '85%' }}></div>
              </div>
              <div className="flex items-center gap-4 text-xs text-slate-500 mb-2">
                <span className="flex items-center gap-1"><span className="material-symbols-outlined text-[14px]">download</span> 56k</span>
                <span className="flex items-center gap-1 text-amber-500"><span className="material-symbols-outlined text-[14px] fill">star</span> 4.6</span>
              </div>
              <button className="w-full bg-white text-blue-600 border border-blue-600 text-sm font-medium py-2 rounded shadow-sm hover:bg-blue-50 transition-colors flex justify-center items-center gap-2">
                <span className="material-symbols-outlined text-[16px]">add_box</span>
                在 VS Code 中安装
              </button>
            </div>
          </Link>
        </div>

        {/* Pagination mock */}
        <div className="flex items-center justify-center gap-2 mt-8">
          <button className="w-8 h-8 rounded border border-slate-200 flex items-center justify-center text-slate-400 disabled" disabled>
            <span className="material-symbols-outlined text-[18px]">chevron_left</span>
          </button>
          <button className="w-8 h-8 rounded bg-blue-600 text-white text-sm font-medium">1</button>
          <button className="w-8 h-8 rounded border border-slate-200 text-slate-600 hover:bg-slate-50 text-sm font-medium transition-colors">2</button>
          <button className="w-8 h-8 rounded border border-slate-200 text-slate-600 hover:bg-slate-50 text-sm font-medium transition-colors">3</button>
          <span className="text-slate-400 px-1">...</span>
          <button className="w-8 h-8 rounded border border-slate-200 text-slate-600 hover:bg-slate-50 flex items-center justify-center transition-colors">
            <span className="material-symbols-outlined text-[18px]">chevron_right</span>
          </button>
        </div>
      </section>
    </div>
  );
}
