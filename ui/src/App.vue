<script setup lang="ts">
import { computed } from "vue";
import { RouterView, useRoute } from "vue-router";
import NavbarLayout from "./layouts/NavbarLayout.vue";
import FullPageLayout from "./layouts/FullPageLayout.vue";

const route = useRoute();

const showNavbar = computed(() => {
  const showNavbar = route?.meta?.showNavbar;
  if (showNavbar) {
    return NavbarLayout;
  }
  return FullPageLayout;
});
</script>

<template>
  <div class="h-full">
    <component :is="showNavbar">
      <RouterView />
    </component>
  </div>
</template>
