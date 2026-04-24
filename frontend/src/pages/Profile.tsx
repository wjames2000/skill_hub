import { useState, useEffect } from "react";
import { useAuth } from "../stores/AuthContext";
import { authApi } from "../lib/api/auth";
import type { ApiKey } from "../types";

type Tab = 'profile' | 'api-keys' | 'favorites' | 'history';

export function Profile() {
  const { user, updateUser, logout } = useAuth();
  const [activeTab, setActiveTab] = useState<Tab>('profile');
  const [apiKeys, setApiKeys] = useState<ApiKey[]>([]);
  const [newKeyName, setNewKeyName] = useState("");
  const [newKeyValue, setNewKeyValue] = useState("");
  const [editing, setEditing] = useState(false);
  const [bio, setBio] = useState(user?.bio || "");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (activeTab === 'api-keys') {
      authApi.getApiKeys().then(setApiKeys).catch(() => {});
    }
  }, [activeTab]);

  const handleSaveProfile = async () => {
    setSaving(true);
    try {
      const updated = await authApi.updateProfile({ bio } as never);
      updateUser(updated);
      setEditing(false);
    } catch { /* ignore */ }
    setSaving(false);
  };

  const handleCreateKey = async () => {
    if (!newKeyName.trim()) return;
    try {
      const key = await authApi.createApiKey(newKeyName);
      setNewKeyValue(key.key);
      setNewKeyName("");
      authApi.getApiKeys().then(setApiKeys);
    } catch { /* ignore */ }
  };

  const handleRevokeKey = async (id: number) => {
    try {
      await authApi.revokeApiKey(id);
      setApiKeys(prev => prev.filter(k => k.id !== id));
    } catch { /* ignore */ }
  };

  const tabs: { key: Tab; label: string; icon: string }[] = [
    { key: 'profile', label: '个人资料', icon: 'person' },
    { key: 'api-keys', label: 'API Key', icon: 'key' },
    { key: 'favorites', label: '我的收藏', icon: 'star' },
    { key: 'history', label: '浏览历史', icon: 'history' },
  ];

  return (
    <div className="w-full max-w-[1280px] mx-auto px-6 py-8">
      <div className="flex items-center gap-4 mb-8">
        <div className="w-16 h-16 rounded-full bg-slate-200 overflow-hidden border border-slate-300">
          <img src={user?.avatar || ''} alt="" className="w-full h-full object-cover" />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-slate-900">{user?.username}</h1>
          <p className="text-sm text-slate-500">{user?.email}</p>
        </div>
      </div>

      <div className="flex gap-8 items-start">
        <nav className="w-52 flex-shrink-0 flex flex-col gap-1 sticky top-24">
          {tabs.map(tab => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key)}
              className={`flex items-center gap-3 px-4 py-2.5 rounded-lg text-sm font-medium text-left transition-all ${
                activeTab === tab.key
                  ? 'bg-brand-50 text-brand-700'
                  : 'text-slate-600 hover:bg-slate-100'
              }`}
            >
              <span className="material-symbols-outlined text-[20px]">{tab.icon}</span>
              {tab.label}
            </button>
          ))}
          <button
            onClick={logout}
            className="flex items-center gap-3 px-4 py-2.5 rounded-lg text-sm font-medium text-left text-red-600 hover:bg-red-50 mt-4 transition-all"
          >
            <span className="material-symbols-outlined text-[20px]">logout</span>
            退出登录
          </button>
        </nav>

        <div className="flex-1 min-w-0">
          {activeTab === 'profile' && (
            <div className="card p-6">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-lg font-bold text-slate-900">个人资料</h2>
                {!editing ? (
                  <button onClick={() => setEditing(true)} className="btn-ghost text-sm">
                    <span className="material-symbols-outlined text-[18px]">edit</span>
                    编辑
                  </button>
                ) : (
                  <div className="flex gap-2">
                    <button onClick={() => setEditing(false)} className="btn-secondary text-sm">取消</button>
                    <button onClick={handleSaveProfile} disabled={saving} className="btn-primary text-sm">
                      {saving ? '保存中...' : '保存'}
                    </button>
                  </div>
                )}
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="text-sm text-slate-500 block mb-1">用户名</label>
                  <p className="font-medium text-slate-900">{user?.username}</p>
                </div>
                <div>
                  <label className="text-sm text-slate-500 block mb-1">邮箱</label>
                  <p className="font-medium text-slate-900">{user?.email}</p>
                </div>
                <div className="md:col-span-2">
                  <label className="text-sm text-slate-500 block mb-1">个人简介</label>
                  {editing ? (
                    <textarea
                      value={bio}
                      onChange={e => setBio(e.target.value)}
                      className="input-field h-24 resize-none"
                      placeholder="介绍一下自己..."
                    />
                  ) : (
                    <p className="text-slate-700">{user?.bio || '暂无简介'}</p>
                  )}
                </div>
                <div>
                  <label className="text-sm text-slate-500 block mb-1">角色</label>
                  <p className="font-medium text-slate-900 flex items-center gap-1">
                    <span className="material-symbols-outlined text-[16px] text-brand-600">
                      {user?.role === 'admin' ? 'admin_panel_settings' : 'person'}
                    </span>
                    {user?.role === 'admin' ? '管理员' : '用户'}
                  </p>
                </div>
                <div>
                  <label className="text-sm text-slate-500 block mb-1">注册时间</label>
                  <p className="font-medium text-slate-900">{user?.createdAt || '-'}</p>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'api-keys' && (
            <div className="card p-6">
              <h2 className="text-lg font-bold text-slate-900 mb-6">API Key 管理</h2>

              <div className="flex items-center gap-3 mb-6">
                <input
                  type="text"
                  value={newKeyName}
                  onChange={e => setNewKeyName(e.target.value)}
                  className="input-field flex-1"
                  placeholder="输入 Key 名称..."
                />
                <button onClick={handleCreateKey} className="btn-primary" disabled={!newKeyName.trim()}>
                  <span className="material-symbols-outlined text-[18px]">add</span>
                  生成
                </button>
              </div>

              {newKeyValue && (
                <div className="bg-green-50 border border-green-200 rounded-lg p-4 mb-6">
                  <p className="text-sm font-medium text-green-800 mb-2">新 API Key 已生成（请立即保存）</p>
                  <div className="flex items-center gap-2">
                    <code className="flex-1 bg-white border border-green-200 rounded px-3 py-2 text-sm font-mono text-green-900 break-all">
                      {newKeyValue}
                    </code>
                    <button
                      onClick={() => navigator.clipboard.writeText(newKeyValue)}
                      className="btn-ghost text-sm"
                    >
                      <span className="material-symbols-outlined text-[18px]">content_copy</span>
                    </button>
                  </div>
                </div>
              )}

              <div className="flex flex-col gap-3">
                {apiKeys.length === 0 ? (
                  <p className="text-sm text-slate-500 text-center py-8">暂无 API Key</p>
                ) : (
                  apiKeys.map(key => (
                    <div key={key.id} className="flex items-center justify-between p-4 bg-slate-50 rounded-lg border border-slate-200">
                      <div>
                        <p className="font-medium text-slate-900">{key.name}</p>
                        <p className="text-xs text-slate-500 mt-0.5">创建于 {key.createdAt}</p>
                      </div>
                      <div className="flex items-center gap-3">
                        <span className="text-xs text-slate-400">{key.lastUsed ? `最近使用: ${key.lastUsed}` : '未使用'}</span>
                        <button
                          onClick={() => handleRevokeKey(key.id)}
                          className="text-red-600 hover:text-red-700 text-sm font-medium"
                        >
                          吊销
                        </button>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>
          )}

          {activeTab === 'favorites' && (
            <div className="card p-6">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-lg font-bold text-slate-900">我的收藏</h2>
              </div>
              <p className="text-sm text-slate-500 text-center py-16">暂无收藏的技能</p>
            </div>
          )}

          {activeTab === 'history' && (
            <div className="card p-6">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-lg font-bold text-slate-900">浏览历史</h2>
              </div>
              <p className="text-sm text-slate-500 text-center py-16">暂无浏览记录</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
