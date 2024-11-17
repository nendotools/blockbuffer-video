<template>
  <slot name="outer-icon-left" />
  <div>
    <slot name="inner-icon-left" />
    <input v-model="modelValue" :type="type" :placeholder="placeholder" :class="inputClasses" @input="onInput" />
    <slot name="inner-icon-right" />
  </div>
  <slot name="outer-icon-right" />
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue';
import Icon from '@/components/ui/Icon.vue';
</script>
