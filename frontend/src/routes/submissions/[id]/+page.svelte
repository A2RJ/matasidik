<script lang="ts">
  import { page } from '$app/stores';
  import { api } from '$lib/api/client';
  import { currentActor } from '$lib/stores/session';
  import PageHeader from '$lib/components/ui/page-header.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Badge from '$lib/components/ui/badge.svelte';
  import Modal from '$lib/components/ui/modal.svelte';
  import DataTable from '$lib/components/ui/data-table.svelte';
  import Spinner from '$lib/components/ui/spinner.svelte';
  import StatusBadge from '$lib/components/ui/status-badge.svelte';
  import ProgressSteps from '$lib/components/ui/progress-steps.svelte';
  import { formatDateTime, formatNumber, formatPercent } from '$lib/utils';
  import {
    REVIEW_ACTION_LABELS,
    SCOPE_TYPE_LABELS,
    WORKFLOW_STEP_LABELS
  } from '$lib/constants/enums';
  import type { SubmissionDetail } from '$lib/types';
  import { ArrowLeft, CheckCircle2, RotateCcw, Download, Send, ShieldCheck } from '@lucide/svelte';

  let detail = $state<SubmissionDetail | null>(null);
  let loading = $state(true);
  let error = $state('');
  let message = $state('');
  let decisionOpen = $state(false);
  let decisionAction = $state<'approve' | 'return'>('approve');
  let notes = $state('');
  let processing = $state(false);

  let submission = $derived(detail?.submission);
  let actor = $derived($currentActor);
  let canSubmit = $derived(
    !!submission &&
      ['draf', 'dikembalikan_petugas', 'dikembalikan_verifikator'].includes(submission.status) &&
      !!actor &&
      (actor.id === submission.uploader_id || actor.role === 'admin')
  );
  let canApprovePKM = $derived(
    !!submission &&
      submission.status === 'dikirim' &&
      submission.step === 'APPROVE_PKM' &&
      !!actor &&
      actor.role === 'pkm_approver'
  );
  let canVerify = $derived(
    !!submission &&
      submission.status === 'dikirim' &&
      submission.step === 'VERIFY_DINAS' &&
      !!actor &&
      ['dinas_verifier', 'admin'].includes(actor.role)
  );
  let canValidate = $derived(
    !!submission &&
      submission.status === 'diverifikasi' &&
      !!actor &&
      ['kabid', 'sekretaris', 'admin'].includes(actor.role)
  );
  let step = $derived(submission?.step ?? 'DRAFT_PKM');
  let columns = [
    { key: 'row', label: 'Baris', class: 'w-16' },
    { key: 'indicator', label: 'Indikator' },
    { key: 'value', label: 'Nilai', class: 'text-right w-28' },
    { key: 'status', label: 'Klasifikasi', class: 'w-40' },
    { key: 'iku', label: 'IKU' },
    { key: 'error', label: 'Catatan' }
  ];

  async function load() {
    const id = $page.params.id;
    if (!id) return;
    loading = true;
    error = '';
    try {
      detail = await api.submission(id);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal memuat submission.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if ($page.params.id) void load();
  });

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

  async function submit() {
    if (!detail) return;
    processing = true;
    error = '';
    try {
      await api.submit(detail.submission.id);
      message = 'Submission dikirim untuk persetujuan Puskesmas.';
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal mengirim.';
    } finally {
      processing = false;
    }
  }

  function openDecision(action: 'approve' | 'return') {
    decisionAction = action;
    notes = '';
    decisionOpen = true;
  }

  async function decide() {
    if (!detail) return;
    processing = true;
    error = '';
    try {
      await api.decision(detail.submission.id, decisionAction, notes, decisionAction === 'return' ? 'REVISI' : 'OK');
      message = decisionAction === 'approve' ? 'Keputusan disimpan.' : 'Dikembalikan dengan catatan.';
      decisionOpen = false;
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal menyimpan keputusan.';
    } finally {
      processing = false;
    }
  }
</script>

<PageHeader
  kicker="Alur / submission"
  title={submission?.kode ?? 'Detail submission'}
  description={submission ? `${submission.faskes_nama} · ${submission.periode_kode} · ${submission.uploader_nama}` : 'Memuat…'}
>
  {#snippet actions()}
    <a href="/submissions" class="inline-flex h-10 items-center gap-2 border-2 border-black bg-card px-4 font-mono text-xs font-bold uppercase shadow-nbsm">
      <ArrowLeft size={14} /> Kembali
    </a>
    {#if canSubmit}
      <Button variant="success" disabled={processing || (detail?.blockers ?? 0) > 0} onclick={submit}>
        <Send size={14} /> Kirim
      </Button>
    {/if}
    {#if canApprovePKM || canVerify || canValidate}
      <Button variant="success" onclick={() => openDecision('approve')}><CheckCircle2 size={14} /> Setujui</Button>
      <Button variant="destructive" onclick={() => openDecision('return')}><RotateCcw size={14} /> Kembalikan</Button>
    {/if}
  {/snippet}
</PageHeader>

{#if error}<div class="mb-4 border-2 border-black bg-destructive p-3 font-mono text-xs font-bold">{error}</div>{/if}
{#if message}<div class="mb-4 border-2 border-black bg-success p-3 font-mono text-xs font-bold">{message}</div>{/if}

{#if loading}
  <div class="flex items-center justify-center py-20"><Spinner label="Memuat submission…" /></div>
{:else if detail}
  <div class="mb-4 grid gap-3 lg:grid-cols-[1fr_340px]">
    <Card>
      <div class="border-b-2 border-black bg-muted px-4 py-3">
        <div class="nb-kicker">Status alur</div>
        <div class="mt-1 flex flex-wrap items-center gap-2">
          <StatusBadge kind="submission" value={detail.submission.status} />
          <Badge tone="neutral">{WORKFLOW_STEP_LABELS[detail.submission.step] ?? detail.submission.step}</Badge>
          <Badge tone={(detail.blockers ?? 0) > 0 ? 'warning' : 'success'}>{detail.blockers ?? 0} baris bermasalah</Badge>
        </div>
      </div>
      <div class="p-4">
        <ProgressSteps
          steps={[
            { key: 'DRAFT_PKM', label: 'Draf PKM' },
            { key: 'APPROVE_PKM', label: 'Approval PKM' },
            { key: 'VERIFY_DINAS', label: 'Verifikasi' },
            { key: 'VALIDATE_KABID', label: 'Validasi' },
            { key: 'FINAL', label: 'Final' }
          ]}
          current={step}
          percent={step === 'DRAFT_PKM' ? 20 : step === 'APPROVE_PKM' ? 40 : step === 'VERIFY_DINAS' ? 60 : step === 'VALIDATE_KABID' ? 80 : 100}
        />
        <div class="mt-4 grid gap-2 font-mono text-[10px] text-muted-foreground sm:grid-cols-2">
          <div>Berkas: BNBA</div>
          <div>Dibuat: {formatDateTime(detail.submission.created_at)}</div>
          <div>Diperbarui: {formatDateTime(detail.submission.updated_at)}</div>
        </div>
      </div>
    </Card>

    <Card>
      <div class="border-b-2 border-black bg-secondary px-4 py-3">
        <div class="nb-kicker">Aksi</div>
        <h2 class="mt-1 text-lg">Keputusan</h2>
      </div>
      <div class="space-y-3 p-4">
        {#if canSubmit}
          <Button class="w-full" variant="success" disabled={processing || (detail.blockers ?? 0) > 0} onclick={submit}>
            <Send size={14} /> Kirim untuk persetujuan
          </Button>
        {/if}
        {#if canApprovePKM || canVerify || canValidate}
          <Button class="w-full" variant="success" onclick={() => openDecision('approve')}>
            <ShieldCheck size={14} /> Setujui
          </Button>
          <Button class="w-full" variant="destructive" onclick={() => openDecision('return')}>
            <RotateCcw size={14} /> Kembalikan + catatan
          </Button>
        {/if}
        {#if !canSubmit && !canApprovePKM && !canVerify && !canValidate}
          <div class="border-2 border-dashed border-black/40 bg-muted p-3 font-mono text-[10px]">
            Tidak ada tindakan untuk akun Anda pada tahap ini.
          </div>
        {/if}
        <Button variant="outline" class="w-full" onclick={() => downloadFile(detail!.submission.file_id)}>
          <Download size={14} /> Unduh berkas asli
        </Button>
      </div>
    </Card>
  </div>

  <Card class="mb-4">
    <div class="border-b-2 border-black bg-muted px-4 py-3">
      <div class="nb-kicker">Rincian baris</div>
      <h2 class="mt-1 text-xl">Baris data yang terbaca</h2>
    </div>
    <DataTable {columns} rows={detail.rows}>
      {#snippet cell(row, key)}
        {#if key === 'row'}
          <span class="font-bold">R{row.row_no}</span>
        {:else if key === 'indicator'}
          <div class="font-semibold">{row.indikator_nama ?? row.kode_indikator_raw ?? '—'}</div>
          <div class="mt-1 font-mono text-[10px] text-muted-foreground">{row.kode_iku_raw ?? '—'} · {SCOPE_TYPE_LABELS[row.scope_type] ?? row.scope_type} {row.scope_code ?? '—'}</div>
        {:else if key === 'value'}
          <span class="font-bold">{formatPercent(row.nilai)}</span>
          <div class="text-[10px] text-muted-foreground">{formatNumber(row.numerator)}/{formatNumber(row.denominator)}</div>
        {:else if key === 'status'}
          <StatusBadge kind="classification" value={row.classification_status} />
        {:else if key === 'iku'}
          <div class="font-mono text-xs">{row.iku_kode ?? '—'}</div>
          <div class="text-[10px] text-muted-foreground">{row.iku_nama ?? 'Belum terklasifikasi'}</div>
        {:else if key === 'error'}
          <span class={row.classification_status === 'MATCHED' ? 'text-[10px] text-muted-foreground' : 'text-[10px] text-destructive'}>{row.error_message ?? '—'}</span>
        {/if}
      {/snippet}
    </DataTable>
  </Card>

  <div class="grid gap-4 lg:grid-cols-2">
    <Card>
      <div class="border-b-2 border-black bg-secondary px-4 py-3">
        <div class="nb-kicker">Timeline</div>
        <h2 class="mt-1 text-xl">Keputusan & catatan</h2>
      </div>
      <div class="space-y-3 p-4">
        {#if !detail.reviews || detail.reviews.length === 0}
          <div class="text-sm text-muted-foreground">Belum ada review.</div>
        {:else}
          {#each detail.reviews as review (review.id)}
            <div class="border-2 border-black bg-card p-3">
              <div class="flex flex-wrap items-center justify-between gap-2">
                <div class="font-mono text-[10px] font-bold uppercase">{WORKFLOW_STEP_LABELS[review.step] ?? review.step} · {REVIEW_ACTION_LABELS[review.action] ?? review.action}</div>
                <div class="font-mono text-[10px] text-muted-foreground">{formatDateTime(review.created_at)}</div>
              </div>
              <div class="mt-1 text-sm font-semibold">{review.reviewer_nama ?? review.reviewer_id ?? 'Sistem'}</div>
              {#if review.notes}<p class="mt-1 text-xs text-muted-foreground">{review.notes}</p>{/if}
            </div>
          {/each}
        {/if}
      </div>
    </Card>

    <Card>
      <div class="border-b-2 border-black bg-muted px-4 py-3">
        <div class="nb-kicker">Revisi</div>
        <h2 class="mt-1 text-xl">Riwayat berkas</h2>
      </div>
      <div class="space-y-2 p-4">
        {#each detail.revisions as revision (revision.id)}
          <div class="flex items-center justify-between border-2 border-black bg-card px-3 py-2 font-mono text-xs">
            <span>Revisi {revision.revision_no}</span>
            <span class="text-muted-foreground">{formatDateTime(revision.created_at)}</span>
          </div>
        {/each}
      </div>
    </Card>
  </div>
{/if}

<Modal
  open={decisionOpen}
  title={decisionAction === 'approve' ? 'Setujui submission' : 'Kembalikan submission'}
  description={decisionAction === 'approve' ? 'Final wajib catatan penerimaan.' : 'Pengembalian wajib catatan perbaikan.'}
  onclose={() => (decisionOpen = false)}
>
  <div class="space-y-3">
    <label class="nb-label" for="notes">Catatan</label>
    <textarea id="notes" bind:value={notes} class="nb-textarea" rows="4" placeholder="Tulis catatan keputusan…"></textarea>
    <div class="flex justify-end gap-2">
      <Button variant="outline" onclick={() => (decisionOpen = false)}>Batal</Button>
      <Button variant={decisionAction === 'approve' ? 'success' : 'destructive'} disabled={processing} onclick={decide}>
        {decisionAction === 'approve' ? 'Setujui' : 'Kembalikan'}
      </Button>
    </div>
  </div>
</Modal>
