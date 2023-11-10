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
        <form class="space-y-6" @submit="onSubmit">
          <div>
            <label for="login" class="block text-sm font-medium text-gray-300">
              {{ $t("auth.email_address_or_username") }}
            </label>
            <div class="mt-1">
              <InputText
                id="login"
                name="login"
                type="text"
                autocomplete="login"
                required
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
              <InputText
                id="password"
                name="password"
                type="password"
                autocomplete="password"
                required
              />
            </div>
          </div>

          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <input
                v-model="remember_me"
                id="remember_me"
                name="remember_me"
                type="checkbox"
                class="h-4 w-4 bg-gray-700 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
              />
              <label for="remember_me" class="ml-2 block text-sm text-gray-100">
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

          <div
            v-if="values.login === lastUnverifiedUser && values.login !== ''"
            class="text-sm flex items-center justify-center"
          >
            <button
              :class="
                'font-medium text-indigo-400 hover:text-indigo-500' +
                (isSubmitting ? ' animate-pulse cursor-wait' : 'cursor-pointer')
              "
              @click="handleResendVerificationEmail"
              :disabled="isSubmitting"
            >
              {{ $t("auth.resend_verification_email") }}
            </button>
          </div>
          <FormButton
            :text="$t('auth.sign_in')"
            :disabled="Object.keys(errors).length > 0"
            :loading="isSubmitting"
          />
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useForm, useField } from "vee-validate";
import { toTypedSchema } from "@vee-validate/zod";
import { z } from "zod";
import { requiredStringSchema } from "@/utils/validation-schemas";

import FormButton from "@/components/common/forms/FormButton.vue";
import InputText from "@/components/common/forms/InputText.vue";

const authStore = useAuthStore();
const router = useRouter();

const schema = toTypedSchema(
  z.object({
    login: requiredStringSchema,
    password: requiredStringSchema,
    remember_me: z.boolean().optional().default(true),
  })
);

const { errors, handleSubmit, isSubmitting, values } = useForm({
  validationSchema: schema,
});
const { value: remember_me } = useField("remember_me");

let sentLogin = "";
const lastUnverifiedUser = ref("");

const onSubmit = handleSubmit(async (values) => {
  sentLogin = values.login;
  await authStore
    .login(values.login, values.password, values.remember_me)
    .then(() => {
      // redirect to previous url or default to home page
      const returnUrl = router.currentRoute.value.query.redirect as string;
      // Dont redirect to logout page
      if (returnUrl === "/logout") {
        router.push("/");
        return;
      }

      router.push(returnUrl || "/");
    })
    .catch((err) => {
      if (err?.response?.data?.error_code === "email_not_verified") {
        lastUnverifiedUser.value = sentLogin;
      }
    });
});

const handleResendVerificationEmail = async () => {
  isSubmitting.value = true;
  await authStore
    .resendVerificationEmail(sentLogin)
    .then(() => {
      lastUnverifiedUser.value = "";
    })
    .catch(() => {});
  isSubmitting.value = false;
};
</script>
