import { useState, type FormEvent } from "react";
import { Link, Navigate, useNavigate } from "react-router-dom";
import { useAuth } from "../stores/AuthContext";

export function Login() {
  const { login, isAuthenticated, loading: authLoading } = useAuth();
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  if (isAuthenticated && !authLoading) {
    return <Navigate to="/" replace />;
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    if (!username || !password) {
      setError("请输入用户名和密码");
      return;
    }
    setLoading(true);
    try {
      await login(username, password);
      navigate("/");
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "登录失败，请重试");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-[80vh] px-4 w-full">
      <div className="w-full max-w-md">
        <div className="card p-8">
          <div className="text-center mb-8">
            <div className="w-14 h-14 rounded-xl bg-brand-100 flex items-center justify-center mx-auto mb-4">
              <span className="material-symbols-outlined text-brand-600 text-[32px] fill">hexagon</span>
            </div>
            <h1 className="text-2xl font-bold text-slate-900">欢迎回来</h1>
            <p className="text-slate-500 text-sm mt-1">登录您的 SkillHub Pro 账户</p>
          </div>

          <form onSubmit={handleSubmit} className="flex flex-col gap-5">
            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 text-sm px-4 py-3 rounded-lg flex items-center gap-2">
                <span className="material-symbols-outlined text-[18px]">error</span>
                {error}
              </div>
            )}

            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium text-slate-700">用户名</label>
              <input
                type="text"
                value={username}
                onChange={e => setUsername(e.target.value)}
                className="input-field"
                placeholder="请输入用户名"
                autoComplete="username"
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium text-slate-700">密码</label>
              <input
                type="password"
                value={password}
                onChange={e => setPassword(e.target.value)}
                className="input-field"
                placeholder="请输入密码"
                autoComplete="current-password"
              />
            </div>

            <button type="submit" disabled={loading} className="btn-primary w-full mt-2">
              {loading ? (
                <span className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
              ) : "登录"}
            </button>

            <div className="relative my-2">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-slate-200" />
              </div>
              <div className="relative flex justify-center text-xs">
                <span className="bg-white px-3 text-slate-400">或</span>
              </div>
            </div>

            <button
              type="button"
              className="btn-secondary w-full"
              onClick={() => {/* GitHub OAuth */}}
            >
              <span className="material-symbols-outlined text-[18px]">lock</span>
              使用 GitHub 登录
            </button>
          </form>

          <p className="text-center text-sm text-slate-500 mt-6">
            还没有账户？
            <Link to="/register" className="text-brand-600 font-medium hover:text-brand-700 ml-1">
              立即注册
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
