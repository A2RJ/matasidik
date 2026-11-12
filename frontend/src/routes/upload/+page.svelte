<script lang="ts">
  import { goto } from '$app/navigation';
  import { api, ApiError } from '$lib/api/client';
  import { bootstrap, currentActor, periods, selectedPeriod } from '$lib/stores/session';
  import PageHeader from '$lib/components/ui/page-header.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Badge from '$lib/components/ui/badge.svelte';
  import Select from '$lib/components/ui/select.svelte';
  import EmptyState from '$lib/components/ui/empty-state.svelte';
  import Spinner from '$lib/components/ui/spinner.svelte';
  import ProgressSteps from '$lib/components/ui/progress-steps.svelte';
  import StatusBadge from '$lib/components/ui/status-badge.svelte';
  import PreviewRow from '$lib/components/upload/preview-row.svelte';
  import { PAGE_COPY, EMPTY_STATE } from '$lib/constants/labels';
  import { formatNumber, parseProgress } from '$lib/utils';
  import type { Job, SubmissionDetail, SubmissionRow } from '$lib/types';
  import { UploadCloud, FileSpreadsheet, CheckCircle2, ArrowRight } from '@lucide/svelte';

  let file = $state<File | null>(null);
  let faskesId = $state('');
  let uploading = $state(false);
  let error = $state('');
  let job = $state<Job | null>(null);
  let detail = $state<SubmissionDetail | null>(null);
  let submitting = $state(false);
  let rowSelections = $state<Record<string, string>>({});
  let savingRow = $state<string | null>(null);
  let message = $state('');

  let periodOptions = $derived($periods.map((p) => ({ value: p.id, label: p.kode })));
  let faskesOptions = $derived(($bootstrap?.faskes ?? []).map((f) => ({ value: f.id, label: f.nama })));
  let progress = $derived(parseProgress(job?.progress_json));
  let allowedRole = $derived(
    $currentActor && ['pkm_operator', 'pkm_pj', 'admin'].includes($currentActor.role)
  );
  let blockers = $derived(detail?.blockers ?? 0);
  let matchedCount = $derived(detail?.rows.filter((r) => r.classification_status === 'MATCHED').length ?? 0);



  $effect(() => {
    if (!faskesId && $currentActor?.faskes_id) faskesId = $currentActor.faskes_id;
  });

  function onFileChange(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    file = input.files?.[0] ?? null;
  }

  async function pollJob(jobId: string) {
    for (let i = 0; i < 120; i++) {
      const result = await api.job(jobId);
      job = result;
      if (result.status === 'done') return result;
      if (result.status === 'failed' || result.status === 'dead') {
        throw new Error(result.last_error ?? 'Job gagal diproses.');
      }
      await new Promise((resolve) => setTimeout(resolve, 1000));
    }
    throw new Error('Job belum selesai setelah 2 menit.');
  }

  async function loadDetail(id: string) {
    detail = await api.submission(id);
    const next: Record<string, string> = {};
    for (const row of detail.rows) next[row.id] = row.iku_kode ?? '';
    rowSelections = next;
  }

  async function upload() {
    if (!file) {
      error = 'Pilih file BNBA terlebih dahulu.';
      return;
    }
    const period = $periods.find((p) => p.id === $selectedPeriod);
    if (!period) {
      error = 'Pilih periode.';
      return;
    }
    if (!faskesId) {
      error = 'Pilih faskes.';
      return;
    }
    uploading = true;
    error = '';
    message = '';
    detail = null;
    job = null;
    try {
      const form = new FormData();
      form.set('file', file);
      form.set('periode', period.kode);
      form.set('faskes_id', faskesId);
      const result = await api.upload(form, crypto.randomUUID());
      if (result.job) {
        await pollJob(result.job.id);
      }
      await loadDetail(result.submission.id);
      message = `File tersimpan sebagai ${result.submission.kode}. Periksa preview sebelum kirim.`;
    } catch (err) {
      if (err instanceof ApiError) error = err.message;
      else error = err instanceof Error ? err.message : 'Upload gagal.';
    } finally {
      uploading = false;
    }
  }

  async function saveRow(row: SubmissionRow) {
    if (!detail) return;
    const selected = rowSelections[row.id];
    if (!selected) {
      error = 'Pilih IKU leaf terlebih dahulu.';
      return;
    }
    savingRow = row.id;
    error = '';
    try {
      await api.updateRow(detail.submission.id, row.id, selected, 'Koreksi IKU di preview');
      await loadDetail(detail.submission.id);
      message = `Baris R${row.row_no} diperbarui.`;
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal menyimpan baris.';
    } finally {
      savingRow = null;
    }
  }

  async function submit() {
    if (!detail) return;
    submitting = true;
    error = '';
    try {
      await api.submit(detail.submission.id);
      await goto(`/submissions/${detail.submission.id}`);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal mengirim submission.';
    } finally {
      submitting = false;
    }
  }
</script>

<PageHeader title={PAGE_COPY.upload.title} description={PAGE_COPY.upload.description} />

{#if !allowedRole}
  <div class="mb-4 border-2 border-black bg-warning p-3 font-mono text-xs font-bold">
    Aktor saat ini tidak punya peran upload. Ganti ke Siti (Puskesmas) atau Admin di kanan atas.
  </div>
{/if}

<div class="grid gap-4 xl:grid-cols-[420px_1fr]">
  <Card>
    <div class="border-b-2 border-black bg-secondary px-4 py-3">
      <div class="nb-kicker">Langkah 1</div>
      <h2 class="mt-1 text-xl">Pilih file & periode</h2>
    </div>
    <CardContent class="space-y-4">
      <div>
        <label class="nb-label" for="file">File BNBA (.csv / .xlsx)</label>
        <input
          id="file"
          type="file"
          accept=".csv,.xlsx"
          onchange={onFileChange}
          class="block w-full border-2 border-black bg-white p-3 font-mono text-xs file:mr-3 file:border-2 file:border-black file:bg-primary file:px-3 file:py-1 file:font-mono file:text-[10px] file:font-bold file:uppercase"
        />
        {#if file}
          <div class="mt-2 flex items-center gap-2 border-2 border-black bg-muted px-2 py-1 font-mono text-[10px]">
            <FileSpreadsheet size={14} /> {file.name} · {formatNumber(file.size / 1024, 1)} KB
          </div>
        {/if}
      </div>
      <div>
        <label class="nb-label" for="period">Periode</label>
        <Select id="period" bind:value={$selectedPeriod} options={periodOptions} placeholder="Pilih periode" />
      </div>
      <div>
        <label class="nb-label" for="faskes">Faskes pengupload</label>
        <Select id="faskes" bind:value={faskesId} options={faskesOptions} placeholder="Pilih faskes" />
      </div>
      <Button class="w-full" size="lg" onclick={upload} disabled={uploading || !allowedRole}>
        {#if uploading}<Spinner label="Mengunggah…" />{:else}<UploadCloud size={16} /> Upload & proses{/if}
      </Button>
      <a href="/api/v1/master/template" class="block text-center font-mono text-[10px] font-bold underline">
        Unduh template CSV
      </a>
      {#if message}<div class="border-2 border-black bg-success p-2 font-mono text-[10px] font-bold">{message}</div>{/if}
      {#if error}<div class="border-2 border-black bg-destructive p-2 font-mono text-[10px] font-bold">{error}</div>{/if}
    </CardContent>
  </Card>

  <div class="space-y-4">
    <Card>
      <div class="border-b-2 border-black bg-primary px-4 py-3">
        <div class="nb-kicker">Langkah 2</div>
        <h2 class="mt-1 text-xl">Proses background</h2>
      </div>
      <CardContent>
        {#if job}
          <ProgressSteps
            steps={[
              { key: 'validasi', label: 'Validasi file' },
              { key: 'klasifikasi', label: 'Klasifikasi' },
              { key: 'simpan', label: 'Simpan baris' },
              { key: 'selesai', label: 'Selesai' }
            ]}
            current={String(progress.step ?? 'queued')}
            percent={Number(progress.percent ?? 0)}
          />
          <div class="mt-3 flex flex-wrap items-center gap-2 font-mono text-[10px]">
            <Badge tone={job.status === 'done' ? 'success' : job.status === 'failed' ? 'danger' : 'info'}>
              Job {job.status}
            </Badge>
            <span>{String(progress.rows ?? 0)} baris · {String(progress.matched ?? 0)} matched</span>
          </div>
        {:else}
          <div class="py-6 text-center text-sm text-muted-foreground">
            Progress akan muncul setelah file diunggah. UI polling setiap ±1 detik.
          </div>
        {/if}
      </CardContent>
    </Card>

    {#if detail}
      <Card>
        <div class="flex flex-wrap items-center justify-between gap-3 border-b-2 border-black bg-muted px-4 py-3">
          <div>
            <div class="nb-kicker">Langkah 3 · Preview</div>
            <h2 class="mt-1 text-xl">{detail.submission.kode}</h2>
            <div class="mt-1 font-mono text-[10px] text-muted-foreground">
              {detail.submission.faskes_nama} · {detail.submission.periode_kode} · {detail.submission.uploader_nama}
            </div>
          </div>
          <div class="flex flex-wrap gap-2">
            <Badge tone="success">{matchedCount} matched</Badge>
            <Badge tone={blockers > 0 ? 'warning' : 'success'}>{blockers} blocker</Badge>
            <StatusBadge kind="submission" value={detail.submission.status} />
          </div>
        </div>

        {#if detail.rows.length === 0}
          <div class="p-4"><EmptyState title="Belum ada baris" description="File mungkin masih diproses atau gagal diparse." /></div>
        {:else}
          <div>
            {#each detail.rows as row (row.id)}
              <PreviewRow
                {row}
                labels={detail.allowed_iku[row.id] ?? []}
                bind:value={rowSelections[row.id]}
                saving={savingRow === row.id}
                onsave={() => saveRow(row)}
              />
            {/each}
          </div>
          <div class="flex flex-wrap items-center justify-between gap-3 border-t-2 border-black bg-card px-4 py-3">
            <div class="font-mono text-[10px] text-muted-foreground">
              Tombol Kirim aktif hanya jika tidak ada blocker. Perubahan IKU tercatat di audit.
            </div>
            <Button variant="success" disabled={submitting || blockers > 0} onclick={submit}>
              {#if submitting}<Spinner label="Mengirim…" />{:else}<CheckCircle2 size={16} /> Kirim ke workflow <ArrowRight size={14} />{/if}
            </Button>
          </div>
        {/if}
      </Card>
    {/if}
  </div>
</div>
