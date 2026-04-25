import { useState, type FormEvent } from "react";
import { Link, Navigate, useNavigate } from "react-router-dom";
import { useAuth } from "../stores/AuthContext";

export function Register() {
  const { register, isAuthenticated, loading: authLoading } = useAuth();
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  if (isAuthenticated && !authLoading) {
    return <Navigate to="/" replace />;
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError("");

    if (!username || !email || !password) {
      setError("请填写所有必填字段");
      return;
    }
    if (password !== confirmPassword) {
      setError("两次密码输入不一致");
      return;
    }
    if (password.length < 6) {
      setError("密码长度至少为 6 个字符");
      return;
    }

    setLoading(true);
    try {
      await register(username, email, password);
      navigate("/");
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "注册失败，请重试");
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
            <h1 className="text-2xl font-bold text-slate-900">创建账户</h1>
            <p className="text-slate-500 text-sm mt-1">注册 SkillHub Pro 开始探索</p>
          </div>

          <form onSubmit={handleSubmit} className="flex flex-col gap-5">
            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 text-sm px-4 py-3 rounded-lg flex items-center gap-2" role="alert">
                <span className="material-symbols-outlined text-[18px]">error</span>
                {error}
              </div>
            )}

            <div className="flex flex-col gap-1.5">
              <label htmlFor="reg-username" className="text-sm font-medium text-slate-700">用户名 *</label>
              <input
                id="reg-username"
                type="text"
                value={username}
                onChange={e => setUsername(e.target.value)}
                className="input-field"
                placeholder="请输入用户名"
                autoComplete="username"
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label htmlFor="reg-email" className="text-sm font-medium text-slate-700">邮箱 *</label>
              <input
                id="reg-email"
                type="email"
                value={email}
                onChange={e => setEmail(e.target.value)}
                className="input-field"
                placeholder="请输入邮箱地址"
                autoComplete="email"
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label htmlFor="reg-password" className="text-sm font-medium text-slate-700">密码 *</label>
              <input
                id="reg-password"
                type="password"
                value={password}
                onChange={e => setPassword(e.target.value)}
                className="input-field"
                placeholder="至少 6 个字符"
                autoComplete="new-password"
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label htmlFor="reg-confirm-password" className="text-sm font-medium text-slate-700">确认密码 *</label>
              <input
                id="reg-confirm-password"
                type="password"
                value={confirmPassword}
                onChange={e => setConfirmPassword(e.target.value)}
                className="input-field"
                placeholder="再次输入密码"
                autoComplete="new-password"
              />
            </div>

            <button type="submit" disabled={loading} className="btn-primary w-full mt-2">
              {loading ? (
                <span className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
              ) : "注册"}
            </button>
          </form>

          <p className="text-center text-sm text-slate-500 mt-6">
            已有账户？
            <Link to="/login" className="text-brand-600 font-medium hover:text-brand-700 ml-1">
              立即登录
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
