import axios from "axios";
import { axiosKey } from "@/plugins/keys";
import type { AxiosInstance } from "axios";
import type { App } from "vue";

interface AxiosOptions {
  baseUrl?: string;
  token?: string;
}

// Note, we use the composition API here
// https://vuejs.org/guide/extras/composition-api-faq.html
export default {
  install: (app: App, options: AxiosOptions) => {
    // Create a new axios instance
    const axiosInstance: AxiosInstance = axios.create({
      baseURL: options.baseUrl,
    });

    // Set the AUTH token for any request
    axiosInstance.interceptors.request.use((config) => {
      if (options.token) {
        config.headers.Authorization = `Bearer ${options.token}`;
      }
      return config;
    });

    // Inject axios to the context as axiosKey (https://vuejs.org/api/composition-api-dependency-injection.html) (https://vuejs.org/guide/components/provide-inject.html#working-with-symbol-keys)
    app.provide(axiosKey, axiosInstance);
  },
};
