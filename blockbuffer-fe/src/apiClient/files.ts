import { useFetch } from "#imports";

export const getFiles = async () => useFetch("/files");
