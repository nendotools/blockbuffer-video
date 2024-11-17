<template>
  <div class="main">
    <div class="menu">
      menu
    </div>
    <div class="file-list">
      <ListFile fileType="video" name="video1" />
      <ListFile fileType="folder" name="folder1" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import Icon from '@/components/ui/Icon.vue';
import ListFile from '@/components/elements/ListFile.vue';
import { storeToRefs } from 'pinia';
import { useGlobalStore } from '@/pinia/global';

const globalStore = useGlobalStore();
const { isMobile } = storeToRefs(globalStore);
</script>

<style scoped lang="scss">
.main {
  height: 100%;
  display: flex;
  flex-direction: row;
}

.menu {
  display: flex;
  flex-direction: column;
  padding: var(--spacing-md) var(--spacing-lg);
  background-color: var(--color-background-secondary);
}

.file-list {
  width: 100%;
  padding: var(--spacing-md) var(--spacing-lg);
  background-color: var(--color-background-secondary);
}
</style>
