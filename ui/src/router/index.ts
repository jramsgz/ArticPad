import { createRouter, createWebHistory } from "vue-router";
import routes from "@/router/routes";
import { useLoadingStore } from "@/stores/loading";
import { useAuthStore } from "@/stores/auth";
import i18n from "@/plugins/i18n";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

router.beforeEach((to, from, next) => {
  const loadingStore = useLoadingStore();
  const authStore = useAuthStore();

  // If this isn't an initial page load.
  if (to.name) {
    // Start the route progress bar.
    loadingStore.setRouteLoading(true);
  }

  // Redirect to login if route requires auth and user isn't logged in.
  if (
    to.matched.some((record) => record.meta.requieresAuth) &&
    !authStore.isLoggedIn
  ) {
    next({
      name: "login",
      query: { redirect: to.fullPath },
    });
    return;
  }
  // Redirect based on redirect query param if user is logged in.
  if (authStore.isLoggedIn && to.name === "login") {
    const redirect = to.query.redirect as string;
    if (redirect) {
      next(redirect);
      return;
    }
    next({ name: "home" });
    return;
  }
  // Redirect to home if user is logged in and tries to access auth routes.
  if (
    authStore.isLoggedIn &&
    to.matched.some((record) => record.meta.requieresLoggedOut)
  ) {
    next({ name: "home" });
    return;
  }

  next();
});

router.afterEach((to, from) => {
  const loadingStore = useLoadingStore();
  // Hide the route progress bar.
  loadingStore.setRouteLoading(false);
  // Set page title.
  document.title = to.meta.title
    ? i18n.global.t(to.meta.title as string) + " - ArticPad"
    : "ArticPad";
});

export default router;
