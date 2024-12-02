<template>
  <div class="main">
    <div v-if="!isMobile" class="menu">
      menu
    </div>
    <div class="file-list">
      <ListFile v-for="file in files" :key="file.id" fileType="video" :file="file" />
    </div>

    <Button variant="primary" class="floating-upload" @click="selectFiles">
      <template #icon-left>
        <Icon name="upload" size="md" />
      </template>
      upload video(s)
    </Button>
  </div>
</template>

<script lang="ts" setup>
import { onMounted } from '#imports';
import Icon from '@/components/ui/Icon.vue';
import Button from '@/components/ui/Button.vue';
import ListFile from '@/components/elements/ListFile.vue';
import { storeToRefs } from 'pinia';
import { useGlobalStore } from '@/pinia/global';
import { useFilesStore } from '@/pinia/files';

const fileStore = useFilesStore();
const globalStore = useGlobalStore();
const { isMobile } = storeToRefs(globalStore);
const { files } = storeToRefs(fileStore);

onMounted(async () => {
  await fileStore.initSocket();
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
}

.menu {
  min-width: 200px;
  display: flex;
  flex-direction: column;
  padding: var(--spacing-md) var(--spacing-lg);
  background-color: var(--color-background-secondary);
}

.file-list {
  width: 100%;
  overflow-y: auto;
  padding: var(--spacing-md) var(--spacing-lg);
  padding-right: 0;
  background-color: var(--color-background-secondary);
}

.floating-upload {
  position: fixed;
  bottom: var(--spacing-lg);
  right: var(--spacing-xl);
  z-index: 100;
}
</style>
