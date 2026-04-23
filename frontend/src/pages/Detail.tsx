import { Link } from "react-router-dom";

export function Detail() {
  return (
    <div className="w-full max-w-[1280px] mx-auto px-6 py-6 pb-24 lg:pb-6 flex flex-col gap-6">
      <nav className="flex text-xs text-slate-500 items-center gap-2">
        <Link to="/" className="hover:text-blue-600 transition-colors">首页</Link>
        <span className="material-symbols-outlined text-[16px]">chevron_right</span>
        <a href="#" className="hover:text-blue-600 transition-colors">数据分析</a>
        <span className="material-symbols-outlined text-[16px]">chevron_right</span>
        <span className="text-slate-900 font-medium">Excel 趋势分析器</span>
      </nav>

      <div className="bg-white border border-slate-200 rounded-lg p-6 shadow-sm flex flex-col md:flex-row md:items-start justify-between gap-6">
        <div className="flex gap-6">
          <div className="w-24 h-24 rounded-lg bg-slate-50 flex items-center justify-center border border-slate-200 shrink-0">
            <span className="material-symbols-outlined text-teal-600 text-[48px] fill">finance</span>
          </div>
          <div className="flex flex-col justify-center">
            <div className="flex items-center gap-3 mb-2">
              <h1 className="text-3xl font-bold text-slate-900">Excel 趋势分析器</h1>
              <span className="bg-green-50 text-green-600 border border-green-200 px-2.5 py-0.5 rounded-full text-xs flex items-center gap-1 font-medium">
                <span className="material-symbols-outlined text-[14px]">verified_user</span> 已通过安全扫描
              </span>
            </div>
            
            <div className="flex flex-wrap items-center gap-4 text-sm text-slate-500">
              <div className="flex items-center gap-1.5">
                <span className="material-symbols-outlined text-[18px]">person</span>
                <span className="font-medium text-slate-900">Anthropic</span>
              </div>
              <div className="w-1 h-1 rounded-full bg-slate-300"></div>
              <div className="flex items-center gap-1.5">
                <span className="material-symbols-outlined text-[18px]">tag</span>
                <span>v1.2.0</span>
              </div>
              <div className="w-1 h-1 rounded-full bg-slate-300"></div>
              <div className="flex items-center gap-1.5">
                <span className="material-symbols-outlined text-[18px]">download</span>
                <span>12.5k 安装</span>
              </div>
              <div className="w-1 h-1 rounded-full bg-slate-300"></div>
              <div className="flex items-center gap-1">
                <span className="material-symbols-outlined text-[16px] text-amber-500 fill">star</span>
                <span className="material-symbols-outlined text-[16px] text-amber-500 fill">star</span>
                <span className="material-symbols-outlined text-[16px] text-amber-500 fill">star</span>
                <span className="material-symbols-outlined text-[16px] text-amber-500 fill">star</span>
                <span className="material-symbols-outlined text-[16px] text-amber-500">star_half</span>
                <span className="ml-1">(4.8)</span>
              </div>
            </div>

            <p className="mt-3 text-base text-slate-600 max-w-2xl leading-relaxed">
              自动分析 Excel 数据集，识别季节性趋势、异常值，并生成执行摘要报告。完美集成至您的数据科学工作流。
            </p>
          </div>
        </div>

        <div className="flex flex-col gap-3 shrink-0 md:items-end w-full md:w-auto">
          <button className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2.5 rounded font-medium transition-colors flex items-center justify-center gap-2 shadow-sm w-full md:w-auto">
            <span className="material-symbols-outlined text-[20px]">code_blocks</span>
            在 VS Code 中安装
          </button>
          <div className="flex gap-2 w-full md:w-auto">
            <button className="flex-1 md:flex-none border border-slate-300 hover:bg-slate-50 text-slate-900 px-4 py-2.5 rounded font-medium transition-colors flex items-center justify-center gap-2 text-sm shadow-sm">
              <span className="material-symbols-outlined text-[18px]">content_copy</span>
              复制安装命令
            </button>
            <button className="border border-slate-300 hover:bg-slate-50 text-slate-900 p-2.5 rounded transition-colors flex items-center justify-center shadow-sm">
              <span className="material-symbols-outlined text-[20px]">star</span>
            </button>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div className="lg:col-span-9 flex flex-col gap-4">
          <div className="flex border-b border-slate-200">
            <button className="px-4 py-3 border-b-2 border-blue-600 text-blue-600 font-medium text-sm flex items-center gap-2">
              <span className="material-symbols-outlined text-[18px]">description</span>
              SKILL.md
            </button>
            <button className="px-4 py-3 border-b-2 border-transparent text-slate-500 hover:text-slate-900 hover:bg-slate-50 transition-colors font-medium text-sm flex items-center gap-2">
              <span className="material-symbols-outlined text-[18px]">folder</span>
              文件目录
            </button>
            <button className="px-4 py-3 border-b-2 border-transparent text-slate-500 hover:text-slate-900 hover:bg-slate-50 transition-colors font-medium text-sm flex items-center gap-2">
              <span className="material-symbols-outlined text-[18px]">info</span>
              详情
            </button>
          </div>

          <div className="bg-[#1e1e1e] rounded-lg border border-slate-800 p-6 shadow-sm overflow-hidden text-slate-300 relative group">
             <button className="absolute top-4 right-4 text-slate-400 hover:text-white bg-slate-800 p-1.5 rounded opacity-0 group-hover:opacity-100 transition-opacity flex items-center">
                <span className="material-symbols-outlined text-[16px]">content_copy</span>
             </button>
             <h2 className="text-2xl font-bold text-white mb-4 border-b border-slate-700 pb-2">Excel Trend Analyzer</h2>
             <p className="mb-4 text-slate-400">A powerful skill to automatically analyze trends in large Excel datasets using advanced LLM reasoning.</p>
             <h3 className="text-xl font-bold text-white mt-6 mb-3">Installation</h3>
             <div className="bg-black/50 border border-slate-700 rounded p-4 mb-4 font-mono text-sm text-blue-300">
               skillhub install anthropic/excel-trend
             </div>
             <h3 className="text-xl font-bold text-white mt-6 mb-3">Usage Example</h3>
             <div className="bg-black/50 border border-slate-700 rounded p-4 font-mono text-sm leading-relaxed overflow-x-auto text-slate-300">
<pre>
<span className="text-purple-400">import</span> {"{ SkillHub }"} <span className="text-purple-400">from</span> <span className="text-green-300">'@skillhub/core'</span>;

<span className="text-slate-500">{"// Initialize the analyzer"}</span>
<span className="text-blue-400">const</span> analyzer = <span className="text-blue-400">new</span> <span className="text-yellow-200">SkillHub</span>(<span className="text-green-300">'anthropic/excel-trend'</span>);

<span className="text-slate-500">{"// Run analysis on sales data"}</span>
<span className="text-blue-400">const</span> results = <span className="text-purple-400">await</span> analyzer.<span className="text-yellow-200">execute</span>({`{
  `}<span className="text-blue-300">file</span>: <span className="text-green-300">'./Q3_Sales.xlsx'</span>,
  <span className="text-blue-300">targetColumn</span>: <span className="text-green-300">'Revenue'</span>,
  <span className="text-blue-300">detectAnomalies</span>: <span className="text-blue-400">true</span>
{`}`});

console.<span className="text-yellow-200">log</span>(results.<span className="text-blue-300">insights</span>);
</pre>
             </div>
          </div>
        </div>

        <aside className="lg:col-span-3 flex flex-col gap-4">
           <h3 className="text-lg font-bold text-slate-900 flex items-center gap-2 mb-2">
             <span className="material-symbols-outlined text-blue-600">psychology</span>
             相似技能推荐
           </h3>
           <div className="flex flex-col gap-3">
             {[
               { title: 'CSV 数据清洗', author: 'Google', v: 'v2.1', icon: 'table_chart', color: 'text-green-600', bg: 'bg-green-50', desc: 'Automatically detects and fixes formatting issues...', s: '4.8', d: '8.2k' },
               { title: 'SQL 查询生成器', author: 'OpenAI', v: 'v1.0', icon: 'query_stats', color: 'text-purple-600', bg: 'bg-purple-50', desc: 'Convert natural language questions into optimized SQL.', s: '4.9', d: '21k' },
               { title: '图表自动生成', author: 'Anthropic', v: 'v0.9', icon: 'insights', color: 'text-orange-600', bg: 'bg-orange-50', desc: 'Analyzes data arrays and suggests best charts.', s: '4.5', d: '3.4k' },
             ].map((skill, idx) => (
               <a key={idx} href="#" className="bg-white border border-slate-200 rounded-lg p-4 hover:border-blue-600 hover:shadow-md transition-all group flex flex-col gap-2">
                 <div className="flex items-start justify-between">
                   <div className="flex items-center gap-2">
                     <div className={`w-8 h-8 rounded flex items-center justify-center shrink-0 ${skill.bg} ${skill.color}`}>
                       <span className="material-symbols-outlined text-[18px]">{skill.icon}</span>
                     </div>
                     <div>
                       <h4 className="font-medium text-slate-900 text-sm group-hover:text-blue-600 transition-colors">{skill.title}</h4>
                       <p className="text-xs text-slate-500">{skill.author}</p>
                     </div>
                   </div>
                   <span className="text-xs text-slate-400 bg-slate-100 px-1.5 py-0.5 rounded font-mono">{skill.v}</span>
                 </div>
                 <p className="text-xs text-slate-600 line-clamp-2">{skill.desc}</p>
                 <div className="flex items-center gap-3 text-xs text-slate-500 mt-1">
                   <span className="flex items-center gap-1"><span className="material-symbols-outlined text-[14px]">star</span> {skill.s}</span>
                   <span className="flex items-center gap-1"><span className="material-symbols-outlined text-[14px]">download</span> {skill.d}</span>
                 </div>
               </a>
             ))}
           </div>
        </aside>
      </div>
    </div>
  );
}
