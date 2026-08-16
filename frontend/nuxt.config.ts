// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  app: {
    head: {
      title: 'Koko - Anime Library',
      meta: [
        { name: 'description', content: 'Koko - Your ultimate anime library to discover and explore anime movies and series.' }
      ]
    }
  },

  modules: [
    '@nuxt/ui',
    '@nuxt/fonts',
    '@nuxt/image',
    '@nuxt/eslint',
  ],

  fonts: {
    families: [
      { name: 'Geist', weights: [400, 500, 600, 700] },
      { name: 'Geist Mono', weights: [400, 500, 600, 700] }
    ]
  },

  css: ['~/assets/css/main.css'],
})
