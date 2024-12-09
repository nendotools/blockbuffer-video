import { defineStore } from "pinia";
import { MessageTypes, type FileMessage, type File as MediaFile } from "~/types/files";
import { getFiles, uploadFiles } from "~/apiClient/files";
import { useWebSocket } from "~/composables/useWebSocket";
import { useLoaderStore } from "./loader";
import type { Encoder, EncoderProfile } from "~/types/encoders";
import { getEncoders } from "~/apiClient/encoder";

export const MEDIA_UPLOAD_KEY = "media-upload";

interface State {
  files: MediaFile[];
  encoders: {
    video: Encoder[];
    audio: Encoder[];
  };
  defaultEncoder: EncoderProfile | null;
  selectedEncoder: EncoderProfile | null;
  ws: WebSocket | null;
}
export const useFilesStore = defineStore("files", {
  state: (): State => ({
    files: [],
    encoders: {
      video: Array<Encoder>(),
      audio: Array<Encoder>(),
    },
    defaultEncoder: null,
    selectedEncoder: null,
    ws: null,
  }),

  getters: {
    encoder: (state) => {
      if (!state.selectedEncoder) {
        return state.defaultEncoder;
      }
      return state.selectedEncoder;
    },
  },

  actions: {
    async initSocket() {
      this.ws = await useWebSocket("/ws", this.updateFiles);
    },
    async fetchEncoders() {
      const data = await getEncoders();
      this.encoders.video = data.videoEncoders;
      this.encoders.audio = data.audioEncoders;
      this.defaultEncoder = data.defaultEncoder;
    },

    async selectEncoder(encoder: EncoderProfile) {
      this.selectedEncoder = encoder;
    },
    async fetchFiles() {
      const data = await getFiles();
      this.files = data.sort((a, b) => a.filePath.localeCompare(b.filePath));
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
      this.files.sort((a, b) => a.filePath.localeCompare(b.filePath));
    }
  },
});
