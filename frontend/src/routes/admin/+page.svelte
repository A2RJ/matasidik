<script lang="ts">
  import { api } from '$lib/api/client';
  import PageHeader from '$lib/components/ui/page-header.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Badge from '$lib/components/ui/badge.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Select from '$lib/components/ui/select.svelte';
  import Modal from '$lib/components/ui/modal.svelte';
  import DataTable from '$lib/components/ui/data-table.svelte';
  import Spinner from '$lib/components/ui/spinner.svelte';
  import { ROLE_LABELS } from '$lib/constants/enums';
  import type { AdminOverview, Faskes, Position, User } from '$lib/types';
  import { Plus, RefreshCw, KeyRound } from '@lucide/svelte';

  let overview = $state<AdminOverview | null>(null);
  let loading = $state(true);
  let error = $state('');
  let message = $state('');
  let tab = $state<'users' | 'faskes' | 'positions'>('users');

  let userOpen = $state(false);
  let editingUser = $state<User | null>(null);
  let userForm = $state({ username: '', nama: '', password: '', position_id: '', faskes_id: '', bidang_id: '' });

  let faskesOpen = $state(false);
  let editingFaskes = $state<Faskes | null>(null);
  let faskesForm = $state({ kode: '', nama: '', wilayah_id: '', jenis: 'puskesmas' });

  let positionOpen = $state(false);
  let editingPosition = $state<Position | null>(null);
  let positionForm = $state({ kode: '', nama: '', role: 'pkm_operator', level: '1', bidang_id: '', faskes_id: '' });

  let roleOptions = $derived(Object.entries(ROLE_LABELS).map(([value, label]) => ({ value, label })));
  let faskesOptions = $derived((overview?.faskes ?? []).map((f) => ({ value: f.id, label: f.nama })));
  let bidangOptions = $derived((overview?.bidang ?? []).map((b) => ({ value: b.id, label: b.nama })));
  let positionOptions = $derived((overview?.positions ?? []).map((p) => ({ value: p.id, label: `${p.nama} · ${ROLE_LABELS[p.role]}` })));

  let userColumns = [
    { key: 'nama', label: 'Nama / username' },
    { key: 'role', label: 'Role efektif', width: '160px' },
    { key: 'position', label: 'Posisi', width: '220px' },
    { key: 'faskes', label: 'Faskes', width: '180px' },
    { key: 'bidang', label: 'Bidang', width: '180px' },
    { key: 'aktif', label: 'Status', width: '110px' },
    { key: 'action', label: 'Aksi', width: '90px' }
  ];
  let faskesColumns = [
    { key: 'kode', label: 'Kode', width: '160px' },
    { key: 'nama', label: 'Nama faskes' },
    { key: 'jenis', label: 'Jenis', width: '120px' },
    { key: 'aktif', label: 'Status', width: '110px' },
    { key: 'action', label: 'Aksi', width: '90px' }
  ];
  let positionColumns = [
    { key: 'kode', label: 'Kode', width: '170px' },
    { key: 'nama', label: 'Nama jabatan' },
    { key: 'role', label: 'Role', width: '160px' },
    { key: 'level', label: 'Level', width: '80px' },
    { key: 'scope', label: 'Scope', width: '180px' },
    { key: 'action', label: 'Aksi', width: '90px' }
  ];

  async function load() {
    loading = true;
    error = '';
    try {
      overview = await api.adminOverview();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal memuat administrasi.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    void load();
  });

  function openUser(user?: User) {
    editingUser = user ?? null;
    userForm = user
      ? { username: user.username, nama: user.nama, password: '', position_id: user.position_id ?? '', faskes_id: user.faskes_id ?? '', bidang_id: user.bidang_id ?? '' }
      : { username: '', nama: '', password: '', position_id: '', faskes_id: '', bidang_id: '' };
    userOpen = true;
  }

  async function saveUser() {
    error = '';
    try {
      if (editingUser) {
        await api.updateUser(editingUser.id, {
          username: userForm.username, nama: userForm.nama, position_id: userForm.position_id,
          faskes_id: userForm.faskes_id, bidang_id: userForm.bidang_id
        });
      } else {
        await api.createUser(userForm);
      }
      userOpen = false;
      message = 'User disimpan.';
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal menyimpan user.';
    }
  }

  async function resetPassword(user: User) {
    const password = prompt(`Password baru untuk ${user.username} (minimal 8 karakter):`);
    if (!password) return;
    try {
      await api.resetPassword(user.id, password);
      message = `Password ${user.username} diperbarui.`;
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal reset password.';
    }
  }

  function openFaskes(f?: Faskes) {
    editingFaskes = f ?? null;
    faskesForm = f ? { kode: f.kode, nama: f.nama, wilayah_id: f.wilayah_id ?? '', jenis: f.jenis } : { kode: '', nama: '', wilayah_id: '', jenis: 'puskesmas' };
    faskesOpen = true;
  }

  async function saveFaskes() {
    error = '';
    try {
      if (editingFaskes) await api.updateFaskes(editingFaskes.id, faskesForm);
      else await api.createFaskes(faskesForm);
      faskesOpen = false;
      message = 'Faskes disimpan.';
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal menyimpan faskes.';
    }
  }

  function openPosition(p?: Position) {
    editingPosition = p ?? null;
    positionForm = p
      ? { kode: p.kode, nama: p.nama, role: p.role, level: String(p.level), bidang_id: p.bidang_id ?? '', faskes_id: p.faskes_id ?? '' }
      : { kode: '', nama: '', role: 'pkm_operator', level: '1', bidang_id: '', faskes_id: '' };
    positionOpen = true;
  }

  async function savePosition() {
    error = '';
    const payload = { ...positionForm, level: Number(positionForm.level) };
    try {
      if (editingPosition) await api.updatePosition(editingPosition.id, payload);
      else await api.createPosition(payload);
      positionOpen = false;
      message = 'Posisi disimpan.';
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal menyimpan posisi.';
    }
  }

  let faskesMap = $derived(new Map((overview?.faskes ?? []).map((f) => [f.id, f])));
  let bidangMap = $derived(new Map((overview?.bidang ?? []).map((b) => [b.id, b])));
</script>

<PageHeader title="Administrasi" description="Puskesmas, user, dan jabatan struktural dikelola terpisah. Perubahan posisi langsung mengubah role efektif dan routing approval.">
  {#snippet actions()}
    <Button variant="outline" onclick={load}><RefreshCw size={14} /> Muat ulang</Button>
    {#if tab === 'users'}<Button onclick={() => openUser()}><Plus size={14} /> User</Button>{/if}
    {#if tab === 'faskes'}<Button onclick={() => openFaskes()}><Plus size={14} /> Faskes</Button>{/if}
    {#if tab === 'positions'}<Button onclick={() => openPosition()}><Plus size={14} /> Posisi</Button>{/if}
  {/snippet}
</PageHeader>

{#if error}<div class="mb-4 border-2 border-black bg-destructive p-3 font-mono text-xs font-bold">{error}</div>{/if}
{#if message}<div class="mb-4 border-2 border-black bg-success p-3 font-mono text-xs font-bold">{message}</div>{/if}

<div class="mb-4 flex flex-wrap gap-2">
  <Button variant={tab === 'users' ? 'default' : 'outline'} onclick={() => (tab = 'users')}>User ({overview?.users.length ?? 0})</Button>
  <Button variant={tab === 'faskes' ? 'default' : 'outline'} onclick={() => (tab = 'faskes')}>Puskesmas ({overview?.faskes.length ?? 0})</Button>
  <Button variant={tab === 'positions' ? 'default' : 'outline'} onclick={() => (tab = 'positions')}>Jabatan ({overview?.positions.length ?? 0})</Button>
</div>

<Card>
  {#if loading}
    <div class="flex items-center justify-center py-16"><Spinner label="Memuat administrasi…" /></div>
  {:else if tab === 'users'}
    <DataTable columns={userColumns} rows={overview?.users ?? []}>
      {#snippet cell(item: User, key)}
        {#if key === 'nama'}<div class="font-bold">{item.nama}</div><div class="text-[10px] text-muted-foreground">@{item.username}</div>
        {:else if key === 'role'}<Badge tone="info">{ROLE_LABELS[item.role]}</Badge>
        {:else if key === 'position'}{item.position_nama ?? '—'}
        {:else if key === 'faskes'}{item.faskes_id ? faskesMap.get(item.faskes_id)?.nama ?? item.faskes_id : '—'}
        {:else if key === 'bidang'}{item.bidang_id ? bidangMap.get(item.bidang_id)?.nama ?? item.bidang_id : '—'}
        {:else if key === 'aktif'}<Badge tone={item.aktif ? 'success' : 'neutral'}>{item.aktif ? 'Aktif' : 'Nonaktif'}</Badge>
        {:else if key === 'action'}<div class="flex gap-1"><Button size="sm" variant="outline" onclick={() => openUser(item)}>Edit</Button><Button size="sm" variant="ghost" onclick={() => resetPassword(item)}><KeyRound size={12} /></Button></div>
        {/if}
      {/snippet}
    </DataTable>
  {:else if tab === 'faskes'}
    <DataTable columns={faskesColumns} rows={overview?.faskes ?? []}>
      {#snippet cell(item: Faskes, key)}
        {#if key === 'kode'}<span class="font-bold">{item.kode}</span>
        {:else if key === 'nama'}{item.nama}
        {:else if key === 'jenis'}{item.jenis}
        {:else if key === 'aktif'}<Badge tone={item.aktif ? 'success' : 'neutral'}>{item.aktif ? 'Aktif' : 'Nonaktif'}</Badge>
        {:else if key === 'action'}<Button size="sm" variant="outline" onclick={() => openFaskes(item)}>Edit</Button>
        {/if}
      {/snippet}
    </DataTable>
  {:else}
    <DataTable columns={positionColumns} rows={overview?.positions ?? []}>
      {#snippet cell(item: Position, key)}
        {#if key === 'kode'}<span class="font-bold">{item.kode}</span>
        {:else if key === 'nama'}{item.nama}
        {:else if key === 'role'}<Badge tone="info">{ROLE_LABELS[item.role]}</Badge>
        {:else if key === 'level'}{item.level}
        {:else if key === 'scope'}{item.bidang_id ? bidangMap.get(item.bidang_id)?.nama : ''}{item.faskes_id ? faskesMap.get(item.faskes_id)?.nama : ''}{!item.bidang_id && !item.faskes_id ? 'Dinas' : ''}
        {:else if key === 'action'}<Button size="sm" variant="outline" onclick={() => openPosition(item)}>Edit</Button>
        {/if}
      {/snippet}
    </DataTable>
  {/if}
</Card>

<Modal open={userOpen} title={editingUser ? 'Edit user' : 'Tambah user'} description="Posisi menentukan role efektif dan routing approval." onclose={() => (userOpen = false)}>
  <div class="grid gap-3 sm:grid-cols-2">
    <div><label class="nb-label" for="un">Username</label><Input id="un" bind:value={userForm.username} /></div>
    <div><label class="nb-label" for="unama">Nama</label><Input id="unama" bind:value={userForm.nama} /></div>
    {#if !editingUser}<div class="sm:col-span-2"><label class="nb-label" for="up">Password awal</label><Input id="up" type="password" bind:value={userForm.password} /></div>{/if}
    <div><label class="nb-label" for="upos">Posisi</label><Select id="upos" bind:value={userForm.position_id} options={positionOptions} placeholder="Pilih posisi" /></div>
    <div><label class="nb-label" for="uf">Faskes</label><Select id="uf" bind:value={userForm.faskes_id} options={faskesOptions} placeholder="Dinas / tanpa faskes" /></div>
    <div class="sm:col-span-2"><label class="nb-label" for="ub">Bidang</label><Select id="ub" bind:value={userForm.bidang_id} options={bidangOptions} placeholder="Pilih bidang" /></div>
  </div>
  <div class="mt-4 flex justify-end gap-2"><Button variant="outline" onclick={() => (userOpen = false)}>Batal</Button><Button variant="success" onclick={saveUser}>Simpan user</Button></div>
</Modal>

<Modal open={faskesOpen} title={editingFaskes ? 'Edit puskesmas' : 'Tambah puskesmas'} description="Daftar puskesmas terpisah dari daftar user." onclose={() => (faskesOpen = false)}>
  <div class="grid gap-3 sm:grid-cols-2">
    <div><label class="nb-label" for="fk">Kode</label><Input id="fk" bind:value={faskesForm.kode} /></div>
    <div><label class="nb-label" for="fn">Nama</label><Input id="fn" bind:value={faskesForm.nama} /></div>
    <div><label class="nb-label" for="fj">Jenis</label><Input id="fj" bind:value={faskesForm.jenis} /></div>
    <div><label class="nb-label" for="fw">Wilayah ID</label><Input id="fw" bind:value={faskesForm.wilayah_id} /></div>
  </div>
  <div class="mt-4 flex justify-end gap-2"><Button variant="outline" onclick={() => (faskesOpen = false)}>Batal</Button><Button variant="success" onclick={saveFaskes}>Simpan faskes</Button></div>
</Modal>

<Modal open={positionOpen} title={editingPosition ? 'Edit jabatan' : 'Tambah jabatan'} description="Jabatan struktural dapat diisi user berbeda tanpa mengubah kode workflow." onclose={() => (positionOpen = false)}>
  <div class="grid gap-3 sm:grid-cols-2">
    <div><label class="nb-label" for="pk">Kode</label><Input id="pk" bind:value={positionForm.kode} /></div>
    <div><label class="nb-label" for="pn">Nama jabatan</label><Input id="pn" bind:value={positionForm.nama} /></div>
    <div><label class="nb-label" for="pr">Role efektif</label><Select id="pr" bind:value={positionForm.role} options={roleOptions} /></div>
    <div><label class="nb-label" for="pl">Level</label><Input id="pl" type="number" bind:value={positionForm.level} /></div>
    <div><label class="nb-label" for="pb">Bidang</label><Select id="pb" bind:value={positionForm.bidang_id} options={bidangOptions} placeholder="Dinas" /></div>
    <div><label class="nb-label" for="pf">Faskes</label><Select id="pf" bind:value={positionForm.faskes_id} options={faskesOptions} placeholder="Dinas" /></div>
  </div>
  <div class="mt-4 flex justify-end gap-2"><Button variant="outline" onclick={() => (positionOpen = false)}>Batal</Button><Button variant="success" onclick={savePosition}>Simpan posisi</Button></div>
</Modal>
