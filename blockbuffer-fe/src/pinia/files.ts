import { defineStore } from "pinia";
import type { File as MediaFile } from "~/types/files";
import { getFiles, uploadFiles } from "~/apiClient/files";
import { useWebSocket } from "~/composables/useWebSocket";

interface State {
  files: MediaFile[];
  ws: WebSocket | null;
}
export const useFilesStore = defineStore("files", {
  state: (): State => ({
    files: [],
    ws: null,
  }),
  actions: {
    async initSocket() {
      this.ws = await useWebSocket("/ws", this.updateFiles);
    },
    async fetchFiles() {
      console.log("fetchFiles");
      const data = await getFiles();
      this.files = data;
    },
    async uploadFiles(files: File[]) {
      const data = await uploadFiles(files);
      console.log(data);
    },

    async updateFiles(files: Record<string, MediaFile>) {
      if (!files) {
        return;
      }

      this.files = Object.values(files).map((f) => ({
        ...f,
        progress: Number(f.progress)
      }));
    }
  },
});
