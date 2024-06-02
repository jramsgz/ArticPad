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
    auth_token: getFromLocalStorage("auth_token", null),
    refresh_token: getFromLocalStorage("refresh_token", null),
    user: tryParseJSON(getFromLocalStorage("user", "{}")) as User,
    lastUpdatedAt: getFromLocalStorage("lastUpdatedAt", null),
  }),
  getters: {
    isLoggedIn: (state) => !!state.auth_token,
  },
  actions: {
    async login(login: string, password: string, rememberMe = false) {
      try {
        const response = await axios.post("/auth/login", {
          login,
          password,
        });

        if (!response.data.auth_token || !response.data.refresh_token) {
          throw "MISSING_TOKEN";
        }

        // store jwt in local storage to keep user logged in between page refreshes if remember me is enabled
        const updatedDate = new Date().toISOString();
        if (rememberMe) {
          saveToLocalStorage("auth_token", response.data.auth_token);
          saveToLocalStorage("refresh_token", response.data.refresh_token);
          saveToLocalStorage("lastUpdatedAt", updatedDate);
          saveToLocalStorage("user", JSON.stringify(response.data.user));
        }

        // update pinia state
        this.auth_token = response.data.auth_token;
        this.refresh_token = response.data.refresh_token;
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
        await axios.post("/auth/logout", {
          session_token: this.refresh_token,
        });
      } catch (error) {
        handleError(error);
      }

      this.auth_token = null;
      this.refresh_token = null;
      this.lastUpdatedAt = null;
      this.user = {} as User;
      removeFromLocalStorage("auth_token");
      removeFromLocalStorage("refresh_token");
      removeFromLocalStorage("lastUpdatedAt");
      removeFromLocalStorage("user");

      router.push("/login");
    },
    async refreshToken() {
      try {
        const response = await axios.post("/auth/refresh", {
          refresh_token: this.refresh_token,
        });

        if (!response.data.auth_token || !response.data.refresh_token) {
          throw "MISSING_TOKEN";
        }

        const updatedDate = new Date().toISOString();
        if (getFromLocalStorage("auth_token", null) !== null) {
          saveToLocalStorage("auth_token", response.data.auth_token);
          saveToLocalStorage("refresh_token", response.data.refresh_token);
          saveToLocalStorage("lastUpdatedAt", updatedDate);
        }

        this.auth_token = response.data.auth_token;
        this.refresh_token = response.data.refresh_token;
        this.lastUpdatedAt = updatedDate;
      } catch (error) {
        return handleError(error);
      }
    },
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
