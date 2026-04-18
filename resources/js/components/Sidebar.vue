<template>
  <aside
    class="fixed inset-y-0 left-0 z-30 -translate-x-full border-r border-slate-200/60 bg-white transition-[transform,width] duration-200 dark:border-white/10 dark:bg-slate-950 md:static md:translate-x-0"
    :class="[
      open ? 'translate-x-0' : '',
      collapsed ? 'w-20' : 'w-72'
    ]"
  >
    <div class="flex h-full flex-col">
      <div class="flex items-center justify-between px-4 py-4">
        <p v-if="!collapsed" class="text-sm font-semibold tracking-tight text-slate-900 dark:text-slate-50">GoVibe Admin</p>
        <p v-else class="text-sm font-semibold tracking-tight text-slate-900 dark:text-slate-50">GV</p>
        <button
          type="button"
          class="rounded-lg p-2 text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-white/10 md:hidden"
          @click="$emit('close')"
        >
          <span class="sr-only">Close sidebar</span>
          <svg viewBox="0 0 24 24" fill="none" class="h-5 w-5" xmlns="http://www.w3.org/2000/svg">
            <path d="M6 6l12 12M18 6L6 18" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
          </svg>
        </button>
      </div>

      <nav class="flex-1 space-y-1 px-3 py-2">
        <RouterLink
          v-for="item in menu"
          :key="item.to"
          :to="item.to"
          :title="collapsed ? item.label : undefined"
          class="group flex items-center gap-3 rounded-xl px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-white/10"
          active-class="bg-slate-100 text-slate-900 dark:bg-white/10 dark:text-white"
          exact-active-class="bg-slate-100 text-slate-900 dark:bg-white/10 dark:text-white"
        >
          <span
            class="inline-flex h-9 w-9 items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-700 shadow-sm dark:border-white/10 dark:bg-slate-900 dark:text-slate-200"
          >
            <Icon :icon="item.icon" class="h-5 w-5" />
          </span>
          <span v-if="!collapsed" class="truncate">{{ item.label }}</span>
        </RouterLink>
      </nav>

      <div
        class="border-t border-slate-200/60 px-4 py-4 text-xs text-slate-500 dark:border-white/10 dark:text-slate-400"
        :class="collapsed ? 'text-center' : ''"
      >
        <span v-if="!collapsed">v0.1 · Admin area</span>
        <span v-else>v0.1</span>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { RouterLink } from "vue-router";

import viewDashboardOutline from "@iconify/icons-mdi/view-dashboard-outline";
import accountMultipleOutline from "@iconify/icons-mdi/account-multiple-outline";
import shieldAccountOutline from "@iconify/icons-mdi/shield-account-outline";
import cogOutline from "@iconify/icons-mdi/cog-outline";

defineProps({
  open: {
    type: Boolean,
    default: false
  },
  collapsed: {
    type: Boolean,
    default: false
  }
});

defineEmits(["close"]);

const menu = [
  { label: "Dashboard", to: "/admin", icon: viewDashboardOutline },
  { label: "User", to: "/admin/users", icon: accountMultipleOutline },
  { label: "Role", to: "/admin/roles", icon: shieldAccountOutline },
  { label: "Setting", to: "/admin/settings", icon: cogOutline }
];
</script>
