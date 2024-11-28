import { useFetch } from "@/composables/useFetch";

export const getFiles = async () => useFetch("/files");
