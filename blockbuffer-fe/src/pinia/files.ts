import { defineStore } from "pinia";
import type { File as MediaFile } from "~/types/files";
import { getFiles, uploadFiles } from "~/apiClient/files";

interface State {
  files: MediaFile[];
}
export const useFilesStore = defineStore("files", {
  state: (): State => ({
    files: [],
  }),
  actions: {
    async fetchFiles() {
      console.log("fetchFiles");
      const data = await getFiles();
      this.files = data;
    },
    async uploadFiles(files: File[]) {
      const data = await uploadFiles(files);
      console.log(data);
    },
  },
});
