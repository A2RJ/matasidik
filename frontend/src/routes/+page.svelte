<script lang="ts">
  import { goto } from '$app/navigation';
  import { api } from '$lib/api/client';
  import { bootstrap, periods, selectedPeriod, currentActor } from '$lib/stores/session';
  import PageHeader from '$lib/components/ui/page-header.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import Select from '$lib/components/ui/select.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Badge from '$lib/components/ui/badge.svelte';
  import DataTable from '$lib/components/ui/data-table.svelte';
  import EmptyState from '$lib/components/ui/empty-state.svelte';
  import Spinner from '$lib/components/ui/spinner.svelte';
  import StatusBadge from '$lib/components/ui/status-badge.svelte';
  import StatCard from '$lib/components/dashboard/stat-card.svelte';
  import MultiBarChart from '$lib/components/dashboard/multi-bar-chart.svelte';
  import { PAGE_COPY, TABLE_LABELS, EMPTY_STATE } from '$lib/constants/labels';
  import { DISPLAY_STATUS_LABELS, STATUS_TARGET_NOTE } from '$lib/constants/enums';
  import { formatPercent, formatNumber } from '$lib/utils';
  import type { DashboardItem, DashboardSummary, IKUChart, MonthlySPMReport, MonthlySPMRow } from '$lib/types';
  import { Download, RefreshCw, Target, Activity, AlertTriangle, Database, Network, FileSpreadsheet, FileText } from '@lucide/svelte';

  let items = $state<DashboardItem[]>([]);
  let summary = $state<DashboardSummary | null>(null);
  let chart = $state<IKUChart | null>(null);
  let loading = $state(true);
  let exporting = $state('');
  let error = $state('');
  let monthlyReport = $state<MonthlySPMReport | null>(null);
  let monthlyLoading = $state(true);
  let monthlyError = $state('');
  let monthlyExporting = $state('');
  let filterFaskes = $state('');
  let filterStatus = $state('');
  let filterCategory = $state('SPM');

  let periodOptions = $derived($periods.map((p) => ({ value: p.id, label: p.kode })));
  let faskesOptions = $derived(($bootstrap?.faskes ?? []).map((f) => ({ value: f.id, label: f.nama })));
  let categoryOptions = $derived([
    { value: 'ALL', label: 'Semua' },
    ...($bootstrap?.categories ?? []).filter((c) => c.aktif).map((c) => ({ value: c.kode, label: c.nama }))
  ]);
  let statusOptions = $derived(
    Object.entries(DISPLAY_STATUS_LABELS).map(([value, label]) => ({ value, label }))
  );
  let canExport = $derived($currentActor && ['pimpinan', 'kabid', 'sekretaris', 'dinas_verifier', 'admin'].includes($currentActor.role));
  let baselineHeader = $derived(
    items[0]?.baseline_year ? `Capaian ${items[0].baseline_year}` : TABLE_LABELS.baselinePrev
  );
  let isSPM = $derived(filterCategory === 'SPM');
  let isAll = $derived(filterCategory === 'ALL');
  const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des'];
  let monthlyYear = $derived(
    Number($periods.find((p) => p.id === $selectedPeriod)?.tahun ?? new Date().getFullYear())
  );
  let monthlyColumns = $derived([
    { key: 'kode', label: 'No', width: '56px' },
    { key: 'indikator', label: 'Indikator', width: '30%' },
    { key: 'target_tahunan', label: `Target ${monthlyYear}`, class: 'text-right', width: '110px' },
    { key: 'target_bulanan', label: 'Target Bulanan', class: 'text-right', width: '100px' },
    ...monthNames.flatMap((name, index) => [
      { key: `m${index + 1}`, label: name, class: 'text-right', width: '64px' },
      { key: `s${index + 1}`, label: 'Status', width: '112px' }
    ])
  ]);

  let columns = $derived(isSPM ? [
    { key: 'kode', label: 'No', width: '56px' },
    { key: 'indikator', label: 'Indikator', width: '32%' },
    { key: 'target', label: TABLE_LABELS.target, class: 'text-right', width: '80px' },
    { key: 'baselinePrev', label: baselineHeader, class: 'text-right', width: '100px' },
    { key: 'baselineTW1', label: TABLE_LABELS.baselineTW1, class: 'text-right', width: '100px' },
    { key: 'value', label: 'Capaian', class: 'text-right', width: '100px' },
    { key: 'status', label: TABLE_LABELS.status, width: '150px' },
    { key: 'gap', label: TABLE_LABELS.gap, class: 'text-right', width: '90px' },
    { key: 'trace', label: TABLE_LABELS.action, width: '96px' }
  ] : [
    { key: 'kode', label: 'Kode', width: '110px' },
    { key: 'indikator', label: 'Label / Indikator', width: '34%' },
    ...(isAll ? [{ key: 'category', label: 'Kelompok', width: '150px' }] : []),
    { key: 'target', label: TABLE_LABELS.target, class: 'text-right', width: '90px' },
    { key: 'value', label: 'Capaian', class: 'text-right', width: '110px' },
    { key: 'status', label: TABLE_LABELS.status, width: '150px' },
    { key: 'gap', label: TABLE_LABELS.gap, class: 'text-right', width: '90px' },
    { key: 'owner', label: 'Bidang/Timker', width: '160px' }
  ]);

  async function load() {
    const period = $selectedPeriod;
    if (!period) return;
    loading = true;
    error = '';
    try {
      if (filterCategory === 'SPM') {
        const [dashboard, chartData] = await Promise.all([
          api.dashboard({ periode: period, faskes_id: filterFaskes || undefined, status: filterStatus || undefined }),
          api.ikuByFaskes({ periode: period, faskes_id: filterFaskes || undefined })
        ]);
        items = dashboard.items;
        summary = dashboard.summary;
        chart = chartData;
      } else {
        const result = await api.dashboardCategory(filterCategory, {
          periode: period,
          faskes_id: filterFaskes || undefined,
          status: filterStatus || undefined
        });
        items = result.items;
        summary = result.summary;
        chart = null;
      }
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal memuat dashboard.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    const period = $selectedPeriod;
    const faskes = filterFaskes;
    const status = filterStatus;
    const category = filterCategory;
    if (period) void load();
  });

  async function loadMonthly() {
    if (!monthlyYear) return;
    monthlyLoading = true;
    monthlyError = '';
    try {
      monthlyReport = await api.spmMonthly(monthlyYear);
    } catch (err) {
      monthlyError = err instanceof Error ? err.message : 'Gagal memuat laporan bulanan.';
    } finally {
      monthlyLoading = false;
    }
  }

  $effect(() => {
    const year = monthlyYear;
    if (year) void loadMonthly();
  });

  async function exportMonthly(format: 'xlsx' | 'pdf') {
    monthlyExporting = format;
    monthlyError = '';
    try {
      const result = await api.spmMonthlyExport(monthlyYear, format);
      const url = URL.createObjectURL(result.blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = result.filename;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      monthlyError = err instanceof Error ? err.message : 'Export laporan bulanan gagal.';
    } finally {
      monthlyExporting = '';
    }
  }

  function openDetail(kode: string) {
    const params = new URLSearchParams({ periode: $selectedPeriod });
    if (filterFaskes) params.set('faskes_id', filterFaskes);
    void goto(`/dashboard/spm/${kode}?${params.toString()}`);
  }

  async function downloadTemplate() {
    try {
      const result = await api.downloadTemplate();
      const url = URL.createObjectURL(result.blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = result.filename;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Template gagal diunduh.';
    }
  }

  async function exportData(format: 'xlsx' | 'pdf') {
    exporting = format;
    error = '';
    try {
      const result = await api.exportSPM(format, $selectedPeriod);
      const url = URL.createObjectURL(result.blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = result.filename;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Export gagal.';
    } finally {
      exporting = '';
    }
  }
</script>

<PageHeader title={PAGE_COPY.dashboard.title} description={PAGE_COPY.dashboard.description}>
  {#snippet actions()}
    <Button variant="outline" onclick={() => load()}><RefreshCw size={14} /> Muat ulang</Button>
    {#if canExport}
      <Button variant="secondary" disabled={exporting !== ''} onclick={() => exportData('xlsx')}><FileSpreadsheet size={14} /> Excel</Button>
      <Button variant="secondary" disabled={exporting !== ''} onclick={() => exportData('pdf')}><FileText size={14} /> PDF</Button>
    {/if}
    <Button variant="default" onclick={downloadTemplate}><Download size={14} /> Template BNBA</Button>
  {/snippet}
</PageHeader>

{#if error}
  <div class="mb-4 border-2 border-black bg-destructive p-3 font-mono text-xs font-bold">{error}</div>
{/if}

<div class="mb-4 border-2 border-black bg-secondary p-3 font-mono text-[11px] font-bold">
  {STATUS_TARGET_NOTE}
</div>

<div class="mb-6 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
  <StatCard
    label={isSPM ? 'IKU induk · Kesehatan Ibu' : 'Total baris kategori'}
    value={isSPM ? (summary?.parent_value !== undefined ? formatPercent(summary.parent_value) : '—') : `${summary?.total ?? 0}`}
    hint={isSPM ? `Status ${summary?.parent_status ?? 'NO_DATA'} · RATIO_OF_SUMS` : `Kategori ${filterCategory}`}
    tone="primary"
  >
    {#snippet icon()}<Target size={18} />{/snippet}
  </StatCard>
  <StatCard
    label="Sesuai target"
    value={`${summary?.on_target ?? 0}`}
    hint={`dari ${summary?.total ?? 0} baris`}
    tone="success"
  >
    {#snippet icon()}<Activity size={18} />{/snippet}
  </StatCard>
  <StatCard
    label="Di bawah target"
    value={`${summary?.below_target ?? 0}`}
    hint={`rata-rata capaian ${formatPercent(summary?.average)}`}
    tone="warning"
  >
    {#snippet icon()}<AlertTriangle size={18} />{/snippet}
  </StatCard>
  <StatCard
    label="Belum ada data"
    value={`${summary?.no_data ?? 0}`}
    hint="Baris tanpa IKU_VALUE final"
    tone="muted"
  >
    {#snippet icon()}<Database size={18} />{/snippet}
  </StatCard>
</div>

<Card class="mb-6">
  <div class="flex flex-wrap items-end justify-between gap-3 border-b-2 border-black bg-muted px-4 py-3">
    <div>
      <div class="nb-kicker">Filter dashboard</div>
      <h2 class="mt-1 text-xl">Periode & scope</h2>
    </div>
    <div class="grid gap-2 sm:grid-cols-4">
      <div class="min-w-44">
        <label class="nb-label" for="period">Periode</label>
        <Select id="period" bind:value={$selectedPeriod} options={periodOptions} placeholder="Pilih periode" />
      </div>
      <div class="min-w-48">
        <label class="nb-label" for="category">Kelompok</label>
        <Select id="category" bind:value={filterCategory} options={categoryOptions} placeholder="Pilih kelompok" />
      </div>
      <div class="min-w-48">
        <label class="nb-label" for="faskes">Faskes/UPTD</label>
        <Select id="faskes" bind:value={filterFaskes} options={faskesOptions} placeholder="Semua faskes" allowEmpty />
      </div>
      <div class="min-w-40">
        <label class="nb-label" for="status">Status</label>
        <Select id="status" bind:value={filterStatus} options={statusOptions} placeholder="Semua status" allowEmpty />
      </div>
    </div>
  </div>
</Card>

{#if isSPM}
  <Card class="mb-6">
    <div class="flex flex-wrap items-center justify-between gap-3 border-b-2 border-black bg-card px-4 py-3">
      <div>
        <div class="nb-kicker">Performa antar puskesmas/UPTD</div>
        <h2 class="mt-1 text-xl">Multi-series IKU</h2>
        <p class="mt-1 text-xs text-muted-foreground">Setiap warna adalah satu IKU leaf; bandingkan capaian antar faskes pada periode aktif.</p>
      </div>
      <Badge tone="info">Legends: seluruh IKU leaf</Badge>
    </div>
    <CardContent>
      <MultiBarChart faskes={chart?.faskes ?? []} series={chart?.series ?? []} />
    </CardContent>
  </Card>
{/if}

<Card>
  <div class="flex flex-wrap items-center justify-between gap-3 border-b-2 border-black bg-muted px-4 py-3">
    <div>
      <div class="nb-kicker">{isSPM ? 'Output SPM' : 'Output kelompok'}</div>
      <h2 class="mt-1 text-xl">{isSPM ? 'Capaian SPM per indikator' : isAll ? 'Capaian semua kelompok' : 'Capaian per label/owner'}</h2>
    </div>
    <div class="flex flex-wrap gap-2">
      <Badge tone="final">Data final</Badge>
      <Badge tone="info">Traceable</Badge>
      {#if summary}<Badge tone="neutral">Total {summary.total} baris</Badge>{/if}
    </div>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-16"><Spinner label="Memuat dashboard…" /></div>
  {:else if items.length === 0}
    <div class="p-4"><EmptyState title="Dashboard kosong" description={EMPTY_STATE.dashboard} /></div>
  {:else}
    <DataTable {columns} rows={items} fixed onrowclick={isSPM ? (row) => openDetail(row.kode) : undefined}>
      {#snippet cell(item: DashboardItem, key)}
        {#if key === 'kode'}
          <span class="font-bold">{item.kode}</span>
        {:else if key === 'indikator'}
          <div class="font-semibold text-foreground">{item.indikator}</div>
        {:else if key === 'category'}
          <span class="font-mono text-[10px]">{item.category_nama ?? item.category_kode ?? '—'}</span>
        {:else if key === 'target'}
          {item.target > 0 ? formatPercent(item.target, 0) : '—'}
        {:else if key === 'baselinePrev'}
          {formatPercent(item.baseline_prev, 0)}
        {:else if key === 'baselineTW1'}
          {formatPercent(item.baseline_tw1, 1)}
        {:else if key === 'value'}
          <span class="font-bold">{formatPercent(item.value)}</span>
          {#if item.numerator !== undefined && item.denominator !== undefined}
            <div class="text-[10px] text-muted-foreground">{formatNumber(item.numerator)}/{formatNumber(item.denominator)}</div>
          {/if}
        {:else if key === 'status'}
          <StatusBadge kind="display" value={item.display_status} />
        {:else if key === 'gap'}
          {item.gap !== undefined ? `${item.gap > 0 ? '+' : ''}${formatNumber(item.gap, 2)} pp` : '—'}
        {:else if key === 'owner'}
          {item.owner_nama ? `${item.owner_nama} (${item.owner_jenis ?? 'owner'})` : '—'}
        {:else if key === 'trace'}
          {#if isSPM}
            <Button size="sm" variant="outline" onclick={(event) => { event.stopPropagation(); openDetail(item.kode); }}>
              <Network size={12} /> Detail
            </Button>
          {/if}
        {/if}
      {/snippet}
    </DataTable>
    <div class="border-t-2 border-black bg-card px-4 py-2 font-mono text-[10px] text-muted-foreground">
      {isSPM ? 'Klik baris untuk melihat provenance per puskesmas, cakupan upload, dan chart performa.' : 'Data mengikuti kategori dan filter yang dipilih; nilai tetap berasal dari label/indikator master.'}
    </div>
  {/if}
</Card>

<Card class="mt-6">
  <div class="flex flex-wrap items-center justify-between gap-3 border-b-2 border-black bg-card px-4 py-3">
    <div>
      <div class="nb-kicker">Laporan turunan IKU/SPM</div>
      <h2 class="mt-1 text-xl">Data Absolut Kabupaten {monthlyYear}</h2>
      <p class="mt-1 text-xs text-muted-foreground">
        Target layanan tahunan, target bulanan, jumlah layanan per bulan, dan status bulanan terhadap target bulanan.
      </p>
    </div>
    <div class="flex flex-wrap gap-2">
      <Button variant="secondary" disabled={monthlyExporting !== '' || monthlyLoading} onclick={() => exportMonthly('xlsx')}>
        <FileSpreadsheet size={14} /> Excel
      </Button>
      <Button variant="secondary" disabled={monthlyExporting !== '' || monthlyLoading} onclick={() => exportMonthly('pdf')}>
        <FileText size={14} /> PDF
      </Button>
    </div>
  </div>

  {#if monthlyError}
    <div class="border-b-2 border-black bg-destructive px-4 py-2 font-mono text-xs font-bold">{monthlyError}</div>
  {/if}
  {#if monthlyLoading}
    <div class="flex items-center justify-center py-12"><Spinner label="Memuat laporan bulanan…" /></div>
  {:else if !monthlyReport || monthlyReport.rows.length === 0}
    <div class="p-4"><EmptyState title="Laporan bulanan kosong" description="Belum ada data IKU/SPM final untuk tahun ini." /></div>
  {:else}
    <DataTable columns={monthlyColumns} rows={monthlyReport.rows} minWidth="2600px" nowrap>
      {#snippet cell(row: MonthlySPMRow, key)}
        {#if key === 'kode'}
          <span class="font-bold">{row.kode}</span>
        {:else if key === 'indikator'}
          <span class="font-semibold">{row.indikator}</span>
        {:else if key === 'target_tahunan'}
          {formatNumber(row.target_tahunan, 0)}
        {:else if key === 'target_bulanan'}
          {formatNumber(row.target_bulanan, 2)}
        {:else if key.startsWith('m')}
          {@const index = Number(key.slice(1)) - 1}
          {formatNumber(row.months[index]?.numerator, 0)}
        {:else if key.startsWith('s')}
          {@const index = Number(key.slice(1)) - 1}
          <StatusBadge kind="display" value={row.months[index]?.status ?? 'NO_DATA'} />
        {/if}
      {/snippet}
    </DataTable>
    <div class="border-t-2 border-black bg-card px-4 py-2 font-mono text-[10px] text-muted-foreground">
      Status bulanan memakai perbandingan jumlah layanan bulan itu dengan target bulanan. Bulan tanpa data ditandai NO_DATA.
    </div>
  {/if}
</Card>
