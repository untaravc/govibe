<template>
  <aside
    class="fixed inset-y-0 left-0 z-30 -translate-x-full border-r border-blue-900/40 bg-blue-950 text-white transition-[transform,width] duration-200 md:static md:translate-x-0"
    :class="[
      open ? 'translate-x-0' : '',
      collapsed ? 'w-20' : 'w-72'
    ]"
  >
    <div class="flex h-full flex-col">
      <div class="flex items-center justify-between px-4 py-4">
        <p v-if="!collapsed" class="text-sm font-semibold tracking-tight text-white">GoVibe Admin</p>
        <p v-else class="text-sm font-semibold tracking-tight text-white">GV</p>
        <button
          type="button"
          class="rounded-lg p-2 text-slate-200 hover:bg-white/10 hover:text-white md:hidden"
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
            class="group flex w-full items-center gap-3 rounded-xl px-3 py-2 text-left text-sm font-medium text-slate-200 hover:bg-white/10 hover:text-white"
            :class="linkClasses(item)"
            :active-class="item.link ? activeClasses : undefined"
            :exact-active-class="item.link ? exactActiveClasses : undefined"
            @click="item.link ? null : toggle(item)"
          >
            <span
              class="inline-flex h-9 w-9 items-center justify-center rounded-xl border border-white/10 bg-white/5 text-sky-200 shadow-sm transition group-hover:border-sky-300/30 group-hover:bg-sky-400/10 group-hover:text-sky-100"
            >
              <Icon :icon="resolveIcon(item.icon)" class="h-5 w-5" />
            </span>

            <span v-if="!collapsed" class="min-w-0 flex-1 truncate">{{ item.name }}</span>

            <span v-if="!collapsed && hasChildren(item)" class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-slate-300/80">
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
              class="group flex items-center rounded-xl px-3 py-2 text-sm font-medium text-slate-200 hover:bg-white/10 hover:text-white"
              :class="childLinkClasses(child)"
            >
              <span class="min-w-0 flex-1 truncate">{{ child.name }}</span>
            </RouterLink>
          </div>
        </div>
      </nav>

      <div
        class="border-t border-white/10 px-4 py-4 text-xs text-slate-300/70"
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
import truckCheckOutline from "@iconify/icons-mdi/truck-check-outline";
import truckDeliveryOutline from "@iconify/icons-mdi/truck-delivery-outline";
import truckFastOutline from "@iconify/icons-mdi/truck-fast-outline";
import truckOutline from "@iconify/icons-mdi/truck-outline";
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

const activeClasses = "bg-white/10 text-white ring-1 ring-sky-400/20";
const exactActiveClasses = "bg-white/10 text-white ring-1 ring-sky-400/20";

const iconMap = {
  "mdi:view-dashboard-outline": viewDashboardOutline,
  "mdi:cog-outline": cogOutline,
  "mdi:account-multiple-outline": accountMultipleOutline,
  "mdi:shield-account-outline": shieldAccountOutline,
  "mdi:shield-key-outline": shieldKeyOutline,
  "mdi:file-document-edit-outline": fileDocumentEditOutline,
  "mdi:truck-outline": truckOutline,
  "mdi:truck-fast-outline": truckFastOutline,
  "mdi:truck-delivery-outline": truckDeliveryOutline,
  "mdi:truck-check-outline": truckCheckOutline,
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

function parseMenuLink(link) {
  const value = typeof link === "string" ? link.trim() : "";
  if (!value || value === "#") return null;

  try {
    const url = new URL(value, "http://govibe.local");
    const query = {};
    url.searchParams.forEach((v, k) => {
      query[k] = v;
    });
    return { path: url.pathname, query };
  } catch {
    return { path: value.split("?")[0], query: {} };
  }
}

function currentQueryValue(key) {
  const value = route.query[key];
  if (Array.isArray(value)) return value.length > 0 && value[0] != null ? String(value[0]) : "";
  return value == null ? "" : String(value);
}

function linkQueryMatches(query) {
  for (const [key, value] of Object.entries(query || {})) {
    if (currentQueryValue(key) !== value) return false;
  }
  return true;
}

function isLinkActive(link) {
  const parsed = parseMenuLink(link);
  if (!parsed) return false;

  const current = route.path;
  const isPathActive =
    parsed.path === "/admin"
      ? current === "/admin" || current === "/admin/"
      : current === parsed.path || current.startsWith(`${parsed.path}/`);

  return isPathActive && linkQueryMatches(parsed.query);
}

function hasActiveDescendant(item) {
  if (isLinkActive(item.link)) return true;
  if (!hasChildren(item)) return false;
  return item.children.some((c) => hasActiveDescendant(c));
}

function linkClasses(item) {
  const active = !item.link && hasActiveDescendant(item);
  return active ? activeClasses : "";
}

function childLinkClasses(item) {
  return isLinkActive(item.link) ? activeClasses : "";
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
  () => route.fullPath,
  () => {
    if (menuTree.value.length > 0) expandToActive(menuTree.value);
  }
);

onMounted(() => {
  loadMenus();
});
</script>
