import { defineStore } from "pinia";

interface State {
  loaders: Set<string>;
}

export const useLoaderStore = defineStore("loader", {
  state: (): State => ({
    loaders: new Set(),
  }),
  getters: {
    isLoading: (state) => (loader: string | string[]) => {
      if (Array.isArray(loader)) {
        return loader.some((l) => state.loaders.has(l));
      }
      return state.loaders.has(loader);
    },
  },
  actions: {
    start(loader: string) {
      this.loaders.add(loader);
    },
    end(loader: string) {
      this.loaders.delete(loader);
    }
  }
});

