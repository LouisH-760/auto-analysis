import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'

const app = createApp(App)


import { loadFonts } from "@/plugins/webfontloader";
import vuetify from "@/plugins/vuetify";
loadFonts();
app.use(vuetify).use(createPinia())

app.mount('#app')
