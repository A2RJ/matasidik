import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]): string {
  return twMerge(clsx(inputs));
}

export function formatPercent(value?: number | null, digits = 2): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '—';
  return `${value.toLocaleString('id-ID', {
    minimumFractionDigits: digits,
    maximumFractionDigits: digits
  })}%`;
}

export function formatNumber(value?: number | null, digits = 0): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '—';
  return value.toLocaleString('id-ID', {
    minimumFractionDigits: digits,
    maximumFractionDigits: digits
  });
}

export function formatMetric(
  value?: number | null,
  metric?: { tipe_nilai?: string | null; satuan?: string | null } | null,
  digits = 2
): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '—';
  if (!metric || !metric.tipe_nilai) {
    return formatPercent(value, digits);
  }
  if (metric.tipe_nilai === 'RATIO') {
    if (metric.satuan?.toLowerCase().startsWith('rasio')) {
      return `${formatNumber(value, digits)}${metric.satuan ? ` ${metric.satuan}` : ''}`;
    }
    return formatPercent(value, digits);
  }
  const number = formatNumber(value, metric.tipe_nilai === 'COUNT' ? 0 : digits);
  return metric.satuan ? `${number} ${metric.satuan}` : number;
}

export function formatGap(
  value?: number | null,
  metric?: { tipe_nilai?: string | null; satuan?: string | null } | null
): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '—';
  if (!metric || !metric.tipe_nilai) {
    return `${value > 0 ? '+' : ''}${formatNumber(value, 2)} pp`;
  }
  if (metric.tipe_nilai === 'RATIO') {
    if (metric.satuan?.toLowerCase().startsWith('rasio')) {
      return `${value > 0 ? '+' : ''}${formatNumber(value, 2)}${metric.satuan ? ` ${metric.satuan}` : ''}`;
    }
    return `${value > 0 ? '+' : ''}${formatNumber(value, 2)} pp`;
  }
  const number = formatNumber(value, metric.tipe_nilai === 'COUNT' ? 0 : 2);
  return `${value > 0 ? '+' : ''}${number}${metric.satuan ? ` ${metric.satuan}` : ''}`;
}

export function formatDateTime(value?: string | null): string {
  if (!value) return '—';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'short'
  });
}

export function formatDate(value?: string | null): string {
  if (!value) return '—';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleDateString('id-ID', { dateStyle: 'medium' });
}

export function truncateHash(hash?: string | null, length = 10): string {
  if (!hash) return '—';
  if (hash.length <= length * 2) return hash;
  return `${hash.slice(0, length)}…${hash.slice(-4)}`;
}

export function parseProgress(json?: string | null): { step?: string; percent?: number; [key: string]: unknown } {
  if (!json) return {};
  try {
    return JSON.parse(json);
  } catch {
    return {};
  }
}

export function statusFromValue(value?: number | null, target = 100): string {
  if (value === null || value === undefined) return 'NO_DATA';
  if (value >= target) return 'OK';
  if (value >= 90) return 'WARNING';
  return 'CRITICAL';
}
