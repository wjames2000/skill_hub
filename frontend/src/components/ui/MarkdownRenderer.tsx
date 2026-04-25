import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import type { Components } from "react-markdown";

interface MarkdownRendererProps {
  content: string;
}

const components: Components = {
  h1: ({ children, ...props }) => (
    <h1 className="text-2xl font-bold text-white mb-4 border-b border-slate-700 pb-2" {...props}>
      {children}
    </h1>
  ),
  h2: ({ children, ...props }) => (
    <h2 className="text-xl font-bold text-white mt-6 mb-3" {...props}>
      {children}
    </h2>
  ),
  h3: ({ children, ...props }) => (
    <h3 className="text-lg font-bold text-white mt-4 mb-2" {...props}>
      {children}
    </h3>
  ),
  h4: ({ children, ...props }) => (
    <h4 className="text-base font-semibold text-white mt-4 mb-2" {...props}>
      {children}
    </h4>
  ),
  p: ({ children, ...props }) => (
    <p className="text-slate-400 mb-3 leading-relaxed" {...props}>
      {children}
    </p>
  ),
  a: ({ children, href, ...props }) => (
    <a
      href={href}
      className="text-blue-400 hover:text-blue-300 underline transition-colors"
      target="_blank"
      rel="noopener noreferrer"
      {...props}
    >
      {children}
    </a>
  ),
  ul: ({ children, ...props }) => (
    <ul className="list-disc list-inside text-slate-400 mb-3 space-y-1 ml-2" {...props}>
      {children}
    </ul>
  ),
  ol: ({ children, ...props }) => (
    <ol className="list-decimal list-inside text-slate-400 mb-3 space-y-1 ml-2" {...props}>
      {children}
    </ol>
  ),
  li: ({ children, ...props }) => (
    <li className="text-slate-400" {...props}>
      {children}
    </li>
  ),
  strong: ({ children, ...props }) => (
    <strong className="text-slate-200 font-semibold" {...props}>
      {children}
    </strong>
  ),
  em: ({ children, ...props }) => (
    <em className="text-slate-300 italic" {...props}>
      {children}
    </em>
  ),
  del: ({ children, ...props }) => (
    <del className="text-slate-500 line-through" {...props}>
      {children}
    </del>
  ),
  code: ({ children, className, ...props }) => {
    const isInline = !className;
    if (isInline) {
      return (
        <code className="bg-slate-800 text-amber-300 px-1.5 py-0.5 rounded text-sm font-mono" {...props}>
          {children}
        </code>
      );
    }
    return (
      <code className={`${className} block text-sm font-mono`} {...props}>
        {children}
      </code>
    );
  },
  pre: ({ children, ...props }) => (
    <pre className="bg-slate-900 border border-slate-700 rounded-lg p-4 mb-4 overflow-x-auto text-sm font-mono text-slate-300" {...props}>
      {children}
    </pre>
  ),
  blockquote: ({ children, ...props }) => (
    <blockquote className="border-l-4 border-brand-600/50 pl-4 py-1 mb-3 text-slate-400 italic bg-slate-800/50 rounded-r" {...props}>
      {children}
    </blockquote>
  ),
  hr: () => <hr className="border-slate-700 my-6" />,
  table: ({ children, ...props }) => (
    <div className="overflow-x-auto mb-4">
      <table className="min-w-full border-collapse border border-slate-700 text-sm" {...props}>
        {children}
      </table>
    </div>
  ),
  thead: ({ children, ...props }) => (
    <thead className="bg-slate-800" {...props}>
      {children}
    </thead>
  ),
  tbody: ({ children, ...props }) => (
    <tbody className="divide-y divide-slate-800" {...props}>
      {children}
    </tbody>
  ),
  tr: ({ children, ...props }) => (
    <tr className="border-b border-slate-800" {...props}>
      {children}
    </tr>
  ),
  th: ({ children, ...props }) => (
    <th className="px-4 py-2 text-left font-semibold text-slate-200 border border-slate-700" {...props}>
      {children}
    </th>
  ),
  td: ({ children, ...props }) => (
    <td className="px-4 py-2 text-slate-400 border border-slate-700" {...props}>
      {children}
    </td>
  ),
  img: ({ src, alt, ...props }) => (
    <img src={src} alt={alt} className="max-w-full rounded-lg my-4" loading="lazy" {...props} />
  ),
  input: ({ type, checked, disabled, ...props }) => {
    if (type === "checkbox") {
      return (
        <input
          type="checkbox"
          checked={checked}
          disabled={disabled}
          className="mr-2 accent-brand-600"
          readOnly
          {...props}
        />
      );
    }
    return <input type={type} checked={checked} disabled={disabled} {...props} />;
  },
};

export function MarkdownRenderer({ content }: MarkdownRendererProps) {
  if (!content) {
    return <p className="text-slate-500 italic">暂无内容</p>;
  }

  return (
    <ReactMarkdown
      remarkPlugins={[remarkGfm]}
      components={components}
    >
      {content}
    </ReactMarkdown>
  );
}
