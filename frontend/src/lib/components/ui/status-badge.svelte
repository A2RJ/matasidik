<script lang="ts">
  import Badge from '$lib/components/ui/badge.svelte';
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

<Badge {tone}>{label}</Badge>
