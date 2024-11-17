import type { File } from "~/types/files";

interface State {
  files: File[];
}
export const useFilesStore = defineStore("files", {
  state: (): State => ({
    files: [],
  }),
  actions: {
    async fetchFiles() {
      console.log("fetchFiles");
      const data = await useFetch<File[]>("/files");
      this.files = data;
    },
  },
});
