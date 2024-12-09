<template>
  <InputField label="Video Codec" :value="codecOpt.input" :placeholder="codecOpt.selected ?? defaultEncoder?.name"
    :suggestions="videoCodecs" @update="updateCodec" @set-value="selectCodec" />

  <InputField label="Video Format" :value="formatOpt.input" :placeholder="formatOpt.selected ?? defaultEncoder?.format"
    :suggestions="videoFormats" @update="updateFormat" @set-value="selectFormat" />

  <InputField label="Audio Codec" :value="audioOpt.input"
    :placeholder="audioOpt.selected ?? defaultEncoder?.audioEncoder.name" :suggestions="audioCodecs"
    @update="updateAudio" @set-value="selectAudio" />

  <div class="opt-group">
    <InputField v-for="option in options" :key="option.name" :label="option.name" :value="''"
      :suggestions="option.options.map((o) => o.option)" />
  </div>

  <Button variant="primary" class="floating-upload" @click="saveEncoder">
    Save Changes
  </Button>
</template>

<script lang="ts" setup>
import { computed, ref, storeToRefs } from '#imports';
import Button from '@/components/ui/Button.vue';
import InputField from '@/components/forms/InputField.vue';
import { useFilesStore } from '@/pinia/files';
import { type EncoderProfile } from '@/types/encoders';

const emit = defineEmits<{
  (e: 'update', value: string): void;
  (e: 'set-value', value: string): void;
  (e: 'save-encoder', value: EncoderProfile): void;
}>();

type Option = { input: string, selected: string | null };
const codecOpt = ref<Option>({
  input: '',
  selected: null
});
const formatOpt = ref<Option>({
  input: '',
  selected: null
});
const audioOpt = ref<Option>({
  input: '',
  selected: null
});

const fileStore = useFilesStore();
const { defaultEncoder } = storeToRefs(fileStore);

const videoCodecs = computed(() => fileStore.encoders.video.map((encoder) => encoder.name));
const videoFormats = computed(() => {
  if (!defaultEncoder.value && !codecOpt.value.selected) return [];
  if (!codecOpt.value.selected || !videoCodecs.value.includes(codecOpt.value.selected)) {
    const codec = fileStore.encoders.video.find((encoder) => encoder.name === defaultEncoder?.value?.name);
    if (!codec) return [];
    return codec.formats;
  } else {
    const codec = fileStore.encoders.video.find((encoder) => encoder.name === codecOpt.value.selected);
    if (!codec) return [];
    return codec.formats;
  }
});

const options = computed(() => {
  if (!defaultEncoder.value) return [];
  if (codecOpt.value.selected) {
    const codec = fileStore.encoders.video.find((encoder) => encoder.name === codecOpt.value.selected);
    if (!codec) return [];
    return codec.options;
  }
  const codec = fileStore.encoders.video.find((encoder) => encoder.name === defaultEncoder.value!.name);
  return codec?.options || [];
});
const audioCodecs = computed(() => fileStore.encoders.audio.map((encoder) => encoder.name));

const selectCodec = (str: string) => codecOpt.value.selected = str;
const updateCodec = (str: string) => {
  codecOpt.value.input = str;
  if (videoCodecs.value.includes(str)) {
    const codec = fileStore.encoders.video.find((encoder) => encoder.name === str);
    if (codec) formatOpt.value.selected = formatOpt.value.input = codec.formats[0];
  }
}

const selectFormat = (str: string) => formatOpt.value.selected = str;
const updateFormat = (str: string) => formatOpt.value.input = str;

const selectAudio = (str: string) => audioOpt.value.selected = str;
const updateAudio = (str: string) => audioOpt.value.input = str;

const saveEncoder = () => {
  if (!codecOpt.value.selected || !formatOpt.value.selected || !audioOpt.value.selected) return;
  const encoder = fileStore.encoders.video.find((encoder) => encoder.name === codecOpt.value.selected);
  if (!encoder) return;
  const audioEncoder = fileStore.encoders.audio.find((encoder) => encoder.name === audioOpt.value.selected);
  if (!audioEncoder) return;

  const newEncoder: EncoderProfile = {
    name: codecOpt.value.selected,
    format: formatOpt.value.selected,
    audioEncoder: { name: audioOpt.value.selected, sampleRate: audioEncoder.sampleRates[0], options: [] },
    options: []
  };
  emit('save-encoder', newEncoder);
}
</script>

<style lang="scss" scoped></style>
