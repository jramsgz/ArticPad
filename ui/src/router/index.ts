import { createRouter, createWebHistory } from "vue-router";
import routes from "./routes";
import { useLoadingStore } from "@/stores/loading";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

router.beforeEach((to, from, next) => {
  const loadingStore = useLoadingStore();
  // If this isn't an initial page load.
  if (to.name) {
    // Start the route progress bar.
    loadingStore.setRouteLoading(true);
  }
  next();
});

router.afterEach((to, from) => {
  const loadingStore = useLoadingStore();
  // Hide the route progress bar.
  loadingStore.setRouteLoading(false);
});

export default router;
