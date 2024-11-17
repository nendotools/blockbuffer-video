<template>
  <i data-feather="name" :class="iconClasses" />
</template>

<script lang="ts" setup>
import { computed } from 'vue';

const props = withDefaults(defineProps<{
  name: string;
  size: "sm" | "md" | "lg";
}>(),
  {
    size: 'md',
  });

const iconClasses = computed(() => {
  return [
    'icon',
    `icon-${props.size}`,
  ];
});
</script>

<style>
@import 'https://unpkg.com/feather-icons';
</style>

<style lang="scss" scoped>
.icon {
  display: inline-block;
  font-size: 1rem;

  &icon-sm {
    font-size: 0.75rem;
  }

  &icon-md {
    font-size: 1rem;
  }

  &icon-lg {
    font-size: 1.25rem;
  }
}
</style>
