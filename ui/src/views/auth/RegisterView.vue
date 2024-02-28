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
        <form class="space-y-4" @submit="onSubmit">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-300">
              {{ $t("auth.email_address") }}
            </label>
            <div class="mt-1">
              <InputText
                id="email"
                name="email"
                type="email"
                autocomplete="email"
                required
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
              <InputText
                id="username"
                name="username"
                type="text"
                autocomplete="username"
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

          <div>
            <label
              for="confirm_password"
              class="block text-sm font-medium text-gray-300"
            >
              {{ $t("auth.confirm_password") }}
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
import { useRouter } from "vue-router";
import { useForm } from "vee-validate";
import { toTypedSchema } from "@vee-validate/zod";
import { z } from "zod";
import {
  emailSchema,
  usernameSchema,
  passwordSchema,
  requiredStringSchema,
} from "@/utils/validation-schemas";

import FormButton from "@/components/common/forms/FormButton.vue";
import InputText from "@/components/common/forms/InputText.vue";

const authStore = useAuthStore();
const router = useRouter();

const schema = toTypedSchema(
  z
    .object({
      email: emailSchema,
      username: usernameSchema,
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

const { errors, handleSubmit, isSubmitting } = useForm({
  validationSchema: schema,
});

const onSubmit = handleSubmit(async (values) => {
  await authStore
    .register(values.username, values.email, values.password)
    .then(() => {
      router.push("/login");
    })
    .catch((error) => {
      if (error.response.data.error_code === "cannot_send_verification_email") {
        router.push("/login");
      }
    });
});
</script>
