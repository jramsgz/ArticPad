<template>
  <div
    class="min-h-full bg-gray-900 flex flex-col justify-center py-12 sm:px-6 lg:px-8"
  >
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <img
        class="mx-auto h-28 w-auto"
        src="@/assets/logo_vertical.svg"
        alt="ArticPad"
      />
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-200">
        {{ $t("auth.create_new_account") }}
      </h2>
      <p class="mt-2 text-center text-sm text-gray-400">
        {{ $t("common.or") }}
        <router-link
          to="/login"
          class="font-medium text-indigo-400 hover:text-indigo-500"
        >
          {{ $t("auth.sign_in_account").toLocaleLowerCase() }}
        </router-link>
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div class="bg-gray-800 py-8 px-4 shadow sm:rounded-lg sm:px-10">
        <form class="space-y-4" @submit.prevent="handleSubmit">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-300">
              {{ $t("auth.email_address") }}
            </label>
            <div class="mt-1">
              <input
                v-model="form.email"
                id="email"
                name="email"
                type="email"
                autocomplete="email"
                required
                class="appearance-none block w-full px-3 py-2 bg-gray-700 border border-gray-500 rounded-md shadow-sm placeholder-gray-300 text-gray-300 focus:text-gray-100 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              />
            </div>
          </div>

          <div>
            <label
              for="username"
              class="block text-sm font-medium text-gray-300"
            >
              {{ $t("auth.username") }}
            </label>
            <div class="mt-1">
              <input
                v-model="form.username"
                id="username"
                name="username"
                type="text"
                autocomplete="username"
                required
                class="appearance-none block w-full px-3 py-2 bg-gray-700 border border-gray-500 rounded-md shadow-sm placeholder-gray-300 text-gray-300 focus:text-gray-100 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              />
            </div>
          </div>

          <div>
            <label
              for="password"
              class="block text-sm font-medium text-gray-300"
            >
              {{ $t("auth.password") }}
            </label>
            <div class="mt-1">
              <PasswordField v-model="form.password" />
            </div>
          </div>

          <div>
            <label
              for="confirm-password"
              class="block text-sm font-medium text-gray-300"
            >
              {{ $t("auth.confirm_password") }}
            </label>
            <div class="mt-1">
              <input
                v-model="form.confirm_password"
                id="confirm-password"
                name="confirm-password"
                type="password"
                autocomplete="confirm-password"
                required
                class="appearance-none block w-full px-3 py-2 bg-gray-700 border border-gray-500 rounded-md shadow-sm placeholder-gray-300 text-gray-300 focus:text-gray-100 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              />
            </div>
          </div>

          <label class="ml-2 block text-sm text-gray-300">
            {{ $t("auth.by_clicking_you_agree_to_our") }}
            <a
              href="dummy-link"
              target="_blank"
              class="font-medium text-indigo-400 hover:text-indigo-500"
            >
              {{ $t("common.privacy_policy") }}
            </a>
          </label>
          <FormButton
            :text="$t('auth.sign_up')"
            :disabled="!form.username"
            :loading="form.loading"
          />
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { reactive } from "vue";

import FormButton from "@/components/common/FormButton.vue";
import PasswordField from "@/components/common/PasswordField.vue";

const authStore = useAuthStore();

const form = reactive({
  username: "",
  email: "",
  password: "",
  confirm_password: "",
  loading: false,
});

const handleSubmit = async (e: Event) => {
  e.preventDefault();
  form.loading = true;
  await authStore.register(form.username, form.email, form.password);
  form.loading = false;
};
</script>
