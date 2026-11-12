<script lang="ts">
  import { api } from '$lib/api/client';
  import { periods, selectedPeriod } from '$lib/stores/session';
  import PageHeader from '$lib/components/ui/page-header.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Select from '$lib/components/ui/select.svelte';
  import DataTable from '$lib/components/ui/data-table.svelte';
  import EmptyState from '$lib/components/ui/empty-state.svelte';
  import Spinner from '$lib/components/ui/spinner.svelte';
  import StatusBadge from '$lib/components/ui/status-badge.svelte';
  import { PAGE_COPY, EMPTY_STATE } from '$lib/constants/labels';
  import { formatDateTime, truncateHash } from '$lib/utils';
  import type { Submission } from '$lib/types';
  import { Upload, RefreshCw, ArrowRight } from '@lucide/svelte';

  let items = $state<Submission[]>([]);
  let loading = $state(true);
  let error = $state('');

  let periodOptions = $derived($periods.map((p) => ({ value: p.id, label: p.kode })));
  let columns = [
    { key: 'kode', label: 'Kode' },
    { key: 'faskes', label: 'Faskes' },
    { key: 'uploader', label: 'Uploader' },
    { key: 'periode', label: 'Periode' },
    { key: 'status', label: 'Status' },
    { key: 'file', label: 'File hash' },
    { key: 'updated', label: 'Diperbarui' },
    { key: 'action', label: 'Aksi', class: 'w-24' }
  ];

  async function load() {
    loading = true;
    error = '';
    try {
      items = await api.submissions({ periode: $selectedPeriod || undefined });
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal memuat submission.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if ($selectedPeriod) void load();
  });
</script>

<PageHeader title={PAGE_COPY.submissions.title} description={PAGE_COPY.submissions.description}>
  {#snippet actions()}
    <Button variant="outline" onclick={load}><RefreshCw size={14} /> Muat ulang</Button>
    <a href="/upload" class="inline-flex h-10 items-center gap-2 border-2 border-black bg-primary px-4 font-mono text-xs font-bold uppercase shadow-nb">
      <Upload size={14} /> Upload BNBA
    </a>
  {/snippet}
</PageHeader>

{#if error}
  <div class="mb-4 border-2 border-black bg-destructive p-3 font-mono text-xs font-bold">{error}</div>
{/if}

<Card>
  <div class="flex flex-wrap items-center justify-between gap-3 border-b-2 border-black bg-muted px-4 py-3">
    <div>
      <div class="nb-kicker">Daftar pengajuan</div>
      <h2 class="mt-1 text-xl">Submission per periode</h2>
    </div>
    <div class="min-w-52">
      <Select bind:value={$selectedPeriod} options={periodOptions} placeholder="Pilih periode" />
    </div>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-16"><Spinner label="Memuat submission…" /></div>
  {:else if items.length === 0}
    <div class="p-4">
      <EmptyState title="Belum ada submission" description={EMPTY_STATE.submissions}>
        {#snippet action()}
          <a href="/upload" class="inline-flex h-10 items-center gap-2 border-2 border-black bg-primary px-4 font-mono text-xs font-bold uppercase shadow-nb">
            <Upload size={14} /> Mulai upload
          </a>
        {/snippet}
      </EmptyState>
    </div>
  {:else}
    <DataTable {columns} rows={items} onrowclick={(row) => (window.location.href = `/submissions/${row.id}`)}>
      {#snippet cell(item, key)}
        {#if key === 'kode'}
          <span class="font-bold">{item.kode}</span>
        {:else if key === 'faskes'}
          {item.faskes_nama ?? item.faskes_id}
        {:else if key === 'uploader'}
          {item.uploader_nama ?? item.uploader_id}
        {:else if key === 'periode'}
          {item.periode_kode ?? item.periode_id}
        {:else if key === 'status'}
          <StatusBadge kind="submission" value={item.status} />
        {:else if key === 'file'}
          <span class="font-mono text-[10px]">{truncateHash(item.file_hash)}</span>
        {:else if key === 'updated'}
          {formatDateTime(item.updated_at)}
        {:else if key === 'action'}
          <a href={`/submissions/${item.id}`} class="inline-flex items-center gap-1 border-2 border-black bg-card px-2 py-1 font-mono text-[10px] font-bold uppercase shadow-nbsm">
            Buka <ArrowRight size={12} />
          </a>
        {/if}
      {/snippet}
    </DataTable>
  {/if}
</Card>
