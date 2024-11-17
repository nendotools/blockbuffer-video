import * as path from "path";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2024-04-03",
  devtools: { enabled: true },
  nitro: {
    output: {
      publicDir: path.join(__dirname, "../public"),
    },
  },
});
