import { writable } from 'svelte/store';
import type { Bootstrap, Periode, User } from '$lib/types';

export const bootstrap = writable<Bootstrap | null>(null);
export const currentActor = writable<User | null>(null);
export const periods = writable<Periode[]>([]);
export const selectedPeriod = writable<string>('');

export function applyBootstrap(data: Bootstrap) {
  bootstrap.set(data);
  periods.set(data.periods);
  if (data.periods.length > 0) {
    const fallback = data.default_period_id || data.periods[0].id;
    selectedPeriod.update((current) => current || fallback);
  }
}
