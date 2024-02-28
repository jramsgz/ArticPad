<template>
  <header class="bg-gray-900 shadow-sm lg:static lg:overflow-y-visible">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div
        class="relative flex justify-between xl:grid xl:grid-cols-12 lg:gap-8"
      >
        <div
          class="hidden sm:flex md:absolute md:left-0 md:inset-y-0 lg:static xl:col-span-3"
        >
          <div class="flex-shrink-0 flex items-center">
            <router-link to="/">
              <img
                class="block h-8 w-auto lg:hidden"
                src="@/assets/logo.svg"
                alt="ArticPad"
              />
              <img
                class="hidden h-8 w-auto lg:block"
                src="@/assets/logo_horizontal.svg"
                alt="ArticPad"
              />
            </router-link>
          </div>
        </div>
        <div class="min-w-0 flex-1 md:px-8 lg:px-0 xl:col-span-6">
          <div
            class="flex items-center px-6 py-4 md:max-w-3xl md:mx-auto lg:max-w-none lg:mx-0 xl:px-0"
          >
            <div class="w-full">
              <label for="search" class="sr-only">Search</label>
              <div class="relative">
                <div
                  class="pointer-events-none absolute inset-y-0 left-0 pl-3 flex items-center"
                >
                  <MagnifyingGlassIcon
                    class="h-5 w-5 text-gray-200"
                    aria-hidden="true"
                  />
                </div>
                <input
                  id="search"
                  name="search"
                  class="block w-full bg-gray-700 border border-gray-500 rounded-md py-2 pl-10 pr-3 text-sm placeholder-gray-300 focus:outline-none text-gray-300 focus:text-gray-100 focus:placeholder-gray-400 focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  placeholder="Search [Unavailable]"
                  type="search"
                />
              </div>
            </div>
          </div>
        </div>
        <div class="flex items-center justify-end xl:col-span-3">
          <!-- Profile dropdown -->
          <Menu as="div" class="flex-shrink-0 relative">
            <MenuButton
              class="rounded-full flex focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <span class="sr-only">Open user menu</span>
              <UserCircleIcon class="text-white h-8 w-8" />
            </MenuButton>
            <transition
              enter-active-class="transition ease-out duration-100"
              enter-from-class="transform opacity-0 scale-95"
              enter-to-class="transform opacity-100 scale-100"
              leave-active-class="transition ease-in duration-75"
              leave-from-class="transform opacity-100 scale-100"
              leave-to-class="transform opacity-0 scale-95"
            >
              <MenuItems
                class="origin-top-right absolute z-10 right-0 mt-2 w-48 rounded-md shadow-lg bg-gray-700 ring-1 ring-black ring-opacity-5 py-1 focus:outline-none"
              >
                <div class="px-4 py-2 text-sm text-gray-300 truncate">
                  Signed in as <br />
                  <strong class="text-gray-100">
                    {{ authStore.user?.username }}
                  </strong>
                </div>
                <MenuItem
                  v-for="item in userNavigation"
                  v-slot="{ active, close }"
                  :key="item.name"
                  as="template"
                >
                  <router-link
                    :to="item.href"
                    :class="[
                      $route.path === item.href ? 'bg-gray-600' : '',
                      active ? 'bg-gray-600' : '',
                      'block text-sm text-gray-300',
                    ]"
                  >
                    <span class="block py-2 px-4" @click="close">
                      {{ $t(item.name) }}
                    </span>
                  </router-link>
                </MenuItem>
              </MenuItems>
            </transition>
          </Menu>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/vue";
import { UserCircleIcon, MagnifyingGlassIcon } from "@heroicons/vue/24/outline";
import { useAuthStore } from "@/stores/auth";

const authStore = useAuthStore();

const userNavigation = [
  { name: "routes.dashboard", href: "/" },
  { name: "routes.settings", href: "/settings/account" },
  { name: "auth.sign_out", href: "/logout" },
];
</script>
