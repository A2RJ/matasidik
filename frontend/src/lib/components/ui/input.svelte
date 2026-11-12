<script lang="ts">
  import { cn } from '$lib/utils';
  import type { FullAutoFill } from 'svelte/elements';
  let {
    value = $bindable(''),
    type = 'text',
    placeholder = '',
    class: className = '',
    name = '',
    id = '',
    required = false,
    autocomplete = undefined,
    disabled = false,
    oninput,
    min,
    max,
    step
  }: {
    value?: string;
    type?: string;
    placeholder?: string;
    class?: string;
    name?: string;
    id?: string;
    required?: boolean;
    autocomplete?: FullAutoFill;
    disabled?: boolean;
    oninput?: (event: Event) => void;
    min?: string | number;
    max?: string | number;
    step?: string | number;
  } = $props();
</script>

<input
  bind:value
  {type}
  {placeholder}
  {name}
  {id}
  {required}
  {autocomplete}
  {disabled}
  {oninput}
  {min}
  {max}
  {step}
  class={cn('nb-input', className)}
/>
