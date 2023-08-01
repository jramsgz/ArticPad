import type { PluginOptions } from "vue-toastification";

import "vue-toastification/dist/index.css";

export const options: PluginOptions = {
  timeout: 10000,
  closeOnClick: false,
  pauseOnFocusLoss: true,
  pauseOnHover: true,
  draggable: false,
};

// Custom toast component with title and body
import { useToast } from "vue-toastification";
import TitleToast from "@/components/common/TitleToast.vue";

type ToastTypes = "info" | "success" | "warning" | "error";

type ToastTitleInterface = {
  [type in ToastTypes]: (title: string, body: string) => void;
};

export const useToastTitle = () => {
  // Get the original toast interface
  const toast = useToast();

  // Helper method that sets up our toasts with title
  const createToastTitleMethod = <M extends ToastTypes>(
    method: M
  ): ToastTitleInterface[M] => {
    return (title, body) =>
      // If method is "error", timeout is set to 1 minute
      toast[method](
        { component: TitleToast, props: { title, body } },
        { timeout: method === "error" ? 60000 : options.timeout }
      );
  };

  // Create and return our new interface
  const ToastTitleInterface: ToastTitleInterface = {
    info: createToastTitleMethod("info"),
    success: createToastTitleMethod("success"),
    warning: createToastTitleMethod("warning"),
    error: createToastTitleMethod("error"),
  };
  return ToastTitleInterface;
};
