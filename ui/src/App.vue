<script setup lang="ts">
import { computed, onBeforeMount } from "vue";
import { RouterView, useRoute } from "vue-router";
import { useI18n } from "vue-i18n";
import { useAuthStore } from "./stores/auth";
import { localStorageAvailable } from "./utils/localStorage";
import NavbarLayout from "./layouts/NavbarLayout.vue";
import FullPageLayout from "./layouts/FullPageLayout.vue";
import LoadingIndicator from "./components/core/LoadingIndicator.vue";

const route = useRoute();
const i18n = useI18n();
const authStore = useAuthStore();

onBeforeMount(() => {
  if (authStore.user.lang) {
    i18n.locale.value = authStore.user.lang;
  }
});

const layout = computed(() => {
  const layout = route?.meta?.layout || "full-page";
  switch (layout) {
    case "navbar":
      return NavbarLayout;
    case "full-page":
      return FullPageLayout;
    default:
      return FullPageLayout;
  }
});

if (!localStorageAvailable()) {
  alert(i18n.t("errors.local_storage_not_available"));
}
</script>

<template>
  <div class="h-full">
    <LoadingIndicator />
    <component :is="layout">
      <RouterView />
    </component>
  </div>
</template>
