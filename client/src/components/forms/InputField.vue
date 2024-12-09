<template>
  <label class="text-input">
    <div v-if="label != ''" class="label">{{ label }}</div>
    <slot name="icon-left" />
    <input ref="textField" :value="value" :type="type" :placeholder="placeholder" @input="onInput"
      @focus="toggleSuggestions(true)" @blur="toggleSuggestions(false)" />
    <div v-if="showSuggestions" class="suggestions">
      <div v-for="suggestion in filteredSuggestions" :key="suggestion" class="option" @click="setValue(suggestion)">
        {{ suggestion }}
      </div>
    </div>
    <slot name="icon-right" />
  </label>
</template>

<script lang="ts" setup>
import { useTemplateRef } from 'vue';
import { ref, computed } from 'vue';

const showSuggestions = ref(false);
const emit = defineEmits(['update', 'set-value']);
const textField = useTemplateRef('textField');

const props = defineProps({
  label: { type: String, default: '' },
  value: { type: String, default: '' },
  type: { type: String, default: 'text' },
  placeholder: { type: String, default: '' },
  suggestions: { type: Array<string>, default: () => [] },
});

const onInput = (event: Event) => {
  emit('update', (event.target as HTMLInputElement).value);
};

const setValue = (value: string) => {
  if (textField.value) textField.value.blur();
  emit('update', value);
};

const overridden = ref(false);
const initialValue = ref<string>('');
const toggleSuggestions = async (show: boolean) => {
  overridden.value = false;
  if (show) {
    initialValue.value = props.value;
    overridden.value = true;
    emit('update', '');
  }
  else await new Promise((resolve) => setTimeout(resolve, 200));
  if (!show) {
    if (!props.suggestions.includes(props.value)) {
      emit('update', initialValue.value);
    }
    emit('set-value', props.value);
  }

  showSuggestions.value = show ? true : overridden.value;
};

const filteredSuggestions = computed(() => {
  return props.suggestions.filter(
    (suggestion: string) => suggestion.toLowerCase().includes(props.value.toLowerCase())
  );
});
</script>

<style lang="scss" scoped>
.label {
  padding: 0.5rem 0;
}

.text-input {
  position: relative;
  width: 100%;
}

.suggestions {
  position: absolute;
  top: 100%;
  width: 100%;
  box-sizing: border-box;
  max-height: 300px;
  overflow-y: auto;
  padding: 0 0.5rem;
  background-color: var(--color-background-secondary-light);
  z-index: 10;

  .option {
    width: 100%;
    cursor: pointer;
    background-color: var(--color-background-secondary-light);
  }

  .option:hover {
    background-color: var(--color-background-primary);
  }
}

input {
  width: 100%;
  box-sizing: border-box;
  padding: 0.5rem;
  border-radius: 0.5rem;
  border: 1px solid var(--color-border-field);
  background-color: var(--color-background-secondary-light);
  color: var(--color-text-primary);
}
</style>
