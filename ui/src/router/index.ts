import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
      meta: {
        showNavbar: true,
      },
    },
    {
      path: "/login",
      name: "login",
      // route level code-splitting
      // this generates a separate chunk (Login.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import("../views/auth/LoginView.vue"),
    },
    {
      path: "/register",
      name: "register",
      component: () => import("../views/auth/RegisterView.vue"),
    },
    {
      path: "/password-reset",
      name: "password-reset",
      component: () => import("../views/auth/PasswordResetView.vue"),
    },
    {
      // will match anything starting with `/settings/` and put it under `$route.params.afterSettings`
      path: "/settings/:afterSettings(.*)",
      name: "settings",
      component: () => import("../views/SettingsView.vue"),
      meta: {
        showNavbar: true,
      },
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import("../views/NotFoundView.vue"),
    },
  ],
});

export default router;
