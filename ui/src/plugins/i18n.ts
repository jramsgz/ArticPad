import { createI18n } from "vue-i18n";
import messages from "@/locales";

const i18n = createI18n({
  legacy: false,
  locale: navigator.language.split("-")[0] || "en",
  fallbackLocale: "en",
  messages,
});

export default i18n;
