<template>
  <div class="item">
    <Icon :name="icon" class="icon" />
    <div class="data">
      <h3>{{ name }}</h3>
    </div>
  </div>
</template>

<script lang="ts" setup>
import Icon from '@/components/ui/Icon.vue';

const props = defineProps<{
  fileType: 'video' | 'folder';
  name: string;
}>();

const icon = props.fileType === 'video' ? 'film' : 'folder';
</script>

<style scoped lang="scss">
.item {
  display: flex;
  flex-direction: row;
  background: linear-gradient(180deg, var(--color-background-secondary), var(--color-background-primary));

  .icon {
    align-self: center;
    justify-self: center;
    padding: var(--spacing-xl);
  }
}

.data {
  display: flex;
  flex-direction: column;
  width: 100%;
}
</style>
