<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from '#imports';
import { useGlobalStore } from '@/pinia/global';

const globalStore = useGlobalStore();
const title = computed(() => 'Home');
onMounted(() => {
  const resize = new ResizeObserver(() => {
    globalStore.windowWidth = window.innerWidth;
  });

  resize.observe(document.body);
});
</script>

<style lang="scss">
html {
  font-size: 62.5%;
}

body {
  font: var(--text-body-md);
  background-color: var(--color-background);
  color: var(--color-primary);
}

h1 {
  font: var(--text-ui-xl);
  padding: 0;
  margin: 0;
}

h2 {
  font: var(--text-ui-lg);
  padding: 0;
  margin: 0;
}

h3 {
  font: var(--text-ui-md);
  padding: 0;
  margin: 0;
}

h4 {
  font: var(--text-ui-sm);
  padding: 0;
  margin: 0;
}

h5 {
  font: var(--text-ui-xs);
  padding: 0;
  margin: 0;
}
</style>
