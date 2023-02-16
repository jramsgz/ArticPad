// Here are all the InjectionKeys for the plugins
// that are used in the app.
// https://vuejs.org/guide/components/provide-inject.html#working-with-symbol-keys

import type { InjectionKey } from "vue";
import type { AxiosInstance } from "axios";

export const axiosKey: InjectionKey<AxiosInstance> = Symbol();
