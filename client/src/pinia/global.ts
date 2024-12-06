import { defineStore } from "pinia";
import { getConfig, updateConfig } from "~/apiClient/config";

interface State {
  windowWidth: number;
  autoConvert: boolean;
  deleteAfterConvert: boolean;
  overwriteExisting: boolean;
}

export const useGlobalStore = defineStore("global", {
  state: (): State => ({
    windowWidth: window.innerWidth,
    autoConvert: true,
    deleteAfterConvert: false,
    overwriteExisting: false,
  }),

  getters: {
    isMobile(state: State) {
      return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i
        .test(navigator.userAgent) || state.windowWidth < 768;
    },
  },

  actions: {
    async fetchSettings() {
      const config = await getConfig();
      this.autoConvert = config.autoConvert;
      this.deleteAfterConvert = config.deleteAfter;
      this.overwriteExisting = config.overwriteExisting;
    },
    async toggleAutoConvert() {
      this.autoConvert = !this.autoConvert;
      await updateConfig({ autoConvert: this.autoConvert });
    },
    async toggleDeleteAfterConvert() {
      this.deleteAfterConvert = !this.deleteAfterConvert;
      await updateConfig({ deleteAfter: this.deleteAfterConvert });
    },
    async toggleIgnoreExisting() {
      this.overwriteExisting = !this.overwriteExisting;
      await updateConfig({ overwriteExisting: this.overwriteExisting });
    }
  }
});
