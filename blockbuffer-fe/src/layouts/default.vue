<template>
  <Header>
    <template #left>
      <h3>Agile Flow</h3>
    </template>
    <template #right>
      <Button size="md">
        menu
      </Button>
    </template>
  </Header>
  <div>
    <slot />
  </div>
</template>

<script lang="ts" setup>
import Header from '@/components/ui/Header.vue';
import Button from '@/components/ui/Button.vue';
</script>
