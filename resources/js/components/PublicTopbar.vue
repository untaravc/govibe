<template>
  <header class="sticky top-0 z-20 border-b border-slate-200/60 bg-white/80 backdrop-blur">
    <div class="flex w-full items-center justify-between gap-3 px-4 py-3">
      <div class="min-w-0">
        <p class="text-sm font-semibold text-slate-900">GoVibe</p>
        <p class="text-xs text-slate-500">Public site</p>
      </div>

      <nav class="flex items-center gap-2">
        <RouterLink
          to="/"
          class="rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100"
          active-class="bg-slate-100 text-slate-900"
          exact-active-class="bg-slate-100 text-slate-900"
        >
          Home
        </RouterLink>
        <RouterLink
          to="/auth/register"
          class="rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100"
          active-class="bg-slate-100 text-slate-900"
        >
          Register
        </RouterLink>
        <RouterLink
          to="/auth/login"
          class="rounded-lg bg-slate-900 px-3 py-2 text-sm font-medium text-white hover:bg-slate-800"
        >
          Login
        </RouterLink>
        <RouterLink
          v-if="hasToken"
          to="/admin"
          class="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
        >
          Admin site
        </RouterLink>
      </nav>
    </div>
  </header>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";

const token = ref("");

const hasToken = computed(() => token.value.trim().length > 0);

function syncTokenFromStorage() {
  try {
    token.value = localStorage.getItem("token") || "";
  } catch {
    token.value = "";
  }
}

function onStorage(event) {
  if (event.key === "token") syncTokenFromStorage();
}

onMounted(() => {
  syncTokenFromStorage();
  window.addEventListener("storage", onStorage);
});

onBeforeUnmount(() => {
  window.removeEventListener("storage", onStorage);
});
</script>
