import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import router from "./router";
import axios from "./plugins/axios";

import "./assets/main.css";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(axios, {
  baseUrl: import.meta.env.VITE_API_URL as string,
});

app.mount("#app");
