<script lang="ts">
  import { page } from '$app/stores';
  import { api } from '$lib/api/client';
  import PageHeader from '$lib/components/ui/page-header.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import Badge from '$lib/components/ui/badge.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import DataTable from '$lib/components/ui/data-table.svelte';
  import Spinner from '$lib/components/ui/spinner.svelte';
  import StatusBadge from '$lib/components/ui/status-badge.svelte';
  import MultiBarChart from '$lib/components/dashboard/multi-bar-chart.svelte';
  import { formatDateTime, formatNumber, formatPercent, truncateHash } from '$lib/utils';
  import type { SPMDetail } from '$lib/types';
  import { ArrowLeft, UploadCloud, XCircle, CheckCircle2 } from '@lucide/svelte';

  let detail = $state<SPMDetail | null>(null);
  let loading = $state(true);
  let error = $state('');

  let kode = $derived($page.params.kode ?? '');
  let periode = $derived($page.url.searchParams.get('periode') ?? '');
  let filterFaskes = $derived($page.url.searchParams.get('faskes_id') ?? '');

  async function load() {
    if (!kode) return;
    loading = true;
    error = '';
    try {
      detail = await api.spmDetail(kode, { periode, faskes_id: filterFaskes || undefined });
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal memuat detail SPM.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    kode;
    periode;
    filterFaskes;
    void load();
  });

  let uploaded = $derived(detail?.faskes.filter((f) => f.uploaded).length ?? 0);
  let notUploaded = $derived(detail?.faskes.filter((f) => !f.uploaded).length ?? 0);
  let chartSeries = $derived(
    detail
      ? [{
          kode: detail.row.label_kode ?? detail.row.kode,
          nama: detail.row.label_nama ?? detail.row.indikator,
          values: detail.faskes.map((f) => f.value ?? null)
        }]
      : []
  );

  async function downloadFile(id: string) {
    try {
      const result = await api.downloadFile(id);
      const url = URL.createObjectURL(result.blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = result.filename;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unduhan gagal.';
    }
  }

  let faskesColumns = [
    { key: 'faskes', label: 'Puskesmas / UPTD', width: '28%' },
    { key: 'status', label: 'Upload', width: '110px' },
    { key: 'value', label: 'Capaian', class: 'text-right', width: '100px' },
    { key: 'numden', label: 'Num/Den', class: 'text-right', width: '110px' },
    { key: 'submission', label: 'Submission / file', width: '220px' },
    { key: 'uploader', label: 'Uploader', width: '140px' },
    { key: 'updated', label: 'Diperbarui', width: '150px' }
  ];
</script>

<PageHeader
  kicker="Dashboard / detail SPM"
  title={detail ? `${detail.row.kode} · ${detail.row.indikator}` : `Detail SPM ${kode}`}
  description="Provenance baris SPM, cakupan upload seluruh puskesmas/UPTD, dan performa antar faskes."
>
  {#snippet actions()}
    <a href="/" class="inline-flex h-10 items-center gap-2 border-2 border-black bg-card px-4 font-mono text-xs font-bold uppercase shadow-nbsm">
      <ArrowLeft size={14} /> Dashboard
    </a>

  {/snippet}
</PageHeader>

{#if error}<div class="mb-4 border-2 border-black bg-destructive p-3 font-mono text-xs font-bold">{error}</div>{/if}

{#if loading}
  <div class="flex items-center justify-center py-20"><Spinner label="Memuat detail SPM…" /></div>
{:else if detail}
  <div class="mb-4 grid gap-3 lg:grid-cols-[1fr_320px]">
    <Card>
      <div class="border-b-2 border-black bg-muted px-4 py-3">
        <div class="nb-kicker">Ringkasan baris</div>
        <div class="mt-2 flex flex-wrap items-center gap-2">
          <StatusBadge kind="display" value={detail.row.display_status} />
          <Badge tone="neutral">{detail.row.label_kode ?? '—'}</Badge>
          {#if detail.row.value_status}<StatusBadge kind="value" value={detail.row.value_status} />{/if}
        </div>
      </div>
      <div class="grid gap-3 p-4 sm:grid-cols-2 lg:grid-cols-4">
        <div><div class="nb-label">Capaian</div><div class="text-2xl font-bold">{formatPercent(detail.row.value)}</div></div>
        <div><div class="nb-label">Target</div><div class="text-2xl font-bold">{formatPercent(detail.row.target, 0)}</div></div>
        <div><div class="nb-label">Gap</div><div class="text-2xl font-bold">{detail.row.gap !== undefined ? `${detail.row.gap > 0 ? '+' : ''}${formatNumber(detail.row.gap, 2)} pp` : '—'}</div></div>
        <div><div class="nb-label">Num/Den</div><div class="text-2xl font-bold">{formatNumber(detail.row.numerator)}/{formatNumber(detail.row.denominator)}</div></div>
      </div>
      <div class="border-t-2 border-black bg-card px-4 py-2 font-mono text-[10px] text-muted-foreground">
        data version {detail.row.data_version ?? '—'} · source hash {truncateHash(detail.row.source_hash)} · calculation run {detail.row.calculation_run_id ?? '—'}
      </div>
    </Card>
    <Card>
      <div class="border-b-2 border-black bg-secondary px-4 py-3">
        <div class="nb-kicker">Cakupan upload</div>
        <h2 class="mt-1 text-lg">Puskesmas / UPTD</h2>
      </div>
      <div class="space-y-2 p-4 font-mono text-xs">
        <div class="flex items-center justify-between border-2 border-black bg-success px-2 py-1"><span><CheckCircle2 size={13} /> Sudah upload</span><b>{uploaded}</b></div>
        <div class="flex items-center justify-between border-2 border-black bg-warning px-2 py-1"><span><XCircle size={13} /> Belum upload</span><b>{notUploaded}</b></div>
        <div class="flex items-center justify-between border-2 border-black bg-card px-2 py-1"><span><UploadCloud size={13} /> Total faskes</span><b>{detail.faskes.length}</b></div>
      </div>
    </Card>
  </div>

  <Card class="mb-4">
    <div class="flex flex-wrap items-center justify-between gap-3 border-b-2 border-black bg-card px-4 py-3">
      <div>
        <div class="nb-kicker">Performa antar puskesmas</div>
        <h2 class="mt-1 text-xl">Capaian {detail.row.label_kode ?? detail.row.kode} per faskes</h2>
      </div>
      <Badge tone="info">{detail.row.indikator}</Badge>
    </div>
    <div class="p-4">
      <MultiBarChart faskes={detail.faskes.map((f) => ({ id: f.faskes_id, kode: f.faskes_kode, nama: f.faskes_nama }))} series={chartSeries} height={320} />
    </div>
  </Card>

  <Card class="mb-4">
    <div class="border-b-2 border-black bg-muted px-4 py-3">
      <div class="nb-kicker">Cakupan & provenance per puskesmas</div>
      <h2 class="mt-1 text-xl">Siapa sudah upload, siapa belum</h2>
    </div>
    <DataTable columns={faskesColumns} rows={detail.faskes}>
      {#snippet cell(f, key)}
        {#if key === 'faskes'}
          <div class="font-bold">{f.faskes_nama}</div>
          <div class="text-[10px] text-muted-foreground">{f.faskes_kode}</div>
        {:else if key === 'status'}
          <Badge tone={f.uploaded ? 'success' : 'warning'}>{f.uploaded ? 'Sudah upload' : 'Belum upload'}</Badge>
        {:else if key === 'value'}
          <span class="font-bold">{formatPercent(f.value)}</span>
        {:else if key === 'numden'}
          {formatNumber(f.numerator)}/{formatNumber(f.denominator)}
        {:else if key === 'submission'}
          {#if f.submission_kode}
            <a href={`/submissions/${f.submission_id}`} class="font-bold underline">{f.submission_kode}</a>
            <div class="mt-1 flex items-center gap-1 text-[10px] text-muted-foreground">
              {f.file_name ?? 'file'}
              {#if f.file_id}<button class="underline" onclick={() => downloadFile(f.file_id!)}>unduh</button>{/if}
            </div>
          {:else}
            <span class="text-[10px] text-muted-foreground">Belum ada submission final</span>
          {/if}
        {:else if key === 'uploader'}
          {f.uploader_nama ?? '—'}
        {:else if key === 'updated'}
          {formatDateTime(f.updated_at)}
        {/if}
      {/snippet}
    </DataTable>
  </Card>

  {#if detail.provenance && detail.provenance.length > 0}
    <Card>
      <div class="border-b-2 border-black bg-muted px-4 py-3">
        <div class="nb-kicker">Provenance komponen</div>
        <h2 class="mt-1 text-xl">Angka ini berasal dari mana</h2>
      </div>
      <div class="overflow-x-auto">
        <table class="nb-table">
          <thead><tr><th>Komponen</th><th>Tipe</th><th>Num/Den</th><th>Nilai</th><th>Bobot</th><th>Kontribusi</th></tr></thead>
          <tbody>
            {#each detail.provenance as p (p.id)}
              <tr>
                <td>{p.indikator_nama ?? p.label_nama ?? p.component_id}</td>
                <td>{p.component_type}</td>
                <td>{formatNumber(p.numerator)}/{formatNumber(p.denominator)}</td>
                <td>{formatPercent(p.component_value)}</td>
                <td>{p.weight}</td>
                <td>{formatPercent(p.contribution)}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </Card>
  {/if}
{/if}
