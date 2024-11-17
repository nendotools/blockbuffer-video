<template>
  <button :class="[size]">
    <slot />
  </button>
</template>

<script lang="ts" setup>
withDefaults(
  defineProps<{
    size: 'sm' | 'md' | 'lg';
  }>(),
  {
    size: 'md',
  });

</script>

<style lang="scss" scoped>
button {
  margin: 0;
  padding: 0;
  border: none;
  cursor: pointer;
  color: var(--color-primary);
  background-color: var(--color-background-lighter);
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
}
</style>
