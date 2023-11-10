<template>
  <section aria-labelledby="account-details-heading">
    <form class="space-y-6" @submit="onSubmit">
      <div class="shadow sm:rounded-md sm:overflow-hidden">
        <div class="bg-gray-900 py-6 px-4 sm:p-6">
          <div>
            <h2
              id="account-details-heading"
              class="text-lg leading-6 font-medium text-gray-100"
            >
              Account details
            </h2>
            <p class="mt-1 text-sm text-gray-500">
              Update your account information. Please note that changing your
              email address will also change the email you use to log into
              ArticPad.
            </p>
          </div>

          <div class="mt-6 grid grid-cols-4 gap-6">
            <div class="col-span-4 sm:col-span-2">
              <label
                for="username"
                class="block text-sm font-medium text-gray-300"
              >
                Username
              </label>
              <InputText
                id="username"
                name="username"
                type="text"
                autocomplete="username"
                required
                input-class="mt-1"
              />
            </div>

            <div class="col-span-4 sm:col-span-2">
              <label
                for="username"
                class="block text-sm font-medium text-gray-300"
              >
                Display name
              </label>
              <input
                id="username"
                name="username"
                type="text"
                autocomplete="username"
                required
                class="mt-1 appearance-none block w-full px-3 py-2 bg-gray-700 border border-gray-500 rounded-md shadow-sm placeholder-gray-300 text-gray-300 focus:ring-1 focus:text-gray-100 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              />
            </div>

            <div class="col-span-4 sm:col-span-2">
              <label
                for="email"
                class="block text-sm font-medium text-gray-300"
              >
                Email
              </label>
              <InputText
                id="email"
                name="email"
                type="email"
                autocomplete="email"
                required
                input-class="mt-1"
              />
            </div>

            <div class="col-span-4 sm:col-span-2">
              <label
                for="language"
                class="block text-sm font-medium text-gray-300"
                >Language</label
              >
              <select
                v-model="language"
                id="language"
                name="language"
                autocomplete="language"
                required
                class="mt-1 appearance-none block w-full px-3 py-2 bg-gray-700 border border-gray-500 rounded-md shadow-sm placeholder-gray-300 text-gray-300 focus:ring-1 focus:text-gray-100 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              >
                <option
                  v-for="(name, code) in availableLocales"
                  :key="`locale-${code}`"
                  :value="code"
                >
                  {{ name }}
                </option>
              </select>
            </div>
          </div>
        </div>
        <div class="px-4 py-3 bg-gray-900 text-right sm:px-6">
          <FormButton
            :text="$t('auth.sign_in')"
            :disabled="Object.keys(errors).length > 0"
            :loading="isSubmitting"
          />
          {{ errors }}
          <button
            type="submit"
            class="bg-gray-800 border border-transparent rounded-md shadow-sm py-2 px-4 inline-flex justify-center text-sm font-medium text-white hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-900"
          >
            Save
          </button>
        </div>
      </div>
    </form>
  </section>
</template>

<script setup lang="ts">
import { watchEffect } from "vue";
import { useI18n } from "vue-i18n";
import { useAuthStore } from "@/stores/auth";
import { useForm, useField } from "vee-validate";
import { toTypedSchema } from "@vee-validate/zod";
import { z } from "zod";
import {
  emailSchema,
  usernameSchema,
  requiredStringSchema,
} from "@/utils/validation-schemas";

import FormButton from "@/components/common/forms/FormButton.vue";
import InputText from "@/components/common/forms/InputText.vue";

// Key is the locale code, value is the name of the language in that locale.
// Note: The key is the same as the exported locale in the locales folder as well as the
// locale code in the backend.
const availableLocales = {
  en: "English",
  es: "Español",
};

const i18n = useI18n();
const authStore = useAuthStore();

const schema = toTypedSchema(
  z.object({
    username: usernameSchema,
    email: emailSchema,
    language: requiredStringSchema.default(
      authStore.user.lang || i18n.locale.value
    ),
  })
);

const { errors, handleSubmit, isSubmitting } = useForm({
  validationSchema: schema,
});

const { value: language } = useField("language");

watchEffect(() => {
  i18n.locale.value = language.value as string;
});

const onSubmit = handleSubmit(async (values) => {
  console.log(values);
  await authStore.login("", "").catch(() => {});
});
</script>
