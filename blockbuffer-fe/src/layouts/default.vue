<template>
  <div class="layout">
    <Header class="header">
      <template #left>
        <div class="logo">
          <h3>BlockBuffer</h3>
          <h5 class="logo-footer">video converter</h5>
        </div>
      </template>
      <template #right>
        <Button size="md" variant="plain" :bottom-highlight="false" class="header-btn">
          manual-queue
        </Button>
        <Button size="md" variant="plain" bottom-highlight class="header-btn">
          auto-queue
        </Button>
      </template>
    </Header>
    <div class="page">
      <slot />
    </div>
  </div>
</template>

<script lang="ts" setup>
import Header from '@/components/ui/Header.vue';
import Button from '@/components/ui/Button.vue';
</script>

<style scoped lang="scss">
.layout {
  position: absolute;
  width: 100%;
  height: 100%;
  display: grid;
  grid-template-rows: min-content 1fr;
}

.header {
  height: min-content;
  padding-bottom: 0;
  background-color: var(--color-background-primary) !important;

  .logo {
    color: var(--color-accent);

    .logo-footer {
      padding-left: var(--spacing-lg);
      font: var(--text-ui-xs);
    }
  }

  .header-btn {
    padding: var(--spacing-xl) var(--spacing-lg);
  }
}

.page {
  align-self: stretch;
  position: relative;
  background-color: var(--color-background-secondary) !important;
}
</style>
