import { defineStore } from "pinia";
import type { FileMessage, File as MediaFile } from "~/types/files";
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

    async updateFiles(message: FileMessage) {
      if (!message) {
        return;
      }

      Object.values(message.data).forEach((file) => {
        const fileIndex = this.files.findIndex((f) => f.id === file.id);
        if (fileIndex === -1) {
          this.files.push(file);
        } else {
          this.files[fileIndex] = file;
        }
      });
    }
  },
});
