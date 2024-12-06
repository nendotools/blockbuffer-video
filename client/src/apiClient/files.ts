import { useFetch } from "@/composables/useFetch";
import { type File as MediaFile } from "~/types/files";

export const getFiles = async () => useFetch<MediaFile[]>("/files");
export const uploadFiles = async (files: File[]) => {
  const formData = new FormData();
  files.forEach(f => formData.append('files', f));
  return useFetch("/upload", { method: "POST", data: formData });
}
