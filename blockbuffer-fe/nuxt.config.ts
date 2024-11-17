import * as path from "path";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2024-04-03",
  devtools: { enabled: true },
  srcDir: "src",
  ssr: false,

  modules: ["@pinia/nuxt"],

  css: ["~/assets/styles/mixins.scss"],
  app: {
    head: {
      link: [{ rel: "preconnect", href: "https://fonts.gstatic.com" }, {
        rel: "preconnect",
        href: "https://fonts.gstatic.com",
      }, {
        rel: "stylesheet",
        href:
          "https://fonts.googleapis.com/css2?family=Lato:ital,wght@0,100;0,300;0,400;0,700;0,900;1,100;1,300;1,400;1,700;1,900&display=swap",
      }],
    },
  },
  nitro: {
    output: {
      publicDir: path.join(__dirname, "../public"),
    },
  },
});
