<template>
  <div class="item">
    <Icon :name="icon" class="icon" />
    <div class="data">
      <h4>{{ name }}</h4>
      <div class="details">12 files: 1:43:10 total runtime</div>
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
  background: linear-gradient(180deg, var(--color-background-secondary), var(--color-background-primary-dark));

  .icon {
    align-self: center;
    justify-self: center;
    padding: var(--spacing-xl);
    background-color: var(--color-background-secondary-dark);
  }
}

.data {
  font: var(--text-body-sm);

  width: 100%;
  display: flex;
  flex-direction: column;
  padding: var(--spacing-md);

  .details {
    font: var(--text-body-xs);
    color: var(--color-text-caption);
  }
}
</style>
