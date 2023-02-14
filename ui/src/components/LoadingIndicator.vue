<script setup lang="ts">
import { useLoadingStore } from "../stores/loading";
const loadingStore = useLoadingStore();
</script>

<template>
  <div
    v-if="loadingStore.isRouteLoading"
    class="fixed overflow-hidden top-0 z-50 w-full bg-indigo-800 h-1"
  >
    <div class="bg-indigo-400 h-1 loading"></div>
  </div>
</template>

<style scoped>
.loading {
  animation: loading 1s infinite;
}

@keyframes loading {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}
</style>
