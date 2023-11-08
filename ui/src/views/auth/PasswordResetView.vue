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
        <form v-if="!token" class="space-y-4" @submit="onSubmit">
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
            :disabled="Object.keys(errors).length > 0"
            :loading="isSubmitting"
          />
        </form>
        <form v-else class="space-y-4" @submit="onSubmit">
          <div>
            <label
              for="password"
              class="block text-sm font-medium text-gray-300"
            >
              {{ $t("auth.new_password") }}
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

          <div>
            <label
              for="confirm-password"
              class="block text-sm font-medium text-gray-300"
            >
              {{ $t("auth.confirm_new_password") }}
            </label>
            <div class="mt-1">
              <InputText
                id="confirm_password"
                name="confirm_password"
                type="password"
                required
              />
            </div>
          </div>

          <FormButton
            :text="$t('auth.reset_password')"
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
import { computed } from "vue";
import { useRouter } from "vue-router";
import { useForm } from "vee-validate";
import { toTypedSchema } from "@vee-validate/zod";
import { z } from "zod";
import {
  passwordSchema,
  requiredStringSchema,
} from "@/utils/validation-schemas";

import FormButton from "@/components/common/forms/FormButton.vue";
import InputText from "@/components/common/forms/InputText.vue";

const authStore = useAuthStore();
const router = useRouter();

const token = router.currentRoute.value.params.token;

const requestSchema = toTypedSchema(
  z.object({
    login: requiredStringSchema,
  })
);

const resetPasswordSchema = toTypedSchema(
  z
    .object({
      password: passwordSchema,
      confirm_password: requiredStringSchema,
    })
    .superRefine(({ password, confirm_password }, ctx) => {
      if (password !== confirm_password) {
        ctx.addIssue({
          code: "custom",
          message: "errors.confirm_password_mismatch",
          path: ["confirm_password"],
        });
      }
    })
);

const currentSchema = computed(() =>
  token && typeof token === "string" ? resetPasswordSchema : requestSchema
);
const { errors, handleSubmit, isSubmitting } = useForm({
  validationSchema: currentSchema,
});

const onSubmit = handleSubmit(async (values) => {
  if (token && typeof token === "string") {
    await authStore
      .resetPassword(token, values.password)
      .then(() => {
        router.push("/login");
      })
      .catch(() => {});
  } else {
    await authStore
      .requestPasswordReset(values.login)
      .then(() => {
        router.push("/login");
      })
      .catch(() => {});
  }
});
</script>
