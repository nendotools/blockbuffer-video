<template>
  <div class="main">
    <div class="menu">
      menu
    </div>
    <div class="file-list">
      <ListFile v-for="file in files" :key="file.id" fileType="video" :file="file" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import Icon from '@/components/ui/Icon.vue';
import ListFile from '@/components/elements/ListFile.vue';
import { storeToRefs } from 'pinia';
import { useGlobalStore } from '@/pinia/global';
import { useFilesStore } from '@/pinia/files';

const fileStore = useFilesStore();
const globalStore = useGlobalStore();
const { isMobile } = storeToRefs(globalStore);
const { files } = storeToRefs(fileStore);

onMounted(async () => {
  await fileStore.fetchFiles();
});
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
</style>
