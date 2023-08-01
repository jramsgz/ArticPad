import App from "./App.vue";

import { createApp } from "vue";
import { createPinia } from "pinia";
import router from "./router";
import Toast from "vue-toastification";
import { options as ToastOptions } from "./plugins/toast";

import "./assets/main.css";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(Toast, ToastOptions);

app.mount("#app");
