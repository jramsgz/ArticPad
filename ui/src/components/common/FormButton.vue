<template>
  <div>
    <button
      type="submit"
      :disabled="disabled || loading"
      :class="{
        'bg-indigo-500 hover:bg-indigo-400 cursor-not-allowed':
          disabled || loading,
        'bg-indigo-600 hover:bg-indigo-700': !disabled && !loading,
      }"
      class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white transition ease-in-out duration-150 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
    >
      <svg
        v-if="loading"
        class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
      >
        <circle
          class="opacity-25"
          cx="12"
          cy="12"
          r="10"
          stroke="currentColor"
          stroke-width="4"
        ></circle>
        <path
          class="opacity-75"
          fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
        ></path>
      </svg>
      {{ text }}
    </button>
  </div>
</template>

<script lang="ts">
export default {
  props: {
    disabled: {
      required: false,
      type: Boolean,
      default: false,
    },
    loading: {
      required: false,
      type: Boolean,
      default: false,
    },
    text: {
      required: false,
      type: String,
      default: "Submit",
    },
  },
};
</script>
