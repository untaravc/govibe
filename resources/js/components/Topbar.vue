<template>
  <header class="sticky top-0 z-20 border-b border-slate-200/60 bg-white/80 backdrop-blur dark:border-white/10 dark:bg-slate-950/70">
    <div class="flex w-full items-center justify-between gap-3 px-4 py-3">
      <div class="flex items-center gap-3">
        <button
          type="button"
          class="inline-flex items-center justify-center rounded-lg p-2 text-slate-600 hover:bg-slate-100 hover:text-slate-900 dark:text-slate-300 dark:hover:bg-white/10 dark:hover:text-white md:hidden"
          @click="$emit('toggleSidebar')"
        >
          <span class="sr-only">Toggle sidebar</span>
          <svg viewBox="0 0 24 24" fill="none" class="h-5 w-5" xmlns="http://www.w3.org/2000/svg">
            <path
              d="M4 6h16M4 12h16M4 18h16"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
        </button>

        <button
          type="button"
          class="hidden items-center justify-center rounded-lg p-2 text-slate-600 hover:bg-slate-100 hover:text-slate-900 dark:text-slate-300 dark:hover:bg-white/10 dark:hover:text-white md:inline-flex"
          @click="$emit('toggleCollapse')"
        >
          <span class="sr-only">Collapse sidebar</span>
          <svg viewBox="0 0 24 24" fill="none" class="h-5 w-5" xmlns="http://www.w3.org/2000/svg">
            <path
              d="M10 6H6v12h4M14 6h4v12h-4"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </button>

        <div>
          <p class="text-sm font-semibold text-slate-900 dark:text-slate-50">Administrator</p>
          <p class="text-xs text-slate-500 dark:text-slate-400">GoVibe dashboard</p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <div class="hidden w-72 md:block">
          <div class="relative">
            <input
              type="search"
              placeholder="Search..."
              class="w-full rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm text-slate-900 shadow-sm outline-none ring-slate-900/10 placeholder:text-slate-400 focus:border-slate-300 focus:ring-4 dark:border-white/10 dark:bg-slate-950 dark:text-slate-50 dark:placeholder:text-slate-500"
            />
            <span class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-slate-400">
              <svg viewBox="0 0 24 24" fill="none" class="h-4 w-4" xmlns="http://www.w3.org/2000/svg">
                <path
                  d="M11 19a8 8 0 1 1 0-16 8 8 0 0 1 0 16Zm10 2-4.3-4.3"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                />
              </svg>
            </span>
          </div>
        </div>

        <RouterLink
          to="/"
          class="rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-white/10"
        >
          Public site
        </RouterLink>

        <button
          type="button"
          class="rounded-lg bg-slate-900 px-3 py-2 text-sm font-medium text-white hover:bg-slate-800 dark:bg-slate-50 dark:text-slate-900 dark:hover:bg-white"
          @click="logout"
        >
          Logout
        </button>
      </div>
    </div>
  </header>
</template>

<script setup>
import { RouterLink, useRouter } from "vue-router";

import { useAuthStore } from "../store/auth.js";

defineEmits(["toggleSidebar", "toggleCollapse"]);

const router = useRouter();
const auth = useAuthStore();

function logout() {
  auth.logout();
  try {
    localStorage.removeItem("token");
    sessionStorage.removeItem("token");
  } catch {
    // ignore
  }
  router.push("/auth/login");
}
</script>
