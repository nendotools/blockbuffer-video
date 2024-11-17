interface State {
  files: any[];
}
export const useFilesStore = defineStore("files", {
  state: (): State => ({
    files: [],
  }),
  actions: {
    async fetchFiles() {
      console.log("fetchFiles");
      const data = await useFetch<any[]>("/files");
      this.files = data;
    },
  },
});
