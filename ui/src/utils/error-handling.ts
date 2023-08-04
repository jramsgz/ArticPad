import { useAuthStore } from "@/stores/auth";
import { useToastWithTitle } from "@/plugins/toast";
import i18n from "@/plugins/i18n";

export const handleError = (err: any) => {
  if (err.code === "ERR_NETWORK") {
    showError("NET_ERR");
  } else if (err.response && err.response.data) {
    if (
      err.response.statusText === "Unauthorized" ||
      err.response.data === "Unauthorized"
    ) {
      // Unauthorized and log out
      const msg = err.response.data.error
        ? err.response.data.error
        : "LOG_IN_AGAIN";

      showError(msg, err.response.data.requestId);

      const authStore = useAuthStore();
      authStore.logout();
    } else if (err.response.data.errors) {
      // Show a notification per error
      const errors = JSON.parse(JSON.stringify(err.response.data.errors));
      for (const i in errors) {
        showError(errors[i][0], err.response.data.requestId);
      }
    } else if (err.response.data.error) {
      showError(err.response.data.error, err.response.data.requestId);
    } else {
      showError(err.response.data, err.response.data.requestId);
    }
  } else {
    // Make sure error is a string
    const error = err.toString();

    showError(error);
  }
  Promise.reject(err);
};

export const showError = (error: string, footer?: string) => {
  switch (error) {
    /* Internal frontend-defined errors */
    case "NET_ERR":
      showToast("errors.net_err", "errors.net_err_msg");
      break;

    case "LOG_IN_AGAIN":
      showToast("errors.unauthorized", "errors.log_in_again", footer);
      break;

    case "MISSING_TOKEN":
      showToast("errors.missing_token", "errors.missing_token_msg");
      break;

    /* External backend-defined errors */
    case "user, email or password is incorrect":
      showToast("errors.login_failed", "errors.invalid_credentials");
      break;

    case "username must be at least 3 characters":
      showToast("errors.signup_failed", "errors.username_too_short");
      break;

    case "username must be at most 32 characters":
      showToast("errors.signup_failed", "errors.username_too_long");
      break;

    case "username must only contain letters, numbers, dashes, underscores and dots":
      showToast("errors.signup_failed", "errors.username_invalid");
      break;

    case "password must be at least 8 characters":
      showToast("errors.signup_failed", "errors.password_too_short");
      break;

    case "password must be at most 64 characters":
      showToast("errors.signup_failed", "errors.password_too_long");
      break;

    case "password must not be too similar to username or email":
      showToast("errors.signup_failed", "errors.password_too_similar");
      break;

    case "invalid email":
      showToast("errors.signup_failed", "errors.email_invalid");
      break;

    case "this email is already in use":
      showToast("errors.signup_failed", "errors.email_in_use");
      break;

    case "username already exists":
      showToast("errors.signup_failed", "errors.username_in_use");
      break;

    case "password must contain at least one uppercase letter, one lowercase letter, one number and one special character":
      showToast("errors.signup_failed", "errors.password_too_weak");
      break;

    default:
      showToast("errors.unknown", error, footer);
      break;
  }
};

export const showToast = (title: string, msg: string, footer?: string) => {
  const toast = useToastWithTitle();
  title = i18n.global.t(title);
  msg = i18n.global.t(msg);
  footer = footer ? "Request ID: " + footer : undefined;

  toast.error(title, msg, footer);
};
