<template>
  <header
    class="sticky top-0 z-20 border-b border-blue-900/40 bg-blue-950 text-white shadow-sm backdrop-blur"
  >
    <div class="flex w-full items-center justify-between gap-3 px-4 py-3">
      <div class="flex items-center gap-3">
        <button
          type="button"
          class="inline-flex items-center justify-center rounded-lg p-2 text-slate-200 hover:bg-white/10 hover:text-white md:hidden"
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
          class="hidden items-center justify-center rounded-lg p-2 text-slate-200 hover:bg-white/10 hover:text-white md:inline-flex"
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
          <p class="text-sm font-semibold tracking-tight text-white">Administrator</p>
          <p class="text-xs text-slate-300/80">GoVibe dashboard</p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <div ref="menuRoot" class="relative">
          <button
            type="button"
            class="inline-flex h-10 w-10 items-center justify-center overflow-hidden rounded-full border border-white/10 bg-white/5 shadow-sm ring-blue-500/20 hover:bg-white/10 focus:outline-none focus:ring-4"
            :aria-expanded="menuOpen ? 'true' : 'false'"
            aria-haspopup="menu"
            @click="toggleMenu"
          >
            <span class="sr-only">Open menu</span>
            <img v-if="avatarUrl" :src="avatarUrl" alt="" class="h-full w-full object-cover" />
            <span v-else class="text-sm font-semibold text-white">{{ initials }}</span>
          </button>

          <div
            v-if="menuOpen"
            role="menu"
            class="absolute right-0 mt-2 w-56 overflow-hidden rounded-xl border border-blue-900/40 bg-slate-950 shadow-lg"
          >
            <div class="px-4 py-3">
              <p class="text-sm font-semibold text-white">{{ displayName }}</p>
              <p class="text-xs text-slate-300/80">{{ displayEmail }}</p>
            </div>
            <div class="h-px bg-white/10"></div>
            <RouterLink
              to="/"
              role="menuitem"
              class="block px-4 py-2.5 text-sm text-slate-200 hover:bg-white/10 hover:text-white"
              @click="closeMenu"
            >
              Public Site
            </RouterLink>
            <RouterLink
              :to="profilePath"
              role="menuitem"
              class="block px-4 py-2.5 text-sm text-slate-200 hover:bg-white/10 hover:text-white"
              @click="closeMenu"
            >
              Profile
            </RouterLink>
            <button
              type="button"
              role="menuitem"
              class="block w-full px-4 py-2.5 text-left text-sm text-rose-300 hover:bg-rose-500/10 hover:text-rose-200"
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
const profilePath = computed(() => "/admin/profile");
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
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    sessionStorage.removeItem("access_token");
    sessionStorage.removeItem("refresh_token");
  } catch {
    // ignore
  }
  router.push("/auth/login");
}
</script>
