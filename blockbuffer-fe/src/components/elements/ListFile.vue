<template>
  <div class="item">
    <Icon :name="icon" class="icon" />
    <div class="data">
      <h4>{{ name }}</h4>
      <div class="details">1 file: 1:43:10 total runtime</div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from '#imports';
import Icon from '@/components/ui/Icon.vue';
import File from '@/types/files';

const props = defineProps<{
  file: File;
  fileType: 'video' | 'folder';
}>();

const icon = props.fileType === 'video' ? 'film' : 'folder';

const name = computed(() => props.file.filePath.split('/').pop());
</script>

<style scoped lang="scss">
.item {
  display: flex;
  flex-direction: row;
  background: linear-gradient(180deg, var(--color-background-secondary), var(--color-background-primary-dark));
  border-bottom: 1px solid var(--color-border-field);

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
