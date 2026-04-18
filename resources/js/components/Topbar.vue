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
        <div ref="menuRoot" class="relative">
          <button
            type="button"
            class="inline-flex h-10 w-10 items-center justify-center overflow-hidden rounded-full border border-slate-200 bg-white shadow-sm ring-slate-900/10 hover:bg-slate-50 focus:outline-none focus:ring-4 dark:border-white/10 dark:bg-slate-950 dark:hover:bg-white/5"
            :aria-expanded="menuOpen ? 'true' : 'false'"
            aria-haspopup="menu"
            @click="toggleMenu"
          >
            <span class="sr-only">Open menu</span>
            <img v-if="avatarUrl" :src="avatarUrl" alt="" class="h-full w-full object-cover" />
            <span v-else class="text-sm font-semibold text-slate-700 dark:text-slate-200">{{ initials }}</span>
          </button>

          <div
            v-if="menuOpen"
            role="menu"
            class="absolute right-0 mt-2 w-56 overflow-hidden rounded-xl border border-slate-200 bg-white shadow-lg dark:border-white/10 dark:bg-slate-950"
          >
            <div class="px-4 py-3">
              <p class="text-sm font-semibold text-slate-900 dark:text-slate-50">{{ displayName }}</p>
              <p class="text-xs text-slate-500 dark:text-slate-400">{{ displayEmail }}</p>
            </div>
            <div class="h-px bg-slate-100 dark:bg-white/10"></div>
            <RouterLink
              to="/"
              role="menuitem"
              class="block px-4 py-2.5 text-sm text-slate-700 hover:bg-slate-50 dark:text-slate-200 dark:hover:bg-white/5"
              @click="closeMenu"
            >
              Public Site
            </RouterLink>
            <RouterLink
              :to="profilePath"
              role="menuitem"
              class="block px-4 py-2.5 text-sm text-slate-700 hover:bg-slate-50 dark:text-slate-200 dark:hover:bg-white/5"
              @click="closeMenu"
            >
              Profile
            </RouterLink>
            <button
              type="button"
              role="menuitem"
              class="block w-full px-4 py-2.5 text-left text-sm text-rose-600 hover:bg-rose-50 dark:text-rose-400 dark:hover:bg-rose-950/30"
              @click="logout"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { RouterLink, useRouter } from "vue-router";

import api from "../api.js";
import { useAuthStore } from "../store/auth.js";

defineEmits(["toggleSidebar", "toggleCollapse"]);

const router = useRouter();
const auth = useAuthStore();

const menuOpen = ref(false);
const menuRoot = ref(null);

const displayName = computed(() => String(auth.user?.name || "Account"));
const displayEmail = computed(() => String(auth.user?.email || ""));
const avatarUrl = computed(() => String(auth.user?.image || "").trim());
const profilePath = computed(() => {
  const id = auth.user?.id;
  if (id) return `/admin/users/${id}/edit`;
  return "/admin/settings";
});
const initials = computed(() => {
  const name = displayName.value.trim();
  if (!name) return "U";
  const parts = name.split(/\s+/).filter(Boolean);
  const first = parts[0]?.[0] || "U";
  const last = parts.length > 1 ? parts[parts.length - 1]?.[0] : "";
  return `${first}${last}`.toUpperCase();
});

function closeMenu() {
  menuOpen.value = false;
}

function toggleMenu() {
  menuOpen.value = !menuOpen.value;
}

function onWindowClick(e) {
  if (!menuOpen.value) return;
  const root = menuRoot.value;
  if (!root) return;
  if (root === e.target || root.contains(e.target)) return;
  closeMenu();
}

function onWindowKeydown(e) {
  if (!menuOpen.value) return;
  if (e.key === "Escape") closeMenu();
}

onMounted(() => {
  window.addEventListener("click", onWindowClick, true);
  window.addEventListener("keydown", onWindowKeydown);
});

onBeforeUnmount(() => {
  window.removeEventListener("click", onWindowClick, true);
  window.removeEventListener("keydown", onWindowKeydown);
});

async function logout() {
  closeMenu();
  try {
    await api.post("/api/logout", {}, { auth: true });
  } catch {
    // ignore
  }

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
