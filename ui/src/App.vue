<script setup lang="ts">
import { computed } from "vue";
import { RouterView, useRoute } from "vue-router";
import NavbarLayout from "./layouts/NavbarLayout.vue";
import FullPageLayout from "./layouts/FullPageLayout.vue";
import LoadingIndicator from "./components/core/LoadingIndicator.vue";

const route = useRoute();

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
</script>

<template>
  <div class="h-full">
    <LoadingIndicator />
    <component :is="layout">
      <RouterView />
    </component>
  </div>
</template>
