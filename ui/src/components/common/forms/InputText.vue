<template>
  <div class="relative">
    <input
      v-model="value"
      :id="name"
      :name="name"
      :type="
        type === 'password'
          ? showPassword
            ? 'text'
            : 'password'
          : type || 'text'
      "
      :autocomplete="autocomplete"
      :required="required || false"
      :class="[
        'appearance-none block w-full px-3 py-2 bg-gray-700 border border-gray-500 rounded-md shadow-sm placeholder-gray-300 text-gray-300 focus:text-gray-100 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm',
        { 'pr-10': type === 'password' },
        inputClass,
        { 'border-red-500': errorMessage },
      ]"
    />
    <button
      v-if="type === 'password'"
      class="absolute inset-y-0 right-0 flex items-center px-4 text-white"
      type="button"
      :aria-label="$t('common.toggle_password_visibility')"
      @click="showPassword = !showPassword"
    >
      <EyeSlashIcon v-if="showPassword" class="w-5 h-5" aria-hidden="true" />
      <EyeIcon v-else class="w-5 h-5" aria-hidden="true" />
    </button>
  </div>
  <div
    v-if="errorMessage"
    :id="`${name}-error`"
    :class="['flex items-center mt-2 text-sm text-red-500', errorClass]"
  >
    <ExclamationCircleIcon class="h-5 w-5 mr-1" aria-hidden="true" />
    <span class="sr-only">{{ $t("errors.error") }}:</span>
    {{ $t(errorMessage) }}
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useField } from "vee-validate";
import {
  EyeIcon,
  EyeSlashIcon,
  ExclamationCircleIcon,
} from "@heroicons/vue/24/outline";

const props = defineProps({
  name: {
    type: String,
    required: true,
  },
  type: String,
  inputClass: String,
  errorClass: String,
  id: String,
  autocomplete: String,
  required: Boolean,
});

const { value, errorMessage } = useField(() => props.name);
const showPassword = ref(false);
</script>
