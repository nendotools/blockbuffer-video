import { defineStore } from "pinia";
import { MessageTypes, type FileMessage, type File as MediaFile } from "~/types/files";
import { getFiles, uploadFiles } from "~/apiClient/files";
import { useWebSocket } from "~/composables/useWebSocket";
import { useLoaderStore } from "./loader";

export const MEDIA_UPLOAD_KEY = "media-upload";

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
      const loader = useLoaderStore();
      loader.start(MEDIA_UPLOAD_KEY);
      const data = await uploadFiles(files);
      loader.end(MEDIA_UPLOAD_KEY);
      console.log(data);
    },

    async updateFiles(message: FileMessage) {
      if (!message) {
        return;
      }

      switch (message.type) {
        case MessageTypes.DELETE_FILE:
          Object.keys(message.data).forEach((id: string) => {
            this.files = this.files.filter((file) => file.id !== id);
          });
          break;
        default:
          Object.values(message.data).forEach((file) => {
            const fileIndex = this.files.findIndex((f) => f.id === file.id);
            if (fileIndex === -1) {
              this.files.push(file);
            } else {
              this.files[fileIndex] = file;
            }
          });
      }
    }
  },
});
