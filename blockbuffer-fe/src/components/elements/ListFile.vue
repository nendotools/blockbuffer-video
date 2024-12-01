<template>
  <div class="item">
    <div class="icon-space">
      <Icon :name="icon" class="icon" />
    </div>
    <div class="item-data">
      <div class="data">
        <h4>{{ name }}</h4>
        <div class="details">1 file: 1:43:10 total runtime</div>
      </div>
    </div>
    <div class="progress" v-if="!completed" :style="{ width: `${file.progress}%` }"></div>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from '#imports';
import Icon from '@/components/ui/Icon.vue';
import { type File } from '@/types/files';

const props = defineProps<{
  file: File;
  fileType: 'video' | 'folder';
}>();

const icon = props.fileType === 'video' ? 'film' : 'folder';
const completed = ref(false);

const name = computed(() => props.file.filePath.split('/').pop());
onMounted(() => {
  if (props.file.progress === 100) {
    completed.value = true;
  }
});
</script>

<style scoped lang="scss">
.item {
  width: 100%;
  display: grid;
  grid-template-columns: auto 1fr;
  grid-template-rows: 1fr min-content;

  .icon-space {
    grid-row: 1/3;
    grid-column: 1/2;

    .icon {
      align-self: center;
      justify-self: center;
      padding: var(--spacing-xl);
      background-color: var(--color-background-secondary-dark);
    }
  }
}

.item-data {
  display: flex;
  flex-direction: row;
  background: linear-gradient(180deg, var(--color-background-secondary), var(--color-background-primary-dark));
  border-bottom: 1px solid var(--color-border-field);
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

.progress {
  grid-column: 2/3;
  position: relative;
  bottom: 0;
  height: 4px;
  margin: 0;
  padding: 0;
  background: linear-gradient(5deg, var(--color-primary), var(--color-tertiary));
  border: none;
}
</style>
