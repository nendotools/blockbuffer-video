<template>
  <div class="header">
    <div class="slot">
      <slot name="left" />
    </div>
    <div class="slot">
    </div>
    <slot />
    <div class="slot">
      <slot name="right" />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  padding: var(--spacing-md) var(--spacing-xxl);
}
</style>
