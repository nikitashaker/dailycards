// File: src/main.js

// 👉 Включаем Vue Devtools (только в dev-режиме)
// Без этого плагин vite-plugin-vue-devtools не сможет поднять Devtools в браузере
// import 'virtual:vue-devtools'

import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

const app = createApp(App)

// Плагин для состояния (Pinia)
app.use(createPinia())

// Роутер (vue-router)
app.use(router)

app.mount('#app')
