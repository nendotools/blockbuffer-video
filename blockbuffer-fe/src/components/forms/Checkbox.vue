<template>
  <label class="checkbox" for="checkbox" @click.prevent="$emit('toggle')">
    <div class="toggle" :class="{ active: checked }"></div>
    <input type="checkbox" id="checkbox" />
    {{ label }}
  </label>
</template>

<script lang="ts" setup>
defineEmits(["toggle"]);

withDefaults(defineProps<{
  label: string;
  checked: boolean;
}>(), {
  label: "Checkbox",
  checked: false,
});
</script>

<style lang="scss" scoped>
.checkbox {
  display: flex;
  flex-direction: row;
  position: relative;
  gap: var(--spacing-sm);
  align-items: center;
  cursor: pointer;

  .toggle {
    top: 0;
    left: 0;
    width: 20px;
    height: 20px;
    background-color: var(--black-tinted-lighter);
    border: 2px solid var(--background-color-lighter);
    border-radius: 4px;
    transition: background-color 0.3s;

    &.active {
      background-color: var(--primary-color-dim);
    }
  }
}

input[type="checkbox"] {
  display: none;
}
</style>
