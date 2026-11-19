import type { Role } from '$lib/types';

export interface RouteAccess {
  href: string;
  label: string;
  roles: Role[] | null;
}

/**
 * Satu sumber aturan akses halaman. `null` berarti semua pengguna yang sudah
 * masuk boleh membuka halaman. Menu sidebar, penjaga rute, dan judul halaman
 * memakai daftar ini agar tidak ada tampilan yang berbeda antar komponen.
 */
export const ROUTE_ACCESS: RouteAccess[] = [
  { href: '/upload', label: 'Upload BNBA', roles: ['pkm_operator', 'pkm_pj', 'admin'] },
  {
    href: '/submissions',
    label: 'Submission',
    roles: ['pkm_operator', 'pkm_pj', 'pkm_approver', 'dinas_verifier', 'kabid', 'sekretaris', 'pimpinan', 'admin']
  },
  { href: '/iku-requests', label: 'Pengajuan IKU', roles: ['pkm_operator', 'pkm_approver', 'admin'] },
  { href: '/master', label: 'Master Data', roles: ['kabid', 'sekretaris', 'admin'] },
  { href: '/factors', label: 'Faktor & TL', roles: ['kabid', 'sekretaris', 'pimpinan', 'admin'] },
  { href: '/audit', label: 'Jejak Audit', roles: ['pimpinan', 'admin'] },
  { href: '/admin', label: 'Administrasi', roles: ['admin'] }
];

export function routeAccessFor(pathname: string): RouteAccess | undefined {
  return ROUTE_ACCESS.find((item) => pathname === item.href || pathname.startsWith(`${item.href}/`));
}

export function canAccessPath(role: Role | undefined, pathname: string): boolean {
  if (pathname === '/login' || pathname === '/user-guide' || pathname.startsWith('/user-guide/')) {
    return true;
  }
  if (!role) return false;
  const access = routeAccessFor(pathname);
  return !access || access.roles === null || access.roles.includes(role);
}

export function accessLabelFor(pathname: string): string {
  return routeAccessFor(pathname)?.label ?? 'halaman ini';
}
