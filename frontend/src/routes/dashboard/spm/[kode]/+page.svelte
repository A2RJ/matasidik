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
  import { formatDateTime, formatNumber, formatMetric, formatGap } from '$lib/utils';
  import type { SPMDetail } from '$lib/types';
  import { ArrowLeft, UploadCloud, XCircle, CheckCircle2 } from '@lucide/svelte';

  let detail = $state<SPMDetail | null>(null);
  let loading = $state(true);
  let error = $state('');

  let kode = $derived($page.params.kode ?? '');
  let periode = $derived($page.url.searchParams.get('periode') ?? '');
  let filterFaskes = $derived($page.url.searchParams.get('faskes_id') ?? '');
  let filterScopeType = $derived($page.url.searchParams.get('scope_type') ?? '');
  let filterScopeID = $derived($page.url.searchParams.get('scope_id') ?? '');

  async function load() {
    if (!kode) return;
    loading = true;
    error = '';
    try {
      detail = await api.spmDetail(kode, {
        periode,
        faskes_id: filterFaskes || undefined,
        scope_type: filterScopeType || undefined,
        scope_id: filterScopeID || undefined
      });
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
    filterScopeType;
    filterScopeID;
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

  let wilayahColumns = [
    { key: 'wilayah', label: 'Wilayah', width: '34%' },
    { key: 'level', label: 'Level', width: '100px' },
    { key: 'status', label: 'Data', width: '110px' },
    { key: 'value', label: 'Capaian', class: 'text-right', width: '100px' },
    { key: 'numden', label: 'Jumlah/Total', class: 'text-right', width: '110px' },
    { key: 'updated', label: 'Diperbarui', width: '150px' }
  ];

  let faskesColumns = [
    { key: 'faskes', label: 'Puskesmas / UPTD', width: '28%' },
    { key: 'status', label: 'Upload', width: '110px' },
    { key: 'value', label: 'Capaian', class: 'text-right', width: '100px' },
    { key: 'numden', label: 'Jumlah/Total', class: 'text-right', width: '110px' },
    { key: 'submission', label: 'Berkas sumber', width: '220px' },
    { key: 'uploader', label: 'Pengunggah', width: '140px' },
    { key: 'updated', label: 'Diperbarui', width: '150px' }
  ];

  const componentLabels: Record<string, string> = {
    IKU: 'IKU',
    INDIKATOR: 'Indikator',
    COMPONENT: 'Komponen'
  };
</script>

<PageHeader
  kicker="Dashboard / detail SPM"
  title={detail ? `${detail.row.kode} · ${detail.row.indikator}` : `Detail SPM ${kode}`}
  description="Rincian sumber angka, cakupan upload seluruh puskesmas/UPTD, dan performa antar faskes."
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
        <div><div class="nb-label">Capaian</div><div class="text-2xl font-bold">{formatMetric(detail.row.value, detail.row)}</div></div>
        <div><div class="nb-label">Target</div><div class="text-2xl font-bold">{formatMetric(detail.row.target, detail.row, 0)}</div></div>
        <div><div class="nb-label">Gap</div><div class="text-2xl font-bold">{formatGap(detail.row.gap, detail.row)}</div></div>
        <div><div class="nb-label">Jumlah/Total</div><div class="text-2xl font-bold">{!detail?.row.tipe_nilai || detail?.row.tipe_nilai === 'RATIO' ? `${formatNumber(detail.row.numerator)}/${formatNumber(detail.row.denominator)}` : formatMetric(detail.row.value, detail.row)}</div></div>
      </div>
      <div class="border-t-2 border-black bg-card px-4 py-2 text-[10px] text-muted-foreground">
        Terakhir diperbarui {formatDateTime(detail.row.updated_at)}
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
    <div class="border-b-2 border-black bg-muted px-4 py-3">
      <div class="nb-kicker">Cakupan per wilayah</div>
      <h2 class="mt-1 text-xl">Desa, kecamatan, dan kabupaten</h2>
    </div>
    <DataTable columns={wilayahColumns} rows={detail.wilayah}>
      {#snippet cell(w, key)}
        {#if key === 'wilayah'}
          <div class="font-bold">{w.scope_nama}</div>
          <div class="text-[10px] text-muted-foreground">{w.scope_kode}</div>
        {:else if key === 'level'}
          <Badge tone="neutral">{w.scope_type}</Badge>
        {:else if key === 'status'}
          <Badge tone={w.uploaded ? 'success' : 'warning'}>{w.uploaded ? 'Ada data' : 'Belum ada'}</Badge>
        {:else if key === 'value'}
          <span class="font-bold">{formatMetric(w.value, detail?.row)}</span>
        {:else if key === 'numden'}
          {!detail?.row.tipe_nilai || detail?.row.tipe_nilai === 'RATIO' ? `${formatNumber(w.numerator)}/${formatNumber(w.denominator)}` : formatMetric(w.value, detail?.row)}
        {:else if key === 'updated'}
          {formatDateTime(w.updated_at)}
        {/if}
      {/snippet}
    </DataTable>
  </Card>

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
      <div class="nb-kicker">Cakupan & sumber data per puskesmas</div>
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
          <span class="font-bold">{formatMetric(f.value, detail?.row)}</span>
        {:else if key === 'numden'}
          {!detail?.row.tipe_nilai || detail?.row.tipe_nilai === 'RATIO' ? `${formatNumber(f.numerator)}/${formatNumber(f.denominator)}` : formatMetric(f.value, detail?.row)}
        {:else if key === 'submission'}
          {#if f.submission_kode}
            <a href={`/submissions/${f.submission_id}`} class="font-bold underline">{f.submission_kode}</a>
            <div class="mt-1 flex items-center gap-1 text-[10px] text-muted-foreground">
              {f.file_name ?? 'Berkas BNBA'}
              {#if f.file_id}<button class="underline" onclick={() => downloadFile(f.file_id!)}>unduh</button>{/if}
            </div>
          {:else}
            <span class="text-[10px] text-muted-foreground">Belum ada berkas final</span>
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
        <div class="nb-kicker">Sumber komponen angka</div>
        <h2 class="mt-1 text-xl">Asal angka</h2>
      </div>
      <div class="overflow-x-auto">
        <table class="nb-table">
          <thead><tr><th>Komponen</th><th>Jenis</th><th>Jumlah/Total</th><th>Nilai</th><th>Bobot</th><th>Kontribusi</th></tr></thead>
          <tbody>
            {#each detail.provenance as p (p.id)}
              <tr>
                <td>{p.indikator_nama ?? p.label_nama ?? p.component_id}</td>
                <td>{componentLabels[p.component_type] ?? p.component_type}</td>
                <td>{!detail?.row.tipe_nilai || detail?.row.tipe_nilai === 'RATIO' ? `${formatNumber(p.numerator)}/${formatNumber(p.denominator)}` : formatMetric(p.component_value, detail?.row)}</td>
                <td>{formatMetric(p.component_value, detail?.row)}</td>
                <td>{p.weight}</td>
                <td>{formatMetric(p.contribution, detail?.row)}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </Card>
  {/if}
{/if}
