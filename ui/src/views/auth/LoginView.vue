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
        {{ $t("auth.sign_in_account") }}
      </h2>
      <p class="mt-2 text-center text-sm text-gray-400">
        {{ $t("common.or") }}
        <router-link
          to="/register"
          class="font-medium text-indigo-400 hover:text-indigo-500"
        >
          {{ $t("auth.create_new_account").toLocaleLowerCase() }}
        </router-link>
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div class="bg-gray-800 py-8 px-4 shadow sm:rounded-lg sm:px-10">
        <form class="space-y-6" @submit.prevent="handleSubmit">
          <div>
            <label for="login" class="block text-sm font-medium text-gray-300">
              {{ $t("auth.email_address_or_username") }}
            </label>
            <div class="mt-1">
              <input
                v-model="form.login"
                id="login"
                name="login"
                type="text"
                autocomplete="email"
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

          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <input
                v-model="form.remember_me"
                id="remember-me"
                name="remember-me"
                type="checkbox"
                class="h-4 w-4 bg-gray-700 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
              />
              <label for="remember-me" class="ml-2 block text-sm text-gray-100">
                {{ $t("auth.remember_me") }}
              </label>
            </div>

            <div class="text-sm">
              <router-link
                to="/password-reset"
                class="font-medium text-indigo-400 hover:text-indigo-500"
              >
                {{ $t("auth.forgot_password") }}
              </router-link>
            </div>
          </div>

          <FormButton
            :text="$t('auth.sign_in')"
            :disabled="!form.login || !form.password"
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
  login: "",
  password: "",
  remember_me: false,
  loading: false,
});

const handleSubmit = async (e: Event) => {
  e.preventDefault();
  form.loading = true;
  await authStore.login(form.login, form.password);
  form.loading = false;
};
</script>
