import type { App } from "vue";

import axios from "axios";
import type { AxiosInstance } from "axios";
import { axiosKey } from "@/plugins/keys";

const axiosInstance: AxiosInstance = axios.create({
  baseURL: `${import.meta.env.VITE_API_URL}/api/v1`,
});

// Set the AUTH token for any request
axiosInstance.interceptors.request.use(
  (config) => {
    const urlsExcludedForBearerHeader = [
      "/login",
      "/forgot",
      "/register",
      "/reset",
    ];

    console.log(config.url);

    if (!urlsExcludedForBearerHeader.includes(config.url as string)) {
      config.headers.Authorization = `Bearer ${localStorage.getItem("token")}`;
    }
    return config;
  },
  (error) => {
    Promise.reject(error);
  }
);

axios.interceptors.response.use(
  (response) => {
    if (!response.config.url?.includes("/refresh")) {
      // TODO: check if the token is expired and refresh it
    }
    return response;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default axiosInstance;
// Install function executed by Vue.use()
export const createAxios = (app: App) => {
  // Inject axios to the context as axiosKey (https://vuejs.org/api/composition-api-dependency-injection.html) (https://vuejs.org/guide/components/provide-inject.html#working-with-symbol-keys)
  app.provide(axiosKey, axiosInstance);
};
