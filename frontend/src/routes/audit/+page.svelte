<script lang="ts">
  import { api } from '$lib/api/client';
  import PageHeader from '$lib/components/ui/page-header.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Select from '$lib/components/ui/select.svelte';
  import DataTable from '$lib/components/ui/data-table.svelte';
  import Spinner from '$lib/components/ui/spinner.svelte';
  import { PAGE_COPY } from '$lib/constants/labels';
  import { formatDateTime, truncateHash } from '$lib/utils';
  import type { AuditLog } from '$lib/types';
  import { RefreshCw } from '@lucide/svelte';

  let logs = $state<AuditLog[]>([]);
  let loading = $state(true);
  let error = $state('');
  let entityType = $state('');

  let columns = [
    { key: 'created', label: 'Waktu', class: 'w-44' },
    { key: 'actor', label: 'Aktor', class: 'w-40' },
    { key: 'action', label: 'Aksi' },
    { key: 'entity', label: 'Entitas', class: 'w-52' },
    { key: 'reason', label: 'Alasan / catatan' },
    { key: 'request', label: 'Request', class: 'w-28' }
  ];

  const entityOptions = [
    { value: 'submission', label: 'Submission' },
    { value: 'submission_row', label: 'Submission row' },
    { value: 'iku_value', label: 'IKU value' },
    { value: 'tindak_lanjut', label: 'Tindak lanjut' },
    { value: 'label', label: 'Label' },
    { value: 'indikator', label: 'Indikator' }
  ];

  async function load() {
    loading = true;
    error = '';
    try {
      logs = await api.audit({ entity_type: entityType || undefined, limit: '200' });
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal memuat audit.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    const filter = entityType;
    void load();
  });
</script>

<PageHeader title={PAGE_COPY.audit.title} description={PAGE_COPY.audit.description}>
  {#snippet actions()}
    <Button variant="outline" onclick={load}><RefreshCw size={14} /> Muat ulang</Button>
  {/snippet}
</PageHeader>

{#if error}<div class="mb-4 border-2 border-black bg-destructive p-3 font-mono text-xs font-bold">{error}</div>{/if}

<Card>
  <div class="flex flex-wrap items-end justify-between gap-3 border-b-2 border-black bg-muted px-4 py-3">
    <div>
      <div class="nb-kicker">Jejak keputusan</div>
      <h2 class="mt-1 text-xl">Audit append-only</h2>
    </div>
    <div class="min-w-56">
      <Select bind:value={entityType} options={entityOptions} placeholder="Semua entitas" />
    </div>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-16"><Spinner label="Memuat audit…" /></div>
  {:else}
    <DataTable {columns} rows={logs}>
      {#snippet cell(item: AuditLog, key)}
        {#if key === 'created'}{formatDateTime(item.created_at)}
        {:else if key === 'actor'}
          <div>{item.actor_nama ?? item.actor_id ?? 'Sistem'}</div>
          {#if item.impersonator_nama}<div class="text-[9px] text-muted-foreground">via {item.impersonator_nama}</div>{/if}
        {:else if key === 'action'}<span class="font-bold">{item.action}</span>
        {:else if key === 'entity'}
          <div>{item.entity_type}</div>
          <div class="text-[10px] text-muted-foreground">{truncateHash(item.entity_id, 8)}</div>
        {:else if key === 'reason'}{item.reason ?? '—'}
        {:else if key === 'request'}<span class="text-[10px]">{truncateHash(item.request_id, 6)}</span>
        {/if}
      {/snippet}
    </DataTable>
  {/if}
</Card>
