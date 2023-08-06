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
        {{ $t("auth.reset_your_password") }}
      </h2>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div class="bg-gray-800 py-8 px-4 shadow sm:rounded-lg sm:px-10">
        <form v-if="!token" class="space-y-4" @submit.prevent="handleSubmit">
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

          <div class="flex items-center justify-between">
            <div class="text-sm">
              <router-link
                to="/login"
                class="font-medium text-indigo-400 hover:text-indigo-500"
              >
                {{ $t("auth.go_back_to_login") }}
              </router-link>
            </div>
          </div>

          <FormButton
            :text="$t('auth.reset_password')"
            :disabled="!form.login"
            :loading="form.loading"
          />
        </form>
        <form v-else class="space-y-4" @submit.prevent="handleSubmit">
          <div>
            <label
              for="password"
              class="block text-sm font-medium text-gray-300"
            >
              {{ $t("auth.new_password") }}
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
              {{ $t("auth.confirm_new_password") }}
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

          <FormButton
            :text="$t('auth.reset_password')"
            :disabled="!form.password || !form.confirm_password"
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
import { useRouter } from "vue-router";

import FormButton from "@/components/common/FormButton.vue";
import PasswordField from "@/components/common/PasswordField.vue";

const authStore = useAuthStore();
const router = useRouter();

const token = router.currentRoute.value.params.token;

const form = reactive({
  login: "",
  password: "",
  confirm_password: "",
  showPassword: false,
  loading: false,
});

const handleSubmit = async (e: Event) => {
  e.preventDefault();
  form.loading = true;
  if (token && typeof token === "string") {
    await authStore.resetPassword(token, form.password);
    router.push("/login");
  } else {
    await authStore.requestPasswordReset(form.login);
  }
  form.loading = false;
};
</script>
