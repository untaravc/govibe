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
            <component :is="item.icon" class="h-5 w-5" />
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

const IconDashboard = {
  template: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M4 13.5V20a1 1 0 0 0 1 1h5v-7.5H4Zm10 7.5h5a1 1 0 0 0 1-1v-10h-6v11ZM4 10h6V4H5a1 1 0 0 0-1 1v5Zm10-6v3.5h6V5a1 1 0 0 0-1-1h-5Z" fill="currentColor"/></svg>`
};
const IconUsers = {
  template: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M16 11a4 4 0 1 0-8 0 4 4 0 0 0 8 0ZM4 20a6 6 0 0 1 12 0v1H4v-1Zm14.5-4.5a3 3 0 0 0-1.1.2 7.9 7.9 0 0 1 2.6 5.3h2v-.8a4.7 4.7 0 0 0-3.5-4.7Z" fill="currentColor"/></svg>`
};
const IconShield = {
  template: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M12 2 4 5.5V11c0 5.2 3.4 9.9 8 11 4.6-1.1 8-5.8 8-11V5.5L12 2Zm0 18c-3.3-1-6-4.8-6-9V6.7l6-2.6 6 2.6V11c0 4.2-2.7 8-6 9Z" fill="currentColor"/></svg>`
};
const IconSettings = {
  template: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M19.4 13a7.7 7.7 0 0 0 .1-1 7.7 7.7 0 0 0-.1-1l2-1.6a.7.7 0 0 0 .2-.9l-1.9-3.3a.7.7 0 0 0-.9-.3l-2.3.9a7.1 7.1 0 0 0-1.7-1l-.3-2.5a.7.7 0 0 0-.7-.6h-3.8a.7.7 0 0 0-.7.6l-.3 2.5a7.1 7.1 0 0 0-1.7 1l-2.3-.9a.7.7 0 0 0-.9.3L2.4 7.5a.7.7 0 0 0 .2.9l2 1.6a7.7 7.7 0 0 0-.1 1 7.7 7.7 0 0 0 .1 1l-2 1.6a.7.7 0 0 0-.2.9l1.9 3.3c.2.3.6.4.9.3l2.3-.9c.5.4 1.1.7 1.7 1l.3 2.5c0 .3.3.6.7.6h3.8c.3 0 .7-.3.7-.6l.3-2.5c.6-.3 1.2-.6 1.7-1l2.3.9c.3.1.7 0 .9-.3l1.9-3.3a.7.7 0 0 0-.2-.9l-2-1.6ZM12 15.5A3.5 3.5 0 1 1 12 8a3.5 3.5 0 0 1 0 7.5Z" fill="currentColor"/></svg>`
};

const menu = [
  { label: "Dashboard", to: "/admin", icon: IconDashboard },
  { label: "User", to: "/admin/users", icon: IconUsers },
  { label: "Role", to: "/admin/roles", icon: IconShield },
  { label: "Setting", to: "/admin/settings", icon: IconSettings }
];
</script>
