<template>
  <button :class="[size, variant, { bottomHighlight }]">
    <slot />
  </button>
</template>

<script lang="ts" setup>
withDefaults(
  defineProps<{
    variant: 'primary' | 'plain';
    size: 'sm' | 'md' | 'lg';
    bottomHighlight?: boolean;
  }>(),
  {
    variant: 'plain',
    size: 'md',
    bottomHighlight: false,
  });

</script>

<style lang="scss" scoped>
button {
  margin: 0;
  padding: 0;
  border: none;
  cursor: pointer;
  color: var(--color-text-primary);
  background-color: color-mix(in srgb, var(--color-background-primary), var(--color-primary));
  border-radius: var(--border-radius-md);

  &.sm {
    font: var(--text-ui-xs);
    padding: var(--spacing-xs) var(--spacing-sm);
  }

  &.md {
    font: var(--text-ui-sm);
    padding: var(--spacing-sm) var(--spacing-md);
  }

  &.lg {
    font: var(--text-ui-md);
    padding: var(--spacing-md) var(--spacing-lg);
  }

  &.primary {
    background-color: var(--color-primary);
  }

  &.plain {
    background-color: unset;
    color: var(--color-text-primary);
  }

  &.bottomHighlight {
    color: color-mix(in srgb, var(--color-primary) 30%, var(--color-text-primary));
    border-bottom: 2px solid var(--color-primary);
  }
}
</style>
