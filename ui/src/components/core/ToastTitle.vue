<template>
  <div>
    <div class="font-bold text-xl">{{ title }}</div>
    <div>{{ body }}</div>
    <div v-if="footer" class="text-xs text-gray-100 pt-1">{{ footer }}</div>
  </div>
</template>
<script setup lang="ts">
defineProps({
  title: {
    required: true,
    type: String,
  },
  body: {
    required: true,
    type: String,
  },
  footer: {
    required: false,
    type: String,
    default: "",
  },
});
</script>
