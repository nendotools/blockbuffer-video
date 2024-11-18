import { defineStore } from "pinia";
import type { File } from "~/types/files";
import { useFetch } from "@/composables/useFetch";

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
