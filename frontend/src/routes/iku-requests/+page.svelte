<script lang="ts">
  import { api } from '$lib/api/client';
  import { bootstrap, currentActor } from '$lib/stores/session';
  import PageHeader from '$lib/components/ui/page-header.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Badge from '$lib/components/ui/badge.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Select from '$lib/components/ui/select.svelte';
  import DataTable from '$lib/components/ui/data-table.svelte';
  import Spinner from '$lib/components/ui/spinner.svelte';
  import type { IKUChangeRequest } from '$lib/types';
  import {
    AGGREGATION_OPTIONS,
    IKU_ACTION_LABELS,
    IKU_STATUS_LABELS
  } from '$lib/constants/enums';
  import { Plus, RefreshCw, Check, X } from '@lucide/svelte';

  let requests = $state<IKUChangeRequest[]>([]);
  let loading = $state(true);
  let saving = $state(false);
  let error = $state('');
  let message = $state('');

  let action = $state<'CREATE' | 'UPDATE' | 'DEACTIVATE'>('CREATE');
  let selectedLabelID = $state('');
  let form = $state({
    kode: '', nama: '', categoryID: '', parentID: '', aggregation_method: 'RATIO_OF_SUMS',
    allow_direct_input: true, is_rollup_only: false, pjBidangID: ''
  });

  let labels = $derived($bootstrap?.labels ?? []);
  let categories = $derived(($bootstrap?.categories ?? []).filter((c) => c.aktif && c.jenis === 'LABEL'));
  let bidangList = $derived($bootstrap?.bidang ?? []);
  let labelOptions = $derived(labels.map((l) => ({ value: l.id, label: `${l.category_kode ?? l.tipe} · ${l.kode} · ${l.nama}` })));
  let categoryOptions = $derived(categories.map((c) => ({ value: c.id, label: `${c.kode} · ${c.nama}` })));
  let parentOptions = $derived([{ value: '', label: 'Tanpa induk' }, ...labelOptions]);
  let bidangOptions = $derived(bidangList.filter((b) => b.aktif && b.jenis !== 'timker').map((b) => ({ value: b.id, label: `${b.kode} · ${b.nama}` })));
  let canApply = $derived($currentActor?.role === 'pkm_operator');
  let isApprover = $derived(!!$currentActor && ['pkm_approver', 'admin'].includes($currentActor.role));

  let columns = $derived(isApprover ? [
    { key: 'action', label: 'Jenis', class: 'w-28' },
    { key: 'iku', label: 'IKU' },
    { key: 'faskes', label: 'Faskes', class: 'w-44' },
    { key: 'requester', label: 'Pengaju', class: 'w-48' },
    { key: 'status', label: 'Status', class: 'w-52' },
    { key: 'decision', label: 'Keputusan', class: 'w-64' }
  ] : [
    { key: 'action', label: 'Jenis', class: 'w-28' },
    { key: 'iku', label: 'IKU' },
    { key: 'faskes', label: 'Faskes', class: 'w-44' },
    { key: 'requester', label: 'Pengaju', class: 'w-48' },
    { key: 'status', label: 'Status', class: 'w-52' }
  ]);

  async function load() {
    loading = true;
    error = '';
    try {
      requests = await api.ikuChanges();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal memuat pengajuan IKU.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    void load();
  });

  function selectLabel(id: string) {
    selectedLabelID = id;
    const label = labels.find((l) => l.id === id);
    if (!label) return;
    form = {
      kode: label.kode, nama: label.nama, categoryID: label.category_id ?? '',
      parentID: label.parent_id ?? '', aggregation_method: label.aggregation_method,
      allow_direct_input: label.allow_direct_input, is_rollup_only: label.is_rollup_only,
      pjBidangID: label.pj_bidang_id ?? ''
    };
  }

  function resetForm() {
    action = 'CREATE';
    selectedLabelID = '';
    form = {
      kode: '', nama: '', categoryID: categoryOptions[0]?.value ?? '', parentID: '',
      aggregation_method: 'RATIO_OF_SUMS', allow_direct_input: true, is_rollup_only: false,
      pjBidangID: ''
    };
  }

  async function submitRequest() {
    if (!canApply) return;
    saving = true;
    error = '';
    try {
      await api.createIKUChange({
        action,
        label_id: action === 'CREATE' ? '' : selectedLabelID,
        kode: form.kode,
        nama: form.nama,
        category_id: form.categoryID,
        parent_id: form.parentID,
        aggregation_method: form.aggregation_method,
        allow_direct_input: form.allow_direct_input,
        is_rollup_only: form.is_rollup_only,
        pj_bidang_id: form.pjBidangID
      });
      message = 'Pengajuan IKU dikirim ke Puskesmas.';
      resetForm();
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal mengirim pengajuan IKU.';
    } finally {
      saving = false;
    }
  }

  async function decide(item: IKUChangeRequest, action: 'approve' | 'reject' | 'return') {
    const notes = action === 'approve' ? '' : prompt('Catatan wajib untuk menolak atau mengembalikan:')?.trim() ?? '';
    if (action !== 'approve' && !notes) {
      error = 'Catatan wajib diisi untuk menolak atau mengembalikan.';
      return;
    }
    error = '';
    try {
      await api.decideIKUChange(item.id, { action, notes });
      message = `Pengajuan ${ikuName(item)} diproses.`;
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal memproses pengajuan.';
    }
  }

  function canDecide(item: IKUChangeRequest): boolean {
    if (!isApprover) return false;
    if ($currentActor?.role === 'pkm_approver') return item.status === 'menunggu_puskesmas';
    if ($currentActor?.role === 'admin') return item.status === 'menunggu_admin';
    return false;
  }

  function ikuName(item: IKUChangeRequest): string {
    if (item.label_kode) return item.label_kode;
    if (item.applied_label_kode) return item.applied_label_kode;
    try {
      const payload = JSON.parse(item.payload_json) as { kode?: string; nama?: string };
      return payload.kode ?? payload.nama ?? item.payload_json;
    } catch {
      return item.payload_json;
    }
  }
</script>

<PageHeader
  title="Pengajuan IKU"
  description="Operator Puskesmas mengajukan IKU; Puskesmas memeriksa; Admin mempublikasikan."
>
  {#snippet actions()}
    <Button variant="outline" onclick={load}><RefreshCw size={14} /> Muat ulang</Button>
  {/snippet}
</PageHeader>

{#if error}<div class="mb-4 border-2 border-black bg-destructive p-3 font-mono text-xs font-bold">{error}</div>{/if}
{#if message}<div class="mb-4 border-2 border-black bg-success p-3 font-mono text-xs font-bold">{message}</div>{/if}

{#if canApply}
  <Card class="mb-4">
    <div class="border-b-2 border-black bg-secondary px-4 py-3">
      <div class="nb-kicker">Form pengajuan</div>
      <h2 class="mt-1 text-xl">Ajukan perubahan IKU</h2>
    </div>
    <div class="grid gap-3 p-4 sm:grid-cols-2">
      <div>
        <label class="nb-label" for="action">Jenis pengajuan</label>
        <Select id="action" bind:value={action} options={[
          { value: 'CREATE', label: 'Tambah IKU' },
          { value: 'UPDATE', label: 'Ubah IKU' },
          { value: 'DEACTIVATE', label: 'Nonaktifkan IKU' }
        ]} />
      </div>
      {#if action !== 'CREATE'}
        <div>
          <label class="nb-label" for="label">IKU yang diajukan</label>
          <Select id="label" bind:value={selectedLabelID} options={labelOptions} placeholder="Pilih IKU" />
          <Button size="sm" variant="ghost" class="mt-1" onclick={() => selectLabel(selectedLabelID)}>Isi dari IKU terpilih</Button>
        </div>
      {/if}
      {#if action !== 'DEACTIVATE'}
        <div><label class="nb-label" for="kode">Kode IKU</label><Input id="kode" bind:value={form.kode} placeholder="IKU-2026-002" /></div>
        <div><label class="nb-label" for="nama">Nama IKU</label><Input id="nama" bind:value={form.nama} /></div>
        <div><label class="nb-label" for="category">Kategori</label><Select id="category" bind:value={form.categoryID} options={categoryOptions} placeholder="Pilih kategori" /></div>
        <div><label class="nb-label" for="parent">IKU induk</label><Select id="parent" bind:value={form.parentID} options={parentOptions} /></div>
        <div><label class="nb-label" for="agg">Cara hitung</label><Select id="agg" bind:value={form.aggregation_method} options={AGGREGATION_OPTIONS} /></div>
        <div><label class="nb-label" for="pj">PJ Bidang</label><Select id="pj" bind:value={form.pjBidangID} options={bidangOptions} placeholder="Pilih bidang" /></div>
        <label class="flex items-center gap-2 font-mono text-[10px] font-bold uppercase"><input type="checkbox" bind:checked={form.allow_direct_input} /> Input langsung</label>
        <label class="flex items-center gap-2 font-mono text-[10px] font-bold uppercase"><input type="checkbox" bind:checked={form.is_rollup_only} /> Nilai dari komponen</label>
      {/if}
    </div>
    <div class="flex justify-end border-t-2 border-black p-4">
      <Button variant="success" disabled={saving} onclick={submitRequest}><Plus size={14} /> Kirim pengajuan</Button>
    </div>
  </Card>
{/if}

<Card>
  <div class="border-b-2 border-black bg-muted px-4 py-3">
    <div class="nb-kicker">Daftar pengajuan</div>
    <h2 class="mt-1 text-xl">Antrean persetujuan IKU</h2>
  </div>
  {#if loading}
    <div class="flex items-center justify-center py-16"><Spinner label="Memuat pengajuan…" /></div>
  {:else}
    <DataTable columns={columns} rows={requests}>
      {#snippet cell(item: IKUChangeRequest, key)}
        {#if key === 'action'}<Badge tone="neutral">{IKU_ACTION_LABELS[item.action] ?? item.action}</Badge>
        {:else if key === 'iku'}<span class="font-semibold">{ikuName(item)}</span>
        {:else if key === 'faskes'}{item.faskes_nama ?? item.faskes_id}
        {:else if key === 'requester'}{item.requested_by_nama ?? item.requested_by}
        {:else if key === 'status'}<Badge tone="info">{IKU_STATUS_LABELS[item.status] ?? item.status}</Badge>
        {:else if key === 'decision'}
          {#if canDecide(item)}
            <div class="flex flex-wrap gap-1">
              <Button size="sm" variant="success" onclick={() => decide(item, 'approve')}><Check size={12} /> Setujui</Button>
              <Button size="sm" variant="outline" onclick={() => decide(item, 'return')}><X size={12} /> Kembalikan</Button>
              <Button size="sm" variant="ghost" onclick={() => decide(item, 'reject')}><X size={12} /> Tolak</Button>
            </div>
          {:else}
            <span class="text-muted-foreground">—</span>
          {/if}
        {/if}
      {/snippet}
    </DataTable>
  {/if}
</Card>
