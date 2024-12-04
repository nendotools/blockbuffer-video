<template>
  <div v-if="show || smoothTransition" class="container" :style="style">
    <progress v-if="isPercentage" class="progress" :value="percentage" max="100"></progress>
    <div v-else-if="countdown" class="timeout" :style="countdownStyle"></div>
    <div v-else class="indeterminate"></div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from "#imports";
import { useLoaderStore } from "~/pinia/loader";
const loaderStore = useLoaderStore();

const props = withDefaults(
  defineProps<{
    value?: number;
    countdown?: number;
    smoothTransition?: boolean;
    loaderKeys?: string[];
  }>(),
  {
    value: -1,
    countdown: 0,
    smoothTransition: false,
    loaderKeys: () => [],
  },
);

const show = computed(
  () =>
    (props.value > 0 && props.value < 1) ||
    props.countdown ||
    (props.loaderKeys.length > 0 && loaderStore.isLoading(props.loaderKeys)),
);
const isPercentage = computed(() => props.value >= 0 && props.value <= 1);
const percentage = computed(() =>
  isPercentage.value
    ? Math.floor(Math.max(0, Math.min(100, props.value * 100)))
    : 0,
);

const style = computed(() => ({
  height: props.smoothTransition ? (show.value ? "5px" : "1px") : "5px",
  opacity: props.smoothTransition ? (show.value ? "1" : "0") : "1",
  transition: props.smoothTransition ? "opacity 0.3s, height 0.5s" : "none",
}));

const countdownStyle = computed(() => ({
  "--countdown": `${props.countdown}s`,
}));
</script>

<style lang="scss" scoped>
.progress,
progress[value] {
  width: 100%;
  border: none;
  display: block;
  appearance: none;
  -webkit-appearance: none;

  &::-webkit-progress-bar {
    background-color: var(--black-tinted-lighter);

  }

  &::-webkit-progress-value {
    background-color: var(--color-primary);
    transition: width 2s linear;
  }
}

.container {
  margin: 0;
  position: relative;
  height: 4px;
  display: block;
  width: 100%;
  background-color: var(--black-tinted-lighter);
  border-radius: 2px;
  overflow: hidden;

  .indeterminate {
    background-color: var(--color-primary);

    &:before {
      content: "";
      position: absolute;
      background-color: inherit;
      top: 0;
      left: 0;
      bottom: 0;
      will-change: left, right;
      animation: indeterminate-short 2s cubic-bezier(0.07, 0.32, 0.19, 0.73) infinite;
    }

    &:after {
      content: "";
      position: absolute;
      background-color: inherit;
      top: 0;
      left: 0;
      bottom: 0;
      will-change: left, right;
      animation: indeterminate-short 2s cubic-bezier(0.165, 0.84, 0.44, 1) infinite;
      animation-delay: 1s;
    }
  }

  .timeout {
    background-color: var(--color-primary);
    animation: countdown var(--countdown) linear;
    animation-fill-mode: forwards;
    height: 4px;
  }
}

@keyframes countdown {
  0% {
    width: 100%;
  }

  100% {
    width: 0;
  }
}

@keyframes inderminate {
  0% {
    left: -15%;
    right: 100%;
  }

  40% {
    left: 100%;
    right: -10%;
  }

  100% {
    left: 100%;
    right: -10%;
  }
}

@keyframes indeterminate-short {
  0% {
    left: -200%;
    right: 100%;
  }

  70% {
    left: 100%;
    right: -8%;
  }

  100% {
    left: 107%;
    right: -8%;
  }
}
</style>
