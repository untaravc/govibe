<template>
  <div class="space-y-6" v-if="showPublicShell">
    <header class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight text-slate-900 dark:text-slate-50">{{ title }}</h1>
        <p class="mt-1 text-sm text-slate-600 dark:text-slate-300">Simple page routing with Vue Router.</p>
      </div>

      <nav class="flex items-center gap-2">
        <RouterLink
          to="/"
          class="rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-white/10"
        >
          Home
        </RouterLink>
        <RouterLink
          to="/auth/register"
          class="rounded-lg px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-white/10"
        >
          Register
        </RouterLink>
        <RouterLink
          to="/auth/login"
          class="rounded-lg bg-slate-900 px-3 py-2 text-sm font-medium text-white hover:bg-slate-800 dark:bg-slate-50 dark:text-slate-900 dark:hover:bg-white"
        >
          Login
        </RouterLink>
        <RouterLink
          to="/admin"
          class="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50 dark:border-white/10 dark:bg-slate-950 dark:text-slate-200 dark:hover:bg-white/5"
        >
          Admin
        </RouterLink>
      </nav>
    </header>

    <RouterView />
  </div>

  <RouterView v-else />
</template>

<script setup>
import { computed } from "vue";
import { RouterLink, RouterView, useRoute } from "vue-router";

const route = useRoute();

defineProps({
  title: {
    type: String,
    required: true
  }
});

const showPublicShell = computed(() => !route.path.startsWith("/admin"));
</script>
