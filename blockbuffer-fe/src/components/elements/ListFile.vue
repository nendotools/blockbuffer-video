<template>
  <div class="item">
    <div class="icon-space">
      <Icon :name="icon" class="main-icon" :class="{ tinted: completed }" />
      <Icon v-if="completed" name="check" class="icon-overlap" />
    </div>
    <div class="item-data">
      <div class="data">
        <h4>{{ name }}</h4>
        <div class="details">1 file: {{ timecode }} total runtime</div>
        <div class="status">Status: {{ file.status }}</div>
      </div>
    </div>
    <ProgressBar :value="file.progress / 100" smooth-transition />
  </div>
</template>

<script lang="ts" setup>
import { watch, computed, onMounted, ref } from '#imports';
import Icon from '@/components/ui/Icon.vue';
import ProgressBar from '@/components/ui/ProgressBar.vue';
import { type File } from '@/types/files';

const props = defineProps<{
  file: File;
  fileType: 'video' | 'folder';
}>();

const icon = props.fileType === 'video' ? 'film' : 'folder';
const completed = ref(false);

const name = computed(() => props.file.filePath.split('/').pop());
const timecode = computed(() =>
  // convert seconds to HH:MM:SS
  new Date(1000 * props.file.duration).toISOString().substring(11, 19),
);
onMounted(() => {
  if (props.file.progress === 100) {
    completed.value = true;
  }
});

watch(() => props.file.progress, async (progress) => {
  if (progress === 100) {
    await new Promise((resolve) => setTimeout(resolve, 500));
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

    .main-icon {
      align-self: center;
      justify-self: center;
      padding: var(--spacing-xl);
      background-color: var(--color-background-secondary-dark);

      &.tinted {
        color: var(--black-color-lighter);
      }
    }

    .icon-overlap {
      position: relative;
      width: 0;
      bottom: -5%;
      right: 50%;
      color: var(--success-color);
      border-radius: 50%;
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

  .status {
    position: relative;
    bottom: 0;
    right: 0;
    text-align: right;
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
  background: linear-gradient(5deg, var(--color-background-primary-dark), var(--color-background-secondary-dark));
  border: none;
}
</style>
