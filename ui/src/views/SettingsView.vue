<template>
  <main class="max-w-7xl mx-auto pb-10 lg:py-12 lg:px-8">
    <div class="lg:grid lg:grid-cols-12 lg:gap-x-5">
      <aside class="py-6 px-2 sm:px-6 lg:py-0 lg:px-0 lg:col-span-3">
        <nav class="space-y-1">
          <router-link
            v-for="item in subNavigation"
            :key="item.name"
            :to="item.href"
            :class="[
              item.current
                ? 'bg-gray-900 text-indigo-400 hover:bg-gray-700'
                : 'text-gray-200 hover:text-gray-100 hover:bg-gray-700',
              'group rounded-md px-3 py-2 flex items-center text-sm font-medium',
            ]"
            :aria-current="item.current ? 'page' : undefined"
          >
            <component
              :is="item.icon"
              :class="[
                item.current
                  ? 'text-indigo-500'
                  : 'text-gray-400 group-hover:text-gray-500',
                'flex-shrink-0 -ml-1 mr-3 h-6 w-6',
              ]"
              aria-hidden="true"
            />
            <span class="truncate">
              {{ item.name }}
            </span>
          </router-link>
        </nav>
      </aside>

      <!-- Settings Tabs -->
      <div class="space-y-6 sm:px-6 lg:px-0 lg:col-span-9">
        <component :is="currentTab" />
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import ToolbarNav from "@/components/ToolbarNav.vue";
import {
  CogIcon,
  KeyIcon,
  UserCircleIcon,
  AdjustmentsHorizontalIcon,
} from "@heroicons/vue/24/outline";
import { ProfileTab, AccountTab, AdminTab } from "@/components/settings";
import { useRouter } from "vue-router";
import { ref, markRaw, watchEffect } from "vue";
const router = useRouter();

const subNavigation = [
  {
    name: "Profile",
    href: "/settings/profile",
    icon: UserCircleIcon,
    current: false,
  },
  { name: "Account", href: "/settings/account", icon: CogIcon, current: false },
  {
    name: "Password",
    href: "/settings/password",
    icon: KeyIcon,
    current: false,
  },
  {
    name: "Admin",
    href: "/settings/admin",
    icon: AdjustmentsHorizontalIcon,
    current: false,
  },
];

const currentTab = ref(null);

watchEffect(() => {
  const tab = subNavigation.find(
    (item) => item.href === router.currentRoute.value.path
  );
  subNavigation.forEach((item) => (item.current = false));
  if (tab) {
    tab.current = true;
    changeTab(tab.name + "Tab");
  } else {
    router.push("/settings/profile");
  }
});

// newTab: component name (string)
function changeTab(newTab: string) {
  const lookup: Record<string, any> = {
    ProfileTab,
    AccountTab,
    AdminTab,
  };
  currentTab.value = markRaw(lookup[newTab]);
}
</script>
