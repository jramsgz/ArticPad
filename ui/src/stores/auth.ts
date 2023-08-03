import { defineStore } from "pinia";

import router from "@/router";
import axios from "@/plugins/axios";
import { handleError } from "@/utils/error-handling";
import { useToastWithTitle } from "@/plugins/toast";
import i18n from "@/plugins/i18n";

export const useAuthStore = defineStore({
  id: "auth",
  state: () => ({
    // Initialize state from local storage to avoid reset on page refresh
    token: localStorage.getItem("token") || null,
    user: localStorage.getItem("user") || null,
    lastUpdatedAt: localStorage.getItem("lastUpdatedAt") || null,
  }),
  actions: {
    async login(login: string, password: string) {
      try {
        const response = await axios.post("/auth/login", {
          login,
          password,
        });

        if (!response.data.token) {
          throw new Error("MISSING_TOKEN");
        }

        // store jwt in local storage to keep user logged in between page refreshes
        localStorage.setItem("token", response.data.token);
        localStorage.setItem("lastUpdatedAt", new Date().toISOString());

        // update pinia state
        this.token = response.data.token;

        // store user details and jwt in local storage to keep user logged in between page refreshes
        //localStorage.setItem("user", JSON.stringify(user));

        // redirect to previous url or default to home page
        router.push(this.returnUrl || "/");
      } catch (error) {
        console.log(error);
        handleError(error);
      }
    },
    async register(username: string, email: string, password: string) {
      try {
        const response = await axios.post("/auth/register", {
          username,
          email,
          password,
        });

        if (!response.data.success) {
          throw response;
        }

        const toast = useToastWithTitle();
        toast.success(
          i18n.global.t("auth.sign_up"),
          i18n.global.t("auth.account_created")
        );

        // Redirect to login page
        router.push("/login");
      } catch (error) {
        console.log(error);
        handleError(error);
      }
    },
  },
});
