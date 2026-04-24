import { describe, it, expect, vi, beforeEach } from 'vitest';

beforeEach(() => {
  localStorage.clear();
});

describe('API Client', () => {
  it('stores and retrieves token', async () => {
    const { setToken, getAuthHeaders } = await import('../lib/api/client');
    setToken('test-token-123');
    const headers = getAuthHeaders();
    expect(headers.Authorization).toBe('Bearer test-token-123');
  });

  it('returns empty headers when no token', async () => {
    const { getAuthHeaders } = await import('../lib/api/client');
    const headers = getAuthHeaders();
    expect(headers.Authorization).toBeUndefined();
  });

  it('clears token on logout', async () => {
    const { setToken, clearToken, getAuthHeaders } = await import('../lib/api/client');
    setToken('test-token');
    clearToken();
    const headers = getAuthHeaders();
    expect(headers.Authorization).toBeUndefined();
  });

  it('makes GET request with correct URL', async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: [] }),
    });
    vi.stubGlobal('fetch', mockFetch);

    const { api } = await import('../lib/api/client');
    await api.get('/skills', { page: 1, pageSize: 10 });

    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('/api/v1/skills?page=1&pageSize=10'),
      expect.any(Object),
    );

    vi.unstubAllGlobals();
  });

  it('makes POST request with JSON body', async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ success: true }),
    });
    vi.stubGlobal('fetch', mockFetch);

    const { api } = await import('../lib/api/client');
    await api.post('/auth/login', { username: 'test', password: 'pass' });

    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('/api/v1/auth/login'),
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ username: 'test', password: 'pass' }),
      }),
    );

    vi.unstubAllGlobals();
  });

  it('throws ApiError on failed request', async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 401,
      json: () => Promise.resolve({ message: 'Unauthorized' }),
    });
    vi.stubGlobal('fetch', mockFetch);

    const { api, ApiError } = await import('../lib/api/client');
    await expect(api.get('/skills')).rejects.toThrow(ApiError);
    await expect(api.get('/skills')).rejects.toThrow('Unauthorized');

    vi.unstubAllGlobals();
  });
});
