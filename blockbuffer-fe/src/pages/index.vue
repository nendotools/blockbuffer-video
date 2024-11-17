<template>
  <h1>page content</h1>
  <h3>
    {{ isMobile }}
  </h3>
</template>

<script lang="ts" setup>
import { storeToRefs } from 'pinia';
import { useGlobalStore } from '@/pinia/global';

const globalStore = useGlobalStore();
const { isMobile } = storeToRefs(globalStore);
</script>
