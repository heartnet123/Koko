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

  css: ['~/assets/css/main.css'],
})