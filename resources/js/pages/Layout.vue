<template>
  <div class="flex min-h-screen flex-col bg-slate-50 text-slate-900 dark:bg-slate-950 dark:text-slate-50">
    <Topbar @toggleSidebar="sidebarOpen = true" @toggleCollapse="toggleCollapsed" />

    <div class="grid w-full flex-1 grid-cols-1 md:grid-cols-[auto_1fr]">
      <Sidebar :open="sidebarOpen" :collapsed="sidebarCollapsed" @close="sidebarOpen = false" />

      <div class="flex min-w-0 flex-col">
        <div v-if="sidebarOpen" class="fixed inset-0 z-20 bg-slate-900/40 md:hidden" @click="sidebarOpen = false"></div>

        <main class="relative z-10 flex-1">
          <div class="mb-6">
            <p class="text-xs font-medium uppercase tracking-wider text-slate-500 dark:text-slate-400">Admin</p>
            <h2 class="mt-1 text-2xl font-semibold tracking-tight">{{ pageTitle }}</h2>
          </div>

          <RouterView />
        </main>

        <Footer />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import { RouterView, useRoute } from "vue-router";

import Footer from "../components/Footer.vue";
import Sidebar from "../components/Sidebar.vue";
import Topbar from "../components/Topbar.vue";

const STORAGE_KEY = "govibe.admin.sidebarCollapsed";

const sidebarOpen = ref(false);
const sidebarCollapsed = ref(false);
const route = useRoute();

try {
  sidebarCollapsed.value = localStorage.getItem(STORAGE_KEY) === "1";
} catch {
  // ignore
}

watch(
  () => route.fullPath,
  () => {
    sidebarOpen.value = false;
  }
);

const pageTitle = computed(() => {
  if (route.path === "/admin" || route.path === "/admin/") return "Dashboard";
  if (route.path.startsWith("/admin/users")) return "User";
  if (route.path.startsWith("/admin/roles")) return "Role";
  if (route.path.startsWith("/admin/settings")) return "Setting";
  return "Dashboard";
});

function toggleCollapsed() {
  sidebarCollapsed.value = !sidebarCollapsed.value;
  try {
    localStorage.setItem(STORAGE_KEY, sidebarCollapsed.value ? "1" : "0");
  } catch {
    // ignore
  }
}
</script>
