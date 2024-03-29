import { defineStore } from "pinia";

import router from "@/router";
import axios from "@/plugins/axios";
import { handleError } from "@/utils/error-handling";
import {
  getFromLocalStorage,
  saveToLocalStorage,
  removeFromLocalStorage,
} from "@/utils/localStorage";
import { useToastWithTitle } from "@/plugins/toast";
import i18n from "@/plugins/i18n";

// User type
interface User {
  id: string;
  username: string;
  displayName: string;
  email: string;
  lang: string;
  role: string;
  created_at: string;
  updated_at: string;
}

function tryParseJSON(json: string) {
  try {
    return JSON.parse(json);
  } catch (error) {
    return JSON.parse("{}");
  }
}

export const useAuthStore = defineStore({
  id: "auth",
  state: () => ({
    // Initialize state from local storage to avoid reset on page refresh
    token: getFromLocalStorage("token", null),
    user: tryParseJSON(getFromLocalStorage("user", "{}")) as User,
    lastUpdatedAt: getFromLocalStorage("lastUpdatedAt", null),
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
  },
  actions: {
    async login(login: string, password: string, rememberMe = false) {
      try {
        const response = await axios.post("/auth/login", {
          login,
          password,
        });

        if (!response.data.token) {
          throw "MISSING_TOKEN";
        }

        // store jwt in local storage to keep user logged in between page refreshes if remember me is enabled
        const updatedDate = new Date().toISOString();
        if (rememberMe) {
          saveToLocalStorage("token", response.data.token);
          saveToLocalStorage("lastUpdatedAt", updatedDate);
          saveToLocalStorage("user", JSON.stringify(response.data.user));
        }

        // update pinia state
        this.token = response.data.token;
        this.lastUpdatedAt = updatedDate;
        this.user = response.data.user;
      } catch (error) {
        return handleError(error);
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
      } catch (error: any) {
        return handleError(error);
      }
    },
    async logout() {
      try {
        await axios.post("/auth/logout");
      } catch (error) {
        handleError(error);
      }

      this.token = null;
      this.lastUpdatedAt = null;
      this.user = {} as User;
      removeFromLocalStorage("token");
      removeFromLocalStorage("lastUpdatedAt");
      removeFromLocalStorage("user");

      router.push("/login");
    },
    async refreshToken() {},
    async requestPasswordReset(login: string) {
      try {
        const response = await axios.post("/auth/forgot", {
          login,
        });

        if (!response.data.success) {
          throw response;
        }

        const toast = useToastWithTitle();
        toast.success(
          i18n.global.t("auth.password_reset_requested"),
          i18n.global.t("auth.password_reset_requested_msg")
        );
      } catch (error) {
        return handleError(error);
      }
    },
    async resetPassword(token: string, password: string) {
      try {
        const response = await axios.post("/auth/reset", {
          token,
          password,
        });

        if (!response.data.success) {
          throw response;
        }

        const toast = useToastWithTitle();
        toast.success(
          i18n.global.t("routes.password_reset"),
          i18n.global.t("auth.password_reset_success")
        );
      } catch (error) {
        return handleError(error);
      }
    },
    async resendVerificationEmail(login: string) {
      try {
        const response = await axios.post("/auth/resend", {
          login,
        });

        if (!response.data.success) {
          throw response;
        }

        const toast = useToastWithTitle();
        toast.success(
          i18n.global.t("auth.account_verification"),
          i18n.global.t("auth.account_verification_resent")
        );
      } catch (error) {
        return handleError(error);
      }
    },
    async verifyAccount(token: string) {
      try {
        const response = await axios.get("/auth/verify/" + token);

        if (!response.data.success) {
          throw response;
        }

        const toast = useToastWithTitle();
        toast.success(
          i18n.global.t("auth.account_verification"),
          i18n.global.t("auth.account_verified_successfully")
        );
      } catch (error) {
        return handleError(error);
      }
    },
  },
});
