import { writable } from 'svelte/store';
import type { AuthResponse, User } from '$lib/types';
import { setAccessToken, setRefreshHandler, getAccessToken } from '$lib/auth/token';
import { currentActor } from '$lib/stores/session';

export const authReady = writable(false);
export const actingAs = writable<User | null>(null);
export const authError = writable('');

const NETWORK_MESSAGE = 'Tidak dapat menghubungi server. Periksa koneksi lalu coba lagi.';
const SERVER_MESSAGE = 'Server sedang bermasalah. Coba lagi beberapa saat.';
const CREDENTIALS_MESSAGE = 'Username atau password salah.';
const SESSION_MESSAGE = 'Sesi berakhir. Silakan login kembali.';

function isSafeMessage(message: unknown): message is string {
  if (typeof message !== 'string') return false;
  const text = message.trim();
  return (
    text.length > 0 &&
    text.length <= 140 &&
    !/internal|panic|sql|constraint|nil pointer|goroutine|runtime error/i.test(text)
  );
}

async function readAuth(response: Response, action: 'login' | 'session' = 'login'): Promise<AuthResponse> {
  const raw = await response.text();
  let payload: { data?: AuthResponse; message?: string } | null = null;
  if (raw) {
    try {
      payload = JSON.parse(raw);
    } catch {
      payload = null;
    }
  }

  if (!response.ok) {
    if (response.status >= 500) throw new Error(SERVER_MESSAGE);
    if (action === 'login' && response.status === 401) throw new Error(CREDENTIALS_MESSAGE);
    if (response.status === 401) throw new Error(SESSION_MESSAGE);
    if (response.status === 403) throw new Error('Akun ini tidak dapat digunakan. Hubungi administrator.');
    if (response.status === 429) throw new Error('Terlalu banyak percobaan. Tunggu sebentar lalu coba lagi.');
    if (isSafeMessage(payload?.message)) throw new Error(payload.message.trim());
    throw new Error('Permintaan gagal. Coba lagi.');
  }

  if (!payload?.data?.access_token || !payload?.data?.user) {
    throw new Error(SERVER_MESSAGE);
  }
  return payload.data;
}

async function postAuth(path: string, body?: unknown): Promise<Response> {
  try {
    return await fetch(path, {
      method: 'POST',
      credentials: 'include',
      headers: body ? { 'Content-Type': 'application/json' } : undefined,
      body: body ? JSON.stringify(body) : undefined
    });
  } catch {
    throw new Error(NETWORK_MESSAGE);
  }
}

function applyAuth(result: AuthResponse) {
  setAccessToken(result.access_token);
  currentActor.set(result.user);
  actingAs.set(result.acting_as ?? null);
  authReady.set(true);
  authError.set('');
}

export async function login(username: string, password: string): Promise<void> {
  const response = await postAuth('/api/v1/auth/login', { username, password });
  applyAuth(await readAuth(response, 'login'));
}

export async function restoreSession(): Promise<boolean> {
  try {
    const response = await fetch('/api/v1/auth/refresh', {
      method: 'POST',
      credentials: 'include'
    });
    if (!response.ok) {
      clearAuth();
      return false;
    }
    applyAuth(await readAuth(response, 'session'));
    return true;
  } catch {
    clearAuth();
    return false;
  }
}

export async function logout(): Promise<void> {
  const token = getAccessToken();
  try {
    await fetch('/api/v1/auth/logout', {
      method: 'POST',
      credentials: 'include',
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    });
  } finally {
    clearAuth();
  }
}

export async function startActingAs(userID: string): Promise<void> {
  const token = getAccessToken();
  const response = await postAuthWithToken('/api/v1/auth/acting-as', token, { user_id: userID });
  applyAuth(await readAuth(response, 'session'));
}

async function postAuthWithToken(path: string, token: string | null, body: unknown): Promise<Response> {
  try {
    return await fetch(path, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {})
      },
      body: JSON.stringify(body)
    });
  } catch {
    throw new Error(NETWORK_MESSAGE);
  }
}

export async function stopActingAs(): Promise<void> {
  await restoreSession();
}

export function clearAuth() {
  setAccessToken(null);
  currentActor.set(null);
  actingAs.set(null);
  authReady.set(true);
}

setRefreshHandler(async () => {
  try {
    const response = await fetch('/api/v1/auth/refresh', {
      method: 'POST',
      credentials: 'include'
    });
    if (!response.ok) {
      clearAuth();
      return null;
    }
    const result = await readAuth(response, 'session');
    applyAuth(result);
    return result.access_token;
  } catch {
    clearAuth();
    return null;
  }
});
