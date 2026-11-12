import { getAccessToken, tryRefresh } from '$lib/auth/token';
import type {
  AdminOverview,
  AuditLog,
  Bootstrap,
  Bidang,
  CapaianTahunan,
  DashboardCategory,
  DashboardCategoryResult,
  DashboardItem,
  DashboardSummary,
  Factor,
  FactorRow,
  IKUChart,
  IKUValue,
  IKUValueDetail,
  Job,
  Label,
  Periode,
  Position,
  SPMDetail,
  SPMTrace,
  Submission,
  SubmissionDetail,
  SubmissionRow,
  Target,
  TindakLanjut,
  User
} from '$lib/types';

const BASE = '/api/v1';
const NETWORK_MESSAGE = 'Tidak dapat menghubungi server. Periksa koneksi lalu coba lagi.';
const SERVER_MESSAGE = 'Server sedang bermasalah. Coba lagi beberapa saat.';

export class ApiError extends Error {
  code: string;
  details?: Record<string, unknown>;
  status: number;
  constructor(status: number, code: string, message: string, details?: Record<string, unknown>) {
    super(message);
    this.status = status;
    this.code = code;
    this.details = details;
  }
}

async function readPayload(response: Response): Promise<any> {
  const raw = await response.text();
  if (!raw) return null;
  try {
    return JSON.parse(raw);
  } catch {
    return null;
  }
}

function fallbackMessage(status: number): string {
  if (status >= 500) return SERVER_MESSAGE;
  if (status === 401) return 'Sesi berakhir. Silakan login kembali.';
  if (status === 403) return 'Anda tidak memiliki akses untuk tindakan ini.';
  if (status === 404) return 'Data yang diminta tidak ditemukan.';
  return 'Permintaan gagal. Coba lagi.';
}

async function request<T>(path: string, options: RequestInit = {}, retry = true): Promise<T> {
  const headers = new Headers(options.headers ?? {});
  const token = getAccessToken();
  if (token) headers.set('Authorization', `Bearer ${token}`);
  if (options.body && !(options.body instanceof FormData) && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }
  let response: Response;
  try {
    response = await fetch(`${BASE}${path}`, { ...options, headers, credentials: 'include' });
  } catch {
    throw new ApiError(0, 'NETWORK_ERROR', NETWORK_MESSAGE);
  }
  if (response.status === 401 && retry) {
    const fresh = await tryRefresh();
    if (fresh) return request<T>(path, options, false);
  }
  const payload = await readPayload(response);
  if (!response.ok) {
    const message = typeof payload?.message === 'string' && payload.message.trim()
      ? payload.message
      : fallbackMessage(response.status);
    throw new ApiError(
      response.status,
      payload?.code ?? 'HTTP_ERROR',
      response.status >= 500 ? SERVER_MESSAGE : message,
      payload?.details
    );
  }
  return (payload?.data ?? payload) as T;
}

function query(params: Record<string, string | undefined>): string {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value) search.set(key, value);
  }
  const qs = search.toString();
  return qs ? `?${qs}` : '';
}

export const api = {
  bootstrap: () => request<Bootstrap>('/meta/bootstrap'),
  me: () => request<{ user: User; acting_as?: User }>('/auth/me'),
  profile: () => request<{ user: User; acting_as?: User }>('/profile'),
  updateProfile: (nama: string) =>
    request<User>('/profile', { method: 'PATCH', body: JSON.stringify({ nama }) }),
  changePassword: (current_password: string, new_password: string) =>
    request<{ ok: boolean }>('/profile/change-password', {
      method: 'POST',
      body: JSON.stringify({ current_password, new_password })
    }),

  dashboard: (params: { periode?: string; faskes_id?: string; status?: string; indikator_id?: string; bidang_id?: string }) =>
    request<{ items: DashboardItem[]; summary: DashboardSummary; meta: Record<string, unknown> }>(
      `/dashboard/spm${query(params)}`
    ),
  trace: (kode: string, periode: string) => request<SPMTrace>(`/dashboard/spm/${kode}/trace${query({ periode })}`),
  spmDetail: (kode: string, params: { periode?: string; faskes_id?: string } = {}) =>
    request<SPMDetail>(`/dashboard/spm/${kode}/detail${query(params)}`),
  ikuByFaskes: (params: { periode?: string; faskes_id?: string } = {}) =>
    request<IKUChart>(`/dashboard/iku-by-faskes${query(params)}`),
  dashboardCategory: (
    kode: string,
    params: { periode?: string; faskes_id?: string; status?: string; bidang_id?: string; timker_id?: string } = {}
  ) => request<DashboardCategoryResult>(`/dashboard/category${query({ kode, ...params })}`),

  submissions: (params: { periode?: string; status?: string; faskes_id?: string } = {}) =>
    request<Submission[]>(`/submissions${query(params)}`),
  submission: (id: string) => request<SubmissionDetail>(`/submissions/${id}`),
  upload: (form: FormData, idempotencyKey?: string) =>
    request<{ submission: Submission; job?: Job }>('/submissions', {
      method: 'POST',
      body: form,
      headers: idempotencyKey ? { 'Idempotency-Key': idempotencyKey } : {}
    }),
  job: (id: string) => request<Job>(`/jobs/${id}`),
  updateRow: (submissionId: string, rowId: string, ikuKode: string, reason: string) =>
    request<SubmissionRow>(`/submissions/${submissionId}/rows/${rowId}`, {
      method: 'PATCH',
      body: JSON.stringify({ iku_kode: ikuKode, reason })
    }),
  submit: (id: string) => request<Submission>(`/submissions/${id}/submit`, { method: 'POST' }),
  decision: (id: string, action: 'approve' | 'return', notes: string, reasonCode = '') =>
    request<Submission>(`/submissions/${id}/decision`, {
      method: 'POST',
      body: JSON.stringify({ action, notes, reason_code: reasonCode })
    }),

  labels: () => request<Label[]>('/master/labels'),
  categories: () => request<DashboardCategory[]>('/master/categories'),
  bidang: () => request<Bidang[]>('/master/bidang'),
  indicators: () => request<Bootstrap['indicators']>('/master/indicators'),
  targets: (periode?: string) => request<Target[]>(`/master/targets${query({ periode })}`),
  periods: () => request<Periode[]>('/master/periods'),
  createLabel: (payload: Record<string, unknown>) =>
    request<Label>('/master/labels', { method: 'POST', body: JSON.stringify(payload) }),
  updateLabel: (id: string, payload: Record<string, unknown>) =>
    request<Label>(`/master/labels/${id}`, { method: 'PATCH', body: JSON.stringify(payload) }),
  createCategory: (payload: Record<string, unknown>) =>
    request<DashboardCategory>('/master/categories', { method: 'POST', body: JSON.stringify(payload) }),
  updateCategory: (id: string, payload: Record<string, unknown>) =>
    request<DashboardCategory>(`/master/categories/${id}`, { method: 'PATCH', body: JSON.stringify(payload) }),
  deleteCategory: (id: string) =>
    request<DashboardCategory>(`/master/categories/${id}`, { method: 'DELETE' }),
  createBidang: (payload: Record<string, unknown>) =>
    request<Bidang>('/master/bidang', { method: 'POST', body: JSON.stringify(payload) }),
  updateBidang: (id: string, payload: Record<string, unknown>) =>
    request<Bidang>(`/master/bidang/${id}`, { method: 'PATCH', body: JSON.stringify(payload) }),
  createIndicator: (payload: Record<string, unknown>) =>
    request<Bootstrap['indicators'][number]>('/master/indicators', {
      method: 'POST',
      body: JSON.stringify(payload)
    }),
  upsertTarget: (payload: Record<string, unknown>) =>
    request<Target>('/master/targets', { method: 'PUT', body: JSON.stringify(payload) }),
  capaianTahunan: (tahun: number) =>
    request<CapaianTahunan[]>(`/master/capaian-tahunan${query({ tahun: String(tahun) })}`),
  createCapaianTahunan: (payload: Record<string, unknown>) =>
    request<CapaianTahunan>('/master/capaian-tahunan', {
      method: 'POST',
      body: JSON.stringify(payload)
    }),
  updateCapaianTahunan: (id: string, payload: Record<string, unknown>) =>
    request<CapaianTahunan>(`/master/capaian-tahunan/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(payload)
    }),
  deleteCapaianTahunan: (id: string) =>
    request<{ ok: boolean }>(`/master/capaian-tahunan/${id}`, { method: 'DELETE' }),

  ikuValues: (periode?: string, scopeType = 'KABUPATEN', scopeId = 'wil-kab-swb') =>
    request<IKUValue[]>(`/iku/values${query({ periode, scope_type: scopeType, scope_id: scopeId })}`),
  ikuByKode: (kode: string, periode?: string) =>
    request<IKUValue>(`/iku/${encodeURIComponent(kode)}/values${query({ periode })}`),
  ikuDetail: (id: string) => request<IKUValueDetail[]>(`/iku-values/${id}/detail`),

  factors: (periode?: string, scopeType = 'KABUPATEN', scopeId = 'wil-kab-swb') =>
    request<FactorRow[]>(`/factors${query({ periode, scope_type: scopeType, scope_id: scopeId })}`),
  upsertFactor: (ikuValueId: string, payload: Record<string, unknown>) =>
    request<Factor>(`/factors/${ikuValueId}`, { method: 'PUT', body: JSON.stringify(payload) }),
  tl: (periode?: string) => request<TindakLanjut[]>(`/tindak-lanjut${query({ periode })}`),
  createTL: (payload: Record<string, unknown>) =>
    request<TindakLanjut>('/tindak-lanjut', { method: 'POST', body: JSON.stringify(payload) }),
  updateTL: (id: string, payload: Record<string, unknown>) =>
    request<TindakLanjut>(`/tindak-lanjut/${id}`, { method: 'PATCH', body: JSON.stringify(payload) }),

  audit: (params: { entity_type?: string; entity_id?: string; limit?: string } = {}) =>
    request<AuditLog[]>(`/audit${query(params)}`),
  fileUrl: (id: string) => `${BASE}/files/${id}`,

  adminOverview: () => request<AdminOverview>('/admin/overview'),
  createUser: (payload: Record<string, unknown>) =>
    request<User>('/admin/users', { method: 'POST', body: JSON.stringify(payload) }),
  updateUser: (id: string, payload: Record<string, unknown>) =>
    request<User>(`/admin/users/${id}`, { method: 'PATCH', body: JSON.stringify(payload) }),
  resetPassword: (id: string, password: string) =>
    request<{ ok: boolean }>(`/admin/users/${id}/reset-password`, {
      method: 'POST',
      body: JSON.stringify({ password })
    }),
  createFaskes: (payload: Record<string, unknown>) =>
    request<Record<string, unknown>>('/admin/faskes', { method: 'POST', body: JSON.stringify(payload) }),
  updateFaskes: (id: string, payload: Record<string, unknown>) =>
    request<Record<string, unknown>>(`/admin/faskes/${id}`, { method: 'PATCH', body: JSON.stringify(payload) }),
  createPosition: (payload: Record<string, unknown>) =>
    request<Position>('/admin/positions', { method: 'POST', body: JSON.stringify(payload) }),
  updatePosition: (id: string, payload: Record<string, unknown>) =>
    request<Position>(`/admin/positions/${id}`, { method: 'PATCH', body: JSON.stringify(payload) }),

  downloadFile: async (id: string): Promise<{ blob: Blob; filename: string }> => {
    const doFetch = async (token: string | null) => {
      try {
        return await fetch(`${BASE}/files/${id}`, {
          headers: token ? { Authorization: `Bearer ${token}` } : {},
          credentials: 'include'
        });
      } catch {
        throw new ApiError(0, 'NETWORK_ERROR', NETWORK_MESSAGE);
      }
    };
    let response = await doFetch(getAccessToken());
    if (response.status === 401) {
      const fresh = await tryRefresh();
      if (fresh) response = await doFetch(fresh);
    }
    if (!response.ok) throw new ApiError(response.status, 'DOWNLOAD_FAILED', 'Unduhan gagal.');
    const disposition = response.headers.get('content-disposition') ?? '';
    const match = disposition.match(/filename="([^"]+)"/);
    return { blob: await response.blob(), filename: match?.[1] ?? 'file' };
  },

  downloadTemplate: async (): Promise<{ blob: Blob; filename: string }> => {
    const token = getAccessToken();
    let response: Response;
    try {
      response = await fetch(`${BASE}/master/template`, {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
        credentials: 'include'
      });
    } catch {
      throw new ApiError(0, 'NETWORK_ERROR', NETWORK_MESSAGE);
    }
    if (!response.ok) throw new ApiError(response.status, 'DOWNLOAD_FAILED', 'Template gagal diunduh.');
    return { blob: await response.blob(), filename: 'template-bnba.csv' };
  },

  exportSPM: async (format: 'xlsx' | 'pdf', periode: string): Promise<{ blob: Blob; filename: string }> => {
    const doFetch = async (token: string | null) => {
      try {
        return await fetch(`${BASE}/dashboard/spm/export${query({ format, periode })}`, {
          headers: token ? { Authorization: `Bearer ${token}` } : {},
          credentials: 'include'
        });
      } catch {
        throw new ApiError(0, 'NETWORK_ERROR', NETWORK_MESSAGE);
      }
    };
    let response = await doFetch(getAccessToken());
    if (response.status === 401) {
      const fresh = await tryRefresh();
      if (fresh) response = await doFetch(fresh);
    }
    if (!response.ok) throw new ApiError(response.status, 'EXPORT_FAILED', 'Export gagal.');
    const disposition = response.headers.get('content-disposition') ?? '';
    const match = disposition.match(/filename="([^"]+)"/);
    return { blob: await response.blob(), filename: match?.[1] ?? `spm.${format}` };
  }
};
