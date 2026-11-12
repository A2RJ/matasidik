<script lang="ts">
  import Badge from '$lib/components/ui/badge.svelte';
  import { Check, X } from '@lucide/svelte';
  import {
    CLASSIFICATION_LABELS,
    CLASSIFICATION_TONE,
    DISPLAY_STATUS_LABELS,
    DISPLAY_STATUS_TONE,
    SUBMISSION_STATUS_LABELS,
    SUBMISSION_STATUS_TONE,
    TL_STATUS_LABELS,
    TL_STATUS_TONE,
    VALUE_STATUS_LABELS,
    VALUE_STATUS_TONE
  } from '$lib/constants/enums';
  import type { ClassificationStatus, SubmissionStatus } from '$lib/types';

  let {
    kind,
    value
  }: {
    kind: 'submission' | 'classification' | 'display' | 'value' | 'tl';
    value: string;
  } = $props();

  let label = $derived(
    kind === 'submission'
      ? SUBMISSION_STATUS_LABELS[value as SubmissionStatus] ?? value
      : kind === 'classification'
        ? CLASSIFICATION_LABELS[value as ClassificationStatus] ?? value
        : kind === 'display'
          ? DISPLAY_STATUS_LABELS[value] ?? value
          : kind === 'value'
            ? VALUE_STATUS_LABELS[value] ?? value
            : TL_STATUS_LABELS[value] ?? value
  );
  let tone = $derived(
    kind === 'submission'
      ? SUBMISSION_STATUS_TONE[value as SubmissionStatus] ?? 'neutral'
      : kind === 'classification'
        ? CLASSIFICATION_TONE[value as ClassificationStatus] ?? 'neutral'
        : kind === 'display'
          ? DISPLAY_STATUS_TONE[value] ?? 'neutral'
          : kind === 'value'
            ? VALUE_STATUS_TONE[value] ?? 'neutral'
            : TL_STATUS_TONE[value] ?? 'neutral'
  );
</script>

{#if kind === 'display'}
  <span
    class={`inline-flex h-6 w-6 items-center justify-center border-2 border-black ${
      value === 'ON_TARGET' ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'
    }`}
    title={label}
    aria-label={label}
  >
    {#if value === 'ON_TARGET'}
      <Check size={14} strokeWidth={3} />
    {:else}
      <X size={14} strokeWidth={3} />
    {/if}
  </span>
{:else}
  <Badge {tone}>{label}</Badge>
{/if}
