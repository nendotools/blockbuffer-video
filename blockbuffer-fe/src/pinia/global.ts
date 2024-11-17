import { defineStore } from "pinia";

interface State {
  windowWidth: number;
}

export const useGlobalStore = defineStore("global", {
  state: (): State => ({
    windowWidth: window.innerWidth,
  }),
  getters: {
    isMobile(state: State) {
      return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i
        .test(navigator.userAgent) || state.windowWidth < 768;
    },
  },
});
