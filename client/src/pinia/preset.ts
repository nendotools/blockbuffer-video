import { defineStore } from 'pinia'
import { getPresets } from '~/apiClient/presets';
import type { Preset } from '~/types/presets';

interface State {
  presets: Preset[];
}

export const usePresetStore = defineStore('preset', {
  state: (): State => ({
    presets: [],
  }),
  actions: {
    async fetchPresets() {
      const response = await getPresets();
      this.presets = response.presets;
      return response;
    }
  },
});
