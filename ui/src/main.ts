import App from "./App.vue";

import { createApp } from "vue";
import { createPinia } from "pinia";
import router from "./router";
import { createAxios } from "./plugins/axios";
import Toast from "vue-toastification";
import { options as ToastOptions } from "./plugins/toast";
import i18n from "./plugins/i18n";

import "./assets/main.css";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(createAxios);
app.use(Toast, ToastOptions);
app.use(i18n);

app.mount("#app");
