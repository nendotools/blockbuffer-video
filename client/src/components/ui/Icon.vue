<template>
  <div :class="iconClasses">
    <i :data-feather="name" />
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted } from 'vue';
import feather from 'feather-icons';

const props = withDefaults(defineProps<{
  name: string;
  size?: "sm" | "md" | "lg";
  animation?: "spin" | null;
}>(),
  {
    size: 'md',
    animation: null
  });

const iconClasses = computed(() => {
  return [
    'icon',
    `icon-${props.size}`,
    { 'icon-spin': props.animation === 'spin' }
  ];
});

onMounted(() => {
  feather.replace();
});
</script>

<style lang="scss" scoped>
.icon {
  display: inline-block;

  svg {
    width: 1rem;
    height: 1rem;
  }

  &.icon-sm {
    svg {
      width: 1rem;
      height: 1rem;
    }
  }

  &.icon-md {
    svg {
      width: 2rem;
      height: 2rem;
    }
  }

  &.icon-lg {
    svg {
      width: 3rem;
      height: 3rem;
    }
  }

  &.icon-spin {
    animation: spin 1s linear infinite;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
