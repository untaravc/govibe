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
        <div v-for="item in menuTree" :key="item.id">
          <component
            :is="item.link ? RouterLink : 'button'"
            :to="item.link || undefined"
            type="button"
            :title="collapsed ? item.name : undefined"
            class="group flex w-full items-center gap-3 rounded-xl px-3 py-2 text-left text-sm font-medium hover:bg-slate-100 dark:hover:bg-white/10"
            :class="linkClasses(item)"
            :active-class="item.link ? activeClasses : undefined"
            :exact-active-class="item.link ? exactActiveClasses : undefined"
            @click="item.link ? null : toggle(item)"
          >
            <span
              class="inline-flex h-9 w-9 items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-700 shadow-sm dark:border-white/10 dark:bg-slate-900 dark:text-slate-200"
            >
              <Icon :icon="resolveIcon(item.icon)" class="h-5 w-5" />
            </span>

            <span v-if="!collapsed" class="min-w-0 flex-1 truncate">{{ item.name }}</span>

            <span v-if="!collapsed && hasChildren(item)" class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-slate-500">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                class="h-4 w-4 transition-transform"
                :class="isExpanded(item.id) ? 'rotate-90' : ''"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path d="M9 6l6 6-6 6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
              </svg>
            </span>
          </component>

          <div v-if="!collapsed && hasChildren(item) && isExpanded(item.id)" class="mt-1 space-y-1 pl-11">
            <RouterLink
              v-for="child in item.children"
              :key="child.id"
              :to="child.link || '#'"
              class="group flex items-center gap-3 rounded-xl px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-white/10"
              active-class="bg-slate-100 text-slate-900 dark:bg-white/10 dark:text-white"
              exact-active-class="bg-slate-100 text-slate-900 dark:bg-white/10 dark:text-white"
            >
              <span class="inline-flex h-8 w-8 items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-700 shadow-sm dark:border-white/10 dark:bg-slate-900 dark:text-slate-200">
                <Icon :icon="resolveIcon(child.icon)" class="h-4 w-4" />
              </span>
              <span class="min-w-0 flex-1 truncate">{{ child.name }}</span>
            </RouterLink>
          </div>
        </div>
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
import { computed, onMounted, ref, watch } from "vue";
import { RouterLink, useRoute } from "vue-router";

import api from "../api.js";

import accountMultipleOutline from "@iconify/icons-mdi/account-multiple-outline";
import circleOutline from "@iconify/icons-mdi/circle-outline";
import cogOutline from "@iconify/icons-mdi/cog-outline";
import fileDocumentEditOutline from "@iconify/icons-mdi/file-document-edit-outline";
import shieldAccountOutline from "@iconify/icons-mdi/shield-account-outline";
import shieldKeyOutline from "@iconify/icons-mdi/shield-key-outline";
import viewDashboardOutline from "@iconify/icons-mdi/view-dashboard-outline";

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

const route = useRoute();

const activeClasses = "bg-slate-100 text-slate-900 dark:bg-white/10 dark:text-white";
const exactActiveClasses = "bg-slate-100 text-slate-900 dark:bg-white/10 dark:text-white";

const iconMap = {
  "mdi:view-dashboard-outline": viewDashboardOutline,
  "mdi:cog-outline": cogOutline,
  "mdi:account-multiple-outline": accountMultipleOutline,
  "mdi:shield-account-outline": shieldAccountOutline,
  "mdi:shield-key-outline": shieldKeyOutline,
  "mdi:file-document-edit-outline": fileDocumentEditOutline,
  "mdi:circle-outline": circleOutline
};

const menuTree = ref([]);
const expandedById = ref({});

function resolveIcon(icon) {
  if (icon && typeof icon === "string" && iconMap[icon]) return iconMap[icon];
  return circleOutline;
}

function hasChildren(item) {
  return Array.isArray(item?.children) && item.children.length > 0;
}

function isExpanded(id) {
  return Boolean(expandedById.value[id]);
}

function toggle(item) {
  expandedById.value = { ...expandedById.value, [item.id]: !expandedById.value[item.id] };
}

function isLinkActive(link) {
  if (!link) return false;
  const current = route.path;
  if (link === "/admin") return current === "/admin" || current === "/admin/";
  return current === link || current.startsWith(`${link}/`);
}

function hasActiveDescendant(item) {
  if (isLinkActive(item.link)) return true;
  if (!hasChildren(item)) return false;
  return item.children.some((c) => hasActiveDescendant(c));
}

function linkClasses(item) {
  const active = !item.link && hasActiveDescendant(item);
  return active ? activeClasses : "text-slate-700 dark:text-slate-200";
}

function expandToActive(tree) {
  const next = { ...expandedById.value };
  for (const item of tree) {
    if (hasChildren(item) && hasActiveDescendant(item)) next[item.id] = true;
  }
  expandedById.value = next;
}

async function loadMenus() {
  try {
    const { res, json } = await api.get("/api/menu", { auth: true });
    if (!res.ok) return;
    menuTree.value = Array.isArray(json?.result?.menus) ? json.result.menus : [];
    expandToActive(menuTree.value);
  } catch {
    // ignore
  }
}

watch(
  () => route.path,
  () => {
    if (menuTree.value.length > 0) expandToActive(menuTree.value);
  }
);

onMounted(() => {
  loadMenus();
});
</script>
