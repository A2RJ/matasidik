<script lang="ts">
  import type { FullAutoFill } from 'svelte/elements';
  import Input from '$lib/components/ui/input.svelte';
  import { Eye, EyeOff } from '@lucide/svelte';

  let {
    value = $bindable(''),
    id = '',
    placeholder = '••••••••',
    autocomplete = undefined,
    required = false,
    disabled = false,
    oninput
  }: {
    value?: string;
    id?: string;
    placeholder?: string;
    autocomplete?: FullAutoFill;
    required?: boolean;
    disabled?: boolean;
    oninput?: (event: Event) => void;
  } = $props();

  let visible = $state(false);
</script>

<div class="relative">
  <Input
    {id}
    type={visible ? 'text' : 'password'}
    bind:value
    {placeholder}
    {autocomplete}
    {required}
    {disabled}
    {oninput}
    class="pr-14"
  />
  <button
    type="button"
    class="absolute right-1.5 top-1/2 flex h-8 w-9 -translate-y-1/2 items-center justify-center border-2 border-black bg-card disabled:opacity-50"
    onclick={() => (visible = !visible)}
    {disabled}
    aria-label={visible ? 'Sembunyikan password' : 'Tampilkan password'}
    aria-pressed={visible}
    title={visible ? 'Sembunyikan password' : 'Tampilkan password'}
  >
    {#if visible}<EyeOff size={16} />{:else}<Eye size={16} />{/if}
  </button>
</div>
