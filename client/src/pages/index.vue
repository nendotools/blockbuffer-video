<template>
  <div class="main">
    <div v-if="!isMobile" class="menu">
      <h4>Conversion Settings</h4>

      <div class="opt-group">
        <Checkbox :checked="globalStore.autoConvert" label="Auto-convert" @toggle="globalStore.toggleAutoConvert" />
        <Checkbox :checked="globalStore.deleteAfterConvert" label="Delete after convert"
          @toggle="globalStore.toggleDeleteAfterConvert" />
        <Checkbox :checked="globalStore.overwriteExisting" label="Overwrite existing files"
          @toggle="globalStore.toggleIgnoreExisting" />
      </div>

      <h4>Encoder Settings</h4>

      <div class="opt-group">
        <div>Video Codec: {{ encoder?.name }}</div>
        <sub>Format: {{ encoder?.format }}</sub>
        <div>Audio Codec: {{ encoder?.audioEncoder?.name }}</div>
        <EncoderMenu @save-encoder="fileStore.selectEncoder" />
      </div>
    </div>
    <div class="file-list">
      <ListFile v-for="file in files" :key="file.id" fileType="video" :file="file" />
    </div>

    <Button variant="primary" class="floating-upload" :loading="uploading" @click="selectFiles">
      <template #icon-left>
        <Icon name="upload" size="md" />
      </template>
      upload video(s)
    </Button>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, storeToRefs } from '#imports';
import Icon from '@/components/ui/Icon.vue';
import Button from '@/components/ui/Button.vue';
import Checkbox from '@/components/forms/Checkbox.vue';
import ListFile from '@/components/elements/ListFile.vue';
import { useGlobalStore } from '@/pinia/global';
import { MEDIA_UPLOAD_KEY, useFilesStore } from '@/pinia/files';
import { useLoaderStore } from '~/pinia/loader';

const fileStore = useFilesStore();
const loaderStore = useLoaderStore();
const globalStore = useGlobalStore();
const { isMobile } = storeToRefs(globalStore);
const { encoder, files } = storeToRefs(fileStore);

const uploading = computed(() => loaderStore.isLoading(MEDIA_UPLOAD_KEY));

onMounted(async () => {
  globalStore.fetchSettings();
  await fileStore.initSocket();
  await fileStore.fetchEncoders();
});

const selectFiles = (event: Event) => {
  if (!event.target) return;
  const input = document.createElement('input');
  input.type = 'file';
  input.multiple = true;
  input.accept = 'video/*';

  input.onchange = () => {
    const files = [...(input.files || [])];
    if (files) {
      fileStore.uploadFiles(files);
    }
  };
  input.click();
};
</script>

<style scoped lang="scss">
.main {
  height: 100%;
  display: flex;
  flex-direction: row;
  overflow: hidden;
}

.menu {
  overflow-x: hidden;
  min-width: 300px;
  display: flex;
  flex-direction: column;
  padding: var(--spacing-md) var(--spacing-lg);
  background-color: var(--color-background-secondary);
}

.file-list {
  width: 100%;
  overflow-y: auto;
  padding: 0 var(--spacing-lg);
  padding-right: 0;
  background-color: var(--color-background-secondary);
}

.floating-upload {
  position: fixed;
  bottom: var(--spacing-lg);
  right: var(--spacing-xl);
  z-index: 100;
}

.opt-group {
  overflow-x: hidden;
  overflow-y: auto;
  min-height: 300px;
  padding: var(--spacing-sm);
  margin-top: var(--spacing-lg);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}
</style>
