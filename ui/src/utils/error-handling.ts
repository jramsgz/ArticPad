import { useAuthStore } from "@/stores/auth";
import { useToastWithTitle } from "@/plugins/toast";
import i18n from "@/plugins/i18n";
import errorMap from "@/utils/error-map";

export const handleError = (err: any) => {
  console.error(err);
  if (err.code === "ERR_NETWORK") {
    showError("NET_ERR");
  } else if (err.response && err.response.data) {
    const msg: string | undefined = err.response.data.error
      ? err.response.data.error
      : undefined;
    if (
      err.response.statusText === "Unauthorized" ||
      err.response.data === "Unauthorized"
    ) {
      const errCode: string = err.response.data.error_code
        ? err.response.data.error_code
        : "LOG_IN_AGAIN";

      showError(errCode, msg, err.response.data.requestId);
      if (window.location.pathname !== "/login") {
        // Unauthorized and log out
        const authStore = useAuthStore();
        authStore.logout();
      }
    } else if (err.response.data.errors) {
      // Show a notification per error
      const errors = JSON.parse(JSON.stringify(err.response.data.errors));
      for (const i in errors) {
        showError(errors[i][0], errors[i][0], err.response.data.requestId);
      }
    } else if (err.response.data.error_code) {
      showError(err.response.data.error_code, msg, err.response.data.requestId);
    } else if (err.response.data.error) {
      showError(err.response.data.error, msg, err.response.data.requestId);
    } else {
      showError(err.response.data, msg, err.response.data.requestId);
    }
  } else {
    // Make sure error is a string
    const error = err.toString();

    showError(error);
  }
  return Promise.reject(err);
};

export const showError = (
  errorCode: string,
  error?: string,
  footer?: string
) => {
  const foundError = errorMap[errorCode];
  if (foundError) {
    if (error === undefined) {
      error = foundError.message;
    }
    showToast(
      foundError.title,
      foundError.setMessage ? error : foundError.message,
      foundError.setFooter ? footer : undefined
    );
    return;
  }

  // If error isn't found, try to find it using the error message
  switch (error) {
    default:
      showToast("errors.unknown", error ? error : errorCode, footer);
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
