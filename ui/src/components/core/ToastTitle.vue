<template>
  <div class="font-bold text-xl">{{ title }}</div>
  <div>{{ body }}</div>
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
});
</script>
