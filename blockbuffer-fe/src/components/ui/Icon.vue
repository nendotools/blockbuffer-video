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
}
</style>
