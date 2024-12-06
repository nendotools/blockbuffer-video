<template>
  <div class="item">
    <div class="icon-space">
      <Icon :name="icon" class="main-icon" :class="{ tinted: completed || warning || error }" />
      <Icon v-if="completed" name="check" class="icon-overlap completed" />
      <Icon v-if="warning" name="alert-triangle" class="icon-overlap warning" />
      <Icon v-if="error" name="alert-octagon" class="icon-overlap error" />
    </div>
    <div class="item-data">
      <div class="data">
        <h4>{{ name }}</h4>
        <div class="details">1 file: {{ timecode }} total runtime</div>
        <div class="status">Status: <span :class="{ completed, warning, error }">{{ file.status }}</span></div>
      </div>
    </div>
    <ProgressBar :value="file.progress / 100" smooth-transition />
  </div>
</template>

<script lang="ts" setup>
import { watch, computed, onMounted, ref } from '#imports';
import Icon from '@/components/ui/Icon.vue';
import ProgressBar from '@/components/ui/ProgressBar.vue';
import { FileStatuses, type File } from '@/types/files';

const props = defineProps<{
  file: File;
  fileType: 'video' | 'folder';
}>();

const icon = props.fileType === 'video' ? 'film' : 'folder';
const completed = ref(false);
const warning = ref(false);
const error = ref(false);

const name = computed(() => props.file.filePath.split('/').pop());
const timecode = computed(() =>
  // convert seconds to HH:MM:SS
  new Date(1000 * props.file.duration).toISOString().substring(11, 19),
);
onMounted(() => {
  statusDisplay(props.file.status);
});

watch(() => props.file.status, async (status) => {
  await new Promise((resolve) => setTimeout(resolve, 500));
  statusDisplay(status);
});

const statusDisplay = (status: FileStatuses) => {
  console.log(status);

  completed.value = false;
  warning.value = false;
  error.value = false;
  switch (status) {
    case FileStatuses.COMPLETED:
      completed.value = true;
      break;
    case FileStatuses.REJECTED:
      warning.value = true;
      break;
    case FileStatuses.FAILED:
      error.value = true;
  }
};
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

.completed {
  color: var(--success-color);
}

.warning {
  color: var(--warning-color);
}

.error {
  color: var(--error-color);
}
</style>
