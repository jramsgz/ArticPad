<script setup lang="ts">
import ToolbarNav from "../components/core/ToolbarNav.vue";
</script>
<template>
  <div class="min-h-full bg-gray-800">
    <ToolbarNav />
    <slot></slot>
  </div>
</template>
