import { defineStore } from "pinia";

export const useLoadingStore = defineStore("loading", {
  state: () => ({
    isRouteLoading: false,
    isApiLoading: false,
  }),
  actions: {
    setRouteLoading(isLoading: boolean) {
      this.isRouteLoading = isLoading;
    },
    setApiLoading(isLoading: boolean) {
      this.isApiLoading = isLoading;
    },
  },
});
